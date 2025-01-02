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
	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
)

// Generate transpiles the domain config to a libvirt xml file. Cthul devices are dynamically resolved with
// the generator attached device controllers. Devices must be attached to the node otherwise lookups will fail.
func (l *Generator) Generate(id string, config *cthulstruct.Domain) (*libvirtstruct.Domain, error) {
	var err error
	domain := &libvirtstruct.Domain{
		MetaType: libvirtstruct.DOMAIN_KVM,
		UUID: id,
		Name: config.Name,
		Title: config.Title,
		Description: config.Description,
		VCPU: l.generateVCPU(&config.ResourceConfig),
		Memory: l.generateMemory(&config.ResourceConfig),
		Devices: []interface{}{},
		Features: []interface{}{},
	}
	
	domain.OS, err = l.generateOS(&config.SystemConfig, &config.FirmwareConfig)
	if err!=nil {
		return nil, err
	}

	for _, videoDevice := range config.VideoDevices {
		device, err := l.generateVideo(&videoDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}
	
	for _, videoAdapter := range config.VideoAdapters {
		device, err := l.generateGraphic(&videoAdapter)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}
	
	for _, serialDevice := range config.SerialDevices {
		device, err := l.generateSerial(&serialDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}
	
	for _, serialDevice := range config.SerialDevices {
		device, err := l.generateSerial(&serialDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}
	
	for _, inputDevice := range config.InputDevices {
		device, err := l.generateInput(&inputDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}
	
	for _, storageDevice := range config.StorageDevices {
		device, err := l.generateDisk(&storageDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, networkDevice := range config.NetworkDevices {
		device, err := l.generateInterface(&networkDevice)
		if err!=nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	return domain, nil
}
