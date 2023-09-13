// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/v1/customer/customer.proto

package __

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

// CustomerServiceClient is the client API for CustomerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomerServiceClient interface {
	CreateCustomer(ctx context.Context, in *CustomerRequest, opts ...grpc.CallOption) (*CustomerResponse, error)
	UpdateVerifiedCustomer(ctx context.Context, in *UpdateVerifiedCustomerRequest, opts ...grpc.CallOption) (*UpdateVerifiedCustomerResponse, error)
	GetCustomerIDByEmail(ctx context.Context, in *GetCustomerIDByEmailRequest, opts ...grpc.CallOption) (*GetCustomerIDByEmailResponse, error)
	GetCustomerDetailsByEmail(ctx context.Context, in *GetCustomerDetailsByEmailRequest, opts ...grpc.CallOption) (*GetCustomerDetailsByEmailResponse, error)
}

type customerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomerServiceClient(cc grpc.ClientConnInterface) CustomerServiceClient {
	return &customerServiceClient{cc}
}

func (c *customerServiceClient) CreateCustomer(ctx context.Context, in *CustomerRequest, opts ...grpc.CallOption) (*CustomerResponse, error) {
	out := new(CustomerResponse)
	err := c.cc.Invoke(ctx, "/proto.v1.customer.CustomerService/CreateCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerServiceClient) UpdateVerifiedCustomer(ctx context.Context, in *UpdateVerifiedCustomerRequest, opts ...grpc.CallOption) (*UpdateVerifiedCustomerResponse, error) {
	out := new(UpdateVerifiedCustomerResponse)
	err := c.cc.Invoke(ctx, "/proto.v1.customer.CustomerService/UpdateVerifiedCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerServiceClient) GetCustomerIDByEmail(ctx context.Context, in *GetCustomerIDByEmailRequest, opts ...grpc.CallOption) (*GetCustomerIDByEmailResponse, error) {
	out := new(GetCustomerIDByEmailResponse)
	err := c.cc.Invoke(ctx, "/proto.v1.customer.CustomerService/GetCustomerIDByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerServiceClient) GetCustomerDetailsByEmail(ctx context.Context, in *GetCustomerDetailsByEmailRequest, opts ...grpc.CallOption) (*GetCustomerDetailsByEmailResponse, error) {
	out := new(GetCustomerDetailsByEmailResponse)
	err := c.cc.Invoke(ctx, "/proto.v1.customer.CustomerService/GetCustomerDetailsByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomerServiceServer is the server API for CustomerService service.
// All implementations must embed UnimplementedCustomerServiceServer
// for forward compatibility
type CustomerServiceServer interface {
	CreateCustomer(context.Context, *CustomerRequest) (*CustomerResponse, error)
	UpdateVerifiedCustomer(context.Context, *UpdateVerifiedCustomerRequest) (*UpdateVerifiedCustomerResponse, error)
	GetCustomerIDByEmail(context.Context, *GetCustomerIDByEmailRequest) (*GetCustomerIDByEmailResponse, error)
	GetCustomerDetailsByEmail(context.Context, *GetCustomerDetailsByEmailRequest) (*GetCustomerDetailsByEmailResponse, error)
	mustEmbedUnimplementedCustomerServiceServer()
}

// UnimplementedCustomerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCustomerServiceServer struct {
}

func (UnimplementedCustomerServiceServer) CreateCustomer(context.Context, *CustomerRequest) (*CustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCustomer not implemented")
}
func (UnimplementedCustomerServiceServer) UpdateVerifiedCustomer(context.Context, *UpdateVerifiedCustomerRequest) (*UpdateVerifiedCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVerifiedCustomer not implemented")
}
func (UnimplementedCustomerServiceServer) GetCustomerIDByEmail(context.Context, *GetCustomerIDByEmailRequest) (*GetCustomerIDByEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomerIDByEmail not implemented")
}
func (UnimplementedCustomerServiceServer) GetCustomerDetailsByEmail(context.Context, *GetCustomerDetailsByEmailRequest) (*GetCustomerDetailsByEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomerDetailsByEmail not implemented")
}
func (UnimplementedCustomerServiceServer) mustEmbedUnimplementedCustomerServiceServer() {}

// UnsafeCustomerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomerServiceServer will
// result in compilation errors.
type UnsafeCustomerServiceServer interface {
	mustEmbedUnimplementedCustomerServiceServer()
}

func RegisterCustomerServiceServer(s grpc.ServiceRegistrar, srv CustomerServiceServer) {
	s.RegisterService(&CustomerService_ServiceDesc, srv)
}

func _CustomerService_CreateCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServiceServer).CreateCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.customer.CustomerService/CreateCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServiceServer).CreateCustomer(ctx, req.(*CustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomerService_UpdateVerifiedCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVerifiedCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServiceServer).UpdateVerifiedCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.customer.CustomerService/UpdateVerifiedCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServiceServer).UpdateVerifiedCustomer(ctx, req.(*UpdateVerifiedCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomerService_GetCustomerIDByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomerIDByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServiceServer).GetCustomerIDByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.customer.CustomerService/GetCustomerIDByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServiceServer).GetCustomerIDByEmail(ctx, req.(*GetCustomerIDByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomerService_GetCustomerDetailsByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomerDetailsByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServiceServer).GetCustomerDetailsByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.v1.customer.CustomerService/GetCustomerDetailsByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServiceServer).GetCustomerDetailsByEmail(ctx, req.(*GetCustomerDetailsByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CustomerService_ServiceDesc is the grpc.ServiceDesc for CustomerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CustomerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.customer.CustomerService",
	HandlerType: (*CustomerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCustomer",
			Handler:    _CustomerService_CreateCustomer_Handler,
		},
		{
			MethodName: "UpdateVerifiedCustomer",
			Handler:    _CustomerService_UpdateVerifiedCustomer_Handler,
		},
		{
			MethodName: "GetCustomerIDByEmail",
			Handler:    _CustomerService_GetCustomerIDByEmail_Handler,
		},
		{
			MethodName: "GetCustomerDetailsByEmail",
			Handler:    _CustomerService_GetCustomerDetailsByEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/customer/customer.proto",
}
