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
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package structure

import (
  adapterstruct "cthul.io/cthul/pkg/adapter/domain/structure"
)

// Domain "should've" state used by cthul.
// It's not representing the actual state but rather the state the domain should be pushed by the operator.
type DOMAIN_STATE string

const (
	DOMAIN_UP          DOMAIN_STATE = "up"
	DOMAIN_PAUSE       DOMAIN_STATE = "pause"
	DOMAIN_DOWN        DOMAIN_STATE = "down"
	DOMAIN_FORCED_DOWN DOMAIN_STATE = "forced_down"
)

// Domain provides the structure of the domain information present on the database.
// Unlike the Domain adapter config, this is not a deterministic blueprint but rather a snapshot
// of the current state of the domain. The error field contains parsing failures of single fields; if the error
// is not nil, the values may be used for informational purposes but should not be relied on.
type Domain struct {
	Reqnode         string                `json:"reqnode"`
	Node            string                `json:"node"`
	Config          *adapterstruct.Domain `json:"config"` // core
	State           DOMAIN_STATE          `json:"state"`
	Affinity        []string              `json:"affinity"`
	AllocatedCPU    float64               `json:"allocated_cpu"`
	AllocatedMemory int64                 `json:"allocated_memory"`
	Error           error                 `json:"-"`
}
