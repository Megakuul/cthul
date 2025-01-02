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
	"strings"
	"sync"
	"time"

	adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	"cthul.io/cthul/pkg/wave/domain/structure"
)

func (o *Operator) synchronize() {
	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		for {
			ctx, cancel := context.WithTimeout(o.workCtx, time.Duration(o.localDomainsCycleTTL) * time.Second)
			defer cancel()
			
			domains, err := o.adapter.List(ctx)
			if err!=nil {
				o.logger.Err("domain-operator", fmt.Sprintf("failed to load local domains: %s", err.Error()))	
			}

			o.localDomainsLock.Lock()
			o.localDomains = domains
			o.localDomainsLock.Unlock()
			
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
		
	}()

	
	domains, err := o.client.GetRange(o.workCtx, "/WAVE/DOMAIN/NODE/")
	if err!=nil {
		o.logger.Err("domain-operator", fmt.Sprintf(
			"failed to load domains: %s; skipping initialization...", err.Error(),
		))
	} else {
		for k, v := range domains {
			uuid := strings.TrimPrefix(k, "/WAVE/DOMAIN/NODE/")
			o.updateStateSyncer(uuid, v)
			o.updateConfigSyncer(uuid, v)
		}
	}
	
	err = o.client.WatchRange(o.workCtx, "/WAVE/DOMAIN/NODE/", func(k, v string, err error) {
		if err!=nil {
			o.logger.Err("domain-operator", fmt.Sprintf("failed to retrieve domains from database: %s", err.Error()))
			return
		}
		uuid := strings.TrimPrefix(k, "/WAVE/DOMAIN/NODE/")
		o.pruneDomain()
		o.updateStateSyncer(uuid, v)
		o.updateConfigSyncer(uuid, v)
	})
	if err!=nil {
		o.logger.Crit("domain-operator", fmt.Sprintf(
			"failed to watch domains: %s; exiting operator...", err.Error(),
		))
	}
	o.operationWg.Wait()
}

// forcePruneDomains compares the local domains with the desired domains on the database. All domains that are
// present on the node but not on the database (managed by this node) are gracefully removed. If this does not
// succeed in specified time, the domains are forcefully removed without releasing external devices.
func (o *Operator) forcePruneDomains() {
	// TODO
}

// pruneDomain checks and removes the domain gracefully based on the new node. If the domain is present on the
// local node and is now managed by another node, its config is pulled to then gracefully destroy the domain
// (this includes releasing stuff like granit, proton and wave devices).
func (o *Operator) pruneDomain(ctx context.Context, uuid, node string) {
	if node==o.nodeId {
		return
	}
	
	if _, ok := o.localDomains[uuid]; ok {
		o.logger.Debug("domain-operator", fmt.Sprintf("starting graceful destruction of domain '%s'", uuid))
		
		rawConfig, err := o.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", uuid))
		if err!=nil {
			o.logger.Err("domain-operator", fmt.Sprintf(
				"failed to load domain '%s' config: %s; skipping graceful destruction...", uuid, err.Error(),
			))
			return
		}
		config := &adapterstruct.Domain{}
		err = json.Unmarshal([]byte(rawConfig), config)
		if err!=nil {
			o.logger.Err("domain-operator", fmt.Sprintf(
				"failed to parse domain '%s' config: %s; skipping graceful destruction...", uuid, err.Error(),
			))
		}
		
		err = o.adapter.Destroy(ctx, uuid, *config)
		if err!=nil {
			o.logger.Err("domain-operator", fmt.Sprintf("failed to destroy domain '%s': %s", uuid, err.Error()))
		}
		
		o.logger.Info("domain-operator", fmt.Sprintf(
			"removed local domain '%s'; domain is now managed by '%s'...", uuid, node,
		))
	}
}

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

// applyConfig tries to apply the domain configuration to the local domain.
func (o *Operator) applyConfig(uuid, rawConfig string) error {
	config := &adapterstruct.Domain{}
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
