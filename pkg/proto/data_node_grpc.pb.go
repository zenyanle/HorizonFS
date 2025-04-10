// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.1
// source: pkg/proto/data_node.proto

package proto

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
	DataNodeService_DownloadData_FullMethodName = "/proto.DataNodeService/DownloadData"
	DataNodeService_UploadData_FullMethodName   = "/proto.DataNodeService/UploadData"
)

// DataNodeServiceClient is the client API for DataNodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataNodeServiceClient interface {
	DownloadData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DownloadDataRequest, DataChunk], error)
	UploadData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DataChunk, UploadDataResponse], error)
}

type dataNodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataNodeServiceClient(cc grpc.ClientConnInterface) DataNodeServiceClient {
	return &dataNodeServiceClient{cc}
}

func (c *dataNodeServiceClient) DownloadData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DownloadDataRequest, DataChunk], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DataNodeService_ServiceDesc.Streams[0], DataNodeService_DownloadData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadDataRequest, DataChunk]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataNodeService_DownloadDataClient = grpc.BidiStreamingClient[DownloadDataRequest, DataChunk]

func (c *dataNodeServiceClient) UploadData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DataChunk, UploadDataResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DataNodeService_ServiceDesc.Streams[1], DataNodeService_UploadData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DataChunk, UploadDataResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataNodeService_UploadDataClient = grpc.BidiStreamingClient[DataChunk, UploadDataResponse]

// DataNodeServiceServer is the server API for DataNodeService service.
// All implementations must embed UnimplementedDataNodeServiceServer
// for forward compatibility.
type DataNodeServiceServer interface {
	DownloadData(grpc.BidiStreamingServer[DownloadDataRequest, DataChunk]) error
	UploadData(grpc.BidiStreamingServer[DataChunk, UploadDataResponse]) error
	mustEmbedUnimplementedDataNodeServiceServer()
}

// UnimplementedDataNodeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataNodeServiceServer struct{}

func (UnimplementedDataNodeServiceServer) DownloadData(grpc.BidiStreamingServer[DownloadDataRequest, DataChunk]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadData not implemented")
}
func (UnimplementedDataNodeServiceServer) UploadData(grpc.BidiStreamingServer[DataChunk, UploadDataResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadData not implemented")
}
func (UnimplementedDataNodeServiceServer) mustEmbedUnimplementedDataNodeServiceServer() {}
func (UnimplementedDataNodeServiceServer) testEmbeddedByValue()                         {}

// UnsafeDataNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataNodeServiceServer will
// result in compilation errors.
type UnsafeDataNodeServiceServer interface {
	mustEmbedUnimplementedDataNodeServiceServer()
}

func RegisterDataNodeServiceServer(s grpc.ServiceRegistrar, srv DataNodeServiceServer) {
	// If the following call pancis, it indicates UnimplementedDataNodeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataNodeService_ServiceDesc, srv)
}

func _DataNodeService_DownloadData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataNodeServiceServer).DownloadData(&grpc.GenericServerStream[DownloadDataRequest, DataChunk]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataNodeService_DownloadDataServer = grpc.BidiStreamingServer[DownloadDataRequest, DataChunk]

func _DataNodeService_UploadData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataNodeServiceServer).UploadData(&grpc.GenericServerStream[DataChunk, UploadDataResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataNodeService_UploadDataServer = grpc.BidiStreamingServer[DataChunk, UploadDataResponse]

// DataNodeService_ServiceDesc is the grpc.ServiceDesc for DataNodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataNodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DataNodeService",
	HandlerType: (*DataNodeServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadData",
			Handler:       _DataNodeService_DownloadData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadData",
			Handler:       _DataNodeService_UploadData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/proto/data_node.proto",
}
