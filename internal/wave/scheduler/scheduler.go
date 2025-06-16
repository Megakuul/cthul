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
  "log/slog"

	"cthul.io/cthul/pkg/db"
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
	// if the channel emits, this indicates that the operator is fully cleaned up.
	finChan chan struct{}

	client db.Client
  logger *slog.Logger
	domainController *domain.Controller
	nodeController *node.Controller

	// leaderStateChan is used to emit the state of the leader.
	leaderStateChan chan bool

	// cycleTTL specifies the interval for scheduler cycles.
	cycleTTL int64
	// rescheduleCycles specifies the number of cycles that a domain must be unmanaged
	// in a row until it is rescheduled.
	rescheduleCycles int64
}

type Option func(*Scheduler)

// New creates a new scheduler instance.
func New(logger *slog.Logger, client db.Client, domain *domain.Controller, node *node.Controller, opts ...Option) *Scheduler {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	scheduler := &Scheduler{
		rootCtx:         rootCtx,
		rootCtxCancel:   rootCtxCancel,
		workCtx:         workCtx,
		workCtxCancel:   workCtxCancel,
		finChan:         make(chan struct{}),
		client:          client,
    logger: logger.WithGroup("scheduler"),
		domainController: domain,
		nodeController: node,
		leaderStateChan: make(chan bool),
		cycleTTL: 5,
		rescheduleCycles: 2,
	}

	for _, opt := range opts {
		opt(scheduler)
	}

	return scheduler
}


// WithCycleTTL defines a custom scheduler cycle interval.
func WithCycleTTL(ttl int64) Option {
	return func(s *Scheduler) {
		s.cycleTTL = ttl
	}
}

// WithRescheduleCycles sets a custom number of scheduler cycles a domain must be unmanaged
// before being rescheduled.
func WithRescheduleCycles(cycles int64) Option {
	return func(s *Scheduler) {
		s.rescheduleCycles = cycles
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
