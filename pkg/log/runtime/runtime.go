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

package runtime

import (
	"runtime"
	"time"

	"cthul.io/cthul/pkg/log/structure"
)

type RuntimeLogger struct {
	logChan chan *structure.LogMessage
	level   structure.LEVEL
	service string
	trace   bool
}

type RuntimeLoggerOption func(*RuntimeLogger)

func NewLogger(level string, service string, opts ...RuntimeLoggerOption) *RuntimeLogger {
	return &RuntimeLogger{
		logChan: make(chan *structure.LogMessage),
		level:   structure.INFO,
		service: service,
		trace:   false,
	}
}

// WithTrace enables tracing, which includes information about file + line.
func WithTrace() RuntimeLoggerOption {
	return func(l *RuntimeLogger) {
		l.trace = true
	}
}

// WithLogBuffer configures a buffer for the logger.
func WithLogBuffer(buffersize int64) RuntimeLoggerOption {
	return func(l *RuntimeLogger) {
		l.logChan = make(chan *structure.LogMessage, buffersize)
	}
}

func (l *RuntimeLogger) Log(level LEVEL, category string, message string) {
	msg := &logMessage{
		level:     level,
		timestamp: time.Now().Unix(),
		category:  &category,
		message:   &message,
		trace:     nil,
	}
	if l.trace {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg.trace = &logTrace{line: line, file: file}
		}
	}
	l.logChan <- msg
}
