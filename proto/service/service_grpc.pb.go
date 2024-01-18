// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: service/service.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	tag "product-service/proto/tag"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Service_GetTags_FullMethodName    = "/service.Service/GetTags"
	Service_GetTagById_FullMethodName = "/service.Service/GetTagById"
	Service_SaveTag_FullMethodName    = "/service.Service/SaveTag"
	Service_UpdateTag_FullMethodName  = "/service.Service/UpdateTag"
	Service_DeleteTag_FullMethodName  = "/service.Service/DeleteTag"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	// obtains tags by name
	GetTags(ctx context.Context, in *tag.GetTagsQuery, opts ...grpc.CallOption) (*tag.GetTagsResponse, error)
	// obtains tag by id
	GetTagById(ctx context.Context, in *tag.TagId, opts ...grpc.CallOption) (*tag.Tag, error)
	// save tag
	SaveTag(ctx context.Context, in *tag.SaveTagRequest, opts ...grpc.CallOption) (*tag.Tag, error)
	// update tag
	UpdateTag(ctx context.Context, in *tag.UpdateTagRequest, opts ...grpc.CallOption) (*tag.Tag, error)
	// deletes a tag
	DeleteTag(ctx context.Context, in *tag.TagId, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) GetTags(ctx context.Context, in *tag.GetTagsQuery, opts ...grpc.CallOption) (*tag.GetTagsResponse, error) {
	out := new(tag.GetTagsResponse)
	err := c.cc.Invoke(ctx, Service_GetTags_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetTagById(ctx context.Context, in *tag.TagId, opts ...grpc.CallOption) (*tag.Tag, error) {
	out := new(tag.Tag)
	err := c.cc.Invoke(ctx, Service_GetTagById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) SaveTag(ctx context.Context, in *tag.SaveTagRequest, opts ...grpc.CallOption) (*tag.Tag, error) {
	out := new(tag.Tag)
	err := c.cc.Invoke(ctx, Service_SaveTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateTag(ctx context.Context, in *tag.UpdateTagRequest, opts ...grpc.CallOption) (*tag.Tag, error) {
	out := new(tag.Tag)
	err := c.cc.Invoke(ctx, Service_UpdateTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteTag(ctx context.Context, in *tag.TagId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Service_DeleteTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	// obtains tags by name
	GetTags(context.Context, *tag.GetTagsQuery) (*tag.GetTagsResponse, error)
	// obtains tag by id
	GetTagById(context.Context, *tag.TagId) (*tag.Tag, error)
	// save tag
	SaveTag(context.Context, *tag.SaveTagRequest) (*tag.Tag, error)
	// update tag
	UpdateTag(context.Context, *tag.UpdateTagRequest) (*tag.Tag, error)
	// deletes a tag
	DeleteTag(context.Context, *tag.TagId) (*emptypb.Empty, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) GetTags(context.Context, *tag.GetTagsQuery) (*tag.GetTagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTags not implemented")
}
func (UnimplementedServiceServer) GetTagById(context.Context, *tag.TagId) (*tag.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTagById not implemented")
}
func (UnimplementedServiceServer) SaveTag(context.Context, *tag.SaveTagRequest) (*tag.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTag not implemented")
}
func (UnimplementedServiceServer) UpdateTag(context.Context, *tag.UpdateTagRequest) (*tag.Tag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTag not implemented")
}
func (UnimplementedServiceServer) DeleteTag(context.Context, *tag.TagId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTag not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_GetTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tag.GetTagsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetTags(ctx, req.(*tag.GetTagsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetTagById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tag.TagId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetTagById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetTagById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetTagById(ctx, req.(*tag.TagId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_SaveTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tag.SaveTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).SaveTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_SaveTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).SaveTag(ctx, req.(*tag.SaveTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tag.UpdateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_UpdateTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateTag(ctx, req.(*tag.UpdateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tag.TagId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_DeleteTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteTag(ctx, req.(*tag.TagId))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTags",
			Handler:    _Service_GetTags_Handler,
		},
		{
			MethodName: "GetTagById",
			Handler:    _Service_GetTagById_Handler,
		},
		{
			MethodName: "SaveTag",
			Handler:    _Service_SaveTag_Handler,
		},
		{
			MethodName: "UpdateTag",
			Handler:    _Service_UpdateTag_Handler,
		},
		{
			MethodName: "DeleteTag",
			Handler:    _Service_DeleteTag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/service.proto",
}