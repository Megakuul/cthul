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

package adapter

// LogAdapter provides a io.Writer implementation that writes to a Logger.
type LogAdapter struct {
	category string
	logFunc func(string, string)
}

// NewLogAdapter creates a new log adapter, the adapter uses the provided log function
// to write incomming events.
func NewLogAdapter(category string, logFunc func(string, string)) *LogAdapter {
	return &LogAdapter{
		category: category,
		logFunc: logFunc,
	}
}

// Write implements the io.Writer interface and is used to write to the adapter.
// It will never return an error, the return types are just in place for to satisfy the interface.
func (l *LogAdapter) Write(input []byte) (int, error) {
	l.logFunc(l.category, string(input))
	return len(input), nil
}
