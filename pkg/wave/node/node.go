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

package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/node/structure"
)

// NodeController provides an interface for wave node related operations.
// The single source of truth for configuration is always the register.toml file on the node.
// The single source of truth for dynamic data is always the database.
type NodeController struct {
	client db.Client
}

type NodeControllerOption func(*NodeController)

func NewNodeController(client db.Client, opts ...NodeControllerOption) *NodeController {
	controller := &NodeController{
		client: client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// affinity provides the node affinity tags structure stored in the database.
type affinity []string

// resources provides the node resource structure stored in the database.
type resources struct {
	AllocatedCpu float64 `json:"allocated_cpu"`
	AvailableCpu float64 `json:"available_cpu"`
	AllocatedMemory int64 `json:"allocated_memory"`
	AvailableMemory int64 `json:"available_memory"`
}

// List returns a map containing node uuids and associated metadata from the database.
func (n *NodeController) List(ctx context.Context) (map[string]structure.Node, error) {
	nodes := map[string]structure.Node{}
	
	nodeStates, err := n.client.GetRange(ctx, "/WAVE/NODE/STATE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node state: %w", err)
	}
	nodeAffinities, err := n.client.GetRange(ctx, "/WAVE/NODE/AFFINITY/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node affinity tags: %w", err)
	}
	nodeResources, err := n.client.GetRange(ctx, "/WAVE/NODE/RESOURCES/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node resources: %w", err)
	}
	
	for key, state := range nodeStates {
		var nodeErr error
		uuid := strings.TrimPrefix(key, "/WAVE/NODE/STATE/")
		
		affinity := &affinity{}
		err := json.Unmarshal([]byte(nodeAffinities[fmt.Sprint("/WAVE/NODE/AFFINITY/", uuid)]), affinity)
		if err!=nil {
			nodeErr = errors.Join(nodeErr, fmt.Errorf("parsing node affinity tags: %w", err))
		}
		
		resources := &resources{}
		err = json.Unmarshal([]byte(nodeResources[fmt.Sprint("/WAVE/NODE/RESOURCES/", uuid)]), resources)
		if err!=nil {
			nodeErr = errors.Join(nodeErr, fmt.Errorf("parsing node resources: %w", err))
		}
		
		nodes[uuid] = structure.Node{
			Affinity: *affinity,
			State: structure.NODE_STATE(state),
			AllocatedCpu: resources.AllocatedCpu,
			AvailableCpu: resources.AvailableCpu,
			AllocatedMemory: resources.AllocatedMemory,
			AvailableMemory: resources.AvailableMemory,
			Error: nodeErr,
		}
	}
	return nodes, nil
}

// Register registers / announces node information to the cluster by adding it to the database.
// The registration expires after ttl and must be renewed in order to ensure the node is part of the cluster.
func (n *NodeController) Register(ctx context.Context, id string, node structure.Node, ttl int64) error {
	_, err := n.client.Set(ctx, fmt.Sprintf("/WAVE/NODE/STATE/%s", id), string(node.State), ttl)
	if err!=nil {
		return err
	}

	rawTags, err := json.Marshal(affinity(node.Affinity))
	if err!=nil {
		return fmt.Errorf("cannot serialize affinity tags: %w", err)
	}
	_, err = n.client.Set(ctx, fmt.Sprintf("/WAVE/NODE/AFFINITY/%s", id), string(rawTags), ttl)
	if err!=nil {
		return err
	}

	rawResources, err := json.Marshal(&resources{
		AllocatedCpu: node.AllocatedCpu,
		AvailableCpu: node.AvailableCpu,
		AllocatedMemory: node.AllocatedMemory,
		AvailableMemory: node.AvailableMemory,
	})
	if err!=nil {
		return fmt.Errorf("cannot serialize resources: %w", err)
	}
	_, err = n.client.Set(ctx, fmt.Sprintf("/WAVE/NODE/RESOURCES/%s", id), string(rawResources), ttl)
	if err!=nil {
		return err
	}

	return nil
}


// Unregister removes an existing node registration entry.
func (n *NodeController) Unregister(ctx context.Context, id string, node string) error {
	err := n.client.Delete(ctx, fmt.Sprintf("/WAVE/NODE/STATE/%s", id))
	if err!=nil {
		return err
	}

	err = n.client.Delete(ctx, fmt.Sprintf("/WAVE/NODE/AFFINITY/%s", id))
	if err!=nil {
		return err
	}

	err = n.client.Delete(ctx, fmt.Sprintf("/WAVE/NODE/RESOURCES/%s", id))
	if err!=nil {
		return err
	}
	
	return nil
}
