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
	"context"
	"os"
	"runtime"
	"sync"
	"time"

	"cthul.io/cthul/pkg/log/structure"
)

// RuntimeLogger provides a logger implementation that sends logs async to a central worker.
// This is useful for runtime logs as slowdowns on the log destination do not directly impact the code.
type RuntimeLogger struct {
	// rootCtx is active for the full lifetime of the logger.
	// cancelling it will immediately shutdown the logger.
	rootCtx context.Context
	rootCtxCancel context.CancelFunc
	
	// workCtx is the context that is active for the runtime of the worker.
	// if it is called, no more logs can be send, and the remaining buffered logs can be flushed.
	workCtx context.Context
	workCtxCancel context.CancelFunc

	// finChan is used to send the absolute exist signal
	// if the channel emits, this indicates that the logger is fully cleaned up.
	finChan chan struct{}
	
	logQueue chan *structure.LogMessage
	level   structure.LEVEL
	component string
	trace   bool
	
	// ioLock is a reference to a shared lock that is used to synchronize stdout/err writes.
	// Some context about why this is used:
	// writing to stdout/err is thread safe for POSIX systems if the message is smaller then PIPE_BUF
	// this will be the case in 95% of the time. However to keep the package as compatible and deterministic
	// as possible, we use regular stdout/err writes (os independent) and lock it with this shared ioLock.
	// Using this lock is still highly performant as usually only one logger is running, which means the lock
	// is mostly uncontestet (on POSIX futex allows the uncontestet lock to be acquired without syscall).
	// If that would grow into a real performance issue, a custom POSIX implementation
	// can be written that checks if the buffer fits into PIPE_BUF and if yes omits the lock.
	ioLock *sync.Mutex
}

type RuntimeLoggerOption func(*RuntimeLogger)

// NewRuntimeLogger generates a RuntimeLogger. The specified component identifier is included in every log.
// The defaults can be overriden using RuntimeLoggerOptions.
func NewRuntimeLogger(component string, opts ...RuntimeLoggerOption) *RuntimeLogger {
	rootCtx, rootCtxCancel := context.WithCancel(context.Background())
	workCtx, workCtxCancel := context.WithCancel(rootCtx)
	logger := &RuntimeLogger{
		rootCtx: rootCtx,
		rootCtxCancel: rootCtxCancel,
		workCtx: workCtx,
		workCtxCancel: workCtxCancel,
		finChan: make(chan struct{}),
		logQueue: make(chan *structure.LogMessage),
		level:   structure.INFO,
		component: component,
		trace:   false,
		ioLock: &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

// WithTrace enables tracing, which includes information about file + line.
func WithTrace(enable bool) RuntimeLoggerOption {
	return func (l *RuntimeLogger) {
		l.trace = enable
	}
}

// WithLevel sets a custom log level. All logs below the specified level are ignored.
// Defaults to 'info' if level is invalid.
func WithLevel(level string) RuntimeLoggerOption {
	return func (l *RuntimeLogger) {
		l.level = structure.Level(level)
	}
}

// WithLogBuffer configures a buffer for the logger.
func WithLogBuffer(buffersize int64) RuntimeLoggerOption {
	return func(l *RuntimeLogger) {
		l.logQueue = make(chan *structure.LogMessage, buffersize)
	}
}

// WithIOLock uses a shared ioLock for stdout/err write operations.
// This must be used if the application uses multiple loggers at the same time.
func WithIOLock(lock *sync.Mutex) RuntimeLoggerOption {
	return func(l *RuntimeLogger) {
		l.ioLock = lock
	}
}

// ServeAndDetach starts the runtime log worker in a new goroutine and immediately returns.
// The worker can be started only once.
func (r *RuntimeLogger) ServeAndDetach() {
	go func() {
		for {
			// Select is not fair and could potentially starve the workCtx.Done().
			// To avoid this, the workCtx is actively rechecked on every iteration.
			select {
			case <-r.workCtx.Done():
				r.flushLogs()
				r.finChan<-struct{}{}
				return
			default:
			}
			
			select {
			case msg := <-r.logQueue:
				r.writeLog(msg)
			case <-r.workCtx.Done():
				// This trigger could be starved, therefore it just signals a reiteration
				// which actively checks the workCtx.
				break
			}
		}
	}()
}

// Terminate stops reading logs and tries to flush out the remaining buffered logs.
// If the provided context exceeds while flushing, flushing is cancelled immediately.
// Function will never return an error, it uses the error to adhere the cthul 'Terminate()' semantics.
func (r *RuntimeLogger) Terminate(ctx context.Context) error {
	r.workCtxCancel()
	defer r.rootCtxCancel()
	select {
	case <-r.finChan:
		return nil
	case <-ctx.Done():
		r.rootCtxCancel()
		<-r.finChan
		return nil
	}
}

// Crit sends a critical error to the logger.
func (r *RuntimeLogger) Crit(category string, message string) {
	r.sendLog(structure.CRITICAL, category, message)
}

// Err sends an error to the logger.
func (r *RuntimeLogger) Err(category string, message string) {
	r.sendLog(structure.ERROR, category, message)
}

// Warn sends a warning to the logger.
func (r *RuntimeLogger) Warn(category string, message string) {
	r.sendLog(structure.WARNING, category, message)
}

// Info sends an information message to the logger.
func (r *RuntimeLogger) Info(category string, message string) {
	r.sendLog(structure.INFO, category, message)
}

// Debug sends a debug message to the logger.
func (r *RuntimeLogger) Debug(category string, message string) {
	r.sendLog(structure.DEBUG, category, message)
}

// sendLog directly sends the log message to the worker queue.
// If the queue buffer is not full it immediately returns, otherwise it waits until a slot in the buffer is free.
func (r *RuntimeLogger) sendLog(level structure.LEVEL, category string, message string) {
	msg := &structure.LogMessage{
		Level:     level,
		Timestamp: time.Now().Unix(),
		Component: r.component,
		Category:  category,
		Message:   message,
		Trace:     nil,
	}
	
	if r.trace {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			msg.Trace = &structure.LogTrace{ File: file, Line: line }
		}
	}

	select {
	case <-r.workCtx.Done():
	case r.logQueue <- msg:
		return
	}
}


// writeLog writes the log message directly to the stdout (level below INFO to stderr).
func (r *RuntimeLogger) writeLog(msg *structure.LogMessage) {
	r.ioLock.Lock()
	defer r.ioLock.Unlock()
	if msg.Level < structure.INFO {
		os.Stderr.Write(msg.Serialize())
	} else {
		os.Stdout.Write(msg.Serialize())
	}
}

// flushLogs takes the current length 'n' of the logQueue and then writes 'n' logs down.
// If the context derived from rootCtx expires it returns after the current log is written.
func (r *RuntimeLogger) flushLogs() {
	flushCtx, flushCtxCancel := context.WithCancel(r.rootCtx)
	defer flushCtxCancel()
	
	flushCount := len(r.logQueue)
	for i:=0; i < flushCount; i++ {
		select {
		case <-flushCtx.Done():
			return
		default:
			r.writeLog(<-r.logQueue)
		}
	}
}
