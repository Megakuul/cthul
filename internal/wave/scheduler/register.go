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
	"time"

	"cthul.io/cthul/internal/wave/scheduler/resource"
)

// registerNode registers the local node periodically in the scheduler space on the database.
// As long as the node is registered, the scheduler assumes that it can move domains to the node.
// On every cycle the node resource capacity is measured and reported to the scheduler.
func (s *Scheduler) registerNode() {
	if s.registerId=="" {
		s.logger.Info("scheduler", "local node will not be registered")
		return
	}

	resourceOperator := resource.NewResourceOperator(s.client, resource.WithLogger(s.logger))
	
	for {
		s.logger.Debug("scheduler", "measuring local node resource capacity...")
		nodeResources, err := resourceOperator.GenerateNodeResources(s.workCtx,
			s.registerCpuFactor, s.registerMemFactor,
		)
		if err!=nil {
			s.logger.Err("scheduler", err.Error())
		}
		
		if nodeResources!=nil {
			ctx, cancel := context.WithTimeout(s.workCtx, time.Second*time.Duration(s.registerTTL))
			defer cancel()
			err := resourceOperator.SetNodeResources(ctx,
				fmt.Sprintf("/WAVE/SCHEDULER/NODE/%s", s.registerId),
				(s.registerTTL * 2),
				nodeResources,
			)
			if err != nil {
				s.logger.Err("scheduler", "failed to register node")
			}
		}

		select {
		case <-time.After(time.Second*time.Duration(s.registerTTL)):
			break
		case <-s.workCtx.Done():
			err := s.client.Delete(s.rootCtx, fmt.Sprintf("/WAVE/SCHEDULER/NODE/%s", s.registerId))
			if err != nil {
				s.logger.Err("scheduler", "failed to unregister node before termination")
			}
			return
		}
	}
}
