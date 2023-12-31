// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: api/proto/secret.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Secret_Create_FullMethodName     = "/gophkeeper.Secret/Create"
	Secret_GetSecrets_FullMethodName = "/gophkeeper.Secret/GetSecrets"
	Secret_Update_FullMethodName     = "/gophkeeper.Secret/Update"
	Secret_Delete_FullMethodName     = "/gophkeeper.Secret/Delete"
)

// SecretClient is the client API for Secret service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetSecrets(ctx context.Context, in *GetSecretsRequest, opts ...grpc.CallOption) (*GetSecretsResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type secretClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretClient(cc grpc.ClientConnInterface) SecretClient {
	return &secretClient{cc}
}

func (c *secretClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Secret_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretClient) GetSecrets(ctx context.Context, in *GetSecretsRequest, opts ...grpc.CallOption) (*GetSecretsResponse, error) {
	out := new(GetSecretsResponse)
	err := c.cc.Invoke(ctx, Secret_GetSecrets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Secret_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Secret_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretServer is the server API for Secret service.
// All implementations must embed UnimplementedSecretServer
// for forward compatibility
type SecretServer interface {
	Create(context.Context, *CreateRequest) (*emptypb.Empty, error)
	GetSecrets(context.Context, *GetSecretsRequest) (*GetSecretsResponse, error)
	Update(context.Context, *UpdateRequest) (*emptypb.Empty, error)
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedSecretServer()
}

// UnimplementedSecretServer must be embedded to have forward compatible implementations.
type UnimplementedSecretServer struct {
}

func (UnimplementedSecretServer) Create(context.Context, *CreateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedSecretServer) GetSecrets(context.Context, *GetSecretsRequest) (*GetSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecrets not implemented")
}
func (UnimplementedSecretServer) Update(context.Context, *UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedSecretServer) Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedSecretServer) mustEmbedUnimplementedSecretServer() {}

// UnsafeSecretServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretServer will
// result in compilation errors.
type UnsafeSecretServer interface {
	mustEmbedUnimplementedSecretServer()
}

func RegisterSecretServer(s grpc.ServiceRegistrar, srv SecretServer) {
	s.RegisterService(&Secret_ServiceDesc, srv)
}

func _Secret_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secret_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secret_GetSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServer).GetSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secret_GetSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServer).GetSecrets(ctx, req.(*GetSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secret_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secret_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secret_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secret_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Secret_ServiceDesc is the grpc.ServiceDesc for Secret service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Secret_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gophkeeper.Secret",
	HandlerType: (*SecretServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Secret_Create_Handler,
		},
		{
			MethodName: "GetSecrets",
			Handler:    _Secret_GetSecrets_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Secret_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Secret_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/secret.proto",
}
