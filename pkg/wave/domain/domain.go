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
	"errors"
	"fmt"
	"strings"

	"cthul.io/cthul/pkg/adapter/domain"
	domainstruct "cthul.io/cthul/pkg/api/wave/v1/domain"
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

// Controller provides an interface for wave domain related operations.
type Controller struct {
	node    string
	runRoot string
	client  db.Client
	adapter domain.Adapter
}

type Option func(*Controller)

func New(node string, client db.Client, adapter domain.Adapter, opts ...Option) *Controller {
	controller := &Controller{
		node:    node,
		runRoot: "/run/cthul/wave/",
		client:  client,
		adapter: adapter,
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

// List returns a map containing domain ids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]*domainstruct.Domain, error) {
	domains := map[string]*domainstruct.Domain{}

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

	for key, rawConfig := range configs {
		var domainErr error
		id := strings.TrimPrefix(key, "/WAVE/DOMAIN/CONFIG/")
		reqnode := reqnodes[fmt.Sprint("/WAVE/DOMAIN/REQNODE/", id)]
		node := nodes[fmt.Sprint("/WAVE/DOMAIN/NODE/", id)]

		config := &domainstruct.DomainConfig{}
		err = proto.Unmarshal([]byte(rawConfig), config)
		if err != nil {
			domainErr = errors.Join(domainErr, fmt.Errorf("parsing domain config: %w", err))
		}

    domain := &domainstruct.Domain{
			Reqnode:         reqnode,
			Node:            node,
			Config:          config,
		}
    if domainErr != nil {
      domain.Error = domainErr.Error()
    }
		domains[id] = domain
	}
	return domains, nil
}

// Apply upserts the domain configuration.
func (c *Controller) Apply(ctx context.Context, id string, config *domainstruct.DomainConfig) error {
	rawConfig, err := proto.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot serialize config: %w", err)
	}

	_, err = c.client.Set(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id), string(rawConfig), 0)
	if err != nil {
		return err
	}
	return nil
}

// Stat returns the current statistics of the domain. The data is read directly from the vmm (e.g. qemu).
func (c *Controller) Stat(ctx context.Context, id string) (*domainstruct.DomainStats, error) {
	node, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err != nil {
		return nil, err
	}
	if node != c.node {
		return nil, &NodeMismatchErr{Message: "domain must be on the same node as the controller", Node: c.node}
	}
	stats, err := c.adapter.GetStats(ctx, id)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// Lookup searches for the domain by id and returns its configuration.
func (c *Controller) Lookup(ctx context.Context, id string) (*domainstruct.Domain, error) {
	reqnode, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/REQNODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain reqnode: %w", err)
	}
	node, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/NODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain node: %w", err)
	}
	rawConfig, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching domain configs: %w", err)
	}

  if rawConfig == "" {
    return nil, fmt.Errorf("domain not found")
  }

	config := &domainstruct.DomainConfig{}
	err = proto.Unmarshal([]byte(rawConfig), config)
	if err != nil {
		return nil, fmt.Errorf("parsing domain config: %w", err)
	}

	return &domainstruct.Domain{
		Reqnode: reqnode,
		Node:    node,
		Config:  config,
		Error:   "",
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
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/DOMAIN/CONFIG/%s", id))
	if err != nil {
		return err
	}
	return nil
}
