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

package disk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cthul.io/cthul/pkg/api/granit/v1/disk"
	"cthul.io/cthul/pkg/db"
	"golang.org/x/sync/errgroup"
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

// Controller provides an interface for wave disk device operations.
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

// List returns a map containing disk device uuids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]*disk.Disk, error) {
	disks := map[string]*disk.Disk{}

	reqnodes, err := c.client.GetRange(ctx, "/GRANIT/DISK/REQNODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device reqnode: %w", err)
	}
	nodes, err := c.client.GetRange(ctx, "/GRANIT/DISK/NODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device node: %w", err)
	}
	clusters, err := c.client.GetRange(ctx, "/GRANIT/DISK/CLUSTER/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device cluster: %w", err)
	}
	configs, err := c.client.GetRange(ctx, "/GRANIT/DISK/CONFIG/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device config: %w", err)
	}

	for key, rawConfig := range configs {
		var diskErr error
		id := strings.TrimPrefix(key, "/GRANIT/DISK/CONFIG/")
		rawCluster := clusters[fmt.Sprint("/GRANIT/DISK/CLUSTER/", id)]
		reqnode := reqnodes[fmt.Sprint("/GRANIT/DISK/REQNODE/", id)]
		node := nodes[fmt.Sprint("/GRANIT/DISK/NODE/", id)]

		config := &disk.DiskConfig{}
		err = proto.Unmarshal([]byte(rawConfig), config)
		if err != nil {
			diskErr = errors.Join(diskErr, fmt.Errorf("parsing device config: %w", err))
		}

		cluster := &disk.DiskCluster{}
		err = proto.Unmarshal([]byte(rawCluster), cluster)
		if err != nil {
			diskErr = errors.Join(diskErr, fmt.Errorf("parsing device cluster: %w", err))
		}

		disk := &disk.Disk{
			Reqnode: reqnode,
			Node:    node,
			Cluster: cluster,
			Config:  config,
		}
		if diskErr != nil {
			disk.Error = diskErr.Error()
		}
		disks[id] = disk
	}
	return disks, nil
}

// Lookup searches for the device by id and returns its configuration.
func (c *Controller) Lookup(ctx context.Context, id string) (*disk.Disk, error) {
	reqnode, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/REQNODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device reqnode: %w", err)
	}
	node, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device node: %w", err)
	}
	rawCluster, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/CLUSTER/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device cluster: %w", err)
	}
	rawConfig, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/CONFIG/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device config: %w", err)
	}

	if rawConfig == "" {
		return nil, fmt.Errorf("device not found")
	}

	cluster := &disk.DiskCluster{}
	err = proto.Unmarshal([]byte(rawCluster), cluster)
	if err != nil {
		return nil, fmt.Errorf("parsing device cluster: %w", err)
	}

	config := &disk.DiskConfig{}
	err = proto.Unmarshal([]byte(rawConfig), config)
	if err != nil {
		return nil, fmt.Errorf("parsing device config: %w", err)
	}

	return &disk.Disk{
		Reqnode: reqnode,
		Node:    node,
		Config:  config,
    Cluster: cluster,
		Error:   "",
	}, nil
}

// Attach requests the device to be relocated to the specified node and waits until it's ready (if wait flag is set).
func (c *Controller) Attach(ctx context.Context, id, node string, wait bool) error {
	if !wait {
		_, err := c.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/REQNODE/%s", id), node, 0)
		if err != nil {
			return err
		}
	}

	pollCtx, pollCtxCancel := context.WithCancel(ctx)
	pollG, pollGCtx := errgroup.WithContext(pollCtx)

	pollG.Go(func() error {
		err := c.client.Watch(pollGCtx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id), func(_, activeNode string, err error) {
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
		activeNode, err := c.client.Get(pollGCtx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id))
		if err != nil {
			return err
		}
		if node == activeNode {
			pollCtxCancel()
		}
		return nil
	})

	_, err := c.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/REQNODE/%s", id), node, 0)
	if err != nil {
		return err
	}

	err = pollG.Wait()
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("context exceeded: device couldn't be attached in the provided context window")
	default:
		return nil
	}
}

// Detach removes the device from the current node. It doesn't wait until the node fully detached it.
func (c *Controller) Detach(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/REQNODE/%s", id))
	if err != nil {
		return nil
	}
	return nil
}

// Delete completely removes a disk device and its associated metadata.
func (c *Controller) Delete(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/REQNODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/CLUSTER/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/CONFIG/%s", id))
	if err != nil {
		return err
	}
	return nil
}
