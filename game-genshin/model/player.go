package model

import (
	"flswld.com/gate-genshin-api/api/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DbInsert = iota
	DbDelete
	DbUpdate
	DbNormal
	DbOffline
)

type Player struct {
	ID                    primitive.ObjectID    `bson:"_id,omitempty"`
	PlayerID              uint32                `bson:"playerID"`     // 玩家uid
	NickName              string                `bson:"nickname"`     // 玩家昵称
	Properties            map[uint16]uint32     `bson:"properties"`   // 玩家自身相关的一些属性
	MpSetting             proto.MpSettingType   `bson:"mpSetting"`    // 世界权限
	RegionId              uint32                `bson:"regionId"`     // regionId
	FlyCloakList          []uint32              `bson:"flyCloakList"` // 风之翼列表
	CostumeList           []uint32              `bson:"costumeList"`  // 角色衣装列表
	SceneId               uint32                `bson:"sceneId"`      // 场景
	Pos                   *Vector               `bson:"pos"`          // 玩家坐标
	Rotation              *Vector               `bson:"rotation"`     // 玩家朝向
	ItemMap               map[uint32]*Item      `bson:"itemMap"`      // 玩家统一大背包仓库
	WeaponMap             map[uint64]*Weapon    `bson:"weaponMap"`    // 玩家武器背包
	ReliquaryMap          map[uint64]*Reliquary `bson:"reliquaryMap"` // 玩家圣遗物背包
	TeamConfig            *TeamInfo             `bson:"teamConfig"`   // 队伍配置
	AvatarMap             map[uint32]*Avatar    `bson:"avatarMap"`    // 角色信息
	DropInfo              *DropInfo             `bson:"dropInfo"`     // 掉落信息
	EnterSceneToken       uint32                `bson:"-"`            // 玩家的世界进入令牌
	DbState               int                   `bson:"-"`            // 数据库存档状态
	WorldId               uint32                `bson:"-"`            // 所在的世界id
	PeerId                uint32                `bson:"-"`
	GameObjectGuidCounter uint64                `bson:"-"` // 游戏对象guid计数器
}

func (p *Player) GetNextGameObjectGuid() uint64 {
	p.GameObjectGuidCounter++
	return uint64(p.PlayerID)<<32 + p.GameObjectGuidCounter
}

func (p *Player) InitAll() {
	p.InitAllAvatar()
	p.InitAllWeapon()
	p.InitAllItem()
	p.InitAllReliquary()
}

func (p *Player) InitAllReliquary() {
	for reliquaryId, reliquary := range p.ReliquaryMap {
		reliquary.Guid = p.GetNextGameObjectGuid()
		p.ReliquaryMap[reliquaryId] = reliquary
		if reliquary.AvatarId != 0 {
			p.AvatarMap[reliquary.AvatarId].EquipGuidList = append(p.AvatarMap[reliquary.AvatarId].EquipGuidList, reliquary.Guid)
			p.AvatarMap[reliquary.AvatarId].EquipReliquaryList = append(p.AvatarMap[reliquary.AvatarId].EquipReliquaryList, reliquary)
		}
	}
}
