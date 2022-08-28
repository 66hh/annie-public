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
// source: TowerCurLevelRecord.proto

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

type TowerCurLevelRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TowerTeamList       []*TowerTeam `protobuf:"bytes,8,rep,name=tower_team_list,json=towerTeamList,proto3" json:"tower_team_list,omitempty"`
	IsEmpty             bool         `protobuf:"varint,6,opt,name=is_empty,json=isEmpty,proto3" json:"is_empty,omitempty"`
	BuffIdList          []uint32     `protobuf:"varint,4,rep,packed,name=buff_id_list,json=buffIdList,proto3" json:"buff_id_list,omitempty"`
	Unk2700_CBPNPEBMPOH bool         `protobuf:"varint,2,opt,name=Unk2700_CBPNPEBMPOH,json=Unk2700CBPNPEBMPOH,proto3" json:"Unk2700_CBPNPEBMPOH,omitempty"`
	CurLevelIndex       uint32       `protobuf:"varint,1,opt,name=cur_level_index,json=curLevelIndex,proto3" json:"cur_level_index,omitempty"`
	CurFloorId          uint32       `protobuf:"varint,15,opt,name=cur_floor_id,json=curFloorId,proto3" json:"cur_floor_id,omitempty"`
}

func (x *TowerCurLevelRecord) Reset() {
	*x = TowerCurLevelRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TowerCurLevelRecord_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TowerCurLevelRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TowerCurLevelRecord) ProtoMessage() {}

func (x *TowerCurLevelRecord) ProtoReflect() protoreflect.Message {
	mi := &file_TowerCurLevelRecord_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TowerCurLevelRecord.ProtoReflect.Descriptor instead.
func (*TowerCurLevelRecord) Descriptor() ([]byte, []int) {
	return file_TowerCurLevelRecord_proto_rawDescGZIP(), []int{0}
}

func (x *TowerCurLevelRecord) GetTowerTeamList() []*TowerTeam {
	if x != nil {
		return x.TowerTeamList
	}
	return nil
}

func (x *TowerCurLevelRecord) GetIsEmpty() bool {
	if x != nil {
		return x.IsEmpty
	}
	return false
}

func (x *TowerCurLevelRecord) GetBuffIdList() []uint32 {
	if x != nil {
		return x.BuffIdList
	}
	return nil
}

func (x *TowerCurLevelRecord) GetUnk2700_CBPNPEBMPOH() bool {
	if x != nil {
		return x.Unk2700_CBPNPEBMPOH
	}
	return false
}

func (x *TowerCurLevelRecord) GetCurLevelIndex() uint32 {
	if x != nil {
		return x.CurLevelIndex
	}
	return 0
}

func (x *TowerCurLevelRecord) GetCurFloorId() uint32 {
	if x != nil {
		return x.CurFloorId
	}
	return 0
}

var File_TowerCurLevelRecord_proto protoreflect.FileDescriptor

var file_TowerCurLevelRecord_proto_rawDesc = []byte{
	0x0a, 0x19, 0x54, 0x6f, 0x77, 0x65, 0x72, 0x43, 0x75, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x54, 0x6f, 0x77,
	0x65, 0x72, 0x54, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x02, 0x0a,
	0x13, 0x54, 0x6f, 0x77, 0x65, 0x72, 0x43, 0x75, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x12, 0x32, 0x0a, 0x0f, 0x74, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x74, 0x65,
	0x61, 0x6d, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x54, 0x6f, 0x77, 0x65, 0x72, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x0d, 0x74, 0x6f, 0x77, 0x65, 0x72,
	0x54, 0x65, 0x61, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x0c, 0x62, 0x75, 0x66, 0x66, 0x5f, 0x69, 0x64, 0x5f, 0x6c,
	0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x0a, 0x62, 0x75, 0x66, 0x66, 0x49,
	0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x13, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30,
	0x5f, 0x43, 0x42, 0x50, 0x4e, 0x50, 0x45, 0x42, 0x4d, 0x50, 0x4f, 0x48, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x12, 0x55, 0x6e, 0x6b, 0x32, 0x37, 0x30, 0x30, 0x43, 0x42, 0x50, 0x4e, 0x50,
	0x45, 0x42, 0x4d, 0x50, 0x4f, 0x48, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x5f, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0d, 0x63, 0x75, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20,
	0x0a, 0x0c, 0x63, 0x75, 0x72, 0x5f, 0x66, 0x6c, 0x6f, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x75, 0x72, 0x46, 0x6c, 0x6f, 0x6f, 0x72, 0x49, 0x64,
	0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_TowerCurLevelRecord_proto_rawDescOnce sync.Once
	file_TowerCurLevelRecord_proto_rawDescData = file_TowerCurLevelRecord_proto_rawDesc
)

func file_TowerCurLevelRecord_proto_rawDescGZIP() []byte {
	file_TowerCurLevelRecord_proto_rawDescOnce.Do(func() {
		file_TowerCurLevelRecord_proto_rawDescData = protoimpl.X.CompressGZIP(file_TowerCurLevelRecord_proto_rawDescData)
	})
	return file_TowerCurLevelRecord_proto_rawDescData
}

var file_TowerCurLevelRecord_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_TowerCurLevelRecord_proto_goTypes = []interface{}{
	(*TowerCurLevelRecord)(nil), // 0: TowerCurLevelRecord
	(*TowerTeam)(nil),           // 1: TowerTeam
}
var file_TowerCurLevelRecord_proto_depIdxs = []int32{
	1, // 0: TowerCurLevelRecord.tower_team_list:type_name -> TowerTeam
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_TowerCurLevelRecord_proto_init() }
func file_TowerCurLevelRecord_proto_init() {
	if File_TowerCurLevelRecord_proto != nil {
		return
	}
	file_TowerTeam_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_TowerCurLevelRecord_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TowerCurLevelRecord); i {
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
			RawDescriptor: file_TowerCurLevelRecord_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_TowerCurLevelRecord_proto_goTypes,
		DependencyIndexes: file_TowerCurLevelRecord_proto_depIdxs,
		MessageInfos:      file_TowerCurLevelRecord_proto_msgTypes,
	}.Build()
	File_TowerCurLevelRecord_proto = out.File
	file_TowerCurLevelRecord_proto_rawDesc = nil
	file_TowerCurLevelRecord_proto_goTypes = nil
	file_TowerCurLevelRecord_proto_depIdxs = nil
}