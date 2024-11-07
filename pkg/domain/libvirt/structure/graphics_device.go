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

type GRAPHICS_TYPE string
const (
	GRAPHICS_SPICE GRAPHICS_TYPE = "spice"
)

type Graphics struct {
	XMLName xml.Name `xml:"graphics"`
	MetaType GRAPHICS_TYPE `xml:"type,attr,omitempty"`
	Listen *GraphicsListen `xml:"listen,omitempty"`
}

type GRAPHICS_LISTEN_TYPE string
const (
	GRAPHICS_LISTEN_SOCKET GRAPHICS_LISTEN_TYPE = "socket"
)

type GraphicsListen struct {
	MetaType GRAPHICS_LISTEN_TYPE `xml:"type,attr,omitempty"`
	MetaPath string `xml:"path,attr,omitempty"`
}
