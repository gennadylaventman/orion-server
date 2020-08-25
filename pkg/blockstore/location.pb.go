// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/blockstore/location.proto

package blockstore

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type BlockLocation struct {
	FileChunkNum         uint64   `protobuf:"varint,1,opt,name=file_chunk_num,json=fileChunkNum,proto3" json:"file_chunk_num,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockLocation) Reset()         { *m = BlockLocation{} }
func (m *BlockLocation) String() string { return proto.CompactTextString(m) }
func (*BlockLocation) ProtoMessage()    {}
func (*BlockLocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fea647f64c5d3aa, []int{0}
}

func (m *BlockLocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockLocation.Unmarshal(m, b)
}
func (m *BlockLocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockLocation.Marshal(b, m, deterministic)
}
func (m *BlockLocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockLocation.Merge(m, src)
}
func (m *BlockLocation) XXX_Size() int {
	return xxx_messageInfo_BlockLocation.Size(m)
}
func (m *BlockLocation) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockLocation.DiscardUnknown(m)
}

var xxx_messageInfo_BlockLocation proto.InternalMessageInfo

func (m *BlockLocation) GetFileChunkNum() uint64 {
	if m != nil {
		return m.FileChunkNum
	}
	return 0
}

func (m *BlockLocation) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func init() {
	proto.RegisterType((*BlockLocation)(nil), "blockstore.BlockLocation")
}

func init() { proto.RegisterFile("pkg/blockstore/location.proto", fileDescriptor_0fea647f64c5d3aa) }

var fileDescriptor_0fea647f64c5d3aa = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0xc8, 0x4e, 0xd7,
	0x4f, 0xca, 0xc9, 0x4f, 0xce, 0x2e, 0x2e, 0xc9, 0x2f, 0x4a, 0xd5, 0xcf, 0xc9, 0x4f, 0x4e, 0x2c,
	0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x48, 0x29, 0xf9, 0x72,
	0xf1, 0x3a, 0x81, 0x78, 0x3e, 0x50, 0x25, 0x42, 0x2a, 0x5c, 0x7c, 0x69, 0x99, 0x39, 0xa9, 0xf1,
	0xc9, 0x19, 0xa5, 0x79, 0xd9, 0xf1, 0x79, 0xa5, 0xb9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41,
	0x3c, 0x20, 0x51, 0x67, 0x90, 0xa0, 0x5f, 0x69, 0xae, 0x90, 0x18, 0x17, 0x5b, 0x7e, 0x5a, 0x5a,
	0x71, 0x6a, 0x89, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x73, 0x10, 0x94, 0xe7, 0x64, 0x1c, 0x65, 0x98,
	0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x99, 0x94, 0xab, 0x97, 0x9c, 0x9f, 0x0b, 0x71, 0x46,
	0x72, 0x46, 0x62, 0x66, 0x5e, 0x4a, 0x92, 0x7e, 0x71, 0x6a, 0x51, 0x59, 0x6a, 0x91, 0x3e, 0xaa,
	0xf3, 0x92, 0xd8, 0xc0, 0xce, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x84, 0x00, 0xea, 0x3b,
	0xb7, 0x00, 0x00, 0x00,
}