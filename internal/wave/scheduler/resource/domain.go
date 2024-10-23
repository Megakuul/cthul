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
	"fmt"
	"strconv"
	"strings"
)


// DomainResources contains resource information about a domain.
type DomainResources struct {
	TotalCpuCores float64 `json:"total_cpu_cores"`
	TotalMemBytes int64   `json:"total_mem_bytes"`
}

// GetDomainResources fetches the domain resources of one domain. Specify the base key (e.g. '/WAVE/DOMAIN') and
// the domain id. The function expects to find the resources under '$key/CPU/$id' and '$key/MEM/$id'.
func (r *ResourceOperator) GetDomainResources(ctx context.Context, key, id string) (*DomainResources, error) {
	domainCpu, err := r.client.Get(ctx, fmt.Sprintf("%s/CPU/%s", key, id))
	if err!=nil {
		return nil, err
	}
	domainMem, err := r.client.Get(ctx, fmt.Sprintf("%s/MEM/%s", key, id))
	if err!=nil {
		return nil, err
	}

	cpuCores, err := strconv.ParseFloat(domainCpu, 64)
	if err!=nil {
		return nil, fmt.Errorf("failed to parse cpu")
	}
	memBytes, err := strconv.Atoi(domainMem)
	if err!=nil {
		return nil, fmt.Errorf("failed to parse memory")
	}

	return &DomainResources{
		TotalCpuCores: cpuCores,
		TotalMemBytes: int64(memBytes),
	}, nil
}

// IndexDomainResources fetches all domain resources of the cluster in one batch. Specify the base key
// (e.g. '/WAVE/DOMAIN'). The function expects to find resources under '$key/NODE/', '$key/CPU/' and '$key/MEM/'.
// Returns a hierarchical map holding nodes->domains->domainResources.
func (r *ResourceOperator) IndexDomainResources(ctx context.Context, key string) (map[string]map[string]DomainResources, error) {
	var (
		NODE_PREFIX = fmt.Sprintf("%s/NODE/", key)
		CPU_PREFIX = fmt.Sprintf("%s/CPU/", key)
		MEM_PREFIX= fmt.Sprintf("%s/MEM/", key)
	)
	
	domainNodes, err := r.client.GetRange(ctx, NODE_PREFIX)
	if err!=nil {
		return nil, err
	}
	domainCpus, err := r.client.GetRange(ctx, CPU_PREFIX)
	if err!=nil {
		return nil, err
	}
	domainMems, err := r.client.GetRange(ctx, MEM_PREFIX)
	if err!=nil {
		return nil, err
	}

	indexMap := map[string]map[string]DomainResources{}
	for domainKey, node := range domainNodes {
		domain := strings.TrimPrefix(NODE_PREFIX, domainKey)
		if indexMap[node] == nil {
			indexMap[node] = make(map[string]DomainResources)
		}
		cpuCores, err := strconv.ParseFloat(domainCpus[CPU_PREFIX + domain], 64)
		if err!=nil {
			r.logger.Warn("scheduler", fmt.Sprintf(
				"failed to parse cpu for domain %s; using zero value...", domain,
			))
			cpuCores = 0
		}
		memBytes, err := strconv.Atoi(domainMems[MEM_PREFIX + domain])
		if err!=nil {
			r.logger.Warn("scheduler", fmt.Sprintf(
				"failed to parse memory for domain %s; using zero value...", domain,
			))
			memBytes = 0
		}
		indexMap[node][domain] = DomainResources{
			TotalCpuCores: cpuCores,
			TotalMemBytes: int64(memBytes),
		}
	}
	return indexMap, nil
}
