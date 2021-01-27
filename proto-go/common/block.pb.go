// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: blockproto/block.proto

package common

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Envelope wraps a Payload with a signature so that the message may be authenticated
// <REQUEST,o,t,c>
// Request
type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A marshaled Payload
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// A signature by the creator specified in the Payload header
	Signature  []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	Timestamp  int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ClientAddr string `protobuf:"bytes,4,opt,name=clientAddr,proto3" json:"clientAddr,omitempty"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Envelope) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Envelope) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Envelope) GetClientAddr() string {
	if x != nil {
		return x.ClientAddr
	}
	return ""
}

// This is finalized block structure to be shared among the orderer and peer
// Note that the BlockHeader chains to the previous BlockHeader, and the BlockData hash is embedded
// in the BlockHeader.  This makes it natural and obvious that the Data is included in the hash, but
// the Metadata is not.
//<<PRE-PREPARE,v,n,d>,m>
//PrePrepare
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *BlockHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Data   *BlockData   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{1}
}

func (x *Block) GetHeader() *BlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Block) GetData() *BlockData {
	if x != nil {
		return x.Data
	}
	return nil
}

// BlockHeader is the element of the block which forms the block chain
// The block header is hashed using the configured chain hashing algorithm
// over the ASN.1 encoding of the BlockHeader
type BlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PreviousHash [][]byte `protobuf:"bytes,1,rep,name=previous_hash,json=previousHash,proto3" json:"previous_hash,omitempty"` // The hash of the previous block header
	DataHash     []byte   `protobuf:"bytes,2,opt,name=data_hash,json=dataHash,proto3" json:"data_hash,omitempty"`             // The hash of the BlockData, by MerkleTree
}

func (x *BlockHeader) Reset() {
	*x = BlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHeader) ProtoMessage() {}

func (x *BlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHeader.ProtoReflect.Descriptor instead.
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{2}
}

func (x *BlockHeader) GetPreviousHash() [][]byte {
	if x != nil {
		return x.PreviousHash
	}
	return nil
}

func (x *BlockHeader) GetDataHash() []byte {
	if x != nil {
		return x.DataHash
	}
	return nil
}

type BlockData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data [][]byte `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *BlockData) Reset() {
	*x = BlockData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockData) ProtoMessage() {}

func (x *BlockData) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockData.ProtoReflect.Descriptor instead.
func (*BlockData) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{3}
}

func (x *BlockData) GetData() [][]byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type PrePrepareMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Block *Block `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	//digest is the blockheader hash
	Digest string `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	Sign   []byte `protobuf:"bytes,3,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (x *PrePrepareMsg) Reset() {
	*x = PrePrepareMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrePrepareMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrePrepareMsg) ProtoMessage() {}

func (x *PrePrepareMsg) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrePrepareMsg.ProtoReflect.Descriptor instead.
func (*PrePrepareMsg) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{4}
}

func (x *PrePrepareMsg) GetBlock() *Block {
	if x != nil {
		return x.Block
	}
	return nil
}

func (x *PrePrepareMsg) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *PrePrepareMsg) GetSign() []byte {
	if x != nil {
		return x.Sign
	}
	return nil
}

type PrepareMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digest string `protobuf:"bytes,1,opt,name=digest,proto3" json:"digest,omitempty"`
	Sign   []byte `protobuf:"bytes,2,opt,name=sign,proto3" json:"sign,omitempty"`
	NodeID string `protobuf:"bytes,3,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *PrepareMsg) Reset() {
	*x = PrepareMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareMsg) ProtoMessage() {}

func (x *PrepareMsg) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareMsg.ProtoReflect.Descriptor instead.
func (*PrepareMsg) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{5}
}

func (x *PrepareMsg) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *PrepareMsg) GetSign() []byte {
	if x != nil {
		return x.Sign
	}
	return nil
}

func (x *PrepareMsg) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type CommitMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digest string `protobuf:"bytes,1,opt,name=digest,proto3" json:"digest,omitempty"`
	Sign   []byte `protobuf:"bytes,2,opt,name=sign,proto3" json:"sign,omitempty"`
	NodeID string `protobuf:"bytes,3,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
}

func (x *CommitMsg) Reset() {
	*x = CommitMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitMsg) ProtoMessage() {}

func (x *CommitMsg) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitMsg.ProtoReflect.Descriptor instead.
func (*CommitMsg) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{6}
}

func (x *CommitMsg) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *CommitMsg) GetSign() []byte {
	if x != nil {
		return x.Sign
	}
	return nil
}

func (x *CommitMsg) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResCode bool `protobuf:"varint,1,opt,name=resCode,proto3" json:"resCode,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockproto_block_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_blockproto_block_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_blockproto_block_proto_rawDescGZIP(), []int{7}
}

func (x *Response) GetResCode() bool {
	if x != nil {
		return x.ResCode
	}
	return false
}

var File_blockproto_block_proto protoreflect.FileDescriptor

var file_blockproto_block_proto_rawDesc = []byte{
	0x0a, 0x16, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x22, 0x80, 0x01, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x22, 0x5b, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x2b, 0x0a, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x4f, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x23, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73,
	0x68, 0x22, 0x1f, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x60, 0x0a, 0x0d, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x4d, 0x73, 0x67, 0x12, 0x23, 0x0a, 0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x52, 0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65,
	0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x73, 0x69, 0x67, 0x6e, 0x22, 0x50, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4d,
	0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69,
	0x67, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x4f, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x4d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x69, 0x67, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x32, 0x43, 0x0a,
	0x0c, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x73, 0x12, 0x33, 0x0a,
	0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x1a, 0x10, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01,
	0x30, 0x01, 0x32, 0xb5, 0x01, 0x0a, 0x04, 0x50, 0x62, 0x66, 0x74, 0x12, 0x3d, 0x0a, 0x10, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12,
	0x15, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70,
	0x61, 0x72, 0x65, 0x4d, 0x73, 0x67, 0x1a, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0d, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x12, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4d, 0x73, 0x67, 0x1a,
	0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x0c, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x12, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2d, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blockproto_block_proto_rawDescOnce sync.Once
	file_blockproto_block_proto_rawDescData = file_blockproto_block_proto_rawDesc
)

func file_blockproto_block_proto_rawDescGZIP() []byte {
	file_blockproto_block_proto_rawDescOnce.Do(func() {
		file_blockproto_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_blockproto_block_proto_rawDescData)
	})
	return file_blockproto_block_proto_rawDescData
}

var file_blockproto_block_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_blockproto_block_proto_goTypes = []interface{}{
	(*Envelope)(nil),      // 0: common.Envelope
	(*Block)(nil),         // 1: common.Block
	(*BlockHeader)(nil),   // 2: common.BlockHeader
	(*BlockData)(nil),     // 3: common.BlockData
	(*PrePrepareMsg)(nil), // 4: common.PrePrepareMsg
	(*PrepareMsg)(nil),    // 5: common.PrepareMsg
	(*CommitMsg)(nil),     // 6: common.CommitMsg
	(*Response)(nil),      // 7: common.Response
}
var file_blockproto_block_proto_depIdxs = []int32{
	2, // 0: common.Block.header:type_name -> common.BlockHeader
	3, // 1: common.Block.data:type_name -> common.BlockData
	1, // 2: common.PrePrepareMsg.block:type_name -> common.Block
	0, // 3: common.SendEnvelops.Request:input_type -> common.Envelope
	4, // 4: common.Pbft.HandlePrePrepare:input_type -> common.PrePrepareMsg
	5, // 5: common.Pbft.HandlePrepare:input_type -> common.PrepareMsg
	6, // 6: common.Pbft.HandleCommit:input_type -> common.CommitMsg
	7, // 7: common.SendEnvelops.Request:output_type -> common.Response
	7, // 8: common.Pbft.HandlePrePrepare:output_type -> common.Response
	7, // 9: common.Pbft.HandlePrepare:output_type -> common.Response
	7, // 10: common.Pbft.HandleCommit:output_type -> common.Response
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_blockproto_block_proto_init() }
func file_blockproto_block_proto_init() {
	if File_blockproto_block_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blockproto_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Envelope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHeader); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrePrepareMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockproto_block_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_blockproto_block_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_blockproto_block_proto_goTypes,
		DependencyIndexes: file_blockproto_block_proto_depIdxs,
		MessageInfos:      file_blockproto_block_proto_msgTypes,
	}.Build()
	File_blockproto_block_proto = out.File
	file_blockproto_block_proto_rawDesc = nil
	file_blockproto_block_proto_goTypes = nil
	file_blockproto_block_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SendEnvelopsClient is the client API for SendEnvelops service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SendEnvelopsClient interface {
	Request(ctx context.Context, opts ...grpc.CallOption) (SendEnvelops_RequestClient, error)
}

type sendEnvelopsClient struct {
	cc grpc.ClientConnInterface
}

func NewSendEnvelopsClient(cc grpc.ClientConnInterface) SendEnvelopsClient {
	return &sendEnvelopsClient{cc}
}

func (c *sendEnvelopsClient) Request(ctx context.Context, opts ...grpc.CallOption) (SendEnvelops_RequestClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SendEnvelops_serviceDesc.Streams[0], "/common.SendEnvelops/Request", opts...)
	if err != nil {
		return nil, err
	}
	x := &sendEnvelopsRequestClient{stream}
	return x, nil
}

type SendEnvelops_RequestClient interface {
	Send(*Envelope) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type sendEnvelopsRequestClient struct {
	grpc.ClientStream
}

func (x *sendEnvelopsRequestClient) Send(m *Envelope) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sendEnvelopsRequestClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SendEnvelopsServer is the server API for SendEnvelops service.
type SendEnvelopsServer interface {
	Request(SendEnvelops_RequestServer) error
}

// UnimplementedSendEnvelopsServer can be embedded to have forward compatible implementations.
type UnimplementedSendEnvelopsServer struct {
}

func (*UnimplementedSendEnvelopsServer) Request(SendEnvelops_RequestServer) error {
	return status.Errorf(codes.Unimplemented, "method Request not implemented")
}

func RegisterSendEnvelopsServer(s *grpc.Server, srv SendEnvelopsServer) {
	s.RegisterService(&_SendEnvelops_serviceDesc, srv)
}

func _SendEnvelops_Request_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SendEnvelopsServer).Request(&sendEnvelopsRequestServer{stream})
}

type SendEnvelops_RequestServer interface {
	Send(*Response) error
	Recv() (*Envelope, error)
	grpc.ServerStream
}

type sendEnvelopsRequestServer struct {
	grpc.ServerStream
}

func (x *sendEnvelopsRequestServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sendEnvelopsRequestServer) Recv() (*Envelope, error) {
	m := new(Envelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _SendEnvelops_serviceDesc = grpc.ServiceDesc{
	ServiceName: "common.SendEnvelops",
	HandlerType: (*SendEnvelopsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Request",
			Handler:       _SendEnvelops_Request_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "blockproto/block.proto",
}

// PbftClient is the client API for Pbft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PbftClient interface {
	HandlePrePrepare(ctx context.Context, in *PrePrepareMsg, opts ...grpc.CallOption) (*Response, error)
	HandlePrepare(ctx context.Context, in *PrepareMsg, opts ...grpc.CallOption) (*Response, error)
	HandleCommit(ctx context.Context, in *CommitMsg, opts ...grpc.CallOption) (*Response, error)
}

type pbftClient struct {
	cc grpc.ClientConnInterface
}

func NewPbftClient(cc grpc.ClientConnInterface) PbftClient {
	return &pbftClient{cc}
}

func (c *pbftClient) HandlePrePrepare(ctx context.Context, in *PrePrepareMsg, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/common.Pbft/HandlePrePrepare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pbftClient) HandlePrepare(ctx context.Context, in *PrepareMsg, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/common.Pbft/HandlePrepare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pbftClient) HandleCommit(ctx context.Context, in *CommitMsg, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/common.Pbft/HandleCommit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PbftServer is the server API for Pbft service.
type PbftServer interface {
	HandlePrePrepare(context.Context, *PrePrepareMsg) (*Response, error)
	HandlePrepare(context.Context, *PrepareMsg) (*Response, error)
	HandleCommit(context.Context, *CommitMsg) (*Response, error)
}

// UnimplementedPbftServer can be embedded to have forward compatible implementations.
type UnimplementedPbftServer struct {
}

func (*UnimplementedPbftServer) HandlePrePrepare(context.Context, *PrePrepareMsg) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandlePrePrepare not implemented")
}
func (*UnimplementedPbftServer) HandlePrepare(context.Context, *PrepareMsg) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandlePrepare not implemented")
}
func (*UnimplementedPbftServer) HandleCommit(context.Context, *CommitMsg) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleCommit not implemented")
}

func RegisterPbftServer(s *grpc.Server, srv PbftServer) {
	s.RegisterService(&_Pbft_serviceDesc, srv)
}

func _Pbft_HandlePrePrepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrePrepareMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PbftServer).HandlePrePrepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Pbft/HandlePrePrepare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PbftServer).HandlePrePrepare(ctx, req.(*PrePrepareMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pbft_HandlePrepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PbftServer).HandlePrepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Pbft/HandlePrepare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PbftServer).HandlePrepare(ctx, req.(*PrepareMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pbft_HandleCommit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PbftServer).HandleCommit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Pbft/HandleCommit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PbftServer).HandleCommit(ctx, req.(*CommitMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pbft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "common.Pbft",
	HandlerType: (*PbftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandlePrePrepare",
			Handler:    _Pbft_HandlePrePrepare_Handler,
		},
		{
			MethodName: "HandlePrepare",
			Handler:    _Pbft_HandlePrepare_Handler,
		},
		{
			MethodName: "HandleCommit",
			Handler:    _Pbft_HandleCommit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blockproto/block.proto",
}
