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

type AccessTokenRequest struct {
	UserID               string   `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccessTokenRequest) Reset()         { *m = AccessTokenRequest{} }
func (m *AccessTokenRequest) String() string { return proto.CompactTextString(m) }
func (*AccessTokenRequest) ProtoMessage()    {}
func (*AccessTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef9903376c2df4a6, []int{3}
}

func (m *AccessTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccessTokenRequest.Unmarshal(m, b)
}
func (m *AccessTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccessTokenRequest.Marshal(b, m, deterministic)
}
func (m *AccessTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessTokenRequest.Merge(m, src)
}
func (m *AccessTokenRequest) XXX_Size() int {
	return xxx_messageInfo_AccessTokenRequest.Size(m)
}
func (m *AccessTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccessTokenRequest proto.InternalMessageInfo

func (m *AccessTokenRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type AccessTokenResponse struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=AccessToken,proto3" json:"AccessToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccessTokenResponse) Reset()         { *m = AccessTokenResponse{} }
func (m *AccessTokenResponse) String() string { return proto.CompactTextString(m) }
func (*AccessTokenResponse) ProtoMessage()    {}
func (*AccessTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef9903376c2df4a6, []int{4}
}

func (m *AccessTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccessTokenResponse.Unmarshal(m, b)
}
func (m *AccessTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccessTokenResponse.Marshal(b, m, deterministic)
}
func (m *AccessTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessTokenResponse.Merge(m, src)
}
func (m *AccessTokenResponse) XXX_Size() int {
	return xxx_messageInfo_AccessTokenResponse.Size(m)
}
func (m *AccessTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AccessTokenResponse proto.InternalMessageInfo

func (m *AccessTokenResponse) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func init() {
	proto.RegisterType((*PermissionData)(nil), "proto.PermissionData")
	proto.RegisterType((*TokenData)(nil), "proto.TokenData")
	proto.RegisterType((*Error)(nil), "proto.Error")
	proto.RegisterType((*AccessTokenRequest)(nil), "proto.AccessTokenRequest")
	proto.RegisterType((*AccessTokenResponse)(nil), "proto.AccessTokenResponse")
}

func init() { proto.RegisterFile("permission_service.proto", fileDescriptor_ef9903376c2df4a6) }

var fileDescriptor_ef9903376c2df4a6 = []byte{
	// 279 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x6d, 0x84, 0x54, 0x3a, 0x5a, 0xd1, 0x29, 0x4a, 0xcc, 0x41, 0xc2, 0x9e, 0x7a, 0xd0, 0x82,
	0x0a, 0x7a, 0x2e, 0x56, 0xd4, 0x43, 0x45, 0x52, 0xf5, 0x2a, 0x6b, 0x3a, 0x48, 0x50, 0x33, 0x71,
	0x67, 0xeb, 0xc1, 0x8f, 0xf4, 0x9b, 0xc4, 0x64, 0xdb, 0xa4, 0x44, 0xb0, 0xa7, 0x65, 0xde, 0xbc,
	0x79, 0xb3, 0xef, 0x0d, 0x04, 0x39, 0x99, 0xf7, 0x54, 0x24, 0xe5, 0xec, 0x49, 0xc8, 0x7c, 0xa6,
	0x09, 0x0d, 0x72, 0xc3, 0x96, 0xd1, 0x2f, 0x1e, 0x15, 0xc3, 0xd6, 0xdd, 0x82, 0x32, 0xd2, 0x56,
	0x63, 0x04, 0x1b, 0xc3, 0x24, 0x21, 0x91, 0x7b, 0x7e, 0xa5, 0x2c, 0xf0, 0x22, 0xaf, 0xdf, 0x89,
	0xeb, 0x10, 0x1e, 0x00, 0x54, 0x33, 0xc1, 0x5a, 0xe4, 0xf5, 0xfd, 0xb8, 0x86, 0xa8, 0x23, 0xe8,
	0x14, 0xc4, 0xd5, 0xe4, 0xd4, 0x10, 0xfc, 0x4b, 0x63, 0xd8, 0x60, 0x00, 0xeb, 0x63, 0x12, 0xd1,
	0x2f, 0xe4, 0x68, 0xf3, 0xf2, 0x77, 0xe3, 0xc4, 0x6a, 0x3b, 0x93, 0x0b, 0x9e, 0xd2, 0x7c, 0x63,
	0x85, 0xa8, 0x43, 0xc0, 0x9a, 0x62, 0x4c, 0x1f, 0x33, 0x12, 0x8b, 0x7b, 0xd0, 0x7e, 0x10, 0x32,
	0x37, 0x23, 0x27, 0xe7, 0x2a, 0x75, 0x0e, 0xbd, 0x25, 0xb6, 0xe4, 0x9c, 0x09, 0xfd, 0xff, 0xd3,
	0x93, 0x6f, 0x0f, 0x76, 0x2a, 0x9f, 0x93, 0x32, 0x4f, 0x3c, 0x86, 0xee, 0xa3, 0x7e, 0x4b, 0xa7,
	0xda, 0x52, 0x99, 0xcf, 0x76, 0x19, 0xf1, 0x60, 0x11, 0x42, 0xb8, 0xe9, 0x90, 0xc2, 0xa7, 0x6a,
	0xe1, 0x19, 0x74, 0xaf, 0xb5, 0x54, 0x52, 0xb8, 0xeb, 0x08, 0xcb, 0xb7, 0x68, 0xcc, 0xdd, 0x42,
	0xef, 0x8a, 0xec, 0x98, 0xb3, 0x2f, 0xae, 0x1f, 0x64, 0xdf, 0xd1, 0x9a, 0x19, 0x84, 0xe1, 0x5f,
	0xad, 0xd2, 0xb0, 0x6a, 0x3d, 0xb7, 0x8b, 0xe6, 0xe9, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x59,
	0xba, 0x47, 0x0a, 0x27, 0x02, 0x00, 0x00,
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
	GetMonzoAccessToken(ctx context.Context, in *AccessTokenRequest, opts ...grpc.CallOption) (*AccessTokenResponse, error)
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

func (c *permissionServiceClient) GetMonzoAccessToken(ctx context.Context, in *AccessTokenRequest, opts ...grpc.CallOption) (*AccessTokenResponse, error) {
	out := new(AccessTokenResponse)
	err := c.cc.Invoke(ctx, "/proto.PermissionService/GetMonzoAccessToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServiceServer is the server API for PermissionService service.
type PermissionServiceServer interface {
	ValidateToken(context.Context, *TokenData) (*Error, error)
	HasPermission(context.Context, *PermissionData) (*Error, error)
	GetMonzoAccessToken(context.Context, *AccessTokenRequest) (*AccessTokenResponse, error)
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
func (*UnimplementedPermissionServiceServer) GetMonzoAccessToken(ctx context.Context, req *AccessTokenRequest) (*AccessTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMonzoAccessToken not implemented")
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

func _PermissionService_GetMonzoAccessToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).GetMonzoAccessToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PermissionService/GetMonzoAccessToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).GetMonzoAccessToken(ctx, req.(*AccessTokenRequest))
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
		{
			MethodName: "GetMonzoAccessToken",
			Handler:    _PermissionService_GetMonzoAccessToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "permission_service.proto",
}