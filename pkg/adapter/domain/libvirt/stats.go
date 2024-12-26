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

package libvirt

import (
	"context"
	"fmt"

	"cthul.io/cthul/pkg/domain/structure"
)

func (l *LibvirtAdapter)GetDomainStats(ctx context.Context, id string) (*structure.DomainStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetCpuStats(ctx context.Context, id string) (*structure.CpuStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetMemoryStats(ctx context.Context, id string) (*structure.MemoryStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetInterfaceStats(ctx context.Context, id string) (*structure.InterfaceStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}

func (l *LibvirtAdapter)GetBlockStats(ctx context.Context, id string) (*structure.BlockStats, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}
