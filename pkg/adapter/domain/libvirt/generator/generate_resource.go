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
	"cthul.io/cthul/pkg/adapter/domain/libvirt/structure"
  "cthul.io/cthul/pkg/api/wave/v1/domain"
)

// generateVCPU generates libvirt VCPUs from resource configuration.
func (l *Generator) generateVCPU(resource *domain.ResourceConfig) *structure.VCPU {
	return &structure.VCPU{
		MetaPlacement: structure.CPU_PLACEMENT_STATIC,
		Data: resource.Vcpus,
	}
}

// generateMemory generates libvirt Memory from resource configuration.
func (l *Generator) generateMemory(resource *domain.ResourceConfig) *structure.Memory {
	return &structure.Memory{
		MetaUnit: structure.MEMORY_UNIT_BYTES,
		Data: resource.Memory,
	}
}
