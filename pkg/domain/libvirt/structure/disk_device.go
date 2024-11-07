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

type DISK_TYPE string

const (
	DISK_BLOCK   DISK_TYPE = "block"
)

type DISK_DEVICE_TYPE string

const (
	DISK_DEVICE_DISK DISK_DEVICE_TYPE = "disk"
	DISK_DEVICE_CDROM DISK_DEVICE_TYPE = "cdrom"
)

type Disk struct {
	XMLName xml.Name `xml:"disk"`
	MetaType DISK_TYPE `xml:"type,attr,omitempty"`
	MetaDevice DISK_DEVICE_TYPE `xml:"device,attr,omitempty"`
	Source *DiskSource `xml:"source,omitempty"`
	Driver *DiskDriver `xml:"driver,omitempty"`
	Target *DiskTarget `xml:"target,omitempty"`
	Readonly *DiskReadonly `xml:"readonly,omitempty"`
	Boot *Boot `xml:"boot,omitempty"`
}

type DiskSource struct {
	MetaDev string `xml:"dev,attr,omitempty"`
}

type DISK_DRIVER_TYPE string

const (
	DISK_DRIVER_QEMU DISK_DRIVER_TYPE = "qemu"
)

type DISK_STORAGE_TYPE string

const (
	DISK_STORAGE_RAW DISK_STORAGE_TYPE = "raw"
	DISK_STORAGE_QCOW2 DISK_STORAGE_TYPE = "qcow2"
)

type DISK_CACHE_TYPE string

const (
	DISK_CACHE_NONE DISK_CACHE_TYPE = "none"
	DISK_CACHE_WRITEBACK DISK_CACHE_TYPE = "writeback"
	DISK_CACHE_WRITETHROUGH = "writethrough"
)

type DiskDriver struct {
	MetaName DISK_DRIVER_TYPE `xml:"name,attr,omitempty"`
	MetaType DISK_STORAGE_TYPE `xml:"type,attr,omitempty"`
	MetaCache DISK_CACHE_TYPE `xml:"cache,attr,omitempty"`
}

type DISK_BUS_TYPE string

const (
	DISK_BUS_IDE DISK_BUS_TYPE = "ide"
	DISK_BUS_VIRTIO DISK_BUS_TYPE = "virtio" 
)

type DISK_TRAY_STATE string

const (
	DISK_TRAY_OPEN DISK_TRAY_STATE = "open"
	DISK_TRAY_CLOSED DISK_TRAY_STATE = "closed"
)

type DiskTarget struct {
	MetaBus DISK_BUS_TYPE `xml:"bus,attr,omitempty"`
	MetaTray DISK_TRAY_STATE `xml:"tray,attr,omitempty"`
}

type DiskReadonly struct {}
