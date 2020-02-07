// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cleanup.proto

package gitalypb

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

type ApplyBfgObjectMapStreamRequest struct {
	// Only available on the first message
	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	// A raw object-map file as generated by BFG: https://rtyley.github.io/bfg-repo-cleaner
	// Each line in the file has two object SHAs, space-separated - the original
	// SHA of the object, and the SHA after BFG has rewritten the object.
	ObjectMap            []byte   `protobuf:"bytes,2,opt,name=object_map,json=objectMap,proto3" json:"object_map,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApplyBfgObjectMapStreamRequest) Reset()         { *m = ApplyBfgObjectMapStreamRequest{} }
func (m *ApplyBfgObjectMapStreamRequest) String() string { return proto.CompactTextString(m) }
func (*ApplyBfgObjectMapStreamRequest) ProtoMessage()    {}
func (*ApplyBfgObjectMapStreamRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b19e990e4662c9c, []int{0}
}

func (m *ApplyBfgObjectMapStreamRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplyBfgObjectMapStreamRequest.Unmarshal(m, b)
}
func (m *ApplyBfgObjectMapStreamRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplyBfgObjectMapStreamRequest.Marshal(b, m, deterministic)
}
func (m *ApplyBfgObjectMapStreamRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplyBfgObjectMapStreamRequest.Merge(m, src)
}
func (m *ApplyBfgObjectMapStreamRequest) XXX_Size() int {
	return xxx_messageInfo_ApplyBfgObjectMapStreamRequest.Size(m)
}
func (m *ApplyBfgObjectMapStreamRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplyBfgObjectMapStreamRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ApplyBfgObjectMapStreamRequest proto.InternalMessageInfo

func (m *ApplyBfgObjectMapStreamRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *ApplyBfgObjectMapStreamRequest) GetObjectMap() []byte {
	if m != nil {
		return m.ObjectMap
	}
	return nil
}

type ApplyBfgObjectMapStreamResponse struct {
	Entries              []*ApplyBfgObjectMapStreamResponse_Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                 `json:"-"`
	XXX_unrecognized     []byte                                   `json:"-"`
	XXX_sizecache        int32                                    `json:"-"`
}

func (m *ApplyBfgObjectMapStreamResponse) Reset()         { *m = ApplyBfgObjectMapStreamResponse{} }
func (m *ApplyBfgObjectMapStreamResponse) String() string { return proto.CompactTextString(m) }
func (*ApplyBfgObjectMapStreamResponse) ProtoMessage()    {}
func (*ApplyBfgObjectMapStreamResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b19e990e4662c9c, []int{1}
}

func (m *ApplyBfgObjectMapStreamResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplyBfgObjectMapStreamResponse.Unmarshal(m, b)
}
func (m *ApplyBfgObjectMapStreamResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplyBfgObjectMapStreamResponse.Marshal(b, m, deterministic)
}
func (m *ApplyBfgObjectMapStreamResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplyBfgObjectMapStreamResponse.Merge(m, src)
}
func (m *ApplyBfgObjectMapStreamResponse) XXX_Size() int {
	return xxx_messageInfo_ApplyBfgObjectMapStreamResponse.Size(m)
}
func (m *ApplyBfgObjectMapStreamResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplyBfgObjectMapStreamResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ApplyBfgObjectMapStreamResponse proto.InternalMessageInfo

func (m *ApplyBfgObjectMapStreamResponse) GetEntries() []*ApplyBfgObjectMapStreamResponse_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

// We send back each parsed entry in the request's object map so the client
// can take action
type ApplyBfgObjectMapStreamResponse_Entry struct {
	Type                 ObjectType `protobuf:"varint,1,opt,name=type,proto3,enum=gitaly.ObjectType" json:"type,omitempty"`
	OldOid               string     `protobuf:"bytes,2,opt,name=old_oid,json=oldOid,proto3" json:"old_oid,omitempty"`
	NewOid               string     `protobuf:"bytes,3,opt,name=new_oid,json=newOid,proto3" json:"new_oid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ApplyBfgObjectMapStreamResponse_Entry) Reset()         { *m = ApplyBfgObjectMapStreamResponse_Entry{} }
func (m *ApplyBfgObjectMapStreamResponse_Entry) String() string { return proto.CompactTextString(m) }
func (*ApplyBfgObjectMapStreamResponse_Entry) ProtoMessage()    {}
func (*ApplyBfgObjectMapStreamResponse_Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b19e990e4662c9c, []int{1, 0}
}

func (m *ApplyBfgObjectMapStreamResponse_Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplyBfgObjectMapStreamResponse_Entry.Unmarshal(m, b)
}
func (m *ApplyBfgObjectMapStreamResponse_Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplyBfgObjectMapStreamResponse_Entry.Marshal(b, m, deterministic)
}
func (m *ApplyBfgObjectMapStreamResponse_Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplyBfgObjectMapStreamResponse_Entry.Merge(m, src)
}
func (m *ApplyBfgObjectMapStreamResponse_Entry) XXX_Size() int {
	return xxx_messageInfo_ApplyBfgObjectMapStreamResponse_Entry.Size(m)
}
func (m *ApplyBfgObjectMapStreamResponse_Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplyBfgObjectMapStreamResponse_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_ApplyBfgObjectMapStreamResponse_Entry proto.InternalMessageInfo

func (m *ApplyBfgObjectMapStreamResponse_Entry) GetType() ObjectType {
	if m != nil {
		return m.Type
	}
	return ObjectType_UNKNOWN
}

func (m *ApplyBfgObjectMapStreamResponse_Entry) GetOldOid() string {
	if m != nil {
		return m.OldOid
	}
	return ""
}

func (m *ApplyBfgObjectMapStreamResponse_Entry) GetNewOid() string {
	if m != nil {
		return m.NewOid
	}
	return ""
}

func init() {
	proto.RegisterType((*ApplyBfgObjectMapStreamRequest)(nil), "gitaly.ApplyBfgObjectMapStreamRequest")
	proto.RegisterType((*ApplyBfgObjectMapStreamResponse)(nil), "gitaly.ApplyBfgObjectMapStreamResponse")
	proto.RegisterType((*ApplyBfgObjectMapStreamResponse_Entry)(nil), "gitaly.ApplyBfgObjectMapStreamResponse.Entry")
}

func init() { proto.RegisterFile("cleanup.proto", fileDescriptor_1b19e990e4662c9c) }

var fileDescriptor_1b19e990e4662c9c = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0xc1, 0x4e, 0x32, 0x31,
	0x18, 0x4c, 0x81, 0x7f, 0xf9, 0x29, 0xc8, 0xa1, 0x17, 0x08, 0x89, 0x4a, 0x38, 0xe0, 0x1e, 0x64,
	0x97, 0xe0, 0xc5, 0xab, 0x18, 0xe3, 0xc9, 0x90, 0x14, 0x4f, 0x5e, 0x48, 0x77, 0xf7, 0x73, 0xad,
	0x59, 0xb6, 0xb5, 0x5b, 0x20, 0x7d, 0x03, 0xdf, 0x40, 0x9f, 0xc8, 0x9b, 0x2f, 0xe4, 0xc9, 0xd0,
	0xba, 0xea, 0x85, 0xe0, 0xad, 0x9d, 0xe9, 0xcc, 0x7c, 0xf3, 0x15, 0x1f, 0xc4, 0x19, 0xb0, 0x7c,
	0x25, 0x03, 0xa9, 0x84, 0x16, 0xc4, 0x4b, 0xb9, 0x66, 0x99, 0xe9, 0xb5, 0x8a, 0x07, 0xa6, 0x20,
	0x71, 0xe8, 0xc0, 0xe0, 0xa3, 0x0b, 0x29, 0x33, 0x33, 0xbd, 0x4f, 0x67, 0xd1, 0x23, 0xc4, 0xfa,
	0x86, 0xc9, 0xb9, 0x56, 0xc0, 0x96, 0x14, 0x9e, 0x56, 0x50, 0x68, 0x72, 0x8e, 0xb1, 0x02, 0x29,
	0x0a, 0xae, 0x85, 0x32, 0x5d, 0xd4, 0x47, 0x7e, 0x73, 0x42, 0x02, 0x67, 0x16, 0xd0, 0x6f, 0x66,
	0x5a, 0x7b, 0x7d, 0x3b, 0x45, 0xf4, 0xd7, 0x5b, 0x72, 0x88, 0xb1, 0xb0, 0x9e, 0x8b, 0x25, 0x93,
	0xdd, 0x4a, 0x1f, 0xf9, 0x2d, 0xda, 0x10, 0x65, 0xca, 0xe0, 0x1d, 0xe1, 0xe3, 0x9d, 0xd9, 0x85,
	0x14, 0x79, 0x01, 0xe4, 0x1a, 0xd7, 0x21, 0xd7, 0x8a, 0x43, 0xd1, 0x45, 0xfd, 0xaa, 0xdf, 0x9c,
	0x8c, 0xca, 0xe4, 0x3d, 0xca, 0xe0, 0x2a, 0xd7, 0xca, 0xd0, 0x52, 0xdd, 0x63, 0xf8, 0x9f, 0x45,
	0xc8, 0x10, 0xd7, 0xb4, 0x91, 0x60, 0x8b, 0xb4, 0x7f, 0x8a, 0x38, 0x9b, 0x5b, 0x23, 0x81, 0x5a,
	0x9e, 0x74, 0x70, 0x5d, 0x64, 0xc9, 0x42, 0xf0, 0xc4, 0x4e, 0xde, 0xa0, 0x9e, 0xc8, 0x92, 0x19,
	0x4f, 0xb6, 0x44, 0x0e, 0x1b, 0x4b, 0x54, 0x1d, 0x91, 0xc3, 0x66, 0xc6, 0x93, 0xc9, 0x33, 0xc2,
	0xed, 0x4b, 0xb7, 0xf2, 0x39, 0xa8, 0x35, 0x8f, 0x81, 0xac, 0x71, 0x67, 0xc7, 0x9c, 0x64, 0xb8,
	0xb7, 0x88, 0x5d, 0x7f, 0xef, 0xe4, 0x8f, 0x85, 0x07, 0xde, 0xc7, 0x8b, 0x5f, 0xf9, 0x8f, 0x7c,
	0x34, 0x46, 0xd3, 0xf1, 0xdd, 0x56, 0x95, 0xb1, 0x28, 0x88, 0xc5, 0x32, 0x74, 0xc7, 0x91, 0x50,
	0x69, 0xe8, 0xbc, 0x42, 0xfb, 0xf7, 0x61, 0x2a, 0xbe, 0xee, 0x32, 0x8a, 0x3c, 0x0b, 0x9d, 0x7d,
	0x06, 0x00, 0x00, 0xff, 0xff, 0x64, 0xf7, 0x38, 0x48, 0x35, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CleanupServiceClient is the client API for CleanupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CleanupServiceClient interface {
	ApplyBfgObjectMapStream(ctx context.Context, opts ...grpc.CallOption) (CleanupService_ApplyBfgObjectMapStreamClient, error)
}

type cleanupServiceClient struct {
	cc *grpc.ClientConn
}

func NewCleanupServiceClient(cc *grpc.ClientConn) CleanupServiceClient {
	return &cleanupServiceClient{cc}
}

func (c *cleanupServiceClient) ApplyBfgObjectMapStream(ctx context.Context, opts ...grpc.CallOption) (CleanupService_ApplyBfgObjectMapStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_CleanupService_serviceDesc.Streams[0], "/gitaly.CleanupService/ApplyBfgObjectMapStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &cleanupServiceApplyBfgObjectMapStreamClient{stream}
	return x, nil
}

type CleanupService_ApplyBfgObjectMapStreamClient interface {
	Send(*ApplyBfgObjectMapStreamRequest) error
	Recv() (*ApplyBfgObjectMapStreamResponse, error)
	grpc.ClientStream
}

type cleanupServiceApplyBfgObjectMapStreamClient struct {
	grpc.ClientStream
}

func (x *cleanupServiceApplyBfgObjectMapStreamClient) Send(m *ApplyBfgObjectMapStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cleanupServiceApplyBfgObjectMapStreamClient) Recv() (*ApplyBfgObjectMapStreamResponse, error) {
	m := new(ApplyBfgObjectMapStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CleanupServiceServer is the server API for CleanupService service.
type CleanupServiceServer interface {
	ApplyBfgObjectMapStream(CleanupService_ApplyBfgObjectMapStreamServer) error
}

// UnimplementedCleanupServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCleanupServiceServer struct {
}

func (*UnimplementedCleanupServiceServer) ApplyBfgObjectMapStream(srv CleanupService_ApplyBfgObjectMapStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ApplyBfgObjectMapStream not implemented")
}

func RegisterCleanupServiceServer(s *grpc.Server, srv CleanupServiceServer) {
	s.RegisterService(&_CleanupService_serviceDesc, srv)
}

func _CleanupService_ApplyBfgObjectMapStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CleanupServiceServer).ApplyBfgObjectMapStream(&cleanupServiceApplyBfgObjectMapStreamServer{stream})
}

type CleanupService_ApplyBfgObjectMapStreamServer interface {
	Send(*ApplyBfgObjectMapStreamResponse) error
	Recv() (*ApplyBfgObjectMapStreamRequest, error)
	grpc.ServerStream
}

type cleanupServiceApplyBfgObjectMapStreamServer struct {
	grpc.ServerStream
}

func (x *cleanupServiceApplyBfgObjectMapStreamServer) Send(m *ApplyBfgObjectMapStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cleanupServiceApplyBfgObjectMapStreamServer) Recv() (*ApplyBfgObjectMapStreamRequest, error) {
	m := new(ApplyBfgObjectMapStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _CleanupService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gitaly.CleanupService",
	HandlerType: (*CleanupServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ApplyBfgObjectMapStream",
			Handler:       _CleanupService_ApplyBfgObjectMapStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "cleanup.proto",
}
