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
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package structure

// Domain represents a cthul domain. This format is used by the underlying domain controller
// to build up the vendor specific config (e.g. libvirt xml).
// Struct is annotated with json, yaml, and toml for easy external serialization/deserialization.
type Domain struct {
	Name        string `json:"name" toml:"name"`
	Title       string `json:"title" toml:"title"`
	Description string `json:"description" toml:"description"`

	SystemConfig   SystemConfig   `json:"system_config" toml:"system_config"`
	FirmwareConfig FirmwareConfig `json:"firmware_config" toml:"firmware_config"`
	ResourceConfig ResourceConfig `json:"resource_config" toml:"resource_config"`
	
	VideoDevices []VideoDevice `json:"video_devices" toml:"video_devices"`
	VideoAdapters []VideoAdapter `json:"video_adapters" toml:"video_adapters"`
	
	InputDevices []InputDevice `json:"input_devices" toml:"input_devices"`
	SerialDevices []SerialDevice `json:"serial_devices" toml:"serial_devices"`
	StorageDevices []StorageDevice `json:"storage_devices" toml:"storage_devices"`
	NetworkDevices []NetworkDevice `json:"network_devices" toml:"network_devices"`
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
	Architecture ARCH_OPTION    `json:"architecture" toml:"architecture"`
	Chipset      CHIPSET_OPTION `json:"chipset" toml:"chipset"`
}

type FIRMWARE_OPTION string

const (
	FIRMWARE_OVMF    FIRMWARE_OPTION = "cthul::firmware::ovmf"
	FIRMWARE_SEABIOS FIRMWARE_OPTION = "cthul::firmware::seabios"
)

type FirmwareConfig struct {
	Firmware       FIRMWARE_OPTION `json:"firmware" toml:"firmware"`
	SecureBoot     bool            `json:"secure_boot" toml:"secure_boot"`
	LoaderDeviceId string          `json:"loader_device_id" toml:"loader_device_id"`
	TmplDeviceId   string          `json:"tmpl_device_id" toml:"tmpl_device_id"`
	NvramDeviceId  string          `json:"nvram_device_id" toml:"nvram_device_id"`
}

type ResourceConfig struct {
	VCPUs  int64 `json:"vcpus" toml:"vcpus"`
	Memory int64 `json:"memory" toml:"memory"`
}

type VIDEO_OPTION string

const (
	VIDEO_VGA  VIDEO_OPTION = "cthul::video::vga"
	VIDEO_QXL  VIDEO_OPTION = "cthul::video::qxl"
	VIDEO_HOST VIDEO_OPTION = "cthul::video::host"
	VIDEO_NONE VIDEO_OPTION = "cthul::video::none"
)

type VideoDevice struct {
	VideoOption       VIDEO_OPTION `json:"video_option" toml:"video_option"`
	CommandBufferSize int64        `json:"commandbuffer_size" toml:"commandbuffer_size"`
	VideoBufferSize   int64        `json:"videobuffer_size" toml:"videobuffer_size"`
	FramebufferSize   int64        `json:"framebuffer_size" toml:"framebuffer_size"`
}

type VideoAdapter struct {
	DeviceId string `json:"device_id" toml:"device_id"`	
}

type SERIAL_BUS string

const (
	SERIAL_ISA SERIAL_BUS = "cthul::serial::isa"
	SERIAL_VIRTIO SERIAL_BUS = "cthul::serial::virtio"
)

type SerialDevice struct {
	DeviceId string `json:"device_id" toml:"device_id"`
	SerialBus SERIAL_BUS `json:"serial_bus" toml:"serial_bus"`
	Port int64 `json:"port" toml:"port"`
}

type INPUT_TYPE string

const (
	INPUT_MOUSE INPUT_TYPE = "cthul::input::mouse"
	INPUT_TABLET INPUT_TYPE = "cthul::input::tablet"
	INPUT_KEYBOARD INPUT_TYPE = "cthul::input::keyboard"
)

type INPUT_BUS string

const (
	INPUT_PS2 INPUT_BUS = "cthul::input::ps2"
	INPUT_USB INPUT_BUS = "cthul::input::usb"
	INPUT_VIRTIO INPUT_BUS = "cthul::input::virtio"
)

type InputDevice struct {
	InputType INPUT_TYPE `json:"input_type" toml:"input_type"`
	InputBus INPUT_BUS `json:"input_bus" toml:"input_bus"`
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
	DeviceId     string       `json:"device_id" toml:"device_id"`
	StorageType  STORAGE_TYPE `json:"storage_type" toml:"storage_type"`
	StorageBus   STORAGE_BUS  `json:"storage_bus" toml:"storage_bus"`
	BootPriority int64        `json:"boot_priority" toml:"boot_priority"`
}

type NETWORK_BUS string

const (
	NETWORK_E1000 NETWORK_BUS = "cthul::network::e1000"
	NETWORK_VIRTIO NETWORK_BUS = "cthul::network::virtio"
)

type NetworkDevice struct {
	DeviceId     string `json:"device_id" toml:"device_id"`
	NetworkBus NETWORK_BUS `json:"network_bus" toml:"network_bus"`
	BootPriority int64  `json:"boot_priority" toml:"boot_priority"`
}
