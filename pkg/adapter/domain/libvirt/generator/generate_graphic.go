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

package generator

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	"cthul.io/cthul/pkg/api/wave/v1/domain"
)

// Explanation: A libvirt graphics device represents a host component that provides an interface for the
// video device on the guest. It is capable of reading data from the video device and writing data to the
// guest keyboard and mouse.

// generateGraphic generates a libvirt graphics device from the cthul video adapter.
func (g *Generator) generateGraphic(ctx context.Context, adapter *domain.VideoAdapter) (*structure.Graphics, error) {
	graphics := &structure.Graphics{
		Listen: &structure.GraphicsListen{},
	}

	graphicDevice, err := g.video.Lookup(ctx, adapter.DeviceId)
	if err != nil {
		return nil, err
	}

	graphics.MetaType = structure.GRAPHICS_SPICE
	graphics.Listen.MetaType = structure.GRAPHICS_LISTEN_SOCKET
	path := filepath.Join(g.waveRoot, graphicDevice.Config.Path)
	if !strings.HasPrefix(filepath.Clean(path), g.waveRoot) {
		return nil, fmt.Errorf("video device uses a socket path that escapes the run root '%s'", g.waveRoot)
	}
	graphics.Listen.MetaPath = path

	return graphics, nil
}
