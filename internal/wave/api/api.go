package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
	golog "log"

	"connectrpc.com/connect"
	"cthul.io/cthul/pkg/api/wave/v1"
	"cthul.io/cthul/pkg/api/wave/v1/wavev1connect"
	"cthul.io/cthul/pkg/log"
	"cthul.io/cthul/pkg/log/adapter"
	"cthul.io/cthul/pkg/log/discard"
)


type ApiEndpoint struct {
	server *http.Server
	logger log.Logger
}

type ApiEndpointOption func(*ApiEndpoint)

func NewApiEndpoint(addr string, tls *tls.Config, logger log.Logger, opts ...ApiEndpointOption) *ApiEndpoint {
	mux := http.NewServeMux()
	mux.Handle(wavev1connect.NewDomainServiceHandler(&domainService{}))
	return &ApiEndpoint{
		logger: discard.NewDiscardLogger(),
		server: &http.Server{
			Addr: addr,
			TLSConfig: tls,
			Handler: mux,
			ErrorLog: golog.New(io.Discard, "", 0),
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
func WithApplicationLog(logger log.Logger) ApiEndpointOption {
	return func (a *ApiEndpoint) {
		a.logger = logger
	}
}

// WithSystemLog enables http system error logs and writes them to the specified logger.
// The logs are written as "error" with the category "api_http_logs".
func WithSystemLog(logger log.Logger) ApiEndpointOption {
	return func (a *ApiEndpoint) {
		a.server.ErrorLog = golog.New(adapter.NewLogAdapter("api_http", logger.Err), "", 0)
	}
}

// ServeAndDetach starts the api endpoint in a seperate goroutine and immediately returns.
// The server can be started only once.
func (a *ApiEndpoint) ServeAndDetach() error {
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

func (d *domainService) GetDomain(
	ctx context.Context,
	req *connect.Request[wavev1.GetDomainRequest],
) (*connect.Response[wavev1.GetDomainResponse], error) {
	
	return nil, fmt.Errorf("Not Implemented")
}
