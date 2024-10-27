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

package domain


// Domain represents a cthul domain. This format is used by the underlying domain controller
// to build up the vendor specific config (e.g. libvirt xml).
// Struct is annotated with json, yaml, and toml for easy external serialization/deserialization.
type Domain struct {
	UUID string
	Name string
	Title string
	Description string

	ResourceConfig ResourceConfig
	BootConfig BootConfig
	BlockDevices []BlockDevice
	NetworkDevices []NetworkDevice
	SerialDevices []SerialDevice
	GraphicDevices []GraphicDevice
	PCIDevices []PCIDevice
	USBDevices []USBDevice
}

type ResourceConfig struct {
	VCPUs int64
	Memory int64
}

type BOOT_OPTION string
const (
	BOOT_HD BOOT_OPTION = "cthul::boot::hd"
	BOOT_CD BOOT_OPTION = "cthul::boot::cd"
	BOOT_NETWORK BOOT_OPTION = "cthul::boot::network"
)

type BootConfig struct {
	SecureBoot bool
	BootOptions []BOOT_OPTION
}

type BlockDevice struct {
	GranitId string
	Virtio bool
}

type NetworkDevice struct {
	ProtonId string
	Virtio bool
}

type SerialDevice struct {
	
}

type GraphicDevice struct {
	
}

type PCIDevice struct {

}

type USBDevice struct {
	
}
