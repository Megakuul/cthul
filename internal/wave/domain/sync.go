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

	"encoding/json"
	"errors"

	"cthul.io/cthul/pkg/api/wave/v1/domain"
	"google.golang.org/protobuf/proto"
)

// synchronize starts the core synchronization which is responsible for applying the database state to the local node.
func (o *Operator) synchronize() {
	o.operationWg.Add(1)
	go func() {
		defer o.operationWg.Done()
		for {
			ctx, cancel := context.WithTimeout(o.rootCtx, time.Duration(o.localCycleTTL)*time.Second)
			defer cancel()

			domains, err := o.adapter.List(ctx)
			if err != nil {
				o.logger.Error(fmt.Sprintf("failed to load local domains: %s", err.Error()))
			}

			o.localDomainsLock.Lock()
			o.localDomains = domains
			o.localDomainsLock.Unlock()

			o.pruneAllDomains(ctx)

			select {
			case <-o.rootCtx.Done():
				return
			case <-ctx.Done():
			}
		}
	}()

	o.syncer.Add("/WAVE/DOMAIN/REQNODE/", o.updateCycleTTL, func(ctx context.Context, k, reqnode string) error {
		id := strings.TrimPrefix(k, "/WAVE/DOMAIN/REQNODE/")
		configKey := fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id)
		if reqnode == o.nodeId {
			o.syncer.Add(configKey, o.configCycleTTL, func(ctx context.Context, k, v string) error {
        err := o.applyConfig(ctx, id, v)
        if err!=nil {
          return err
        }
				_, err = o.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id), reqnode, 0)
				if err != nil {
					return err
				}
        return nil 
			})
		} else {
			o.syncer.Remove(configKey, false)
			o.pruneDomain(ctx, id, reqnode)
      return nil
		}
		return nil
	})
}

// applyConfig tries to apply the domain configuration to the local domain.
func (o *Operator) applyConfig(ctx context.Context, id, rawConfig string) error {
	config := &domain.DomainConfig{}
	err := proto.Unmarshal([]byte(rawConfig), config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	err = o.adapter.Apply(ctx, id, config)
	if err != nil {
		return fmt.Errorf("failed to apply config: %w", err)
	}

	switch config.State {
	case domain.DomainState_DOMAIN_STATE_UP:
		err := o.adapter.Start(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to start domain: %w", err)
		}
	case domain.DomainState_DOMAIN_STATE_PAUSE:
		err := o.adapter.Pause(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to pause domain: %w", err)
		}
  case domain.DomainState_DOMAIN_STATE_DOWN:
		err := o.adapter.Shutdown(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to shutdown domain: %w", err)
		}
	case domain.DomainState_DOMAIN_STATE_FORCED_DOWN:
		err := o.adapter.Kill(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to kill domain: %w", err)
		}
	}

	return nil
}

// pruneAllDomains compares the local domains with the desired domains on the database. All domains that are
// present on the node but not on the database (managed by this node) are removed with pruneDomain().
func (o *Operator) pruneAllDomains(ctx context.Context) {
	domains, err := o.client.GetRange(o.rootCtx, "/WAVE/DOMAIN/REQNODE/")
	if err != nil {
		o.logger.Error(fmt.Sprintf(
			"failed to load domains: %s; skipping prune process...", err.Error(),
		))
		return
	}

	// create localdomain snapshot to avoid blocking the mutex while pruning.
	// Map is shallow copied, but because there are only immutable strings, it is practically a deep copy.
	o.localDomainsLock.RLock()
	localDomains := o.localDomains
	o.localDomainsLock.RUnlock()

	for id := range localDomains {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if node, ok := domains[id]; !ok || node != o.nodeId {
			o.pruneDomain(ctx, id, node)
		}
	}
}

// pruneDomain checks and removes the domain based on the new node. If the domain is present on the local node
// and is now managed by another node, its config is pulled to then gracefully destroy the domain (this includes
// releasing stuff like granit, proton and wave devices). If graceful destruction fails, the domain is
// forcefully removed from the host.
func (o *Operator) pruneDomain(ctx context.Context, id, node string) {
	if node == o.nodeId {
		return
	}

	o.localDomainsLock.RLock()
	_, ok := o.localDomains[id]
	o.localDomainsLock.RUnlock()
	if !ok {
		return
	}

	o.logger.Debug(fmt.Sprintf("starting graceful destruction of domain '%s'...", id))

	rawConfig, err := o.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			o.logger.Warn(fmt.Sprintf(
				"failed to load domain '%s' config: %s; skipping forceful destruction...", id, err.Error(),
			))
			return
		} else {
			o.logger.Warn(fmt.Sprintf(
				"failed to load domain '%s' config: %s; starting forceful destruction...", id, err.Error(),
			))
			// defaulting to empty config means Destroy() is "forceful" as no cthul devices must be deallocated.
			rawConfig = "{}"
		}
	}
	config := &domain.DomainConfig{}
	err = json.Unmarshal([]byte(rawConfig), config)
	if err != nil {
		o.logger.Warn(fmt.Sprintf(
			"failed to parse domain '%s' config: %s; starting forceful destruction...", id, err.Error(),
		))
	}

	err = o.adapter.Destroy(ctx, id, config)
	if err != nil {
		o.logger.Error(fmt.Sprintf("failed to destroy domain '%s': %s", id, err.Error()))
		return
	}

	// technically removing the domain from the cache is not strictly required (it's not avoiding race conditions!)
	// however, it avoids unnecessary calls to pruneDomain() until the cache is renewed.
	o.localDomainsLock.Lock()
	delete(o.localDomains, id)
	o.localDomainsLock.Unlock()

	o.logger.Info(fmt.Sprintf(
		"removed local domain '%s'; domain is now managed by '%s'...", id, node,
	))
}
