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

package video

import (
	"context"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
	"cthul.io/cthul/pkg/syncer"
)

// Operator is responsible for applying the video device database state to the local node.
// Because the main video server is provided by the vmm, the wave video device is mainly responsible for
// ensuring the environment allows bootstrap of the video server endpoint (e.g. mkdirall() the socket path).
type Operator struct {
	client db.Client
	syncer *syncer.Syncer
	logger log.Logger

	// waveRunRoot specifies the wave base path for runtime files (unix-sockets and stuff).
	waveRunRoot string

	// nodeId specifies the id of the node, this is used to determine which video devices must be applieo.
	nodeId string

	// updateCycleTTL specifies the ttl of the cycle that updates the devices syncers
	// (cycle essentially finds out what device must be synced by this node).
	updateCycleTTL int64

	// pathCycleTTL defines the ttl of the cycle that prepares the unix socket path for the video device.
	pathCycleTTL int64
}

type Option func(*Operator)

func New(client db.Client, opts ...Option) *Operator {
	operator := &Operator{
		client:        client,
		syncer:        syncer.New(client),
		logger:        discard.NewDiscardLogger(),
		waveRunRoot:   "/run/cthul/wave/",
		nodeId:        "undefined",
    updateCycleTTL: 30,
    pathCycleTTL: 30,
	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}

// WithLogger sets a custom logger for the video operator.
func WithLogger(logger log.Logger) Option {
	return func(o *Operator) {
		o.logger = logger
	}
}

// WithWaveRunRoot defines a custom root path for wave runtime files (sockets, etc.).
func WithWaveRunRoot(path string) Option {
	return func(o *Operator) {
		o.waveRunRoot = path
	}
}

// WithNodeId specifies the id of the local node. This id is used to identify which videos
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

// WithPathCycleTTL defines a custom cycle interval for the device path syncer.
// Every cycle prepares the path for the device unix socket.
func WithPathCycleTTL(ttl int64) Option {
	return func(o *Operator) {
		o.pathCycleTTL = ttl
	}
}

func (o *Operator) ServeAndDetach() {
  o.synchronize()
}

func (o *Operator) Terminate(ctx context.Context) error {
  o.syncer.Shutdown()
  return nil
}
