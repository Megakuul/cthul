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

package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/domain/structure"
	"github.com/google/uuid"
)

type DomainController struct {
	client db.Client
}

type DomainControllerOption func(*DomainController)

func NewDomainController(client db.Client, opts ...DomainControllerOption) *DomainController {
	controller := &DomainController{
		client: client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// List returns a map containing domain uuids and associated metadata from the database.
func (d *DomainController) List(ctx context.Context) (map[string]structure.Domain, error) {
	domains := map[string]structure.Domain{}
	
	domainNodes, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/NODE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain node: %w", err)
	}
	domainStates, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/STATE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain state: %w", err)
	}
	domainAffinities, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/AFFINITY/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain affinity: %w", err)
	}
	domainCPUs, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/ALCPU/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain cpu: %w", err)
	}
	domainMems, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/ALMEM/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain memory: %w", err)
	}
	
	for key, node := range domainNodes {
		uuid := strings.TrimPrefix(key, "/WAVE/DOMAIN/NODE/")
		state := domainStates[fmt.Sprint("/WAVE/DOMAIN/STATE/", key)]
		affinity := strings.Split(domainAffinities[fmt.Sprint("/WAVE/DOMAIN/AFFINITY/", key)],  "|")
		
		cpu, err := strconv.ParseFloat(domainCPUs[fmt.Sprint("/WAVE/DOMAIN/ALCPU/", key)], 64)
		if err!=nil {
			return nil, fmt.Errorf("parsing allocated domain cpu '%.2f': %w", cpu, err)
		}
		memory, err := strconv.ParseInt(domainMems[fmt.Sprint("/WAVE/DOMAIN/ALMEM/", key)], 10, 64)
		if err!=nil {
			return nil, fmt.Errorf("parsing allocated domain memory '%d': %w", memory, err)
		}
		
		domains[uuid] = structure.Domain{
			Node: node,
			Affinity: affinity,
			State: structure.DOMAIN_STATE(state),
			AllocatedCPU: cpu,
			AllocatedMemory: memory,
		}
	}
	return domains, nil
}

// GetConfig returns the deserialized configuration of one node.
func (d *DomainController) GetConfig(ctx context.Context, id string) (*adapterstruct.Domain, error) {
	result, err := d.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
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

func (d *DomainController) GetDomainStats(ctx context.Context, id string) (string, error) {
	// call domain adapter 
}

// Create creates a domain with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (d *DomainController) Create(ctx context.Context, config *adapterstruct.Domain) (string, error) {
	uuid := uuid.New()
	err := d.SetConfig(ctx, uuid, config)
	if err!=nil {
		return "", errors.Join(err, d.Delete(ctx, uuid))
	}
	err = d.SetState(ctx, uuid, structure.DOMAIN_DOWN)
	if err!=nil {
		return "", errors.Join(err, d.Delete(ctx, uuid))
	}
	err = d.SetAffinity(ctx, uuid, []string{})
	if err!=nil {
		return "", errors.Join(err, d.Delete(ctx, uuid))
	}
	err = d.SetNode(ctx, uuid, "")
	if err!=nil {
		return "", errors.Join(err, d.Delete(ctx, uuid))
	}
	return uuid, nil
}

// SetNode updates the node that is responsible for the domain.
func (d *DomainController) SetNode(ctx context.Context, id string, node string) error {
	_, err := d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id), node, 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetAffinity updates the affinity of the domain. The affinity specifies a list of affinity tags; any node
// tagged with at least one of those tags is considered eligible for hosting. An empty list means all nodes are
// eligible.
func (d *DomainController) SetAffinity(ctx context.Context, id string, affinity []string) error {
	_, err := d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id), strings.Join(affinity, "|"), 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetState updates the desired state of the domain.
func (d *DomainController) SetState(ctx context.Context, id string, state structure.DOMAIN_STATE) error {
	_, err := d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id), string(state), 0)
	if err!=nil {
		return err
	}
	return nil
}

// SetConfig updates the domain configuration and updates associated metadata on the database. 
func (d *DomainController) SetConfig(ctx context.Context, id string, config *adapterstruct.Domain) error {
	rawConfig, err := json.Marshal(config)
	if err!=nil {
		return fmt.Errorf("cannot parse config: %w", err)
	}
	
	_, err = d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/ALCPU/%s", id),
		strconv.Itoa(int(config.ResourceConfig.VCPUs)), 0,
	)
	if err!=nil {
		return err
	}
	_, err = d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/ALMEM/%s", id),
		strconv.Itoa(int(config.ResourceConfig.Memory)), 0,
	)
	if err!=nil {
		return err
	}
	
	_, err = d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id), string(rawConfig), 0)
	if err!=nil {
		return err
	}
	return nil
}

// Delete completely removes a domain and its associated metadata.
func (d *DomainController) Delete(ctx context.Context, id string) error {
	err := d.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err!=nil {
		return err
	}
	err = d.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/AFFINITY/%s", id))
	if err!=nil {
		return err
	}
	err = d.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id))
	if err!=nil {
		return err
	}
	err = d.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/ALCPU/%s", id))
	if err!=nil {
		return err
	}
	err = d.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/ALMEM/%s", id))
	if err!=nil {
		return err
	}
	err = d.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err!=nil {
		return err
	}
	return nil
}
