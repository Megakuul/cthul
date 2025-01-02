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
	"fmt"

	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
)

// Explanation: A libvirt input device can represent multiple human interface devices, like a mouse, tablet or
// keyboard. With a basic bus like PS2 data is transferred via PIO, the host device (qemu) performs ISR
// interrupts to send data to the guest os.
// Using the VIRTIO bus leverages virtqueues to transfer data between the guest and host os.

// generateInterface generates a libvirt network interface device from the cthul network device.
func (l *Generator) generateInput(device *cthulstruct.InputDevice) (*libvirtstruct.Input, error) {
	input := &libvirtstruct.Input{}

	switch device.InputType {
	case cthulstruct.INPUT_MOUSE:
		input.MetaType = libvirtstruct.INPUT_MOUSE
	case cthulstruct.INPUT_TABLET:
		input.MetaType = libvirtstruct.INPUT_TABLET
	case cthulstruct.INPUT_KEYBOARD:
		input.MetaType = libvirtstruct.INPUT_KEYBOARD
	default:
		return nil, fmt.Errorf("unknown input type: %s", device.InputType)
	}
	
	switch device.InputBus {
	case cthulstruct.INPUT_PS2:
		input.MetaBus = libvirtstruct.INPUT_PS2
	case cthulstruct.INPUT_USB:
		input.MetaBus = libvirtstruct.INPUT_USB
	case cthulstruct.INPUT_VIRTIO:
		input.MetaBus = libvirtstruct.INPUT_VIRTIO
	default:
		return nil, fmt.Errorf("unknown input bus: %s", device.InputBus)
	}
	
	return input, nil
}
