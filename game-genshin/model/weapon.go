package model

import (
	gdc "game-genshin/config"
)

type Weapon struct {
	WeaponId    uint64   `bson:"weaponId"`    // 武器的唯一id
	ItemId      uint32   `bson:"itemId"`      // 武器的道具id
	Level       uint8    `bson:"level"`       // 等级
	Exp         uint32   `bson:"exp"`         // 当前经验值
	TotalExp    uint32   `bson:"totalExp"`    // 升级所需总经验值
	Promote     uint8    `bson:"promote"`     // 突破等阶
	Lock        bool     `bson:"lock"`        // 锁定状态
	AffixIdList []uint32 `bson:"affixIdList"` // 词缀
	Refinement  uint8    `bson:"refinement"`  // 精炼等阶
	MainPropId  uint32   `bson:"mainPropId"`  // 主词条id
	AvatarId    uint32   `bson:"avatarId"`    // 装备角色id
	Guid        uint64   `bson:"-"`
}

func (p *Player) InitWeapon(weapon *Weapon) *Weapon {
	weapon.Guid = p.GetNextGameObjectGuid()
	return weapon
}

func (p *Player) InitAllWeapon() {
	for weaponId, weapon := range p.WeaponMap {
		p.WeaponMap[weaponId] = p.InitWeapon(weapon)
		if weapon.AvatarId != 0 {
			p.AvatarMap[weapon.AvatarId].EquipGuidList = append(p.AvatarMap[weapon.AvatarId].EquipGuidList, weapon.Guid)
			p.AvatarMap[weapon.AvatarId].EquipWeapon = weapon
		}
	}
}

func (p *Player) GetWeaponGuid(weaponId uint64) uint64 {
	weaponInfo := p.WeaponMap[weaponId]
	if weaponInfo == nil {
		return 0
	}
	return weaponInfo.Guid
}

func (p *Player) GetWeapon(weaponId uint64) *Weapon {
	return p.WeaponMap[weaponId]
}

func (p *Player) AddWeapon(itemId uint32, weaponId uint64) {
	weapon := &Weapon{
		WeaponId:    weaponId,
		ItemId:      itemId,
		Level:       1,
		Exp:         0,
		TotalExp:    0,
		Promote:     0,
		Lock:        false,
		AffixIdList: make([]uint32, 0),
		Refinement:  0,
		MainPropId:  0,
		Guid:        0,
	}
	itemDataConfig := gdc.CONF.ItemDataMap[int32(itemId)]
	if itemDataConfig.SkillAffix != nil {
		for _, skillAffix := range itemDataConfig.SkillAffix {
			if skillAffix > 0 {
				weapon.AffixIdList = append(weapon.AffixIdList, uint32(skillAffix))
			}
		}
	}
	p.InitWeapon(weapon)
	p.WeaponMap[weaponId] = weapon
}

func (p *Player) EquipWeaponToAvatar(avatarId uint32, weaponId uint64) {
	avatar := p.AvatarMap[avatarId]
	weapon := p.WeaponMap[weaponId]
	avatar.EquipWeapon = weapon
	weapon.AvatarId = avatarId
}
