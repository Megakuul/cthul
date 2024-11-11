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

package generator

import (
	"fmt"

	libvirtstruct "cthul.io/cthul/pkg/domain/libvirt/structure"
	cthulstruct "cthul.io/cthul/pkg/domain/structure"
)

// Explanation: A libvirt graphics device represents a host component that provides an interface for the
// video device on the guest. It is capable of reading data from the video device and writing data to the
// guest keyboard and mouse.

// generateGraphic generates a libvirt graphics device from the cthul video adapter.
func (l *LibvirtGenerator) generateGraphic(adapter *cthulstruct.VideoAdapter) (*libvirtstruct.Graphics, error) {
	graphics := &libvirtstruct.Graphics{
		Listen: &libvirtstruct.GraphicsListen{},
	}

	graphicDevice, err := l.wave.LookupGraphic(adapter.DeviceId)
	if err!=nil {
		return nil, err
	}

	switch graphicDevice.Type {
	case wave.GRAPHIC_SPICE:
		graphics.MetaType = libvirtstruct.GRAPHICS_SPICE
		graphics.Listen.MetaType = libvirtstruct.GRAPHICS_LISTEN_SOCKET
		graphics.Listen.MetaPath = graphicDevice.Path
	default:
		return nil, fmt.Errorf("unsupported device type: %s", graphicDevice.Type)
	}

	return graphics, nil
}
