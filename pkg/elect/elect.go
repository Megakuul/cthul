/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package elect

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)


// ElectController is used in every component to participate in the leader election system.
// It provides functions to update the leader status and also contest the leader to given conditions.
//
// Election is managed over a contestKey on the database which is usually something like WAVE/LEADER.
// The election (internally called `mr. crabs election system`) works like this:
// If the leader value is empty or invalid (not in 'nodeid|cash' format) the local node is contesting the leader.
// If the current leader == local node, the leader is NOT updated (important to avoid a race condition).
// If the local node has more cash then the leader node, the local node is contesting the leader.
// If the local node has less cash then the leader node, the leader node is set as leader.
type ElectController struct {
	// rootCtx is active for the full lifetime of the controller.
	// closing it leads to the forceful / immediate shutdown of the controller.
	rootCtx context.Context
	rootCtxCancel context.CancelFunc

	// workCtx is active for the lifetime of the background operations
	// closing it leads to the graceful shutdown of the controller.
	workCtx context.Context
	workCtxCancel context.CancelFunc
	
	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the controller is fully cleaned up.
	finChan chan struct{}

	
	client db.Client
	logger log.Logger

	// contestKey marks the database key which is used to contest the leader.
	contestKey string
	// contentTTL specifies the ttl of one contest cycle.
	// After this the leader re-contested (if elected).
	contestTTL int64

	// localNode holds the leader information of this node.
	localNode node
	// leaderNode holds the leader information of the current leader.
	leaderNode node
	leaderNodeLock sync.RWMutex
}

// node provides information about a node relevant to contest the leader.
type node struct {
	active bool
	id string
	cash int64
}

type ElectControllerOption func(*ElectController)

// NewElectController creates a new election controller. The contestKey is the database key used for election.
func NewElectController(client db.Client, contestKey string, opts ...ElectControllerOption) *ElectController {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	controller := &ElectController{
		rootCtx: rootCtx,
		rootCtxCancel: rootCtxCancel,
		workCtx: workCtx,
		workCtxCancel: workCtxCancel,
		finChan: make(chan struct{}),
		client: client,
		logger: discard.NewDiscardLogger(),
		contestKey: contestKey,
		contestTTL: 5,
		localNode: node{ active: false, id: "", cash: -1 },
		leaderNode: node{ active: false, id: "", cash: -1 },
		leaderNodeLock: sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// WithLocalLeader enables the local node to contest the leader. The specified nodeId will be used as
// leader id if this node contests. NodeCash determines the importance of this node; more cash = more important.
func WithLocalLeader(nodeId string, nodeCash int64) ElectControllerOption {
	return func (e *ElectController) {
		e.localNode = node{
			active: true,
			id: nodeId,
			cash: nodeCash,
		}
	}
}

// WithContestTTL specifies a custom ttl for the leader contest cycle. If the leader is contested by this node
// it does this in cycles based on this ttl.
func WithContestTTL(ttl int64) ElectControllerOption {
	return func (e *ElectController) {
		e.contestTTL = ttl
	}
}

// WithLogger sets a custom logger to the election controller.
func WithLogger(logger log.Logger) ElectControllerOption {
	return func (e *ElectController) {
		e.logger = logger
	}
}

// GetLeader returns the current leader id. If the leader is currently not determined it returns "".
func (e *ElectController) GetLeader() string {
	e.leaderNodeLock.RLock()
	defer e.leaderNodeLock.RUnlock()
	if e.leaderNode.active {
		return e.leaderNode.id
	}
	return ""
}

// ServeAndDetach launches two routines to check the current leader and contest it under given conditions.
func (e *ElectController) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	
	go func () {
		defer wg.Done()
		e.checkLeader()
	}()
	go func() {
		defer wg.Done()
		e.contestLeader()
	}()
	go func() {
		wg.Wait()
		e.finChan<-struct{}{}
	}()
}

// checkLeader watches the contestKey for changes. On every change, the controller performs a local election
// that determines wheter the local node should contest the leader or not.
func (e *ElectController) checkLeader() {
	ctx, cancel := context.WithCancel(e.workCtx)
	defer cancel()
	
	leaderStr, err := e.client.Get(ctx, e.contestKey)
	if err!=nil {
		e.logger.Err("elect_controller", fmt.Sprintf("initial check failed: %v", err))
	}
	leaderNode := e.electLeader(leaderStr)
	if leaderNode!=nil {
		e.leaderNodeLock.Lock()
		e.leaderNode = *leaderNode
		e.leaderNodeLock.Unlock()
	}
	
	err = e.client.Watch(ctx, e.contestKey, func(_, value string, err error) {
		e.logger.Debug("elect_controller", "received new leader update: " + value)
		if err!=nil {
			e.logger.Err("elect_controller", err.Error())
			return
		}
		leaderNode := e.electLeader(value)
		if leaderNode!=nil {
			e.leaderNodeLock.Lock()
			e.leaderNode = *leaderNode
			e.leaderNodeLock.Unlock()
		}
	})
	if err!=nil {
		e.logger.Crit("elect_controller", "unrecoverable check error occured: " + err.Error())
	}
}

// electLeader analyzes the leaderStr and returns the new leader or nil if it should not be changed.
func (e *ElectController) electLeader(leaderStr string) *node {
	if leaderStr=="" {
		e.logger.Debug("elect_controller", "contesting leader; reason: leader is uncontested")
		return &node{ active: e.localNode.active, id: e.localNode.id, cash: e.localNode.cash }
	}
	newLeaderNode, err := e.parseLeaderString(leaderStr)
	if err!=nil {
		e.logger.Debug("elect_controller", "contesting leader; reason: " + err.Error())
		return &node{ active: e.localNode.active, id: e.localNode.id, cash: e.localNode.cash }
	}

	// Important: If the local node == new leader the leader node should NOT be changed.
	// This is very important because otherwise this causes a rare race condition that works like this:
	// Node1 with 10$ contests for leader in a 5 second cycle, it believes it is the leader.
	// Node2 with 20$ has the exact same schedule and knows that in the next cycle he must contest.
	// Node2 elects itself as the leader, it believes to have more cash then Node1.
	// Node1 elects itself also as the leader, it doesn't know that Node2 contests yet.
	// Node2 sets itself as leader.
	// Node1 sets itself as leader (even if Node2 has more cash).
	// Node1 & Node2 elect Node2 as the leader, Node2 has more cash then Node1.
	// Node1 & Node2 elect Node1 as the leader, this is only because the Node1 leader write occured after
	// the Node2 write. This write race condition can lead to the situation where the watch update of Node2's
	// election commes in later.
	// Because usually both controllers have the same schedule this will repeat itself.
	// To avoid this, elections that have the local node as candidate are skipped instead of overwritting the
	// current leader. With this, the leader is only overwritten if he actually has more cash.
	if e.localNode.id == newLeaderNode.id {
		e.logger.Debug("elect_controller", "skipping leader; reason: local node is already leader")
		return nil
	}
	if e.localNode.cash <= newLeaderNode.cash {
		e.logger.Debug("elect_controller", "skipping leader; reason: local node has not enough cash")
		return &node{ active: true, id: newLeaderNode.id, cash: newLeaderNode.cash }
	}
	
	e.logger.Debug("elect_controller", "contesting leader; reason: local node has more cash")
	return &node{ active: e.localNode.active, id: e.localNode.id, cash: e.localNode.cash }
}


// contestLeader checks if the current leader (reported by checkLeader) matches the local node.
// If the local node is active and the reported leader, it will set the local node as leader
// and repeat this step in the provided contestTTL interval.
func (e *ElectController) contestLeader() {
	for {
		e.leaderNodeLock.RLock()
		if e.leaderNode.active && e.leaderNode.id == e.localNode.id {
			ctx, cancel := context.WithTimeout(e.workCtx, time.Second * time.Duration(e.contestTTL))
			defer cancel()
			err := e.client.Set(ctx, e.contestKey,
				fmt.Sprintf("%s|%d", e.leaderNode.id, e.leaderNode.cash), (e.contestTTL * 2),
			)
			if err!=nil {
				e.logger.Err("elect_controller", err.Error())
			}
		}
		e.leaderNodeLock.RUnlock()
		select {
		case <-time.After(time.Second * time.Duration(e.contestTTL)):
			break
		case <-e.workCtx.Done():
			if e.leaderNode.active && e.leaderNode.id == e.localNode.id {
				// If the node is currently contesting leader, it sets leader explicitly to "" before termination
				// so that other nodes can immediately contest the leader.
				err := e.client.Set(e.rootCtx, e.contestKey, "", 0)
				if err!=nil {
					e.logger.Err("elect_controller", "failed to reset leader before termination")
				}
			}
			return
		}
	}
}

// parseLeaderString parses the raw leader string into a node struct.
func (e *ElectController) parseLeaderString(leaderStr string) (*node, error) {
	leaderSub := strings.SplitN(leaderStr, "|", 2)
	if len(leaderSub) != 2 {
		return nil, fmt.Errorf("invalid leader string; expected '|'.")
	}

	id := leaderSub[0]
	cash, err := strconv.Atoi(leaderSub[1])
	if err!=nil {
		return nil, fmt.Errorf("invalid leader string; cannot parse cash.")
	}	
	
	return &node {
		id: id,
		cash: int64(cash),
	}, nil
}

// serializeLeaderString serializes the node struct into a raw leader string.
func (e *ElectController) serializeLeaderString(leaderNode *node) string {
	if leaderNode!=nil {
		return fmt.Sprintf("%s|%d", leaderNode.id, leaderNode.cash)
	}
	return ""
}

// Terminate stops the election controller gracefully. If this node currently contested the leader
// it tries to reset the contestKey in order to make other nodes immediately contest the leader.
// If this does not succeed in the provided context window, it terminates forcefully.
func (e *ElectController) Terminate(ctx context.Context) error {
	e.workCtxCancel()
	defer e.rootCtxCancel()
	select {
	case <-e.finChan:
		return nil
	case <-ctx.Done():
		e.rootCtxCancel()
		<-e.finChan
		return nil
	}
}
