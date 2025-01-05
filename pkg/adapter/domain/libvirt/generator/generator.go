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
	granitdevice "cthul.io/cthul/pkg/granit/device"
	protondevice "cthul.io/cthul/pkg/proton/device"
	wavedevice "cthul.io/cthul/pkg/wave/device"
	"cthul.io/cthul/pkg/wave/video"
)

// Generator provides operations to generate libvirt xml from cthul domain configurations.
// The generator uses provided device controllers to lookup information provided by external cthul devices
// (resolving things like 'GranitBlockDeviceId').
// It also provides operations to attach and release those required devices.
type Generator struct {
	// waveRunRoot specifies the wave base path for runtime files (unix-sockets and stuff).
	waveRunRoot string
	
	video *video.Controller
	granit *granitdevice.DeviceController
	proton *protondevice.DeviceController
}

type GeneratorOption func(*Generator)

func NewGenerator(
	videoDevice *video.Controller,
	granitDevice *granitdevice.DeviceController,
	protonDevice *protondevice.DeviceController,
	opts ...GeneratorOption) *Generator {

	generator := &Generator{
		waveRunRoot: "/run/cthul/wave/",
		video: videoDevice,
		granit: granitDevice,
		proton: protonDevice,
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

// WithWaveRunRoot defines a custom root path for wave runtime files (sockets, etc.).
func WithWaveRunRoot(path string) GeneratorOption {
	return func(g *Generator) {
		g.waveRunRoot = path
	}
}
