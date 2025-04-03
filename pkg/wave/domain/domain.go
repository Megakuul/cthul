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
	"strconv"
	"strings"

	"cthul.io/cthul/pkg/adapter/domain"
	adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/domain/structure"
	"github.com/google/uuid"
)

// Controller provides an interface for wave domain related operations.
type Controller struct {
	client db.Client
	adapter domain.Adapter
}

type Option func(*Controller)

func New(client db.Client, adapter domain.Adapter, opts ...Option) *Controller {
	controller := &Controller{
		client: client,
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
	AllocatedCpu float64 `json:"allocated_cpu"`
	AllocatedMemory int64 `json:"allocated_memory"`
}

// List returns a map containing domain uuids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]structure.Domain, error) {
	domains := map[string]structure.Domain{}
	
	domainNodes, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/NODE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain node: %w", err)
	}
	domainStates, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/STATE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain state: %w", err)
	}
	domainStatePushIntervals, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/STATEPUSHINT/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain state push interval: %w", err)
	}
	domainConfigPushIntervals, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/CONFIGPUSHINT/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain config push interval: %w", err)
	}
	domainAffinities, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/AFFINITY/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain affinity: %w", err)
	}
	domainResources, err := c.client.GetRange(ctx, "/WAVE/DOMAIN/RESOURCES/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain resources: %w", err)
	}
	
	for key, node := range domainNodes {
		var domainErr error
		uuid := strings.TrimPrefix(key, "/WAVE/DOMAIN/NODE/")
		state := domainStates[fmt.Sprint("/WAVE/DOMAIN/STATE/", uuid)]

		statePushInterval, err := strconv.Atoi(
			domainStatePushIntervals[fmt.Sprint("/WAVE/DOMAIN/STATEPUSHINT/", uuid)],
		)
		if err!=nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing state push interval: %w", err))
		}

		configPushInterval, err := strconv.Atoi(
			domainConfigPushIntervals[fmt.Sprint("/WAVE/DOMAIN/CONFIGPUSHINT/", uuid)],
		)
		if err!=nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing state push interval: %w", err))
		}

		affinity := &affinity{}
		err = json.Unmarshal([]byte(domainAffinities[fmt.Sprint("/WAVE/DOMAIN/AFFINITY/", uuid)]), affinity)
		if err!=nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing node affinity tags: %w", err))
		}
		
		resources := &resources{}
		err = json.Unmarshal([]byte(domainResources[fmt.Sprint("/WAVE/DOMAIN/RESOURCES/", uuid)]), resources)
		if err!=nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing domain resources: %w", err))
		}
		
		domains[uuid] = structure.Domain{
			Node: node,
			Affinity: *affinity,
			State: structure.DOMAIN_STATE(state),
			StatePushInterval: int64(statePushInterval),
			ConfigPushInterval: int64(configPushInterval),
			AllocatedCPU: resources.AllocatedCpu,
			AllocatedMemory: resources.AllocatedMemory,
			Error: domainErr,
		}
	}
	return domains, nil
}

// GetConfig returns the deserialized configuration of one node.
func (c *Controller) GetConfig(ctx context.Context, id string) (*adapterstruct.Domain, error) {
	result, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err!=nil {
		return nil, err
	}
	config := &adapterstruct.Domain{}
	err = json.Unmarshal([]byte(result), config)
	if err!=nil {
		return nil, err
	}
	return config, nil
}

// GetStats returns the current statistics of the domain. The data is read directly from the vmm (e.g. qemu).
func (c *Controller) GetStats(ctx context.Context, id string) (*adapterstruct.DomainStats, error) {
	stats, err := c.adapter.GetStats(ctx, id)
	if err!=nil {
		return nil, err
	}
	return stats, nil
}

// Create creates a domain with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (c *Controller) Create(ctx context.Context, config *adapterstruct.Domain) (string, error) {
	uuid := uuid.New().String()
	err := c.SetConfig(ctx, uuid, config)
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	err = c.SetState(ctx, uuid, structure.DOMAIN_DOWN)
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	err = c.SetStatePushInterval(ctx, uuid, 20)
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	err = c.SetConfigPushInterval(ctx, uuid, 30)
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	err = c.SetAffinity(ctx, uuid, []string{})
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	err = c.SetNode(ctx, uuid, "")
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	return uuid, nil
}

// SetNode updates the node that is responsible for the domain.
func (c *Controller) SetNode(ctx context.Context, id string, node string) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id), node, 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetAffinity updates the affinity of the domain. The affinity specifies a list of affinity tags; any node
// tagged with at least one of those tags is considered eligible for hosting. An empty list means all nodes are
// eligible.
func (c *Controller) SetAffinity(ctx context.Context, id string, tags []string) error {
	rawTags, err := json.Marshal(affinity(tags))
	if err!=nil {
		return fmt.Errorf("cannot serialize affinity tags: %w", err)
	}
	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id), string(rawTags), 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetState updates the desired state of the domain.
func (c *Controller) SetState(ctx context.Context, id string, state structure.DOMAIN_STATE) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id), string(state), 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetStatePushInterval updates the state push interval of the domain.
func (c *Controller) SetStatePushInterval(ctx context.Context, id string, interval int) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATEPUSHINT/%s", id), strconv.Itoa(interval), 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetConfigPushInterval updates the config push interval of the domain.
func (c *Controller) SetConfigPushInterval(ctx context.Context, id string, interval int) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIGPUSHINT/%s", id), strconv.Itoa(interval), 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetConfig updates the domain configuration and updates associated metadata on the database. 
func (c *Controller) SetConfig(ctx context.Context, id string, config *adapterstruct.Domain) error {
	rawConfig, err := json.Marshal(config)
	if err!=nil {
		return fmt.Errorf("cannot serialize config: %w", err)
	}

	rawResources, err := json.Marshal(resources{
		AllocatedCpu: float64(config.ResourceConfig.VCPUs),
		AllocatedMemory: config.ResourceConfig.Memory,
	})
	if err!=nil {
		return fmt.Errorf("cannot serialize resource config: %w", err)
	}
	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/RESOURCES/%s", id), string(rawResources), 0)
	if err!=nil {
		return err
	}
	
	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id), string(rawConfig), 0)
	if err!=nil {
		return err
	}
	return nil
}

// Delete completely removes a domain and its associated metadata.
func (c *Controller) Delete(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIGPUSHINT/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATEPUSHINT/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/RESOURCES/%s", id))
	if err!=nil {
		return err
	}
	return nil
}
