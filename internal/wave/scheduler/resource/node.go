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

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// NodeResources contains resource information about a node.
type NodeResources struct {
	AvailableCpuCores float64 `json:"available_cpu_cores"`
	AvailableMemBytes int64   `json:"available_mem_bytes"`
	TotalCpuCores     float64 `json:"total_cpu_cores"`
	TotalMemBytes     int64   `json:"total_mem_bytes"`
}

// GetNodeResources fetches the node resources of one single node.
func (r *ResourceOperator) GetNodeResources(ctx context.Context, key string) (*NodeResources, error) {
	nodeResourceStr, err := r.client.Get(ctx, key)
	if err!=nil {
		return nil, err
	}
	
	resources := NodeResources{}
	err = json.Unmarshal([]byte(nodeResourceStr), &resources)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node resources")
	}

	return &resources, nil
}

// IndexNodeResources fetches all node resources in one batch. Specify the key that can be used to find
// the node resources. Returns a map holding nodes->nodeResources.
func (r *ResourceOperator) IndexNodeResources(ctx context.Context, key string) (map[string]NodeResources, error){
	nodeResources, err := r.client.GetRange(ctx, key)
	if err!=nil {
		return nil, err
	}

	indexMap := map[string]NodeResources{}
	for nodeKey, resourceStr := range nodeResources {
		node := strings.TrimPrefix(nodeKey, key)
		resources := NodeResources{}
		err := json.Unmarshal([]byte(resourceStr), &resources)
		if err!=nil {
			r.logger.Warn("scheduler", "failed to parse node resources; skipping node...")
			continue
		}
		indexMap[node] = resources
	}
	return indexMap, nil
}

// GenerateNodeResources generates a node resources from the specs of the local machine.
// The cpu and mem factor is applied to all total resources before calculating further.
func (r *ResourceOperator) GenerateNodeResources(ctx context.Context, cpuFactor float64, memFactor float64) (*NodeResources, error) {
	resources := NodeResources{}

	cpuCores, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire cpu information")
	}
	var totalCpuCores int32 = 0
	for _, core := range cpuCores {
		totalCpuCores += core.Cores
	}
	// total cpu cores * factor (e.g. 10 * 0.8 = 8 cores)
	resources.TotalCpuCores = float64(totalCpuCores) * cpuFactor

	cpuLoad, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to measure cpu load")
	}
	// cpu factor - cpu load * total cpu cores (e.g. (0.8 - 0.4) * 10 = 4 cores)
	resources.AvailableCpuCores = (cpuFactor * 100 - cpuLoad[0]) / 100 * float64(totalCpuCores)

	memUsage, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire virtual memory information")
	}
	// total mem bytes * memfactor (e.g. 4096 * 0.8 = 3276 bytes)
	resources.TotalMemBytes = int64(float64(memUsage.Total) * memFactor)
	// factored mem bytes * used mem bytes (e.g. 3276 - 2000 = 1276 bytes)
	resources.AvailableMemBytes = resources.TotalMemBytes - int64(memUsage.Used)

	return &resources, nil
}

// SetNodeResources sets the node resources for one single node.
func (r *ResourceOperator) SetNodeResources(ctx context.Context, key string, ttl int64, resources *NodeResources) error {
	nodeResourceStr, err := json.Marshal(resources)
	if err != nil {
		return fmt.Errorf("failed to serialize node resources")
	}

	_, err = r.client.Set(ctx, key, string(nodeResourceStr), ttl)
	if err!=nil {
		return err
	}

	return nil
}
