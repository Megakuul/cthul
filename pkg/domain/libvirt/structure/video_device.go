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

package structure

type Video struct {
	Model *VideoModel `xml:"model,omitempty"`
}

type VIDEO_MODEL_TYPE string

const (
	VIDEO_MODEL_VGA    VIDEO_MODEL_TYPE = "vga"    // just writes raw pixel data to the framebuffer
	VIDEO_MODEL_VIRTIO VIDEO_MODEL_TYPE = "virtio" // sends unrendered data to the host (framebuffer on host)
	VIDEO_MODEL_QX1    VIDEO_MODEL_TYPE = "qx1"    // writes spice protocol data to the device buffers
)

type VideoModel struct {
	MetaType   VIDEO_MODEL_TYPE `xml:"type,attr,omitempty"`
	MetaRam    int64            `xml:"ram,attr,omitempty"` // qx1: vga fallback buffer and command / release ring.
	MetaVRam   int64            `xml:"vram,attr,omitempty"` // qx1: cache buffer. vga: framebuffer.
	MetaVGAMem int64            `xml:"vgamem,attr,omitempty"` // qx1: size of vga fallback buffer inside 'ram'.
	// virtio drivers use virtqueue ringbuffers that are self allocated by the guest not by the device.
}
