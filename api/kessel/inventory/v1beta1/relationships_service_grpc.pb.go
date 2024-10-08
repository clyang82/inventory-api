// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: kessel/inventory/v1beta1/relationships_service.proto

package v1beta1

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
	KesselPolicyRelationshipService_CreatePolicyRelationship_FullMethodName          = "/kessel.inventory.v1beta1.KesselPolicyRelationshipService/CreatePolicyRelationship"
	KesselPolicyRelationshipService_UpdateResourceRelationshipByUrnHs_FullMethodName = "/kessel.inventory.v1beta1.KesselPolicyRelationshipService/UpdateResourceRelationshipByUrnHs"
	KesselPolicyRelationshipService_DeleteResourceRelationshipByUrn_FullMethodName   = "/kessel.inventory.v1beta1.KesselPolicyRelationshipService/DeleteResourceRelationshipByUrn"
)

// KesselPolicyRelationshipServiceClient is the client API for KesselPolicyRelationshipService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KesselPolicyRelationshipServiceClient interface {
	CreatePolicyRelationship(ctx context.Context, in *CreatePolicyRelationshipRequest, opts ...grpc.CallOption) (*CreatePolicyRelationshipResponse, error)
	UpdateResourceRelationshipByUrnHs(ctx context.Context, in *UpdateResourceRelationshipByUrnHsRequest, opts ...grpc.CallOption) (*UpdateResourceRelationshipByUrnHsResponse, error)
	DeleteResourceRelationshipByUrn(ctx context.Context, in *DeleteResourceRelationshipByUrnRequest, opts ...grpc.CallOption) (*DeleteResourceRelationshipByUrnResponse, error)
}

type kesselPolicyRelationshipServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKesselPolicyRelationshipServiceClient(cc grpc.ClientConnInterface) KesselPolicyRelationshipServiceClient {
	return &kesselPolicyRelationshipServiceClient{cc}
}

func (c *kesselPolicyRelationshipServiceClient) CreatePolicyRelationship(ctx context.Context, in *CreatePolicyRelationshipRequest, opts ...grpc.CallOption) (*CreatePolicyRelationshipResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePolicyRelationshipResponse)
	err := c.cc.Invoke(ctx, KesselPolicyRelationshipService_CreatePolicyRelationship_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kesselPolicyRelationshipServiceClient) UpdateResourceRelationshipByUrnHs(ctx context.Context, in *UpdateResourceRelationshipByUrnHsRequest, opts ...grpc.CallOption) (*UpdateResourceRelationshipByUrnHsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateResourceRelationshipByUrnHsResponse)
	err := c.cc.Invoke(ctx, KesselPolicyRelationshipService_UpdateResourceRelationshipByUrnHs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kesselPolicyRelationshipServiceClient) DeleteResourceRelationshipByUrn(ctx context.Context, in *DeleteResourceRelationshipByUrnRequest, opts ...grpc.CallOption) (*DeleteResourceRelationshipByUrnResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResourceRelationshipByUrnResponse)
	err := c.cc.Invoke(ctx, KesselPolicyRelationshipService_DeleteResourceRelationshipByUrn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KesselPolicyRelationshipServiceServer is the server API for KesselPolicyRelationshipService service.
// All implementations must embed UnimplementedKesselPolicyRelationshipServiceServer
// for forward compatibility.
type KesselPolicyRelationshipServiceServer interface {
	CreatePolicyRelationship(context.Context, *CreatePolicyRelationshipRequest) (*CreatePolicyRelationshipResponse, error)
	UpdateResourceRelationshipByUrnHs(context.Context, *UpdateResourceRelationshipByUrnHsRequest) (*UpdateResourceRelationshipByUrnHsResponse, error)
	DeleteResourceRelationshipByUrn(context.Context, *DeleteResourceRelationshipByUrnRequest) (*DeleteResourceRelationshipByUrnResponse, error)
	mustEmbedUnimplementedKesselPolicyRelationshipServiceServer()
}

// UnimplementedKesselPolicyRelationshipServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKesselPolicyRelationshipServiceServer struct{}

func (UnimplementedKesselPolicyRelationshipServiceServer) CreatePolicyRelationship(context.Context, *CreatePolicyRelationshipRequest) (*CreatePolicyRelationshipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePolicyRelationship not implemented")
}
func (UnimplementedKesselPolicyRelationshipServiceServer) UpdateResourceRelationshipByUrnHs(context.Context, *UpdateResourceRelationshipByUrnHsRequest) (*UpdateResourceRelationshipByUrnHsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResourceRelationshipByUrnHs not implemented")
}
func (UnimplementedKesselPolicyRelationshipServiceServer) DeleteResourceRelationshipByUrn(context.Context, *DeleteResourceRelationshipByUrnRequest) (*DeleteResourceRelationshipByUrnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteResourceRelationshipByUrn not implemented")
}
func (UnimplementedKesselPolicyRelationshipServiceServer) mustEmbedUnimplementedKesselPolicyRelationshipServiceServer() {
}
func (UnimplementedKesselPolicyRelationshipServiceServer) testEmbeddedByValue() {}

// UnsafeKesselPolicyRelationshipServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KesselPolicyRelationshipServiceServer will
// result in compilation errors.
type UnsafeKesselPolicyRelationshipServiceServer interface {
	mustEmbedUnimplementedKesselPolicyRelationshipServiceServer()
}

func RegisterKesselPolicyRelationshipServiceServer(s grpc.ServiceRegistrar, srv KesselPolicyRelationshipServiceServer) {
	// If the following call pancis, it indicates UnimplementedKesselPolicyRelationshipServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KesselPolicyRelationshipService_ServiceDesc, srv)
}

func _KesselPolicyRelationshipService_CreatePolicyRelationship_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePolicyRelationshipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KesselPolicyRelationshipServiceServer).CreatePolicyRelationship(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KesselPolicyRelationshipService_CreatePolicyRelationship_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KesselPolicyRelationshipServiceServer).CreatePolicyRelationship(ctx, req.(*CreatePolicyRelationshipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KesselPolicyRelationshipService_UpdateResourceRelationshipByUrnHs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceRelationshipByUrnHsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KesselPolicyRelationshipServiceServer).UpdateResourceRelationshipByUrnHs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KesselPolicyRelationshipService_UpdateResourceRelationshipByUrnHs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KesselPolicyRelationshipServiceServer).UpdateResourceRelationshipByUrnHs(ctx, req.(*UpdateResourceRelationshipByUrnHsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KesselPolicyRelationshipService_DeleteResourceRelationshipByUrn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteResourceRelationshipByUrnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KesselPolicyRelationshipServiceServer).DeleteResourceRelationshipByUrn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KesselPolicyRelationshipService_DeleteResourceRelationshipByUrn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KesselPolicyRelationshipServiceServer).DeleteResourceRelationshipByUrn(ctx, req.(*DeleteResourceRelationshipByUrnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KesselPolicyRelationshipService_ServiceDesc is the grpc.ServiceDesc for KesselPolicyRelationshipService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KesselPolicyRelationshipService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kessel.inventory.v1beta1.KesselPolicyRelationshipService",
	HandlerType: (*KesselPolicyRelationshipServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePolicyRelationship",
			Handler:    _KesselPolicyRelationshipService_CreatePolicyRelationship_Handler,
		},
		{
			MethodName: "UpdateResourceRelationshipByUrnHs",
			Handler:    _KesselPolicyRelationshipService_UpdateResourceRelationshipByUrnHs_Handler,
		},
		{
			MethodName: "DeleteResourceRelationshipByUrn",
			Handler:    _KesselPolicyRelationshipService_DeleteResourceRelationshipByUrn_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kessel/inventory/v1beta1/relationships_service.proto",
}
