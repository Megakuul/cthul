/**
 * Cthul System
 *
 * Copyright (C) 2025 Linus Ilian Moser <linus.moser@megakuul.ch>
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
	"cthul.io/cthul/pkg/granit/disk"
	"cthul.io/cthul/pkg/proton/inter"
	"cthul.io/cthul/pkg/wave/serial"
	"cthul.io/cthul/pkg/wave/video"
)

// Generator provides operations to generate drbd resource configurations 
// from granit disk and cluster definitions.
type Generator struct {
  deviceRoot string
}

type Option func(*Generator)

func New(opts ...Option) *Generator {

	generator := &Generator{
    deviceRoot: "/dev/cthul/granit/",
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

// Generate generates a resource configuration for one device from teh specified structs.
func (g *Generator) Generate(ctx context.Context, id string, config *domain.DomainConfig) (*structure.Domain, error) {
}
