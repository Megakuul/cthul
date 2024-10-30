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

type OS struct {
	Type *OSType `xml:"type,omitempty"`
	Loader *OSLoader `xml:"loader,omitempty"`
	Nvram *OSNvram `xml:"nvram,omitempty"`
	Nvrams *OSNvram `xml:"nvrams,omitempty"`
}

type OS_ARCH string
const (
	X86_64 OS_ARCH = "x86_64"
	AARCH64 OS_ARCH = "aarch64"
)

type OS_CHIPSET string
const (
	I440FX OS_CHIPSET = "pc"
	Q35 OS_CHIPSET = "q35"
	VIRT OS_CHIPSET = "virt"
)

type OSType struct {
	Arch OS_ARCH `xml:"arch,attr,omitempty"`
	Machine OS_CHIPSET `xml:"machine,attr,omitempty"`
	Data string `xml:",charset"`
}

type OS_LOADER_TYPE string
const (
	SEABIOS OS_LOADER_TYPE = "seabios"
	OVMF OS_CHIPSET = "pflash"
)

type OSLoader struct {
	MetaReadonly bool `xml:"readonly,attr,omitempty"`
	MetaSecure bool `xml:"secure,attr,omitempty"`
	MetaType OS_LOADER_TYPE `xml:"type,attr,omitempty"`
	Data string `xml:",charset"`
}

type OS_NVRAM_TYPE string
const (
	FILE OS_NVRAM_TYPE = "file"
)

type OSNvram struct {
	MetaType OS_NVRAM_TYPE `xml:"type,attr,omitempty"`
	MetaTemplate string `xml:"template,attr,omitempty"`
	Source OSNvramSource `xml:"source"`
}

type OSNvramSource struct {
	MetaFile string `xml:"file,attr,omitempty"`
}
