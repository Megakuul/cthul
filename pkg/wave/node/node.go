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

package node

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cthul.io/cthul/pkg/api/wave/v1/node"
	"cthul.io/cthul/pkg/db"
	"google.golang.org/protobuf/proto"
)

// NodeMismatchErr indicates that the action cannot be executed on this node.
type NodeMismatchErr struct {
	Node    string
	Message string
}

func (n *NodeMismatchErr) Error() string {
	return n.Message
}

// Controller provides an interface for wave node related operations.
type Controller struct {
	node   string
	client db.Client
}

type Option func(*Controller)

func New(node string, client db.Client, opts ...Option) *Controller {
	controller := &Controller{
		node:   node,
		client: client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// List returns a map containing node uuids and associated metadata from the database.
func (n *Controller) List(ctx context.Context) (map[string]*node.Node, error) {
	nodes := map[string]*node.Node{}

	configs, err := n.client.GetRange(ctx, "/WAVE/NODE/CONFIG/")
	if err != nil {
		return nil, fmt.Errorf("fetching node config: %w", err)
	}

	for key, rawConfig := range configs {
		var nodeErr error
		id := strings.TrimPrefix(key, "/WAVE/NODE/CONFIG/")

		config := &node.NodeConfig{}
		err := proto.Unmarshal([]byte(rawConfig), config)
		if err != nil {
			nodeErr = errors.Join(nodeErr, fmt.Errorf("parsing node config %w", err))
		}

		node := &node.Node{
			Config: config,
		}
		if nodeErr != nil {
			node.Error = nodeErr.Error()
		}
		nodes[id] = node
	}
	return nodes, nil
}

// Lookup finds the specified node and returns its associated metadata from the database.
func (n *Controller) Lookup(ctx context.Context, id string) (*node.Node, error) {
	rawConfig, err := n.client.Get(ctx, fmt.Sprintf("/WAVE/NODE/CONFIG/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching node config: %w", err)
	}

  if rawConfig == "" {
    return nil, fmt.Errorf("node not found")
  }

	config := &node.NodeConfig{}
	err = proto.Unmarshal([]byte(rawConfig), config)
	if err != nil {
    return nil, fmt.Errorf("failed to parse node config %w", err)
	}

  return &node.Node{
    Config: config,
  }, nil
}

// Register registers / announces node information to the cluster by adding it to the database.
// The registration expires after ttl and must be renewed in order to ensure the node is part of the cluster.
func (n *Controller) Register(ctx context.Context, id string, node *node.Node, ttl int64) error {
  rawConfig, err := proto.Marshal(node.Config)
  if err!=nil {
    return fmt.Errorf("failed to serialize config: %w", err)
  }
	_, err = n.client.Set(ctx, fmt.Sprintf("/WAVE/NODE/CONFIG/%s", id), string(rawConfig), ttl)
	if err != nil {
		return err
	}

	return nil
}

// Unregister removes an existing node registration entry.
func (n *Controller) Unregister(ctx context.Context, id string) error {
	err := n.client.Delete(ctx, fmt.Sprintf("/WAVE/NODE/CONFIG/%s", id))
	if err != nil {
		return err
	}

	return nil
}
