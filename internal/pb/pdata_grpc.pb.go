// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: pdata.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PrivateDataClient is the client API for PrivateData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrivateDataClient interface {
	AddPdata(ctx context.Context, in *AddPdataRequest, opts ...grpc.CallOption) (*AddPdataResponse, error)
	GetPdata(ctx context.Context, in *GetPdataRequest, opts ...grpc.CallOption) (*GetPdataResponse, error)
	UpdatePdata(ctx context.Context, in *UpdatePdataRequest, opts ...grpc.CallOption) (*UpdatePdataResponse, error)
}

type privateDataClient struct {
	cc grpc.ClientConnInterface
}

func NewPrivateDataClient(cc grpc.ClientConnInterface) PrivateDataClient {
	return &privateDataClient{cc}
}

func (c *privateDataClient) AddPdata(ctx context.Context, in *AddPdataRequest, opts ...grpc.CallOption) (*AddPdataResponse, error) {
	out := new(AddPdataResponse)
	err := c.cc.Invoke(ctx, "/PrivateData/AddPdata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateDataClient) GetPdata(ctx context.Context, in *GetPdataRequest, opts ...grpc.CallOption) (*GetPdataResponse, error) {
	out := new(GetPdataResponse)
	err := c.cc.Invoke(ctx, "/PrivateData/GetPdata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateDataClient) UpdatePdata(ctx context.Context, in *UpdatePdataRequest, opts ...grpc.CallOption) (*UpdatePdataResponse, error) {
	out := new(UpdatePdataResponse)
	err := c.cc.Invoke(ctx, "/PrivateData/UpdatePdata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrivateDataServer is the server API for PrivateData service.
// All implementations must embed UnimplementedPrivateDataServer
// for forward compatibility
type PrivateDataServer interface {
	AddPdata(context.Context, *AddPdataRequest) (*AddPdataResponse, error)
	GetPdata(context.Context, *GetPdataRequest) (*GetPdataResponse, error)
	UpdatePdata(context.Context, *UpdatePdataRequest) (*UpdatePdataResponse, error)
	mustEmbedUnimplementedPrivateDataServer()
}

// UnimplementedPrivateDataServer must be embedded to have forward compatible implementations.
type UnimplementedPrivateDataServer struct {
}

func (UnimplementedPrivateDataServer) AddPdata(context.Context, *AddPdataRequest) (*AddPdataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPdata not implemented")
}
func (UnimplementedPrivateDataServer) GetPdata(context.Context, *GetPdataRequest) (*GetPdataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPdata not implemented")
}
func (UnimplementedPrivateDataServer) UpdatePdata(context.Context, *UpdatePdataRequest) (*UpdatePdataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePdata not implemented")
}
func (UnimplementedPrivateDataServer) mustEmbedUnimplementedPrivateDataServer() {}

// UnsafePrivateDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrivateDataServer will
// result in compilation errors.
type UnsafePrivateDataServer interface {
	mustEmbedUnimplementedPrivateDataServer()
}

func RegisterPrivateDataServer(s grpc.ServiceRegistrar, srv PrivateDataServer) {
	s.RegisterService(&PrivateData_ServiceDesc, srv)
}

func _PrivateData_AddPdata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPdataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateDataServer).AddPdata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PrivateData/AddPdata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateDataServer).AddPdata(ctx, req.(*AddPdataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateData_GetPdata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPdataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateDataServer).GetPdata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PrivateData/GetPdata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateDataServer).GetPdata(ctx, req.(*GetPdataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateData_UpdatePdata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePdataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateDataServer).UpdatePdata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PrivateData/UpdatePdata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateDataServer).UpdatePdata(ctx, req.(*UpdatePdataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrivateData_ServiceDesc is the grpc.ServiceDesc for PrivateData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrivateData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PrivateData",
	HandlerType: (*PrivateDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddPdata",
			Handler:    _PrivateData_AddPdata_Handler,
		},
		{
			MethodName: "GetPdata",
			Handler:    _PrivateData_GetPdata_Handler,
		},
		{
			MethodName: "UpdatePdata",
			Handler:    _PrivateData_UpdatePdata_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pdata.proto",
}
