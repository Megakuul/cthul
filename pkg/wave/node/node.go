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
	"fmt"
	"strconv"
	"strings"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/node/structure"
)

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

// List returns a map containing node uuids and associated metadata from the database.
func (d *NodeController) List(ctx context.Context) (map[string]structure.Node, error) {
	nodes := map[string]structure.Node{}
	
	nodeNodes, err := d.client.GetRange(ctx, "/WAVE/NODE/NODE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node node: %w", err)
	}
	nodeAffinity, err := d.client.GetRange(ctx, "/WAVE/NODE/AFFINITY/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node affinity tags: %w", err)
	}
	nodeStates, err := d.client.GetRange(ctx, "/WAVE/NODE/STATE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node state: %w", err)
	}
	nodeCPUs, err := d.client.GetRange(ctx, "/WAVE/NODE/CPU/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node cpu: %w", err)
	}
	nodeMems, err := d.client.GetRange(ctx, "/WAVE/NODE/MEM/")
	if err!=nil {
		return nil, fmt.Errorf("fetching node memory: %w", err)
	}
	
	for key, node := range nodeNodes {
		uuid := strings.TrimPrefix(key, "/WAVE/NODE/NODE/")
		state := nodeStates[key]
		cpu, err := strconv.ParseFloat(nodeCPUs[key], 64)
		if err!=nil {
			return nil, fmt.Errorf("parsing node cpu '%s': %w", cpu, err)
		}
		memory, err := strconv.ParseFloat(nodeMems[key], 64)
		if err!=nil {
			return nil, fmt.Errorf("parsing node memory '%s': %w", memory, err)
		}
		nodes[uuid] = structure.Node{
			Node: node,
			State: structure.NODE_STATE(state),
			CPU: cpu,
			Memory: memory,
		}
	}
	return nodes, nil
}
