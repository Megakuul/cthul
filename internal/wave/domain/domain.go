/**
 * Cthul System
 *
 * Copyright (C) 2025 Linus Ilian Moser <linus.moser@megakuul.ch>
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

package domain

import (
	"context"
	"sync"

	domadapter "cthul.io/cthul/pkg/adapter/domain"
	domctrl "cthul.io/cthul/pkg/wave/domain"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)

// Operator is responsible for applying the domains database state to the local virtual machine monitor.
type Operator struct {
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc

	workCtx       context.Context
	workCtxCancel context.CancelFunc

	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the operator is fully cleaned up.
	finChan chan struct{}

	controller domctrl.Controller
	adapter domadapter.Adapter
	client db.Client
	logger log.Logger

	// nodeId specifies the id of the node, this is used to determine which domains must be applied.
	nodeId string

	// localDomains is a buffer holding all domains installed on the local machine.
	// Map is used to avoid many libvirt list requests.
	localDomains map[string]string
	localDomainsLock sync.RWMutex
	localDomainsCycleTTL int64 // ttl of the cycle that renews the buffer.

	// stateSyncers is a simple register for tracking running state synchronizers.
	stateSyncers map[string]context.CancelFunc
	stateSyncersLock sync.RWMutex
	
	// configSyncers is a simple register for tracking running config synchronizers.
	configSyncers map[string]context.CancelFunc
	configSyncersLock sync.RWMutex

	// operationWg is a waitgroup used to track all detached operations (like syncers, pruners, etc).
	// It is used to correctly wait for all running operation before exiting the operator.
	operationWg sync.WaitGroup
}

type OperatorOption func(*Operator)

func NewOperator(client db.Client, controller domctrl.Controller, adapter domadapter.Adapter, opts ...OperatorOption) *Operator {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	operator := &Operator{
		rootCtx:         rootCtx,
		rootCtxCancel:   rootCtxCancel,
		workCtx:         workCtx,
		workCtxCancel:   workCtxCancel,
		finChan:         make(chan struct{}),
		controller: controller,
		adapter: adapter,
		client: client,
		logger:          discard.NewDiscardLogger(),
	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}

// WithLogger sets a custom logger for the domain operator.
func WithLogger(logger log.Logger) OperatorOption {
	return func(d *Operator) {
		d.logger = logger
	}
}

// ServeAndDetach starts the Operator reporting process in a detached goroutine.
func (d *Operator) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.apply()
	}()

	go func() {
		wg.Wait()
		d.finChan <- struct{}{}
	}()
}

// Terminate shuts down the domain operator gracefully, if shutdown did not complete in the provided context
// window the operator is terminated forcefully.
// Never returns an error (just there to match termination pattern).
func (d *Operator) Terminate(ctx context.Context) error {
	d.workCtxCancel()
	defer d.rootCtxCancel()
	select {
	case <-d.finChan:
		return nil
	case <-ctx.Done():
		d.rootCtxCancel()
		<-d.finChan
		return nil
	}
}
