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
	"fmt"
	"strconv"
	"strings"
	"time"
)

// startSchedulerCycle starts a scheduler cycle. This cycle executes periodically based next schedule stored
// in the database. When starting a cycle, the scheduler updates this schedule with the current time + cycleTTL,
// if multiple schedulers run at the same time (usually not the case) the first one updating the schedule key
// will executes the cycle, all others wait till the next cycle.
// One cycle captures all domains that use nodes that are not registered in the scheduler. If those domains
// are captured in the subsequent request as well, the scheduler assigns them to one of the active nodes
// based on their current capacity.
// The schedulerCtx can be cancelled to stop the scheduler, this will stop the scheduler AFTER the current cycle.
func (s *Scheduler) startSchedulerCycle(schedulerCtx context.Context) {
	unmanagedDomains := map[string]int{}
	
	next, err := s.client.Get(schedulerCtx, "/WAVE/SCHEDULER/NEXT")
	if err!=nil {
		s.logger.Err("scheduler", "failed to fetch next scheduler cycle initially; initiating schedule...")
	}
	nextSchedule := parseTime(next)
	
	for {
		select {
		case <- schedulerCtx.Done():
			return
		case <- s.workCtx.Done():
			return
		case <- time.After(time.Until(nextSchedule)):
			break
		}

		nextSchedule = time.Now().Add(time.Second * time.Duration(s.cycleTTL))
		prevNext, err := s.client.Set(s.workCtx, "/WAVE/SCHEDULER/NEXT",
			serializeTime(nextSchedule), 0,
		)
		if err!=nil {
			s.logger.Err("scheduler", "failed to update scheduler cycle; waiting for next cycle...")
			continue
		}
		if parseTime(prevNext).After(time.Now()) {
			s.logger.Debug("scheduler", "scheduler possibly double contested; waiting for next cycle...")
			continue
		}

		schedulerNodes, err := s.client.GetRange(s.workCtx, "/WAVE/SCHEDULER/NODE/")
		if err!=nil {
			s.logger.Err("scheduler", "failed to load scheduler nodes: " + err.Error())
			continue
		}
		
		domainNodes, err := s.client.GetRange(s.workCtx, "/WAVE/DOMAIN/NODE/")
		if err!=nil {
			s.logger.Err("scheduler", "failed to load domain nodes: " + err.Error())
			continue
		}

		for domain, node := range domainNodes {
			_, ok := schedulerNodes[node]
			if !ok {
				retries := unmanagedDomains[domain]
				unmanagedDomains[domain] = retries + 1
			} else {
				unmanagedDomains[domain] = 0
			}

			if unmanagedDomains[domain] >= int(s.rescheduleCycles) {
				newNode, newCap, err := s.findNode(domain, schedulerNodes)
				if err!=nil {
					s.logger.Warn("scheduler", "skipping reschedule: " + err.Error())
					continue
				}
				_, err = s.client.Set(s.workCtx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", domain), newNode, 0)
				if err!=nil {
					s.logger.Err("scheduler", "failed to reschedule node: " + err.Error())
					continue
				}
				schedulerNodes[newNode] = newCap
			}
		}
	}
}

// findNode evaluates the optimal node to move the domain to.
// Returns the node id and its assumed new capacity (guessed based on heuristics).
func (s *Scheduler) findNode(ctx context.Context, domain string, nodes map[string]nodeResources) (string, *nodeResources, error) {
	domainCpuStr, err := s.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CPU/%s", domain))
	if err != nil {
		return "", nil, err
	}
	domainCpu, err := strconv.ParseFloat(domainCpuStr, 64)
	if err!=nil {
		return "", nil, fmt.Errorf("failed to parse domain cpu cores")
	}
	domainMemStr, err := s.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/MEM/%s", domain))
	if err != nil {
		return "", nil, err
	}
	domainMem, err := strconv.Atoi(domainMemStr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse domain memory")
	}

	// constants define the usage factor that a domain is assumed to consume.
	// this is a heuristic to "guess" how much cpu/mem the domain will actually consume on the cluster node.
	// defaulting to 100% is a pretty dumb idea because most domains don't use 100% of their provisioned capacity.
	const DOMAIN_CPU_USAGE_FACTOR_HEURISTIC = 0.3
	const DOMAIN_MEM_USAGE_FACTOR_HEURISTIC = 0.6
	
	domainCpu = domainCpu * DOMAIN_CPU_USAGE_FACTOR_HEURISTIC
	domainMem = int(float64(domainMem) * DOMAIN_MEM_USAGE_FACTOR_HEURISTIC)

	// filter out nodes that currently do not provide enough available resources.
	// This is done to prevent moving a domain to a node that has currently not sufficient capacity.
	// For example if a high load cluster node failsover, instead of moving all domains at once to another
	// available node, every cycle just moves the amount of nodes the currently fit within the current capacity.
	availableNodes := map[string]nodeResources{}
	for node, cap := range nodes {
		if cap.AvailableCpuCores > domainCpu && cap.AvailableMemBytes > int64(domainMem) {
			availableNodes[node] = cap
		}
	}
	if len(availableNodes) < 1 {
		return "", nil, fmt.Errorf("no cluster node provides sufficient available resources for this domain")
	}

	theChosenNode := ""
	for node, cap := range availableNodes {
		if chosen, ok := availableNodes[theChosenNode]; ok {

		}
	}
	
	return "", nil, nil
}

// indexDomainResources performs a full resource index of all domains in the cluster.
// This operation is used to avoid high pressure on etcd; it only performs few lighweight range requests.
// It returns a map with all nodes and their managed domains.
func (s *Scheduler) indexDomainResources(ctx context.Context) (map[string]map[string]domainResources, error) {
	domainNodes, err := s.client.GetRange(ctx, "/WAVE/DOMAIN/NODE/")
	if err!=nil {
		return nil, err
	}

	domainCpus, err := s.client.GetRange(ctx, "/WAVE/DOMAIN/CPU/")
	if err!=nil {
		return nil, err
	}

	domainMems, err := s.client.GetRange(ctx, "/WAVE/DOMAIN/MEM/")
	if err!=nil {
		return nil, err
	}

	indexMap := map[string]map[string]domainResources{}
	for key, node := range domainNodes {
		domain := strings.TrimPrefix("/WAVE/DOMAIN/NODE/", key)
		if indexMap[node] == nil {
			indexMap[node] = make(map[string]domainResources)
		}
		cpuCores, err := strconv.ParseFloat(domainCpus["/WAVE/DOMAIN/CPU/" + domain], 64)
		if err!=nil {
			s.logger.Warn("scheduler", fmt.Sprintf(
				"failed to parse cpu for domain %s; using zero value...", domain,
			))
			cpuCores = 0
		}
		memBytes, err := strconv.Atoi(domainMems["/WAVE/DOMAIN/MEM/" + domain])
		if err!=nil {
			s.logger.Warn("scheduler", fmt.Sprintf(
				"failed to parse memory for domain %s; using zero value...", domain,
			))
			memBytes = 0
		}
		resources, err := generateDomainResources(cpuCores, int64(memBytes))
		if err!=nil {
			s.logger.Warn("scheduler", fmt.Sprintf(
				"failed to generate resources for domain %s; using zero values...", domain,
			))
			continue
		}
		indexMap[node][domain] = *resources
	}

	return indexMap, nil
}

// parseTime converts a unix timestamp (sec) as string to time.Time. Returns 01.01.1970 if it fails to parse.
func parseTime(unixString string) time.Time {
	unixInt, err := strconv.Atoi(unixString)
	if err!=nil {
		return time.Unix(0, 0)
	}
	return time.Unix(int64(unixInt), 0)
}

// serializeTime converts a time.Time struct to a unix timestamp (sec) as string.
func serializeTime(unixTime time.Time) string {
	return strconv.Itoa(int(unixTime.Unix()))
}
