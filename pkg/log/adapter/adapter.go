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

import (
	"strconv"
)

// sanitizeLog removes escape characters that can break the loggers json representation.
// This is not very efficient but ensures consistent logging for adapters that use external formats.
func sanitizeLog(log string) string {
	// currently Quote covers all characters that break the json.
	// this function is employed to add other sanitization or faster implementations in the future.
	sanitizedLog := strconv.Quote(log)
	return sanitizedLog[1:len(sanitizedLog)-1] // trim the '"' added by Quote.
}
