/**
 * Cthul System
 *
 * Copyright (C) 2025 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	golog "log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"cthul.io/cthul/internal/wave/api/domain"
	"cthul.io/cthul/internal/wave/api/node"
	"cthul.io/cthul/internal/wave/api/serial"
	"cthul.io/cthul/internal/wave/api/video"
	"cthul.io/cthul/pkg/api/wave/v1/domain/domainconnect"
	"cthul.io/cthul/pkg/api/wave/v1/node/nodeconnect"
	"cthul.io/cthul/pkg/api/wave/v1/serial/serialconnect"
	"cthul.io/cthul/pkg/api/wave/v1/video/videoconnect"
	domctrl "cthul.io/cthul/pkg/wave/domain"
	nodectrl "cthul.io/cthul/pkg/wave/node"
	serialctrl "cthul.io/cthul/pkg/wave/serial"
	videoctrl "cthul.io/cthul/pkg/wave/video"
)

type Endpoint struct {
	addr       string
	tlsConfig  *tls.Config
	logger     *slog.Logger
	mux        *http.ServeMux
	serverLock sync.Mutex
	server     *http.Server
}

type Option func(*Endpoint)

func New(addr string, cert tls.Certificate, opts ...Option) *Endpoint {
	endpoint := &Endpoint{
		addr: addr,
		tlsConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		logger:     slog.Default().WithGroup("api-endpoint"),
		mux:        http.NewServeMux(),
		serverLock: sync.Mutex{},
		server:     nil,
	}

	for _, opt := range opts {
		opt(endpoint)
	}

	return endpoint
}

// WithLogger adds a custom slog instance for this endpoint.
func WithLogger(logger *slog.Logger) Option {
	return func(e *Endpoint) {
		e.logger = logger.WithGroup("api-endpoint")
	}
}

func WithDomain(controller *domctrl.Controller) Option {
	return func(e *Endpoint) {
		e.mux.Handle(domainconnect.NewDomainServiceHandler(domain.New(controller)))
	}
}

func WithVideo(controller *videoctrl.Controller) Option {
	return func(e *Endpoint) {
		e.mux.Handle(videoconnect.NewVideoServiceHandler(video.New(controller)))
	}
}

func WithSerial(controller *serialctrl.Controller) Option {
	return func(e *Endpoint) {
		e.mux.Handle(serialconnect.NewSerialServiceHandler(serial.New(controller)))
	}
}

func WithNode(controller *nodectrl.Controller) Option {
	return func(e *Endpoint) {
		e.mux.Handle(nodeconnect.NewNodeServiceHandler(node.New(controller)))
	}
}

// ServeAndDetach starts the api endpoint in a seperate goroutine and immediately returns.
// The server can be started only once.
func (e *Endpoint) ServeAndDetach() error {
	e.serverLock.Lock()
  defer e.serverLock.Unlock()
	if e.server != nil {
		return fmt.Errorf("server cannot be started twice")
	}
	e.server = &http.Server{
		Handler:     e.mux,
		ErrorLog:    golog.New(io.Discard, "", 0),
		IdleTimeout: 10 * time.Minute,
	}

	listener, err := tls.Listen("tcp", e.addr, e.tlsConfig)
	if err != nil {
		return err
	}

	go func() {
		if err := e.server.Serve(listener); err != nil {
			e.logger.Error(fmt.Sprintf("unrecoverable api error: %s", err.Error()))
		}
	}()
	return nil
}

// Terminate tries to gracefully shutdown the api endpoint (waiting for connections to finish)
// if this fails or exceeds the provided context window, the connection is forcefully closed.
// If forcefully closing the connection fails too, an error is returned.
func (e *Endpoint) Terminate(ctx context.Context) error {
  e.serverLock.Lock()
  defer e.serverLock.Unlock()
  if e.server==nil {
    return nil
  }
	if err := e.server.Shutdown(ctx); err != nil {
		return e.server.Close()
	}
	return nil
}
