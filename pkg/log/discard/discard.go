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

package discard

// DiscardLogger is a logger implementation that discards all written logs immediately.
// This can be used as default value (if the default behavior should be to NOT log).
type DiscardLogger struct {}

// NewDiscardLogger creates a new discard logger.
func NewDiscardLogger() *DiscardLogger {
	return &DiscardLogger{}
}

// Crit sends a critical error to the logger.
func (d *DiscardLogger) Crit(category string, message string) {
	return
}

// Err sends an error to the logger.
func (d *DiscardLogger) Err(category string, message string) {
	return
}

// Warn sends a warning to the logger.
func (d *DiscardLogger) Warn(category string, message string) {
	return
}

// Info sends an information message to the logger.
func (d *DiscardLogger) Info(category string, message string) {
	return
}

// Debug sends a debug message to the logger.
func (d *DiscardLogger) Debug(category string, message string) {
	return
}

