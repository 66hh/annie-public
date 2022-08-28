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
// source: Unk2700_MOFABPNGIKP.proto

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

type Unk2700_MOFABPNGIKP int32

const (
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_DGJFKKIBLCJ Unk2700_MOFABPNGIKP = 0
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_IANMLLDEIJH Unk2700_MOFABPNGIKP = 1
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_CCBNMEBCOKM Unk2700_MOFABPNGIKP = 2
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_CABFGAEJAIA Unk2700_MOFABPNGIKP = 3
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_JFPKBELPINO Unk2700_MOFABPNGIKP = 4
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_ECHKDKLKPLH Unk2700_MOFABPNGIKP = 5
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_NALBIHIEGAF Unk2700_MOFABPNGIKP = 6
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_KNAHCHDLEOM Unk2700_MOFABPNGIKP = 7
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_CAIOEECIPIM Unk2700_MOFABPNGIKP = 8
	Unk2700_MOFABPNGIKP_Unk2700_MOFABPNGIKP_Unk2700_IEICHGLOIAL Unk2700_MOFABPNGIKP = 9
)

// Enum value maps for Unk2700_MOFABPNGIKP.
var (
	Unk2700_MOFABPNGIKP_name = map[int32]string{
		0: "Unk2700_MOFABPNGIKP_Unk2700_DGJFKKIBLCJ",
		1: "Unk2700_MOFABPNGIKP_Unk2700_IANMLLDEIJH",
		2: "Unk2700_MOFABPNGIKP_Unk2700_CCBNMEBCOKM",
		3: "Unk2700_MOFABPNGIKP_Unk2700_CABFGAEJAIA",
		4: "Unk2700_MOFABPNGIKP_Unk2700_JFPKBELPINO",
		5: "Unk2700_MOFABPNGIKP_Unk2700_ECHKDKLKPLH",
		6: "Unk2700_MOFABPNGIKP_Unk2700_NALBIHIEGAF",
		7: "Unk2700_MOFABPNGIKP_Unk2700_KNAHCHDLEOM",
		8: "Unk2700_MOFABPNGIKP_Unk2700_CAIOEECIPIM",
		9: "Unk2700_MOFABPNGIKP_Unk2700_IEICHGLOIAL",
	}
	Unk2700_MOFABPNGIKP_value = map[string]int32{
		"Unk2700_MOFABPNGIKP_Unk2700_DGJFKKIBLCJ": 0,
		"Unk2700_MOFABPNGIKP_Unk2700_IANMLLDEIJH": 1,
		"Unk2700_MOFABPNGIKP_Unk2700_CCBNMEBCOKM": 2,
		"Unk2700_MOFABPNGIKP_Unk2700_CABFGAEJAIA": 3,
		"Unk2700_MOFABPNGIKP_Unk2700_JFPKBELPINO": 4,
		"Unk2700_MOFABPNGIKP_Unk2700_ECHKDKLKPLH": 5,
		"Unk2700_MOFABPNGIKP_Unk2700_NALBIHIEGAF": 6,
		"Unk2700_MOFABPNGIKP_Unk2700_KNAHCHDLEOM": 7,
		"Unk2700_MOFABPNGIKP_Unk2700_CAIOEECIPIM": 8,
		"Unk2700_MOFABPNGIKP_Unk2700_IEICHGLOIAL": 9,
	}
)

func (x Unk2700_MOFABPNGIKP) Enum() *Unk2700_MOFABPNGIKP {
	p := new(Unk2700_MOFABPNGIKP)
	*p = x
	return p
}

func (x Unk2700_MOFABPNGIKP) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Unk2700_MOFABPNGIKP) Descriptor() protoreflect.EnumDescriptor {
	return file_Unk2700_MOFABPNGIKP_proto_enumTypes[0].Descriptor()
}

func (Unk2700_MOFABPNGIKP) Type() protoreflect.EnumType {
	return &file_Unk2700_MOFABPNGIKP_proto_enumTypes[0]
}

func (x Unk2700_MOFABPNGIKP) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Unk2700_MOFABPNGIKP.Descriptor instead.
func (Unk2700_MOFABPNGIKP) EnumDescriptor() ([]byte, []int) {
	return file_Unk2700_MOFABPNGIKP_proto_rawDescGZIP(), []int{0}
}

var File_Unk2700_MOFABPNGIKP_proto protoreflect.FileDescriptor

var file_Unk2700_MOFABPNGIKP_proto_rawDesc = []byte{
	0x0a, 0x19, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50,
	0x4e, 0x47, 0x49, 0x4b, 0x50, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xd7, 0x03, 0x0a, 0x13,
	0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47,
	0x49, 0x4b, 0x50, 0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d,
	0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37,
	0x30, 0x30, 0x5f, 0x44, 0x47, 0x4a, 0x46, 0x4b, 0x4b, 0x49, 0x42, 0x4c, 0x43, 0x4a, 0x10, 0x00,
	0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41,
	0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f,
	0x49, 0x41, 0x4e, 0x4d, 0x4c, 0x4c, 0x44, 0x45, 0x49, 0x4a, 0x48, 0x10, 0x01, 0x12, 0x2b, 0x0a,
	0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e,
	0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x43, 0x43, 0x42,
	0x4e, 0x4d, 0x45, 0x42, 0x43, 0x4f, 0x4b, 0x4d, 0x10, 0x02, 0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e,
	0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b,
	0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x43, 0x41, 0x42, 0x46, 0x47, 0x41,
	0x45, 0x4a, 0x41, 0x49, 0x41, 0x10, 0x03, 0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37,
	0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55,
	0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4a, 0x46, 0x50, 0x4b, 0x42, 0x45, 0x4c, 0x50, 0x49,
	0x4e, 0x4f, 0x10, 0x04, 0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f,
	0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32,
	0x37, 0x30, 0x30, 0x5f, 0x45, 0x43, 0x48, 0x4b, 0x44, 0x4b, 0x4c, 0x4b, 0x50, 0x4c, 0x48, 0x10,
	0x05, 0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46,
	0x41, 0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30,
	0x5f, 0x4e, 0x41, 0x4c, 0x42, 0x49, 0x48, 0x49, 0x45, 0x47, 0x41, 0x46, 0x10, 0x06, 0x12, 0x2b,
	0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50,
	0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4b, 0x4e,
	0x41, 0x48, 0x43, 0x48, 0x44, 0x4c, 0x45, 0x4f, 0x4d, 0x10, 0x07, 0x12, 0x2b, 0x0a, 0x27, 0x55,
	0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47, 0x49,
	0x4b, 0x50, 0x5f, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x43, 0x41, 0x49, 0x4f, 0x45,
	0x45, 0x43, 0x49, 0x50, 0x49, 0x4d, 0x10, 0x08, 0x12, 0x2b, 0x0a, 0x27, 0x55, 0x6e, 0x6b, 0x32,
	0x37, 0x30, 0x30, 0x5f, 0x4d, 0x4f, 0x46, 0x41, 0x42, 0x50, 0x4e, 0x47, 0x49, 0x4b, 0x50, 0x5f,
	0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x5f, 0x49, 0x45, 0x49, 0x43, 0x48, 0x47, 0x4c, 0x4f,
	0x49, 0x41, 0x4c, 0x10, 0x09, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Unk2700_MOFABPNGIKP_proto_rawDescOnce sync.Once
	file_Unk2700_MOFABPNGIKP_proto_rawDescData = file_Unk2700_MOFABPNGIKP_proto_rawDesc
)

func file_Unk2700_MOFABPNGIKP_proto_rawDescGZIP() []byte {
	file_Unk2700_MOFABPNGIKP_proto_rawDescOnce.Do(func() {
		file_Unk2700_MOFABPNGIKP_proto_rawDescData = protoimpl.X.CompressGZIP(file_Unk2700_MOFABPNGIKP_proto_rawDescData)
	})
	return file_Unk2700_MOFABPNGIKP_proto_rawDescData
}

var file_Unk2700_MOFABPNGIKP_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_Unk2700_MOFABPNGIKP_proto_goTypes = []interface{}{
	(Unk2700_MOFABPNGIKP)(0), // 0: Unk2700_MOFABPNGIKP
}
var file_Unk2700_MOFABPNGIKP_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_Unk2700_MOFABPNGIKP_proto_init() }
func file_Unk2700_MOFABPNGIKP_proto_init() {
	if File_Unk2700_MOFABPNGIKP_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_Unk2700_MOFABPNGIKP_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Unk2700_MOFABPNGIKP_proto_goTypes,
		DependencyIndexes: file_Unk2700_MOFABPNGIKP_proto_depIdxs,
		EnumInfos:         file_Unk2700_MOFABPNGIKP_proto_enumTypes,
	}.Build()
	File_Unk2700_MOFABPNGIKP_proto = out.File
	file_Unk2700_MOFABPNGIKP_proto_rawDesc = nil
	file_Unk2700_MOFABPNGIKP_proto_goTypes = nil
	file_Unk2700_MOFABPNGIKP_proto_depIdxs = nil
}