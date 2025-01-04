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
	"errors"
	"fmt"

	adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
)

// pruneAllDomains compares the local domains with the desired domains on the database. All domains that are
// present on the node but not on the database (managed by this node) are removed with pruneDomain().
func (o *Operator) pruneAllDomains(ctx context.Context) {
	domains, err := o.client.GetRange(o.workCtx, "/WAVE/DOMAIN/NODE/")
	if err!=nil {
		o.logger.Err("domain-operator", fmt.Sprintf(
			"failed to load domains: %s; skipping prune process...", err.Error(),
		))
		return
	}

	// create localdomain snapshot to avoid blocking the mutex while pruning.
	// Map is shallow copied, but because there are only immutable strings, it is practically a deep copy.
	o.localDomainsLock.RLock()
	localDomains := o.localDomains
	o.localDomainsLock.RUnlock()
	
	for uuid := range localDomains {
		select {
		case <-ctx.Done():
			return
		default:
		}
		
		if node, ok := domains[uuid]; !ok || node!=o.nodeId {
			o.pruneDomain(ctx, uuid, node)
		}
	}
}

// pruneDomain checks and removes the domain based on the new node. If the domain is present on the local node
// and is now managed by another node, its config is pulled to then gracefully destroy the domain (this includes
// releasing stuff like granit, proton and wave devices). If graceful destruction fails, the domain is
// forcefully removed from the host.
func (o *Operator) pruneDomain(ctx context.Context, uuid, node string) {
	if node==o.nodeId {
		return
	}

	o.localDomainsLock.RLock()
	_, ok := o.localDomains[uuid]
	o.localDomainsLock.RUnlock()
	if !ok {
		return
	}
	
	o.logger.Debug("domain-operator", fmt.Sprintf("starting graceful destruction of domain '%s'...", uuid))
	
	rawConfig, err := o.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", uuid))
	if err!=nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			o.logger.Warn("domain-operator", fmt.Sprintf(
				"failed to load domain '%s' config: %s; skipping forceful destruction...", uuid, err.Error(),
			))
			return
		} else {
			o.logger.Warn("domain-operator", fmt.Sprintf(
				"failed to load domain '%s' config: %s; starting forceful destruction...", uuid, err.Error(),
			))
			// defaulting to empty config means Destroy() is "forceful" as no cthul devices must be deallocated.
			rawConfig = "{}"
		} 
	}
	config := adapterstruct.Domain{}
	err = json.Unmarshal([]byte(rawConfig), &config)
	if err!=nil {
		o.logger.Warn("domain-operator", fmt.Sprintf(
			"failed to parse domain '%s' config: %s; starting forceful destruction...", uuid, err.Error(),
		))
	}
	
	err = o.adapter.Destroy(ctx, uuid, config)
	if err!=nil {
		o.logger.Err("domain-operator", fmt.Sprintf("failed to destroy domain '%s': %s", uuid, err.Error()))
		return
	}
	
	o.logger.Info("domain-operator", fmt.Sprintf(
		"removed local domain '%s'; domain is now managed by '%s'...", uuid, node,
	))
}
