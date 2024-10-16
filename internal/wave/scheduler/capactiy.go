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

// nodeCapacity contains resource capacity information about a node
// which are used for scheduler decisions.
type nodeCapacity struct {
	// primarly used for "can I even move?" decisions.
	AvailableCpuCores float64 `json:"available_cpu_cores"`
	AvailableMemBytes int64   `json:"available_mem_bytes"`
	// primarly used for "where can I move?" decisions.
	CpuCores float64 `json:"cpu_cores"`
	MemBytes int64   `json:"mem_bytes"`
}

// parseNodeCapacity converts the capacity string into a capacity struct.
func parseNodeCapacity(capacityStr string) (*nodeCapacity, error) {
	capacity := nodeCapacity{}
	err := json.Unmarshal([]byte(capacityStr), &capacity)
	if err != nil {
		return nil, fmt.Errorf("failed to parse capacity")
	}

	return &capacity, nil
}

// serializeNodeCapacity serializes the capacity struct into a string.
func serializeNodeCapacity(capacity *nodeCapacity) string {
	capacityStr, err := json.Marshal(capacity)
	if err != nil {
		return ""
	}

	return string(capacityStr)
}

// generateNodeCapacity generates a nodeCapacity struct from the resources of the local machine.
// The cpu and mem factor is applied to all total resources before calculating further.
func generateNodeCapacity(ctx context.Context, cpuFactor float64, memFactor float64) (*nodeCapacity, error) {
	capacity := nodeCapacity{}

	cpuCores, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire cpu information")
	}
	var totalCpuCores int32 = 0
	for _, core := range cpuCores {
		totalCpuCores += core.Cores
	}
	// total cpu cores * factor (e.g. 10 * 0.8 = 8 cores)
	capacity.CpuCores = float64(totalCpuCores) * cpuFactor

	cpuLoad, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to measure cpu load")
	}
	// cpu factor - cpu load * total cpu cores (e.g. (0.8 - 0.4) * 10 = 4 cores)
	capacity.AvailableCpuCores = (cpuFactor * 100 - cpuLoad[0]) / 100 * float64(totalCpuCores)

	memUsage, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire virtual memory information")
	}
	// total mem bytes * memfactor (e.g. 4096 * 0.8 = 3276 bytes)
	capacity.MemBytes = int64(float64(memUsage.Total) * memFactor)
	// factored mem bytes * used mem bytes (e.g. 3276 - 2000 = 1276 bytes)
	capacity.AvailableMemBytes = capacity.MemBytes - int64(memUsage.Used)

	return &capacity, nil
}
