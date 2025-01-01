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

	params, _, err := l.client.DomainGetCPUStats(domain, 0, 0, 0, 0)
	if err!=nil {

	}

	switch params[0].Field {
	case "cpu_time":

	case "system_time":

	case "user_time":
		
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

	stats, err := l.client.DomainMemoryStats(domain, uint32(libvirt.DomainMemoryStatNr), 0)
	if err!=nil {
		return nil, err
	}

	for _, stat := range stats {
		switch libvirt.DomainMemoryStatTags(stat.Tag) {
		case libvirt.DomainMemoryStatActualBalloon:
		case libvirt.DomainMemoryStatAvailable:
		case libvirt.DomainMemoryStatDiskCaches:
		case libvirt.DomainMemoryStatHugetlbPgalloc:
		case libvirt.DomainMemoryStatHugetlbPgfail:
		case libvirt.DomainMemoryStatLastUpdate:
		case libvirt.DomainMemoryStatMajorFault:
		case libvirt.DomainMemoryStatMinorFault:
		case libvirt.DomainMemoryStatNr:
		case libvirt.DomainMemoryStatRss:
		case libvirt.DomainMemoryStatSwapIn:
		case libvirt.DomainMemoryStatSwapOut:
		case libvirt.DomainMemoryStatUnused:
		case libvirt.DomainMemoryStatUsable:
		}
	}
	
	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetInterfaceStats(ctx context.Context, id string) (*structure.InterfaceStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetBlockStats(ctx context.Context, id string) (*structure.BlockStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}
