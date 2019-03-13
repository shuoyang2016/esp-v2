// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/envoy/http/path_matcher/config.proto

package path_matcher

import (
	common "cloudesf.googlesource.com/gcpproxy/src/go/proto/api/envoy/http/common"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/lyft/protoc-gen-validate/validate"
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

type PathMatcherRule struct {
	Pattern               *common.Pattern `protobuf:"bytes,1,opt,name=pattern,proto3" json:"pattern,omitempty"`
	Operation             string          `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	ExtractPathParameters bool            `protobuf:"varint,3,opt,name=extract_path_parameters,json=extractPathParameters,proto3" json:"extract_path_parameters,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}        `json:"-"`
	XXX_unrecognized      []byte          `json:"-"`
	XXX_sizecache         int32           `json:"-"`
}

func (m *PathMatcherRule) Reset()         { *m = PathMatcherRule{} }
func (m *PathMatcherRule) String() string { return proto.CompactTextString(m) }
func (*PathMatcherRule) ProtoMessage()    {}
func (*PathMatcherRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_bebe96b55b7e4dec, []int{0}
}

func (m *PathMatcherRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PathMatcherRule.Unmarshal(m, b)
}
func (m *PathMatcherRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PathMatcherRule.Marshal(b, m, deterministic)
}
func (m *PathMatcherRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PathMatcherRule.Merge(m, src)
}
func (m *PathMatcherRule) XXX_Size() int {
	return xxx_messageInfo_PathMatcherRule.Size(m)
}
func (m *PathMatcherRule) XXX_DiscardUnknown() {
	xxx_messageInfo_PathMatcherRule.DiscardUnknown(m)
}

var xxx_messageInfo_PathMatcherRule proto.InternalMessageInfo

func (m *PathMatcherRule) GetPattern() *common.Pattern {
	if m != nil {
		return m.Pattern
	}
	return nil
}

func (m *PathMatcherRule) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *PathMatcherRule) GetExtractPathParameters() bool {
	if m != nil {
		return m.ExtractPathParameters
	}
	return false
}

type SegmentName struct {
	SnakeName            string   `protobuf:"bytes,1,opt,name=snake_name,json=snakeName,proto3" json:"snake_name,omitempty"`
	JsonName             string   `protobuf:"bytes,2,opt,name=json_name,json=jsonName,proto3" json:"json_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SegmentName) Reset()         { *m = SegmentName{} }
func (m *SegmentName) String() string { return proto.CompactTextString(m) }
func (*SegmentName) ProtoMessage()    {}
func (*SegmentName) Descriptor() ([]byte, []int) {
	return fileDescriptor_bebe96b55b7e4dec, []int{1}
}

func (m *SegmentName) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SegmentName.Unmarshal(m, b)
}
func (m *SegmentName) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SegmentName.Marshal(b, m, deterministic)
}
func (m *SegmentName) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SegmentName.Merge(m, src)
}
func (m *SegmentName) XXX_Size() int {
	return xxx_messageInfo_SegmentName.Size(m)
}
func (m *SegmentName) XXX_DiscardUnknown() {
	xxx_messageInfo_SegmentName.DiscardUnknown(m)
}

var xxx_messageInfo_SegmentName proto.InternalMessageInfo

func (m *SegmentName) GetSnakeName() string {
	if m != nil {
		return m.SnakeName
	}
	return ""
}

func (m *SegmentName) GetJsonName() string {
	if m != nil {
		return m.JsonName
	}
	return ""
}

type FilterConfig struct {
	Rules                []*PathMatcherRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
	SegmentNames         []*SegmentName     `protobuf:"bytes,2,rep,name=segment_names,json=segmentNames,proto3" json:"segment_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *FilterConfig) Reset()         { *m = FilterConfig{} }
func (m *FilterConfig) String() string { return proto.CompactTextString(m) }
func (*FilterConfig) ProtoMessage()    {}
func (*FilterConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_bebe96b55b7e4dec, []int{2}
}

func (m *FilterConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilterConfig.Unmarshal(m, b)
}
func (m *FilterConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilterConfig.Marshal(b, m, deterministic)
}
func (m *FilterConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilterConfig.Merge(m, src)
}
func (m *FilterConfig) XXX_Size() int {
	return xxx_messageInfo_FilterConfig.Size(m)
}
func (m *FilterConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_FilterConfig.DiscardUnknown(m)
}

var xxx_messageInfo_FilterConfig proto.InternalMessageInfo

func (m *FilterConfig) GetRules() []*PathMatcherRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

func (m *FilterConfig) GetSegmentNames() []*SegmentName {
	if m != nil {
		return m.SegmentNames
	}
	return nil
}

func init() {
	proto.RegisterType((*PathMatcherRule)(nil), "google.api.envoy.http.path_matcher.PathMatcherRule")
	proto.RegisterType((*SegmentName)(nil), "google.api.envoy.http.path_matcher.SegmentName")
	proto.RegisterType((*FilterConfig)(nil), "google.api.envoy.http.path_matcher.FilterConfig")
}

func init() {
	proto.RegisterFile("api/envoy/http/path_matcher/config.proto", fileDescriptor_bebe96b55b7e4dec)
}

var fileDescriptor_bebe96b55b7e4dec = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0xd9, 0xd6, 0x3f, 0xcd, 0xb6, 0x22, 0x04, 0xa4, 0xa1, 0x22, 0x84, 0x88, 0x98, 0xd3,
	0x06, 0x5a, 0xf0, 0x2a, 0x28, 0x08, 0x3d, 0x28, 0x65, 0xf5, 0x5e, 0xd6, 0x38, 0x26, 0xd1, 0x64,
	0x37, 0xec, 0x4e, 0x8b, 0x3e, 0x95, 0x07, 0x6f, 0x9e, 0x7c, 0x1d, 0xdf, 0x42, 0xb2, 0xdb, 0xda,
	0x52, 0x04, 0xbd, 0x25, 0xf3, 0xed, 0xf7, 0x9b, 0x6f, 0x66, 0x68, 0x2c, 0xea, 0x22, 0x01, 0x39,
	0x57, 0xaf, 0x49, 0x8e, 0x58, 0x27, 0xb5, 0xc0, 0x7c, 0x5a, 0x09, 0x4c, 0x73, 0xd0, 0x49, 0xaa,
	0xe4, 0x63, 0x91, 0xb1, 0x5a, 0x2b, 0x54, 0x7e, 0x94, 0x29, 0x95, 0x95, 0xc0, 0x44, 0x5d, 0x30,
	0x6b, 0x60, 0x8d, 0x81, 0xad, 0x1b, 0x06, 0xc7, 0x1b, 0xb4, 0x54, 0x55, 0x95, 0x92, 0x0d, 0x14,
	0x41, 0x4b, 0x07, 0x1a, 0xf4, 0xe7, 0xa2, 0x2c, 0x1e, 0x04, 0x42, 0xb2, 0xfc, 0x70, 0x42, 0xf4,
	0x4e, 0xe8, 0xfe, 0x44, 0x60, 0x7e, 0xed, 0x68, 0x7c, 0x56, 0x82, 0x7f, 0x4e, 0x77, 0x17, 0xee,
	0x80, 0x84, 0x24, 0xee, 0x0e, 0x4f, 0xd8, 0xef, 0x39, 0x5c, 0x2b, 0x36, 0x71, 0x8f, 0xf9, 0xd2,
	0xe5, 0x9f, 0x52, 0x4f, 0xd5, 0xa0, 0x05, 0x16, 0x4a, 0x06, 0xad, 0x90, 0xc4, 0xde, 0x85, 0xf7,
	0xf1, 0xf5, 0xd9, 0xde, 0xd2, 0xad, 0x90, 0xf0, 0x95, 0xe6, 0x9f, 0xd1, 0x3e, 0xbc, 0xa0, 0x16,
	0x29, 0x4e, 0xed, 0x4c, 0xb5, 0xd0, 0xa2, 0x02, 0x04, 0x6d, 0x82, 0x76, 0x48, 0xe2, 0x0e, 0x3f,
	0x58, 0xc8, 0x4d, 0xc4, 0xc9, 0x8f, 0x18, 0x8d, 0x69, 0xf7, 0x16, 0xb2, 0x0a, 0x24, 0xde, 0x88,
	0x0a, 0xfc, 0x23, 0x4a, 0x8d, 0x14, 0xcf, 0x30, 0x95, 0xa2, 0x02, 0x9b, 0xd9, 0xe3, 0x9e, 0xad,
	0x58, 0xf9, 0x90, 0x7a, 0x4f, 0x46, 0x49, 0xa7, 0xda, 0x38, 0xbc, 0xd3, 0x14, 0x1a, 0x31, 0x7a,
	0x23, 0xb4, 0x77, 0x55, 0x94, 0x08, 0xfa, 0xd2, 0x6e, 0xde, 0x1f, 0xd3, 0x6d, 0x3d, 0x2b, 0xc1,
	0x04, 0x24, 0x6c, 0xc7, 0xdd, 0xe1, 0x88, 0xfd, 0x7d, 0x03, 0xb6, 0xb1, 0x41, 0xee, 0x08, 0xfe,
	0x1d, 0xdd, 0x33, 0x2e, 0xa6, 0xed, 0x6d, 0x82, 0x96, 0x45, 0x26, 0xff, 0x41, 0xae, 0xcd, 0xc7,
	0x7b, 0x66, 0xf5, 0x63, 0xee, 0x77, 0xec, 0xe5, 0x46, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x12,
	0xdb, 0x92, 0xbe, 0x47, 0x02, 0x00, 0x00,
}
