package bootstrap

import "cthul.io/cthul/pkg/log/structure"

type BootstrapLogger struct {
	level structure.LEVEL
	service string
	trace bool
}

type BootstrapLoggerOption func(*BootstrapLogger)

func NewBootstrapLogger(service string, opts ...BootstrapLoggerOption) *BootstrapLogger {
	return &BootstrapLogger{
		level: structure.INFO,
		service: service,
		trace: false,
	}
}

// WithTrace enables tracing, which includes information about file + line.
func WithTrace() BootstrapLoggerOption {
	return func (l *BootstrapLogger) {
		l.trace = true
	}
}

// WithLevel sets a custom log level. All logs below the specified level are ignored.
// Defaults to 'info' if level is invalid.
func WithLevel(level string) BootstrapLoggerOption {
	return func (l *BootstrapLogger) {
		l.level = structure.Level(level)
	}
}

func (l *Logger) Log(level LEVEL, category string, message string) {
	msg := &bootstrapLogMessage{
		level: level,
		timestamp: time.Now().Unix(),
		category: &category,
		message: &message,
		trace: nil, 
	}
	if l.trace {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg.trace = &logTrace{ line: line, file: file }
		}
	}
	l.logChan <- msg
}
