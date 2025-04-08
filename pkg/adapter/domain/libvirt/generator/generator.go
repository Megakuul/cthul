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
	"cthul.io/cthul/pkg/wave/video"
)

// Generator provides operations to generate libvirt xml from cthul domain configurations.
// The generator uses provided device controllers to lookup information provided by external cthul devices
// (resolving things like 'GranitBlockDeviceId').
// It also provides operations to attach and release those required devices.
type Generator struct {
  nodeId string

	video *video.Controller
	granit *granitdevice.DeviceController
	proton *protondevice.DeviceController

  videoRoot string
  granitRoot string
  protonRoot string
}

type Option func(*Generator)

func New(
  nodeId string,
	videoDevice *video.Controller,
	granitDevice *granitdevice.DeviceController,
	protonDevice *protondevice.DeviceController,
	opts ...Option) *Generator {

	generator := &Generator{
    nodeId: nodeId,
		video: videoDevice,
		granit: granitDevice,
		proton: protonDevice,
    videoRoot: "/run/cthul/wave/video/",
    granitRoot: "/run/cthul/granit/",
    protonRoot: "/run/cthul/proton/",
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

