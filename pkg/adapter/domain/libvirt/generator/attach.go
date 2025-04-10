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

	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
)

// Attach installs / locks all devices that are required by the domain config.
func (g *Generator) Attach(ctx context.Context, config *cthulstruct.Domain) error {
	for _, device := range config.VideoAdapters {
    err := g.video.Attach(ctx, device.DeviceId, g.nodeId, true)
		if err!=nil {
			return err
		}
	}

	for _, device := range config.SerialDevices {
    err := g.serial.Attach(ctx, device.DeviceId, g.nodeId, true)
    if err!=nil {
      return err
    }
	}
	
	for _, device := range config.StorageDevices {
    err := g.disk.Attach(ctx, device.DeviceId, g.nodeId, true)
    if err!=nil {
      return err
    }
	}

	for _, device := range config.NetworkDevices {
    err := g.inter.Attach(ctx, device.DeviceId, g.nodeId, true)
    if err!=nil {
      return err
    }
	}
	
	return fmt.Errorf("not implemented biatch")
}
