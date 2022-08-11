package model

import (
	gdc "game-genshin/config"
	"game-genshin/constant"
	"time"
)

type Avatar struct {
	AvatarId            uint32             `bson:"avatarId"`         // 角色id
	Level               uint8              `bson:"level"`            // 等级
	Exp                 uint32             `bson:"exp"`              // 经验值
	Promote             uint8              `bson:"promote"`          // 突破等阶
	Satiation           uint32             `bson:"satiation"`        // 饱食度
	SatiationPenalty    uint32             `bson:"satiationPenalty"` // 饱食度溢出
	CurrHP              float64            `bson:"currHP"`           // 当前生命值
	CurrEnergy          float64            `bson:"currEnergy"`       // 当前元素能量值
	FetterList          []uint32           `bson:"fetterList"`       // 资料解锁条目
	SkillLevelMap       map[uint32]uint32  `bson:"skillLevelMap"`    // 技能等级
	SkillExtraChargeMap map[uint32]uint32  `bson:"skillExtraChargeMap"`
	ProudSkillBonusMap  map[uint32]uint32  `bson:"proudSkillBonusMap"`
	SkillDepotId        uint32             `bson:"skillDepotId"`
	CoreProudSkillLevel uint8              `bson:"coreProudSkillLevel"` // 已解锁命之座层数
	TalentIdList        []uint32           `bson:"talentIdList"`        // 已解锁命之座技能列表
	ProudSkillList      []uint32           `bson:"proudSkillList"`      // 被动技能列表
	FlyCloak            uint32             `bson:"flyCloak"`            // 当前风之翼
	Costume             uint32             `bson:"costume"`             // 当前衣装
	BornTime            int64              `bson:"bornTime"`            // 获得时间
	FetterLevel         uint8              `bson:"fetterLevel"`         // 好感度等级
	FetterExp           uint32             `bson:"fetterExp"`           // 好感度经验
	NameCardRewardId    uint32             `bson:"nameCardRewardId"`
	NameCardId          uint32             `bson:"nameCardId"`
	Guid                uint64             `bson:"-"`
	EquipGuidList       []uint64           `bson:"-"`
	EquipWeapon         *Weapon            `bson:"-"`
	EquipReliquaryList  []*Reliquary       `bson:"-"`
	FightPropMap        map[uint32]float32 `bson:"-"`
	ExtraAbilityEmbryos map[string]bool    `bson:"-"`
}

func (p *Player) InitAllAvatar() {
	for avatarId, avatar := range p.AvatarMap {
		p.AvatarMap[avatarId] = p.InitAvatar(avatar)
	}
}

func (p *Player) InitAvatar(avatar *Avatar) *Avatar {
	avatarDataConfig := gdc.CONF.AvatarDataMap[int32(avatar.AvatarId)]
	// 角色战斗属性
	fightPropertyConst := constant.GetFightPropertyConst()
	avatar.FightPropMap = make(map[uint32]float32)
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_NONE)] = 0.0
	// 白字攻防血
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)] = float32(avatarDataConfig.GetBaseAttackByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)] = float32(avatarDataConfig.GetBaseDefenseByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)] = float32(avatarDataConfig.GetBaseHpByLevel(avatar.Level))
	// 白字+绿字攻防血
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)] = float32(avatarDataConfig.GetBaseAttackByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)] = float32(avatarDataConfig.GetBaseDefenseByLevel(avatar.Level))
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)] = float32(avatarDataConfig.GetBaseHpByLevel(avatar.Level))
	// 当前血量
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)] = float32(avatar.CurrHP)
	// 双暴
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)] = float32(avatarDataConfig.Critical)
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)] = float32(avatarDataConfig.CriticalHurt)
	// 元素充能
	avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)] = 1.0
	p.SetCurrEnergy(avatar, avatar.CurrEnergy, true)
	// guid
	avatar.Guid = p.GetNextGameObjectGuid()
	return avatar
}

func (p *Player) AddAvatar(avatarId uint32) {
	avatarDataConfig := gdc.CONF.AvatarDataMap[int32(avatarId)]
	skillDepotId := int32(0)
	// 主角要单独设置
	if avatarId == 10000005 {
		skillDepotId = 504
	} else if avatarId == 10000007 {
		skillDepotId = 704
	} else {
		skillDepotId = avatarDataConfig.SkillDepotId
	}
	avatarSkillDepotDataConfig := gdc.CONF.AvatarSkillDepotDataMap[skillDepotId]
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
		ExtraAbilityEmbryos: nil,
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
	avatar.CurrHP = avatarDataConfig.GetBaseHpByLevel(avatar.Level)

	p.InitAvatar(avatar)
	p.AvatarMap[avatarId] = avatar
}

func (p *Player) SetCurrEnergy(avatar *Avatar, value float64, max bool) {
	avatarDataConfig := gdc.CONF.AvatarDataMap[int32(avatar.AvatarId)]
	avatarSkillDepotDataConfig := gdc.CONF.AvatarSkillDepotDataMap[avatarDataConfig.SkillDepotId]
	if avatarSkillDepotDataConfig == nil || avatarSkillDepotDataConfig.EnergySkillData == nil {
		return
	}
	element := avatarSkillDepotDataConfig.ElementType
	avatar.FightPropMap[uint32(element.MaxEnergyProp)] = float32(avatarSkillDepotDataConfig.EnergySkillData.CostElemVal)
	if max {
		avatar.FightPropMap[uint32(element.CurrEnergyProp)] = float32(avatarSkillDepotDataConfig.EnergySkillData.CostElemVal)
	} else {
		avatar.FightPropMap[uint32(element.CurrEnergyProp)] = float32(value)
	}
}
