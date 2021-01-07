// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

package test_proto

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

type Test struct {
	BusinessId           *uint64  `protobuf:"varint,1,req,name=business_id,json=businessId" json:"business_id,omitempty"`
	ModuleId             *uint64  `protobuf:"varint,2,req,name=module_id,json=moduleId" json:"module_id,omitempty"`
	Picture              *string  `protobuf:"bytes,3,req,name=picture" json:"picture,omitempty"`
	RequestId            *string  `protobuf:"bytes,4,req,name=request_id,json=requestId" json:"request_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (m *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(m, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

func (m *Test) GetBusinessId() uint64 {
	if m != nil && m.BusinessId != nil {
		return *m.BusinessId
	}
	return 0
}

func (m *Test) GetModuleId() uint64 {
	if m != nil && m.ModuleId != nil {
		return *m.ModuleId
	}
	return 0
}

func (m *Test) GetPicture() string {
	if m != nil && m.Picture != nil {
		return *m.Picture
	}
	return ""
}

func (m *Test) GetRequestId() string {
	if m != nil && m.RequestId != nil {
		return *m.RequestId
	}
	return ""
}

func init() {
	proto.RegisterType((*Test)(nil), "test_proto.Test")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x02, 0xb3, 0xe3, 0xc1, 0x6c, 0xa5, 0x5a, 0x2e, 0x96,
	0x90, 0xd4, 0xe2, 0x12, 0x21, 0x79, 0x2e, 0xee, 0xa4, 0xd2, 0xe2, 0xcc, 0xbc, 0xd4, 0xe2, 0xe2,
	0xf8, 0xcc, 0x14, 0x09, 0x46, 0x05, 0x26, 0x0d, 0x96, 0x20, 0x2e, 0x98, 0x90, 0x67, 0x8a, 0x90,
	0x34, 0x17, 0x67, 0x6e, 0x7e, 0x4a, 0x69, 0x4e, 0x2a, 0x48, 0x9a, 0x09, 0x2c, 0xcd, 0x01, 0x11,
	0xf0, 0x4c, 0x11, 0x92, 0xe0, 0x62, 0x2f, 0xc8, 0x4c, 0x2e, 0x29, 0x2d, 0x4a, 0x95, 0x60, 0x56,
	0x60, 0xd2, 0xe0, 0x0c, 0x82, 0x71, 0x85, 0x64, 0xb9, 0xb8, 0x8a, 0x52, 0x0b, 0x4b, 0x41, 0x16,
	0x66, 0xa6, 0x48, 0xb0, 0x80, 0x25, 0x39, 0xa1, 0x22, 0x9e, 0x29, 0x4e, 0x7c, 0x51, 0x3c, 0x7a,
	0xd6, 0x08, 0xe7, 0x00, 0x02, 0x00, 0x00, 0xff, 0xff, 0x1d, 0x8e, 0xfe, 0x7d, 0xa7, 0x00, 0x00,
	0x00,
}