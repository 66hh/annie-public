package entity

import (
	"flswld.com/gate-genshin-api/api/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DbInsert = iota
	DbDelete
	DbUpdate
	DbNormal
)

type Player struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"`
	PlayerID        uint32              `bson:"playerID"`     // 玩家uid
	NickName        string              `bson:"nickname"`     // 玩家昵称
	Properties      map[uint16]uint32   `bson:"properties"`   // 玩家自身相关的一些属性
	MpSetting       proto.MpSettingType `bson:"mpSetting"`    // 世界权限
	RegionId        uint16              `bson:"regionId"`     // regionId
	FlyCloakList    []uint32            `bson:"flyCloakList"` // 风之翼列表
	SceneId         uint16              `bson:"sceneId"`      // 场景
	Pos             *Vector             `bson:"pos"`          // 玩家坐标
	Rotation        *Vector             `bson:"rotation"`     // 玩家朝向
	EnterSceneToken uint32              `bson:"-"`            // 玩家的世界进入令牌
	DbState         int                 `bson:"-"`            // 数据库存档状态
	AvatarEntityId  uint32              `bson:"-"`
	WeaponEntityId  uint32              `bson:"-"`
}
