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

package generator

import (
	"fmt"

	cthulstruct "cthul.io/cthul/pkg/domain/structure"
	libvirtstruct "cthul.io/cthul/pkg/domain/libvirt/structure"
)

// generateOS generates the libvirt operating system configuration from system and firmware information.
func (l *LibvirtGenerator) generateOS(system *cthulstruct.SystemConfig, firmware *cthulstruct.FirmwareConfig) (*libvirtstruct.OS, error) {
	os := &libvirtstruct.OS{
		Type: &libvirtstruct.OSType{
			Data: "hvm",
		},
		Loader: &libvirtstruct.OSLoader{
			MetaReadonly: true,
			MetaSecure: firmware.SecureBoot,
		},
	}

	switch system.Architecture {
	case cthulstruct.ARCH_AMD64:
		os.Type.Arch = libvirtstruct.OS_ARCH_X86_64
	case cthulstruct.ARCH_AARCH64:
		os.Type.Arch = libvirtstruct.OS_ARCH_AARCH64
	default:
		return nil, fmt.Errorf("unknown system architecture: %s", system.Architecture)
	}

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

	loaderDevice, err := l.granit.LookupBlock(firmware.LoaderDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware loader device lookup: %s", err.Error())
	}
	
	templateDevice, err := l.granit.LookupBlock(firmware.TemplateDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware template device lookup: %s", err.Error())
	}

	nvramDevice, err := l.granit.LookupBlock(firmware.NvramDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware nvram device lookup: %s", err.Error())
	}

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
