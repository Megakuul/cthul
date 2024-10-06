package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"cthul.io/cthul/pkg/api/wave/v1"
	"cthul.io/cthul/pkg/api/wave/v1/wavev1connect"
)


type ApiEndpoint struct {
	server *http.Server
	logger *log.Logger
}

type ApiEndpointOption func(*ApiEndpoint)

func NewApiEndpoint(addr string, tls *tls.Config, opts ...ApiEndpointOption) *ApiEndpoint {
	mux := http.NewServeMux()
	mux.Handle(wavev1connect.NewDomainServiceHandler(&domainService{}))
	return &ApiEndpoint{
		logger: log.New(io.Discard, "", 0),
		server: &http.Server{
			Addr: addr,
			TLSConfig: tls,
			Handler: mux,
			ErrorLog: log.New(io.Discard, "", 0),
			IdleTimeout: 10 * time.Minute,
		},
	}
}

// WithIdleTimeout sets a custom timeout for idle http connections.
func WithIdleTimeout(timeout time.Duration) ApiEndpointOption {
	return func (a *ApiEndpoint) {
		a.server.IdleTimeout = timeout
	}
}

// WithApplicationLog enables api logs and writes them to the specified logger.
func WithApplicationLog(logger *log.Logger) ApiEndpointOption {
	return func (a *ApiEndpoint) {
		a.logger = logger
	}
}

// WithSystemLog enables http system logs and writes them to the specified logger.
func WithSystemLog(logger *log.Logger) ApiEndpointOption {
	return func (a *ApiEndpoint) {
		a.server.ErrorLog = logger
	}
}

// Serve starts the api endpoint and blocks until closed.
// The server can be started only once.
func (a *ApiEndpoint) Serve() error {
	return a.server.ListenAndServeTLS("", "")
}

// Close tries to gracefully close the api endpoint (waiting for connections to finish)
// if this fails or exceeds the provided context window, the connection is forcefully closed.
// If forcefully closing the connection fails too, an error is returned.
func (a *ApiEndpoint) Close(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err!=nil {
		return a.server.Close()
	}
	return nil
}


type domainService struct{}

func (d *domainService) GetDomain(ctx context.Context, req *wavev1.GetDomainRequest) (*wavev1.GetDomainResponse, error) {

	return nil, fmt.Errorf("Not Implemented")
}
