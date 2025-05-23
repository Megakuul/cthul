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

package serial

import (
	"context"
	"log/slog"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/syncer"
)

// Operator is responsible for applying the serial device database state to the local node.
// Because the main serial server is provided by the vmm, the wave serial device is mainly responsible for
// ensuring the environment allows bootstrap of the serial server endpoint (e.g. mkdirall() the socket path).
type Operator struct {
	client db.Client
	logger *slog.Logger
	syncer *syncer.Syncer

	// runRoot specifies the base path for runtime files (unix-sockets and stuff).
	runRoot string

	// nodeId specifies the id of the node, this is used to determine which serial devices must be applieo.
	nodeId string

	// updateCycleTTL specifies the ttl of the cycle that updates the devices syncers
	// (cycle essentially finds out what device must be synced by this node).
	updateCycleTTL int64

	// syncCycleTTL defines the ttl of the cycle that prepares the unix socket path for the serial device.
	syncCycleTTL int64
}

type Option func(*Operator)

func New(logger *slog.Logger, client db.Client, opts ...Option) *Operator {
	operator := &Operator{
		client:         client,
		logger:         logger.WithGroup("serial-operator"),
		syncer:         syncer.New(logger.WithGroup("serial-operator"), client),
		runRoot:        "/run/cthul/wave/",
		nodeId:         "undefined",
		updateCycleTTL: 30,
		syncCycleTTL:   30,
	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}

// WithRunRoot defines a custom root path for runtime files (sockets, etc.).
func WithRunRoot(path string) Option {
	return func(o *Operator) {
		o.runRoot = path
	}
}

// WithNodeId specifies the id of the local node. This id is used to identify which serials
// must be synced to this node.
func WithNodeId(id string) Option {
	return func(n *Operator) {
		n.nodeId = id
	}
}

// WithUpdateCylceTTL defines a custom update cycle interval.
// Every cycle checks which devices are managed by the local node and fully refreshes all running syncers.
// Syncers are also incrementally updated in realtime when updated on database, this cycle just fully resyncs.
func WithUpdateCylceTTL(ttl int64) Option {
	return func(o *Operator) {
		o.updateCycleTTL = ttl
	}
}

// WithSyncCycleTTL defines a custom cycle interval for the device syncer.
// Every cycle prepares u.a. the path for the device unix socket.
func WithSyncCycleTTL(ttl int64) Option {
	return func(o *Operator) {
		o.syncCycleTTL = ttl
	}
}

func (o *Operator) ServeAndDetach() {
	o.synchronize()
}

func (o *Operator) Terminate(ctx context.Context) error {
	o.syncer.Shutdown()
	return nil
}
