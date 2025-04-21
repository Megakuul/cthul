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

// domain adapter provides an abstraction layer for virtual machines.
// it is primarily used by wave to provision and control virtual machines 'domains'.
package domain

import (
	"context"

  "cthul.io/cthul/pkg/api/wave/v1/domain"
)

type Adapter interface {
	// List returns a map with uuids & name of all domains on the host.
	List(context.Context) (map[string]string, error)
	// GetStats fetches all domain stats directly from the underlying vmm.
	GetStats(context.Context, string) (*domain.DomainStats, error)
	// Apply updates the domain to the specified state. Updates that can be hotplugged are hotplugged, other
	// updates are applied at next reboot. Operation is idempotent.
	Apply(context.Context, string, *domain.DomainConfig) error
	// Destroy removes a domain from the local machine. Operation is idempotent.
	Destroy(context.Context, string, *domain.DomainConfig) error
	// Start starts the domain or resumes it if it was paused.
	Start(context.Context, string) error
	// Reboot reboots the domain if in running state.
	Reboot(context.Context, string) error
	// Pause freezes the domain state if in running state.
	Pause(context.Context, string) error
	// Shutdown stops the domain gracefully.
	Shutdown(context.Context, string) error
	// Kill stops the domain forcefully.
	Kill(context.Context, string) error
}

