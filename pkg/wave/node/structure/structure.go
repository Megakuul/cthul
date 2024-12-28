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

type NODE_STATE string

const (
	NODE_HEALTHY          NODE_STATE = "healthy"
	NODE_DEGRADED       NODE_STATE = "degraded"
	NODE_MAINTENANCE       NODE_STATE = "maintenance"
)

// Node provides the structure of the node information present on the database.
type Node struct {
	Affinity []string `json:"affinity"`
	State NODE_STATE `json:"state"`
	AllocatedCpu float64 `json:"allocated_cpu"`
	AvailableCpu    float64        `json:"available_cpu"`
	AllocatedMemory int64 `json:"allocated_memory"`
	AvailableMemory int64 `json:"available_memory"`
}


