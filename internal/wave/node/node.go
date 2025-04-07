/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
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

package node

import (
	"context"
	"sync"
	"log/slog"

	"cthul.io/cthul/pkg/db"
)

// Operator is responsible to monitor and measure the state and resources of the host node.
// Evaluated data is reported to the cluster via NodeController. This allows components to discover this node.
type Operator struct {
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc

	workCtx       context.Context
	workCtxCancel context.CancelFunc

	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the operator is fully cleaned up.
	finChan chan struct{}

	client db.Client
	logger *slog.Logger

	// nodeId specifies the id of the node that is reported to the cluster.
	nodeId string
	// cycleTTL specifies the interval for scheduler cycles.
	cycleTTL int64
	// maintenance specifies whether maintenance mode is enabled.
	maintenance bool
	// affinity holds node affinity tags used for scheduling decisions.
	affinity []string
	// cpuFactor specifies how much host cpu is incorporated to the reported values.
	cpuFactor float64
	// memoryFactor specifies how much host memory is incorporated to the reported values.
	memoryFactor float64
}

type OperatorOption func(*Operator)

func NewOperator(logger *slog.Logger, client db.Client, opts ...OperatorOption) *Operator {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	operator := &Operator{
		rootCtx:       rootCtx,
		rootCtxCancel: rootCtxCancel,
		workCtx:       workCtx,
		workCtxCancel: workCtxCancel,
		finChan:       make(chan struct{}),
		client:        client,
		logger:        logger.WithGroup("node-operator"),
		nodeId:        "undefined",
		cycleTTL:      5,
		maintenance:   false,
		affinity:      []string{},
		cpuFactor:     1,
		memoryFactor:  1,
	}

	for _, opt := range opts {
		opt(operator)
	}

	return operator
}


// WithNodeId specifies the id of the node that is reported to the cluster.
func WithNodeId(id string) OperatorOption {
	return func(n *Operator) {
		n.nodeId = id
	}
}

// WithCycleTTL defines a custom cycle interval. Every cycle measures resources and reports the state to the
// cluster.
func WithCycleTTL(ttl int64) OperatorOption {
	return func(n *Operator) {
		n.cycleTTL = ttl
	}
}

// WithMaintenance enables the maintenance mode. If enabled the node is reported with maintenance mode to
// the cluster.
func WithMaintenance(maintenance bool) OperatorOption {
	return func(n *Operator) {
		n.maintenance = maintenance
	}
}

// WithAffinity defines custom affinity tags. The defined affinity tags are reported to the cluster.
func WithAffinity(tags []string) OperatorOption {
	return func(n *Operator) {
		n.affinity = tags
	}
}

// WithResourceFactor defines custom resource factors. The resource factors specify how much host resources
// are incorporated to the reported resource values.
// (e.g. 10 cores with 8 available and a factor of 0.8 = 8 cores with 6 available).
func WithResourceFactor(cpuFactor, memoryFactor float64) OperatorOption {
	return func(n *Operator) {
		n.cpuFactor = cpuFactor
		n.memoryFactor = memoryFactor
	}
}

// ServeAndDetach starts the Operator reporting process in a detached goroutine.
func (n *Operator) ServeAndDetach() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		n.register()
	}()

	go func() {
		wg.Wait()
		n.finChan <- struct{}{}
	}()
}

// Terminate shuts down the node operator gracefully, if shutdown did not complete in the provided context window
// the operator is terminated forcefully. Never returns an error (just there to match termination pattern).
func (n *Operator) Terminate(ctx context.Context) error {
	n.workCtxCancel()
	defer n.rootCtxCancel()
	select {
	case <-n.finChan:
		return nil
	case <-ctx.Done():
		n.rootCtxCancel()
		<-n.finChan
		return nil
	}
}
