// Code generated by protoc-gen-go. DO NOT EDIT.
// source: permission_service.proto

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

type PermissionData struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=AccessToken,proto3" json:"AccessToken,omitempty"`
	Permission           int32    `protobuf:"varint,2,opt,name=Permission,proto3" json:"Permission,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PermissionData) Reset()         { *m = PermissionData{} }
func (m *PermissionData) String() string { return proto.CompactTextString(m) }
func (*PermissionData) ProtoMessage()    {}
func (*PermissionData) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef9903376c2df4a6, []int{0}
}

func (m *PermissionData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PermissionData.Unmarshal(m, b)
}
func (m *PermissionData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PermissionData.Marshal(b, m, deterministic)
}
func (m *PermissionData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PermissionData.Merge(m, src)
}
func (m *PermissionData) XXX_Size() int {
	return xxx_messageInfo_PermissionData.Size(m)
}
func (m *PermissionData) XXX_DiscardUnknown() {
	xxx_messageInfo_PermissionData.DiscardUnknown(m)
}

var xxx_messageInfo_PermissionData proto.InternalMessageInfo

func (m *PermissionData) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *PermissionData) GetPermission() int32 {
	if m != nil {
		return m.Permission
	}
	return 0
}

type TokenData struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=AccessToken,proto3" json:"AccessToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenData) Reset()         { *m = TokenData{} }
func (m *TokenData) String() string { return proto.CompactTextString(m) }
func (*TokenData) ProtoMessage()    {}
func (*TokenData) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef9903376c2df4a6, []int{1}
}

func (m *TokenData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenData.Unmarshal(m, b)
}
func (m *TokenData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenData.Marshal(b, m, deterministic)
}
func (m *TokenData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenData.Merge(m, src)
}
func (m *TokenData) XXX_Size() int {
	return xxx_messageInfo_TokenData.Size(m)
}
func (m *TokenData) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenData.DiscardUnknown(m)
}

var xxx_messageInfo_TokenData proto.InternalMessageInfo

func (m *TokenData) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type Error struct {
	Message              string   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	StatusCode           int32    `protobuf:"varint,2,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef9903376c2df4a6, []int{2}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Error) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func init() {
	proto.RegisterType((*PermissionData)(nil), "proto.PermissionData")
	proto.RegisterType((*TokenData)(nil), "proto.TokenData")
	proto.RegisterType((*Error)(nil), "proto.Error")
}

func init() { proto.RegisterFile("permission_service.proto", fileDescriptor_ef9903376c2df4a6) }

var fileDescriptor_ef9903376c2df4a6 = []byte{
	// 213 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x28, 0x48, 0x2d, 0xca,
	0xcd, 0x2c, 0x2e, 0xce, 0xcc, 0xcf, 0x8b, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x41, 0x5c, 0x7c, 0x01, 0x70, 0x25, 0x2e,
	0x89, 0x25, 0x89, 0x42, 0x0a, 0x5c, 0xdc, 0x8e, 0xc9, 0xc9, 0xa9, 0xc5, 0xc5, 0x21, 0xf9, 0xd9,
	0xa9, 0x79, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0xc8, 0x42, 0x42, 0x72, 0x5c, 0x5c, 0x08,
	0x3d, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xac, 0x41, 0x48, 0x22, 0x4a, 0xba, 0x5c, 0x9c, 0x60, 0x85,
	0xc4, 0x19, 0xa7, 0xe4, 0xc8, 0xc5, 0xea, 0x5a, 0x54, 0x94, 0x5f, 0x24, 0x24, 0xc1, 0xc5, 0xee,
	0x9b, 0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0x0a, 0x55, 0x06, 0xe3, 0x82, 0x6c, 0x0c, 0x2e, 0x49, 0x2c,
	0x29, 0x2d, 0x76, 0xce, 0x4f, 0x49, 0x85, 0xd9, 0x88, 0x10, 0x31, 0xaa, 0xe3, 0x12, 0x44, 0xd8,
	0x1f, 0x0c, 0xf1, 0xa7, 0x90, 0x21, 0x17, 0x6f, 0x58, 0x62, 0x4e, 0x66, 0x4a, 0x62, 0x49, 0x2a,
	0xc4, 0xdd, 0x02, 0x10, 0xaf, 0xeb, 0xc1, 0x1d, 0x27, 0xc5, 0x03, 0x15, 0x01, 0xdb, 0xaf, 0xc4,
	0x20, 0x64, 0xc6, 0xc5, 0xeb, 0x91, 0x58, 0x8c, 0x30, 0x4a, 0x48, 0x14, 0xaa, 0x00, 0x35, 0x8c,
	0xd0, 0xf5, 0x25, 0xb1, 0x81, 0xb9, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3a, 0xbc, 0x01,
	0x78, 0x6f, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PermissionServiceClient is the client API for PermissionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PermissionServiceClient interface {
	ValidateToken(ctx context.Context, in *TokenData, opts ...grpc.CallOption) (*Error, error)
	HasPermission(ctx context.Context, in *PermissionData, opts ...grpc.CallOption) (*Error, error)
}

type permissionServiceClient struct {
	cc *grpc.ClientConn
}

func NewPermissionServiceClient(cc *grpc.ClientConn) PermissionServiceClient {
	return &permissionServiceClient{cc}
}

func (c *permissionServiceClient) ValidateToken(ctx context.Context, in *TokenData, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/proto.PermissionService/ValidateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionServiceClient) HasPermission(ctx context.Context, in *PermissionData, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/proto.PermissionService/HasPermission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServiceServer is the server API for PermissionService service.
type PermissionServiceServer interface {
	ValidateToken(context.Context, *TokenData) (*Error, error)
	HasPermission(context.Context, *PermissionData) (*Error, error)
}

// UnimplementedPermissionServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPermissionServiceServer struct {
}

func (*UnimplementedPermissionServiceServer) ValidateToken(ctx context.Context, req *TokenData) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (*UnimplementedPermissionServiceServer) HasPermission(ctx context.Context, req *PermissionData) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasPermission not implemented")
}

func RegisterPermissionServiceServer(s *grpc.Server, srv PermissionServiceServer) {
	s.RegisterService(&_PermissionService_serviceDesc, srv)
}

func _PermissionService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PermissionService/ValidateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).ValidateToken(ctx, req.(*TokenData))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionService_HasPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PermissionData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).HasPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PermissionService/HasPermission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).HasPermission(ctx, req.(*PermissionData))
	}
	return interceptor(ctx, in, info, handler)
}

var _PermissionService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.PermissionService",
	HandlerType: (*PermissionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValidateToken",
			Handler:    _PermissionService_ValidateToken_Handler,
		},
		{
			MethodName: "HasPermission",
			Handler:    _PermissionService_HasPermission_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "permission_service.proto",
}
