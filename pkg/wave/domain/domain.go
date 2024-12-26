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
	"fmt"
	"strconv"
	"strings"

	adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/domain/structure"
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
	domainCPUs, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/CPU/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain cpu: %w", err)
	}
	domainMems, err := d.client.GetRange(ctx, "/WAVE/DOMAIN/MEM/")
	if err!=nil {
		return nil, fmt.Errorf("fetching domain memory: %w", err)
	}
	
	for key, node := range domainNodes {
		uuid := strings.TrimPrefix(key, "/WAVE/DOMAIN/NODE/")
		state := domainStates[key]
		cpu, err := strconv.ParseFloat(domainCPUs[key], 64)
		if err!=nil {
			return nil, fmt.Errorf("parsing domain cpu '%s': %w", cpu, err)
		}
		memory, err := strconv.ParseFloat(domainMems[key], 64)
		if err!=nil {
			return nil, fmt.Errorf("parsing domain memory '%s': %w", memory, err)
		}
		domains[uuid] = structure.Domain{
			Node: node,
			State: structure.DOMAIN_STATE(state),
			CPU: cpu,
			Memory: memory,
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

func (d *DomainController) Apply(ctx context.Context, id string, config *adapterstruct.Domain) error {
	_, err := d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id), string(state), 0)
	if err!=nil {
		return err
	}
	return nil
}

func (d *DomainController) Destroy(ctx context.Context, id string) error {
	_, err := d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id), string(state), 0)
	if err!=nil {
		return err
	}
	return nil
}

func (d *DomainController) ChangeState(ctx context.Context, id string, state structure.DOMAIN_STATE) error {
	_, err := d.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/STATE/%s", id), string(state), 0)
	if err!=nil {
		return err
	}
	return nil
}

func (d *DomainController) GetDomainStats(ctx context.Context, id string) (string, error) {
	// call domain adapter 
}
