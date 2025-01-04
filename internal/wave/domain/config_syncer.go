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

package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"cthul.io/cthul/pkg/adapter/domain/structure"
)

// updateConfigSyncer updates the config syncer map based on the new node. If the domain is located on the local
// node, a config syncer is started. If a config syncer associated with the domain is running, but the new node
// is not == the local one, the config syncer is stopped.
// The config syncer consists of two routines, one routine operates in periodical ticks pushing the config in a
// configured interval. The other routine watches the database and pushes config changes immediately upon update.
func (o *Operator) updateConfigSyncer(uuid, node string) {
	o.configSyncersLock.Lock()
	defer o.configSyncersLock.Unlock()

	if node!=o.nodeId {
		if cancel, ok := o.configSyncers[uuid]; ok {
			o.logger.Debug("domain-operator", fmt.Sprintf(
				"removing config synchronizer for domain '%s' on node '%s'...", uuid, node,
			))
			cancel()
			delete(o.configSyncers, uuid)
		}
		return
	}

	o.logger.Debug("domain-operator", fmt.Sprintf(
		"setting up config synchronizer for domain '%s' on node '%s'...", uuid, node,
	))
	
	if _, ok := o.configSyncers[uuid]; ok {
		o.logger.Debug("domain-operator", fmt.Sprintf(
			"config synchronization for domain '%s' is already running on '%s'; skipping setup...", uuid, node,
		))
		return
	}

	syncerCtx, cancel := context.WithCancel(o.workCtx)
	o.configSyncers[uuid] = cancel
	
	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		syncerWg := sync.WaitGroup{}

		intervalStr, err := o.client.Get(syncerCtx, fmt.Sprintf("/WAVE/DOMAIN/CONFIGPUSHINT/%s", uuid))
		if err!=nil {
			o.logger.Err("domain-operator", fmt.Sprintf("failed to load domain config push interval: %s", err.Error()))
		}
		interval, err := strconv.Atoi(intervalStr)
		if err!=nil {
			interval = 30
			o.logger.Err("domain-operator", fmt.Sprintf(
				"failed to parse domain config push interval: %s; defaulting to '%ds'...", err.Error(), interval,
			))
		}
	
		syncerWg.Add(1)
		go func() {
			defer syncerWg.Done()
			for {
				ctx, cancel := context.WithTimeout(syncerCtx, time.Duration(interval) * time.Second)
				defer cancel()
				
				config, err := o.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", uuid))
				if err!=nil {
					o.logger.Err("domain-operator", fmt.Sprintf("failed to load domain config: %s", err.Error()))
				} else {
					err = o.applyConfig(uuid, config)
					if err!=nil {
						o.logger.Err("domain-operator", fmt.Sprintf(
							"cannot apply domain '%s' config: %s", uuid, err.Error(),
						))
					} else {
						o.logger.Info("domain-operator", fmt.Sprintf("updated config of domain '%s'", uuid))
					}
				}

				select {
				case <-syncerCtx.Done():
					return
				case <-ctx.Done():
				}
			}
		}()

		syncerWg.Add(1)
		go func() {
			defer syncerWg.Done()
			err := o.client.WatchRange(syncerCtx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", uuid), func(_, v string, err error) {
				if err!=nil {
					o.logger.Err("domain-operator", fmt.Sprintf(
						"failed to load domain '%s' config: %s", uuid, err.Error(),
					))
					return
				}
				err = o.applyConfig(uuid, v)
				if err!=nil {
					o.logger.Err("domain-operator", fmt.Sprintf("cannot apply domain '%s' config: %s", uuid, err.Error()))
				} else {
					o.logger.Info("domain-operator", fmt.Sprintf("updated config of domain '%s'", uuid))
				}
			})
			if err!=nil {
				o.logger.Crit("domain-operator", fmt.Sprintf(
					"failed to watch domain '%s' config: %s; exiting config watcher...", uuid, err.Error(),
				))
			}
		}()

		syncerWg.Wait()
		
		o.configSyncersLock.Lock()
		defer o.configSyncersLock.Unlock()
		if cancel, ok := o.configSyncers[uuid]; ok {
			cancel()
			delete(o.configSyncers, uuid)
		}
	}()		
}

// applyConfig tries to apply the domain configuration to the local domain.
func (o *Operator) applyConfig(uuid, rawConfig string) error {
	config := &structure.Domain{}
	err := json.Unmarshal([]byte(rawConfig), config)
	if err!=nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	
	err = o.adapter.Apply(o.workCtx, uuid, *config)
	if err!=nil {
		return fmt.Errorf("failed to apply domain config: %w", err)
	}

	return nil
}
