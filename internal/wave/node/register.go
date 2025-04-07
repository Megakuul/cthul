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
	"fmt"
	"time"

	"cthul.io/cthul/pkg/wave/node"
	"cthul.io/cthul/pkg/wave/node/structure"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// register registers the local node periodically in the cluster. This process reports the nodes
// associated state and allows other wave components like the scheduler to discover it.
// On every cycle the node state & resources are measured and reevaluated.
func (n *Operator) register() {
	nodeController := node.NewController(n.client)
	
	for {
		ctx, cancel := context.WithTimeout(n.workCtx, time.Second*time.Duration(n.cycleTTL))
		defer cancel()
		
		n.logger.Debug("measuring local node resource capacity...")
		node, err := n.acquireNodeInfo(ctx)
		if err!=nil {
			n.logger.Error(fmt.Sprintf("cannot report node state: %s", err.Error()))
		} else {
			err = nodeController.Register(ctx, n.nodeId, *node, n.cycleTTL*2)
			if err != nil {
				n.logger.Error(fmt.Sprintf("failed to register node: %s", err.Error()))
			}
		}

		select {
		case <-time.After(time.Second*time.Duration(n.cycleTTL)):
			break
		case <-n.workCtx.Done():
			err = nodeController.Unregister(n.rootCtx, n.nodeId)
			if err != nil {
				n.logger.Error("failed to unregister node before termination")
			}
			return
		}
	}
}

// acquireNodeInfo acquires an informational node structure by reading local machine specs (cpu, mem, etc)
// and further attributes statically defined on the operator.
func (n *Operator) acquireNodeInfo(ctx context.Context) (*structure.Node, error) {
	node := structure.Node{
		Affinity: n.affinity,
		State: structure.NODE_HEALTHY,
	}

	if n.maintenance {
		node.State = structure.NODE_MAINTENANCE
	}

	cpuCores, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire cpu information")
	}
	var totalCpuCores int32 = 0
	for _, core := range cpuCores {
		totalCpuCores += core.Cores
	}
	// total cpu cores * factor (e.g. 10 * 0.8 = 8 cores)
	node.AllocatedCpu = float64(totalCpuCores) * n.cpuFactor

	cpuLoad, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to measure cpu load")
	}
	// cpu factor - cpu load * total cpu cores (e.g. (0.8 - 0.4) * 10 = 4 cores)
	node.AvailableCpu = (n.cpuFactor * 100 - cpuLoad[0]) / 100 * float64(totalCpuCores)

	memoryUsage, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire virtual memory information")
	}
	// total mem bytes * memfactor (e.g. 4096 * 0.8 = 3276 bytes)
	node.AllocatedMemory = int64(float64(memoryUsage.Total) * n.memoryFactor)
	// factored mem bytes * used mem bytes (e.g. 3276 - 2000 = 1276 bytes)
	node.AvailableMemory = node.AllocatedMemory - int64(memoryUsage.Used)

	return &node, nil
}
