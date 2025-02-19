// Code generated by protoc-gen-go.
// source: golang.yacloud.eu/apis/urlmapper/urlmapper.proto
// DO NOT EDIT!

/*
Package urlmapper is a generated protocol buffer package.

It is generated from these files:
	golang.yacloud.eu/apis/urlmapper/urlmapper.proto

It has these top-level messages:
	JsonMapping
	GetJsonMappingRequest
	JsonMappingResponse
	JsonMappingResponseList
	DomainList
	ServiceID
	AnyHostMapping
	AnyMappingRequest
	AnyMappingResponse
	AllMapping
	AllMappingList
	RPCMappingRequest
	RPCMapping
	MappingRequest
	Mapping
*/
package urlmapper

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
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

type JsonMapping struct {
	ID              uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Domain          string `protobuf:"bytes,2,opt,name=Domain" json:"Domain,omitempty"`
	Path            string `protobuf:"bytes,3,opt,name=Path" json:"Path,omitempty"`
	ServiceID       string `protobuf:"bytes,4,opt,name=ServiceID" json:"ServiceID,omitempty"`
	GroupID         string `protobuf:"bytes,5,opt,name=GroupID" json:"GroupID,omitempty"`
	FQDNServiceName string `protobuf:"bytes,6,opt,name=FQDNServiceName" json:"FQDNServiceName,omitempty"`
	ServiceName     string `protobuf:"bytes,7,opt,name=ServiceName" json:"ServiceName,omitempty"`
	RPC             string `protobuf:"bytes,8,opt,name=RPC" json:"RPC,omitempty"`
}

func (m *JsonMapping) Reset()                    { *m = JsonMapping{} }
func (m *JsonMapping) String() string            { return proto.CompactTextString(m) }
func (*JsonMapping) ProtoMessage()               {}
func (*JsonMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *JsonMapping) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *JsonMapping) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *JsonMapping) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *JsonMapping) GetServiceID() string {
	if m != nil {
		return m.ServiceID
	}
	return ""
}

func (m *JsonMapping) GetGroupID() string {
	if m != nil {
		return m.GroupID
	}
	return ""
}

func (m *JsonMapping) GetFQDNServiceName() string {
	if m != nil {
		return m.FQDNServiceName
	}
	return ""
}

func (m *JsonMapping) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *JsonMapping) GetRPC() string {
	if m != nil {
		return m.RPC
	}
	return ""
}

type GetJsonMappingRequest struct {
	Domain string `protobuf:"bytes,1,opt,name=Domain" json:"Domain,omitempty"`
	Path   string `protobuf:"bytes,2,opt,name=Path" json:"Path,omitempty"`
}

func (m *GetJsonMappingRequest) Reset()                    { *m = GetJsonMappingRequest{} }
func (m *GetJsonMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*GetJsonMappingRequest) ProtoMessage()               {}
func (*GetJsonMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetJsonMappingRequest) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *GetJsonMappingRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type JsonMappingResponse struct {
	Mapping     *JsonMapping `protobuf:"bytes,1,opt,name=Mapping" json:"Mapping,omitempty"`
	GRPCService string       `protobuf:"bytes,2,opt,name=GRPCService" json:"GRPCService,omitempty"`
}

func (m *JsonMappingResponse) Reset()                    { *m = JsonMappingResponse{} }
func (m *JsonMappingResponse) String() string            { return proto.CompactTextString(m) }
func (*JsonMappingResponse) ProtoMessage()               {}
func (*JsonMappingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *JsonMappingResponse) GetMapping() *JsonMapping {
	if m != nil {
		return m.Mapping
	}
	return nil
}

func (m *JsonMappingResponse) GetGRPCService() string {
	if m != nil {
		return m.GRPCService
	}
	return ""
}

type JsonMappingResponseList struct {
	Responses []*JsonMappingResponse `protobuf:"bytes,1,rep,name=Responses" json:"Responses,omitempty"`
}

func (m *JsonMappingResponseList) Reset()                    { *m = JsonMappingResponseList{} }
func (m *JsonMappingResponseList) String() string            { return proto.CompactTextString(m) }
func (*JsonMappingResponseList) ProtoMessage()               {}
func (*JsonMappingResponseList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *JsonMappingResponseList) GetResponses() []*JsonMappingResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

type DomainList struct {
	Domains []string `protobuf:"bytes,1,rep,name=Domains" json:"Domains,omitempty"`
}

func (m *DomainList) Reset()                    { *m = DomainList{} }
func (m *DomainList) String() string            { return proto.CompactTextString(m) }
func (*DomainList) ProtoMessage()               {}
func (*DomainList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DomainList) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

type ServiceID struct {
	ID string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
}

func (m *ServiceID) Reset()                    { *m = ServiceID{} }
func (m *ServiceID) String() string            { return proto.CompactTextString(m) }
func (*ServiceID) ProtoMessage()               {}
func (*ServiceID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ServiceID) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

// via _api this can get mapped in
type AnyHostMapping struct {
	ID              uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Path            string `protobuf:"bytes,2,opt,name=Path" json:"Path,omitempty"`
	ServiceID       string `protobuf:"bytes,3,opt,name=ServiceID" json:"ServiceID,omitempty"`
	ServiceName     string `protobuf:"bytes,4,opt,name=ServiceName" json:"ServiceName,omitempty"`
	FQDNServiceName string `protobuf:"bytes,5,opt,name=FQDNServiceName" json:"FQDNServiceName,omitempty"`
}

func (m *AnyHostMapping) Reset()                    { *m = AnyHostMapping{} }
func (m *AnyHostMapping) String() string            { return proto.CompactTextString(m) }
func (*AnyHostMapping) ProtoMessage()               {}
func (*AnyHostMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AnyHostMapping) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *AnyHostMapping) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *AnyHostMapping) GetServiceID() string {
	if m != nil {
		return m.ServiceID
	}
	return ""
}

func (m *AnyHostMapping) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *AnyHostMapping) GetFQDNServiceName() string {
	if m != nil {
		return m.FQDNServiceName
	}
	return ""
}

type AnyMappingRequest struct {
	Path        string `protobuf:"bytes,1,opt,name=Path" json:"Path,omitempty"`
	ServiceName string `protobuf:"bytes,2,opt,name=ServiceName" json:"ServiceName,omitempty"`
}

func (m *AnyMappingRequest) Reset()                    { *m = AnyMappingRequest{} }
func (m *AnyMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*AnyMappingRequest) ProtoMessage()               {}
func (*AnyMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *AnyMappingRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *AnyMappingRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

type AnyMappingResponse struct {
	Path string `protobuf:"bytes,1,opt,name=Path" json:"Path,omitempty"`
}

func (m *AnyMappingResponse) Reset()                    { *m = AnyMappingResponse{} }
func (m *AnyMappingResponse) String() string            { return proto.CompactTextString(m) }
func (*AnyMappingResponse) ProtoMessage()               {}
func (*AnyMappingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *AnyMappingResponse) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type AllMapping struct {
	ServiceID   string `protobuf:"bytes,1,opt,name=ServiceID" json:"ServiceID,omitempty"`
	ServiceName string `protobuf:"bytes,2,opt,name=ServiceName" json:"ServiceName,omitempty"`
	Domain      string `protobuf:"bytes,3,opt,name=Domain" json:"Domain,omitempty"`
	Path        string `protobuf:"bytes,4,opt,name=Path" json:"Path,omitempty"`
	RPC         string `protobuf:"bytes,5,opt,name=RPC" json:"RPC,omitempty"`
}

func (m *AllMapping) Reset()                    { *m = AllMapping{} }
func (m *AllMapping) String() string            { return proto.CompactTextString(m) }
func (*AllMapping) ProtoMessage()               {}
func (*AllMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *AllMapping) GetServiceID() string {
	if m != nil {
		return m.ServiceID
	}
	return ""
}

func (m *AllMapping) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *AllMapping) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *AllMapping) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *AllMapping) GetRPC() string {
	if m != nil {
		return m.RPC
	}
	return ""
}

type AllMappingList struct {
	AllMappings []*AllMapping `protobuf:"bytes,1,rep,name=AllMappings" json:"AllMappings,omitempty"`
}

func (m *AllMappingList) Reset()                    { *m = AllMappingList{} }
func (m *AllMappingList) String() string            { return proto.CompactTextString(m) }
func (*AllMappingList) ProtoMessage()               {}
func (*AllMappingList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *AllMappingList) GetAllMappings() []*AllMapping {
	if m != nil {
		return m.AllMappings
	}
	return nil
}

type RPCMappingRequest struct {
	Expose      bool   `protobuf:"varint,1,opt,name=Expose" json:"Expose,omitempty"`
	ServiceName string `protobuf:"bytes,2,opt,name=ServiceName" json:"ServiceName,omitempty"`
	RPCName     string `protobuf:"bytes,3,opt,name=RPCName" json:"RPCName,omitempty"`
}

func (m *RPCMappingRequest) Reset()                    { *m = RPCMappingRequest{} }
func (m *RPCMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*RPCMappingRequest) ProtoMessage()               {}
func (*RPCMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *RPCMappingRequest) GetExpose() bool {
	if m != nil {
		return m.Expose
	}
	return false
}

func (m *RPCMappingRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *RPCMappingRequest) GetRPCName() string {
	if m != nil {
		return m.RPCName
	}
	return ""
}

type RPCMapping struct {
	ID          uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	ServiceName string `protobuf:"bytes,2,opt,name=ServiceName" json:"ServiceName,omitempty"`
	FQDNService string `protobuf:"bytes,3,opt,name=FQDNService" json:"FQDNService,omitempty"`
	RPCName     string `protobuf:"bytes,4,opt,name=RPCName" json:"RPCName,omitempty"`
}

func (m *RPCMapping) Reset()                    { *m = RPCMapping{} }
func (m *RPCMapping) String() string            { return proto.CompactTextString(m) }
func (*RPCMapping) ProtoMessage()               {}
func (*RPCMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *RPCMapping) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *RPCMapping) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *RPCMapping) GetFQDNService() string {
	if m != nil {
		return m.FQDNService
	}
	return ""
}

func (m *RPCMapping) GetRPCName() string {
	if m != nil {
		return m.RPCName
	}
	return ""
}

type MappingRequest struct {
	Path string `protobuf:"bytes,1,opt,name=Path" json:"Path,omitempty"`
}

func (m *MappingRequest) Reset()                    { *m = MappingRequest{} }
func (m *MappingRequest) String() string            { return proto.CompactTextString(m) }
func (*MappingRequest) ProtoMessage()               {}
func (*MappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *MappingRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type Mapping struct {
	MappingFound    bool   `protobuf:"varint,1,opt,name=MappingFound" json:"MappingFound,omitempty"`
	RegistryName    string `protobuf:"bytes,2,opt,name=RegistryName" json:"RegistryName,omitempty"`
	FQDNServiceName string `protobuf:"bytes,3,opt,name=FQDNServiceName" json:"FQDNServiceName,omitempty"`
	RPCName         string `protobuf:"bytes,4,opt,name=RPCName" json:"RPCName,omitempty"`
}

func (m *Mapping) Reset()                    { *m = Mapping{} }
func (m *Mapping) String() string            { return proto.CompactTextString(m) }
func (*Mapping) ProtoMessage()               {}
func (*Mapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *Mapping) GetMappingFound() bool {
	if m != nil {
		return m.MappingFound
	}
	return false
}

func (m *Mapping) GetRegistryName() string {
	if m != nil {
		return m.RegistryName
	}
	return ""
}

func (m *Mapping) GetFQDNServiceName() string {
	if m != nil {
		return m.FQDNServiceName
	}
	return ""
}

func (m *Mapping) GetRPCName() string {
	if m != nil {
		return m.RPCName
	}
	return ""
}

func init() {
	proto.RegisterType((*JsonMapping)(nil), "urlmapper.JsonMapping")
	proto.RegisterType((*GetJsonMappingRequest)(nil), "urlmapper.GetJsonMappingRequest")
	proto.RegisterType((*JsonMappingResponse)(nil), "urlmapper.JsonMappingResponse")
	proto.RegisterType((*JsonMappingResponseList)(nil), "urlmapper.JsonMappingResponseList")
	proto.RegisterType((*DomainList)(nil), "urlmapper.DomainList")
	proto.RegisterType((*ServiceID)(nil), "urlmapper.ServiceID")
	proto.RegisterType((*AnyHostMapping)(nil), "urlmapper.AnyHostMapping")
	proto.RegisterType((*AnyMappingRequest)(nil), "urlmapper.AnyMappingRequest")
	proto.RegisterType((*AnyMappingResponse)(nil), "urlmapper.AnyMappingResponse")
	proto.RegisterType((*AllMapping)(nil), "urlmapper.AllMapping")
	proto.RegisterType((*AllMappingList)(nil), "urlmapper.AllMappingList")
	proto.RegisterType((*RPCMappingRequest)(nil), "urlmapper.RPCMappingRequest")
	proto.RegisterType((*RPCMapping)(nil), "urlmapper.RPCMapping")
	proto.RegisterType((*MappingRequest)(nil), "urlmapper.MappingRequest")
	proto.RegisterType((*Mapping)(nil), "urlmapper.Mapping")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for URLMapper service

type URLMapperClient interface {
	// DEPRECATED - get json mapping for URL and user
	GetJsonMappingWithUser(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error)
	// DEPRECATED - get json mapping for URL
	GetJsonMapping(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error)
	// DEPRECATED - get all specific host json mappings (excluding anyhost mappings)
	GetJsonMappings(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*JsonMappingResponseList, error)
	// DEPRECATED - add a json mapping
	AddJsonMapping(ctx context.Context, in *JsonMapping, opts ...grpc.CallOption) (*common.Void, error)
	// DEPRECATED - get all domains "mapped" to jsonmultiplexer
	GetJsonDomains(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*DomainList, error)
	// add any host mapping
	AddAnyHostMapping(ctx context.Context, in *AnyMappingRequest, opts ...grpc.CallOption) (*AnyMappingResponse, error)
	// get all mappings, including anyhost and specific host mappings (useful for rendering list)
	GetAllMappings(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*AllMappingList, error)
	// Expose (or unexpose) a single RPC on any URL
	SetRPCMapping(ctx context.Context, in *RPCMappingRequest, opts ...grpc.CallOption) (*common.Void, error)
	// get the mapping for a full path, e.g. "/_api/golang.yacloud.eu/apis/vuehelper/Log"
	GetMapping(ctx context.Context, in *MappingRequest, opts ...grpc.CallOption) (*Mapping, error)
}

type uRLMapperClient struct {
	cc *grpc.ClientConn
}

func NewURLMapperClient(cc *grpc.ClientConn) URLMapperClient {
	return &uRLMapperClient{cc}
}

func (c *uRLMapperClient) GetJsonMappingWithUser(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error) {
	out := new(JsonMappingResponse)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonMappingWithUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetJsonMapping(ctx context.Context, in *GetJsonMappingRequest, opts ...grpc.CallOption) (*JsonMappingResponse, error) {
	out := new(JsonMappingResponse)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetJsonMappings(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*JsonMappingResponseList, error) {
	out := new(JsonMappingResponseList)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonMappings", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) AddJsonMapping(ctx context.Context, in *JsonMapping, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/AddJsonMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetJsonDomains(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*DomainList, error) {
	out := new(DomainList)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetJsonDomains", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) AddAnyHostMapping(ctx context.Context, in *AnyMappingRequest, opts ...grpc.CallOption) (*AnyMappingResponse, error) {
	out := new(AnyMappingResponse)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/AddAnyHostMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetAllMappings(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*AllMappingList, error) {
	out := new(AllMappingList)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetAllMappings", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) SetRPCMapping(ctx context.Context, in *RPCMappingRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/SetRPCMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLMapperClient) GetMapping(ctx context.Context, in *MappingRequest, opts ...grpc.CallOption) (*Mapping, error) {
	out := new(Mapping)
	err := grpc.Invoke(ctx, "/urlmapper.URLMapper/GetMapping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for URLMapper service

type URLMapperServer interface {
	// DEPRECATED - get json mapping for URL and user
	GetJsonMappingWithUser(context.Context, *GetJsonMappingRequest) (*JsonMappingResponse, error)
	// DEPRECATED - get json mapping for URL
	GetJsonMapping(context.Context, *GetJsonMappingRequest) (*JsonMappingResponse, error)
	// DEPRECATED - get all specific host json mappings (excluding anyhost mappings)
	GetJsonMappings(context.Context, *common.Void) (*JsonMappingResponseList, error)
	// DEPRECATED - add a json mapping
	AddJsonMapping(context.Context, *JsonMapping) (*common.Void, error)
	// DEPRECATED - get all domains "mapped" to jsonmultiplexer
	GetJsonDomains(context.Context, *common.Void) (*DomainList, error)
	// add any host mapping
	AddAnyHostMapping(context.Context, *AnyMappingRequest) (*AnyMappingResponse, error)
	// get all mappings, including anyhost and specific host mappings (useful for rendering list)
	GetAllMappings(context.Context, *common.Void) (*AllMappingList, error)
	// Expose (or unexpose) a single RPC on any URL
	SetRPCMapping(context.Context, *RPCMappingRequest) (*common.Void, error)
	// get the mapping for a full path, e.g. "/_api/golang.yacloud.eu/apis/vuehelper/Log"
	GetMapping(context.Context, *MappingRequest) (*Mapping, error)
}

func RegisterURLMapperServer(s *grpc.Server, srv URLMapperServer) {
	s.RegisterService(&_URLMapper_serviceDesc, srv)
}

func _URLMapper_GetJsonMappingWithUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJsonMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonMappingWithUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonMappingWithUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonMappingWithUser(ctx, req.(*GetJsonMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetJsonMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJsonMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonMapping(ctx, req.(*GetJsonMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetJsonMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonMappings(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_AddJsonMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JsonMapping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).AddJsonMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/AddJsonMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).AddJsonMapping(ctx, req.(*JsonMapping))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetJsonDomains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetJsonDomains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetJsonDomains",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetJsonDomains(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_AddAnyHostMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnyMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).AddAnyHostMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/AddAnyHostMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).AddAnyHostMapping(ctx, req.(*AnyMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetAllMappings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetAllMappings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetAllMappings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetAllMappings(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_SetRPCMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RPCMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).SetRPCMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/SetRPCMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).SetRPCMapping(ctx, req.(*RPCMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLMapper_GetMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLMapperServer).GetMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlmapper.URLMapper/GetMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLMapperServer).GetMapping(ctx, req.(*MappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _URLMapper_serviceDesc = grpc.ServiceDesc{
	ServiceName: "urlmapper.URLMapper",
	HandlerType: (*URLMapperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetJsonMappingWithUser",
			Handler:    _URLMapper_GetJsonMappingWithUser_Handler,
		},
		{
			MethodName: "GetJsonMapping",
			Handler:    _URLMapper_GetJsonMapping_Handler,
		},
		{
			MethodName: "GetJsonMappings",
			Handler:    _URLMapper_GetJsonMappings_Handler,
		},
		{
			MethodName: "AddJsonMapping",
			Handler:    _URLMapper_AddJsonMapping_Handler,
		},
		{
			MethodName: "GetJsonDomains",
			Handler:    _URLMapper_GetJsonDomains_Handler,
		},
		{
			MethodName: "AddAnyHostMapping",
			Handler:    _URLMapper_AddAnyHostMapping_Handler,
		},
		{
			MethodName: "GetAllMappings",
			Handler:    _URLMapper_GetAllMappings_Handler,
		},
		{
			MethodName: "SetRPCMapping",
			Handler:    _URLMapper_SetRPCMapping_Handler,
		},
		{
			MethodName: "GetMapping",
			Handler:    _URLMapper_GetMapping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.yacloud.eu/apis/urlmapper/urlmapper.proto",
}

func init() { proto.RegisterFile("golang.yacloud.eu/apis/urlmapper/urlmapper.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 752 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x56, 0xcd, 0x6e, 0xd3, 0x4a,
	0x14, 0x96, 0x13, 0xb7, 0x69, 0x4e, 0x7a, 0xd3, 0xdb, 0xb9, 0x6a, 0xaf, 0x1b, 0x0a, 0x8a, 0x46,
	0x08, 0x65, 0xe5, 0x56, 0x45, 0x80, 0x50, 0x61, 0x91, 0x26, 0x34, 0x04, 0xb5, 0x55, 0x98, 0xaa,
	0x54, 0xea, 0xce, 0xc4, 0xa3, 0xd4, 0x52, 0xe2, 0x31, 0xfe, 0x01, 0xc2, 0x92, 0x25, 0x7b, 0xde,
	0x01, 0xf1, 0x52, 0xbc, 0x05, 0x6b, 0x94, 0xc9, 0x38, 0x3e, 0xfe, 0x09, 0x61, 0xc1, 0x2a, 0x9e,
	0xf3, 0xff, 0xcd, 0xf9, 0xce, 0xc9, 0xc0, 0xe1, 0x48, 0x8c, 0x2d, 0x77, 0x64, 0x4e, 0xad, 0xe1,
	0x58, 0x44, 0xb6, 0xc9, 0xa3, 0x03, 0xcb, 0x73, 0x82, 0x83, 0xc8, 0x1f, 0x4f, 0x2c, 0xcf, 0xe3,
	0x7e, 0xf2, 0x65, 0x7a, 0xbe, 0x08, 0x05, 0xa9, 0x2e, 0x04, 0x0d, 0x53, 0x39, 0x0f, 0x85, 0xeb,
	0x5b, 0xf6, 0x07, 0x21, 0x6c, 0xd3, 0xe5, 0xe1, 0x3c, 0xc0, 0x50, 0x4c, 0x26, 0xc2, 0x55, 0x3f,
	0x73, 0x57, 0xfa, 0x43, 0x83, 0xda, 0xab, 0x40, 0xb8, 0xe7, 0x96, 0xe7, 0x39, 0xee, 0x88, 0xd4,
	0xa1, 0xd4, 0xef, 0x1a, 0x5a, 0x53, 0x6b, 0xe9, 0xac, 0xd4, 0xef, 0x92, 0x5d, 0x58, 0xef, 0x8a,
	0x89, 0xe5, 0xb8, 0x46, 0xa9, 0xa9, 0xb5, 0xaa, 0x4c, 0x9d, 0x08, 0x01, 0x7d, 0x60, 0x85, 0xb7,
	0x46, 0x59, 0x4a, 0xe5, 0x37, 0xd9, 0x87, 0xea, 0x25, 0xf7, 0xdf, 0x3b, 0x43, 0xde, 0xef, 0x1a,
	0xba, 0x54, 0x24, 0x02, 0x62, 0x40, 0xa5, 0xe7, 0x8b, 0xc8, 0xeb, 0x77, 0x8d, 0x35, 0xa9, 0x8b,
	0x8f, 0xa4, 0x05, 0x5b, 0xa7, 0xaf, 0xbb, 0x17, 0xca, 0xf4, 0xc2, 0x9a, 0x70, 0x63, 0x5d, 0x5a,
	0x64, 0xc5, 0xa4, 0x09, 0x35, 0x6c, 0x55, 0x91, 0x56, 0x58, 0x44, 0xfe, 0x85, 0x32, 0x1b, 0x74,
	0x8c, 0x0d, 0xa9, 0x99, 0x7d, 0xd2, 0x0e, 0xec, 0xf4, 0x78, 0x88, 0x30, 0x32, 0xfe, 0x2e, 0xe2,
	0x41, 0x88, 0xa0, 0x69, 0x85, 0xd0, 0x4a, 0x09, 0x34, 0xea, 0xc0, 0x7f, 0xa9, 0x08, 0x81, 0x27,
	0xdc, 0x80, 0x93, 0x43, 0xa8, 0x28, 0x91, 0x8c, 0x51, 0x3b, 0xda, 0x35, 0x93, 0xde, 0x60, 0x87,
	0xd8, 0x6c, 0x86, 0xa0, 0xc7, 0x06, 0x1d, 0x55, 0xb2, 0xca, 0x81, 0x45, 0xf4, 0x1a, 0xfe, 0x2f,
	0x48, 0x75, 0xe6, 0x04, 0x21, 0x79, 0x06, 0xd5, 0xf8, 0x1c, 0x18, 0x5a, 0xb3, 0xdc, 0xaa, 0x1d,
	0xdd, 0x5b, 0x92, 0x50, 0x99, 0xb1, 0xc4, 0x81, 0x3e, 0x00, 0x98, 0x23, 0x94, 0xb1, 0x0c, 0xa8,
	0xcc, 0x4f, 0xf3, 0x48, 0x55, 0x16, 0x1f, 0xe9, 0x1d, 0xd4, 0x46, 0xc4, 0x87, 0xea, 0x8c, 0x0f,
	0xf4, 0x9b, 0x06, 0xf5, 0xb6, 0x3b, 0x7d, 0x29, 0x82, 0x70, 0x19, 0x65, 0xf6, 0xf1, 0xfd, 0x9d,
	0x6c, 0x7c, 0xff, 0xbc, 0xa7, 0x87, 0x7e, 0xc4, 0x8b, 0x48, 0x52, 0xce, 0x92, 0x24, 0xd3, 0x60,
	0x3d, 0xdf, 0xe0, 0x02, 0xb2, 0xac, 0x15, 0x92, 0x85, 0xf6, 0x61, 0xbb, 0xed, 0x4e, 0x33, 0x4d,
	0x8f, 0x9b, 0xab, 0x21, 0xde, 0x66, 0x92, 0x96, 0x72, 0x49, 0x69, 0x0b, 0x08, 0x0e, 0xa5, 0xba,
	0x5f, 0x10, 0x8b, 0x7e, 0xd1, 0x00, 0xda, 0xe3, 0x71, 0x7c, 0x37, 0x29, 0xb4, 0xda, 0x0a, 0xb4,
	0xf9, 0xc4, 0x88, 0xa3, 0xe5, 0x42, 0x8e, 0xea, 0x08, 0x86, 0xa2, 0xfe, 0x5a, 0x42, 0xfd, 0x3e,
	0xd4, 0x93, 0x5a, 0x64, 0xd7, 0x9f, 0x40, 0x2d, 0x91, 0xc4, 0x1c, 0xda, 0x41, 0x1c, 0x4a, 0xb4,
	0x0c, 0x5b, 0xd2, 0x11, 0x6c, 0xb3, 0x41, 0x27, 0x3f, 0x41, 0x2f, 0x3e, 0x7a, 0x22, 0xe0, 0x12,
	0xda, 0x06, 0x53, 0xa7, 0x3f, 0xc0, 0x65, 0x40, 0x85, 0x0d, 0x3a, 0x52, 0x3b, 0x07, 0x16, 0x1f,
	0xe9, 0x27, 0x80, 0x24, 0x51, 0x8e, 0x5b, 0xab, 0x23, 0x37, 0xa1, 0x86, 0x88, 0xa0, 0xa2, 0x63,
	0x11, 0xce, 0xad, 0xa7, 0x73, 0xdf, 0x87, 0xfa, 0x6a, 0xba, 0xd0, 0xaf, 0xda, 0x62, 0xea, 0x09,
	0x85, 0x4d, 0xf5, 0x79, 0x2a, 0x22, 0xd7, 0x56, 0xf7, 0x90, 0x92, 0xcd, 0x6c, 0x18, 0x1f, 0x39,
	0x41, 0xe8, 0x4f, 0x51, 0xd1, 0x29, 0x59, 0x11, 0xab, 0xcb, 0xc5, 0x2b, 0x70, 0x69, 0xf5, 0x47,
	0x3f, 0x75, 0xa8, 0x5e, 0xb1, 0xb3, 0x73, 0xd9, 0x48, 0x72, 0x03, 0xbb, 0xe9, 0xb5, 0x77, 0xed,
	0x84, 0xb7, 0x57, 0x01, 0xf7, 0x49, 0x13, 0xb5, 0xbb, 0x70, 0x33, 0x36, 0x56, 0x2c, 0x15, 0xc2,
	0xa0, 0x9e, 0x76, 0xfc, 0x0b, 0x31, 0x3b, 0xb0, 0x95, 0x76, 0x0c, 0xc8, 0xa6, 0xa9, 0xfe, 0xaa,
	0xde, 0x08, 0xc7, 0x6e, 0xd0, 0xdf, 0x07, 0x90, 0xf4, 0x7e, 0x0c, 0xf5, 0xb6, 0x6d, 0xe3, 0xc2,
	0x96, 0x2c, 0xe4, 0x46, 0x2a, 0x36, 0x79, 0xb4, 0x00, 0xa4, 0x96, 0x60, 0x26, 0x37, 0x9e, 0x10,
	0xb4, 0x43, 0x07, 0xb0, 0xdd, 0xb6, 0xed, 0xcc, 0x3a, 0xdc, 0xc7, 0xd3, 0x94, 0xdd, 0x3f, 0x8d,
	0xbb, 0x4b, 0xb4, 0xea, 0x16, 0x9e, 0xca, 0x42, 0xd0, 0xe0, 0x65, 0x0a, 0xd9, 0x2b, 0x1c, 0x55,
	0x59, 0xcc, 0x31, 0xfc, 0x73, 0xc9, 0x43, 0x34, 0x3b, 0xb8, 0x90, 0xdc, 0xec, 0x66, 0x2e, 0xe0,
	0x18, 0xa0, 0xc7, 0x17, 0x10, 0x70, 0x96, 0x8c, 0x1b, 0xc9, 0xab, 0x4e, 0x9e, 0xc3, 0x1e, 0x8f,
	0x16, 0xcf, 0x95, 0xd9, 0x53, 0x23, 0x31, 0xba, 0x69, 0xae, 0x7a, 0xcd, 0xbc, 0x5d, 0x97, 0x2f,
	0x91, 0x87, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x6e, 0x1b, 0x9c, 0xf8, 0x08, 0x00, 0x00,
}
