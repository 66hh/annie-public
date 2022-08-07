package model

import (
	"flswld.com/gate-genshin-api/api/proto"
	"game-genshin/config"
	"game-genshin/game/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	DbInsert = iota
	DbDelete
	DbUpdate
	DbNormal
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

func (p *Player) AddAvatar(avatarId uint32, avatarDataConfig *config.AvatarData, avatarSkillDepotDataConfig *config.AvatarSkillDepotData) {
	avatar := &Avatar{
		AvatarId:            avatarId,
		Level:               1,
		Exp:                 0,
		Promote:             0,
		Satiation:           0,
		SatiationPenalty:    0,
		CurrHP:              0,
		CurrEnergy:          0,
		FetterList:          nil,
		SkillLevelMap:       make(map[uint32]uint32),
		SkillExtraChargeMap: make(map[uint32]uint32),
		ProudSkillBonusMap:  nil,
		SkillDepotId:        uint32(avatarSkillDepotDataConfig.Id),
		CoreProudSkillLevel: 0,
		TalentIdList:        make([]uint32, 0),
		ProudSkillList:      make([]uint32, 0),
		FlyCloak:            140001,
		Costume:             0,
		BornTime:            time.Now().Unix(),
		FetterLevel:         1,
		FetterExp:           0,
		NameCardRewardId:    0,
		NameCardId:          0,
		Guid:                0,
		EquipGuidList:       nil,
		EquipWeapon:         nil,
		EquipReliquaryList:  nil,
		FightPropMap:        nil,
	}

	if avatarSkillDepotDataConfig.EnergySkill > 0 {
		avatar.SkillLevelMap[uint32(avatarSkillDepotDataConfig.EnergySkill)] = 1
	}
	for _, skillId := range avatarSkillDepotDataConfig.Skills {
		if skillId > 0 {
			avatar.SkillLevelMap[uint32(skillId)] = 1
		}
	}
	for _, openData := range avatarSkillDepotDataConfig.InherentProudSkillOpens {
		if openData.ProudSkillGroupId == 0 {
			continue
		}
		if openData.NeedAvatarPromoteLevel <= int32(avatar.Promote) {
			proudSkillId := (openData.ProudSkillGroupId * 100) + 1
			// TODO if GameData.getProudSkillDataMap().containsKey(proudSkillId) java
			avatar.ProudSkillList = append(avatar.ProudSkillList, uint32(proudSkillId))
		}
	}

	p.InitAvatar(avatar, avatarDataConfig)

	fightPropertyConst := constant.GetFightPropertyConst()
	avatar.CurrHP = float64(avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)])

	p.AvatarMap[avatarId] = avatar
}

func (p *Player) InitAvatar(avatar *Avatar, avatarDataConfig *config.AvatarData) *Avatar {
	// 战斗属性
	fightPropertyConst := constant.GetFightPropertyConst()
	avatar.FightPropMap = make(map[uint32]float32)
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_NONE)] = 0.0
	// 白字
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)] = float32(avatarDataConfig.GetBaseHpByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)] = float32(avatarDataConfig.GetBaseAttackByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)] = float32(avatarDataConfig.GetBaseDefenseByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)] = float32(avatarDataConfig.Critical)
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)] = float32(avatarDataConfig.CriticalHurt)
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)] = 1.0
	// 绿字
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)] = float32(avatarDataConfig.GetBaseHpByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)] = float32(avatarDataConfig.GetBaseHpByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)] = float32(avatarDataConfig.GetBaseAttackByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)] = float32(avatarDataConfig.GetBaseDefenseByLevel(avatar.Level))
	// guid
	avatar.Guid = p.GetNextGameObjectGuid()
	return avatar
}

func (p *Player) InitAllAvatar(avatarDataMapConfig map[int32]*config.AvatarData) {
	for avatarId, avatar := range p.AvatarMap {
		p.AvatarMap[avatarId] = p.InitAvatar(avatar, avatarDataMapConfig[int32(avatarId)])
	}
}

func (p *Player) AddWeapon(itemId uint32, weaponId uint64) {
	p.WeaponMap[weaponId] = &Weapon{
		WeaponId:   weaponId,
		ItemId:     itemId,
		Level:      1,
		Exp:        0,
		TotalExp:   0,
		Promote:    0,
		Lock:       false,
		Refinement: 0,
		MainPropId: 0,
	}
}

func (p *Player) InitAllWeapon() {
	for weaponId, weapon := range p.WeaponMap {
		weapon.Guid = p.GetNextGameObjectGuid()
		p.WeaponMap[weaponId] = weapon
		if weapon.AvatarId != 0 {
			p.AvatarMap[weapon.AvatarId].EquipGuidList = append(p.AvatarMap[weapon.AvatarId].EquipGuidList, weapon.Guid)
			p.AvatarMap[weapon.AvatarId].EquipWeapon = weapon
		}
	}
}

func (p *Player) InitAllItem() {
	for itemId, item := range p.ItemMap {
		item.Guid = p.GetNextGameObjectGuid()
		p.ItemMap[itemId] = item
	}
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

func (p *Player) InitAll(gameDataConfig *config.GameDataConfig) {
	p.InitAllAvatar(gameDataConfig.AvatarDataMap)
	p.InitAllWeapon()
	p.InitAllItem()
	p.InitAllReliquary()
}

func (p *Player) AvatarEquipWeapon(avatarId uint32, weaponId uint64) {
	avatar := p.AvatarMap[avatarId]
	weapon := p.WeaponMap[weaponId]
	avatar.EquipWeapon = weapon
	weapon.AvatarId = avatarId
}
