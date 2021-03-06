// Code generated by protoc-gen-go. DO NOT EDIT.
// source: consensus/consensus.proto

package consensus

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import ledger "github.com/msaldanha/realChain/ledger"

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

type VoteRequest struct {
	SendTx               *ledger.Transaction `protobuf:"bytes,1,opt,name=sendTx,proto3" json:"sendTx,omitempty"`
	ReceiveTx            *ledger.Transaction `protobuf:"bytes,2,opt,name=receiveTx,proto3" json:"receiveTx,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *VoteRequest) Reset()         { *m = VoteRequest{} }
func (m *VoteRequest) String() string { return proto.CompactTextString(m) }
func (*VoteRequest) ProtoMessage()    {}
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_224b1349252bcb32, []int{0}
}
func (m *VoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteRequest.Unmarshal(m, b)
}
func (m *VoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteRequest.Marshal(b, m, deterministic)
}
func (dst *VoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteRequest.Merge(dst, src)
}
func (m *VoteRequest) XXX_Size() int {
	return xxx_messageInfo_VoteRequest.Size(m)
}
func (m *VoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VoteRequest proto.InternalMessageInfo

func (m *VoteRequest) GetSendTx() *ledger.Transaction {
	if m != nil {
		return m.SendTx
	}
	return nil
}

func (m *VoteRequest) GetReceiveTx() *ledger.Transaction {
	if m != nil {
		return m.ReceiveTx
	}
	return nil
}

type Vote struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	PubKey               []byte   `protobuf:"bytes,3,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	Signature            []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vote) Reset()         { *m = Vote{} }
func (m *Vote) String() string { return proto.CompactTextString(m) }
func (*Vote) ProtoMessage()    {}
func (*Vote) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_224b1349252bcb32, []int{1}
}
func (m *Vote) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vote.Unmarshal(m, b)
}
func (m *Vote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vote.Marshal(b, m, deterministic)
}
func (dst *Vote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vote.Merge(dst, src)
}
func (m *Vote) XXX_Size() int {
	return xxx_messageInfo_Vote.Size(m)
}
func (m *Vote) XXX_DiscardUnknown() {
	xxx_messageInfo_Vote.DiscardUnknown(m)
}

var xxx_messageInfo_Vote proto.InternalMessageInfo

func (m *Vote) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *Vote) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Vote) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Vote) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type VoteResult struct {
	Vote                 *Vote    `protobuf:"bytes,1,opt,name=vote,proto3" json:"vote,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteResult) Reset()         { *m = VoteResult{} }
func (m *VoteResult) String() string { return proto.CompactTextString(m) }
func (*VoteResult) ProtoMessage()    {}
func (*VoteResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_224b1349252bcb32, []int{2}
}
func (m *VoteResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteResult.Unmarshal(m, b)
}
func (m *VoteResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteResult.Marshal(b, m, deterministic)
}
func (dst *VoteResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteResult.Merge(dst, src)
}
func (m *VoteResult) XXX_Size() int {
	return xxx_messageInfo_VoteResult.Size(m)
}
func (m *VoteResult) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteResult.DiscardUnknown(m)
}

var xxx_messageInfo_VoteResult proto.InternalMessageInfo

func (m *VoteResult) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

type AcceptRequest struct {
	SendTx               *ledger.Transaction `protobuf:"bytes,1,opt,name=sendTx,proto3" json:"sendTx,omitempty"`
	ReceiveTx            *ledger.Transaction `protobuf:"bytes,2,opt,name=receiveTx,proto3" json:"receiveTx,omitempty"`
	Votes                []*Vote             `protobuf:"bytes,3,rep,name=votes,proto3" json:"votes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *AcceptRequest) Reset()         { *m = AcceptRequest{} }
func (m *AcceptRequest) String() string { return proto.CompactTextString(m) }
func (*AcceptRequest) ProtoMessage()    {}
func (*AcceptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_224b1349252bcb32, []int{3}
}
func (m *AcceptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptRequest.Unmarshal(m, b)
}
func (m *AcceptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptRequest.Marshal(b, m, deterministic)
}
func (dst *AcceptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptRequest.Merge(dst, src)
}
func (m *AcceptRequest) XXX_Size() int {
	return xxx_messageInfo_AcceptRequest.Size(m)
}
func (m *AcceptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptRequest proto.InternalMessageInfo

func (m *AcceptRequest) GetSendTx() *ledger.Transaction {
	if m != nil {
		return m.SendTx
	}
	return nil
}

func (m *AcceptRequest) GetReceiveTx() *ledger.Transaction {
	if m != nil {
		return m.ReceiveTx
	}
	return nil
}

func (m *AcceptRequest) GetVotes() []*Vote {
	if m != nil {
		return m.Votes
	}
	return nil
}

type AcceptResult struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcceptResult) Reset()         { *m = AcceptResult{} }
func (m *AcceptResult) String() string { return proto.CompactTextString(m) }
func (*AcceptResult) ProtoMessage()    {}
func (*AcceptResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_224b1349252bcb32, []int{4}
}
func (m *AcceptResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcceptResult.Unmarshal(m, b)
}
func (m *AcceptResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcceptResult.Marshal(b, m, deterministic)
}
func (dst *AcceptResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcceptResult.Merge(dst, src)
}
func (m *AcceptResult) XXX_Size() int {
	return xxx_messageInfo_AcceptResult.Size(m)
}
func (m *AcceptResult) XXX_DiscardUnknown() {
	xxx_messageInfo_AcceptResult.DiscardUnknown(m)
}

var xxx_messageInfo_AcceptResult proto.InternalMessageInfo

func init() {
	proto.RegisterType((*VoteRequest)(nil), "VoteRequest")
	proto.RegisterType((*Vote)(nil), "Vote")
	proto.RegisterType((*VoteResult)(nil), "VoteResult")
	proto.RegisterType((*AcceptRequest)(nil), "AcceptRequest")
	proto.RegisterType((*AcceptResult)(nil), "AcceptResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConsensusClient is the client API for Consensus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConsensusClient interface {
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResult, error)
	Accept(ctx context.Context, in *AcceptRequest, opts ...grpc.CallOption) (*AcceptResult, error)
}

type consensusClient struct {
	cc *grpc.ClientConn
}

func NewConsensusClient(cc *grpc.ClientConn) ConsensusClient {
	return &consensusClient{cc}
}

func (c *consensusClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResult, error) {
	out := new(VoteResult)
	err := c.cc.Invoke(ctx, "/Consensus/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) Accept(ctx context.Context, in *AcceptRequest, opts ...grpc.CallOption) (*AcceptResult, error) {
	out := new(AcceptResult)
	err := c.cc.Invoke(ctx, "/Consensus/Accept", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusServer is the server API for Consensus service.
type ConsensusServer interface {
	Vote(context.Context, *VoteRequest) (*VoteResult, error)
	Accept(context.Context, *AcceptRequest) (*AcceptResult, error)
}

func RegisterConsensusServer(s *grpc.Server, srv ConsensusServer) {
	s.RegisterService(&_Consensus_serviceDesc, srv)
}

func _Consensus_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Consensus_Accept_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).Accept(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus/Accept",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).Accept(ctx, req.(*AcceptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Consensus_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Consensus",
	HandlerType: (*ConsensusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Vote",
			Handler:    _Consensus_Vote_Handler,
		},
		{
			MethodName: "Accept",
			Handler:    _Consensus_Accept_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "consensus/consensus.proto",
}

func init() {
	proto.RegisterFile("consensus/consensus.proto", fileDescriptor_consensus_224b1349252bcb32)
}

var fileDescriptor_consensus_224b1349252bcb32 = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x52, 0x4f, 0x4f, 0xbb, 0x40,
	0x10, 0x2d, 0x85, 0x92, 0x1f, 0x03, 0xe5, 0xb0, 0xbf, 0xc4, 0x6c, 0xd1, 0x03, 0x59, 0x0f, 0xd6,
	0x98, 0x6c, 0x63, 0xfd, 0x04, 0xc6, 0xa3, 0x37, 0x42, 0xbc, 0x78, 0xa2, 0x30, 0x69, 0x48, 0xeb,
	0x2e, 0xee, 0x2e, 0x4d, 0xfd, 0x0e, 0x7e, 0x68, 0xc3, 0x42, 0xa5, 0x1e, 0xbc, 0x7a, 0xdb, 0x79,
	0xf3, 0xe7, 0xbd, 0x37, 0xb3, 0xb0, 0x28, 0xa5, 0xd0, 0x28, 0x74, 0xab, 0x57, 0xdf, 0x2f, 0xde,
	0x28, 0x69, 0x64, 0x42, 0xf7, 0x58, 0x6d, 0x51, 0xad, 0x8c, 0x2a, 0x84, 0x2e, 0x4a, 0x53, 0x4b,
	0xd1, 0x67, 0xd8, 0x1b, 0x84, 0x2f, 0xd2, 0x60, 0x86, 0xef, 0x2d, 0x6a, 0x43, 0xee, 0xc0, 0xd7,
	0x28, 0xaa, 0xfc, 0x48, 0x9d, 0xd4, 0x59, 0x86, 0xeb, 0xff, 0xbc, 0xef, 0xe4, 0xf9, 0xd8, 0x99,
	0x0d, 0x25, 0xe4, 0x1e, 0x02, 0x85, 0x25, 0xd6, 0x07, 0xcc, 0x8f, 0x74, 0xfa, 0x7b, 0xfd, 0x58,
	0xc5, 0x2a, 0xf0, 0x3a, 0x3a, 0x12, 0xc3, 0x54, 0xee, 0x2c, 0xc7, 0xbf, 0x6c, 0x2a, 0x77, 0xe4,
	0x02, 0x7c, 0x85, 0x85, 0x96, 0xc2, 0xce, 0x09, 0xb2, 0x21, 0xea, 0xf0, 0xa6, 0xdd, 0x3c, 0xe3,
	0x07, 0x75, 0x53, 0x67, 0x19, 0x65, 0x43, 0x44, 0xae, 0x20, 0xd0, 0xf5, 0x56, 0x14, 0xa6, 0x55,
	0x48, 0x3d, 0x9b, 0x1a, 0x01, 0x76, 0x03, 0xd0, 0x9b, 0xd2, 0xed, 0xde, 0x90, 0x05, 0x78, 0x07,
	0x69, 0x70, 0x70, 0x34, 0xe3, 0x36, 0x65, 0x21, 0xf6, 0xe9, 0xc0, 0xfc, 0xb1, 0x2c, 0xb1, 0x31,
	0x7f, 0xb4, 0x00, 0x72, 0x09, 0xb3, 0x8e, 0x59, 0x53, 0x37, 0x75, 0x47, 0x35, 0x3d, 0xc6, 0x62,
	0x88, 0x4e, 0x6a, 0x3a, 0xe5, 0xeb, 0x57, 0x08, 0x9e, 0x4e, 0x97, 0x24, 0xd7, 0xc3, 0xea, 0x22,
	0x7e, 0x76, 0xb0, 0x24, 0xe4, 0xa3, 0x53, 0x36, 0x21, 0xb7, 0xe0, 0xf7, 0x13, 0x48, 0xcc, 0x7f,
	0x18, 0x4b, 0xe6, 0xfc, 0x7c, 0x34, 0x9b, 0x6c, 0x7c, 0xfb, 0x01, 0x1e, 0xbe, 0x02, 0x00, 0x00,
	0xff, 0xff, 0xff, 0x4f, 0xce, 0x63, 0x37, 0x02, 0x00, 0x00,
}
