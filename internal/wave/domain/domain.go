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
  "log/slog"
	"context"
	"sync"

	domadapter "cthul.io/cthul/pkg/adapter/domain"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/syncer"
)

// Operator is responsible for applying the domains database state to the local virtual machine monitor.
type Operator struct {
  rootCtx context.Context
  rootCtxCancel context.CancelFunc

	adapter domadapter.Adapter
	client db.Client
  logger *slog.Logger
  syncer *syncer.Syncer

	// nodeId specifies the id of the node, this is used to determine which domains must be applieo.
	nodeId string

	// updateCycleTTL specifies the ttl of the cycle that updates the domain syncers
	// (cycle essentially finds out what domains must be synced by this node).
	updateCycleTTL int64
	
  // localCycleTTL specifies the ttl of the cycle that manually resyncs the local domains.
  // This includes caching local domains and removing orphaned ones. 
  localCycleTTL int64

  // stateCycleTTL specifies the ttl of the cycle that manually syncs the domain state.
  stateCycleTTL int64

  // configCycleTTL specifies the ttl of teh cylce that manually syncs the domain config.
  configCycleTTL int64

	// localDomains is a buffer holding all domains installed on the local machine.
	// Map is used to avoid many libvirt list requests.
	localDomains map[string]string
	localDomainsLock sync.RWMutex

  // operationWg is a waitgroup that captures workers that are not managed by the syncer.
  operationWg sync.WaitGroup
}

type Option func(*Operator)

func New(logger *slog.Logger, client db.Client, adapter domadapter.Adapter, opts ...Option) *Operator {
  rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	operator := &Operator{
    rootCtx: rootCtx,
    rootCtxCancel: rootCtxCancel,
		adapter: adapter,
		client: client,
		logger: logger.WithGroup("domain-operator"),
    syncer: syncer.New(logger.WithGroup("domain-operator"), client),
		nodeId: "undefined",
		updateCycleTTL: 10,
    localCycleTTL: 60,
    stateCycleTTL: 30,
    configCycleTTL: 30,
		localDomains: map[string]string{},
		localDomainsLock: sync.RWMutex{},
    operationWg: sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}


// WithNodeId specifies the id of the local node. This id is used to identify which domains
// must be synced to this node.
func WithNodeId(id string) Option {
	return func(n *Operator) {
		n.nodeId = id
	}
}

// WithUpdateCylceTTL defines a custom update cycle interval.
// Every cycle checks which domains are managed by the local node and fully refreshes all running syncers.
// Syncers are also incrementally updated in realtime when updated on database, this cycle just fully resyncs.
func WithUpdateCylceTTL(ttl int64) Option {
	return func(o *Operator) {
		o.updateCycleTTL = ttl
	}
}

// WithLocalCycle defines a custom local cycle interval.
// The cycle is responsible for synchronizing the local domains (e.g. caching or removing orphaned)
func WithLocalCycle(ttl int64) Option {
	return func(o *Operator) {
    o.localCycleTTL = ttl
	}
}

// WithStateCycleTTL defines a custom cycle interval for manually syncing the domain state (up, down, paused, etc.)
func WithStateCycleTTL(ttl int64) Option {
	return func(o *Operator) {
		o.stateCycleTTL = ttl
	}
}

// WithConfigCycleTTL defines a custom cycle interval for manually syncing the domain config.
func WithConfigCycleTTL(ttl int64) Option {
	return func(o *Operator) {
		o.configCycleTTL = ttl
	}
}

func (o *Operator) ServeAndDetach() {
  o.synchronize()
}

func (o *Operator) Terminate(ctx context.Context) error {
  o.syncer.Shutdown()
  o.rootCtxCancel()
  o.operationWg.Wait()
  return nil
}
