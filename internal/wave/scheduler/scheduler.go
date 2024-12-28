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
	"cthul.io/cthul/pkg/wave/domain"
	"cthul.io/cthul/pkg/wave/node"
)

// Scheduler provides a component responsible for advertising the local node and its resources to the cluster.
// If the scheduler is set as the leader scheduler, it indexes the advertised nodes and moves unmanaged domains
// (domains located on nodes that are NOT advertised) to advertised nodes based on available resources.
type Scheduler struct {
	// root context runs until the scheduler is fully terminated.
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc

	// work context runs as long as the scheduler is operating.
	// cancel it to initiate scheduler termination.
	workCtx       context.Context
	workCtxCancel context.CancelFunc

	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the controller is fully cleaned up.
	finChan chan struct{}

	client db.Client
	logger log.Logger

	// leaderStateChan is used to emit the state of the leader.
	leaderStateChan chan bool

	// cycleTTL specifies the interval for scheduler cycles.
	cycleTTL int64
	// rescheduleCycles specifies the number of cycles that a domain must be unmanaged
	// in a row until it is rescheduled.
	rescheduleCycles int64

	// register paramters specify how the local node registers itself on the scheduler
	// this happens independent of the leader state.
	registerId string
	registerTTL int64
	registerCpuFactor float64
	registerMemFactor float64
}

type SchedulerOption func(*Scheduler)

// NewScheduler creates a new scheduler instance.
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
		leaderStateChan: make(chan bool),
		cycleTTL: 5,
		rescheduleCycles: 2,
		registerId: "",
		registerTTL: 5,
		registerCpuFactor: 1,
		registerMemFactor: 1,
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
			s.registerId = nodeId
		} else {
			s.registerId = ""
		}
	}
}

// WithRegisterTTL sets a custom ttl for the register cycle (essentially the keepalive interval).
func WithRegisterTTL(registerTTL int64) SchedulerOption {
	return func(s *Scheduler) {
		s.registerTTL = registerTTL
	}
}

// WithLocalResourceThreshold defines a custom resource threshold for this node.
// The resource threshold must be provided as percentage and essentially serves as filter applied to the raw
// host resources which are reported to the scheduler for deciding how to schedule domains.
// E.g. mem/cpuThreshold of 80 with 10 cores and 10GB memory converts to reported 8 cores and 8GB memory.
func WithLocalResourceThreshold(cpuThreshold, memThreshold int64) SchedulerOption {
	return func(s *Scheduler) {
		s.registerCpuFactor = float64(cpuThreshold) / 100
		s.registerMemFactor = float64(memThreshold) / 100
	}
}

// WithLogger sets a custom logger for the scheduler.
func WithLogger(logger log.Logger) SchedulerOption {
	return func(s *Scheduler) {
		s.logger = logger
	}
}


// SetLeaderState changes the state of the schedule leader. Local set to true enables the scheduler leader mode.
// Setting it to false disables the scheduler leader mode gracefully. This operation is idempotent.
// Currently the leader id (first param) is unused, it is in place for future expansion.
func (s *Scheduler) SetLeaderState(_ string, local bool) {
	select {
	case <-s.workCtx.Done():
	case s.leaderStateChan <- local:
	}
}

// ServeAndDetach starts the scheduler registration and the leader scheduler in a detached goroutine.
func (s *Scheduler) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.registerNode()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-s.workCtx.Done():
				return
			case leaderState := <-s.leaderStateChan:
				if leaderState {
					s.runSchedulerCycle()
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		s.finChan <- struct{}{}
	}()
}

// runSchedulerCycle starts a scheduler leader cycle and blocks until the leaderStateChan emits false.
func (s *Scheduler) runSchedulerCycle() {
	cycleCtx, cycleCtxCancel := context.WithCancel(s.workCtx)
	defer cycleCtxCancel()
	go s.startSchedulerCycle(cycleCtx)
	for {
		select {
		case <-cycleCtx.Done():
			return
		case leaderState := <-s.leaderStateChan:
			if !leaderState {
				return
			}
		}
	}
}

// Terminate shuts down the scheduler gracefully, if shutdown did not complete in the provided context window
// the scheduler is terminated forcefully. Never returns an error (just there to match termination pattern).
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
