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

	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
  "cthul.io/cthul/pkg/api/wave/v1/domain"
)

// Explanation: A libvirt serial device is used to emulate a real serial device. The serial device does just
// transfer raw data from host to guest and vice versa. With the default isa-serial bus,
// it does fully operate on top of PIO (with PMIO calls from guest os and ISR interrupts triggered by qemu).
// Using virtio driver enables more efficient communication by transfering data over virtqueues instead of
// sending every chunk via separate PIO instruction.

// generateSerial generates a libvirt serial device from the cthul serial device.
func (g *Generator) generateSerial(ctx context.Context, device *domain.SerialDevice) (*structure.Serial, error) {
	serial := &structure.Serial{
		MetaType: structure.SERIAL_UNIX,
		Source:   &structure.SerialSource{},
		Target: &structure.SerialTarget{
			MetaPort: device.Port,
		},
	}

	// Source (on host)
	serial.MetaType = structure.SERIAL_UNIX
	serial.Source.MetaMode = structure.SERIAL_SOURCE_BIND
	serial.Source.MetaPath = filepath.Join(g.waveRoot, "serial", device.DeviceId)

	// Target (on guest)
	switch device.SerialBus {
	case domain.SerialBus_SERIAL_BUS_ISA:
		serial.Target.MetaType = structure.SERIAL_BUS_ISA
	case domain.SerialBus_SERIAL_BUS_VIRTIO:
		serial.Target.MetaType = structure.SERIAL_BUS_VIRTIO
	default:
		return nil, fmt.Errorf("unknown serial bus type: %s", device.SerialBus)
	}

	return serial, nil
}
