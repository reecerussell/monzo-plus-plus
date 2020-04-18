// Code generated by protoc-gen-go. DO NOT EDIT.
// source: plugin_service.proto

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

type SendRequest struct {
	UserID               string   `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	JSONData             string   `protobuf:"bytes,2,opt,name=JSONData,proto3" json:"JSONData,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendRequest) Reset()         { *m = SendRequest{} }
func (m *SendRequest) String() string { return proto.CompactTextString(m) }
func (*SendRequest) ProtoMessage()    {}
func (*SendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0b792d6030e82303, []int{0}
}

func (m *SendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendRequest.Unmarshal(m, b)
}
func (m *SendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendRequest.Marshal(b, m, deterministic)
}
func (m *SendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendRequest.Merge(m, src)
}
func (m *SendRequest) XXX_Size() int {
	return xxx_messageInfo_SendRequest.Size(m)
}
func (m *SendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendRequest proto.InternalMessageInfo

func (m *SendRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *SendRequest) GetJSONData() string {
	if m != nil {
		return m.JSONData
	}
	return ""
}

type EmptySendResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptySendResponse) Reset()         { *m = EmptySendResponse{} }
func (m *EmptySendResponse) String() string { return proto.CompactTextString(m) }
func (*EmptySendResponse) ProtoMessage()    {}
func (*EmptySendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0b792d6030e82303, []int{1}
}

func (m *EmptySendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptySendResponse.Unmarshal(m, b)
}
func (m *EmptySendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptySendResponse.Marshal(b, m, deterministic)
}
func (m *EmptySendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptySendResponse.Merge(m, src)
}
func (m *EmptySendResponse) XXX_Size() int {
	return xxx_messageInfo_EmptySendResponse.Size(m)
}
func (m *EmptySendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptySendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmptySendResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SendRequest)(nil), "proto.SendRequest")
	proto.RegisterType((*EmptySendResponse)(nil), "proto.EmptySendResponse")
}

func init() { proto.RegisterFile("plugin_service.proto", fileDescriptor_0b792d6030e82303) }

var fileDescriptor_0b792d6030e82303 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xc8, 0x29, 0x4d,
	0xcf, 0xcc, 0x8b, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x05, 0x53, 0x4a, 0x8e, 0x5c, 0xdc, 0xc1, 0xa9, 0x79, 0x29, 0x41, 0xa9, 0x85, 0xa5,
	0xa9, 0xc5, 0x25, 0x42, 0x62, 0x5c, 0x6c, 0xa1, 0xc5, 0xa9, 0x45, 0x9e, 0x2e, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0x9c, 0x41, 0x50, 0x9e, 0x90, 0x14, 0x17, 0x87, 0x57, 0xb0, 0xbf, 0x9f, 0x4b, 0x62,
	0x49, 0xa2, 0x04, 0x13, 0x58, 0x06, 0xce, 0x57, 0x12, 0xe6, 0x12, 0x74, 0xcd, 0x2d, 0x28, 0xa9,
	0x84, 0x98, 0x53, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x6a, 0xe4, 0xce, 0xc5, 0x1b, 0x00, 0xb6, 0x36,
	0x18, 0x62, 0xab, 0x90, 0x19, 0x17, 0x0b, 0x48, 0x81, 0x90, 0x10, 0xc4, 0x7e, 0x3d, 0x24, 0x5b,
	0xa5, 0x24, 0xa0, 0x62, 0x18, 0xc6, 0x28, 0x31, 0x24, 0xb1, 0x81, 0xa5, 0x8c, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xf3, 0x1d, 0xce, 0xbe, 0xc6, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PluginServiceClient is the client API for PluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PluginServiceClient interface {
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*EmptySendResponse, error)
}

type pluginServiceClient struct {
	cc *grpc.ClientConn
}

func NewPluginServiceClient(cc *grpc.ClientConn) PluginServiceClient {
	return &pluginServiceClient{cc}
}

func (c *pluginServiceClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*EmptySendResponse, error) {
	out := new(EmptySendResponse)
	err := c.cc.Invoke(ctx, "/proto.PluginService/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServiceServer is the server API for PluginService service.
type PluginServiceServer interface {
	Send(context.Context, *SendRequest) (*EmptySendResponse, error)
}

// UnimplementedPluginServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPluginServiceServer struct {
}

func (*UnimplementedPluginServiceServer) Send(ctx context.Context, req *SendRequest) (*EmptySendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}

func RegisterPluginServiceServer(s *grpc.Server, srv PluginServiceServer) {
	s.RegisterService(&_PluginService_serviceDesc, srv)
}

func _PluginService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServiceServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PluginService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.PluginService",
	HandlerType: (*PluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _PluginService_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin_service.proto",
}
