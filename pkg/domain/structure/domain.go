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

// Domain represents a cthul domain. This format is used by the underlying domain controller
// to build up the vendor specific config (e.g. libvirt xml).
// Struct is annotated with json, yaml, and toml for easy external serialization/deserialization.
type Domain struct {
	UUID        string `json:"uuid" yaml:"uuid" toml:"uuid"`
	Name        string `json:"name" yaml:"name" toml:"name"`
	Title       string `json:"title" yaml:"title" toml:"title"`
	Description string `json:"description" yaml:"description" toml:"description"`

	ResourceConfig ResourceConfig `json:"resource_config" yaml:"resource_config" toml:"resource_config"`
	BootConfig     BootConfig     `json:"boot_config" yaml:"boot_config" toml:"boot_config"`

	BlockDevices   []BlockDevice   `json:"block_devices" yaml:"block_devices" toml:"block_devices"`
	NetworkDevices []NetworkDevice `json:"network_devices" yaml:"network_devices" toml:"network_devices"`
	SerialDevices  []SerialDevice  `json:"serial_devices" yaml:"serial_devices" toml:"serial_devices"`
	GraphicDevices []GraphicDevice `json:"graphic_devices" yaml:"graphic_devices" toml:"graphic_devices"`
}

type ARCH_OPTION string

const (
	ARCH_AMD64      ARCH_OPTION = "cthul::arch::amd64"
	ARCH_AARCH64      ARCH_OPTION = "cthul::arch::aarch64"
)

type CHIPSET_OPTION string

const (
	CHIPSET_I440FX CHIPSET_OPTION = "cthul::chipset::i440fx"
	CHIPSET_Q35      CHIPSET_OPTION = "cthul::chipset::q35"
	CHIPSET_VIRT      CHIPSET_OPTION = "cthul::chipset::virt"
)

type SystemConfig struct {
	Architecture ARCH_OPTION `json:"architecture" yaml:"architecture" toml:"architecture"`
	Chipset CHIPSET_OPTION `json:"chipset" yaml:"chipset" toml:"chipset"`
}

type ResourceConfig struct {
	VCPUs  int64 `json:"vcpus" yaml:"vcpus" toml:"vcpus"`
	Memory int64 `json:"memory" yaml:"memory" toml:"memory"`
}

type BOOT_OPTION string

const (
	BOOT_HD      BOOT_OPTION = "cthul::boot::hd"
	BOOT_CD      BOOT_OPTION = "cthul::boot::cd"
	BOOT_NETWORK BOOT_OPTION = "cthul::boot::network"
)

type BootConfig struct {
	SecureBoot  bool          `json:"secure_boot" yaml:"secure_boot" toml:"secure_boot"`
	BootOptions []BOOT_OPTION `json:"boot_options" yaml:"boot_options" toml:"boot_options"`
}

type BlockDevice struct {
	GranitBlockDeviceId string `json:"device_id" yaml:"device_id" toml:"device_id"`
	Virtio              bool   `json:"virtio" yaml:"virtio" toml:"virtio"`
}

type NetworkDevice struct {
	ProtonNetworkDeviceId string `json:"device_id" yaml:"device_id" toml:"device_id"`
	Virtio                bool   `json:"virtio" yaml:"virtio" toml:"virtio"`
}

type VIDEO_OPTION string

const (
	VIDEO_VGA    VIDEO_OPTION = "cthul::video::vga"
	VIDEO_QX1    VIDEO_OPTION = "cthul::video::qx1"
	VIDEO_HOST VIDEO_OPTION = "cthul::video::host"
	VIDEO_NONE VIDEO_OPTION = "cthul::video::none"
)

type VideoDevice struct {
	VideoOption VIDEO_OPTION `json:"video_option" yaml:"video_option" toml:"video_option"`
	ProtocolbufferSize int64 `json:"protocolbuffer_size" yaml:"protocolbuffer_size" toml:"protocolbuffer_size"`
	FramebufferSize int64 `json:"framebuffer_size" yaml:"framebuffer_size" toml:"framebuffer_size"`
}

type SerialDevice struct {
	WaveSerialDeviceId string `json:"device_id" yaml:"device_id" toml:"device_id"`
}

type GraphicDevice struct {
	WaveGraphicDeviceId string `json:"device_id" yaml:"device_id" toml:"device_id"`
}
