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

	SystemConfig   SystemConfig   `json:"system_config" yaml:"system_config" toml:"system_config"`
	FirmwareConfig FirmwareConfig `json:"firmware_config" yaml:"firmware_config" toml:"firmware_config"`
	ResourceConfig ResourceConfig `json:"resource_config" yaml:"resource_config" toml:"resource_config"`

	StorageDevices []StorageDevice `json:"storage_devices" yaml:"storage_devices" toml:"storage_devices"`
	NetworkDevices []NetworkDevice `json:"network_devices" yaml:"network_devices" toml:"network_devices"`

	VideoDevices   []VideoDevice   `json:"video_devices" yaml:"video_devices" toml:"video_devices"`
	GraphicDevices []GraphicDevice `json:"graphic_devices" yaml:"graphic_devices" toml:"graphic_devices"`
	SerialDevices  []SerialDevice  `json:"serial_devices" yaml:"serial_devices" toml:"serial_devices"`
}

type ARCH_OPTION string

const (
	ARCH_AMD64   ARCH_OPTION = "cthul::arch::amd64"
	ARCH_AARCH64 ARCH_OPTION = "cthul::arch::aarch64"
)

type CHIPSET_OPTION string

const (
	CHIPSET_I440FX CHIPSET_OPTION = "cthul::chipset::i440fx"
	CHIPSET_Q35    CHIPSET_OPTION = "cthul::chipset::q35"
	CHIPSET_VIRT   CHIPSET_OPTION = "cthul::chipset::virt"
)

type SystemConfig struct {
	Architecture ARCH_OPTION    `json:"architecture" yaml:"architecture" toml:"architecture"`
	Chipset      CHIPSET_OPTION `json:"chipset" yaml:"chipset" toml:"chipset"`
}

type FIRMWARE_OPTION string

const (
	FIRMWARE_OVMF    FIRMWARE_OPTION = "cthul::firmware::ovmf"
	FIRMWARE_SEABIOS FIRMWARE_OPTION = "cthul::firmware::seabios"
)

type FirmwareConfig struct {
	Firmware       FIRMWARE_OPTION `json:"firmware" yaml:"firmware" toml:"firmware"`
	SecureBoot     bool            `json:"secure_boot" yaml:"secure_boot" toml:"secure_boot"`
	LoaderDeviceId string          `json:"loader_device_id" yaml:"loader_device_id" toml:"loader_device_id"`
	TmplDeviceId   string          `json:"tmpl_device_id" yaml:"tmpl_device_id" toml:"tmpl_device_id"`
	NvramDeviceId  string          `json:"nvram_device_id" yaml:"nvram_device_id" toml:"nvram_device_id"`
}

type ResourceConfig struct {
	VCPUs  int64 `json:"vcpus" yaml:"vcpus" toml:"vcpus"`
	Memory int64 `json:"memory" yaml:"memory" toml:"memory"`
}

type STORAGE_TYPE string

const (
	STORAGE_CDROM STORAGE_TYPE = "cthul::storage::cdrom"
	STORAGE_DISK  STORAGE_TYPE = "cthul::storage::disk"
)

type STORAGE_BUS string

const (
	STORAGE_IDE    STORAGE_BUS = "cthul::storage::ide"
	STORAGE_SATA   STORAGE_BUS = "cthul::storage::sata"
	STORAGE_VIRTIO STORAGE_BUS = "cthul::storage::virtio"
)

type StorageDevice struct {
	DeviceId     string       `json:"device_id" yaml:"device_id" toml:"device_id"`
	StorageType  STORAGE_TYPE `json:"storage_type" yaml:"storage_type" toml:"storage_type"`
	StorageBus   STORAGE_BUS  `json:"storage_bus" yaml:"storage_bus" toml:"storage_bus"`
	BootPriority int64        `json:"boot_priority" yaml:"boot_priority" toml:"boot_priority"`
}

type NetworkDevice struct {
	DeviceId     string `json:"device_id" yaml:"device_id" toml:"device_id"`
	Virtio       bool   `json:"virtio" yaml:"virtio" toml:"virtio"`
	BootPriority int64  `json:"boot_priority" yaml:"boot_priority" toml:"boot_priority"`
}

type VIDEO_OPTION string

const (
	VIDEO_VGA  VIDEO_OPTION = "cthul::video::vga"
	VIDEO_QXL  VIDEO_OPTION = "cthul::video::qxl"
	VIDEO_HOST VIDEO_OPTION = "cthul::video::host"
	VIDEO_NONE VIDEO_OPTION = "cthul::video::none"
)

type VideoDevice struct {
	VideoOption       VIDEO_OPTION `json:"video_option" yaml:"video_option" toml:"video_option"`
	CommandBufferSize int64        `json:"commandbuffer_size" yaml:"commandbuffer_size" toml:"commandbuffer_size"`
	VideoBufferSize   int64        `json:"videobuffer_size" yaml:"videobuffer_size" toml:"videobuffer_size"`
	FramebufferSize   int64        `json:"framebuffer_size" yaml:"framebuffer_size" toml:"framebuffer_size"`
}

type SerialDevice struct {
	WaveSerialDeviceId string `json:"device_id" yaml:"device_id" toml:"device_id"`
}

type GraphicDevice struct {
	WaveGraphicDeviceId string `json:"device_id" yaml:"device_id" toml:"device_id"`
}
