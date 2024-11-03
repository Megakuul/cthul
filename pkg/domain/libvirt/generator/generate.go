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

	cthulstruct "cthul.io/cthul/pkg/domain/structure"
	libvirtstruct "cthul.io/cthul/pkg/domain/libvirt/structure"
)

// Generate transpiles the domain config to a libvirt xml file. Cthul devices are dynamically resolved with
// the generator attached device controllers. Devices must be attached to the node otherwise lookups will fail.
func (l *LibvirtGenerator) Generate(config *cthulstruct.Domain) (*libvirtstruct.Domain, error) {
	domain := &libvirtstruct.Domain{
		MetaType: libvirtstruct.KVM,
		UUID: config.UUID,
		Name: config.Name,
		Title: config.Title,
		Description: config.Description,
		VCPU: generateVCPU(&config.ResourceConfig),
		Memory: generateMemory(&config.ResourceConfig),
		Devices: []interface{}{},
		Features: []interface{}{},
	}

	domain.OS, err := generateOS(&config.BootConfig)
	if err!=nil {
		return nil, err
	}
	
	for _, blockDevice := range config.BlockDevices {
		device, err := generateBlockDevice(blockDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, networkDevice := range config.NetworkDevices {
		device, err := generateNetworkDevice(networkDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, serialDevice := range config.SerialDevices {
		device, err := generateSerialDevice(serialDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, videoDevice := range config.VideoDevices {
		device, err := generateVideoDevice(videoDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}
	
	for _, graphicDevice := range config.GraphicDevices {
		device, err := generateGraphicDevice(graphicDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	return domain, nil
}

func generateVCPU(config *cthulstruct.ResourceConfig) *libvirtstruct.VCPU {
	return &libvirtstruct.VCPU{
		MetaPlacement: libvirtstruct.STATIC,
		Data: config.VCPUs,
	}
}

func generateMemory(config *cthulstruct.ResourceConfig) *libvirtstruct.Memory {
	return &libvirtstruct.Memory{
		MetaUnit: libvirtstruct.BYTES,
		Data: config.Memory,
	}
}
