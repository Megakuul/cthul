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

import "context"

// DomainController provides a domain abstraction layer.
// It ensures that the underlying domain (vm) system can be replaced without much effort (even if not planned).
type DomainController interface {
	// Apply updates the domain to the specified state. Updates that can be hotplugged are hotplugged, other
	// updates are applied at next reboot. Operation is idempotent.
	Apply(context.Context, Domain) error
	// Destroy removes a domain from the local machine. The domain is identified by uuid. Operation is idempotent.
	Destroy(context.Context, string) error

	// Start starts the domain.
	Start(string) error
	// Reboot reboots the domain if in running state.
	Reboot(string) error
	// Pause freezes the domain state if in running state.
	Pause(string) error
	// Resume unfreezes the domain state if in paused state.
	Resume(string) error
	// Shutdown stops the domain gracefully.
	Shutdown(string) error
	// Kill stops the domain forcefully.
	Kill(string) error

	// CreateSnapshot creates a domain snapshot based on the specified config.
	CreateSnapshot(context.Context, Snapshot) error
	// RevertSnapshot reverts the domain to a previous snapshot. The snapshot is identified by uuid.
	RevertSnapshot(context.Context, string) error
	// ConsolidateSnapshot consolidates the specified snapshots into the base image.
	// The snapshot is identified by uuid. Operation is idempotent.
	ConsolidateSnapshot(context.Context, string) error

	// GetTextConsole starts a tty console session to the domain. Returns a send and recv channel, closing both
	// channels deallocates the session. Data is transfered in raw tty chunks and must be handled manually.
	GetTextConsole(context.Context, string) (chan<-[]byte, <-chan []byte, error)
	// GetSpiceConsole starts a spice session to the domain. Returns a send and recv channel, closing both
	// channels deallocates the session. Data is transfered in raw spice chunks and must be handled manually.
	GetSpiceConsole(context.Context, string) (chan<-[]byte, <-chan []byte, error)

	// GetDomainStats fetches overall domain stats. The domain is identified by uuid.
	GetDomainStats(context.Context, string) (DomainStats, error)
	// GetCpuStats fetches cpu stats of the domain. The domain is identified by uuid.
	GetCpuStats(context.Context, string) (CpuStats, error)
	// GetMemoryStats fetches memory stats of the domain. The domain is identified by uuid.
	GetMemoryStats(context.Context, string) (MemoryStats, error)
	// GetInterfaceStats fetches interface stats of the domain. The domain is identified by uuid.
	GetInterfaceStats(context.Context, string) (InterfaceStats, error)
	// GetBlockStats fetches block device stats of the domain. The domain is identified by uuid.
	GetBlockStats(context.Context, string) (BlockStats, error)
}
