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

	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
  "cthul.io/cthul/pkg/api/wave/v1/domain"
)

// Explanation: A libvirt input device can represent multiple human interface devices, like a mouse, tablet or
// keyboard. With a basic bus like PS2 data is transferred via PIO, the host device (qemu) performs ISR
// interrupts to send data to the guest os.
// Using the VIRTIO bus leverages virtqueues to transfer data between the guest and host os.

// generateInterface generates a libvirt network interface device from the cthul network device.
func (l *Generator) generateInput(device *domain.InputDevice) (*structure.Input, error) {
	input := &structure.Input{}

	switch device.InputType {
	case domain.InputType_INPUT_TYPE_MOUSE:
		input.MetaType = structure.INPUT_MOUSE
	case domain.InputType_INPUT_TYPE_TABLET:
		input.MetaType = structure.INPUT_TABLET
	case domain.InputType_INPUT_TYPE_KEYBOARD:
		input.MetaType = structure.INPUT_KEYBOARD
	default:
		return nil, fmt.Errorf("unknown input type: %s", device.InputType)
	}
	
	switch device.InputBus {
	case domain.InputBus_INPUT_BUS_PS2:
		input.MetaBus = structure.INPUT_PS2
	case domain.InputBus_INPUT_BUS_USB:
		input.MetaBus = structure.INPUT_USB
	case domain.InputBus_INPUT_BUS_VIRTIO:
		input.MetaBus = structure.INPUT_VIRTIO
	default:
		return nil, fmt.Errorf("unknown input bus: %s", device.InputBus)
	}
	
	return input, nil
}
