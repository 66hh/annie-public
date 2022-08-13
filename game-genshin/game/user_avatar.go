package game

import (
	"flswld.com/logger"
	gdc "game-genshin/config"
)

func (g *GameManager) GetAllAvatarDataConfig() map[int32]*gdc.AvatarData {
	allAvatarDataConfig := make(map[int32]*gdc.AvatarData)
	for avatarId, avatarData := range gdc.CONF.AvatarDataMap {
		if avatarId < 10000002 || avatarId >= 11000000 {
			// 跳过无效角色
			continue
		}
		if avatarId == 10000005 || avatarId == 10000007 {
			// 跳过主角
			continue
		}
		allAvatarDataConfig[avatarId] = avatarData
	}
	return allAvatarDataConfig
}

func (g *GameManager) AddUserAvatar(userId uint32, avatarId uint32) {
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
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
