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

// Explanation: A libvirt serial device is used to emulate a real serial device. The serial device does just
// transfer raw data from host to guest and vice versa. With the default isa-serial bus,
// it does fully operate on top of PIO (with PMIO calls from guest os and ISR interrupts triggered by qemu).
// Using virtio driver enables more efficient communication by transfering data over virtqueues instead of
// sending every chunk via separate PIO instruction.

// generateSerial generates a libvirt serial device from the cthul serial device and adapter.
func (l *LibvirtGenerator) generateSerial(device *cthulstruct.SerialDevice, adapter *cthulstruct.SerialAdapter) (*libvirtstruct.Serial, error) {
	serial := &libvirtstruct.Serial{
		MetaType: libvirtstruct.SERIAL_UNIX,
		Source: &libvirtstruct.SerialSource{MetaMode: libvirtstruct.SERIAL_SOURCE_BIND},
		Target: &libvirtstruct.SerialTarget{
			MetaPort: device.Port,
		},
	}

	serialDevice, err := l.wave.LookupSerial(adapter.DeviceId)
	if err!=nil {
		return nil, err
	}

	// Source (on host)
	serial.MetaType = libvirtstruct.SERIAL_UNIX
	serial.Source.MetaMode = libvirtstruct.SERIAL_SOURCE_BIND
	serial.Source.MetaPath = serialDevice.Path

	// Target (on guest)
	switch device.SerialBus {
	case cthulstruct.SERIAL_ISA:
		serial.Target.MetaType = libvirtstruct.SERIAL_BUS_ISA
	case cthulstruct.SERIAL_VIRTIO:
		serial.Target.MetaType = libvirtstruct.SERIAL_BUS_VIRTIO
	default:
		return nil, fmt.Errorf("unknown serial bus type: %s", device.SerialBus)
	}
	
	return serial, nil
}
