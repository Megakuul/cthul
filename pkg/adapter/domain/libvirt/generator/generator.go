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
	"cthul.io/cthul/pkg/granit/disk"
	"cthul.io/cthul/pkg/proton/inter"
	"cthul.io/cthul/pkg/wave/serial"
	"cthul.io/cthul/pkg/wave/video"
)

// Generator provides operations to generate libvirt xml from cthul domain configurations.
// The generator uses provided device controllers to lookup information provided by external cthul devices
// (resolving things like 'GranitBlockDeviceId').
// It also provides operations to attach and release those required devices.
type Generator struct {
  nodeId string

	video *video.Controller
  serial *serial.Controller
	disk *disk.Controller
	inter *inter.Controller

  waveRoot string
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
    nodeId: nodeId,
		video: videoController,
    serial: serialController,
    disk: diskController,
    inter: interController,
    waveRoot: "/run/cthul/wave",
    granitRoot: "/run/cthul/granit/",
    protonRoot: "/run/cthul/proton/",
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

