/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
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
	"strings"

	"cthul.io/cthul/pkg/adapter/domain"
	adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/domain/structure"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Controller provides an interface for wave domain related operations.
type Controller struct {
	client  db.Client
	adapter domain.Adapter
}

type Option func(*Controller)

func New(client db.Client, adapter domain.Adapter, opts ...Option) *Controller {
	controller := &Controller{
		client:  client,
		adapter: adapter,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// affinity provides the domain affinity tags structure stored in the database.
type affinity []string

// resources provides the domain resource structure stored in the database.
type resources struct {
	AllocatedCpu    float64 `json:"allocated_cpu"`
	AllocatedMemory int64   `json:"allocated_memory"`
}

// List returns a map containing domain ids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]structure.Domain, error) {
	domains := map[string]structure.Domain{}

	reqnodes, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/REQNODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching domain reqnode: %w", err)
	}
	nodes, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/NODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching domain node: %w", err)
	}
	configs, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/CONFIG/")
	if err != nil {
		return nil, fmt.Errorf("fetching domain configs: %w", err)
	}
	states, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/STATE/")
	if err != nil {
		return nil, fmt.Errorf("fetching domain state: %w", err)
	}
	affinities, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/AFFINITY/")
	if err != nil {
		return nil, fmt.Errorf("fetching domain affinity: %w", err)
	}
	domResources, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/RESOURCES/")
	if err != nil {
		return nil, fmt.Errorf("fetching domain resources: %w", err)
	}

	for key, domConfig := range configs {
		var domainErr error
		id := strings.TrimPrefix(key, "/WAVE/DOMAIN/NODE/")
		reqnode := reqnodes[fmt.Sprint("/WAVE/DOMAIN/REQNODE/", id)]
		node := nodes[fmt.Sprint("/WAVE/DOMAIN/NODE/", id)]
		state := states[fmt.Sprint("/WAVE/DOMAIN/STATE/", id)]

		config := &adapterstruct.Domain{}
		err = json.Unmarshal([]byte(domConfig), config)
		if err != nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing node config: %w", err))
		}

		affinity := &affinity{}
		err = json.Unmarshal([]byte(affinities[fmt.Sprint("/WAVE/DOMAIN/AFFINITY/", id)]), affinity)
		if err != nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing node affinity tags: %w", err))
		}

		resources := &resources{}
		err = json.Unmarshal([]byte(domResources[fmt.Sprint("/WAVE/DOMAIN/RESOURCES/", id)]), resources)
		if err != nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing domain resources: %w", err))
		}

		domains[id] = structure.Domain{
			Reqnode:         reqnode,
			Node:            node,
			Config:          config,
			Affinity:        *affinity,
			State:           structure.DOMAIN_STATE(state),
			AllocatedCPU:    resources.AllocatedCpu,
			AllocatedMemory: resources.AllocatedMemory,
			Error:           domainErr,
		}
	}
	return domains, nil
}

// GetStats returns the current statistics of the domain. The data is read directly from the vmm (e.g. qemu).
func (c *Controller) GetStats(ctx context.Context, id string) (*adapterstruct.DomainStats, error) {
	stats, err := c.adapter.GetStats(ctx, id)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// Create creates a domain with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (c *Controller) Create(ctx context.Context, config *adapterstruct.Domain) (string, error) {
	id := uuid.New().String()
	err := c.SetConfig(ctx, id, config)
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
	err = c.SetState(ctx, id, structure.DOMAIN_DOWN)
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
	err = c.SetAffinity(ctx, id, []string{})
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
	return id, nil
}

// SetAffinity updates the affinity of the domain. The affinity specifies a list of affinity tags; any node
// tagged with at least one of those tags is considered eligible for hosting. An empty list means all nodes are
// eligible.
func (c *Controller) SetAffinity(ctx context.Context, id string, tags []string) error {
	rawTags, err := json.Marshal(affinity(tags))
	if err != nil {
		return fmt.Errorf("cannot serialize affinity tags: %w", err)
	}
	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id), string(rawTags), 0)
	if err != nil {
		return err
	}
	return nil
}

// SetState updates the desired state of the domain.
func (c *Controller) SetState(ctx context.Context, id string, state structure.DOMAIN_STATE) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id), string(state), 0)
	if err != nil {
		return err
	}
	return nil
}

// SetConfig updates the domain configuration and updates associated metadata on the database.
func (c *Controller) SetConfig(ctx context.Context, id string, config *adapterstruct.Domain) error {
	rawConfig, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot serialize config: %w", err)
	}

	rawResources, err := json.Marshal(resources{
		AllocatedCpu:    float64(config.ResourceConfig.VCPUs),
		AllocatedMemory: config.ResourceConfig.Memory,
	})
	if err != nil {
		return fmt.Errorf("cannot serialize resource config: %w", err)
	}
	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/RESOURCES/%s", id), string(rawResources), 0)
	if err != nil {
		return err
	}

	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id), string(rawConfig), 0)
	if err != nil {
		return err
	}
	return nil
}

// Lookup searches for the domain by id and returns its configuration.
func (c *Controller) Lookup(ctx context.Context, id string) (*structure.Domain, error) {
	reqnode, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/REQNODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain reqnode: %w", err)
	}
	node, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain node: %w", err)
	}
	domConfig, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain configs: %w", err)
	}
	state, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain state: %w", err)
	}
	domAffinity, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain affinity: %w", err)
	}
	domResources, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/RESOURCES/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain resources: %w", err)
	}

	config := &adapterstruct.Domain{}
	err = json.Unmarshal([]byte(domConfig), config)
	if err != nil {
		return nil, fmt.Errorf("parsing node config: %w", err)
	}

	affinity := &affinity{}
	err = json.Unmarshal([]byte(domAffinity), affinity)
	if err != nil {
		return nil, fmt.Errorf("parsing node affinity tags: %w", err)
	}

	resources := &resources{}
	err = json.Unmarshal([]byte(domResources), resources)
	if err != nil {
		return nil, fmt.Errorf("parsing domain resources: %w", err)
	}

	return &structure.Domain{
		Reqnode:         reqnode,
		Node:            node,
		Config:          config,
		Affinity:        *affinity,
		State:           structure.DOMAIN_STATE(state),
		AllocatedCPU:    resources.AllocatedCpu,
		AllocatedMemory: resources.AllocatedMemory,
	}, nil
}

// Attach requests the domain to be relocated to the specified node and waits until it's ready (if wait flag is set).
func (c *Controller) Attach(ctx context.Context, id, node string, wait bool) error {
	if !wait {
		_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/REQNODE/%s", id), node, 0)
		if err != nil {
			return err
		}
		return nil
	}

	pollCtx, pollCtxCancel := context.WithCancel(ctx)
	pollG, pollGCtx := errgroup.WithContext(pollCtx)

	pollG.Go(func() error {
		err := c.client.Watch(pollGCtx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id), func(_, activeNode string, err error) {
			if err == nil && node == activeNode {
				pollCtxCancel()
			}
		})
		if err != nil {
			return err
		}
		return nil
	})

	// initial check, required in case the node is already set to the requested node (watch will not trigger in this case)
	pollG.Go(func() error {
		activeNode, err := c.client.Get(pollGCtx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
		if err != nil {
			return err
		}
		if node == activeNode {
			pollCtxCancel()
		}
		return nil
	})

	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/REQNODE/%s", id), node, 0)
	if err != nil {
		return err
	}

	err = pollG.Wait()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("context exceeded: domain couldn't be attached in the provided context window")
	default:
		return nil
	}
}

// Detach removes the domain from the current node. It doesn't wait until the node fully detached it.
func (c *Controller) Detach(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/REQNODE/%s", id))
	if err != nil {
		return nil
	}
	return nil
}

// Delete completely removes a domain and its associated metadata.
func (c *Controller) Delete(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/REQNODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/RESOURCES/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err != nil {
		return err
	}
	return nil
}
