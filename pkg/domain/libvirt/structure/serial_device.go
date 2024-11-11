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

type SERIAL_TYPE string
const (
	SERIAL_UNIX SERIAL_TYPE = "unix"
)

type Serial struct {
	XMLName xml.Name `xml:"serial"`
	MetaType SERIAL_TYPE `xml:"type,attr,omitempty"`
	Source *SerialSource `xml:"source,omitempty"`
	Target *SerialTarget `xml:"target,omitempty"`
}

type SERIAL_SOURCE_MODE string
const (
	SERIAL_SOURCE_BIND SERIAL_SOURCE_MODE = "bind"
)

type SerialSource struct {
	MetaMode SERIAL_SOURCE_MODE `xml:"mode,attr,omitempty"`
	MetaPath string `xml:"path,attr,omitempty"`
}

type SERIAL_BUS_TYPE string
const (
	SERIAL_BUS_ISA SERIAL_BUS_TYPE = "isa-serial"
	SERIAL_BUS_VIRTIO SERIAL_BUS_TYPE = "virtio"
)

type SerialTarget struct {
	MetaType SERIAL_BUS_TYPE `xml:"type,attr,omitempty"`
	MetaPort int64 `xml:"port,attr,omitempty"`
}
