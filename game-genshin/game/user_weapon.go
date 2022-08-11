package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
)

func (g *GameManager) AddUserWeapon(userId uint32, itemId uint32) {
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		logger.LOG.Error("player not found, user id: %v", userId)
		return
	}
	weaponId := uint64(g.snowflake.GenId())
	player.AddWeapon(itemId, weaponId)
	weapon := player.GetWeapon(weaponId)

	// PacketStoreItemChangeNotify
	storeItemChangeNotify := new(proto.StoreItemChangeNotify)
	storeItemChangeNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	affixMap := make(map[uint32]uint32)
	for _, affixId := range weapon.AffixIdList {
		affixMap[affixId] = uint32(weapon.Refinement)
	}
	pbItem := &proto.Item{
		ItemId: itemId,
		Guid:   player.GetWeaponGuid(weaponId),
		Detail: &proto.Item_Equip{
			Equip: &proto.Equip{
				Detail: &proto.Equip_Weapon{
					Weapon: &proto.Weapon{
						Level:        uint32(weapon.Level),
						Exp:          weapon.Exp,
						PromoteLevel: uint32(weapon.Promote),
						// key:武器效果id value:精炼等阶
						AffixMap: affixMap,
					},
				},
				IsLocked: weapon.Lock,
			},
		},
	}
	storeItemChangeNotify.ItemList = append(storeItemChangeNotify.ItemList, pbItem)
	g.SendMsg(api.ApiStoreItemChangeNotify, userId, nil, storeItemChangeNotify)
}

func (g *GameManager) EquipUserWeaponToAvatar(userId uint32, avatarId uint32, weaponId uint64) {
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		logger.LOG.Error("player not found, user id: %v", userId)
		return
	}
	player.EquipWeaponToAvatar(avatarId, weaponId)
}
