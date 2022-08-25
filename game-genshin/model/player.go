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

type GameObject interface {
}

type Player struct {
	// 离线数据
	ID              primitive.ObjectID    `bson:"_id,omitempty"`
	PlayerID        uint32                `bson:"playerID"`        // 玩家uid
	NickName        string                `bson:"nickname"`        // 玩家昵称
	Signature       string                `bson:"signature"`       // 玩家签名
	HeadImage       uint32                `bson:"headImage"`       // 玩家头像
	NameCard        uint32                `bson:"nameCard"`        // 当前名片
	NameCardList    []uint32              `bson:"nameCardList"`    // 已解锁名片列表
	FriendList      []uint32              `bson:"friendList"`      // 好友uid列表
	FriendApplyList []uint32              `bson:"friendApplyList"` // 好友申请uid列表
	OfflineTime     uint32                `bson:"offlineTime"`     // 离线时间
	Properties      map[uint16]uint32     `bson:"properties"`      // 玩家自身相关的一些属性
	MpSetting       proto.MpSettingType   `bson:"mpSetting"`       // 世界权限
	RegionId        uint32                `bson:"regionId"`        // regionId
	FlyCloakList    []uint32              `bson:"flyCloakList"`    // 风之翼列表
	CostumeList     []uint32              `bson:"costumeList"`     // 角色衣装列表
	SceneId         uint32                `bson:"sceneId"`         // 场景
	Pos             *Vector               `bson:"pos"`             // 玩家坐标
	Rot             *Vector               `bson:"rot"`             // 玩家朝向
	ItemMap         map[uint32]*Item      `bson:"itemMap"`         // 玩家统一大背包仓库
	WeaponMap       map[uint64]*Weapon    `bson:"weaponMap"`       // 玩家武器背包
	ReliquaryMap    map[uint64]*Reliquary `bson:"reliquaryMap"`    // 玩家圣遗物背包
	TeamConfig      *TeamInfo             `bson:"teamConfig"`      // 队伍配置
	AvatarMap       map[uint32]*Avatar    `bson:"avatarMap"`       // 角色信息
	DropInfo        *DropInfo             `bson:"dropInfo"`        // 掉落信息
	// 在线数据
	EnterSceneToken       uint32                `bson:"-"` // 玩家的世界进入令牌
	DbState               int                   `bson:"-"` // 数据库存档状态
	WorldId               uint32                `bson:"-"` // 所在的世界id
	PeerId                uint32                `bson:"-"` // 多人世界的玩家编号 待确定
	GameObjectGuidCounter uint64                `bson:"-"` // 游戏对象guid计数器
	ClientTime            uint32                `bson:"-"` // 玩家客户端的本地时钟
	ClientRTT             uint32                `bson:"-"` // 玩家客户端往返时延
	GameObjectGuidMap     map[uint64]GameObject `bson:"-"` // 游戏对象guid映射表
	Online                bool                  `bson:"-"` // 在线状态
	Pause                 bool                  `bson:"-"` // 暂停状态
	BornInScene           bool                  `bson:"-"` // 是否在场景内发送过出生通知
	CoopApplyMap          map[uint32]int64      `bson:"-"` // 敲门申请的玩家uid及时间
}

func (p *Player) GetNextGameObjectGuid() uint64 {
	p.GameObjectGuidCounter++
	return uint64(p.PlayerID)<<32 + p.GameObjectGuidCounter
}

func (p *Player) InitAll() {
	p.GameObjectGuidMap = make(map[uint64]GameObject)
	p.CoopApplyMap = make(map[uint32]int64)
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
			avatar := p.AvatarMap[reliquary.AvatarId]
			avatar.EquipGuidList[reliquary.Guid] = reliquary.Guid
			avatar.EquipReliquaryList = append(avatar.EquipReliquaryList, reliquary)
		}
	}
}
