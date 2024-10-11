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

package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"cthul.io/cthul/pkg/db"
)

type Scheduler struct {
	rootCtx context.Context
	rootCtxCancel context.CancelFunc

	workCtx context.Context
	workCtxCancel context.CancelFunc
	
	client db.Client

	registerTTL int64
	
	localNode node
	localNodeLock sync.RWMutex

	leaderState bool
	leaderStateLock sync.RWMutex
}

type node struct {
	Id string `json:"id"`
	Capacity capacity `json:"capacity"`
}

type capacity struct {
	CPUs int64 `json:"cpus"`
	Memory int64 `json:"memory"`
	Storage int64 `json:"storage"`
}

type SchedulerOption func(*Scheduler)

func NewScheduler(client db.Client, opts ...SchedulerOption) *Scheduler {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	scheduler := &Scheduler{
		rootCtx: rootCtx,
		rootCtxCancel: rootCtxCancel,
		workCtx: workCtx,
		workCtxCancel: workCtxCancel,
		client: client,
		leaderState: false,
		leaderStateLock: sync.RWMutex{},
		registerTTL: 5,
		localNode: node{ active: false, id: "", capacity: -1 },
	}

	for _, opt := range opts {
		opt(scheduler)
	}

	return scheduler
}

func (s *Scheduler) SetState(leader bool) {
	s.leaderStateLock.Lock()
	defer s.leaderStateLock.Unlock()
	s.leaderState = leader
}


func (s *Scheduler) registerNode() {	
	for {
		if s.localNode.active && s.localNode.id != "" {
			ctx, cancel := context.WithTimeout(s.workCtx, time.Second * time.Duration(s.registerTTL))
			err := s.client.Set(ctx,
				fmt.Sprintf("/WAVE/SCHEDULER/NODE/%s", s.localNode.Id),
				s.serializeNode(&s.localNode),
				(s.registerTTL * 2),
			)		
		}
	}
}


// serializeNode serializes the node into a raw string.
func (s *Scheduler) parseNode(nodeStr string) (*node, error) {
	var node node
	err := json.Unmarshal([]byte(nodeStr), node)
	if err!=nil {
		return nil, fmt.Errorf("cannot parse node information")
	}
	return &node, nil
}

// serializeNode serializes the node into a raw string.
func (s *Scheduler) serializeNode(node *node) string {
	nodeStr, err := json.Marshal(node)
	if err!=nil {
		return ""
	}
	return string(nodeStr)
}
