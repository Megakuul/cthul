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
	"sync/atomic"
	"time"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)


type ElectController struct {
	rootCtx context.Context
	rootCtxCancel context.CancelFunc
	client db.Client
	logger log.Logger
	
	leaderKey string
	localNodeId string
	localNodeCash int64
	contestTTL int64
	contestLeader atomic.Bool
}

type ElectControllerOption func(*ElectController)

func NewElectController(client db.Client, leaderKey string, opts ...ElectControllerOption) *ElectController {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	controller := &ElectController{
		rootCtx: rootCtx,
		rootCtxCancel: rootCtxCancel,
		client: client,
		logger: discard.NewDiscardLogger(),
		leaderKey: leaderKey,
		localNodeId: "",
		localNodeCash: -1,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

func WithContestor(nodeId string, nodeCash int64) ElectControllerOption {
	return func (e *ElectController) {
		e.localNodeId = nodeId
		e.localNodeCash = nodeCash
	}
}

func WithLogger(logger log.Logger) ElectControllerOption {
	return func (e *ElectController) {
		e.logger = logger
	}
}

func (e *ElectController) CheckLeader() {
	ctx, cancel := context.WithCancel(e.rootCtx)
	go func() {
		defer cancel()

		leaderStr, err := e.client.Get(ctx, e.leaderKey)
		if err!=nil {
			e.logger.Err("elect_controller", fmt.Sprintf("initial check failed: %v", err))
		}
		if e.checkCurrentLeader(leaderStr) {
			e.contestLeader.Store(true)
		}
		
		err = e.client.Watch(ctx, e.leaderKey, func(_, value string, err error) {
			e.logger.Debug("elect_controller", "received new leader update: " + value)
			if err!=nil {
				e.logger.Err("elect_controller", err.Error())
				return
			}
			if e.checkCurrentLeader(value) {
				e.contestLeader.Store(true)
			} else {
				e.contestLeader.Store(false)
			}
		})
		if err!=nil {
			e.logger.Crit("elect_controller", fmt.Sprintf("unrecoverable check error occured: %v", err))
		}
	}()
}

// checkCurrentLeader parses the leader value string and determines if the leader is valid
// and if yes, if he has more cash then the local node. If this function concludes that the local node
// should contest the leader, it returns true.
func (e *ElectController) checkCurrentLeader(leaderStr string) bool {
	if leaderStr == "" {
		return true
	}
	
	leaderSub := strings.SplitN(leaderStr, "|", 2)
	if len(leaderSub) != 2 {
		return true
	}
	
	leaderId := leaderSub[0]
	leaderCash, err := strconv.Atoi(leaderSub[1])
	if err!=nil {
		return true
	}

	if e.localNodeId == leaderId {
		return true
	}
	if e.localNodeCash > int64(leaderCash) {
		return true
	}
	return false
}

func (e *ElectController) ContestLeader() {
	go func() {
		for {
			if e.contestLeader.Load() {
				ctx, cancel := context.WithTimeout(e.rootCtx, time.Second * time.Duration(e.contestTTL))
				defer cancel()
				err := e.client.Set(ctx, e.leaderKey,
					fmt.Sprintf("%s|%d", e.localNodeId, e.localNodeCash), (e.contestTTL * 2),
				)
				if err!=nil {
					e.logger.Err("elect_controller", err.Error())
				}
			}
			select {
			case <-time.After(time.Second * time.Duration(e.contestTTL)):
				break
			case <-e.rootCtx.Done():
				return
			}
		}
	}()
}
