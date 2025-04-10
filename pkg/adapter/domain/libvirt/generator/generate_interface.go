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

	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	interstruct "cthul.io/cthul/pkg/proton/inter/structure"
)

// Explanation: A libvirt interface device is a network adapter for the guest os. It intercepts MMIO calls
// from the guest and interacts with the host network (based on the type, bridge for example just sends
// the data to a host os bridge). With E1000 the data is transferred with an emulated DMA buffer.
// The VIRTIO bus on the other hand uses virtqueues to transfer data.

// generateInterface generates a libvirt network interface device from the cthul network device.
func (g *Generator) generateInterface(ctx context.Context, device *cthulstruct.NetworkDevice) (*libvirtstruct.Interface, error) {
	inter := &libvirtstruct.Interface{
		Model: &libvirtstruct.InterfaceModel{},
		Source: &libvirtstruct.InterfaceSource{},
		Boot: &libvirtstruct.Boot{MetaOrder: device.BootPriority},
	}

	interDevice, err := g.inter.Lookup(ctx, device.DeviceId)
	if err!=nil {
		return nil, err
	}

	switch device.NetworkBus {
	case cthulstruct.NETWORK_E1000:
		inter.Model.MetaType = libvirtstruct.INTERFACE_MODEL_E1000
	case cthulstruct.NETWORK_VIRTIO:
		inter.Model.MetaType = libvirtstruct.INTERFACE_MODEL_VIRTIO
	default:
		return nil, fmt.Errorf("unknown network bus type: %s", device.NetworkBus)
	}

	switch interDevice.Type {
	case interstruct.INTER_BRIDGE:
		inter.MetaType = libvirtstruct.INTERFACE_BRIDGE
		inter.Source.MetaBridge = interDevice.Device
	default:
		return nil, fmt.Errorf("unsupported interface type: %s", interDevice.Type)
	}
	
	return inter, nil
}
