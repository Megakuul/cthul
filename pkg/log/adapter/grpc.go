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

import "fmt"

// GrpcLogAdapter provides a grpclog.LoggerV2 implementation that writes to a logger.
// This exists, because using grpclog.NewLoggerV2 adds a lot of additional information like timestamps
// that are redundant just generate noise. Currently this grpc behavior is hardcoded and can not be changed.
// Btw thanks to google for implementing the single worst logging approach the world has ever seen.
type GrpcLogAdapter struct {
	category string
	newLineChar string
	infoLogFunc func(string, string)
	warnLogFunc func(string, string)
	errLogFunc func(string, string)
	critLogFunc func(string, string)
}

type GrpcLogAdapterOption func(*GrpcLogAdapter)

// NewGrpcLogAdapter creates a new grpc log adapter.
func NewGrpcLogAdapter(category string, opts ...GrpcLogAdapterOption) *GrpcLogAdapter {
	adapter := &GrpcLogAdapter{
		category: category,
		newLineChar: "",
		infoLogFunc: func(_, _ string) {},
		warnLogFunc: func(_, _ string) {},
		errLogFunc: func(_, _ string) {},
		critLogFunc: func(_, _ string) {},
	}

	for _, opt := range opts {
		opt(adapter)
	}

	return adapter
}

// WithNewLineChar defines a custom string that is inserted instead of a newline.
func WithNewLineChar(newLineChar string) GrpcLogAdapterOption {
	return func (g *GrpcLogAdapter) {
		g.newLineChar = newLineChar
	}
}

// WithInfoLog adds a log function that is fired on info events.
// Important: there seems to be a fundamental difference between what grpc and cthul considers as INFO log.
// The grpc INFO logs fire every crap event; don't be surprised if it tells you what you ate for breakfast.
func WithInfoLog(infoLogFunc func(string, string)) GrpcLogAdapterOption {
	return func (g *GrpcLogAdapter) {
		g.infoLogFunc = infoLogFunc
	}
}

// WithWarnLog adds a log function that is fired on warning events.
func WithWarnLog(warnLogFunc func(string, string)) GrpcLogAdapterOption {
	return func (g *GrpcLogAdapter) {
		g.warnLogFunc = warnLogFunc
	}
}

// WithErrLog adds a log function that is fired on error events.
func WithErrLog(errLogFunc func(string, string)) GrpcLogAdapterOption {
	return func (g *GrpcLogAdapter) {
		g.errLogFunc = errLogFunc
	}
}

// WithCritLog adds a log function that is fired on fatal events.
func WithCritLog(critLogFunc func(string, string)) GrpcLogAdapterOption {
	return func (g *GrpcLogAdapter) {
		g.critLogFunc = critLogFunc
	}
}

// Not commenting the functions below, they simply implement any of the 420 function interfaces
// that underemployed grpc developers thought were necessary.

func (g *GrpcLogAdapter) Info(args ...interface{}) {
	g.infoLogFunc(g.category, sanitizeLog(fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Infoln(args ...interface{}) {
	g.infoLogFunc(g.category, sanitizeLog(g.newLineChar + fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Infof(format string, args ...interface{}) {
	g.infoLogFunc(g.category, sanitizeLog(fmt.Sprintf(format, args...)))
}

func (g *GrpcLogAdapter) Warning(args ...interface{}) {
	g.warnLogFunc(g.category, sanitizeLog(fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Warningln(args ...interface{}) {
	g.warnLogFunc(g.category, sanitizeLog(g.newLineChar + fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Warningf(format string, args ...interface{}) {
	g.warnLogFunc(g.category, sanitizeLog(fmt.Sprintf(format, args...)))
}

func (g *GrpcLogAdapter) Error(args ...interface{}) {
	g.errLogFunc(g.category, sanitizeLog(fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Errorln(args ...interface{}) {
	g.errLogFunc(g.category, sanitizeLog(g.newLineChar + fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Errorf(format string, args ...interface{}) {
	g.errLogFunc(g.category, sanitizeLog(fmt.Sprintf(format, args...)))
}

func (g *GrpcLogAdapter) Fatal(args ...interface{}) {
	g.critLogFunc(g.category, sanitizeLog(fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Fatalln(args ...interface{}) {
	g.critLogFunc(g.category, sanitizeLog(g.newLineChar + fmt.Sprint(args...)))
}

func (g *GrpcLogAdapter) Fatalf(format string, args ...interface{}) {
	g.critLogFunc(g.category, sanitizeLog(fmt.Sprintf(format, args...)))
}

func (g *GrpcLogAdapter) V(_ int) bool {
	// Verbosity is managed by the cthul logger, therefore verbosity level always mets the requested level.
	return true
}
