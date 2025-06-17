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
	"context"
	"fmt"
	"path/filepath"

	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	"cthul.io/cthul/pkg/api/granit/v1/disk"
	"cthul.io/cthul/pkg/api/wave/v1/domain"
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
func (g *Generator) generateOS(ctx context.Context, system *domain.SystemConfig, firmware *domain.FirmwareConfig) (*structure.OS, error) {
	// just provisorisch
	return &structure.OS{
		Type: &structure.OSType{
			Arch: structure.OS_ARCH_X86_64,
			Machine: structure.OS_CHIPSET_I440FX,
			Data: "hvm",
		},
	}, nil

	os := &structure.OS{
		Type: &structure.OSType{
			Data: "hvm",
		},
		Loader: &structure.OSLoader{
			MetaReadonly: true,
			MetaSecure: firmware.SecureBoot,
		},
	}

	loaderDevice, err := g.disk.Lookup(ctx, firmware.LoaderDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware loader device lookup: %s", err.Error())
	}
	if loaderDevice.Config.Format != disk.DiskFormat_DISK_FORMAT_RAW {
		return nil, fmt.Errorf("only supported firmware loader device format is 'raw'")
	}
	
	templateDevice, err := g.disk.Lookup(ctx, firmware.TmplDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware template device lookup: %s", err.Error())
	}
	if templateDevice.Config.Format != disk.DiskFormat_DISK_FORMAT_RAW {
		return nil, fmt.Errorf("only supported firmware template device format is 'raw'")
	}

	nvramDevice, err := g.disk.Lookup(ctx, firmware.NvramDeviceId)
	if err!=nil {
		return nil, fmt.Errorf("firmware nvram device lookup: %s", err.Error())
	}
	if nvramDevice.Config.Format != disk.DiskFormat_DISK_FORMAT_RAW {
		return nil, fmt.Errorf("only supported firmware nvram device type is 'raw'")
	}

	// Architecture
	switch system.Architecture {
	case domain.Arch_ARCH_AMD64:
		os.Type.Arch = structure.OS_ARCH_X86_64
	case domain.Arch_ARCH_AARCH64:
		os.Type.Arch = structure.OS_ARCH_AARCH64
	default:
		return nil, fmt.Errorf("unknown system architecture: %s", system.Architecture)
	}

	// Chipset
	switch system.Chipset {
	case domain.Chipset_CHIPSET_I440FX:
		os.Type.Machine = structure.OS_CHIPSET_I440FX
	case domain.Chipset_CHIPSET_Q35:
		os.Type.Machine = structure.OS_CHIPSET_Q35
	case domain.Chipset_CHIPSET_VIRT:
		os.Type.Machine = structure.OS_CHIPSET_VIRT
	default:
		return nil, fmt.Errorf("unknown system chipset: %s", system.Chipset)
	}

	// Firmware
	loaderPath := filepath.Join(g.granitRoot, "disk", firmware.LoaderDeviceId)
	templatePath := filepath.Join(g.granitRoot, "disk", firmware.TmplDeviceId)
	nvramPath := filepath.Join(g.granitRoot, "disk", firmware.NvramDeviceId)

	os.Loader.Data = loaderPath
	
	switch firmware.Firmware {
	case domain.Firmware_FIRMWARE_OVMF:
		os.Loader.MetaType = structure.OS_LOADER_PFLASH
		os.Nvram = &structure.OSNvram{
			MetaType: structure.OS_NVRAM_FILE,
			MetaTemplate: templatePath,
			Source: structure.OSNvramSource{MetaFile: nvramPath},
		}
	case domain.Firmware_FIRMWARE_SEABIOS:
		// TODO
	default:
		return nil, fmt.Errorf("unknown firmware type: %s", firmware.Firmware)
	}
	
	return os, nil
}
