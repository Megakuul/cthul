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

package db

import "context"

// Client provides a database abstraction layer.
// It ensures that the underlying database can be replaced without much effort (even if not planned).
type Client interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string) error
	GetRange(context.Context, string, string) ([]string, error)
	SetRange(context.Context, []string) error
	Watch(context.Context, func(string, error)) error 
}
