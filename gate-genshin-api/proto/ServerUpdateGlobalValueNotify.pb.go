// Sorapointa - A server software re-implementation for a certain anime game, and avoid sorapointa.
// Copyright (C) 2022  Sorapointa Team
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.7.0
// source: ServerUpdateGlobalValueNotify.proto

package proto

import (
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

type ServerUpdateGlobalValueNotify_UpdateType int32

const (
	ServerUpdateGlobalValueNotify_UPDATE_TYPE_INVALUE ServerUpdateGlobalValueNotify_UpdateType = 0
	ServerUpdateGlobalValueNotify_UPDATE_TYPE_ADD     ServerUpdateGlobalValueNotify_UpdateType = 1
	ServerUpdateGlobalValueNotify_UPDATE_TYPE_SET     ServerUpdateGlobalValueNotify_UpdateType = 2
)

// Enum value maps for ServerUpdateGlobalValueNotify_UpdateType.
var (
	ServerUpdateGlobalValueNotify_UpdateType_name = map[int32]string{
		0: "UPDATE_TYPE_INVALUE",
		1: "UPDATE_TYPE_ADD",
		2: "UPDATE_TYPE_SET",
	}
	ServerUpdateGlobalValueNotify_UpdateType_value = map[string]int32{
		"UPDATE_TYPE_INVALUE": 0,
		"UPDATE_TYPE_ADD":     1,
		"UPDATE_TYPE_SET":     2,
	}
)

func (x ServerUpdateGlobalValueNotify_UpdateType) Enum() *ServerUpdateGlobalValueNotify_UpdateType {
	p := new(ServerUpdateGlobalValueNotify_UpdateType)
	*p = x
	return p
}

func (x ServerUpdateGlobalValueNotify_UpdateType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerUpdateGlobalValueNotify_UpdateType) Descriptor() protoreflect.EnumDescriptor {
	return file_ServerUpdateGlobalValueNotify_proto_enumTypes[0].Descriptor()
}

func (ServerUpdateGlobalValueNotify_UpdateType) Type() protoreflect.EnumType {
	return &file_ServerUpdateGlobalValueNotify_proto_enumTypes[0]
}

func (x ServerUpdateGlobalValueNotify_UpdateType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServerUpdateGlobalValueNotify_UpdateType.Descriptor instead.
func (ServerUpdateGlobalValueNotify_UpdateType) EnumDescriptor() ([]byte, []int) {
	return file_ServerUpdateGlobalValueNotify_proto_rawDescGZIP(), []int{0, 0}
}

// CmdId: 1148
// EnetChannelId: 0
// EnetIsReliable: true
type ServerUpdateGlobalValueNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityId   uint32                                   `protobuf:"varint,9,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
	UpdateType ServerUpdateGlobalValueNotify_UpdateType `protobuf:"varint,13,opt,name=update_type,json=updateType,proto3,enum=ServerUpdateGlobalValueNotify_UpdateType" json:"update_type,omitempty"`
	Delta      float32                                  `protobuf:"fixed32,3,opt,name=delta,proto3" json:"delta,omitempty"`
	KeyHash    uint32                                   `protobuf:"varint,10,opt,name=key_hash,json=keyHash,proto3" json:"key_hash,omitempty"`
	Value      float32                                  `protobuf:"fixed32,6,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ServerUpdateGlobalValueNotify) Reset() {
	*x = ServerUpdateGlobalValueNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ServerUpdateGlobalValueNotify_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerUpdateGlobalValueNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerUpdateGlobalValueNotify) ProtoMessage() {}

func (x *ServerUpdateGlobalValueNotify) ProtoReflect() protoreflect.Message {
	mi := &file_ServerUpdateGlobalValueNotify_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerUpdateGlobalValueNotify.ProtoReflect.Descriptor instead.
func (*ServerUpdateGlobalValueNotify) Descriptor() ([]byte, []int) {
	return file_ServerUpdateGlobalValueNotify_proto_rawDescGZIP(), []int{0}
}

func (x *ServerUpdateGlobalValueNotify) GetEntityId() uint32 {
	if x != nil {
		return x.EntityId
	}
	return 0
}

func (x *ServerUpdateGlobalValueNotify) GetUpdateType() ServerUpdateGlobalValueNotify_UpdateType {
	if x != nil {
		return x.UpdateType
	}
	return ServerUpdateGlobalValueNotify_UPDATE_TYPE_INVALUE
}

func (x *ServerUpdateGlobalValueNotify) GetDelta() float32 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *ServerUpdateGlobalValueNotify) GetKeyHash() uint32 {
	if x != nil {
		return x.KeyHash
	}
	return 0
}

func (x *ServerUpdateGlobalValueNotify) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_ServerUpdateGlobalValueNotify_proto protoreflect.FileDescriptor

var file_ServerUpdateGlobalValueNotify_proto_rawDesc = []byte{
	0x0a, 0x23, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x6c,
	0x6f, 0x62, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x02, 0x0a, 0x1d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x49, 0x64, 0x12, 0x4a, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x4f, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x00, 0x12, 0x13,
	0x0a, 0x0f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x44,
	0x44, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x53, 0x45, 0x54, 0x10, 0x02, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ServerUpdateGlobalValueNotify_proto_rawDescOnce sync.Once
	file_ServerUpdateGlobalValueNotify_proto_rawDescData = file_ServerUpdateGlobalValueNotify_proto_rawDesc
)

func file_ServerUpdateGlobalValueNotify_proto_rawDescGZIP() []byte {
	file_ServerUpdateGlobalValueNotify_proto_rawDescOnce.Do(func() {
		file_ServerUpdateGlobalValueNotify_proto_rawDescData = protoimpl.X.CompressGZIP(file_ServerUpdateGlobalValueNotify_proto_rawDescData)
	})
	return file_ServerUpdateGlobalValueNotify_proto_rawDescData
}

var file_ServerUpdateGlobalValueNotify_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ServerUpdateGlobalValueNotify_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ServerUpdateGlobalValueNotify_proto_goTypes = []interface{}{
	(ServerUpdateGlobalValueNotify_UpdateType)(0), // 0: ServerUpdateGlobalValueNotify.UpdateType
	(*ServerUpdateGlobalValueNotify)(nil),         // 1: ServerUpdateGlobalValueNotify
}
var file_ServerUpdateGlobalValueNotify_proto_depIdxs = []int32{
	0, // 0: ServerUpdateGlobalValueNotify.update_type:type_name -> ServerUpdateGlobalValueNotify.UpdateType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ServerUpdateGlobalValueNotify_proto_init() }
func file_ServerUpdateGlobalValueNotify_proto_init() {
	if File_ServerUpdateGlobalValueNotify_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ServerUpdateGlobalValueNotify_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerUpdateGlobalValueNotify); i {
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
			RawDescriptor: file_ServerUpdateGlobalValueNotify_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ServerUpdateGlobalValueNotify_proto_goTypes,
		DependencyIndexes: file_ServerUpdateGlobalValueNotify_proto_depIdxs,
		EnumInfos:         file_ServerUpdateGlobalValueNotify_proto_enumTypes,
		MessageInfos:      file_ServerUpdateGlobalValueNotify_proto_msgTypes,
	}.Build()
	File_ServerUpdateGlobalValueNotify_proto = out.File
	file_ServerUpdateGlobalValueNotify_proto_rawDesc = nil
	file_ServerUpdateGlobalValueNotify_proto_goTypes = nil
	file_ServerUpdateGlobalValueNotify_proto_depIdxs = nil
}