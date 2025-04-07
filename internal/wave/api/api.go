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
  "log/slog"
	golog "log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"cthul.io/cthul/pkg/api/wave/v1"
	"cthul.io/cthul/pkg/api/wave/v1/wavev1connect"
)


type Endpoint struct {
	addr string
	tlsConfig *tls.Config
  logger *slog.Logger
	server *http.Server
}

type Option func(*Endpoint)

func NewEndpoint(logger *slog.Logger, addr string, cert tls.Certificate, opts ...Option) *Endpoint {
	mux := http.NewServeMux()
	mux.Handle(wavev1connect.NewDomainServiceHandler(&domainService{}))
	endpoint := &Endpoint{
		addr: addr,
		tlsConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		logger: logger.WithGroup("api-endpoint"),
		server: &http.Server{
			Handler: mux,
			ErrorLog: golog.New(io.Discard, "", 0),
			IdleTimeout: 10 * time.Minute,
		},
	}

	for _, opt := range opts {
		opt(endpoint)
	}

	return endpoint
}

// WithIdleTimeout sets a custom timeout for idle http connections.
func WithIdleTimeout(timeout time.Duration) Option {
	return func (e *Endpoint) {
		e.server.IdleTimeout = timeout
	}
}

// WithSkipInsecure enables skipping of insecure public certificates when mTLS is used.
func WithSkipInsecure(skip bool) Option {
	return func (e *Endpoint) {
		e.server.TLSConfig.InsecureSkipVerify = skip
	}
}

// WithSystemLog enables http system error logs and writes them to the specified logger.
// The logs are written as "error" with the category "api_server".
func WithSystemLog(logger log.Logger) Option {
	return func (e *Endpoint) {
		e.server.ErrorLog = golog.New(adapter.NewCommonLogAdapter("api_server", logger.Err), "", 0)
	}
}

// ServeAndDetach starts the api endpoint in a seperate goroutine and immediately returns.
// The server can be started only once.
func (e *Endpoint) ServeAndDetach() error {
	listener, err := tls.Listen("tcp", e.addr, e.tlsConfig)
	if err!=nil {
		return err
	}
	go func() {
		if err := e.server.Serve(listener); err!=nil {
			e.logger.Error(fmt.Sprintf("unrecoverable api error: %s", err.Error())) 
		}
	}()
	return nil
}

// Terminate tries to gracefully shutdown the api endpoint (waiting for connections to finish)
// if this fails or exceeds the provided context window, the connection is forcefully closed.
// If forcefully closing the connection fails too, an error is returned.
func (e *Endpoint) Terminate(ctx context.Context) error {
	if err := e.server.Shutdown(ctx); err!=nil {
		return e.server.Close()
	}
	return nil
}


type domainService struct{}

func (d *domainService) GetDomain(
	ctx context.Context,
	req *connect.Request[wavev1.GetDomainRequest],
) (*connect.Response[wavev1.GetDomainResponse], error) {
	
	return nil, fmt.Errorf("Not Implemented")
}
