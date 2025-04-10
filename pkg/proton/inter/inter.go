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

package inter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/proton/inter/structure"
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

// Controller provides an interface for wave interface device operations.
type Controller struct {
	node    string
	runRoot string
	client  db.Client
}

type Option func(*Controller)

func New(node string, client db.Client, opts ...Option) *Controller {
	controller := &Controller{
		node:    node,
		runRoot: "/run/cthul/proton/",
		client:  client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// WithRunRoot defines a custom root for runtime files (bsd sockets etc.).
// The controller needs this information to understand where to find those files (usually created by operators).
func WithRunRoot(device string) Option {
	return func(c *Controller) {
		c.runRoot = device
	}
}

// List returns a map containing interface device uuids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]structure.Inter, error) {
	inters := map[string]structure.Inter{}

	reqnodes, err := c.client.GetRange(ctx, "/PROTON/INTER/REQNODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching interface device reqnode: %w", err)
	}
	nodes, err := c.client.GetRange(ctx, "/PROTON/INTER/NODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching interface device node: %w", err)
	}
	types, err := c.client.GetRange(ctx, "/PROTON/INTER/TYPE/")
	if err != nil {
		return nil, fmt.Errorf("fetching interface type: %w", err)
	}
	devices, err := c.client.GetRange(ctx, "/PROTON/INTER/DEVICE/")
	if err != nil {
		return nil, fmt.Errorf("fetching interface device device: %w", err)
	}

	for key, device := range devices {
		var interErr error
		id := strings.TrimPrefix(key, "/PROTON/INTER/DEVICE/")
		reqnode := reqnodes[fmt.Sprint("/PROTON/INTER/REQNODE/", id)]
		node := nodes[fmt.Sprint("/PROTON/INTER/NODE/", id)]
		typ := types[fmt.Sprint("/PROTON/INTER/TYPE/", id)]

		inters[id] = structure.Inter{
			Reqnode:  reqnode,
			Node:     node,
			Type:     structure.INTER_TYPE(typ),
			Device:     device,
			Error:    interErr,
		}
	}
	return inters, nil
}

// Create creates a interface device with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (c *Controller) Create(ctx context.Context, device string, typ structure.INTER_TYPE) (string, error) {
	id := uuid.New().String()
	err := c.SetDevice(ctx, id, device)
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
  _, err = c.client.Set(ctx, fmt.Sprintf("/PROTON/INTER/TYPE/%s", id), string(typ), 0)
  if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, id))
  }
	return id, nil
}

func (c *Controller) SetDevice(ctx context.Context, id, device string) error {
	if device == "" {
		return fmt.Errorf("device must be non-empty because it is the device core-property")
	}
	_, err := c.client.Set(ctx, fmt.Sprintf("/PROTON/INTER/DEVICE/%s", id), device, 0)
	if err != nil {
		return err
	}
	return nil
}

// Lookup searches for the device by id and returns its configuration.
func (c *Controller) Lookup(ctx context.Context, id string) (*structure.Inter, error) {
	reqnode, err := c.client.Get(ctx, fmt.Sprintf("/PROTON/INTER/REQNODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching interface device reqnode: %w", err)
	}
	node, err := c.client.Get(ctx, fmt.Sprintf("/PROTON/INTER/NODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching interface device node: %w", err)
	}
	typ, err := c.client.Get(ctx, fmt.Sprintf("/PROTON/INTER/TYPE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching interface type: %w", err)
	}
	device, err := c.client.Get(ctx, fmt.Sprintf("/PROTON/INTER/DEVICE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching interface device device: %w", err)
	}

	if device == "" {
		return nil, fmt.Errorf("device not found")
	}

	return &structure.Inter{
		Reqnode:  reqnode,
		Node:     node,
		Type:     structure.INTER_TYPE(typ),
		Device:     device,
	}, nil
}

// Attach requests the device to be relocated to the specified node and waits until it's ready (if wait flag is set).
func (c *Controller) Attach(ctx context.Context, id, node string, wait bool) error {
	if !wait {
		_, err := c.client.Set(ctx, fmt.Sprintf("/PROTON/INTER/REQNODE/%s", id), node, 0)
		if err != nil {
			return err
		}
	}

	pollCtx, pollCtxCancel := context.WithCancel(ctx)
	pollG, pollGCtx := errgroup.WithContext(pollCtx)

	pollG.Go(func() error {
		err := c.client.Watch(pollGCtx, fmt.Sprintf("/PROTON/INTER/NODE/%s", id), func(_, activeNode string, err error) {
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
		activeNode, err := c.client.Get(pollGCtx, fmt.Sprintf("/PROTON/INTER/NODE/%s", id))
		if err != nil {
			return err
		}
		if node == activeNode {
			pollCtxCancel()
		}
		return nil
	})

	_, err := c.client.Set(ctx, fmt.Sprintf("/PROTON/INTER/REQNODE/%s", id), node, 0)
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
	err := c.client.Delete(ctx, fmt.Sprintf("/PROTON/INTER/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/PROTON/INTER/REQNODE/%s", id))
	if err != nil {
		return nil
	}
	return nil
}

// Delete completely removes a interface device and its associated metadata.
func (c *Controller) Delete(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/PROTON/INTER/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/PROTON/INTER/REQNODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/PROTON/INTER/TYPE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/PROTON/INTER/DEVICE/%s", id))
	if err != nil {
		return err
	}
	return nil
}
