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
	"sync"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/discard"
)

// Operator is responsible for applying the video device database state to the local node.
// Because the main video server is provided by the vmm, the wave video device is mainly responsible for
// ensuring the environment allows bootstrap of the video server endpoint (e.g. mkdirall() the socket path).
type Operator struct {
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc

	workCtx       context.Context
	workCtxCancel context.CancelFunc

	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the operator is fully cleaned up.
	finChan chan struct{}

	client db.Client
	logger log.Logger

	// waveRunRoot specifies the wave base path for runtime files (unix-sockets and stuff).
	waveRunRoot string

	// nodeId specifies the id of the node, this is used to determine which video devices must be applied.
	nodeId string

	// operationWg is a waitgroup used to track all detached operations (like syncers, pruners, etc).
	// It is used to correctly wait for all running operation before exiting the operator.
	operationWg sync.WaitGroup
}

type OperatorOption func(*Operator)

func New(client db.Client, opts ...OperatorOption) *Operator {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	operator := &Operator{
		rootCtx:         rootCtx,
		rootCtxCancel:   rootCtxCancel,
		workCtx:         workCtx,
		workCtxCancel:   workCtxCancel,
		finChan:         make(chan struct{}),
		client: client,
		logger:          discard.NewDiscardLogger(),
		waveRunRoot: "/run/cthul/wave/",
		nodeId: "undefined",
		operationWg: sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(operator)
	}
	

	return operator
}

// WithLogger sets a custom logger for the video operator.
func WithLogger(logger log.Logger) OperatorOption {
	return func(d *Operator) {
		d.logger = logger
	}
}

// WithWaveRunRoot defines a custom root path for wave runtime files (sockets, etc.).
func WithWaveRunRoot(path string) OperatorOption {
	return func(d *Operator) {
		d.waveRunRoot = path
	}
}

// WithNodeId specifies the id of the local node. This id is used to identify which videos
// must be synced to this node.
func WithNodeId(id string) OperatorOption {
	return func(n *Operator) {
		n.nodeId = id
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

// Terminate shuts down the video operator gracefully, if shutdown did not complete in the provided context
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
