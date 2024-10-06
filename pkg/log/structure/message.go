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

import (
	"bytes"
	"strconv"
)

// LogMessage represents one single log event / entry.
type LogMessage struct {
	Level LEVEL
	Timestamp int64
	Category string
	Message string
	Trace *LogTrace // If trace==nil it is not further processed.
}

// Serialize uses a hardcoded high performance json serializer
// to convert the LogMessage into a structured json string.
func (l *LogMessage) Serialize() []byte {
	// Function code looks terrible, but in this one situation it is fine
	// because it gains a notable performance boost as there is only 1 heap allocation (2 with trace).
	// Using reflection would lead to many micro allocations which can starve performance.
	estimatedCapacity := len(l.Category) + len(l.Message) + 256
	buffer := bytes.NewBuffer(make([]byte, 0, estimatedCapacity))

	// object prefix
	buffer.WriteByte(byte('{'))

	// field "timestamp"
	buffer.WriteString(`"timestamp":"`)
	buffer = bytes.NewBuffer(strconv.AppendInt(buffer.Bytes(), l.Timestamp, 10))
	buffer.WriteString(`", `)
	
	// field "level"
	buffer.WriteString(`"level":"`)
	buffer.WriteString(l.Level.String())
	buffer.WriteString(`", `)

	// field "category"
	buffer.WriteString(`"category":"`)
	buffer.WriteString(l.Category)
	buffer.WriteString(`", `)

	// field "message"
	buffer.WriteString(`"message":"`)
	buffer.WriteString(l.Message)
	buffer.WriteString(`"`)

	// field "trace"
	if l.Trace!=nil {
		buffer.WriteString(`, "trace":`)
		buffer.Write(l.Trace.serialize())
		buffer.WriteString(`"`)
	}

	// object suffix
	buffer.WriteByte(byte('}'))

	return buffer.Bytes()
}

// LogTrace holds trace information about the log.
type LogTrace struct {
	Line int
	File string
}

// Serialize uses a hardcoded high performance json serializer
// to convert the LogTrace into a structured json string.
func (l *LogTrace) serialize() []byte {
	estimatedCapacity := len(l.File) + 128
	buffer := bytes.NewBuffer(make([]byte, 0, estimatedCapacity))

	// object prefix
	buffer.WriteByte(byte('{'))

	// field "file"
	buffer.WriteString(`"file": "`)
	buffer.WriteString(l.File)
	buffer.WriteString(`", `)

	// field "line"
	buffer.WriteString(`"line": "`)
	buffer = bytes.NewBuffer(strconv.AppendInt(buffer.Bytes(), int64(l.Line), 10))
	buffer.WriteString(`"`)

	// object suffix
	buffer.WriteByte(byte('}'))

	return buffer.Bytes()
}

