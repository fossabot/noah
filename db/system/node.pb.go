// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node.proto

package system

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NNode struct {
	NodeId               uint64   `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Port                 int32    `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	Database             string   `protobuf:"bytes,4,opt,name=database,proto3" json:"database,omitempty"`
	User                 string   `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
	Password             string   `protobuf:"bytes,6,opt,name=password,proto3" json:"password,omitempty"`
	ReplicaOf            uint64   `protobuf:"varint,7,opt,name=replica_of,json=replicaOf,proto3" json:"replica_of,omitempty"`
	Region               string   `protobuf:"bytes,8,opt,name=region,proto3" json:"region,omitempty"`
	Zone                 string   `protobuf:"bytes,9,opt,name=zone,proto3" json:"zone,omitempty"`
	IsAlive              bool     `protobuf:"varint,10,opt,name=is_alive,json=isAlive,proto3" json:"is_alive,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NNode) Reset()         { *m = NNode{} }
func (m *NNode) String() string { return proto.CompactTextString(m) }
func (*NNode) ProtoMessage()    {}
func (*NNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

func (m *NNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NNode.Unmarshal(m, b)
}
func (m *NNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NNode.Marshal(b, m, deterministic)
}
func (m *NNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NNode.Merge(m, src)
}
func (m *NNode) XXX_Size() int {
	return xxx_messageInfo_NNode.Size(m)
}
func (m *NNode) XXX_DiscardUnknown() {
	xxx_messageInfo_NNode.DiscardUnknown(m)
}

var xxx_messageInfo_NNode proto.InternalMessageInfo

func (m *NNode) GetNodeId() uint64 {
	if m != nil {
		return m.NodeId
	}
	return 0
}

func (m *NNode) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *NNode) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *NNode) GetDatabase() string {
	if m != nil {
		return m.Database
	}
	return ""
}

func (m *NNode) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *NNode) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *NNode) GetReplicaOf() uint64 {
	if m != nil {
		return m.ReplicaOf
	}
	return 0
}

func (m *NNode) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *NNode) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

func (m *NNode) GetIsAlive() bool {
	if m != nil {
		return m.IsAlive
	}
	return false
}

func init() {
	proto.RegisterType((*NNode)(nil), "system.NNode")
}

func init() { proto.RegisterFile("node.proto", fileDescriptor_0c843d59d2d938e7) }

var fileDescriptor_0c843d59d2d938e7 = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x90, 0x3d, 0x4e, 0xc4, 0x30,
	0x10, 0x85, 0xe5, 0x25, 0x71, 0x92, 0x29, 0xa7, 0x80, 0x01, 0x09, 0x29, 0xa2, 0x4a, 0x45, 0xc3,
	0x09, 0x28, 0x69, 0x16, 0x29, 0x17, 0x88, 0xbc, 0xcc, 0x2c, 0xb2, 0xb4, 0x64, 0x22, 0x4f, 0x00,
	0x41, 0xc5, 0xd1, 0x91, 0xbd, 0x3f, 0xdd, 0xfb, 0x9e, 0x3f, 0xeb, 0x59, 0x06, 0x98, 0x95, 0xe5,
	0x71, 0x49, 0xba, 0x2a, 0x7a, 0xfb, 0xb1, 0x55, 0x3e, 0x1e, 0xfe, 0x36, 0x50, 0x6f, 0xb7, 0xca,
	0x82, 0x37, 0xd0, 0xe4, 0xf3, 0x29, 0x32, 0xb9, 0xde, 0x0d, 0xd5, 0xe8, 0x33, 0xbe, 0x30, 0x12,
	0x34, 0x81, 0x39, 0x89, 0x19, 0x6d, 0x7a, 0x37, 0x74, 0xe3, 0x19, 0x11, 0xa1, 0x5a, 0x34, 0xad,
	0x74, 0xd5, 0xbb, 0xa1, 0x1e, 0x4b, 0xc6, 0x3b, 0x68, 0x39, 0xac, 0x61, 0x17, 0x4c, 0xa8, 0x2a,
	0xfa, 0x85, 0xb3, 0xff, 0x69, 0x92, 0xa8, 0x2e, 0x7d, 0xc9, 0xd9, 0x5f, 0x82, 0xd9, 0xb7, 0x26,
	0x26, 0x7f, 0xf4, 0xcf, 0x8c, 0xf7, 0x00, 0x49, 0x96, 0x43, 0x7c, 0x0b, 0x93, 0xee, 0xa9, 0x29,
	0xaf, 0xea, 0x4e, 0xcd, 0xeb, 0x1e, 0xaf, 0xc1, 0x27, 0x79, 0x8f, 0x3a, 0x53, 0x5b, 0x2e, 0x9e,
	0x28, 0xcf, 0xfc, 0xea, 0x2c, 0xd4, 0x1d, 0x67, 0x72, 0xc6, 0x5b, 0x68, 0xa3, 0x4d, 0xe1, 0x10,
	0xbf, 0x84, 0xa0, 0x77, 0x43, 0x3b, 0x36, 0xd1, 0x9e, 0x33, 0xee, 0x7c, 0xf9, 0x91, 0xa7, 0xff,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x75, 0xf9, 0xb2, 0x0d, 0x1f, 0x01, 0x00, 0x00,
}
