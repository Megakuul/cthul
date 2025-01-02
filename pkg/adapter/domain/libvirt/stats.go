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

package libvirt

import (
	"context"
	"fmt"

	"cthul.io/cthul/pkg/adapter/domain/structure"
	"github.com/digitalocean/go-libvirt"
)

// GetStats collects all domain statistics from the vmm (e.g. qemu) in one batch. Some returned values
// may be set to 0 or -1 if required guest os drivers or vmm features are missing.
func (l *Adapter)GetStats(ctx context.Context, id string) (*structure.DomainStats, error) {
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

	domainRecords, err := l.client.ConnectGetAllDomainStats([]libvirt.Domain{domain}, uint32(
		libvirt.DomainStatsState |
			libvirt.DomainStatsCPUTotal |
			libvirt.DomainStatsBalloon |
			libvirt.DomainStatsVCPU |
			libvirt.DomainStatsInterface |
			libvirt.DomainStatsBlock,
	), 0)
	if err!=nil {
		return nil, err
	}

	if len(domainRecords) < 1 {
		return nil, fmt.Errorf("domain with id '%s' not found", id)
	}

	domainStats := &structure.DomainStats{}

	// TODO: Test how this undocumented weird structure actually works and then implement the parsing step.
	for _, param := range domainRecords[0].Params {
		val := int64(param.Value.D)
		switch param.Field {
		case "state.state":
			domainStats.State = structure.DOMAIN_STATE(val)
		case "cpu.time":
			domainStats.Cpu.CpuTime = val
		case "cpu.user":
			domainStats.Cpu.UserTime = val
		case "cpu.system":
			domainStats.Cpu.KernelTime = val
		case "vcpu.current":
			domainStats.Cpu.VCpus = []structure.VCpuStats{}
		case "balloon.current":
			domainStats.Memory.Balloned = val * 1024
		case "balloon.available":
			domainStats.Memory.Available = val * 1024
		case "balloon.usable":
			domainStats.Memory.Usable = val * 1024
		case "balloon.unused":
			domainStats.Memory.Unused = val * 1024
		case "balloon.rss":
			domainStats.Memory.HostRSS = val * 1024
		case "balloon.swap_in":
			domainStats.Memory.SwapIn = val * 1024
		case "balloon.swap_out":
			domainStats.Memory.SwapOut = val * 1024
		case "balloon.major_fault":
			domainStats.Memory.MajorFaults = val
		case "balloon.minor_fault":
			domainStats.Memory.MinorFaults = val
		case "balloon.hugetlb_pgalloc":
			domainStats.Memory.HugepageAllocations = val
		case "balloon.hugetlb_pgfail":
			domainStats.Memory.HugepageFailures = val
		}
	}

	return nil, fmt.Errorf("not implemented")
}
