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

package libvirt

import (
	"context"
	"fmt"
	"time"

	"cthul.io/cthul/pkg/adapter/domain/structure"
	"github.com/digitalocean/go-libvirt"
)

func (l *LibvirtAdapter)GetDomainStats(ctx context.Context, id string) (*structure.DomainStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return nil, err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return nil, err
	}

	params, err := l.client.ConnectGetAllDomainStats([]libvirt.Domain{domain}, 0, 0)
	if err!=nil {
		return nil, err
	}

	if len(params) < 1 {
		return nil, fmt.Errorf("domain with id '%s' not found", id)
	}

	for _, param := range params {
		switch param {
		case libvirt.DomainCPUStatsCputime:
			cpuStats.CpuTime = int64(param.Value.D)
		case libvirt.DomainCPUStatsSystemtime:
			cpuStats.KernelTime = int64(param.Value.D)
		case libvirt.DomainCPUStatsUsertime:
			cpuStats.UserTime = int64(param.Value.D)
		}		
	}

	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetCpuStats(ctx context.Context, id string) (*structure.CpuStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	
	uuid, err := l.parseUUID(id)
	if err!=nil {
		return nil, err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return nil, err
	}

	params, _, err := l.client.DomainGetCPUStats(domain, 0, -1, 1, 0)
	if err!=nil {
		return nil, err
	}

	cpuStats := &structure.CpuStats{
		Timestamp: time.Now().Unix(),
	}

	for _, param := range params {
		switch param.Field {
		case libvirt.DomainCPUStatsCputime:
			cpuStats.CpuTime = int64(param.Value.D)
		case libvirt.DomainCPUStatsSystemtime:
			cpuStats.KernelTime = int64(param.Value.D)
		case libvirt.DomainCPUStatsUsertime:
			cpuStats.UserTime = int64(param.Value.D)
		}		
	}
	
	return cpuStats, nil
}

func (l *LibvirtAdapter)GetMemoryStats(ctx context.Context, id string) (*structure.MemoryStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
		uuid, err := l.parseUUID(id)
	if err!=nil {
		return nil, err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return nil, err
	}

	params, err := l.client.DomainMemoryStats(domain, uint32(libvirt.DomainMemoryStatNr), 0)
	if err!=nil {
		return nil, err
	}

	memoryStats := &structure.MemoryStats{
		Timestamp: time.Now().Unix(),
	}

	for _, param := range params {
		switch libvirt.DomainMemoryStatTags(param.Tag) {
		case libvirt.DomainMemoryStatSwapIn:
			memoryStats.SwapIn = int64(param.Val * 1000) // convert from kB
		case libvirt.DomainMemoryStatSwapOut:
			memoryStats.SwapOut = int64(param.Val * 1000) // convert from kB
		case libvirt.DomainMemoryStatMinorFault:
			memoryStats.MinorFaults = int64(param.Val)
		case libvirt.DomainMemoryStatMajorFault:
			memoryStats.MajorFaults = int64(param.Val)
		case libvirt.DomainMemoryStatHugetlbPgalloc:
			memoryStats.HugepageAllocations = int64(param.Val)
		case libvirt.DomainMemoryStatHugetlbPgfail:
			memoryStats.HugepageFailures = int64(param.Val)
		case libvirt.DomainMemoryStatActualBalloon:
			memoryStats.Balloned = int64(param.Val * 1000) // convert from kB
		case libvirt.DomainMemoryStatAvailable:
			memoryStats.Available = int64(param.Val * 1000) // convert from kB
		case libvirt.DomainMemoryStatUsable:
			memoryStats.Usable = int64(param.Val * 1000) // convert from kB
		case libvirt.DomainMemoryStatUnused:
			memoryStats.Unused = int64(param.Val * 1000) // convert from kB
		case libvirt.DomainMemoryStatRss:
			memoryStats.HostRSS = int64(param.Val * 1000) // convert from kB
		}
	}
	
	return memoryStats, nil
}
