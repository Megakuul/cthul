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
	"strconv"
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
			panic("alarm")
		}
		
		domainNodes, err := s.client.GetRange(s.workCtx, "/WAVE/DOMAIN/NODE/")
		if err!=nil {
			s.logger.Err("scheduler", "failed to load domain nodes: " + err.Error())
		}

		for domain, node := range domainNodes {
			_, ok := schedulerNodes[node]
			if !ok {
				retries := unmanagedDomains[domain]
				unmanagedDomains[domain] = retries + 1
			} else {
				unmanagedDomains[domain] = 0
			}
		}

		// TODO: Perform the re scheduling for domains with more then 1 retry (or config)
	}
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
