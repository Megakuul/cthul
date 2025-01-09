/**
 * Cthul System
 *
 * Copyright (C) 2025 Linus Ilian Moser <linus.moser@megakuul.ch>
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

package video

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (o *Operator) synchronize() {

	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		for {
			ctx, cancel := context.WithTimeout(o.workCtx, time.Duration(o.updateCycleTTL) * time.Second)
			defer cancel()
			
			deviceNodes, err := o.client.GetRange(ctx, "/WAVE/VIDEO/NODE")
			if err!=nil {
				o.logger.Err("video-operator", fmt.Sprintf(
					"failed to load requested video device nodes: %s; skipping update cycle...", err.Error(),
				))
			}
			deviceNodeRequests, err := o.client.GetRange(ctx, "/WAVE/VIDEO/REQNODE/")
			if err!=nil {
				o.logger.Err("video-operator", fmt.Sprintf(
					"failed to load requested video device nodes: %s; skipping update cycle...", err.Error(),
				))
			}

			for key, nodeRequest := range deviceNodeRequests {
				uuid := strings.TrimPrefix(key, "/WAVE/VIDEO/REQNODE/")
				node := deviceNodes[fmt.Sprint("/WAVE/VIDEO/NODE/", uuid)]

				o.updatePathSyncer(ctx, uuid, node, nodeRequest)
			}

			select {
			case <-o.workCtx.Done():
				return
			case <-ctx.Done():
			}
		}
	}()

	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		err := o.client.WatchRange(o.workCtx, "/WAVE/VIDEO/REQNODE/", func(key, nodeRequest string, err error) {
			if err!=nil {
				o.logger.Err("video-operator", fmt.Sprintf(
					"failed to retrieve video requested device from database: %s", err.Error(),
				))
				return
			}

			ctx, cancel := context.WithTimeout(o.workCtx, time.Duration(o.updateCycleTTL) * time.Second)
			defer cancel()
			
			uuid := strings.TrimPrefix(key, "/WAVE/VIDEO/REQNODE/")
			node, err := o.client.Get(ctx, fmt.Sprint("/WAVE/VIDEO/NODE/", uuid))
			if err!=nil {
				o.logger.Err("video-operator", fmt.Sprintf(
					"failed to retrieve video device '%s' node from database: %s", uuid, err.Error(),
				))
				return
			}

			o.updatePathSyncer(ctx, uuid, node, nodeRequest)
		})
		if err!=nil {
			o.logger.Crit("video-operator", fmt.Sprintf(
				"failed to watch video devices: %s; exiting update watcher...", err.Error(),
			))
		}
	}()

	o.operationWg.Wait()
}
