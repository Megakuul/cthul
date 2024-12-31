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

// report registers the local node periodically in the cluster. This process announces the nodes
// associated state and allows other wave components like the scheduler to discover it.
// On every cycle the node state & resources are measured and reevaluated.
func (n *NodeOperator) report() {
	nodeController := node.NewNodeController(n.client)
	
	for {
		ctx, cancel := context.WithTimeout(n.workCtx, time.Second*time.Duration(n.cycleTTL))
		defer cancel()
		
		n.logger.Debug("scheduler", "measuring local node resource capacity...")
		report, err := n.generateReport(ctx)
		if err!=nil {
			n.logger.Err("scheduler", fmt.Sprintf("cannot report node state: %s", err.Error()))
		} else {
			err = nodeController.Register(ctx, n.nodeId, *report, n.cycleTTL*2)
			if err != nil {
				n.logger.Err("scheduler", fmt.Sprintf("failed to register node: %s", err.Error()))
			}
		}

		select {
		case <-time.After(time.Second*time.Duration(n.cycleTTL)):
			break
		case <-n.workCtx.Done():
			err = nodeController.Unregister(n.rootCtx, n.nodeId)
			if err != nil {
				n.logger.Err("scheduler", "failed to unregister node before termination")
			}
			return
		}
	}
}

// GenerateNodeResources generates a node resources from the specs of the local machine.
// The cpu and mem factor is applied to all total resources before calculating further.
func (n *NodeOperator) generateReport(ctx context.Context) (*structure.Node, error) {
	report := structure.Node{
		Affinity: n.affinity,
		State: structure.NODE_HEALTHY,
	}

	if n.maintenance {
		report.State = structure.NODE_MAINTENANCE
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
	report.AllocatedCpu = float64(totalCpuCores) * n.cpuFactor

	cpuLoad, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to measure cpu load")
	}
	// cpu factor - cpu load * total cpu cores (e.g. (0.8 - 0.4) * 10 = 4 cores)
	report.AvailableCpu = (n.cpuFactor * 100 - cpuLoad[0]) / 100 * float64(totalCpuCores)

	memoryUsage, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire virtual memory information")
	}
	// total mem bytes * memfactor (e.g. 4096 * 0.8 = 3276 bytes)
	report.AllocatedMemory = int64(float64(memoryUsage.Total) * n.memoryFactor)
	// factored mem bytes * used mem bytes (e.g. 3276 - 2000 = 1276 bytes)
	report.AvailableMemory = report.AllocatedMemory - int64(memoryUsage.Used)

	return &report, nil
}
