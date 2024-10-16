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
	"sync"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)

type Scheduler struct {
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc

	workCtx       context.Context
	workCtxCancel context.CancelFunc

	finChan chan struct{}

	client db.Client
	logger log.Logger

	localNode node

	
	leaderState     bool
	leaderStateLock sync.RWMutex

	cycleTTL int64
}

type node struct {
	id          string
	registerTTL int64
	cpuFactor   float64
	memFactor   float64
}

type SchedulerOption func(*Scheduler)

func NewScheduler(client db.Client, opts ...SchedulerOption) *Scheduler {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	scheduler := &Scheduler{
		rootCtx:         rootCtx,
		rootCtxCancel:   rootCtxCancel,
		workCtx:         workCtx,
		workCtxCancel:   workCtxCancel,
		finChan:         make(chan struct{}),
		client:          client,
		logger:          discard.NewDiscardLogger(),
		localNode:       node{id: "", registerTTL: 5, cpuFactor: 1, memFactor: 1},
		leaderState:     false,
		leaderStateLock: sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(scheduler)
	}

	return scheduler
}

// WithLocalNode registers this local node in the scheduler, allowing it to allocate domains to this node.
func WithLocalNode(register bool, nodeId string) SchedulerOption {
	return func(s *Scheduler) {
		if register {
			s.localNode.id = nodeId
		} else {
			s.localNode.id = ""
		}
	}
}

// WithRegisterTTL sets a custom ttl for the register cycle (essentially the keepalive interval).
func WithRegisterTTL(registerTTL int64) SchedulerOption {
	return func(s *Scheduler) {
		s.localNode.registerTTL = registerTTL
	}
}

// WithLocalResourceThreshold defines a custom resource threshold for this node.
// The resource threshold must be provided as percentage and essentially serves as filter applied to the raw
// host resources which are reported to the scheduler for deciding how to schedule domains.
// E.g. mem/cpuThreshold of 80 with 10 cores and 10GB memory converts to reported 8 cores and 8GB memory.
func WithLocalResourceThreshold(cpuThreshold, memThreshold int64) SchedulerOption {
	return func(s *Scheduler) {
		s.localNode.cpuFactor = float64(cpuThreshold) / 100
		s.localNode.memFactor = float64(memThreshold) / 100
	}
}

// WithLogger sets a custom logger for the scheduler.
func WithLogger(logger log.Logger) SchedulerOption {
	return func(s *Scheduler) {
		s.logger = logger
	}
}

func (s *Scheduler) SetState(leader bool) {
	s.leaderStateLock.Lock()
	defer s.leaderStateLock.Unlock()
	s.leaderState = leader
}

func (s *Scheduler) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.registerNode()
	}()

	go func() {
		wg.Wait()
		s.finChan <- struct{}{}
	}()
}

func (s *Scheduler) Nodes() {
	ctx, cancel := context.WithCancel(s.workCtx)
	defer cancel()

	err := s.client.WatchRange(ctx, "/WAVE/SCHEDULER/NODE/", func(key, value string, err error) {

	})
	if err != nil {
		s.logger.Crit("scheduler", "unrecoverable watch error occured: "+err.Error())
	}
}

func (s *Scheduler) Terminate(ctx context.Context) error {
	s.workCtxCancel()
	defer s.rootCtxCancel()
	select {
	case <-s.finChan:
		return nil
	case <-ctx.Done():
		s.rootCtxCancel()
		<-s.finChan
		return nil
	}
}
