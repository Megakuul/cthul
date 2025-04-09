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
	"slices"
	"strconv"
	"time"

	domstruct "cthul.io/cthul/pkg/wave/domain/structure"
	nodestruct "cthul.io/cthul/pkg/wave/node/structure"
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
	// unmanagedDomain holds unmanaged domains and the number of cycles they were already unmanaged.
	// it is used to avoid immediate rescheduling of unmanagedDomains.
	unmanagedDomains := map[string]int{}
	
	next, err := s.client.Get(schedulerCtx, "/WAVE/SCHEDULER/NEXT")
	if err!=nil {
		s.logger.Error("failed to fetch next scheduler cycle initially; initiating schedule...")
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
			s.logger.Error("failed to update scheduler cycle; waiting for next cycle...")
			continue
		}
		if parseTime(prevNext).After(time.Now()) {
			s.logger.Debug("scheduler possibly double contested; waiting for next cycle...")
			continue
		}
		

		domains, err := s.domainController.List(s.workCtx)
		if err!=nil {
			s.logger.Error("failed to load domains: " + err.Error())
			continue
		}

		nodes, err := s.nodeController.List(s.workCtx)
		if err!=nil {
			s.logger.Error("failed to load nodes: " + err.Error())
			continue
		}
		
		for domainId, domain  := range domains {
			if domain.Error != nil {
				s.logger.Warn(fmt.Sprintf(
					"skipping scheduler analysis for '%s': domain information is malformed: %s", domainId, domain.Error,
				))
				continue
			}
			
			_, ok := nodes[domain.Reqnode]
			if !ok {
				retries := unmanagedDomains[domainId]
				unmanagedDomains[domainId] = retries + 1
			} else {
				unmanagedDomains[domainId] = 0
			}

			if unmanagedDomains[domainId] < int(s.rescheduleCycles) {
				continue
			}
			
			targetNodeId, targetNode, err := s.findNode(domain, domains, nodes)
			if err!=nil {
				s.logger.Warn(fmt.Sprintf(
					"skipping reschedule for '%s': %s", domainId, err.Error(),
				))
				continue
			}

			err = s.domainController.Attach(s.workCtx, domainId, domain.Node, false)
			if err!=nil {
				s.logger.Error(fmt.Sprintf(
					"failed to reschedule '%s': %s", domainId, err.Error(),
				))
				continue
			}
			
			domain.Reqnode = targetNodeId			
			domains[domainId] = domain
			nodes[targetNodeId] = *targetNode
		}
	}
}

// findNode evaluates the optimal node to move the domain to. Returns the new target node id and its associated
// node information. The assumed resource impact of the new domain is already factored in.
func (s *Scheduler) findNode(
	domain domstruct.Domain,
	domains map[string]domstruct.Domain,
	nodes map[string]nodestruct.Node,
) (string, *nodestruct.Node, error) {

	eligibleNodes := map[string]nodestruct.Node{}
	for nodeId, node := range nodes {
		if node.Error != nil {
			continue
		}
		if node.State != nodestruct.NODE_HEALTHY {
			continue
		}
		if !checkAffinity(domain.Affinity, node.Affinity) {
			continue
		}
		eligibleNodes[nodeId] = node
	}
	if len(eligibleNodes) < 1 {
		return "", nil, fmt.Errorf("no healthy cluster node with matching affinity tags available")
	}
	
	// constants define the usage factor that a domain is assumed to consume.
	// this is a heuristic to "guess" how much cpu/mem the domain will actually consume on the cluster node.
	// defaulting to 100% is a pretty dumb idea because most domains don't use 100% of their provisioned capacity.
	const DOMAIN_CPU_USAGE_FACTOR_HEURISTIC = 0.3
	const DOMAIN_MEM_USAGE_FACTOR_HEURISTIC = 0.6
	
	assumedCpuUsage := domain.AllocatedCPU * DOMAIN_CPU_USAGE_FACTOR_HEURISTIC
	assumedMemUsage := int64(float64(domain.AllocatedMemory) * DOMAIN_MEM_USAGE_FACTOR_HEURISTIC)

	// filter out nodes that currently do not provide enough available resources.
	// This is done to prevent moving a domain to a node that has currently not sufficient capacity.
	// For example if a high load cluster node failsover, instead of moving all domains at once to another
	// available node, every cycle just moves the amount of nodes the currently fit within the current capacity.
	availableNodes := map[string]nodestruct.Node{}
	for nodeId, node := range nodes {
		if node.AvailableCpu > assumedCpuUsage && node.AvailableMemory > assumedMemUsage {
			availableNodes[nodeId] = node
		}
	}
	if len(availableNodes) < 1 {
		return "", nil, fmt.Errorf("no cluster node provides sufficient available resources for this domain")
	}

	// constants define what factor is applied to the resources to normalize them as rating points.
	// this defines the relation between cpu cores and memory bytes which is set to 1 core = 1 gb = 1 point
	const (
		CPU_POINT_FACTOR = 1.0
		MEM_POINT_FACTOR = 0.000000001
	)

	// variables define a weight for cpu/mem that is multiplied with the cpu/mem points.
	// This influences the algorithm to weight either cpu or mem stronger, therefore it is set to the resource
	// requirements of the domain (e.g. if domain has 8 cores and 2GB memory, the cpu weight is 8 and the mem 2)
	var (
		CPU_POINT_WEIGHT = domain.AllocatedCPU * CPU_POINT_FACTOR
		MEM_POINT_WEIGHT = float64(domain.AllocatedMemory) * MEM_POINT_FACTOR
	)

	// search the node with the highest rating, rating points are calculated based on a simple formula:
	// ((nodeTotalCpu - nodeAllocatedCpu) * CPU_WEIGHT) + ((nodeTotalMem - nodeAllocatedMem) * MEM_WEIGHT)
	// This finds the node that best fits the capacity of the domain in question.
	chosenNodeId, chosenRating := "", 0.0
	for nodeId, node := range availableNodes {
		availableCpu := node.AllocatedCpu
		availableMem := node.AllocatedMemory
		for _, dom := range domains {
			if dom.Reqnode == nodeId {
				availableCpu -= dom.AllocatedCPU
				availableMem -= dom.AllocatedMemory
			}
		}

		cpuRating := availableCpu * CPU_POINT_FACTOR * CPU_POINT_WEIGHT
		memRating := float64(availableMem) * MEM_POINT_FACTOR * MEM_POINT_WEIGHT
		nodeRating := cpuRating + memRating
		
		if chosenNodeId == "" {
			chosenNodeId, chosenRating = nodeId, nodeRating
		} else if chosenRating < nodeRating {
			chosenNodeId, chosenRating = nodeId, nodeRating
		}
	}

	chosenNode := availableNodes[chosenNodeId]
	chosenNode.AvailableCpu -= assumedCpuUsage
	chosenNode.AvailableMemory -= assumedMemUsage
	return chosenNodeId, &chosenNode, nil
}


// checkAffinity checks if any affinity searchTag is contained in the targetTags.
func checkAffinity(searchTags, targetTags []string) bool {
	for _, searchTag := range searchTags {
    if slices.Contains(targetTags, searchTag) {
      return true
    }
	}
	return false
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
