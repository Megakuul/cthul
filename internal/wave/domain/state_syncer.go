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
	"fmt"
	"strconv"
	"sync"
	"time"
	
	"cthul.io/cthul/pkg/wave/domain/structure"
)

// updateStateSyncer updates the state syncer map based on the new node. If the domain is located on the local
// node, a state syncer is started. If a state syncer associated with the domain is running, but the new node is
// not == the local one, the state syncer is stopped.
// The state syncer consists of two routines, one routine operates in periodical ticks pushing the state in a
// configured interval. The other routine watches the database and pushes state changes immediately upon update.
func (o *Operator) updateStateSyncer(uuid, node string) {
	o.stateSyncersLock.Lock()
	defer o.stateSyncersLock.Unlock()
	
	if node!=o.nodeId {
		if cancel, ok := o.stateSyncers[uuid]; ok {
			o.logger.Debug("domain-operator", fmt.Sprintf(
				"removing state synchronizer for domain '%s' on node '%s'...", uuid, node,
			))
			cancel()
			delete(o.stateSyncers, uuid)
		}
		return
	}

	o.logger.Debug("domain-operator", fmt.Sprintf(
		"setting up state synchronizer for domain '%s' on node '%s'...", uuid, node,
	))
	
	if _, ok := o.stateSyncers[uuid]; ok {
		o.logger.Debug("domain-operator", fmt.Sprintf(
			"state synchronization for domain '%s' is already running on '%s'; skipping setup...", uuid, node,
		))
		return
	}

	syncerCtx, cancel := context.WithCancel(o.workCtx)
	o.stateSyncers[uuid] = cancel
	
	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		syncerWg := sync.WaitGroup{}

		intervalStr, err := o.client.Get(syncerCtx, fmt.Sprintf("/WAVE/DOMAIN/STATEPUSHINT/%s", uuid))
		if err!=nil {
			o.logger.Err("domain-operator", fmt.Sprintf("failed to load domain state push interval: %s", err.Error()))
		}
		interval, err := strconv.Atoi(intervalStr)
		if err!=nil {
			interval = 30
			o.logger.Err("domain-operator", fmt.Sprintf(
				"failed to parse domain state push interval: %s; defaulting to '%ds'...", err.Error(), interval,
			))
		}
		
		syncerWg.Add(1)
		go func() {
			defer syncerWg.Done()
			for {
				ctx, cancel := context.WithTimeout(syncerCtx, time.Duration(interval) * time.Second)
				defer cancel()
				
				state, err := o.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", uuid))
				if err!=nil {
					o.logger.Err("domain-operator", fmt.Sprintf("failed to load domain state: %s", err.Error()))
				} else {
					err = o.applyState(uuid, state)
					if err!=nil {
						o.logger.Err("domain-operator", fmt.Sprintf("cannot apply domain '%s' state: %s", uuid, err.Error()))
					} else {
						o.logger.Info("domain-operator", fmt.Sprintf("updated state of domain '%s' to '%s'", uuid, state))
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
			err := o.client.WatchRange(syncerCtx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", uuid), func(_, v string, err error) {
				if err!=nil {
					o.logger.Err("domain-operator", fmt.Sprintf(
						"failed to load domain '%s' state: %s", uuid, err.Error(),
					))
					return
				}
				err = o.applyState(uuid, v)
				if err!=nil {
					o.logger.Err("domain-operator", fmt.Sprintf("cannot apply domain '%s' state: %s", uuid, err.Error()))
				} else {
					o.logger.Info("domain-operator", fmt.Sprintf("updated state of domain '%s' to '%s'", uuid, v))
				}
			})
			if err!=nil {
				o.logger.Crit("domain-operator", fmt.Sprintf(
					"failed to watch domain '%s' state: %s; exiting state watcher...", uuid, err.Error(),
				))
			}
		}()

		syncerWg.Wait()

		o.stateSyncersLock.Lock()
		defer o.stateSyncersLock.Unlock()
		if cancel, ok := o.stateSyncers[uuid]; ok {
			cancel()
			delete(o.stateSyncers, uuid)
		}
	}()
}

// applyState tries to apply the desired power state to the local domain.
func (o *Operator) applyState(uuid, state string) error {
	switch structure.DOMAIN_STATE(state) {
	case structure.DOMAIN_UP:
		err := o.adapter.Start(o.workCtx, uuid)
		if err!=nil {
			return fmt.Errorf("failed to start domain: %w", err)
		}
	case structure.DOMAIN_PAUSE:
		err := o.adapter.Pause(o.workCtx, uuid)
		if err!=nil {
			return fmt.Errorf("failed to pause domain: %w", err)
		}
	case structure.DOMAIN_DOWN:
		err := o.adapter.Shutdown(o.workCtx, uuid)
		if err!=nil {
			return fmt.Errorf("failed to shutdown domain: %w", err)
		}
	case structure.DOMAIN_FORCED_DOWN:
		err := o.adapter.Kill(o.workCtx, uuid)
		if err!=nil {
			return fmt.Errorf("failed to kill domain: %w", err)
		}
	}
	return nil
}
