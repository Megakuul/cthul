// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: wave/v1/video/service.proto

package videoconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	video "cthul.io/cthul/pkg/api/wave/v1/video"
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
	// VideoServiceName is the fully-qualified name of the VideoService service.
	VideoServiceName = "wave.v1.video.VideoService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// VideoServiceGetProcedure is the fully-qualified name of the VideoService's Get RPC.
	VideoServiceGetProcedure = "/wave.v1.video.VideoService/Get"
	// VideoServiceListProcedure is the fully-qualified name of the VideoService's List RPC.
	VideoServiceListProcedure = "/wave.v1.video.VideoService/List"
	// VideoServiceCreateProcedure is the fully-qualified name of the VideoService's Create RPC.
	VideoServiceCreateProcedure = "/wave.v1.video.VideoService/Create"
	// VideoServiceUpdateProcedure is the fully-qualified name of the VideoService's Update RPC.
	VideoServiceUpdateProcedure = "/wave.v1.video.VideoService/Update"
	// VideoServiceDeleteProcedure is the fully-qualified name of the VideoService's Delete RPC.
	VideoServiceDeleteProcedure = "/wave.v1.video.VideoService/Delete"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	videoServiceServiceDescriptor      = video.File_wave_v1_video_service_proto.Services().ByName("VideoService")
	videoServiceGetMethodDescriptor    = videoServiceServiceDescriptor.Methods().ByName("Get")
	videoServiceListMethodDescriptor   = videoServiceServiceDescriptor.Methods().ByName("List")
	videoServiceCreateMethodDescriptor = videoServiceServiceDescriptor.Methods().ByName("Create")
	videoServiceUpdateMethodDescriptor = videoServiceServiceDescriptor.Methods().ByName("Update")
	videoServiceDeleteMethodDescriptor = videoServiceServiceDescriptor.Methods().ByName("Delete")
)

// VideoServiceClient is a client for the wave.v1.video.VideoService service.
type VideoServiceClient interface {
	Get(context.Context, *connect.Request[video.GetRequest]) (*connect.Response[video.GetResponse], error)
	List(context.Context, *connect.Request[video.ListRequest]) (*connect.Response[video.ListResponse], error)
	Create(context.Context, *connect.Request[video.CreateRequest]) (*connect.Response[video.CreateResponse], error)
	Update(context.Context, *connect.Request[video.UpdateRequest]) (*connect.Response[video.UpdateResponse], error)
	Delete(context.Context, *connect.Request[video.DeleteRequest]) (*connect.Response[video.DeleteResponse], error)
}

// NewVideoServiceClient constructs a client for the wave.v1.video.VideoService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewVideoServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) VideoServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &videoServiceClient{
		get: connect.NewClient[video.GetRequest, video.GetResponse](
			httpClient,
			baseURL+VideoServiceGetProcedure,
			connect.WithSchema(videoServiceGetMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		list: connect.NewClient[video.ListRequest, video.ListResponse](
			httpClient,
			baseURL+VideoServiceListProcedure,
			connect.WithSchema(videoServiceListMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		create: connect.NewClient[video.CreateRequest, video.CreateResponse](
			httpClient,
			baseURL+VideoServiceCreateProcedure,
			connect.WithSchema(videoServiceCreateMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		update: connect.NewClient[video.UpdateRequest, video.UpdateResponse](
			httpClient,
			baseURL+VideoServiceUpdateProcedure,
			connect.WithSchema(videoServiceUpdateMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		delete: connect.NewClient[video.DeleteRequest, video.DeleteResponse](
			httpClient,
			baseURL+VideoServiceDeleteProcedure,
			connect.WithSchema(videoServiceDeleteMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// videoServiceClient implements VideoServiceClient.
type videoServiceClient struct {
	get    *connect.Client[video.GetRequest, video.GetResponse]
	list   *connect.Client[video.ListRequest, video.ListResponse]
	create *connect.Client[video.CreateRequest, video.CreateResponse]
	update *connect.Client[video.UpdateRequest, video.UpdateResponse]
	delete *connect.Client[video.DeleteRequest, video.DeleteResponse]
}

// Get calls wave.v1.video.VideoService.Get.
func (c *videoServiceClient) Get(ctx context.Context, req *connect.Request[video.GetRequest]) (*connect.Response[video.GetResponse], error) {
	return c.get.CallUnary(ctx, req)
}

// List calls wave.v1.video.VideoService.List.
func (c *videoServiceClient) List(ctx context.Context, req *connect.Request[video.ListRequest]) (*connect.Response[video.ListResponse], error) {
	return c.list.CallUnary(ctx, req)
}

// Create calls wave.v1.video.VideoService.Create.
func (c *videoServiceClient) Create(ctx context.Context, req *connect.Request[video.CreateRequest]) (*connect.Response[video.CreateResponse], error) {
	return c.create.CallUnary(ctx, req)
}

// Update calls wave.v1.video.VideoService.Update.
func (c *videoServiceClient) Update(ctx context.Context, req *connect.Request[video.UpdateRequest]) (*connect.Response[video.UpdateResponse], error) {
	return c.update.CallUnary(ctx, req)
}

// Delete calls wave.v1.video.VideoService.Delete.
func (c *videoServiceClient) Delete(ctx context.Context, req *connect.Request[video.DeleteRequest]) (*connect.Response[video.DeleteResponse], error) {
	return c.delete.CallUnary(ctx, req)
}

// VideoServiceHandler is an implementation of the wave.v1.video.VideoService service.
type VideoServiceHandler interface {
	Get(context.Context, *connect.Request[video.GetRequest]) (*connect.Response[video.GetResponse], error)
	List(context.Context, *connect.Request[video.ListRequest]) (*connect.Response[video.ListResponse], error)
	Create(context.Context, *connect.Request[video.CreateRequest]) (*connect.Response[video.CreateResponse], error)
	Update(context.Context, *connect.Request[video.UpdateRequest]) (*connect.Response[video.UpdateResponse], error)
	Delete(context.Context, *connect.Request[video.DeleteRequest]) (*connect.Response[video.DeleteResponse], error)
}

// NewVideoServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewVideoServiceHandler(svc VideoServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	videoServiceGetHandler := connect.NewUnaryHandler(
		VideoServiceGetProcedure,
		svc.Get,
		connect.WithSchema(videoServiceGetMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	videoServiceListHandler := connect.NewUnaryHandler(
		VideoServiceListProcedure,
		svc.List,
		connect.WithSchema(videoServiceListMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	videoServiceCreateHandler := connect.NewUnaryHandler(
		VideoServiceCreateProcedure,
		svc.Create,
		connect.WithSchema(videoServiceCreateMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	videoServiceUpdateHandler := connect.NewUnaryHandler(
		VideoServiceUpdateProcedure,
		svc.Update,
		connect.WithSchema(videoServiceUpdateMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	videoServiceDeleteHandler := connect.NewUnaryHandler(
		VideoServiceDeleteProcedure,
		svc.Delete,
		connect.WithSchema(videoServiceDeleteMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/wave.v1.video.VideoService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case VideoServiceGetProcedure:
			videoServiceGetHandler.ServeHTTP(w, r)
		case VideoServiceListProcedure:
			videoServiceListHandler.ServeHTTP(w, r)
		case VideoServiceCreateProcedure:
			videoServiceCreateHandler.ServeHTTP(w, r)
		case VideoServiceUpdateProcedure:
			videoServiceUpdateHandler.ServeHTTP(w, r)
		case VideoServiceDeleteProcedure:
			videoServiceDeleteHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedVideoServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedVideoServiceHandler struct{}

func (UnimplementedVideoServiceHandler) Get(context.Context, *connect.Request[video.GetRequest]) (*connect.Response[video.GetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wave.v1.video.VideoService.Get is not implemented"))
}

func (UnimplementedVideoServiceHandler) List(context.Context, *connect.Request[video.ListRequest]) (*connect.Response[video.ListResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wave.v1.video.VideoService.List is not implemented"))
}

func (UnimplementedVideoServiceHandler) Create(context.Context, *connect.Request[video.CreateRequest]) (*connect.Response[video.CreateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wave.v1.video.VideoService.Create is not implemented"))
}

func (UnimplementedVideoServiceHandler) Update(context.Context, *connect.Request[video.UpdateRequest]) (*connect.Response[video.UpdateResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wave.v1.video.VideoService.Update is not implemented"))
}

func (UnimplementedVideoServiceHandler) Delete(context.Context, *connect.Request[video.DeleteRequest]) (*connect.Response[video.DeleteResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("wave.v1.video.VideoService.Delete is not implemented"))
}
