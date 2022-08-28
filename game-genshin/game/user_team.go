package game

import (
	"flswld.com/common/utils/endec"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/constant"
	pb "google.golang.org/protobuf/proto"
)

func (g *GameManager) ChangeAvatarReq(userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message) {
	logger.LOG.Debug("user change avatar, user id: %v", userId)
	req := payloadMsg.(*proto.ChangeAvatarReq)
	targetAvatarGuid := req.Guid

	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	world := g.worldManager.GetWorldByID(player.WorldId)
	scene := world.GetSceneById(player.SceneId)
	playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)

	oldAvatarId := player.TeamConfig.GetActiveAvatarId()
	oldAvatar := player.AvatarMap[oldAvatarId]
	if oldAvatar.Guid == targetAvatarGuid {
		logger.LOG.Error("can not change to the same avatar, user id: %v, oldAvatarId: %v, oldAvatarGuid: %v", userId, oldAvatarId, oldAvatar.Guid)
		return
	}
	activeTeam := player.TeamConfig.GetActiveTeam()
	index := -1
	for avatarIndex, avatarId := range activeTeam.AvatarIdList {
		if avatarId == 0 {
			break
		}
		if targetAvatarGuid == player.AvatarMap[avatarId].Guid {
			index = avatarIndex
		}
	}
	if index == -1 {
		logger.LOG.Error("can not find the target avatar in team, user id: %v, target avatar guid: %v", userId, targetAvatarGuid)
		return
	}
	player.TeamConfig.CurrAvatarIndex = uint8(index)

	entity := scene.GetEntity(playerTeamEntity.avatarEntityMap[oldAvatarId])
	entity.moveState = uint16(proto.MotionState_MOTION_STATE_STANDBY)

	// TODO 目前多人游戏可能会存在问题 可能需要将原来的队伍里的角色实体放到世界里去才行 只是可能而已 待验证

	// PacketSceneEntityDisappearNotify
	sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REPLACE
	sceneEntityDisappearNotify.EntityList = []uint32{playerTeamEntity.avatarEntityMap[oldAvatarId]}
	g.SendMsg(proto.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)

	// PacketSceneEntityAppearNotify
	sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REPLACE
	sceneEntityAppearNotify.Param = playerTeamEntity.avatarEntityMap[oldAvatarId]
	sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{g.PacketSceneEntityInfoAvatar(scene, player, player.TeamConfig.GetActiveAvatarId())}
	g.SendMsg(proto.ApiSceneEntityAppearNotify, userId, nil, sceneEntityAppearNotify)

	// PacketChangeAvatarRsp
	changeAvatarRsp := new(proto.ChangeAvatarRsp)
	changeAvatarRsp.Retcode = int32(proto.Retcode_RETCODE_RET_SUCC)
	changeAvatarRsp.CurGuid = targetAvatarGuid
	g.SendMsg(proto.ApiChangeAvatarRsp, userId, nil, changeAvatarRsp)
}

func (g *GameManager) SetUpAvatarTeamReq(userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message) {
	logger.LOG.Debug("user set up avatar team, user id: %v", userId)
	req := payloadMsg.(*proto.SetUpAvatarTeamReq)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	teamId := req.TeamId
	avatarGuidList := req.AvatarTeamGuidList
	world := g.worldManager.GetWorldByID(player.WorldId)
	selfTeam := teamId == uint32(player.TeamConfig.GetActiveTeamId())
	if (selfTeam && len(avatarGuidList) == 0) || len(avatarGuidList) > 4 || world.multiplayer {
		return
	}
	avatarIdList := make([]uint32, 0)
	for _, avatarGuid := range avatarGuidList {
		for avatarId, avatar := range player.AvatarMap {
			if avatarGuid == avatar.Guid {
				avatarIdList = append(avatarIdList, avatarId)
			}
		}
	}
	player.TeamConfig.ClearTeamAvatar(uint8(teamId - 1))
	for _, avatarId := range avatarIdList {
		player.TeamConfig.AddAvatarToTeam(avatarId, uint8(teamId-1))
	}
	if world.multiplayer {
		// TODO 多人世界队伍
	} else {
		// PacketAvatarTeamUpdateNotify
		avatarTeamUpdateNotify := new(proto.AvatarTeamUpdateNotify)
		avatarTeamUpdateNotify.AvatarTeamMap = make(map[uint32]*proto.AvatarTeam)
		for teamIndex, team := range player.TeamConfig.TeamList {
			avatarTeam := new(proto.AvatarTeam)
			avatarTeam.TeamName = team.Name
			for _, avatarId := range team.AvatarIdList {
				if avatarId == 0 {
					break
				}
				avatarTeam.AvatarGuidList = append(avatarTeam.AvatarGuidList, player.AvatarMap[avatarId].Guid)
			}
			avatarTeamUpdateNotify.AvatarTeamMap[uint32(teamIndex)+1] = avatarTeam
		}
		g.SendMsg(proto.ApiAvatarTeamUpdateNotify, userId, nil, avatarTeamUpdateNotify)

		if selfTeam {
			player.TeamConfig.CurrAvatarIndex = 0
			player.TeamConfig.UpdateTeam()
			scene := world.GetSceneById(player.SceneId)
			scene.UpdatePlayerTeamEntity(player)
			// TODO 还有一大堆没写 SceneTeamUpdateNotify
			// PacketSceneTeamUpdateNotify
			sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world)
			g.SendMsg(proto.ApiSceneTeamUpdateNotify, userId, nil, sceneTeamUpdateNotify)

			// PacketSetUpAvatarTeamRsp
			setUpAvatarTeamRsp := new(proto.SetUpAvatarTeamRsp)
			setUpAvatarTeamRsp.TeamId = teamId
			setUpAvatarTeamRsp.CurAvatarGuid = player.AvatarMap[player.TeamConfig.GetActiveAvatarId()].Guid
			team := player.TeamConfig.GetTeamByIndex(uint8(teamId - 1))
			for _, avatarId := range team.AvatarIdList {
				if avatarId == 0 {
					break
				}
				setUpAvatarTeamRsp.AvatarTeamGuidList = append(setUpAvatarTeamRsp.AvatarTeamGuidList, player.AvatarMap[avatarId].Guid)
			}
			g.SendMsg(proto.ApiSetUpAvatarTeamRsp, userId, nil, setUpAvatarTeamRsp)
		} else {
			// PacketSetUpAvatarTeamRsp
			setUpAvatarTeamRsp := new(proto.SetUpAvatarTeamRsp)
			setUpAvatarTeamRsp.TeamId = teamId
			setUpAvatarTeamRsp.CurAvatarGuid = player.AvatarMap[player.TeamConfig.GetActiveAvatarId()].Guid
			team := player.TeamConfig.GetTeamByIndex(uint8(teamId - 1))
			for _, avatarId := range team.AvatarIdList {
				if avatarId == 0 {
					break
				}
				setUpAvatarTeamRsp.AvatarTeamGuidList = append(setUpAvatarTeamRsp.AvatarTeamGuidList, player.AvatarMap[avatarId].Guid)
			}
			g.SendMsg(proto.ApiSetUpAvatarTeamRsp, userId, nil, setUpAvatarTeamRsp)
		}
	}
}

func (g *GameManager) ChooseCurAvatarTeamReq(userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message) {
	logger.LOG.Debug("user switch team, user id: %v", userId)
	req := payloadMsg.(*proto.ChooseCurAvatarTeamReq)
	teamId := req.TeamId
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world.multiplayer {
		return
	}
	team := player.TeamConfig.GetTeamByIndex(uint8(teamId) - 1)
	if team == nil || len(team.AvatarIdList) == 0 {
		return
	}
	player.TeamConfig.CurrTeamIndex = uint8(teamId) - 1
	player.TeamConfig.CurrAvatarIndex = 0
	player.TeamConfig.UpdateTeam()
	scene := world.GetSceneById(player.SceneId)
	scene.UpdatePlayerTeamEntity(player)

	// TODO 还有一大堆没写 SceneTeamUpdateNotify
	// PacketSceneTeamUpdateNotify
	sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world)
	g.SendMsg(proto.ApiSceneTeamUpdateNotify, userId, nil, sceneTeamUpdateNotify)

	// PacketChooseCurAvatarTeamRsp
	chooseCurAvatarTeamRsp := new(proto.ChooseCurAvatarTeamRsp)
	chooseCurAvatarTeamRsp.CurTeamId = teamId
	g.SendMsg(proto.ApiChooseCurAvatarTeamRsp, userId, nil, chooseCurAvatarTeamRsp)
}

func (g *GameManager) PacketSceneTeamUpdateNotify(world *World) *proto.SceneTeamUpdateNotify {
	sceneTeamUpdateNotify := new(proto.SceneTeamUpdateNotify)
	sceneTeamUpdateNotify.IsInMp = world.multiplayer
	empty := new(proto.AbilitySyncStateInfo)
	for _, worldPlayer := range world.playerMap {
		worldPlayerScene := world.GetSceneById(worldPlayer.SceneId)
		worldPlayerTeamEntity := worldPlayerScene.GetPlayerTeamEntity(worldPlayer.PlayerID)
		team := worldPlayer.TeamConfig.GetActiveTeam()
		for _, avatarId := range team.AvatarIdList {
			if avatarId == 0 {
				break
			}
			worldPlayerAvatar := worldPlayer.AvatarMap[avatarId]
			equipIdList := make([]uint32, 0)
			weapon := worldPlayerAvatar.EquipWeapon
			equipIdList = append(equipIdList, weapon.ItemId)
			for _, reliquary := range worldPlayerAvatar.EquipReliquaryList {
				equipIdList = append(equipIdList, reliquary.ItemId)
			}
			sceneTeamAvatar := &proto.SceneTeamAvatar{
				PlayerUid:           worldPlayer.PlayerID,
				AvatarGuid:          worldPlayerAvatar.Guid,
				SceneId:             worldPlayer.SceneId,
				EntityId:            worldPlayerTeamEntity.avatarEntityMap[avatarId],
				SceneEntityInfo:     g.PacketSceneEntityInfoAvatar(worldPlayerScene, worldPlayer, avatarId),
				WeaponGuid:          worldPlayerAvatar.EquipWeapon.Guid,
				WeaponEntityId:      worldPlayerTeamEntity.weaponEntityMap[worldPlayerAvatar.EquipWeapon.WeaponId],
				IsPlayerCurAvatar:   worldPlayer.TeamConfig.GetActiveAvatarId() == avatarId,
				IsOnScene:           worldPlayer.TeamConfig.GetActiveAvatarId() == avatarId,
				AvatarAbilityInfo:   empty,
				WeaponAbilityInfo:   empty,
				AbilityControlBlock: new(proto.AbilityControlBlock),
			}
			if world.multiplayer {
				sceneTeamAvatar.AvatarInfo = g.PacketAvatarInfo(worldPlayerAvatar)
				sceneTeamAvatar.SceneAvatarInfo = g.PacketSceneAvatarInfo(worldPlayerScene, worldPlayer, avatarId)
			}
			// add AbilityControlBlock
			avatarDataConfig := gdc.CONF.AvatarDataMap[int32(avatarId)]
			acb := sceneTeamAvatar.AbilityControlBlock
			embryoId := 0
			gameConstant := constant.GetGameConstant()
			// add avatar abilities
			for _, abilityId := range avatarDataConfig.Abilities {
				embryoId++
				emb := &proto.AbilityEmbryo{
					AbilityId:               uint32(embryoId),
					AbilityNameHash:         uint32(abilityId),
					AbilityOverrideNameHash: uint32(gameConstant.DEFAULT_ABILITY_NAME),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, emb)
			}
			// add default abilities
			for _, abilityId := range gameConstant.DEFAULT_ABILITY_HASHES {
				embryoId++
				emb := &proto.AbilityEmbryo{
					AbilityId:               uint32(embryoId),
					AbilityNameHash:         uint32(abilityId),
					AbilityOverrideNameHash: uint32(gameConstant.DEFAULT_ABILITY_NAME),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, emb)
			}
			// add team resonances
			for id := range worldPlayer.TeamConfig.TeamResonancesConfig {
				embryoId++
				emb := &proto.AbilityEmbryo{
					AbilityId:               uint32(embryoId),
					AbilityNameHash:         uint32(id),
					AbilityOverrideNameHash: uint32(gameConstant.DEFAULT_ABILITY_NAME),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, emb)
			}
			// add skill depot abilities
			skillDepot := gdc.CONF.AvatarSkillDepotDataMap[int32(worldPlayerAvatar.SkillDepotId)]
			if skillDepot != nil && len(skillDepot.Abilities) != 0 {
				for _, id := range skillDepot.Abilities {
					embryoId++
					emb := &proto.AbilityEmbryo{
						AbilityId:               uint32(embryoId),
						AbilityNameHash:         uint32(id),
						AbilityOverrideNameHash: uint32(gameConstant.DEFAULT_ABILITY_NAME),
					}
					acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, emb)
				}
			}
			// add equip abilities
			for skill := range worldPlayerAvatar.ExtraAbilityEmbryos {
				embryoId++
				emb := &proto.AbilityEmbryo{
					AbilityId:               uint32(embryoId),
					AbilityNameHash:         uint32(endec.GenshinAbilityHashCode(skill)),
					AbilityOverrideNameHash: uint32(gameConstant.DEFAULT_ABILITY_NAME),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, emb)
			}
			sceneTeamUpdateNotify.SceneTeamAvatarList = append(sceneTeamUpdateNotify.SceneTeamAvatarList, sceneTeamAvatar)
		}
	}
	return sceneTeamUpdateNotify
}
