// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/services/mobile_app_category_constant_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v1/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Request message for
// [MobileAppCategoryConstantService.GetMobileAppCategoryConstant][google.ads.googleads.v1.services.MobileAppCategoryConstantService.GetMobileAppCategoryConstant].
type GetMobileAppCategoryConstantRequest struct {
	// Resource name of the mobile app category constant to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMobileAppCategoryConstantRequest) Reset()         { *m = GetMobileAppCategoryConstantRequest{} }
func (m *GetMobileAppCategoryConstantRequest) String() string { return proto.CompactTextString(m) }
func (*GetMobileAppCategoryConstantRequest) ProtoMessage()    {}
func (*GetMobileAppCategoryConstantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_57065493a9a2f64e, []int{0}
}

func (m *GetMobileAppCategoryConstantRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMobileAppCategoryConstantRequest.Unmarshal(m, b)
}
func (m *GetMobileAppCategoryConstantRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMobileAppCategoryConstantRequest.Marshal(b, m, deterministic)
}
func (m *GetMobileAppCategoryConstantRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMobileAppCategoryConstantRequest.Merge(m, src)
}
func (m *GetMobileAppCategoryConstantRequest) XXX_Size() int {
	return xxx_messageInfo_GetMobileAppCategoryConstantRequest.Size(m)
}
func (m *GetMobileAppCategoryConstantRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMobileAppCategoryConstantRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetMobileAppCategoryConstantRequest proto.InternalMessageInfo

func (m *GetMobileAppCategoryConstantRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetMobileAppCategoryConstantRequest)(nil), "google.ads.googleads.v1.services.GetMobileAppCategoryConstantRequest")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/services/mobile_app_category_constant_service.proto", fileDescriptor_57065493a9a2f64e)
}

var fileDescriptor_57065493a9a2f64e = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x31, 0x4b, 0xfb, 0x40,
	0x18, 0xc6, 0x49, 0xfe, 0xf0, 0x07, 0x83, 0x2e, 0x99, 0xa4, 0x74, 0x08, 0x6d, 0x05, 0x71, 0xb8,
	0x98, 0xba, 0xc8, 0xa9, 0x43, 0x5a, 0xa5, 0xa2, 0x28, 0xa5, 0x42, 0x07, 0x09, 0x84, 0x6b, 0x72,
	0x84, 0x40, 0x73, 0x77, 0xe6, 0xbd, 0x16, 0x44, 0x5c, 0x3a, 0xbb, 0xf9, 0x0d, 0x1c, 0xfd, 0x28,
	0xae, 0xce, 0x6e, 0x4e, 0x7e, 0x0a, 0x49, 0x2f, 0x17, 0x70, 0x88, 0x71, 0x7b, 0xb8, 0x3c, 0xf9,
	0x3d, 0xf7, 0x3e, 0xf7, 0x5a, 0x97, 0x09, 0xe7, 0xc9, 0x9c, 0xba, 0x24, 0x06, 0x57, 0xc9, 0x42,
	0x2d, 0x3d, 0x17, 0x68, 0xbe, 0x4c, 0x23, 0x0a, 0x6e, 0xc6, 0x67, 0xe9, 0x9c, 0x86, 0x44, 0x88,
	0x30, 0x22, 0x92, 0x26, 0x3c, 0xbf, 0x0f, 0x23, 0xce, 0x40, 0x12, 0x26, 0xc3, 0xd2, 0x85, 0x44,
	0xce, 0x25, 0xb7, 0x1d, 0x45, 0x40, 0x24, 0x06, 0x54, 0xc1, 0xd0, 0xd2, 0x43, 0x1a, 0xd6, 0x3a,
	0xad, 0x8b, 0xcb, 0x29, 0xf0, 0x45, 0xde, 0x94, 0xa7, 0x72, 0x5a, 0x6d, 0x4d, 0x11, 0xa9, 0x4b,
	0x18, 0xe3, 0x92, 0xc8, 0x94, 0x33, 0x50, 0x5f, 0x3b, 0x17, 0x56, 0x77, 0x44, 0xe5, 0xd5, 0x1a,
	0xe3, 0x0b, 0x31, 0x2c, 0x21, 0xc3, 0x92, 0x31, 0xa1, 0x77, 0x0b, 0x0a, 0xd2, 0xee, 0x5a, 0x5b,
	0x3a, 0x34, 0x64, 0x24, 0xa3, 0xdb, 0x86, 0x63, 0xec, 0x6e, 0x4c, 0x36, 0xf5, 0xe1, 0x35, 0xc9,
	0x68, 0x7f, 0x65, 0x5a, 0x4e, 0x2d, 0xe9, 0x46, 0x4d, 0x65, 0x7f, 0x18, 0x56, 0xfb, 0xb7, 0x44,
	0xfb, 0x0c, 0x35, 0x15, 0x83, 0xfe, 0x70, 0xe3, 0xd6, 0x71, 0x2d, 0xa6, 0x6a, 0x0f, 0xd5, 0x42,
	0x3a, 0x87, 0xab, 0xf7, 0xcf, 0x67, 0xb3, 0x6f, 0xef, 0x17, 0x75, 0x3f, 0xfc, 0x18, 0xfd, 0x24,
	0xab, 0xfb, 0x0b, 0xdc, 0xbd, 0xc7, 0xc1, 0x93, 0x69, 0xf5, 0x22, 0x9e, 0x35, 0x0e, 0x31, 0xd8,
	0x69, 0xaa, 0x6a, 0x5c, 0x3c, 0xd0, 0xd8, 0xb8, 0x3d, 0x2f, 0x51, 0x09, 0x9f, 0x13, 0x96, 0x20,
	0x9e, 0x27, 0x6e, 0x42, 0xd9, 0xfa, 0xf9, 0xf4, 0x5a, 0x88, 0x14, 0xea, 0x97, 0xf2, 0x48, 0x8b,
	0x17, 0xf3, 0xdf, 0xc8, 0xf7, 0x5f, 0x4d, 0x67, 0xa4, 0x80, 0x7e, 0x0c, 0x48, 0xc9, 0x42, 0x4d,
	0x3d, 0x54, 0x06, 0xc3, 0x9b, 0xb6, 0x04, 0x7e, 0x0c, 0x41, 0x65, 0x09, 0xa6, 0x5e, 0xa0, 0x2d,
	0x5f, 0x66, 0x4f, 0x9d, 0x63, 0xec, 0xc7, 0x80, 0x71, 0x65, 0xc2, 0x78, 0xea, 0x61, 0xac, 0x6d,
	0xb3, 0xff, 0xeb, 0x7b, 0x1e, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0xfd, 0xfe, 0xdd, 0xd9, 0x3b,
	0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MobileAppCategoryConstantServiceClient is the client API for MobileAppCategoryConstantService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MobileAppCategoryConstantServiceClient interface {
	// Returns the requested mobile app category constant.
	GetMobileAppCategoryConstant(ctx context.Context, in *GetMobileAppCategoryConstantRequest, opts ...grpc.CallOption) (*resources.MobileAppCategoryConstant, error)
}

type mobileAppCategoryConstantServiceClient struct {
	cc *grpc.ClientConn
}

func NewMobileAppCategoryConstantServiceClient(cc *grpc.ClientConn) MobileAppCategoryConstantServiceClient {
	return &mobileAppCategoryConstantServiceClient{cc}
}

func (c *mobileAppCategoryConstantServiceClient) GetMobileAppCategoryConstant(ctx context.Context, in *GetMobileAppCategoryConstantRequest, opts ...grpc.CallOption) (*resources.MobileAppCategoryConstant, error) {
	out := new(resources.MobileAppCategoryConstant)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.MobileAppCategoryConstantService/GetMobileAppCategoryConstant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MobileAppCategoryConstantServiceServer is the server API for MobileAppCategoryConstantService service.
type MobileAppCategoryConstantServiceServer interface {
	// Returns the requested mobile app category constant.
	GetMobileAppCategoryConstant(context.Context, *GetMobileAppCategoryConstantRequest) (*resources.MobileAppCategoryConstant, error)
}

// UnimplementedMobileAppCategoryConstantServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMobileAppCategoryConstantServiceServer struct {
}

func (*UnimplementedMobileAppCategoryConstantServiceServer) GetMobileAppCategoryConstant(ctx context.Context, req *GetMobileAppCategoryConstantRequest) (*resources.MobileAppCategoryConstant, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMobileAppCategoryConstant not implemented")
}

func RegisterMobileAppCategoryConstantServiceServer(s *grpc.Server, srv MobileAppCategoryConstantServiceServer) {
	s.RegisterService(&_MobileAppCategoryConstantService_serviceDesc, srv)
}

func _MobileAppCategoryConstantService_GetMobileAppCategoryConstant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMobileAppCategoryConstantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MobileAppCategoryConstantServiceServer).GetMobileAppCategoryConstant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.MobileAppCategoryConstantService/GetMobileAppCategoryConstant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MobileAppCategoryConstantServiceServer).GetMobileAppCategoryConstant(ctx, req.(*GetMobileAppCategoryConstantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MobileAppCategoryConstantService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v1.services.MobileAppCategoryConstantService",
	HandlerType: (*MobileAppCategoryConstantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMobileAppCategoryConstant",
			Handler:    _MobileAppCategoryConstantService_GetMobileAppCategoryConstant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v1/services/mobile_app_category_constant_service.proto",
}
