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
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package generator

import (
	"fmt"

	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
)

// Explanation: A libvirt os tag represents configuration options about the emulated cpu, mainboard and firmware
// components. The cpu architecture is rather irrelevant as we do hardware assisted full virtualisation, this
// makes it impossible to emulate cpu architectures as the instructions are executed directly on the cpu.
// The chipset configuration is more important, it defines the chipset interface provided to the guest os.
// Sending instructions to the chipset is intercepted by the cpu virtualization-extension, which then
// invokes the host device (e.g. qemu) that handles the request. Depending on the chipset configuration
// qemu will use different communication methods and protocols.
// The firmware configuration defines what firmware is launched by the hypervisor upon starting a domain,
// the hypervisor essentially moves the compiled firmware to memory and points the program counter to it.
// To start the firmware, loader, template and nvram files are required. All of those are not used directly, but
// rather moved to the memory by the hypervisor (nvram is also written back on domain shutdown).

// Loader, template and nvram files are provided by granit devices that are loaded as blockdevs on the host.

// generateOS generates the libvirt operating system configuration from system and firmware information.
func (l *Generator) generateOS(system *cthulstruct.SystemConfig, firmware *cthulstruct.FirmwareConfig) (*libvirtstruct.OS, error) {
	os := &libvirtstruct.OS{
		Type: &libvirtstruct.OSType{
			Data: "hvm",
		},
		Loader: &libvirtstruct.OSLoader{
			MetaReadonly: true,
			MetaSecure: firmware.SecureBoot,
		},
	}

	loaderDevice, err := l.granit.LookupStorage(firmware.LoaderDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware loader device lookup: %s", err.Error())
	}
	if loaderDevice.Type != granit.STORAGE_FILE {
		return nil, fmt.Errorf("only supported firmware loader device type is: %s", granit.STORAGE_FILE)
	}
	
	templateDevice, err := l.granit.LookupStorage(firmware.TmplDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware template device lookup: %s", err.Error())
	}
	if templateDevice.Type != granit.STORAGE_FILE {
		return nil, fmt.Errorf("only supported firmware template device type is: %s", granit.STORAGE_FILE)
	}

	nvramDevice, err := l.granit.LookupStorage(firmware.NvramDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware nvram device lookup: %s", err.Error())
	}
	if nvramDevice.Type != granit.STORAGE_FILE {
		return nil, fmt.Errorf("only supported firmware nvram device type is: %s", granit.STORAGE_FILE)
	}

	// Architecture
	switch system.Architecture {
	case cthulstruct.ARCH_AMD64:
		os.Type.Arch = libvirtstruct.OS_ARCH_X86_64
	case cthulstruct.ARCH_AARCH64:
		os.Type.Arch = libvirtstruct.OS_ARCH_AARCH64
	default:
		return nil, fmt.Errorf("unknown system architecture: %s", system.Architecture)
	}

	// Chipset
	switch system.Chipset {
	case cthulstruct.CHIPSET_I440FX:
		os.Type.Machine = libvirtstruct.OS_CHIPSET_I440FX
	case cthulstruct.CHIPSET_Q35:
		os.Type.Machine = libvirtstruct.OS_CHIPSET_Q35
	case cthulstruct.CHIPSET_VIRT:
		os.Type.Machine = libvirtstruct.OS_CHIPSET_VIRT
	default:
		return nil, fmt.Errorf("unknown system chipset: %s", system.Chipset)
	}

	// Firmware
	os.Loader.Data = loaderDevice.Path
	
	switch firmware.Firmware {
	case cthulstruct.FIRMWARE_OVMF:
		os.Loader.MetaType = libvirtstruct.OS_LOADER_OVMF
		os.Nvram = &libvirtstruct.OSNvram{
			MetaType: libvirtstruct.OS_NVRAM_FILE,
			MetaTemplate: templateDevice.Path,
			Source: nvramDevice.Path,
		}
	case cthulstruct.FIRMWARE_SEABIOS:
		os.Loader.MetaType = libvirtstruct.OS_LOADER_SEABIOS
		os.Nvrams = &libvirtstruct.OSNvram{
			MetaType: libvirtstruct.OS_NVRAM_FILE,
			MetaTemplate: templateDevice.Path,
			Source: nvramDevice.Path,
		}
	default:
		return nil, fmt.Errorf("unknown firmware type: %s", firmware.Firmware)
	}
	
	return os, nil
}
