package game

import (
	"flswld.com/logger"
	gdc "game-genshin/config"
)

func (g *GameManager) AddUserAvatar(userId uint32, avatarId uint32) {
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		logger.LOG.Error("player not found, user id: %v", userId)
		return
	}
	player.AddAvatar(avatarId)
	// 添加初始武器
	avatarDataConfig := gdc.CONF.AvatarDataMap[int32(avatarId)]
	weaponId := uint64(g.snowflake.GenId())
	player.AddWeapon(uint32(avatarDataConfig.InitialWeapon), weaponId)
	// 角色装上初始武器
	player.EquipWeaponToAvatar(avatarId, weaponId)
}
