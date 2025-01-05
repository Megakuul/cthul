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
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"cthul.io/cthul/pkg/adapter/video"
	adapterstruct "cthul.io/cthul/pkg/adapter/video/structure"
	"cthul.io/cthul/pkg/db"
	"cthul.io/cthul/pkg/wave/video/structure"
	"github.com/google/uuid"
)

// Controller provides an interface for wave video device operations.
type Controller struct {
	client db.Client
}

type ControllerOption func(*Controller)

func NewController(client db.Client, opts ...ControllerOption) *Controller {
	controller := &Controller{
		client: client,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// List returns a map containing video device uuids and associated metadata from the database.
func (c *Controller) List(ctx context.Context) (map[string]structure.Video, error) {
	videos := map[string]structure.Video{}
	
	videoNodes, err := c.client.GetRange(ctx, "/WAVE/VIDEO/NODE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching video device node: %w", err)
	}
	videoTypes, err := c.client.GetRange(ctx, "/WAVE/VIDEO/TYPE/")
	if err!=nil {
		return nil, fmt.Errorf("fetching video device type: %w", err)
	}
	videoPaths, err := c.client.GetRange(ctx, "/WAVE/VIDEO/PATH/")
	if err!=nil {
		return nil, fmt.Errorf("fetching video device path: %w", err)
	}
	
	for key, node := range videoNodes {
		uuid := strings.TrimPrefix(key, "/WAVE/VIDEO/NODE/")
		typ := videoTypes[fmt.Sprint("/WAVE/VIDEO/TYPE/", uuid)]
		path := videoPaths[fmt.Sprint("/WAVE/VIDEO/PATH", uuid)]
		
		videos[uuid] = structure.Video{
			Node: node,
			Type: structure.VIDEO_TYPE(typ),
			Path: path,
		}
	}
	return videos, nil
}

// Create creates a video device with the specified configuration and default metadata values.
// If the creation fails, the function tries to remove already created resources from the database.
func (c *Controller) Create(ctx context.Context, id string) (string, error) {
	uuid := uuid.New().String()
	err := c.SetType(ctx, uuid, structure.VIDEO_SPICE)
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	err = c.SetNode(ctx, uuid, "")
	if err!=nil {
		return "", errors.Join(err, c.Delete(ctx, uuid))
	}
	return uuid, nil
}

// SetType updates the type of the video device and .
func (c *Controller) SetType(ctx context.Context, id string, typ structure.VIDEO_TYPE) error {
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/TYPE/%s", id), string(typ), 0)
	if err!=nil {
		return err
	}
	return nil
}

func (c *Controller) Attach(ctx context.Context, id string, node string) error {
	cooking, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/COOKING/%s", id))
	
	cooking, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/COOKING/%s", id), "yes", 5)
	if err!=nil {
		return err
	}
	if cooking!="" {
		return fmt.Errorf("someone is cooking")
	}
	
	deviceNode, err := c.client.Get(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
	if err!=nil {
		return err
	}
	if deviceNode == node {
		return nil
	}

	go func() {
		err = c.client.Watch(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id), func(_, _ string, err error) {
			
		})
	}()
	_, err := c.client.Set(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id), node, 0)
	if err!=nil {
		return err
	}


}


// Delete completely removes a video device and its associated metadata.
func (c *Controller) Delete(ctx context.Context, id string) error {
	err := c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/NODE/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/REQNODE/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/TYPE/%s", id))
	if err!=nil {
		return err
	}
	err = c.client.Delete(ctx, fmt.Sprintf("/WAVE/VIDEO/PATH/%s", id))
	if err!=nil {
		return err
	}
	return nil
}
