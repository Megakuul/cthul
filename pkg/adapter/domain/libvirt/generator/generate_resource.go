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
	cthulstruct "cthul.io/cthul/pkg/adapter/domain/structure"
	libvirtstruct "cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
)

// generateVCPU generates libvirt VCPUs from resource configuration.
func (l *Generator) generateVCPU(resource *cthulstruct.ResourceConfig) *libvirtstruct.VCPU {
	return &libvirtstruct.VCPU{
		MetaPlacement: libvirtstruct.CPU_PLACEMENT_STATIC,
		Data: resource.VCPUs,
	}
}

// generateMemory generates libvirt Memory from resource configuration.
func (l *Generator) generateMemory(resource *cthulstruct.ResourceConfig) *libvirtstruct.Memory {
	return &libvirtstruct.Memory{
		MetaUnit: libvirtstruct.MEMORY_UNIT_BYTES,
		Data: resource.Memory,
	}
}
