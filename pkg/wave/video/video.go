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

package video

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/video/structure"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// NodeMismatchErr indicates that the action cannot be executed on this node.
type NodeMismatchErr struct {
  Node string
  Message string
}

func (n *NodeMismatchErr) Error() string {
  return n.Message
}

// Controller provides an interface for wave video device operations.
type Controller struct {
  node string
  runRoot string
	client db.Client
}

type Option func(*Controller)

func New(node string, client db.Client, opts ...Option) *Controller {
	controller := &Controller{
    node: node,
    runRoot: "/run/cthul/wave/",
		client: client,
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

// List returns a map containing video device uuids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]structure.Video, error) {
	videos := map[string]structure.Video{}

	reqnodes, err := c.client.GetRange(ctx, "/WAVE/VIDEO/REQNODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching video device reqnode: %w", err)
	}
	nodes, err := c.client.GetRange(ctx, "/WAVE/VIDEO/NODE/")
	if err != nil {
		return nil, fmt.Errorf("fetching video device node: %w", err)
	}
	paths, err := c.client.GetRange(ctx, "/WAVE/VIDEO/PATH/")
	if err != nil {
		return nil, fmt.Errorf("fetching video device path: %w", err)
	}
	types, err := c.client.GetRange(ctx, "/WAVE/VIDEO/TYPE/")
	if err != nil {
		return nil, fmt.Errorf("fetching video device type: %w", err)
	}

	for key, path := range paths {
		uuid := strings.TrimPrefix(key, "/WAVE/VIDEO/PATH/")
		reqnode := reqnodes[fmt.Sprint("/WAVE/VIDEO/REQNODE/", uuid)]
		node := nodes[fmt.Sprint("/WAVE/VIDEO/NODE/", uuid)]
		typ := types[fmt.Sprint("/WAVE/VIDEO/TYPE/", uuid)]

		videos[uuid] = structure.Video{
			Reqnode: reqnode,
			Node:    node,
			Type:    structure.VIDEO_TYPE(typ),
			Path:    path,
		}
	}
	return videos, nil
}

// Create creates a video device with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (c *Controller) Create(ctx context.Context, typ structure.VIDEO_TYPE, path string) (string, error) {
	id := uuid.New().String()
	err := c.SetType(ctx, id, typ)
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
	err = c.SetPath(ctx, id, path)
	if err != nil {
		return "", errors.Join(err, c.Delete(ctx, id))
	}
	return id, nil
}

func (c *Controller) SetType(ctx context.Context, id string, typ structure.VIDEO_TYPE) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/TYPE/%s", id), string(typ), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) SetPath(ctx context.Context, id, path string) error {
	if path == "" {
		return fmt.Errorf("path must be non-empty because it is the device core-property")
	}
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/PATH/%s", id), path, 0)
	if err != nil {
		return err
	}
	return nil
}

// Connect creates a bidirectional communication bridge to the video device socket.
// Input and output is not manipulated, the format depends on the device type. Runs until the context is cancelled.
func (c *Controller) Connect(ctx context.Context, id string, reader chan<-[]byte, writer <-chan []byte) error {
  node, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
  if err!=nil {
    return err
  }
  if node != c.node {
    return &NodeMismatchErr{Message: "device must be on the same node as the controller", Node: c.node}
  }
  path, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/PATH/%s", id))
  if err!=nil {
    return err
  }

  path = filepath.Join(c.runRoot, path)
  if !strings.HasPrefix(path, c.runRoot) {
    return fmt.Errorf("device socket path escapes the run root '%s'", c.runRoot)
  }
  conn, err := net.Dial("unix", path) 
  if err!=nil {
    return err
  }
  defer conn.Close()

  loopG, loopGCtx := errgroup.WithContext(ctx)

  loopG.Go(func() error {
    buffer := make([]byte, 1024)
    for {
        err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
        if err!=nil {
        return err
        }
      n, err := conn.Read(buffer)
      if err!=nil {
          if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
          return err
        }
      }
      if n < 1 {
        continue
      }
      select {
      case reader<-buffer[:n]: 
      case <-loopGCtx.Done():
        return nil
      }
    }
  })
  loopG.Go(func() error {
    for {
      select {
      case data, ok := <-writer: 
        if !ok {
          return fmt.Errorf("writer got closed")
        } 
        err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
        if err!=nil {
          return err
        }
        _, err = conn.Write(data)
        if err!=nil {
          return err
        }
      
      case <-loopGCtx.Done():
        return nil
      }
    }
  })

  return loopG.Wait()
}

// Lookup searches for the device by id and returns its configuration.
func (c *Controller) Lookup(ctx context.Context, id string) (*structure.Video, error) {
	reqnode, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching video device reqnode: %w", err)
	}
	node, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching video device node: %w", err)
	}
	typ, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/TYPE/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching video device type: %w", err)
	}
	path, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/PATH/%s", id))
	if err != nil {
		return nil, fmt.Errorf("fetching video device path: %w", err)
	}

	if path == "" {
		return nil, fmt.Errorf("device not found")
	}

	return &structure.Video{
		Reqnode: reqnode,
		Node:    node,
		Type:    structure.VIDEO_TYPE(typ),
		Path:    path,
	}, nil
}

// Attach requests the device to be relocated to the specified node and waits until it's ready (if wait flag is set).
func (c *Controller) Attach(ctx context.Context, id, node string, wait bool) error {
	if !wait {
		_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id), node, 0)
		if err != nil {
			return err
		}
	}

	pollCtx, pollCtxCancel := context.WithCancel(ctx)
	pollG, pollGCtx := errgroup.WithContext(pollCtx)

	pollG.Go(func() error {
		err := c.client.Watch(pollGCtx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id), func(_, activeNode string, err error) {
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
		activeNode, err := c.client.Get(pollGCtx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
		if err != nil {
			return err
		}
		if node == activeNode {
			pollCtxCancel()
		}
		return nil
	})

	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id), node, 0)
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
	err := c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id))
	if err != nil {
		return nil
	}
	return nil
}

// Delete completely removes a video device and its associated metadata.
func (c *Controller) Delete(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/TYPE/%s", id))
	if err != nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/PATH/%s", id))
	if err != nil {
		return err
	}
	return nil
}
