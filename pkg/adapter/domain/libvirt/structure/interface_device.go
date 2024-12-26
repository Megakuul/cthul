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

type INTERFACE_TYPE string

const (
	INTERFACE_BRIDGE   INTERFACE_TYPE = "bridge"
)

type Interface struct {
	XMLName xml.Name `xml:"interface"`
	MetaType INTERFACE_TYPE `xml:"type,attr,omitempty"`
	Source *InterfaceSource `xml:"source,omitempty"`
	Model *InterfaceModel `xml:"model,omitempty"`
	Boot *Boot `xml:"boot,omitempty"`
}

type InterfaceSource struct {
	MetaBridge string `xml:"bridge,attr,omitempty"`
}

type INTERFACE_MODEL_TYPE string

const (
	INTERFACE_MODEL_E1000 INTERFACE_MODEL_TYPE = "e1000"
	INTERFACE_MODEL_VIRTIO INTERFACE_MODEL_TYPE = "virtio"
)

type InterfaceModel struct {
	MetaType INTERFACE_MODEL_TYPE `xml:"type,attr,omitempty"`
}
