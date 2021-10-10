// Code generated by protoc-gen-go.
// source: golang.singingcat.net/apis/usercommand/usercommand.proto
// DO NOT EDIT!

/*
Package usercommand is a generated protocol buffer package.

It is generated from these files:
	golang.singingcat.net/apis/usercommand/usercommand.proto

It has these top-level messages:
	UserCommand
*/
package usercommand

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import singingcat "golang.singingcat.net/apis/singingcat"
import common "golang.conradwood.net/apis/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// a usercommand was received
type UserCommand struct {
	Module *singingcat.ModuleRef `protobuf:"bytes,1,opt,name=Module" json:"Module,omitempty"`
	Args   [][]byte              `protobuf:"bytes,2,rep,name=Args,proto3" json:"Args,omitempty"`
}

func (m *UserCommand) Reset()                    { *m = UserCommand{} }
func (m *UserCommand) String() string            { return proto.CompactTextString(m) }
func (*UserCommand) ProtoMessage()               {}
func (*UserCommand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UserCommand) GetModule() *singingcat.ModuleRef {
	if m != nil {
		return m.Module
	}
	return nil
}

func (m *UserCommand) GetArgs() [][]byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func init() {
	proto.RegisterType((*UserCommand)(nil), "usercommand.UserCommand")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UserCommandService service

type UserCommandServiceClient interface {
	// received a new one
	CommandReceived(ctx context.Context, in *UserCommand, opts ...grpc.CallOption) (*common.Void, error)
}

type userCommandServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserCommandServiceClient(cc *grpc.ClientConn) UserCommandServiceClient {
	return &userCommandServiceClient{cc}
}

func (c *userCommandServiceClient) CommandReceived(ctx context.Context, in *UserCommand, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/usercommand.UserCommandService/CommandReceived", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserCommandService service

type UserCommandServiceServer interface {
	// received a new one
	CommandReceived(context.Context, *UserCommand) (*common.Void, error)
}

func RegisterUserCommandServiceServer(s *grpc.Server, srv UserCommandServiceServer) {
	s.RegisterService(&_UserCommandService_serviceDesc, srv)
}

func _UserCommandService_CommandReceived_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserCommandServiceServer).CommandReceived(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usercommand.UserCommandService/CommandReceived",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserCommandServiceServer).CommandReceived(ctx, req.(*UserCommand))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserCommandService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "usercommand.UserCommandService",
	HandlerType: (*UserCommandServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommandReceived",
			Handler:    _UserCommandService_CommandReceived_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.singingcat.net/apis/usercommand/usercommand.proto",
}

func init() {
	proto.RegisterFile("golang.singingcat.net/apis/usercommand/usercommand.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x90, 0x31, 0x4b, 0x04, 0x31,
	0x10, 0x85, 0x39, 0x95, 0x2b, 0xb2, 0x07, 0x42, 0x40, 0x58, 0xb6, 0x3a, 0x2c, 0x64, 0x1b, 0x73,
	0x70, 0x82, 0x68, 0xa9, 0xd6, 0xa2, 0x44, 0xb4, 0x8f, 0xc9, 0x18, 0x02, 0x77, 0x33, 0x47, 0x92,
	0x5d, 0xff, 0xbe, 0xec, 0x26, 0xb0, 0x53, 0x59, 0x65, 0x78, 0x99, 0xf7, 0xe5, 0xbd, 0x88, 0x07,
	0x4f, 0x07, 0x83, 0x5e, 0xa5, 0x80, 0x3e, 0xa0, 0xb7, 0x26, 0x2b, 0x84, 0xbc, 0x33, 0xa7, 0x90,
	0x76, 0x43, 0x82, 0x68, 0xe9, 0x78, 0x34, 0xe8, 0xf8, 0xac, 0x4e, 0x91, 0x32, 0xc9, 0x86, 0x49,
	0xdd, 0xfd, 0x3f, 0x98, 0x45, 0x63, 0x63, 0x81, 0x74, 0xaa, 0xfa, 0x2c, 0x61, 0x34, 0xee, 0x97,
	0xc8, 0x2d, 0xbe, 0x89, 0x4d, 0x58, 0x8f, 0xb2, 0x7f, 0xfd, 0x2e, 0x9a, 0xcf, 0x04, 0xf1, 0xa5,
	0x3c, 0x2b, 0x6f, 0xc5, 0xfa, 0x95, 0xdc, 0x70, 0x80, 0x76, 0xb5, 0x5d, 0xf5, 0xcd, 0xfe, 0x8a,
	0x07, 0x28, 0x37, 0x1a, 0x7e, 0x74, 0x5d, 0x92, 0x52, 0x5c, 0x3c, 0x45, 0x9f, 0xda, 0xb3, 0xed,
	0x79, 0xbf, 0xd1, 0xf3, 0xbc, 0x7f, 0x13, 0x92, 0x11, 0x3f, 0x20, 0x8e, 0xc1, 0x82, 0x7c, 0x14,
	0x97, 0x55, 0xd1, 0x60, 0x21, 0x8c, 0xe0, 0x64, 0xab, 0xf8, 0x1f, 0x30, 0x4f, 0xb7, 0x51, 0x35,
	0xe3, 0x17, 0x05, 0xf7, 0xdc, 0x8b, 0x1b, 0x84, 0xcc, 0x83, 0xd4, 0x8e, 0x53, 0x29, 0xce, 0xf8,
	0x5e, 0xcf, 0x9d, 0xee, 0xfe, 0x02, 0x00, 0x00, 0xff, 0xff, 0xd2, 0x25, 0x9f, 0x95, 0x84, 0x01,
	0x00, 0x00,
}
