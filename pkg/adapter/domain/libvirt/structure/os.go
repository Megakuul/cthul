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

type OS_FIRMWARE string

const (
	OS_FIRMWARE_EFI  OS_FIRMWARE = "efi"
	OS_FIRMWARE_BIOS OS_FIRMWARE = "bios"
)

type OS struct {
	MetaFirmware OS_FIRMWARE `xml:"firmware,attr,omitempty"`
	Type         *OSType     `xml:"type,omitempty"`
	Loader       *OSLoader   `xml:"loader,omitempty"`
	Nvram        *OSNvram    `xml:"nvram,omitempty"`
	Boot         *OSBoot     `xml:"boot,omitempty"`
	Smbios       *OSSMBios   `xml:"smbios,omitempty"`
	Bios         *OSBios     `xml:"bios,omitempty"`
}

type OS_ARCH string

const (
	OS_ARCH_X86_64  OS_ARCH = "x86_64"
	OS_ARCH_AARCH64 OS_ARCH = "aarch64"
)

type OS_CHIPSET string

const (
	OS_CHIPSET_I440FX OS_CHIPSET = "pc"
	OS_CHIPSET_Q35    OS_CHIPSET = "q35"
	OS_CHIPSET_VIRT   OS_CHIPSET = "virt"
)

type OSType struct {
	Arch    OS_ARCH    `xml:"arch,attr,omitempty"`
	Machine OS_CHIPSET `xml:"machine,attr,omitempty"`
	Data    string     `xml:",chardata"`
}

type OS_LOADER_TYPE string

const (
	OS_LOADER_ROM    OS_LOADER_TYPE = "rom"
	OS_LOADER_PFLASH OS_LOADER_TYPE = "pflash"
)

type OSLoader struct {
	MetaReadonly bool           `xml:"readonly,attr,omitempty"`
	MetaSecure   bool           `xml:"secure,attr,omitempty"`
	MetaType     OS_LOADER_TYPE `xml:"type,attr,omitempty"`
	Data         string         `xml:",chardata"`
}

type OS_NVRAM_TYPE string

const (
	OS_NVRAM_FILE OS_NVRAM_TYPE = "file"
)

type OSNvram struct {
	MetaType     OS_NVRAM_TYPE `xml:"type,attr,omitempty"`
	MetaTemplate string        `xml:"template,attr,omitempty"`
	Source       OSNvramSource `xml:"source"`
}

type OSNvramSource struct {
	MetaFile string `xml:"file,attr,omitempty"`
}

type OS_BOOT_OPTION string

const (
	OS_BOOT_HD      OS_BOOT_OPTION = "hd"
	OS_BOOT_CDROM   OS_BOOT_OPTION = "cdrom"
	OS_BOOT_NETWORK OS_BOOT_OPTION = "network"
)

type OSBoot struct {
	MetaDev OS_BOOT_OPTION `xml:"dev,attr,omitempty"`
}


type OSSMBios struct {
	MetaMode string `xml:"mode,attr,omitempty"`
}

type OSBios struct {
	MetaUseserial     string `xml:"useserial,attr,omitempty"`
	MetaRebootTimeout int    `xml:"rebootTimeout,attr,omitempty"`
}
