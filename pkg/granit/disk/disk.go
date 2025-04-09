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
	"strconv"
	"strings"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/granit/disk/structure"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// TODO: maybe add more functions + implement operator (currently just for testing)

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
	node    string
	runRoot string
	client  db.Client
}

type Option func(*Controller)

func New(node string, client db.Client, opts ...Option) *Controller {
	controller := &Controller{
		node:    node,
		runRoot: "/run/cthul/granit/",
		client:  client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// WithRunRoot defines a custom root for runtime files (bsd sockets etc.).
// The controller needs this information to understand where to find those files (usually created by operators).
func WithRunRoot(path string) Option {
	return func(c *Controller) {
		c.runRoot = path
	}
}

// List returns a map containing disk device uuids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]structure.Disk, error) {
	disks := map[string]structure.Disk{}

	reqnodes, err := c.client.GetRange(ctx, "/GRANIT/DISK/REQNODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device reqnode: %w", err)
	}
	nodes, err := c.client.GetRange(ctx, "/GRANIT/DISK/NODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device node: %w", err)
	}
	types, err := c.client.GetRange(ctx, "/GRANIT/DISK/TYPE/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk type: %w", err)
	}
	formats, err := c.client.GetRange(ctx, "/GRANIT/DISK/FORMAT/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk format: %w", err)
	}
	paths, err := c.client.GetRange(ctx, "/GRANIT/DISK/PATH/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk device path: %w", err)
	}
	readonlys, err := c.client.GetRange(ctx, "/GRANIT/DISK/READONLY/")
	if err != nil {
		return nil, fmt.Errorf("fetching disk readonly status: %w", err)
	}

	for key, path := range paths {
		var diskErr error
		id := strings.TrimPrefix(key, "/GRANIT/DISK/PATH/")
		reqnode := reqnodes[fmt.Sprint("/GRANIT/DISK/REQNODE/", id)]
		node := nodes[fmt.Sprint("/GRANIT/DISK/NODE/", id)]
		typ := types[fmt.Sprint("/GRANIT/DISK/TYPE/", id)]
		format := formats[fmt.Sprint("/GRANIT/DISK/FORMAT/", id)]

		readonly, err := strconv.ParseBool(readonlys[fmt.Sprint("/GRANIT/DISK/READONLY/", id)])
		if err != nil {
			diskErr = errors.Join(diskErr, fmt.Errorf("parsing disk readonly state: %w", err))
		}

		disks[id] = structure.Disk{
			Reqnode:  reqnode,
			Node:     node,
			Type:     structure.DISK_TYPE(typ),
			Format:   structure.DISK_FORMAT(format),
			Path:     path,
			Readonly: readonly,
			Error:    diskErr,
		}
	}
	return disks, nil
}

// Create creates a disk device with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (c *Controller) Create(ctx context.Context, path string, typ structure.DISK_TYPE, format structure.DISK_FORMAT, readonly bool) (string, error) {
	id := uuid.New().String()
	err := c.SetPath(ctx, id, path)
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
  _, err = c.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/TYPE/%s", id), string(typ), 0)
  if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, id))
  }
  _, err = c.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/FORMAT/%s", id), string(format), 0)
  if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, id))
  }
  _, err = c.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/READONLY/%s", id), strconv.FormatBool(readonly), 0)
  if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, id))
  }
	return id, nil
}

func (c *Controller) SetPath(ctx context.Context, id, path string) error {
	if path == "" {
		return fmt.Errorf("path must be non-empty because it is the device core-property")
	}
	_, err := c.client.Set(ctx, fmt.Sprintf("/GRANIT/DISK/PATH/%s", id), path, 0)
	if err != nil {
		return err
	}
	return nil
}

// Lookup searches for the device by id and returns its configuration.
func (c *Controller) Lookup(ctx context.Context, id string) (*structure.Disk, error) {
	reqnode, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/REQNODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device reqnode: %w", err)
	}
	node, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/NODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device node: %w", err)
	}
	typ, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/TYPE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk type: %w", err)
	}
	format, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/FORMAT/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk format: %w", err)
	}
	path, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/PATH/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk device path: %w", err)
	}
	diskReadonly, err := c.client.Get(ctx, fmt.Sprintf("/GRANIT/DISK/READONLY/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching disk readonly status: %w", err)
	}

	readonly, err := strconv.ParseBool(diskReadonly)
	if err != nil {
		return nil, fmt.Errorf("parsing disk readonly state: %w", err)
	}

	if path == "" {
		return nil, fmt.Errorf("device not found")
	}

	return &structure.Disk{
		Reqnode:  reqnode,
		Node:     node,
		Type:     structure.DISK_TYPE(typ),
		Format:   structure.DISK_FORMAT(format),
		Path:     path,
		Readonly: readonly,
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
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/READONLY/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/TYPE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/FORMAT/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/GRANIT/DISK/PATH/%s", id))
	if err != nil {
		return err
	}
	return nil
}
