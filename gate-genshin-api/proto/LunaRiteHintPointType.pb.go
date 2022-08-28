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
// source: LunaRiteHintPointType.proto

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

type LunaRiteHintPointType int32

const (
	LunaRiteHintPointType_LUNA_RITE_HINT_POINT_TYPE_NONE  LunaRiteHintPointType = 0
	LunaRiteHintPointType_LUNA_RITE_HINT_POINT_TYPE_RUNE  LunaRiteHintPointType = 1
	LunaRiteHintPointType_LUNA_RITE_HINT_POINT_TYPE_CHEST LunaRiteHintPointType = 2
)

// Enum value maps for LunaRiteHintPointType.
var (
	LunaRiteHintPointType_name = map[int32]string{
		0: "LUNA_RITE_HINT_POINT_TYPE_NONE",
		1: "LUNA_RITE_HINT_POINT_TYPE_RUNE",
		2: "LUNA_RITE_HINT_POINT_TYPE_CHEST",
	}
	LunaRiteHintPointType_value = map[string]int32{
		"LUNA_RITE_HINT_POINT_TYPE_NONE":  0,
		"LUNA_RITE_HINT_POINT_TYPE_RUNE":  1,
		"LUNA_RITE_HINT_POINT_TYPE_CHEST": 2,
	}
)

func (x LunaRiteHintPointType) Enum() *LunaRiteHintPointType {
	p := new(LunaRiteHintPointType)
	*p = x
	return p
}

func (x LunaRiteHintPointType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LunaRiteHintPointType) Descriptor() protoreflect.EnumDescriptor {
	return file_LunaRiteHintPointType_proto_enumTypes[0].Descriptor()
}

func (LunaRiteHintPointType) Type() protoreflect.EnumType {
	return &file_LunaRiteHintPointType_proto_enumTypes[0]
}

func (x LunaRiteHintPointType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LunaRiteHintPointType.Descriptor instead.
func (LunaRiteHintPointType) EnumDescriptor() ([]byte, []int) {
	return file_LunaRiteHintPointType_proto_rawDescGZIP(), []int{0}
}

var File_LunaRiteHintPointType_proto protoreflect.FileDescriptor

var file_LunaRiteHintPointType_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x4c, 0x75, 0x6e, 0x61, 0x52, 0x69, 0x74, 0x65, 0x48, 0x69, 0x6e, 0x74, 0x50, 0x6f,
	0x69, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x84, 0x01,
	0x0a, 0x15, 0x4c, 0x75, 0x6e, 0x61, 0x52, 0x69, 0x74, 0x65, 0x48, 0x69, 0x6e, 0x74, 0x50, 0x6f,
	0x69, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x1e, 0x4c, 0x55, 0x4e, 0x41, 0x5f,
	0x52, 0x49, 0x54, 0x45, 0x5f, 0x48, 0x49, 0x4e, 0x54, 0x5f, 0x50, 0x4f, 0x49, 0x4e, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x22, 0x0a, 0x1e, 0x4c,
	0x55, 0x4e, 0x41, 0x5f, 0x52, 0x49, 0x54, 0x45, 0x5f, 0x48, 0x49, 0x4e, 0x54, 0x5f, 0x50, 0x4f,
	0x49, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x55, 0x4e, 0x45, 0x10, 0x01, 0x12,
	0x23, 0x0a, 0x1f, 0x4c, 0x55, 0x4e, 0x41, 0x5f, 0x52, 0x49, 0x54, 0x45, 0x5f, 0x48, 0x49, 0x4e,
	0x54, 0x5f, 0x50, 0x4f, 0x49, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x48, 0x45,
	0x53, 0x54, 0x10, 0x02, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_LunaRiteHintPointType_proto_rawDescOnce sync.Once
	file_LunaRiteHintPointType_proto_rawDescData = file_LunaRiteHintPointType_proto_rawDesc
)

func file_LunaRiteHintPointType_proto_rawDescGZIP() []byte {
	file_LunaRiteHintPointType_proto_rawDescOnce.Do(func() {
		file_LunaRiteHintPointType_proto_rawDescData = protoimpl.X.CompressGZIP(file_LunaRiteHintPointType_proto_rawDescData)
	})
	return file_LunaRiteHintPointType_proto_rawDescData
}

var file_LunaRiteHintPointType_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_LunaRiteHintPointType_proto_goTypes = []interface{}{
	(LunaRiteHintPointType)(0), // 0: LunaRiteHintPointType
}
var file_LunaRiteHintPointType_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_LunaRiteHintPointType_proto_init() }
func file_LunaRiteHintPointType_proto_init() {
	if File_LunaRiteHintPointType_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_LunaRiteHintPointType_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_LunaRiteHintPointType_proto_goTypes,
		DependencyIndexes: file_LunaRiteHintPointType_proto_depIdxs,
		EnumInfos:         file_LunaRiteHintPointType_proto_enumTypes,
	}.Build()
	File_LunaRiteHintPointType_proto = out.File
	file_LunaRiteHintPointType_proto_rawDesc = nil
	file_LunaRiteHintPointType_proto_goTypes = nil
	file_LunaRiteHintPointType_proto_depIdxs = nil
}