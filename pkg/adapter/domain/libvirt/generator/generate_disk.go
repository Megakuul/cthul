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

	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
  diskstruct "cthul.io/cthul/pkg/granit/disk/structure"
)

// Explanation: A libvirt disk device represents a device which is provided to the guest os via a custom
// bus and managed via a host driver (e.g. qemu). The host driver is invoked whenever the cpu's
// virtualization-extension detects writes to either PMIO or MMIO regions of a device, the driver will then
// essentially emulating the devices behavior. With a IDE or SATA bus this may involve writing data from the
// source (e.g. a block device on the host) to the emulated DMA buffer on the guest.
// With the VIRTIO bus the device writes to a used virtqueue.

// The source (e.g. block device) is provided by granit, granit manages the underlying replication.
// For the host driver this feels like reading/writing to a regular source (e.g. a block device).

// generateDisk generates a libvirt disk device from the cthul storage device.
func (g *Generator) generateDisk(ctx context.Context, device *cthulstruct.StorageDevice) (*libvirtstruct.Disk, error) {
	disk := &libvirtstruct.Disk{
		Source: &libvirtstruct.DiskSource{},
		Driver: &libvirtstruct.DiskDriver{MetaName: libvirtstruct.DISK_DRIVER_QEMU},
		Target: &libvirtstruct.DiskTarget{},
		Boot: &libvirtstruct.Boot{MetaOrder: device.BootPriority},
	}

	storageDevice, err := g.disk.Lookup(ctx, device.DeviceId)
	if err!=nil {
		return nil, err
	}

	// Metadata
	switch device.StorageType {
	case cthulstruct.STORAGE_DISK:
		disk.MetaDevice = libvirtstruct.DISK_DEVICE_DISK
	case cthulstruct.STORAGE_CDROM:
		disk.MetaDevice = libvirtstruct.DISK_DEVICE_CDROM
	default:
		return nil, fmt.Errorf("unknown storage type: %s", device.StorageType)
	}

	if storageDevice.Readonly {
		disk.Readonly = &libvirtstruct.DiskReadonly{}
	}

	// Target (on guest)
	switch device.StorageBus {
	case cthulstruct.STORAGE_IDE:
		disk.Target.MetaBus = libvirtstruct.DISK_BUS_IDE
	case cthulstruct.STORAGE_SATA:
		disk.Target.MetaBus = libvirtstruct.DISK_BUS_SATA
	case cthulstruct.STORAGE_VIRTIO:
		disk.Target.MetaBus = libvirtstruct.DISK_BUS_VIRTIO
	default:
		return nil, fmt.Errorf("unknown storage bus: %s", device.StorageBus)
	}

	// Driver (transfering data between guest and host)
	switch storageDevice.Format {
	case diskstruct.DISK_RAW:
		disk.Driver.MetaType = libvirtstruct.DISK_STORAGE_RAW
	case diskstruct.DISK_QCOW2:
		disk.Driver.MetaType = libvirtstruct.DISK_STORAGE_QCOW2
	default:
		return nil, fmt.Errorf("unsupported device format: %s", storageDevice.Format)
	}

	// Source (on host)
	switch storageDevice.Type {
	case diskstruct.DISK_BLOCK:
		disk.MetaType = libvirtstruct.DISK_BLOCK
		disk.Source = &libvirtstruct.DiskSource{
			MetaDev: storageDevice.Path,
		}
	case diskstruct.DISK_FILE:
		disk.MetaType = libvirtstruct.DISK_FILE
		disk.Source = &libvirtstruct.DiskSource{
			MetaFile: storageDevice.Path,
		}
	default:
		return nil, fmt.Errorf("unsupported device type: %s", storageDevice.Type)
	}

	return disk, nil
}
