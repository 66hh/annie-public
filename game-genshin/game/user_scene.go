package game

import (
	"flswld.com/common/utils/object"
	"flswld.com/common/utils/random"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/constant"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

func (g *GameManager) EnterSceneReadyReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user enter scene ready, user id: %v", userId)

	// PacketEnterScenePeerNotify
	enterScenePeerNotify := new(proto.EnterScenePeerNotify)
	enterScenePeerNotify.DestSceneId = player.SceneId
	world := g.worldManager.GetWorldByID(player.WorldId)
	enterScenePeerNotify.PeerId = player.PeerId
	enterScenePeerNotify.HostPeerId = world.owner.PeerId
	enterScenePeerNotify.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(proto.ApiEnterScenePeerNotify, userId, player.ClientSeq, enterScenePeerNotify)

	// PacketEnterSceneReadyRsp
	enterSceneReadyRsp := new(proto.EnterSceneReadyRsp)
	enterSceneReadyRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(proto.ApiEnterSceneReadyRsp, userId, player.ClientSeq, enterSceneReadyRsp)
}

func (g *GameManager) SceneInitFinishReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user scene init finish, user id: %v", userId)

	// PacketServerTimeNotify
	serverTimeNotify := new(proto.ServerTimeNotify)
	serverTimeNotify.ServerTime = uint64(time.Now().UnixMilli())
	g.SendMsg(proto.ApiServerTimeNotify, userId, player.ClientSeq, serverTimeNotify)

	// PacketWorldPlayerInfoNotify
	worldPlayerInfoNotify := new(proto.WorldPlayerInfoNotify)
	world := g.worldManager.GetWorldByID(player.WorldId)
	scene := world.GetSceneById(player.SceneId)
	playerPropertyConst := constant.GetPlayerPropertyConst()
	for _, worldPlayer := range world.playerMap {
		onlinePlayerInfo := new(proto.OnlinePlayerInfo)
		onlinePlayerInfo.Uid = worldPlayer.PlayerID
		onlinePlayerInfo.Nickname = worldPlayer.NickName
		onlinePlayerInfo.PlayerLevel = worldPlayer.PropertiesMap[playerPropertyConst.PROP_PLAYER_LEVEL]
		onlinePlayerInfo.MpSettingType = proto.MpSettingType(worldPlayer.PropertiesMap[playerPropertyConst.PROP_PLAYER_MP_SETTING_TYPE])
		onlinePlayerInfo.NameCardId = worldPlayer.NameCard
		onlinePlayerInfo.Signature = worldPlayer.Signature
		onlinePlayerInfo.ProfilePicture = &proto.ProfilePicture{AvatarId: worldPlayer.HeadImage}
		onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(world.playerMap))
		worldPlayerInfoNotify.PlayerInfoList = append(worldPlayerInfoNotify.PlayerInfoList, onlinePlayerInfo)
		worldPlayerInfoNotify.PlayerUidList = append(worldPlayerInfoNotify.PlayerUidList, worldPlayer.PlayerID)
	}
	g.SendMsg(proto.ApiWorldPlayerInfoNotify, userId, player.ClientSeq, worldPlayerInfoNotify)

	// PacketWorldDataNotify
	worldDataNotify := new(proto.WorldDataNotify)
	worldDataNotify.WorldPropMap = make(map[uint32]*proto.PropValue)
	// 世界等级
	worldDataNotify.WorldPropMap[1] = &proto.PropValue{
		Type:  1,
		Val:   int64(world.worldLevel),
		Value: &proto.PropValue_Ival{Ival: int64(world.worldLevel)},
	}
	// 是否多人游戏
	worldDataNotify.WorldPropMap[2] = &proto.PropValue{
		Type:  2,
		Val:   object.ConvBoolToInt64(world.multiplayer),
		Value: &proto.PropValue_Ival{Ival: object.ConvBoolToInt64(world.multiplayer)},
	}
	g.SendMsg(proto.ApiWorldDataNotify, userId, player.ClientSeq, worldDataNotify)

	// PacketPlayerWorldSceneInfoListNotify
	playerWorldSceneInfoListNotify := new(proto.PlayerWorldSceneInfoListNotify)
	playerWorldSceneInfoListNotify.InfoList = []*proto.PlayerWorldSceneInfo{
		{SceneId: 1, IsLocked: false, SceneTagIdList: []uint32{}},
		{SceneId: 3, IsLocked: false, SceneTagIdList: []uint32{102, 113, 117}},
		{SceneId: 4, IsLocked: false, SceneTagIdList: []uint32{106, 109, 117}},
		{SceneId: 5, IsLocked: false, SceneTagIdList: []uint32{}},
		{SceneId: 6, IsLocked: false, SceneTagIdList: []uint32{}},
		{SceneId: 7, IsLocked: false, SceneTagIdList: []uint32{}},
	}
	xumi := &proto.PlayerWorldSceneInfo{
		SceneId:        9,
		IsLocked:       false,
		SceneTagIdList: []uint32{},
	}
	for i := 0; i < 3000; i++ {
		xumi.SceneTagIdList = append(xumi.SceneTagIdList, uint32(i))
	}
	playerWorldSceneInfoListNotify.InfoList = append(playerWorldSceneInfoListNotify.InfoList, xumi)
	g.SendMsg(proto.ApiPlayerWorldSceneInfoListNotify, userId, player.ClientSeq, playerWorldSceneInfoListNotify)

	// SceneForceUnlockNotify
	g.SendMsg(proto.ApiSceneForceUnlockNotify, userId, player.ClientSeq, new(proto.SceneForceUnlockNotify))

	// PacketHostPlayerNotify
	hostPlayerNotify := new(proto.HostPlayerNotify)
	hostPlayerNotify.HostUid = world.owner.PlayerID
	hostPlayerNotify.HostPeerId = world.owner.PeerId
	g.SendMsg(proto.ApiHostPlayerNotify, userId, player.ClientSeq, hostPlayerNotify)

	// PacketSceneTimeNotify
	sceneTimeNotify := new(proto.SceneTimeNotify)
	sceneTimeNotify.SceneId = player.SceneId
	sceneTimeNotify.SceneTime = 0
	g.SendMsg(proto.ApiSceneTimeNotify, userId, player.ClientSeq, sceneTimeNotify)

	// PacketPlayerGameTimeNotify
	playerGameTimeNotify := new(proto.PlayerGameTimeNotify)
	playerGameTimeNotify.GameTime = scene.time
	playerGameTimeNotify.Uid = player.PlayerID
	g.SendMsg(proto.ApiPlayerGameTimeNotify, userId, player.ClientSeq, playerGameTimeNotify)

	// PacketPlayerEnterSceneInfoNotify
	empty := new(proto.AbilitySyncStateInfo)
	playerEnterSceneInfoNotify := new(proto.PlayerEnterSceneInfoNotify)
	activeAvatarId := player.TeamConfig.GetActiveAvatarId()
	playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)
	playerEnterSceneInfoNotify.CurAvatarEntityId = playerTeamEntity.avatarEntityMap[activeAvatarId]
	playerEnterSceneInfoNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneInfoNotify.TeamEnterInfo = &proto.TeamEnterSceneInfo{
		TeamEntityId:        playerTeamEntity.teamEntityId,
		TeamAbilityInfo:     empty,
		AbilityControlBlock: new(proto.AbilityControlBlock),
	}
	playerEnterSceneInfoNotify.MpLevelEntityInfo = &proto.MPLevelEntityInfo{
		EntityId:        g.worldManager.GetWorldByID(player.WorldId).mpLevelEntityId,
		AuthorityPeerId: g.worldManager.GetWorldByID(player.WorldId).owner.PeerId,
		AbilityInfo:     empty,
	}
	activeTeam := player.TeamConfig.GetActiveTeam()
	for _, avatarId := range activeTeam.AvatarIdList {
		if avatarId == 0 {
			break
		}
		avatar := player.AvatarMap[avatarId]
		avatarEnterSceneInfo := new(proto.AvatarEnterSceneInfo)
		avatarEnterSceneInfo.AvatarGuid = avatar.Guid
		avatarEnterSceneInfo.AvatarEntityId = playerTeamEntity.avatarEntityMap[avatarId]
		avatarEnterSceneInfo.WeaponGuid = avatar.EquipWeapon.Guid
		avatarEnterSceneInfo.WeaponEntityId = playerTeamEntity.weaponEntityMap[avatar.EquipWeapon.WeaponId]
		avatarEnterSceneInfo.AvatarAbilityInfo = empty
		avatarEnterSceneInfo.WeaponAbilityInfo = empty
		playerEnterSceneInfoNotify.AvatarEnterInfo = append(playerEnterSceneInfoNotify.AvatarEnterInfo, avatarEnterSceneInfo)
	}
	g.SendMsg(proto.ApiPlayerEnterSceneInfoNotify, userId, player.ClientSeq, playerEnterSceneInfoNotify)

	// PacketSceneAreaWeatherNotify
	sceneAreaWeatherNotify := new(proto.SceneAreaWeatherNotify)
	sceneAreaWeatherNotify.WeatherAreaId = 0
	climateTypeConst := constant.GetClimateTypeConst()
	sceneAreaWeatherNotify.ClimateType = uint32(climateTypeConst.CLIMATE_SUNNY)
	g.SendMsg(proto.ApiSceneAreaWeatherNotify, userId, player.ClientSeq, sceneAreaWeatherNotify)

	// PacketScenePlayerInfoNotify
	scenePlayerInfoNotify := new(proto.ScenePlayerInfoNotify)
	for _, worldPlayer := range world.playerMap {
		onlinePlayerInfo := new(proto.OnlinePlayerInfo)
		onlinePlayerInfo.Uid = worldPlayer.PlayerID
		onlinePlayerInfo.Nickname = worldPlayer.NickName
		onlinePlayerInfo.PlayerLevel = worldPlayer.PropertiesMap[playerPropertyConst.PROP_PLAYER_LEVEL]
		onlinePlayerInfo.MpSettingType = proto.MpSettingType(worldPlayer.PropertiesMap[playerPropertyConst.PROP_PLAYER_MP_SETTING_TYPE])
		onlinePlayerInfo.NameCardId = worldPlayer.NameCard
		onlinePlayerInfo.Signature = worldPlayer.Signature
		onlinePlayerInfo.ProfilePicture = &proto.ProfilePicture{AvatarId: worldPlayer.HeadImage}
		onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(world.playerMap))
		scenePlayerInfoNotify.PlayerInfoList = append(scenePlayerInfoNotify.PlayerInfoList, &proto.ScenePlayerInfo{
			Uid:              worldPlayer.PlayerID,
			PeerId:           worldPlayer.PeerId,
			Name:             worldPlayer.NickName,
			SceneId:          worldPlayer.SceneId,
			OnlinePlayerInfo: onlinePlayerInfo,
		})
	}
	g.SendMsg(proto.ApiScenePlayerInfoNotify, userId, player.ClientSeq, scenePlayerInfoNotify)

	// PacketSceneTeamUpdateNotify
	sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world)
	g.SendMsg(proto.ApiSceneTeamUpdateNotify, userId, player.ClientSeq, sceneTeamUpdateNotify)

	// PacketSyncTeamEntityNotify
	syncTeamEntityNotify := new(proto.SyncTeamEntityNotify)
	syncTeamEntityNotify.SceneId = player.SceneId
	syncTeamEntityNotify.TeamEntityInfoList = make([]*proto.TeamEntityInfo, 0)
	if world.multiplayer {
		for _, worldPlayer := range world.playerMap {
			if worldPlayer.PlayerID == player.PlayerID {
				continue
			}
			worldPlayerScene := world.GetSceneById(worldPlayer.SceneId)
			worldPlayerTeamEntity := worldPlayerScene.GetPlayerTeamEntity(worldPlayer.PlayerID)
			teamEntityInfo := &proto.TeamEntityInfo{
				TeamEntityId:    worldPlayerTeamEntity.teamEntityId,
				AuthorityPeerId: worldPlayer.PeerId,
				TeamAbilityInfo: new(proto.AbilitySyncStateInfo),
			}
			syncTeamEntityNotify.TeamEntityInfoList = append(syncTeamEntityNotify.TeamEntityInfoList, teamEntityInfo)
		}
	}
	g.SendMsg(proto.ApiSyncTeamEntityNotify, userId, player.ClientSeq, syncTeamEntityNotify)

	// PacketSyncScenePlayTeamEntityNotify
	syncScenePlayTeamEntityNotify := new(proto.SyncScenePlayTeamEntityNotify)
	syncScenePlayTeamEntityNotify.SceneId = player.SceneId
	g.SendMsg(proto.ApiSyncScenePlayTeamEntityNotify, userId, player.ClientSeq, syncScenePlayTeamEntityNotify)

	// PacketSceneInitFinishRsp
	SceneInitFinishRsp := new(proto.SceneInitFinishRsp)
	SceneInitFinishRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(proto.ApiSceneInitFinishRsp, userId, player.ClientSeq, SceneInitFinishRsp)

	player.SceneLoadState = model.SceneInitFinish
}

func (g *GameManager) EnterSceneDoneReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user enter scene done, user id: %v", userId)

	// PacketEnterSceneDoneRsp
	enterSceneDoneRsp := new(proto.EnterSceneDoneRsp)
	enterSceneDoneRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(proto.ApiEnterSceneDoneRsp, userId, player.ClientSeq, enterSceneDoneRsp)

	// PacketPlayerTimeNotify
	playerTimeNotify := new(proto.PlayerTimeNotify)
	playerTimeNotify.IsPaused = player.Pause
	playerTimeNotify.PlayerTime = uint64(player.TotalOnlineTime)
	playerTimeNotify.ServerTime = uint64(time.Now().UnixMilli())
	g.SendMsg(proto.ApiPlayerTimeNotify, userId, player.ClientSeq, playerTimeNotify)

	player.SceneLoadState = model.SceneEnterDone
	g.AddSceneEntityAvatarBroadcastNotify(player)
	g.MeetSceneEntityNotify(player)
}

func (g *GameManager) PostEnterSceneReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user post enter scene, user id: %v", userId)

	// PacketPostEnterSceneRsp
	postEnterSceneRsp := new(proto.PostEnterSceneRsp)
	postEnterSceneRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(proto.ApiPostEnterSceneRsp, userId, player.ClientSeq, postEnterSceneRsp)
}

func (g *GameManager) EnterWorldAreaReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user enter world area, user id: %v", userId)
	req := payloadMsg.(*proto.EnterWorldAreaReq)

	// PacketEnterWorldAreaRsp
	enterWorldAreaRsp := new(proto.EnterWorldAreaRsp)
	enterWorldAreaRsp.AreaType = req.AreaType
	enterWorldAreaRsp.AreaId = req.AreaId
	g.SendMsg(proto.ApiEnterWorldAreaRsp, userId, player.ClientSeq, enterWorldAreaRsp)
}

func (g *GameManager) ChangeGameTimeReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user change game time, user id: %v", userId)
	req := payloadMsg.(*proto.ChangeGameTimeReq)
	gameTime := req.GameTime
	world := g.worldManager.GetWorldByID(player.WorldId)
	scene := world.GetSceneById(player.SceneId)
	scene.ChangeTime(gameTime)

	for _, scenePlayer := range scene.playerMap {
		// PacketPlayerGameTimeNotify
		playerGameTimeNotify := new(proto.PlayerGameTimeNotify)
		playerGameTimeNotify.GameTime = scene.time
		playerGameTimeNotify.Uid = scenePlayer.PlayerID
		g.SendMsg(proto.ApiPlayerGameTimeNotify, scenePlayer.PlayerID, scenePlayer.ClientSeq, playerGameTimeNotify)
	}

	// PacketChangeGameTimeRsp
	changeGameTimeRsp := new(proto.ChangeGameTimeRsp)
	changeGameTimeRsp.CurGameTime = scene.time
	g.SendMsg(proto.ApiChangeGameTimeRsp, userId, player.ClientSeq, changeGameTimeRsp)
}

func (g *GameManager) PacketPlayerEnterSceneNotify(player *model.Player) *proto.PlayerEnterSceneNotify {
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
	playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
	playerEnterSceneNotify.SceneId = player.SceneId
	playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
	playerEnterSceneNotify.Type = proto.EnterType_ENTER_TYPE_SELF
	playerEnterSceneNotify.TargetUid = player.PlayerID
	playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneNotify.WorldLevel = player.PropertiesMap[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.Login)
	// 刚登录进入场景的时候才为true
	playerEnterSceneNotify.IsFirstLoginEnterScene = true
	playerEnterSceneNotify.WorldType = 1
	playerEnterSceneNotify.SceneTransaction = strconv.Itoa(int(player.SceneId)) + "-" +
		strconv.Itoa(int(player.PlayerID)) + "-" +
		strconv.Itoa(int(time.Now().Unix())) + "-" +
		"18402"
	return playerEnterSceneNotify
}

func (g *GameManager) PacketPlayerEnterSceneNotifyTp(
	player *model.Player,
	enterType proto.EnterType,
	enterReason uint32,
	prevSceneId uint32,
	prevPos *model.Vector,
) *proto.PlayerEnterSceneNotify {
	return g.PacketPlayerEnterSceneNotifyMp(player, player, enterType, enterReason, prevSceneId, prevPos)
}

func (g *GameManager) PacketPlayerEnterSceneNotifyMp(
	player *model.Player,
	targetPlayer *model.Player,
	enterType proto.EnterType,
	enterReason uint32,
	prevSceneId uint32,
	prevPos *model.Vector,
) *proto.PlayerEnterSceneNotify {
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
	playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
	playerEnterSceneNotify.PrevSceneId = prevSceneId
	playerEnterSceneNotify.PrevPos = &proto.Vector{X: float32(prevPos.X), Y: float32(prevPos.Y), Z: float32(prevPos.Z)}
	playerEnterSceneNotify.SceneId = player.SceneId
	playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
	playerEnterSceneNotify.Type = enterType
	playerEnterSceneNotify.TargetUid = targetPlayer.PlayerID
	playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneNotify.WorldLevel = targetPlayer.PropertiesMap[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
	playerEnterSceneNotify.EnterReason = enterReason
	playerEnterSceneNotify.WorldType = 1
	playerEnterSceneNotify.SceneTransaction = strconv.Itoa(int(player.SceneId)) + "-" +
		strconv.Itoa(int(targetPlayer.PlayerID)) + "-" +
		strconv.Itoa(int(time.Now().Unix())) + "-" +
		"18402"

	//playerEnterSceneNotify.SceneTagIdList = []uint32{102, 107, 109, 113, 117}
	playerEnterSceneNotify.SceneTagIdList = make([]uint32, 0)
	for sceneTagId := uint32(0); sceneTagId < 3000; sceneTagId++ {
		playerEnterSceneNotify.SceneTagIdList = append(playerEnterSceneNotify.SceneTagIdList, sceneTagId)
	}

	return playerEnterSceneNotify
}

func (g *GameManager) AddSceneEntityAvatarBroadcastNotify(player *model.Player) {
	world := g.worldManager.GetWorldByID(player.WorldId)
	scene := world.GetSceneById(player.SceneId)

	// PacketSceneEntityAppearNotify
	sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
	sceneEntityAppearNotify.AppearType = proto.VisionType_VISION_TYPE_BORN
	sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{g.PacketSceneEntityInfoAvatar(scene, player, player.TeamConfig.GetActiveAvatarId())}
	for _, scenePlayer := range scene.playerMap {
		g.SendMsg(proto.ApiSceneEntityAppearNotify, scenePlayer.PlayerID, scenePlayer.ClientSeq, sceneEntityAppearNotify)
		logger.LOG.Debug("SceneEntityAppearNotify, uid: %v, data: %v", scenePlayer.PlayerID, sceneEntityAppearNotify)
	}
}

func (g *GameManager) RemoveSceneEntityAvatarBroadcastNotify(player *model.Player) {
	world := g.worldManager.GetWorldByID(player.WorldId)
	scene := world.GetSceneById(player.SceneId)

	// PacketSceneEntityDisappearNotify
	sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
	activeAvatarId := player.TeamConfig.GetActiveAvatarId()
	playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)
	sceneEntityDisappearNotify.EntityList = []uint32{playerTeamEntity.avatarEntityMap[activeAvatarId]}
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REMOVE
	for _, scenePlayer := range scene.playerMap {
		g.SendMsg(proto.ApiSceneEntityDisappearNotify, scenePlayer.PlayerID, scenePlayer.ClientSeq, sceneEntityDisappearNotify)
		logger.LOG.Debug("SceneEntityDisappearNotify, uid: %v, data: %v", scenePlayer.PlayerID, sceneEntityDisappearNotify)
	}
}

func (g *GameManager) MeetSceneEntityNotify(player *model.Player) {
	world := g.worldManager.GetWorldByID(player.WorldId)
	scene := world.GetSceneById(player.SceneId)

	// PacketSceneEntityAppearNotify
	sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
	sceneEntityAppearNotify.AppearType = proto.VisionType_VISION_TYPE_MEET
	sceneEntityAppearNotify.EntityList = make([]*proto.SceneEntityInfo, 0)

	for _, entity := range scene.entityMap {
		switch entity.entityType {
		case uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_AVATAR):
			if entity.uid == player.PlayerID {
				continue
			}
			scenePlayer := g.userManager.GetOnlineUser(entity.uid)
			if scenePlayer == nil {
				logger.LOG.Error("get scene player is nil, world id: %v, scene id: %v", world.id, scene.id)
				continue
			}
			if scenePlayer.SceneLoadState != model.SceneEnterDone {
				continue
			}
			if entity.avatarId != scenePlayer.TeamConfig.GetActiveAvatarId() {
				continue
			}
			sceneEntityInfoAvatar := g.PacketSceneEntityInfoAvatar(scene, scenePlayer, scenePlayer.TeamConfig.GetActiveAvatarId())
			sceneEntityAppearNotify.EntityList = append(sceneEntityAppearNotify.EntityList, sceneEntityInfoAvatar)
		case uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_WEAPON):
		case uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_MONSTER):
			sceneEntityInfoMonster := g.PacketSceneEntityInfoMonster(scene, entity.id)
			sceneEntityAppearNotify.EntityList = append(sceneEntityAppearNotify.EntityList, sceneEntityInfoMonster)
		}
	}

	g.SendMsg(proto.ApiSceneEntityAppearNotify, player.PlayerID, player.ClientSeq, sceneEntityAppearNotify)
	logger.LOG.Debug("SceneEntityAppearNotify, uid: %v, data: %v", player.PlayerID, sceneEntityAppearNotify)
}

func (g *GameManager) PacketSceneEntityInfoAvatar(scene *Scene, player *model.Player, avatarId uint32) *proto.SceneEntityInfo {
	playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)
	entity := scene.GetEntity(playerTeamEntity.avatarEntityMap[avatarId])
	playerPropertyConst := constant.GetPlayerPropertyConst()
	fightPropertyConst := constant.GetFightPropertyConst()
	sceneEntityInfo := &proto.SceneEntityInfo{
		EntityType: proto.ProtEntityType_PROT_ENTITY_TYPE_AVATAR,
		EntityId:   entity.id,
		MotionInfo: &proto.MotionInfo{
			Pos: &proto.Vector{
				X: float32(entity.pos.X),
				Y: float32(entity.pos.Y),
				Z: float32(entity.pos.Z),
			},
			Rot: &proto.Vector{
				X: float32(entity.rot.X),
				Y: float32(entity.rot.Y),
				Z: float32(entity.rot.Z),
			},
			Speed: &proto.Vector{},
			State: proto.MotionState(entity.moveState),
		},
		PropList: []*proto.PropPair{{Type: uint32(playerPropertyConst.PROP_LEVEL), PropValue: &proto.PropValue{
			Type:  uint32(playerPropertyConst.PROP_LEVEL),
			Value: &proto.PropValue_Ival{Ival: int64(entity.level)},
			Val:   int64(entity.level),
		}}},
		FightPropList: []*proto.FightPropPair{
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_HP),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_HP),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_MAX_HP),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)],
			},
		},
		LifeState:        1,
		AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
		Entity: &proto.SceneEntityInfo_Avatar{
			Avatar: g.PacketSceneAvatarInfo(scene, player, avatarId),
		},
		EntityClientData: new(proto.EntityClientData),
		EntityAuthorityInfo: &proto.EntityAuthorityInfo{
			AbilityInfo:         new(proto.AbilitySyncStateInfo),
			RendererChangedInfo: new(proto.EntityRendererChangedInfo),
			AiInfo: &proto.SceneEntityAiInfo{
				IsAiOpen: true,
				BornPos:  new(proto.Vector),
			},
			BornPos: new(proto.Vector),
		},
		LastMoveSceneTimeMs: entity.lastMoveSceneTimeMs,
		LastMoveReliableSeq: entity.lastMoveReliableSeq,
	}
	return sceneEntityInfo
}

func (g *GameManager) PacketSceneEntityInfoMonster(scene *Scene, entityId uint32) *proto.SceneEntityInfo {
	entity := scene.GetEntity(entityId)
	playerPropertyConst := constant.GetPlayerPropertyConst()
	fightPropertyConst := constant.GetFightPropertyConst()
	pos := &proto.Vector{
		X: float32(entity.pos.X),
		Y: float32(entity.pos.Y),
		Z: float32(entity.pos.Z),
	}
	sceneEntityInfo := &proto.SceneEntityInfo{
		EntityType: proto.ProtEntityType_PROT_ENTITY_TYPE_MONSTER,
		EntityId:   entity.id,
		MotionInfo: &proto.MotionInfo{
			Pos: pos,
			Rot: &proto.Vector{
				X: float32(entity.rot.X),
				Y: float32(entity.rot.Y),
				Z: float32(entity.rot.Z),
			},
			Speed: &proto.Vector{},
			State: proto.MotionState(entity.moveState),
		},
		PropList: []*proto.PropPair{{Type: uint32(playerPropertyConst.PROP_LEVEL), PropValue: &proto.PropValue{
			Type:  uint32(playerPropertyConst.PROP_LEVEL),
			Value: &proto.PropValue_Ival{Ival: int64(entity.level)},
			Val:   int64(entity.level),
		}}},
		FightPropList: []*proto.FightPropPair{
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_HP),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_HP),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_MAX_HP),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE),
				PropValue: entity.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)],
			},
		},
		LifeState:        1,
		AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
		Entity: &proto.SceneEntityInfo_Monster{
			Monster: g.PacketSceneMonsterInfo(),
		},
		EntityClientData: new(proto.EntityClientData),
		EntityAuthorityInfo: &proto.EntityAuthorityInfo{
			AbilityInfo:         new(proto.AbilitySyncStateInfo),
			RendererChangedInfo: new(proto.EntityRendererChangedInfo),
			AiInfo: &proto.SceneEntityAiInfo{
				IsAiOpen: true,
				BornPos:  pos,
			},
			BornPos: pos,
		},
	}
	return sceneEntityInfo
}

func (g *GameManager) PacketSceneAvatarInfo(scene *Scene, player *model.Player, avatarId uint32) *proto.SceneAvatarInfo {
	activeAvatarId := player.TeamConfig.GetActiveAvatarId()
	activeAvatar := player.AvatarMap[activeAvatarId]
	playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)
	equipIdList := make([]uint32, 0)
	weapon := player.AvatarMap[avatarId].EquipWeapon
	equipIdList = append(equipIdList, weapon.ItemId)
	for _, reliquary := range player.AvatarMap[avatarId].EquipReliquaryList {
		equipIdList = append(equipIdList, reliquary.ItemId)
	}
	sceneAvatarInfo := &proto.SceneAvatarInfo{
		Uid:          player.PlayerID,
		AvatarId:     avatarId,
		Guid:         player.AvatarMap[avatarId].Guid,
		PeerId:       player.PeerId,
		EquipIdList:  equipIdList,
		SkillDepotId: player.AvatarMap[avatarId].SkillDepotId,
		Weapon: &proto.SceneWeaponInfo{
			EntityId:    playerTeamEntity.weaponEntityMap[activeAvatar.EquipWeapon.WeaponId],
			GadgetId:    uint32(gdc.CONF.ItemDataMap[int32(weapon.ItemId)].GadgetId),
			ItemId:      weapon.ItemId,
			Guid:        weapon.Guid,
			Level:       uint32(weapon.Level),
			AbilityInfo: new(proto.AbilitySyncStateInfo),
		},
		ReliquaryList:     nil,
		SkillLevelMap:     player.AvatarMap[avatarId].SkillLevelMap,
		WearingFlycloakId: player.AvatarMap[avatarId].FlyCloak,
		CostumeId:         player.AvatarMap[avatarId].Costume,
		BornTime:          uint32(player.AvatarMap[avatarId].BornTime),
		TeamResonanceList: make([]uint32, 0),
	}
	for id := range player.TeamConfig.TeamResonances {
		sceneAvatarInfo.TeamResonanceList = append(sceneAvatarInfo.TeamResonanceList, uint32(id))
	}
	return sceneAvatarInfo
}

func (g *GameManager) PacketSceneMonsterInfo() *proto.SceneMonsterInfo {
	sceneMonsterInfo := &proto.SceneMonsterInfo{
		MonsterId:       21010101,
		AuthorityPeerId: 1,
		BornType:        proto.MonsterBornType_MONSTER_BORN_TYPE_DEFAULT,
		BlockId:         3001,
		TitleId:         3001,
		SpecialNameId:   40,
	}
	return sceneMonsterInfo
}

func (g *GameManager) PacketDelTeamEntityNotify(scene *Scene, player *model.Player) *proto.DelTeamEntityNotify {
	delTeamEntityNotify := new(proto.DelTeamEntityNotify)
	delTeamEntityNotify.SceneId = player.SceneId
	playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)
	delTeamEntityNotify.DelEntityIdList = []uint32{playerTeamEntity.teamEntityId}
	return delTeamEntityNotify
}
