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

package domain

import (
	"context"

	"cthul.io/cthul/pkg/adapter/domain/structure"
)

// Domain provides the direct domain abstraction layer.
// It ensures that the underlying domain (vm) system can be replaced without much effort (even if not planned).
type DomainAdapter interface {
	// List returns a list with uuids from all domains on the host.
	ListDomains(context.Context) ([]string, error)
	// Apply updates the domain to the specified state. Updates that can be hotplugged are hotplugged, other
	// updates are applied at next reboot. Operation is idempotent.
	ApplyDomain(context.Context, string, structure.Domain) error
	// Destroy removes a domain from the local machine. Operation is idempotent.
	DestroyDomain(context.Context, string, structure.Domain) error
	// Start starts the domain.
	StartDomain(context.Context, string) error
	// Reboot reboots the domain if in running state.
	RebootDomain(context.Context, string) error
	// Pause freezes the domain state if in running state.
	PauseDomain(context.Context, string) error
	// Resume unfreezes the domain state if in paused state.
	ResumeDomain(context.Context, string) error
	// Shutdown stops the domain gracefully.
	ShutdownDomain(context.Context, string) error
	// Kill stops the domain forcefully.
	KillDomain(context.Context, string) error

	// GetDomainStats fetches overall domain stats. The domain is identified by uuid.
	GetDomainStats(context.Context, string) (*structure.DomainStats, error)
	// GetCpuStats fetches cpu stats of the domain. The domain is identified by uuid.
	GetCpuStats(context.Context, string) (*structure.CpuStats, error)
	// GetMemoryStats fetches memory stats of the domain. The domain is identified by uuid.
	GetMemoryStats(context.Context, string) (*structure.MemoryStats, error)
	// GetInterfaceStats fetches interface stats of the domain. The domain is identified by uuid.
	GetInterfaceStats(context.Context, string) (*structure.InterfaceStats, error)
	// GetBlockStats fetches block device stats of the domain. The domain is identified by uuid.
	GetBlockStats(context.Context, string) (*structure.BlockStats, error)
}

