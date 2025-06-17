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

type DOMAIN_TYPE string

const (
	DOMAIN_KVM DOMAIN_TYPE = "kvm"
)

// Domain structure holds the relevant libvirt xml structure. It generally follows the rule that everything
// currently not required by cthul is not defined in this configuration.
type Domain struct {
	XMLName     xml.Name      `xml:"domain"`
	MetaType    DOMAIN_TYPE   `xml:"type,attr"`
	UUID        string        `xml:"uuid,omitempty"`
	Name        string        `xml:"name,omitempty"`
	Title       string        `xml:"title,omitempty"`
	Description string        `xml:"description,omitempty"`
	Memory      *Memory       `xml:"memory,omitempty"`
	VCPU        *VCPU         `xml:"vcpu,omitempty"`
	OS          *OS           `xml:"os,omitempty"`
	Features    *Features `xml:"features,omitempty"`
	Devices     *Devices `xml:"devices,omitempty"`
}

type CPU_PLACEMENT string

const (
	CPU_PLACEMENT_STATIC CPU_PLACEMENT = "static"
)

type VCPU struct {
	MetaPlacement CPU_PLACEMENT `xml:"placement,attr"`
	Data          int64         `xml:",chardata"`
}

type MEMORY_UNIT string

const (
	MEMORY_UNIT_BYTES MEMORY_UNIT = "bytes"
)

type Memory struct {
	MetaUnit MEMORY_UNIT `xml:"unit,attr"`
	Data     int64       `xml:",chardata"`
}

type Devices struct {
	Devices []any `xml:"devices"`
}

type Features struct {
	Features []any `xml:"features"`
}
