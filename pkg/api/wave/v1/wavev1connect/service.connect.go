// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: wave/v1/service.proto

package wavev1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	v1 "cthul.io/cthul/pkg/api/wave/v1"
	errors "errors"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// DomainServiceName is the fully-qualified name of the DomainService service.
	DomainServiceName = "wave.v1.DomainService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// DomainServiceGetDomainProcedure is the fully-qualified name of the DomainService's GetDomain RPC.
	DomainServiceGetDomainProcedure = "/wave.v1.DomainService/GetDomain"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	domainServiceServiceDescriptor         = v1.File_wave_v1_service_proto.Services().ByName("DomainService")
	domainServiceGetDomainMethodDescriptor = domainServiceServiceDescriptor.Methods().ByName("GetDomain")
)

// DomainServiceClient is a client for the wave.v1.DomainService service.
type DomainServiceClient interface {
	GetDomain(context.Context, *connect.Request[v1.GetDomainRequest]) (*connect.Response[v1.GetDomainResponse], error)
}

// NewDomainServiceClient constructs a client for the wave.v1.DomainService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewDomainServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) DomainServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &domainServiceClient{
		getDomain: connect.NewClient[v1.GetDomainRequest, v1.GetDomainResponse](
			httpClient,
			baseURL+DomainServiceGetDomainProcedure,
			connect.WithSchema(domainServiceGetDomainMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// domainServiceClient implements DomainServiceClient.
type domainServiceClient struct {
	getDomain *connect.Client[v1.GetDomainRequest, v1.GetDomainResponse]
}

// GetDomain calls wave.v1.DomainService.GetDomain.
func (c *domainServiceClient) GetDomain(ctx context.Context, req *connect.Request[v1.GetDomainRequest]) (*connect.Response[v1.GetDomainResponse], error) {
	return c.getDomain.CallUnary(ctx, req)
}

// DomainServiceHandler is an implementation of the wave.v1.DomainService service.
type DomainServiceHandler interface {
	GetDomain(context.Context, *connect.Request[v1.GetDomainRequest]) (*connect.Response[v1.GetDomainResponse], error)
}

// NewDomainServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewDomainServiceHandler(svc DomainServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	domainServiceGetDomainHandler := connect.NewUnaryHandler(
		DomainServiceGetDomainProcedure,
		svc.GetDomain,
		connect.WithSchema(domainServiceGetDomainMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/wave.v1.DomainService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case DomainServiceGetDomainProcedure:
			domainServiceGetDomainHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedDomainServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedDomainServiceHandler struct{}

func (UnimplementedDomainServiceHandler) GetDomain(context.Context, *connect.Request[v1.GetDomainRequest]) (*connect.Response[v1.GetDomainResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wave.v1.DomainService.GetDomain is not implemented"))
}
