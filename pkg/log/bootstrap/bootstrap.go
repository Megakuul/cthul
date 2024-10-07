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

package bootstrap

import (
	"os"
	"runtime"
	"sync"
	"time"

	"cthul.io/cthul/pkg/log/structure"
)

// BootstrapLogger provides a logger implementation that directly writes logs to stdout/err.
// This is useful for bootstrap or cleanup actions that should not rely on asynchron loggers.
type BootstrapLogger struct {
	level structure.LEVEL
	component string
	trace bool

	// ioLock is a reference to a shared lock that is used to synchronize stdout/err writes.
	// For more context about why this is used, check runtime/runtime.go.
	ioLock *sync.Mutex
}

type BootstrapLoggerOption func(*BootstrapLogger)

// NewBootstrapLogger generates a BootstrapLogger. The specified component identifier is included in every log.
// The defaults can be overriden using BootstrapLoggerOptions.
func NewBootstrapLogger(component string, opts ...BootstrapLoggerOption) *BootstrapLogger {
	logger := &BootstrapLogger{
		level: structure.INFO,
		component: component,
		trace: false,
		ioLock: &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

// WithTrace enables tracing, which includes information about file + line.
func WithTrace(enable bool) BootstrapLoggerOption {
	return func (l *BootstrapLogger) {
		l.trace = enable
	}
}

// WithLevel sets a custom log level. All logs below the specified level are ignored.
// Defaults to 'info' if level is invalid.
func WithLevel(level string) BootstrapLoggerOption {
	return func (l *BootstrapLogger) {
		l.level = structure.Level(level)
	}
}

// WithIOLock uses a shared ioLock for stdout/err write operations.
// This must be used if the application uses multiple loggers at the same time.
func WithIOLock(lock *sync.Mutex) BootstrapLoggerOption {
	return func(l *BootstrapLogger) {
		l.ioLock = lock
	}
}

// Crit sends a critical error to the logger.
func (b *BootstrapLogger) Crit(category string, message string) {
	b.writeLog(structure.CRITICAL, category, message)
}

// Err sends an error to the logger.
func (b *BootstrapLogger) Err(category string, message string) {
	b.writeLog(structure.ERROR, category, message)
}

// Warn sends a warning to the logger.
func (b *BootstrapLogger) Warn(category string, message string) {
	b.writeLog(structure.WARNING, category, message)
}

// Info sends an information message to the logger.
func (b *BootstrapLogger) Info(category string, message string) {
	b.writeLog(structure.INFO, category, message)
}

// Debug sends a debug message to the logger.
func (b *BootstrapLogger) Debug(category string, message string) {
	b.writeLog(structure.DEBUG, category, message)
}

// writeLog writes the specified log directly to the stdout (if below INFO it writes to stderr).
func (b *BootstrapLogger) writeLog(level structure.LEVEL, category string, message string) {
	if level > b.level {
		return
	}

	msg := &structure.LogMessage{
		Level: level,
		Timestamp: time.Now().Unix(),
		Component: b.component,
		Category: category,
		Message: message,
		Trace: nil,
	}

	if b.trace {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			msg.Trace = &structure.LogTrace{ File: file, Line: line	}
		}
	}

	b.ioLock.Lock()
	defer b.ioLock.Unlock()
	if level < structure.INFO {
		os.Stderr.Write(msg.Serialize())
	} else {
		os.Stdout.Write(msg.Serialize())
	}
}
