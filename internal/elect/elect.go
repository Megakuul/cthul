/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package elect

import (
	"context"
	"fmt"
	"sync"
	"time"
  "log/slog"

	"cthul.io/cthul/pkg/db"
)

// Operator is used in every component to participate in the leader election system.
// It provides functions to update the leader status and also contest the leader to given conditions.
//
// Election is managed over a contestKey on the database which is usually something like WAVE/LEADER.
// The election (internally called `mr. crabs election system`) works like this:
// If the leader value is empty or in an invalid format the local node is contesting the leader.
// If the current leader == local node, the leader is NOT updated (important to avoid a race condition).
// If the local node has more cash then the leader node, the local node is contesting the leader.
// If the local node has less cash then the leader node, the leader node is set as leader.
type Operator struct {
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
  logger *slog.Logger

	// contestKey marks the database key which is used to contest the leader.
	contestKey string
	// contestTTL specifies the time to live of one contest cycle.
	contestTTL int64
	// contestHook holds a callback executed every contest cycle with the contest state.
	contestHook func(string, bool)

	// localNode holds leader information about the local node.
	localNode clusterLeader
	
	// leaderNode holds the leader information of the current leader.
	leaderNode clusterLeader
	leaderNodeLock sync.RWMutex
}

type Option func(*Operator)

// New creates a new election controller. The contestKey is the database key used for election.
func New(logger *slog.Logger, client db.Client, contestKey string, opts ...Option) *Operator {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	controller := &Operator{
		rootCtx: rootCtx,
		rootCtxCancel: rootCtxCancel,
		workCtx: workCtx,
		workCtxCancel: workCtxCancel,
		finChan: make(chan struct{}),
		client: client,
		logger: logger.WithGroup("elect-controller"),
		contestKey: contestKey,
		contestTTL: 5,
		contestHook: func(_ string, _ bool) {},
		localNode: clusterLeader{ Id: "", Cash: -1 },
		leaderNode: clusterLeader{ Id: "", Cash: -1 },
		leaderNodeLock: sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// WithLocalLeader enables the local node to contest the leader. The specified nodeId will be used as
// leader id if this node contests. NodeCash determines the importance of this node; more cash = more important.
func WithLocalLeader(contest bool, nodeId string, nodeCash int64) Option {
	return func (o *Operator) {
		if contest {
			o.localNode.Id, o.localNode.Cash = nodeId, nodeCash
		} else {
			o.localNode.Id, o.localNode.Cash = "", -1
		}
	}
}

// WithContestTTL specifies a custom ttl for the leader contest cycle. If the leader is contested by this node
// it does this in cycles based on this ttl.
func WithContestTTL(ttl int64) Option {
	return func (o *Operator) {
		o.contestTTL = ttl
	}
}

// WithContestHooks adds one or more callback hooks that are executed on every contest cycle.
// Callback informs about the state of the contestant on every contest cylce. It provides the id of the
// current contestant and a bool stating whether this is the local node or not.
// The callback functions must NOT block, they block the contest cycle. Callbacks should be idempotent.
func WithContestHooks(callbacks ...func(string, bool)) Option {
	return func (o *Operator) {
		o.contestHook = func(contestantId string, local bool) {
			for _, callback := range callbacks {
				callback(contestantId, local)
			}
		}
	}
}

// ServeAndDetach launches two routines to check the current leader and contest it under given conditions.
func (o *Operator) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	
	go func () {
		defer wg.Done()
		o.checkLeader()
	}()
	go func() {
		defer wg.Done()
		o.contestLeader()
	}()
	go func() {
		wg.Wait()
		o.finChan<-struct{}{}
	}()
}

// checkLeader watches the contestKey for changes. On every change, the controller performs a local election
// that determines wheter the local node should contest the leader or not.
func (o *Operator) checkLeader() {
	ctx, cancel := context.WithCancel(o.workCtx)
	defer cancel()
	
	leaderStr, err := o.client.Get(ctx, o.contestKey)
	if err!=nil {
		o.logger.Error(fmt.Sprintf("initial check failed: %v", err))
	}
	leaderNode := o.electLeader(leaderStr)
	if leaderNode!=nil {
		o.leaderNodeLock.Lock()
		o.leaderNode = *leaderNode
		o.leaderNodeLock.Unlock()
	}
	
	err = o.client.Watch(ctx, o.contestKey, func(_, value string, err error) {
		if err!=nil {
			o.logger.Error(err.Error())
			return
		}
		leaderNode := o.electLeader(value)
		if leaderNode!=nil {
			o.leaderNodeLock.Lock()
			o.leaderNode = *leaderNode
			o.leaderNodeLock.Unlock()
		}
	})
	if err!=nil {
		o.logger.Error("unrecoverable check error occured: " + err.Error())
	}
}

// electLeader analyzes the leaderStr and returns the new leader or nil if it should not be changed.
func (o *Operator) electLeader(leaderStr string) *clusterLeader {
	if leaderStr=="" {
		o.logger.Debug("contesting leader; reason: leader is uncontested")
		return &clusterLeader{ Id: o.localNode.Id, Cash: o.localNode.Cash }
	}
	newLeaderNode, err := parseClusterLeader(leaderStr)
	if err!=nil {
		o.logger.Debug("contesting leader; reason: " + err.Error())
		return &clusterLeader{ Id: o.localNode.Id, Cash: o.localNode.Cash }
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
	if o.localNode.Id == newLeaderNode.Id {
		o.logger.Debug("skipping leader; reason: local node is already leader")
		return nil
	}
	if o.localNode.Cash <= newLeaderNode.Cash {
		o.logger.Debug("skipping leader; reason: local node has not enough cash")
		return &clusterLeader{ Id: newLeaderNode.Id, Cash: newLeaderNode.Cash }
	}
	
	o.logger.Debug("contesting leader; reason: local node has more cash")
	return &clusterLeader{ Id: o.localNode.Id, Cash: o.localNode.Cash }
}


// contestLeader checks if the current leader (reported by checkLeader) matches the local node.
// If the local node is the reported leader, it will set the local node as leader
// and repeat this step in the provided contestTTL interval.
func (o *Operator) contestLeader() {
	if o.localNode.Id == "" {
		o.logger.Info("local node will not serve as leader")
		return
	}
	
	for {
		o.leaderNodeLock.RLock()
		if o.leaderNode.Id != "" && o.leaderNode.Id == o.localNode.Id {
			ctx, cancel := context.WithTimeout(o.workCtx, time.Second * time.Duration(o.contestTTL))
			defer cancel()
			_, err := o.client.Set(ctx, o.contestKey,
				serializeClusterLeader(&o.leaderNode), (o.contestTTL * 2),
			)
			if err!=nil {
				o.logger.Error("failed to contest leader")
			}
			o.contestHook(o.leaderNode.Id, true)
		} else {
			o.contestHook(o.leaderNode.Id, false)
		}
		o.leaderNodeLock.RUnlock()
		
		select {
		case <-time.After(time.Second * time.Duration(o.contestTTL)):
			break
		case <-o.workCtx.Done():
			o.leaderNodeLock.RLock()
			defer o.leaderNodeLock.RUnlock()
			if o.leaderNode.Id != "" && o.leaderNode.Id == o.localNode.Id {
				// If the node is currently contesting leader, it sets leader explicitly to "" before termination
				// so that other nodes can immediately contest the leader.
				err := o.client.Delete(o.rootCtx, o.contestKey)
				if err!=nil {
					o.logger.Error("failed to reset leader before termination")
				}
			}
			return
		}
	}
}


// Terminate stops the election controller gracefully. If this node currently contested the leader
// it tries to reset the contestKey in order to make other nodes immediately contest the leader.
// If this does not succeed in the provided context window, it terminates forcefully.
func (o *Operator) Terminate(ctx context.Context) error {
	o.workCtxCancel()
	defer o.rootCtxCancel()
	select {
	case <-o.finChan:
		return nil
	case <-ctx.Done():
		o.rootCtxCancel()
		<-o.finChan
		return nil
	}
}
