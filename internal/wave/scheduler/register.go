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
	"fmt"
	"time"
	"context"
)

// registerNode registers the local node periodically in the scheduler space on the database.
// As long as the node is registered, the scheduler assumes that it can move domains to the node.
// On every cycle the node resource capacity is measured and reported to the scheduler.
func (s *Scheduler) registerNode() {
	if s.localNode.id=="" {
		s.logger.Info("scheduler", "local node will not be registered")
		return
	}
	
	for {
		s.logger.Debug("scheduler", "measuring local node resource capacity...")
		localNodeCapacity, err := generateNodeCapacity(s.workCtx, 0, s.localNode.cpuFactor, s.localNode.memFactor)
		if err!=nil {
			s.logger.Err("scheduler", err.Error())
		}
		
		if localNodeCapacity!=nil {
			ctx, cancel := context.WithTimeout(s.workCtx, time.Second*time.Duration(s.localNode.registerTTL))
			defer cancel()
			err := s.client.Set(ctx,
				fmt.Sprintf("/WAVE/SCHEDULER/NODE/%s", s.localNode.id),
				serializeNodeCapacity(localNodeCapacity),
				(s.localNode.registerTTL * 2),
			)
			if err != nil {
				s.logger.Err("scheduler", "failed to register node")
			}
		}

		select {
		case <-time.After(time.Second*time.Duration(s.localNode.registerTTL)):
			break
		case <-s.workCtx.Done():
			err := s.client.Delete(s.rootCtx, fmt.Sprintf("/WAVE/SCHEDULER/NODE/%s", s.localNode.id))
			if err != nil {
				s.logger.Err("scheduler", "failed to unregister node before termination")
			}
			return
		}
	}
}
