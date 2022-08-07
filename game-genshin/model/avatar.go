package model

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
