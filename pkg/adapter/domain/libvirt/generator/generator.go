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

	"cthul.io/cthul/pkg/granit/disk"
	"cthul.io/cthul/pkg/proton/inter"
	"cthul.io/cthul/pkg/wave/serial"
	"cthul.io/cthul/pkg/wave/video"

	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	"cthul.io/cthul/pkg/api/wave/v1/domain"
)

// Generator provides operations to generate libvirt xml from cthul domain configurations.
// The generator uses provided device controllers to lookup information provided by external cthul devices
// (resolving things like 'GranitBlockDeviceId').
// It also provides operations to attach and release those required devices.
type Generator struct {
	nodeId string

	video  *video.Controller
	serial *serial.Controller
	disk   *disk.Controller
	inter  *inter.Controller

	waveRoot   string
	granitRoot string
	protonRoot string
}

type Option func(*Generator)

func New(
	nodeId string,
	videoController *video.Controller,
	serialController *serial.Controller,
	diskController *disk.Controller,
	interController *inter.Controller,
	opts ...Option) *Generator {

	generator := &Generator{
		nodeId:     nodeId,
		video:      videoController,
		serial:     serialController,
		disk:       diskController,
		inter:      interController,
		waveRoot:   "/run/cthul/wave",
		granitRoot: "/run/cthul/granit/",
		protonRoot: "/run/cthul/proton/",
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

// Attach installs / locks all devices that are required by the domain config.
func (g *Generator) Attach(ctx context.Context, config *domain.DomainConfig) error {
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
	
  return nil
}

// Generate transpiles the domain config to a libvirt xml file. Cthul devices are dynamically resolved with
// the generator attached device controllers. Devices must be attached to the node otherwise lookups will fail.
func (l *Generator) Generate(ctx context.Context, id string, config *domain.DomainConfig) (*structure.Domain, error) {
	var err error
	domain := &structure.Domain{
		MetaType:    structure.DOMAIN_KVM,
		UUID:        id,
		Name:        config.Name,
		Title:       config.Title,
		Description: config.Description,
		VCPU:        l.generateVCPU(config.ResourceConfig),
		Memory:      l.generateMemory(config.ResourceConfig),
		Devices:     []any{},
		Features:    []any{},
	}

	if config.GetSystemConfig() == nil || config.GetFirmwareConfig() == nil {
		return nil, fmt.Errorf("TODO: check failed, nil")
	}
	domain.OS, err = l.generateOS(ctx, config.GetSystemConfig(), config.GetFirmwareConfig())
	if err != nil {
		return nil, err
	}

	for _, videoDevice := range config.VideoDevices {
		device, err := l.generateVideo(videoDevice)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, videoAdapter := range config.VideoAdapters {
		device, err := l.generateGraphic(ctx, videoAdapter)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, serialDevice := range config.SerialDevices {
		device, err := l.generateSerial(ctx, serialDevice)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, serialDevice := range config.SerialDevices {
		device, err := l.generateSerial(ctx, serialDevice)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, inputDevice := range config.InputDevices {
		device, err := l.generateInput(inputDevice)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, storageDevice := range config.StorageDevices {
		device, err := l.generateDisk(ctx, storageDevice)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	for _, networkDevice := range config.NetworkDevices {
		device, err := l.generateInterface(ctx, networkDevice)
		if err != nil {
			return nil, err
		}
		domain.Devices = append(domain.Devices, device)
	}

	return domain, nil
}
