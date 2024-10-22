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

package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// nodeResources contains resource information about a node.
type nodeResources struct {
	// Specifies the currently available cpu cores / memory.
	// primarly used for "can I even move?" decisions.
	AvailableCpuCores float64 `json:"available_cpu_cores"`
	AvailableMemBytes int64   `json:"available_mem_bytes"`
	// Specifies the total cpu cores / memory.
	// primarly used for "where can I move?" decisions.
	TotalCpuCores float64 `json:"total_cpu_cores"`
	TotalMemBytes int64   `json:"total_mem_bytes"`
}

// parseNodeResources converts the resource string into a resource struct.
func parseNodeResources(resourceStr string) (*nodeResources, error) {
	resources := nodeResources{}
	err := json.Unmarshal([]byte(resourceStr), &resources)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node resources")
	}

	return &resources, nil
}

// serializeNodeResources serializes the resource struct into a string.
func serializeNodeResources(resources *nodeResources) string {
	resourceStr, err := json.Marshal(resources)
	if err != nil {
		return ""
	}

	return string(resourceStr)
}

// generateNodeResources generates a nodeResources struct from the resources of the local machine.
// The cpu and mem factor is applied to all total resources before calculating further.
func generateNodeResources(ctx context.Context, cpuFactor float64, memFactor float64) (*nodeResources, error) {
	resources := nodeResources{}

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


// domainResources contains resource information about a domain.
type domainResources struct {
	TotalCpuCores float64 `json:"total_cpu_cores"`
	TotalMemBytes int64 `json:"total_mem_bytes"`
}

// parseDomainResources converts the resource string into a resource struct.
func parseDomainResources(resourceStr string) (*domainResources, error) {
	resources := domainResources{}
	err := json.Unmarshal([]byte(resourceStr), &resources)
	if err != nil {
		return nil, fmt.Errorf("failed to parse domain resources")
	}

	return &resources, nil
}

// serializeDomainResources serializes the resource struct into a string.
func serializeDomainResources(resources *domainResources) string {
	resourceStr, err := json.Marshal(resources)
	if err != nil {
		return ""
	}

	return string(resourceStr)
}

// generateDomainResources generates a domainResources struct. This function does currently not load any data,
// because it is usually more efficient to batch load the domain resource data (e.g. indexDomainResources()).
// Currently this function never returns an error, however it is still in place for future expansions.
func generateDomainResources(cpuCores float64, memBytes int64) (*domainResources, error) {
	return &domainResources{
		// Specifies total cpu/mem of one domain.
		// primarly used for "can I even move?" decisions.
		TotalCpuCores: cpuCores,
		TotalMemBytes: memBytes,
	}, nil
}
