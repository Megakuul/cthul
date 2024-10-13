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

// CommonLogAdapter provides a io.Writer implementation that writes to a Logger.
type CommonLogAdapter struct {
	category string
	logFunc func(string, string)
}

// NewCommonLogAdapter creates a new common log adapter, the adapter uses the provided log function
// to write incomming events.
func NewCommonLogAdapter(category string, logFunc func(string, string)) *CommonLogAdapter {
	return &CommonLogAdapter{
		category: category,
		logFunc: logFunc,
	}
}

// Write implements the io.Writer interface and is used to write to the adapter.
// It will never return an error, the return types are just in place for to satisfy the interface.
func (l *CommonLogAdapter) Write(input []byte) (int, error) {
	l.logFunc(l.category, sanitizeLog(string(input)))
	return len(input), nil
}
