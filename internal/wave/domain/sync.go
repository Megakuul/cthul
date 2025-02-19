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
	"strings"
	"time"
)

// synchronize starts the core synchronization, this includes various routine operations all tracked with the
// operationWg.
// 1. Starts a simple routine to load the localDomains buffer periodically.
// 2. Starts an operation for periodically pruning old unused domains.
// 3. Starts a routine that periodically reads all existing cluster domains and starts/stops "syncers" based
// on their managing node (if the local node is managing, the syncers are started otherwise they are stopped).
// Once started syncers will synchronize their respective attribute from the database to the local node.
// 4. Starts a watcher that also updates "syncers" but does so immediately after a change has been committed.
// Additionally this watcher also prunes domains that are manually removed from the database immediately.
//
// The reason for this approach with two separate routines (periodic/watcher) is their responsibility:
// Periodic routines are responsible for deterministically synchronizing the database state to the local node.
// Watcher routines operate incrementally and can theoretically miss events (e.g. on database interruption or
// when the state updates on the node without being reported to the database), their purpose is to avoid
// waiting for the next interval on manual updates; a user starting a domain wants this to happen immediately.
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
		for {
			ctx, cancel := context.WithTimeout(o.workCtx, time.Duration(o.pruneCycleTTL) * time.Second)
			defer cancel()

			o.pruneAllDomains(ctx)

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
		for {
			ctx, cancel := context.WithTimeout(o.workCtx, time.Duration(o.updateCycleTTL) * time.Second)
			defer cancel()

			domains, err := o.client.GetRange(ctx, "/WAVE/DOMAIN/NODE/")
			if err!=nil {
				o.logger.Err("domain-operator", fmt.Sprintf(
					"failed to load domains: %s; skipping update cycle...", err.Error(),
				))
			} else {
				for k, v := range domains {
					uuid := strings.TrimPrefix(k, "/WAVE/DOMAIN/NODE/")
					o.updateStateSyncer(uuid, v)
					o.updateConfigSyncer(uuid, v)
				}
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
		err := o.client.WatchRange(o.workCtx, "/WAVE/DOMAIN/NODE/", func(k, v string, err error) {
			if err!=nil {
				o.logger.Err("domain-operator", fmt.Sprintf("failed to retrieve domains from database: %s", err.Error()))
				return
			}
			uuid := strings.TrimPrefix(k, "/WAVE/DOMAIN/NODE/")
			o.updateStateSyncer(uuid, v)
			o.updateConfigSyncer(uuid, v)
			
			pruneCtx, pruneCancel := context.WithTimeout(o.workCtx, time.Duration(o.pruneCycleTTL) * time.Second)
			defer pruneCancel()
			o.pruneDomain(pruneCtx, uuid, v)
		})
		if err!=nil {
			o.logger.Crit("domain-operator", fmt.Sprintf(
				"failed to watch domains: %s; exiting update watcher...", err.Error(),
			))
		}
	}()

	o.operationWg.Wait()
}
