// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: internal/api/api.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GoMetrics_BatchUpdate_FullMethodName = "/api.GoMetrics/BatchUpdate"
)

// GoMetricsClient is the client API for GoMetrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoMetricsClient interface {
	BatchUpdate(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type goMetricsClient struct {
	cc grpc.ClientConnInterface
}

func NewGoMetricsClient(cc grpc.ClientConnInterface) GoMetricsClient {
	return &goMetricsClient{cc}
}

func (c *goMetricsClient) BatchUpdate(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, GoMetrics_BatchUpdate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoMetricsServer is the server API for GoMetrics service.
// All implementations must embed UnimplementedGoMetricsServer
// for forward compatibility.
type GoMetricsServer interface {
	BatchUpdate(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedGoMetricsServer()
}

// UnimplementedGoMetricsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGoMetricsServer struct{}

func (UnimplementedGoMetricsServer) BatchUpdate(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpdate not implemented")
}
func (UnimplementedGoMetricsServer) mustEmbedUnimplementedGoMetricsServer() {}
func (UnimplementedGoMetricsServer) testEmbeddedByValue()                   {}

// UnsafeGoMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoMetricsServer will
// result in compilation errors.
type UnsafeGoMetricsServer interface {
	mustEmbedUnimplementedGoMetricsServer()
}

func RegisterGoMetricsServer(s grpc.ServiceRegistrar, srv GoMetricsServer) {
	// If the following call pancis, it indicates UnimplementedGoMetricsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GoMetrics_ServiceDesc, srv)
}

func _GoMetrics_BatchUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoMetricsServer).BatchUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoMetrics_BatchUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoMetricsServer).BatchUpdate(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// GoMetrics_ServiceDesc is the grpc.ServiceDesc for GoMetrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoMetrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.GoMetrics",
	HandlerType: (*GoMetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BatchUpdate",
			Handler:    _GoMetrics_BatchUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/api/api.proto",
}
