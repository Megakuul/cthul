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

	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
	diskstruct "cthul.io/cthul/pkg/api/granit/v1/disk"
	"cthul.io/cthul/pkg/api/wave/v1/domain"
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
func (g *Generator) generateDisk(ctx context.Context, device *domain.StorageDevice) (*structure.Disk, error) {
	disk := &structure.Disk{
		Source: &structure.DiskSource{},
		Driver: &structure.DiskDriver{MetaName: structure.DISK_DRIVER_QEMU},
		Target: &structure.DiskTarget{},
		Boot:   &structure.Boot{MetaOrder: device.BootPriority},
	}

	storageDevice, err := g.disk.Lookup(ctx, device.DeviceId)
	if err != nil {
		return nil, err
	}

	// Metadata
	switch device.StorageType {
	case domain.StorageType_STORAGE_TYPE_DISK:
		disk.MetaDevice = structure.DISK_DEVICE_DISK
	case domain.StorageType_STORAGE_TYPE_CDROM:
		disk.MetaDevice = structure.DISK_DEVICE_CDROM
	default:
		return nil, fmt.Errorf("unknown storage type: %s", device.StorageType)
	}

	if storageDevice.Config.Readonly {
		disk.Readonly = &structure.DiskReadonly{}
	}

	// Target (on guest)
	switch device.StorageBus {
	case domain.StorageBus_STORAGE_BUS_IDE:
		disk.Target.MetaBus = structure.DISK_BUS_IDE
	case domain.StorageBus_STORAGE_BUS_SATA:
		disk.Target.MetaBus = structure.DISK_BUS_SATA
	case domain.StorageBus_STORAGE_BUS_VIRTIO:
		disk.Target.MetaBus = structure.DISK_BUS_VIRTIO
	default:
		return nil, fmt.Errorf("unknown storage bus: %s", device.StorageBus)
	}

	// Driver (transfering data between guest and host)
	switch storageDevice.Config.Format {
	case diskstruct.DiskFormat_DISK_FORMAT_RAW:
		disk.Driver.MetaType = structure.DISK_STORAGE_RAW
	case diskstruct.DiskFormat_DISK_FORMAT_QCOW2:
		disk.Driver.MetaType = structure.DISK_STORAGE_QCOW2
	default:
		return nil, fmt.Errorf("unsupported device format: %s", storageDevice.Config.Format)
	}

	// Source (on host)
	disk.MetaType = structure.DISK_BLOCK
	disk.Source = &structure.DiskSource{
		MetaDev: storageDevice.Config.Path,
	}

	return disk, nil
}
