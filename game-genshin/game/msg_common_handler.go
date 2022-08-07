package game

import (
	"flswld.com/common/utils/endec"
	"flswld.com/common/utils/object"
	"flswld.com/common/utils/random"
	"flswld.com/common/utils/reflection"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"game-genshin/game/constant"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

func (g *GameManager) PlayerSetPauseReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Debug("user pause, user id: %v", userId)
	if headMsg != nil {
		g.log.Debug("client sequence id: %v", headMsg.ClientSequenceId)
	}
	if payloadMsg != nil {
		req := payloadMsg.(*proto.PlayerSetPauseReq)
		g.log.Debug("is paused: %v", req.IsPaused)
	}
}

func (g *GameManager) SetPlayerBornDataReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user set born data, user id: %v", userId)
	if headMsg != nil {
		g.log.Debug("client sequence id: %v", headMsg.ClientSequenceId)
	}
	user := g.userManager.GetTargetUser(userId)
	if user != nil {
		g.log.Error("recv set born data req, but user is already exist, userId: %v", userId)
		return
	}
	if payloadMsg != nil {
		req := payloadMsg.(*proto.SetPlayerBornDataReq)
		g.log.Debug("avatar id: %v, nickname: %v", req.AvatarId, req.NickName)

		mainCharAvatarId := req.GetAvatarId()
		if mainCharAvatarId != 10000005 && mainCharAvatarId != 10000007 {
			g.log.Error("invalid main char avatar id: %v", mainCharAvatarId)
			return
		}

		player := new(model.Player)
		player.PlayerID = userId
		player.NickName = req.NickName

		player.RegionId = 1
		player.SceneId = 3

		player.Properties = make(map[uint16]uint32)
		playerPropertyConst := constant.GetPlayerPropertyConst()
		// 初始化所有属性
		propList := reflection.ConvStructToMap(playerPropertyConst)
		for fieldName, fieldValue := range propList {
			if fieldName == "PROP_EXP" ||
				fieldName == "PROP_BREAK_LEVEL" ||
				fieldName == "PROP_SATIATION_VAL" ||
				fieldName == "PROP_SATIATION_PENALTY_TIME" ||
				fieldName == "PROP_LEVEL" {
				continue
			}
			value := fieldValue.(uint16)
			player.Properties[value] = 0
		}
		player.Properties[playerPropertyConst.PROP_PLAYER_LEVEL] = 60
		player.Properties[playerPropertyConst.PROP_IS_SPRING_AUTO_USE] = 1
		player.Properties[playerPropertyConst.PROP_SPRING_AUTO_USE_PERCENT] = 50
		player.Properties[playerPropertyConst.PROP_IS_FLYABLE] = 1
		player.Properties[playerPropertyConst.PROP_IS_TRANSFERABLE] = 1
		player.Properties[playerPropertyConst.PROP_MAX_STAMINA] = 24000
		player.Properties[playerPropertyConst.PROP_CUR_PERSIST_STAMINA] = 24000
		player.Properties[playerPropertyConst.PROP_PLAYER_RESIN] = 160

		player.FlyCloakList = make([]uint32, 0)
		player.FlyCloakList = append(player.FlyCloakList, 140001)

		player.CostumeList = make([]uint32, 0)

		player.Pos = &model.Vector{X: 2747, Y: 194, Z: -1719}
		player.Rotation = &model.Vector{X: 0, Y: 307, Z: 0}

		player.MpSetting = proto.MpSettingType_MP_SETTING_TYPE_ENTER_AFTER_APPLY

		player.ItemMap = make(map[uint32]*model.Item)
		player.WeaponMap = make(map[uint64]*model.Weapon)
		player.ReliquaryMap = make(map[uint64]*model.Reliquary)
		player.AvatarMap = make(map[uint32]*model.Avatar)

		// 添加主角
		{
			avatarDataConfig := g.gameDataConfig.AvatarDataMap[int32(mainCharAvatarId)]
			skillDepotId := int32(0)
			// 主角要单独设置
			if mainCharAvatarId == 10000005 {
				skillDepotId = 504
			} else if mainCharAvatarId == 10000007 {
				skillDepotId = 704
			} else {
				skillDepotId = avatarDataConfig.SkillDepotId
			}
			avatarSkillDepotDataConfig := g.gameDataConfig.AvatarSkillDepotDataMap[skillDepotId]
			player.AddAvatar(mainCharAvatarId, avatarDataConfig, avatarSkillDepotDataConfig)
			weaponId := uint64(g.snowflake.GenId())
			// 雾切
			player.AddWeapon(11509, weaponId)
			player.AvatarEquipWeapon(mainCharAvatarId, weaponId)
		}

		// 添加所有角色
		for avatarId, _ := range g.gameDataConfig.AvatarDataMap {
			if avatarId == 10000005 || avatarId == 10000007 {
				continue
			}
			if avatarId < 10000002 || avatarId >= 11000000 {
				continue
			}
			avatarDataConfig := g.gameDataConfig.AvatarDataMap[avatarId]
			avatarSkillDepotDataConfig := g.gameDataConfig.AvatarSkillDepotDataMap[avatarDataConfig.SkillDepotId]
			player.AddAvatar(uint32(avatarId), avatarDataConfig, avatarSkillDepotDataConfig)
			weaponId := uint64(g.snowflake.GenId())
			player.AddWeapon(uint32(avatarDataConfig.InitialWeapon), weaponId)
			player.AvatarEquipWeapon(uint32(avatarId), weaponId)
		}

		player.TeamConfig = model.NewTeamInfo()
		player.TeamConfig.AddAvatarToTeam(mainCharAvatarId, 0)

		g.userManager.AddUser(player)

		g.SendMsg(api.ApiSetPlayerBornDataRsp, userId, nil, nil)
		g.OnLoginOk(userId)
	}
}

func (g *GameManager) GetPlayerSocialDetailReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user get player social detail, user id: %v", userId)
	// TODO 构造社交信息管理器
	socialDetail := new(proto.SocialDetail)
	socialDetail.Uid = userId
	socialDetail.AvatarId = (&proto.HeadImage{AvatarId: 10000007}).GetAvatarId()
	socialDetail.Nickname = "flswld"
	socialDetail.Level = 60
	socialDetail.Birthday = &proto.Birthday{Month: 2, Day: 13}
	socialDetail.NameCardId = 210001
	getPlayerSocialDetailRsp := new(proto.GetPlayerSocialDetailRsp)
	getPlayerSocialDetailRsp.DetailData = socialDetail
	g.SendMsg(api.ApiGetPlayerSocialDetailRsp, userId, nil, getPlayerSocialDetailRsp)
}

func (g *GameManager) EnterSceneReadyReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user enter scene ready, user id: %v", userId)
	req := payloadMsg.(*proto.EnterSceneReadyReq)
	g.log.Debug("EnterSceneReadyReq: %v", req)
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		g.log.Error("player is nil, userId: %v", userId)
		return
	}
	g.log.Info("player.EnterSceneToken: %v", player.EnterSceneToken)
	enterScenePeerNotify := new(proto.EnterScenePeerNotify)
	enterScenePeerNotify.DestSceneId = player.SceneId
	world := g.worldManager.GetWorldByID(player.WorldId)
	enterScenePeerNotify.PeerId = player.PeerId
	enterScenePeerNotify.HostPeerId = world.owner.PeerId
	enterScenePeerNotify.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(api.ApiEnterScenePeerNotify, userId, nil, enterScenePeerNotify)

	enterSceneReadyRsp := new(proto.EnterSceneReadyRsp)
	enterSceneReadyRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(api.ApiEnterSceneReadyRsp, userId, g.getHeadMsg(11), enterSceneReadyRsp)
}

func (g *GameManager) PathfindingEnterSceneReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user pathfinding enter scene, user id: %v", userId)
	g.SendMsg(api.ApiPathfindingEnterSceneRsp, userId, g.getHeadMsg(headMsg.ClientSequenceId), nil)
}

func (g *GameManager) GetScenePointReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user get scene point, user id: %v", userId)
	req := payloadMsg.(*proto.GetScenePointReq)
	getScenePointRsp := new(proto.GetScenePointRsp)
	getScenePointRsp.SceneId = req.SceneId
	getScenePointRsp.UnlockedPointList = make([]uint32, 0)
	for i := uint32(1); i < 1000; i++ {
		getScenePointRsp.UnlockedPointList = append(getScenePointRsp.UnlockedPointList, i)
	}
	getScenePointRsp.UnlockAreaList = make([]uint32, 0)
	for i := uint32(1); i < 9; i++ {
		getScenePointRsp.UnlockAreaList = append(getScenePointRsp.UnlockAreaList, i)
	}
	g.SendMsg(api.ApiGetScenePointRsp, userId, nil, getScenePointRsp)
}

func (g *GameManager) GetSceneAreaReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user get scene area, user id: %v", userId)
	req := payloadMsg.(*proto.GetSceneAreaReq)
	getSceneAreaRsp := new(proto.GetSceneAreaRsp)
	getSceneAreaRsp.SceneId = req.SceneId
	getSceneAreaRsp.AreaIdList = []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 17, 18, 19, 20, 21, 22, 23, 24, 25, 100, 101, 102, 103, 200, 210, 300, 400, 401, 402, 403}
	getSceneAreaRsp.CityInfoList = make([]*proto.CityInfo, 0)
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 1, Level: 1})
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 2, Level: 1})
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 3, Level: 1})
	g.SendMsg(api.ApiGetSceneAreaRsp, userId, g.getHeadMsg(0), getSceneAreaRsp)
}

func (g *GameManager) SceneInitFinishReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user scene init finish, user id: %v", userId)

	// PacketServerTimeNotify
	serverTimeNotify := new(proto.ServerTimeNotify)
	serverTimeNotify.ServerTime = uint64(time.Now().UnixMilli())
	g.SendMsg(api.ApiServerTimeNotify, userId, nil, serverTimeNotify)

	// PacketWorldPlayerInfoNotify
	worldPlayerInfoNotify := new(proto.WorldPlayerInfoNotify)
	player := g.userManager.GetTargetUser(userId)
	world := g.worldManager.GetWorldByID(player.WorldId)
	for _, worldPlayer := range world.playerMap {
		onlinePlayerInfo := new(proto.OnlinePlayerInfo)
		onlinePlayerInfo.Uid = worldPlayer.PlayerID
		onlinePlayerInfo.Nickname = worldPlayer.NickName
		onlinePlayerInfo.PlayerLevel = 60
		onlinePlayerInfo.MpSettingType = worldPlayer.MpSetting
		onlinePlayerInfo.NameCardId = 210001
		onlinePlayerInfo.Signature = "惟愿时光记忆，一路繁花千树。"
		// 头像
		onlinePlayerInfo.ProfilePicture = &proto.ProfilePicture{AvatarId: 10000007}
		// TODO 待确定 这个到底是1p 2p 3p 4p的意思 还是世界内玩家数量的意思
		onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(world.playerMap))
		worldPlayerInfoNotify.PlayerInfoList = append(worldPlayerInfoNotify.PlayerInfoList, onlinePlayerInfo)
		worldPlayerInfoNotify.PlayerUidList = append(worldPlayerInfoNotify.PlayerUidList, worldPlayer.PlayerID)
	}
	g.SendMsg(api.ApiWorldPlayerInfoNotify, userId, nil, worldPlayerInfoNotify)

	// PacketWorldDataNotify
	worldDataNotify := new(proto.WorldDataNotify)
	worldDataNotify.WorldPropMap = make(map[uint32]*proto.PropValue)
	// 世界等级
	worldDataNotify.WorldPropMap[1] = &proto.PropValue{
		Type:  1,
		Value: &proto.PropValue_Ival{Ival: int64(world.worldLevel)},
	}
	// 是否多人游戏
	worldDataNotify.WorldPropMap[2] = &proto.PropValue{
		Type:  2,
		Value: &proto.PropValue_Ival{Ival: object.ConvBoolToInt64(world.multiplayer)},
	}
	g.SendMsg(api.ApiWorldDataNotify, userId, nil, worldDataNotify)

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
	g.SendMsg(api.ApiPlayerWorldSceneInfoListNotify, userId, nil, playerWorldSceneInfoListNotify)

	// SceneForceUnlockNotify
	g.SendMsg(api.ApiSceneForceUnlockNotify, userId, nil, nil)

	// PacketHostPlayerNotify
	hostPlayerNotify := new(proto.HostPlayerNotify)
	hostPlayerNotify.HostUid = world.owner.PlayerID
	hostPlayerNotify.HostPeerId = world.owner.PeerId
	g.SendMsg(api.ApiHostPlayerNotify, userId, nil, hostPlayerNotify)

	// PacketSceneTimeNotify
	sceneTimeNotify := new(proto.SceneTimeNotify)
	sceneTimeNotify.SceneId = player.SceneId
	sceneTimeNotify.SceneTime = 0
	g.SendMsg(api.ApiSceneTimeNotify, userId, nil, sceneTimeNotify)

	// PacketPlayerGameTimeNotify
	playerGameTimeNotify := new(proto.PlayerGameTimeNotify)
	playerGameTimeNotify.GameTime = 8 * 60 // 游戏内时间 time % 1440
	playerGameTimeNotify.Uid = player.PlayerID
	g.SendMsg(api.ApiPlayerGameTimeNotify, userId, nil, playerGameTimeNotify)

	// PacketPlayerEnterSceneInfoNotify
	empty := new(proto.AbilitySyncStateInfo)
	playerEnterSceneInfoNotify := new(proto.PlayerEnterSceneInfoNotify)
	playerEnterSceneInfoNotify.CurAvatarEntityId = player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId // 世界里面的实体id
	playerEnterSceneInfoNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneInfoNotify.TeamEnterInfo = &proto.TeamEnterSceneInfo{
		TeamEntityId:        player.TeamConfig.TeamEntityId, // 世界里面的实体id
		TeamAbilityInfo:     empty,
		AbilityControlBlock: new(proto.AbilityControlBlock),
	}
	playerEnterSceneInfoNotify.MpLevelEntityInfo = &proto.MPLevelEntityInfo{
		EntityId:        g.worldManager.GetWorldByID(player.WorldId).mpLevelEntityId, // 世界里面的实体id
		AuthorityPeerId: g.worldManager.GetWorldByID(player.WorldId).owner.PeerId,
		AbilityInfo:     empty,
	}
	activeTeam := player.TeamConfig.GetActiveTeam()
	for avatarIndex, avatarId := range activeTeam.AvatarIdList {
		if avatarId == 0 {
			break
		}
		avatarEnterSceneInfo := new(proto.AvatarEnterSceneInfo)
		avatarEnterSceneInfo.AvatarGuid = player.AvatarMap[avatarId].Guid
		avatarEnterSceneInfo.AvatarEntityId = player.TeamConfig.GetAvatarEntityByIndex(uint8(avatarIndex)).AvatarEntityId
		avatarEnterSceneInfo.WeaponGuid = player.AvatarMap[avatarId].EquipWeapon.Guid
		avatarEnterSceneInfo.WeaponEntityId = player.TeamConfig.GetAvatarEntityByIndex(uint8(avatarIndex)).WeaponEntityId
		avatarEnterSceneInfo.AvatarAbilityInfo = empty
		avatarEnterSceneInfo.WeaponAbilityInfo = empty
		playerEnterSceneInfoNotify.AvatarEnterInfo = append(playerEnterSceneInfoNotify.AvatarEnterInfo, avatarEnterSceneInfo)
	}
	g.SendMsg(api.ApiPlayerEnterSceneInfoNotify, userId, nil, playerEnterSceneInfoNotify)
	g.userManager.UpdateUser(player)

	// PacketSceneAreaWeatherNotify
	sceneAreaWeatherNotify := new(proto.SceneAreaWeatherNotify)
	sceneAreaWeatherNotify.WeatherAreaId = 0
	climateTypeConst := constant.GetClimateTypeConst()
	sceneAreaWeatherNotify.ClimateType = uint32(climateTypeConst.CLIMATE_SUNNY)
	g.SendMsg(api.ApiSceneAreaWeatherNotify, userId, nil, sceneAreaWeatherNotify)

	// PacketScenePlayerInfoNotify
	scenePlayerInfoNotify := new(proto.ScenePlayerInfoNotify)
	for _, worldPlayer := range world.playerMap {
		onlinePlayerInfo := new(proto.OnlinePlayerInfo)
		onlinePlayerInfo.Uid = worldPlayer.PlayerID
		onlinePlayerInfo.Nickname = worldPlayer.NickName
		onlinePlayerInfo.PlayerLevel = 60
		onlinePlayerInfo.MpSettingType = worldPlayer.MpSetting
		onlinePlayerInfo.NameCardId = 210001
		onlinePlayerInfo.Signature = "惟愿时光记忆，一路繁花千树。"
		// 头像
		onlinePlayerInfo.ProfilePicture = &proto.ProfilePicture{AvatarId: 10000007}
		// TODO 待确定 这个到底是1p 2p 3p 4p的意思 还是世界内玩家数量的意思
		onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(world.playerMap))
		scenePlayerInfoNotify.PlayerInfoList = append(scenePlayerInfoNotify.PlayerInfoList, &proto.ScenePlayerInfo{
			Uid:              worldPlayer.PlayerID,
			PeerId:           worldPlayer.PeerId,
			Name:             worldPlayer.NickName,
			SceneId:          worldPlayer.SceneId,
			OnlinePlayerInfo: onlinePlayerInfo,
		})
	}
	g.SendMsg(api.ApiScenePlayerInfoNotify, userId, nil, scenePlayerInfoNotify)

	// PacketSceneTeamUpdateNotify
	sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world)
	g.SendMsg(api.ApiSceneTeamUpdateNotify, userId, nil, sceneTeamUpdateNotify)

	// PacketSyncTeamEntityNotify
	syncTeamEntityNotify := new(proto.SyncTeamEntityNotify)
	syncTeamEntityNotify.SceneId = player.SceneId
	syncTeamEntityNotify.TeamEntityInfoList = make([]*proto.TeamEntityInfo, 0)
	g.SendMsg(api.ApiSyncTeamEntityNotify, userId, nil, syncTeamEntityNotify)

	// PacketSyncScenePlayTeamEntityNotify
	syncScenePlayTeamEntityNotify := new(proto.SyncScenePlayTeamEntityNotify)
	syncScenePlayTeamEntityNotify.SceneId = player.SceneId
	g.SendMsg(api.ApiSyncScenePlayTeamEntityNotify, userId, nil, syncScenePlayTeamEntityNotify)

	// PacketSceneInitFinishRsp
	SceneInitFinishRsp := new(proto.SceneInitFinishRsp)
	SceneInitFinishRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(api.ApiSceneInitFinishRsp, userId, g.getHeadMsg(11), SceneInitFinishRsp)
}

func (g *GameManager) EnterSceneDoneReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user enter scene done, user id: %v", userId)
	player := g.userManager.GetTargetUser(userId)

	// PacketEnterSceneDoneRsp
	enterSceneDoneRsp := new(proto.EnterSceneDoneRsp)
	enterSceneDoneRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(api.ApiEnterSceneDoneRsp, userId, nil, enterSceneDoneRsp)

	// PacketPlayerTimeNotify
	playerTimeNotify := new(proto.PlayerTimeNotify)
	playerTimeNotify.IsPaused = false
	playerTimeNotify.PlayerTime = uint64(0) // 客户端ping包的时间
	playerTimeNotify.ServerTime = uint64(time.Now().UnixMilli())
	g.SendMsg(api.ApiPlayerTimeNotify, userId, nil, playerTimeNotify)

	// PacketSceneEntityAppearNotify
	sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
	sceneEntityAppearNotify.AppearType = proto.VisionType_VISION_TYPE_BORN
	avatarId := player.TeamConfig.GetActiveAvatarId()
	playerPropertyConst := constant.GetPlayerPropertyConst()
	fightPropertyConst := constant.GetFightPropertyConst()
	equipIdList := make([]uint32, 0)
	weapon := player.AvatarMap[avatarId].EquipWeapon
	equipIdList = append(equipIdList, weapon.ItemId)
	for _, reliquary := range player.AvatarMap[avatarId].EquipReliquaryList {
		equipIdList = append(equipIdList, reliquary.ItemId)
	}
	sceneEntityInfo := &proto.SceneEntityInfo{
		EntityType: proto.ProtEntityType_PROT_ENTITY_TYPE_AVATAR,
		EntityId:   player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId,
		MotionInfo: &proto.MotionInfo{
			Pos: &proto.Vector{
				X: float32(player.Pos.X),
				Y: float32(player.Pos.Y),
				Z: float32(player.Pos.Z),
			},
			Rot: &proto.Vector{
				X: float32(player.Rotation.X),
				Y: float32(player.Rotation.Y),
				Z: float32(player.Rotation.Z),
			},
			Speed: &proto.Vector{},
		},
		PropList: []*proto.PropPair{{Type: uint32(playerPropertyConst.PROP_LEVEL), PropValue: &proto.PropValue{
			Type:  uint32(playerPropertyConst.PROP_LEVEL),
			Value: &proto.PropValue_Ival{Ival: int64(player.AvatarMap[avatarId].Level)},
			Val:   int64(player.AvatarMap[avatarId].Level),
		}}},
		FightPropList: []*proto.FightPropPair{
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_HP),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_HP),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_MAX_HP),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)],
			},
		},
		LifeState:        1,
		AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
		Entity: &proto.SceneEntityInfo_Avatar{
			Avatar: &proto.SceneAvatarInfo{
				Uid:          player.PlayerID,
				AvatarId:     avatarId,
				Guid:         player.AvatarMap[avatarId].Guid,
				PeerId:       player.PeerId,
				EquipIdList:  equipIdList,
				SkillDepotId: player.AvatarMap[avatarId].SkillDepotId,
				Weapon: &proto.SceneWeaponInfo{
					EntityId:    player.TeamConfig.GetActiveAvatarEntity().WeaponEntityId,
					GadgetId:    uint32(g.gameDataConfig.ItemDataMap[int32(weapon.ItemId)].GadgetId),
					ItemId:      weapon.ItemId,
					Guid:        weapon.Guid,
					Level:       uint32(weapon.Level),
					AbilityInfo: new(proto.AbilitySyncStateInfo),
				},
				ReliquaryList:     nil,
				SkillLevelMap:     player.AvatarMap[avatarId].SkillLevelMap,
				WearingFlycloakId: player.AvatarMap[avatarId].FlyCloak,
				BornTime:          uint32(player.AvatarMap[avatarId].BornTime),
			},
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
	}
	sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{sceneEntityInfo}
	g.SendMsg(api.ApiSceneEntityAppearNotify, userId, g.getHeadMsg(11), sceneEntityAppearNotify)

	// PacketWorldPlayerLocationNotify
	worldPlayerLocationNotify := new(proto.WorldPlayerLocationNotify)
	worldPlayerLocationNotify.PlayerWorldLocList = []*proto.PlayerWorldLocationInfo{{
		SceneId: player.SceneId,
		PlayerLoc: &proto.PlayerLocationInfo{
			Uid: player.PlayerID,
			Pos: &proto.Vector{
				X: float32(player.Pos.X),
				Y: float32(player.Pos.Y),
				Z: float32(player.Pos.Z),
			},
			Rot: &proto.Vector{
				X: float32(player.Rotation.X),
				Y: float32(player.Rotation.Y),
				Z: float32(player.Rotation.Z),
			},
		},
	}}
	g.SendMsg(api.ApiWorldPlayerLocationNotify, userId, nil, worldPlayerLocationNotify)

	// PacketScenePlayerLocationNotify
	scenePlayerLocationNotify := new(proto.ScenePlayerLocationNotify)
	scenePlayerLocationNotify.SceneId = player.SceneId
	scenePlayerLocationNotify.PlayerLocList = []*proto.PlayerLocationInfo{{
		Uid: player.PlayerID,
		Pos: &proto.Vector{
			X: float32(player.Pos.X),
			Y: float32(player.Pos.Y),
			Z: float32(player.Pos.Z),
		},
		Rot: &proto.Vector{
			X: float32(player.Rotation.X),
			Y: float32(player.Rotation.Y),
			Z: float32(player.Rotation.Z),
		},
	}}
	g.SendMsg(api.ApiScenePlayerLocationNotify, userId, nil, scenePlayerLocationNotify)

	// PacketWorldPlayerRTTNotify
	worldPlayerRTTNotify := new(proto.WorldPlayerRTTNotify)
	worldPlayerRTTNotify.PlayerRttList = []*proto.PlayerRTTInfo{{Uid: player.PlayerID, Rtt: 10}}
	g.SendMsg(api.ApiWorldPlayerRTTNotify, userId, nil, worldPlayerRTTNotify)
}

func (g *GameManager) EnterWorldAreaReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user enter world area, user id: %v", userId)
	req := payloadMsg.(*proto.EnterWorldAreaReq)
	enterWorldAreaRsp := new(proto.EnterWorldAreaRsp)
	enterWorldAreaRsp.AreaType = req.AreaType
	enterWorldAreaRsp.AreaId = req.AreaId
	g.SendMsg(api.ApiEnterWorldAreaRsp, userId, g.getHeadMsg(headMsg.ClientSequenceId), enterWorldAreaRsp)
}

func (g *GameManager) PostEnterSceneReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user post enter scene, user id: %v", userId)
	player := g.userManager.GetTargetUser(userId)
	postEnterSceneRsp := new(proto.PostEnterSceneRsp)
	postEnterSceneRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(api.ApiPostEnterSceneRsp, userId, nil, postEnterSceneRsp)
}

func (g *GameManager) TowerAllDataReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user get tower all data, user id: %v", userId)
	towerAllDataRsp := new(proto.TowerAllDataRsp)
	towerAllDataRsp.TowerScheduleId = 29
	towerAllDataRsp.TowerFloorRecordList = []*proto.TowerFloorRecord{{FloorId: 1001}}
	towerAllDataRsp.CurLevelRecord = &proto.TowerCurLevelRecord{IsEmpty: true}
	towerAllDataRsp.NextScheduleChangeTime = 4294967295
	towerAllDataRsp.FloorOpenTimeMap = make(map[uint32]uint32)
	towerAllDataRsp.FloorOpenTimeMap[1024] = 1630486800
	towerAllDataRsp.FloorOpenTimeMap[1025] = 1630486800
	towerAllDataRsp.FloorOpenTimeMap[1026] = 1630486800
	towerAllDataRsp.FloorOpenTimeMap[1027] = 1630486800
	towerAllDataRsp.ScheduleStartTime = 1630486800
	g.SendMsg(api.ApiTowerAllDataRsp, userId, nil, towerAllDataRsp)
}

func (g *GameManager) SceneTransToPointReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user get scene trans to point, user id: %v", userId)
	req := payloadMsg.(*proto.SceneTransToPointReq)

	transPointId := strconv.Itoa(int(req.SceneId)) + "_" + strconv.Itoa(int(req.PointId))
	transPointConfig, exist := g.gameDataConfig.ScenePointEntries[transPointId]
	if !exist {
		// PacketSceneTransToPointRsp
		sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
		// TODO Retcode.proto
		sceneTransToPointRsp.Retcode = 1 // RET_SVR_ERROR_VALUE
		g.SendMsg(api.ApiSceneTransToPointRsp, userId, nil, sceneTransToPointRsp)
		return
	}

	// 传送玩家
	player := g.userManager.GetTargetUser(userId)
	oldSceneId := player.SceneId
	newSceneId := req.SceneId
	world := g.worldManager.GetWorldByID(player.WorldId)
	oldScene := world.GetSceneById(oldSceneId)
	oldScene.RemovePlayer(player)
	newScene := world.GetSceneById(newSceneId)
	newScene.AddPlayer(player)
	player.TeamConfig.UpdateTeam(world.GetNextWorldEntityId, g.gameDataConfig.AvatarSkillDepotDataMap)
	player.Pos.X = transPointConfig.PointData.TranPos.X
	player.Pos.Y = transPointConfig.PointData.TranPos.Y
	player.Pos.Z = transPointConfig.PointData.TranPos.Z
	player.SceneId = newSceneId
	g.log.Info("player goto scene: %v, pos x: %v, y: %v, z: %v", newSceneId, player.Pos.X, player.Pos.Y, player.Pos.Z)
	g.userManager.UpdateUser(player)

	// PacketSceneEntityDisappearNotify
	sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
	sceneEntityDisappearNotify.EntityList = []uint32{player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId}
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REMOVE
	g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)

	// PacketPlayerEnterSceneNotify
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
	playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
	playerEnterSceneNotify.PrevSceneId = newSceneId
	playerEnterSceneNotify.PrevPos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneId = newSceneId
	playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
	if newSceneId == oldSceneId {
		playerEnterSceneNotify.Type = proto.EnterType_ENTER_TYPE_GOTO
	} else {
		playerEnterSceneNotify.Type = proto.EnterType_ENTER_TYPE_JUMP
	}
	playerEnterSceneNotify.TargetUid = player.PlayerID
	playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneNotify.WorldLevel = player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.TransPoint)
	playerEnterSceneNotify.SceneTagIdList = []uint32{102, 107, 109, 113, 117}
	playerEnterSceneNotify.WorldType = 1
	playerEnterSceneNotify.SceneTransaction = strconv.FormatInt(int64(newSceneId), 10) + "-" + strconv.FormatInt(int64(player.PlayerID), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + "18402"
	g.SendMsg(api.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)
	g.userManager.UpdateUser(player)

	// PacketSceneTransToPointRsp
	sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
	sceneTransToPointRsp.Retcode = 0
	sceneTransToPointRsp.PointId = req.PointId
	sceneTransToPointRsp.SceneId = req.SceneId
	g.SendMsg(api.ApiSceneTransToPointRsp, userId, nil, sceneTransToPointRsp)
}

func (g *GameManager) CombatInvocationsNotify(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user combat invocations, user id: %v", userId)
	req := payloadMsg.(*proto.CombatInvocationsNotify)
	player := g.userManager.GetTargetUser(userId)
	for _, v := range req.InvokeList {
		switch v.ArgumentType {
		case proto.CombatTypeArgument_COMBAT_TYPE_ARGUMENT_ENTITY_MOVE:
			{
				entityMoveInfo := new(proto.EntityMoveInfo)
				err := pb.Unmarshal(v.CombatData, entityMoveInfo)
				if err != nil {
					g.log.Error("parse combat invocations entity move info error: %v", err)
					continue
				}
				if entityMoveInfo.EntityId == player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId {
					// 玩家在移动
					if entityMoveInfo.MotionInfo.Pos == nil {
						g.log.Error("parse motion info pos is nil, entityMoveInfo: %v", entityMoveInfo)
						continue
					}
					player.Pos.X = float64(entityMoveInfo.MotionInfo.Pos.X)
					player.Pos.Y = float64(entityMoveInfo.MotionInfo.Pos.Y)
					player.Pos.Z = float64(entityMoveInfo.MotionInfo.Pos.Z)
					player.Rotation.X = float64(entityMoveInfo.MotionInfo.Rot.X)
					player.Rotation.Y = float64(entityMoveInfo.MotionInfo.Rot.Y)
					player.Rotation.Z = float64(entityMoveInfo.MotionInfo.Rot.Z)
				}
			}
		}
	}
}

func (g *GameManager) MarkMapReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user mark map, user id: %v", userId)
	req := payloadMsg.(*proto.MarkMapReq)
	operation := req.Op
	if operation == proto.MarkMapReq_OPERATION_ADD {
		g.log.Debug("user mark type: %v", req.Mark.PointType)
		if req.Mark.PointType == proto.MapMarkPointType_MAP_MARK_POINT_TYPE_NPC {
			posYInt, err := strconv.ParseInt(req.Mark.Name, 10, 64)
			if err != nil {
				g.log.Error("parse pos y error: %v", err)
				posYInt = 0
			}

			// 传送玩家
			player := g.userManager.GetTargetUser(userId)
			oldSceneId := player.SceneId
			newSceneId := req.Mark.SceneId
			world := g.worldManager.GetWorldByID(player.WorldId)
			oldScene := world.GetSceneById(oldSceneId)
			oldScene.RemovePlayer(player)
			newScene := world.GetSceneById(newSceneId)
			newScene.AddPlayer(player)
			player.TeamConfig.UpdateTeam(world.GetNextWorldEntityId, g.gameDataConfig.AvatarSkillDepotDataMap)
			x := float64(req.Mark.Pos.X)
			y := float64(posYInt)
			z := float64(req.Mark.Pos.Z)
			player.Pos.X = x
			player.Pos.Y = y
			player.Pos.Z = z
			player.SceneId = newSceneId
			g.log.Info("player goto scene: %v, pos x: %v, y: %v, z: %v", newSceneId, x, y, z)
			g.userManager.UpdateUser(player)

			// PacketSceneEntityDisappearNotify
			sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
			sceneEntityDisappearNotify.EntityList = []uint32{player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId}
			sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REMOVE
			g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)
			g.userManager.UpdateUser(player)

			// PacketPlayerEnterSceneNotify
			playerPropertyConst := constant.GetPlayerPropertyConst()
			player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
			playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
			playerEnterSceneNotify.PrevSceneId = oldSceneId
			playerEnterSceneNotify.PrevPos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
			playerEnterSceneNotify.SceneId = newSceneId
			playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
			playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
			if newSceneId == oldSceneId {
				playerEnterSceneNotify.Type = proto.EnterType_ENTER_TYPE_GOTO
			} else {
				playerEnterSceneNotify.Type = proto.EnterType_ENTER_TYPE_JUMP
			}
			playerEnterSceneNotify.TargetUid = player.PlayerID
			playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
			playerEnterSceneNotify.WorldLevel = player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
			enterReasonConst := constant.GetEnterReasonConst()
			playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.TransPoint)
			playerEnterSceneNotify.SceneTagIdList = []uint32{102, 107, 109, 113, 117}
			playerEnterSceneNotify.WorldType = 1
			playerEnterSceneNotify.SceneTransaction = strconv.FormatInt(int64(newSceneId), 10) + "-" + strconv.FormatInt(int64(player.PlayerID), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + "18402"
			g.SendMsg(api.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)
			g.userManager.UpdateUser(player)
		}
	}
}

func (g *GameManager) ChangeAvatarReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user change avatar, user id: %v", userId)
	req := payloadMsg.(*proto.ChangeAvatarReq)
	targetAvatarGuid := req.Guid
	player := g.userManager.GetTargetUser(userId)
	oldAvatarId := player.TeamConfig.GetActiveAvatarId()
	oldAvatarEntity := player.TeamConfig.GetActiveAvatarEntity()
	oldAvatar := player.AvatarMap[oldAvatarId]
	if oldAvatar.Guid == targetAvatarGuid {
		g.log.Error("can not change to the same avatar, user id: %v, oldAvatarId: %v, oldAvatarGuid: %v", userId, oldAvatarId, oldAvatar.Guid)
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
		g.log.Error("can not find the target avatar in team, user id: %v, target avatar guid: %v", userId, targetAvatarGuid)
		return
	}
	player.TeamConfig.CurrAvatarIndex = uint8(index)

	// TODO 目前多人游戏可能会存在问题 可能需要将原来的队伍里的角色实体放到世界里去才行 只是可能而已 待验证

	// PacketSceneEntityDisappearNotify
	sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REPLACE
	sceneEntityDisappearNotify.EntityList = []uint32{oldAvatarEntity.AvatarEntityId}
	g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)

	// PacketSceneEntityAppearNotify
	sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REPLACE
	sceneEntityAppearNotify.Param = oldAvatarEntity.AvatarEntityId
	sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{g.PacketAvatarSceneEntityInfo(player, player.TeamConfig.GetActiveAvatarId())}
	g.SendMsg(api.ApiSceneEntityAppearNotify, userId, g.getHeadMsg(11), sceneEntityAppearNotify)

	// PacketChangeAvatarRsp
	changeAvatarRsp := new(proto.ChangeAvatarRsp)
	changeAvatarRsp.Retcode = int32(proto.Retcode_RET_SUCC)
	changeAvatarRsp.CurGuid = targetAvatarGuid
	g.SendMsg(api.ApiChangeAvatarRsp, userId, nil, changeAvatarRsp)
}

func (g *GameManager) SetUpAvatarTeamReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user set up avatar team, user id: %v", userId)
	req := payloadMsg.(*proto.SetUpAvatarTeamReq)
	player := g.userManager.GetTargetUser(userId)
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
		g.SendMsg(api.ApiAvatarTeamUpdateNotify, userId, nil, avatarTeamUpdateNotify)

		if selfTeam {
			player.TeamConfig.CurrAvatarIndex = 0
			player.TeamConfig.UpdateTeam(world.GetNextWorldEntityId, g.gameDataConfig.AvatarSkillDepotDataMap)
			// TODO 还有一大堆没写 SceneTeamUpdateNotify
			// PacketSceneTeamUpdateNotify
			sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world)
			g.SendMsg(api.ApiSceneTeamUpdateNotify, userId, nil, sceneTeamUpdateNotify)

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
			g.SendMsg(api.ApiSetUpAvatarTeamRsp, userId, nil, setUpAvatarTeamRsp)
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
			g.SendMsg(api.ApiSetUpAvatarTeamRsp, userId, nil, setUpAvatarTeamRsp)
		}
	}
}

func (g *GameManager) ChooseCurAvatarTeamReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	g.log.Info("user switch team, user id: %v", userId)
	req := payloadMsg.(*proto.ChooseCurAvatarTeamReq)
	teamId := req.TeamId
	player := g.userManager.GetTargetUser(userId)
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
	player.TeamConfig.UpdateTeam(world.GetNextWorldEntityId, g.gameDataConfig.AvatarSkillDepotDataMap)

	// TODO 还有一大堆没写 SceneTeamUpdateNotify
	// PacketSceneTeamUpdateNotify
	sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world)
	g.SendMsg(api.ApiSceneTeamUpdateNotify, userId, nil, sceneTeamUpdateNotify)

	// PacketChooseCurAvatarTeamRsp
	chooseCurAvatarTeamRsp := new(proto.ChooseCurAvatarTeamRsp)
	chooseCurAvatarTeamRsp.CurTeamId = teamId
	g.SendMsg(api.ApiChooseCurAvatarTeamRsp, userId, nil, chooseCurAvatarTeamRsp)
}

func (g *GameManager) PacketAvatarSceneEntityInfo(player *model.Player, avatarId uint32) *proto.SceneEntityInfo {
	playerPropertyConst := constant.GetPlayerPropertyConst()
	fightPropertyConst := constant.GetFightPropertyConst()
	equipIdList := make([]uint32, 0)
	weapon := player.AvatarMap[avatarId].EquipWeapon
	equipIdList = append(equipIdList, weapon.ItemId)
	for _, reliquary := range player.AvatarMap[avatarId].EquipReliquaryList {
		equipIdList = append(equipIdList, reliquary.ItemId)
	}
	sceneEntityInfo := &proto.SceneEntityInfo{
		EntityType: proto.ProtEntityType_PROT_ENTITY_TYPE_AVATAR,
		EntityId:   player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId,
		MotionInfo: &proto.MotionInfo{
			Pos: &proto.Vector{
				X: float32(player.Pos.X),
				Y: float32(player.Pos.Y),
				Z: float32(player.Pos.Z),
			},
			Rot: &proto.Vector{
				X: float32(player.Rotation.X),
				Y: float32(player.Rotation.Y),
				Z: float32(player.Rotation.Z),
			},
			Speed: &proto.Vector{},
		},
		PropList: []*proto.PropPair{{Type: uint32(playerPropertyConst.PROP_LEVEL), PropValue: &proto.PropValue{
			Type:  uint32(playerPropertyConst.PROP_LEVEL),
			Value: &proto.PropValue_Ival{Ival: int64(player.AvatarMap[avatarId].Level)},
			Val:   int64(player.AvatarMap[avatarId].Level),
		}}},
		FightPropList: []*proto.FightPropPair{
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_HP),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_HP),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_MAX_HP),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)],
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE),
				PropValue: player.AvatarMap[avatarId].FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)],
			},
		},
		LifeState:        1,
		AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
		Entity: &proto.SceneEntityInfo_Avatar{
			Avatar: &proto.SceneAvatarInfo{
				Uid:          player.PlayerID,
				AvatarId:     avatarId,
				Guid:         player.AvatarMap[avatarId].Guid,
				PeerId:       player.PeerId,
				EquipIdList:  equipIdList,
				SkillDepotId: player.AvatarMap[avatarId].SkillDepotId,
				Weapon: &proto.SceneWeaponInfo{
					EntityId:    player.TeamConfig.GetActiveAvatarEntity().WeaponEntityId,
					GadgetId:    uint32(g.gameDataConfig.ItemDataMap[int32(weapon.ItemId)].GadgetId),
					ItemId:      weapon.ItemId,
					Guid:        weapon.Guid,
					Level:       uint32(weapon.Level),
					AbilityInfo: new(proto.AbilitySyncStateInfo),
				},
				ReliquaryList:     nil,
				SkillLevelMap:     player.AvatarMap[avatarId].SkillLevelMap,
				WearingFlycloakId: player.AvatarMap[avatarId].FlyCloak,
				BornTime:          uint32(player.AvatarMap[avatarId].BornTime),
			},
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
	}
	return sceneEntityInfo
}

func (g *GameManager) PacketSceneTeamUpdateNotify(world *World) *proto.SceneTeamUpdateNotify {
	sceneTeamUpdateNotify := new(proto.SceneTeamUpdateNotify)
	sceneTeamUpdateNotify.IsInMp = world.multiplayer
	playerPropertyConst := constant.GetPlayerPropertyConst()
	fightPropertyConst := constant.GetFightPropertyConst()
	empty := new(proto.AbilitySyncStateInfo)
	for _, worldPlayer := range world.playerMap {
		team := worldPlayer.TeamConfig.GetActiveTeam()
		for avatarIndex, avatarId := range team.AvatarIdList {
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
				PlayerUid:  worldPlayer.PlayerID,
				AvatarGuid: worldPlayerAvatar.Guid,
				SceneId:    worldPlayer.SceneId,
				EntityId:   worldPlayer.TeamConfig.GetAvatarEntityByIndex(uint8(avatarIndex)).AvatarEntityId,
				SceneEntityInfo: &proto.SceneEntityInfo{
					EntityType: proto.ProtEntityType_PROT_ENTITY_TYPE_AVATAR,
					EntityId:   worldPlayer.TeamConfig.GetAvatarEntityByIndex(uint8(avatarIndex)).AvatarEntityId,
					MotionInfo: &proto.MotionInfo{
						Pos: &proto.Vector{
							X: float32(worldPlayer.Pos.X),
							Y: float32(worldPlayer.Pos.Y),
							Z: float32(worldPlayer.Pos.Z),
						},
						Rot: &proto.Vector{
							X: float32(worldPlayer.Rotation.X),
							Y: float32(worldPlayer.Rotation.Y),
							Z: float32(worldPlayer.Rotation.Z),
						},
						Speed: &proto.Vector{},
					},
					PropList: []*proto.PropPair{{Type: uint32(playerPropertyConst.PROP_LEVEL), PropValue: &proto.PropValue{
						Type:  uint32(playerPropertyConst.PROP_LEVEL),
						Value: &proto.PropValue_Ival{Ival: int64(worldPlayerAvatar.Level)},
						Val:   int64(worldPlayerAvatar.Level),
					}}},
					FightPropList: []*proto.FightPropPair{
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_HP),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_HP)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL_HURT)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CHARGE_EFFICIENCY)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_HP),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_MAX_HP),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_MAX_HP)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)],
						},
						{
							PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE),
							PropValue: worldPlayerAvatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE)],
						},
					},
					LifeState:        1,
					AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
					Entity: &proto.SceneEntityInfo_Avatar{
						Avatar: &proto.SceneAvatarInfo{
							Uid:          worldPlayer.PlayerID,
							AvatarId:     avatarId,
							Guid:         worldPlayerAvatar.Guid,
							PeerId:       worldPlayer.PeerId,
							EquipIdList:  equipIdList,
							SkillDepotId: worldPlayerAvatar.SkillDepotId,
							Weapon: &proto.SceneWeaponInfo{
								EntityId:    worldPlayer.TeamConfig.GetAvatarEntityByIndex(uint8(avatarIndex)).WeaponEntityId,
								GadgetId:    uint32(g.gameDataConfig.ItemDataMap[int32(weapon.ItemId)].GadgetId),
								ItemId:      weapon.ItemId,
								Guid:        weapon.Guid,
								Level:       uint32(weapon.Level),
								AbilityInfo: new(proto.AbilitySyncStateInfo),
							},
							ReliquaryList:     nil,
							SkillLevelMap:     worldPlayerAvatar.SkillLevelMap,
							WearingFlycloakId: worldPlayerAvatar.FlyCloak,
							BornTime:          uint32(worldPlayerAvatar.BornTime),
						},
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
				},
				WeaponGuid:          worldPlayerAvatar.EquipWeapon.Guid,
				WeaponEntityId:      worldPlayer.TeamConfig.GetAvatarEntityByIndex(uint8(avatarIndex)).WeaponEntityId,
				IsPlayerCurAvatar:   worldPlayer.TeamConfig.GetActiveAvatarId() == avatarId,
				IsOnScene:           worldPlayer.TeamConfig.GetActiveAvatarId() == avatarId,
				AvatarAbilityInfo:   empty,
				WeaponAbilityInfo:   empty,
				AbilityControlBlock: new(proto.AbilityControlBlock),
			}
			// add AbilityControlBlock
			avatarDataConfig := g.gameDataConfig.AvatarDataMap[int32(avatarId)]
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
			for id, _ := range worldPlayer.TeamConfig.TeamResonancesConfig {
				embryoId++
				emb := &proto.AbilityEmbryo{
					AbilityId:               uint32(embryoId),
					AbilityNameHash:         uint32(id),
					AbilityOverrideNameHash: uint32(gameConstant.DEFAULT_ABILITY_NAME),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, emb)
			}
			// add skill depot abilities
			skillDepot := g.gameDataConfig.AvatarSkillDepotDataMap[int32(worldPlayerAvatar.SkillDepotId)]
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
			for skill, _ := range worldPlayerAvatar.ExtraAbilityEmbryos {
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
