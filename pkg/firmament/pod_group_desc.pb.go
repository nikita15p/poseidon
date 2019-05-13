/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pod_group_desc.proto

package firmament

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

type PodGroupDescriptor struct {
	Uid                  uint64   `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	MinMember            int32    `protobuf:"varint,2,opt,name=minMember,proto3" json:"minMember,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Queue                string   `protobuf:"bytes,4,opt,name=queue,proto3" json:"queue,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PodGroupDescriptor) Reset()         { *m = PodGroupDescriptor{} }
func (m *PodGroupDescriptor) String() string { return proto.CompactTextString(m) }
func (*PodGroupDescriptor) ProtoMessage()    {}
func (*PodGroupDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_337a310aa4536dd3, []int{0}
}

func (m *PodGroupDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PodGroupDescriptor.Unmarshal(m, b)
}
func (m *PodGroupDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PodGroupDescriptor.Marshal(b, m, deterministic)
}
func (m *PodGroupDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PodGroupDescriptor.Merge(m, src)
}
func (m *PodGroupDescriptor) XXX_Size() int {
	return xxx_messageInfo_PodGroupDescriptor.Size(m)
}
func (m *PodGroupDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_PodGroupDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_PodGroupDescriptor proto.InternalMessageInfo

func (m *PodGroupDescriptor) GetUid() uint64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *PodGroupDescriptor) GetMinMember() int32 {
	if m != nil {
		return m.MinMember
	}
	return 0
}

func (m *PodGroupDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PodGroupDescriptor) GetQueue() string {
	if m != nil {
		return m.Queue
	}
	return ""
}

func init() {
	proto.RegisterType((*PodGroupDescriptor)(nil), "firmament.PodGroupDescriptor")
}

func init() { proto.RegisterFile("pod_group_desc.proto", fileDescriptor_337a310aa4536dd3) }

var fileDescriptor_337a310aa4536dd3 = []byte{
	// 151 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xc8, 0x4f, 0x89,
	0x4f, 0x2f, 0xca, 0x2f, 0x2d, 0x88, 0x4f, 0x49, 0x2d, 0x4e, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x4c, 0xcb, 0x2c, 0xca, 0x4d, 0xcc, 0x4d, 0xcd, 0x2b, 0x51, 0xca, 0xe3, 0x12, 0x0a,
	0xc8, 0x4f, 0x71, 0x07, 0xa9, 0x70, 0x49, 0x2d, 0x4e, 0x2e, 0xca, 0x2c, 0x28, 0xc9, 0x2f, 0x12,
	0x12, 0xe0, 0x62, 0x2e, 0xcd, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x09, 0x02, 0x31, 0x85,
	0x64, 0xb8, 0x38, 0x73, 0x33, 0xf3, 0x7c, 0x53, 0x73, 0x93, 0x52, 0x8b, 0x24, 0x98, 0x14, 0x18,
	0x35, 0x58, 0x83, 0x10, 0x02, 0x42, 0x42, 0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0xcc, 0x0a,
	0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x90, 0x08, 0x17, 0x6b, 0x61, 0x69, 0x6a, 0x69, 0xaa, 0x04,
	0x0b, 0x58, 0x10, 0xc2, 0x49, 0x62, 0x03, 0xbb, 0xc0, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x1a,
	0xed, 0x3a, 0xef, 0x99, 0x00, 0x00, 0x00,
}
