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

package structure

import "strings"

// LEVEL is used to indicate the severity of a log message.
type LEVEL int64
// External constants holding the LEVEL representation of the levels.
const (
	CRITICAL LEVEL = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

// Internal constants holding the string representation of the levels.
const (
	sCRITICAL string = "critical"
	sERROR    string = "error"
	sWARNING  string = "warning"
	sINFO     string = "info"
	sDEBUG    string = "debug"
)

// Level returns the level representation of the loglevel.
// Defaults to INFO if the specified level is invalid.
func Level(name string) LEVEL {
	switch strings.ToLower(name) {
	case sCRITICAL:
		return CRITICAL
	case sERROR:
		return ERROR
	case sWARNING:
		return WARNING
	case sINFO:
		return INFO
	case sDEBUG:
		return DEBUG
	default:
		return INFO
	}
}

// String returns the string representation of the loglevel.
func (l *LEVEL) String() string {
	switch *l {
	case CRITICAL:
		return sCRITICAL
	case ERROR:
		return sERROR
	case WARNING:
		return sWARNING
	case INFO:
		return sINFO
	case DEBUG:
		return sDEBUG
	default:
		return ""
	}
}
