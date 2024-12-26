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

import "encoding/xml"

type INPUT_TYPE string
const (
	INPUT_MOUSE INPUT_TYPE = "mouse"
	INPUT_KEYBOARD INPUT_TYPE = "keyboard"
	INPUT_TABLET INPUT_TYPE = "tablet"
)

type INPUT_BUS string
const (
	INPUT_PS2 INPUT_BUS  = "ps2"
	INPUT_USB INPUT_BUS = "usb"
	INPUT_VIRTIO INPUT_BUS = "virtio"
)

type Input struct {
	XMLName xml.Name `xml:"input"`
	MetaType INPUT_TYPE `xml:"type,attr,omitempty"`
	MetaBus INPUT_BUS `xml:"bus,attr,omitempty"`
}
