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

// Explanation: A libvirt video device provides frame- and commandbuffers to the guest os.
// The video device is usually provided as PCI device to the guest. Buffers are allocated by the hypervisor
// inside the guests main memory and their locations can be queried via MMIO instructions.
// Using a simple VGA video device will simply provide the guest with a framebuffer it can write pixel data to.
// The QXL video device is optimized for the spice system, it contains a fallback vga framebuffer, but
// mainly leverages a client qxl driver which intercepts 2d graphic instructions in the window system (e.g. x11)
// sending them directly to a command ringbuffer. With this system, the lightweight 2d instructions
// (and rendered 3d bitmaps) are sent to the client, which then renders the instructions.
// The VIRTIO (also called virtio-gpu) device requires virtio drivers installed on the guest. Unlike the other
// options, the framebuffer is not allocated, instead the driver itself uses virtqueues to communicate with the
// hypervisor. Instead of rendering 2d and 3d objects on the guest, OpenGL calls are translated to an
// intermediate format that is written to the avail virtqueue. The virglrenderer on the host is now responsible
// for rendering this intermediate representation leveraging the hosts GPU.

// generateVideo generates a libvirt video device from the cthul video device.
func (l *Generator) generateVideo(device *cthulstruct.VideoDevice) (*libvirtstruct.Video, error) {
	video := &libvirtstruct.Video{
		Model: &libvirtstruct.VideoModel{},
	}

	switch device.VideoOption {
	case cthulstruct.VIDEO_NONE:
		video.Model.MetaType = libvirtstruct.VIDEO_MODEL_NONE
	case cthulstruct.VIDEO_VGA:
		video.Model.MetaType = libvirtstruct.VIDEO_MODEL_VGA
		video.Model.MetaVRam = device.VideoBufferSize + device.FramebufferSize // framebuffer is inside vram in vga
		video.Model.MetaVGAMem = device.FramebufferSize
	case cthulstruct.VIDEO_QXL:
		video.Model.MetaType = libvirtstruct.VIDEO_MODEL_QXL
		video.Model.MetaRam = device.CommandBufferSize + device.FramebufferSize // framebuffer is inside ram in qxl
		video.Model.MetaVRam = device.VideoBufferSize
		video.Model.MetaVGAMem = device.FramebufferSize
	case cthulstruct.VIDEO_HOST:
		video.Model.MetaType = libvirtstruct.VIDEO_MODEL_VIRTIO
	default:
		return nil, fmt.Errorf("unknown video option: %s", device.VideoOption)
	}
	
	return video, nil
}
