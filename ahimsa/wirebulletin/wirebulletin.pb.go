// Code generated by protoc-gen-go.
// source: wirebulletin.proto
// DO NOT EDIT!

/*
Package wirebulletin is a generated protocol buffer package.

It is generated from these files:
	wirebulletin.proto

It has these top-level messages:
	WireBulletin
*/
package wirebulletin

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type WireBulletin struct {
	Board            *string `protobuf:"bytes,1,opt,name=board" json:"board,omitempty"`
	Message          *string `protobuf:"bytes,2,req,name=message" json:"message,omitempty"`
	Timestamp        *int64  `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *WireBulletin) Reset()         { *m = WireBulletin{} }
func (m *WireBulletin) String() string { return proto.CompactTextString(m) }
func (*WireBulletin) ProtoMessage()    {}

func (m *WireBulletin) GetBoard() string {
	if m != nil && m.Board != nil {
		return *m.Board
	}
	return ""
}

func (m *WireBulletin) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

func (m *WireBulletin) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func init() {
}
