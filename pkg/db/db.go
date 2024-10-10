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
	// Get returns a single value from database. Returns "" if the key is emtpy OR does not exist.
	Get(context.Context, string) (string, error)
	// GetRange returns a map of kvs based on the provided prefix.
	GetRange(context.Context, string) (map[string]string, error)
	// Set upserts a kv with the specified ttl. If ttl is 0 the kv does not expire.
	Set(context.Context, string, string, int64) error
	// Delete removes a kv from the database.
	Delete(context.Context, string) error
	// DeleteRange removes all kvs from the database by prefix.
	DeleteRange(context.Context, string) error
	// Watch calls the specified function on every update of the specified key.
	// The callback provides the updated key, value and an error in case of a failure.
	// The function is blocking, to stop it cancel the context.
	Watch(context.Context, string, func(string, string, error)) error
	// WatchRange calls the specified function on every update of a key in the prefix range.
	// The callback provides the updated key, value and an error in case of a failure.
	// The function is blocking, to stop it cancel the context.
	WatchRange(context.Context, string, func(string, string, error)) error
}
