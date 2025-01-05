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

	adapter domadapter.Adapter
	client db.Client
	logger log.Logger

	// nodeId specifies the id of the node, this is used to determine which domains must be applied.
	nodeId string

	// updateCycleTTL specifies the ttl of the cycle that updates the domain syncers
	// (cycle essentially finds out what domains must be synced by this node).
	updateCycleTTL int64
	
	// pruneCycleTTL specifies the ttl of the cycle that prunes unused domains from the local node.
	// Value is also used as context timeout of prune watcher calls (prune calls made by update watcher).
	pruneCycleTTL int64

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

func NewOperator(client db.Client, adapter domadapter.Adapter, opts ...OperatorOption) *Operator {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	operator := &Operator{
		rootCtx:         rootCtx,
		rootCtxCancel:   rootCtxCancel,
		workCtx:         workCtx,
		workCtxCancel:   workCtxCancel,
		finChan:         make(chan struct{}),
		adapter: adapter,
		client: client,
		logger:          discard.NewDiscardLogger(),
		nodeId: "undefined",
		updateCycleTTL: 10,
		pruneCycleTTL: 30,
		localDomains: map[string]string{},
		localDomainsLock: sync.RWMutex{},
		localDomainsCycleTTL: 30,
		stateSyncers: map[string]context.CancelFunc{},
		stateSyncersLock: sync.RWMutex{},
		configSyncers: map[string]context.CancelFunc{},
		configSyncersLock: sync.RWMutex{},
		operationWg: sync.WaitGroup{},
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

// WithNodeId specifies the id of the local node. This id is used to identify which domains
// must be synced to this node.
func WithNodeId(id string) OperatorOption {
	return func(n *Operator) {
		n.nodeId = id
	}
}

// WithUpdateCylceTTL defines a custom update cycle interval.
// Every cycle checks which domains are managed by the local node and fully refreshes all running syncers.
// Syncers are also incrementally updated in realtime when updated on database, this cycle just fully resyncs.
func WithUpdateCylceTTL(ttl int64) OperatorOption {
	return func(d *Operator) {
		d.updateCycleTTL = ttl
	}
}

// WithPruneCylceTTL defines a custom prune cycle interval.
// Every cycle checks which local domains can be removed and destroys them.
// Domains are also incrementally pruned in realtime when deleted, this cycle just rechecks all domains.
func WithPruneCylceTTL(ttl int64) OperatorOption {
	return func(d *Operator) {
		d.pruneCycleTTL = ttl
	}
}

// WithLocalDomainsCylceTTL defines a custom local domains cycle interval.
// Every cycle reads the local domains into an internal buffer (avoids listing all domains on every operation).
func WithLocalDomainsCylceTTL(ttl int64) OperatorOption {
	return func(d *Operator) {
		d.localDomainsCycleTTL = ttl
	}
}

// ServeAndDetach starts the Operator reporting process in a detached goroutine.
func (d *Operator) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.synchronize()
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
