// Code generated by protoc-gen-go. DO NOT EDIT.
// source: jobs_service.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

type PushRequest struct {
	UserID               string   `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	PluginID             string   `protobuf:"bytes,2,opt,name=PluginID,proto3" json:"PluginID,omitempty"`
	Data                 string   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushRequest) Reset()         { *m = PushRequest{} }
func (m *PushRequest) String() string { return proto.CompactTextString(m) }
func (*PushRequest) ProtoMessage()    {}
func (*PushRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fab68e5a84c56fd6, []int{0}
}

func (m *PushRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRequest.Unmarshal(m, b)
}
func (m *PushRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRequest.Marshal(b, m, deterministic)
}
func (m *PushRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRequest.Merge(m, src)
}
func (m *PushRequest) XXX_Size() int {
	return xxx_messageInfo_PushRequest.Size(m)
}
func (m *PushRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushRequest proto.InternalMessageInfo

func (m *PushRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *PushRequest) GetPluginID() string {
	if m != nil {
		return m.PluginID
	}
	return ""
}

func (m *PushRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type EmptyPushResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyPushResponse) Reset()         { *m = EmptyPushResponse{} }
func (m *EmptyPushResponse) String() string { return proto.CompactTextString(m) }
func (*EmptyPushResponse) ProtoMessage()    {}
func (*EmptyPushResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fab68e5a84c56fd6, []int{1}
}

func (m *EmptyPushResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyPushResponse.Unmarshal(m, b)
}
func (m *EmptyPushResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyPushResponse.Marshal(b, m, deterministic)
}
func (m *EmptyPushResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyPushResponse.Merge(m, src)
}
func (m *EmptyPushResponse) XXX_Size() int {
	return xxx_messageInfo_EmptyPushResponse.Size(m)
}
func (m *EmptyPushResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyPushResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyPushResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PushRequest)(nil), "proto.PushRequest")
	proto.RegisterType((*EmptyPushResponse)(nil), "proto.EmptyPushResponse")
}

func init() { proto.RegisterFile("jobs_service.proto", fileDescriptor_fab68e5a84c56fd6) }

var fileDescriptor_fab68e5a84c56fd6 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xca, 0x4f, 0x2a,
	0x8e, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x05, 0x53, 0x4a, 0xa1, 0x5c, 0xdc, 0x01, 0xa5, 0xc5, 0x19, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5,
	0x25, 0x42, 0x62, 0x5c, 0x6c, 0xa1, 0xc5, 0xa9, 0x45, 0x9e, 0x2e, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0x9c, 0x41, 0x50, 0x9e, 0x90, 0x14, 0x17, 0x47, 0x40, 0x4e, 0x69, 0x7a, 0x66, 0x9e, 0xa7, 0x8b,
	0x04, 0x13, 0x58, 0x06, 0xce, 0x17, 0x12, 0xe2, 0x62, 0x71, 0x49, 0x2c, 0x49, 0x94, 0x60, 0x06,
	0x8b, 0x83, 0xd9, 0x4a, 0xc2, 0x5c, 0x82, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x10, 0xb3, 0x8b, 0x0b,
	0xf2, 0xf3, 0x8a, 0x53, 0x8d, 0x5c, 0xb9, 0xb8, 0xbd, 0xf2, 0x93, 0x8a, 0x83, 0x21, 0xee, 0x10,
	0x32, 0xe3, 0x62, 0x01, 0x49, 0x0b, 0x09, 0x41, 0x5c, 0xa4, 0x87, 0xe4, 0x0e, 0x29, 0x09, 0xa8,
	0x18, 0x86, 0x21, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0x29, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xc7, 0xe2, 0xc1, 0xff, 0xd6, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JobsServiceClient is the client API for JobsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JobsServiceClient interface {
	Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*EmptyPushResponse, error)
}

type jobsServiceClient struct {
	cc *grpc.ClientConn
}

func NewJobsServiceClient(cc *grpc.ClientConn) JobsServiceClient {
	return &jobsServiceClient{cc}
}

func (c *jobsServiceClient) Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*EmptyPushResponse, error) {
	out := new(EmptyPushResponse)
	err := c.cc.Invoke(ctx, "/proto.JobsService/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JobsServiceServer is the server API for JobsService service.
type JobsServiceServer interface {
	Push(context.Context, *PushRequest) (*EmptyPushResponse, error)
}

// UnimplementedJobsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedJobsServiceServer struct {
}

func (*UnimplementedJobsServiceServer) Push(ctx context.Context, req *PushRequest) (*EmptyPushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}

func RegisterJobsServiceServer(s *grpc.Server, srv JobsServiceServer) {
	s.RegisterService(&_JobsService_serviceDesc, srv)
}

func _JobsService_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobsServiceServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.JobsService/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobsServiceServer).Push(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _JobsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.JobsService",
	HandlerType: (*JobsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _JobsService_Push_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "jobs_service.proto",
}
