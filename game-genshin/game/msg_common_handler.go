package game

import (
	"flswld.com/common/utils/random"
	"flswld.com/common/utils/reflection"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"game-genshin/entity"
	"game-genshin/game/constant"
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
		// TODO 把初始选择的主角的角色信息写入
		g.log.Debug("avatar id: %v, nickname: %v", req.AvatarId, req.NickName)
		player := new(entity.Player)
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
		player.Properties[playerPropertyConst.PROP_PLAYER_LEVEL] = 1
		player.Properties[playerPropertyConst.PROP_IS_SPRING_AUTO_USE] = 1
		player.Properties[playerPropertyConst.PROP_SPRING_AUTO_USE_PERCENT] = 50
		player.Properties[playerPropertyConst.PROP_IS_FLYABLE] = 1
		player.Properties[playerPropertyConst.PROP_IS_TRANSFERABLE] = 1
		player.Properties[playerPropertyConst.PROP_MAX_STAMINA] = 24000
		player.Properties[playerPropertyConst.PROP_CUR_PERSIST_STAMINA] = 24000
		player.Properties[playerPropertyConst.PROP_PLAYER_RESIN] = 160

		player.FlyCloakList = make([]uint32, 0)
		player.FlyCloakList = append(player.FlyCloakList, 140001)

		player.Pos = &entity.Vector{X: 2747, Y: 194, Z: -1719}
		player.Rotation = &entity.Vector{X: 0, Y: 307, Z: 0}

		player.MpSetting = proto.MpSettingType_MP_SETTING_ENTER_AFTER_APPLY
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
	socialDetail.Level = 1
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
	enterScenePeerNotify.DestSceneId = uint32(player.SceneId)
	// TODO 要做世界管理器
	enterScenePeerNotify.PeerId = 1
	enterScenePeerNotify.HostPeerId = 1
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
	getSceneAreaRsp.AreaIdList = []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 17, 18, 19, 100, 101, 102, 103, 200, 210, 300}
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
	// TODO 世界管理器
	player := g.userManager.GetTargetUser(userId)
	onlinePlayerInfo := new(proto.OnlinePlayerInfo)
	onlinePlayerInfo.Uid = player.PlayerID
	onlinePlayerInfo.Nickname = player.NickName
	onlinePlayerInfo.PlayerLevel = 1
	onlinePlayerInfo.MpSettingType = player.MpSetting
	onlinePlayerInfo.NameCardId = 210001
	onlinePlayerInfo.AvatarId = (&proto.HeadImage{AvatarId: 10000007}).GetAvatarId()
	onlinePlayerInfo.CurPlayerNumInWorld = 1 // 1p 2p 3p 4p 的意思
	worldPlayerInfoNotify.PlayerInfoList = []*proto.OnlinePlayerInfo{onlinePlayerInfo}
	worldPlayerInfoNotify.PlayerUidList = []uint32{player.PlayerID}
	g.SendMsg(api.ApiWorldPlayerInfoNotify, userId, nil, worldPlayerInfoNotify)
	// PacketWorldDataNotify
	worldDataNotify := new(proto.WorldDataNotify)
	worldDataNotify.WorldPropMap = make(map[uint32]*proto.PropValue)
	worldDataNotify.WorldPropMap[1] = &proto.PropValue{Type: 1, Value: &proto.PropValue_Ival{Ival: 0}} // 世界等级
	worldDataNotify.WorldPropMap[2] = &proto.PropValue{Type: 2, Value: &proto.PropValue_Ival{Ival: 0}} // 是否多人游戏
	g.SendMsg(api.ApiWorldDataNotify, userId, nil, worldDataNotify)
	// PacketSceneUnlockInfoNotify
	sceneUnlockInfoNotify := new(proto.SceneUnlockInfoNotify)
	sceneUnlockInfoNotify.UnlockInfos = []*proto.SceneUnlockInfo{
		{SceneId: 1, IsLocked: false, SceneTagIdList: []uint32{}},
		{SceneId: 3, IsLocked: false, SceneTagIdList: []uint32{102, 113, 117}},
		{SceneId: 4, IsLocked: false, SceneTagIdList: []uint32{106, 109}},
		{SceneId: 5, IsLocked: false, SceneTagIdList: []uint32{}},
		{SceneId: 6, IsLocked: false, SceneTagIdList: []uint32{}},
		{SceneId: 7, IsLocked: false, SceneTagIdList: []uint32{}},
	}
	g.SendMsg(api.ApiSceneUnlockInfoNotify, userId, nil, sceneUnlockInfoNotify)
	// SceneForceUnlockNotify
	g.SendMsg(api.ApiSceneForceUnlockNotify, userId, nil, nil)
	// PacketHostPlayerNotify
	hostPlayerNotify := new(proto.HostPlayerNotify)
	hostPlayerNotify.HostUid = player.PlayerID
	hostPlayerNotify.HostPeerId = 1
	g.SendMsg(api.ApiHostPlayerNotify, userId, nil, hostPlayerNotify)
	// PacketSceneTimeNotify
	sceneTimeNotify := new(proto.SceneTimeNotify)
	sceneTimeNotify.SceneId = uint32(player.SceneId)
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
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	player.AvatarEntityId = g.worldManager.GetNextWorldEntityID(entityIdTypeConst.AVATAR)
	player.WeaponEntityId = g.worldManager.GetNextWorldEntityID(entityIdTypeConst.WEAPON)
	playerEnterSceneInfoNotify.CurAvatarEntityId = player.AvatarEntityId // 世界里面的实体id
	playerEnterSceneInfoNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneInfoNotify.TeamEnterInfo = &proto.TeamEnterSceneInfo{
		TeamEntityId:        g.worldManager.GetNextWorldEntityID(entityIdTypeConst.TEAM), // 世界里面的实体id
		TeamAbilityInfo:     empty,
		AbilityControlBlock: new(proto.AbilityControlBlock),
	}
	playerEnterSceneInfoNotify.MpLevelEntityInfo = &proto.MPLevelEntityInfo{
		EntityId:        g.worldManager.GetNextWorldEntityID(entityIdTypeConst.MPLEVEL), // 世界里面的实体id
		AuthorityPeerId: 1,
		AbilityInfo:     empty,
	}
	avatarEnterSceneInfo := new(proto.AvatarEnterSceneInfo)
	avatarEnterSceneInfo.AvatarGuid = 429496733894967297
	avatarEnterSceneInfo.AvatarEntityId = player.AvatarEntityId
	avatarEnterSceneInfo.WeaponGuid = 429496733894967298
	avatarEnterSceneInfo.WeaponEntityId = player.WeaponEntityId
	avatarEnterSceneInfo.AvatarAbilityInfo = empty
	avatarEnterSceneInfo.WeaponAbilityInfo = empty
	playerEnterSceneInfoNotify.AvatarEnterInfo = []*proto.AvatarEnterSceneInfo{avatarEnterSceneInfo}
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
	scenePlayerInfoNotify.PlayerInfoList = []*proto.ScenePlayerInfo{{
		Uid:              player.PlayerID,
		PeerId:           1,
		Name:             player.NickName,
		SceneId:          uint32(player.SceneId),
		OnlinePlayerInfo: onlinePlayerInfo,
	}}
	g.SendMsg(api.ApiScenePlayerInfoNotify, userId, nil, scenePlayerInfoNotify)
	// PacketSceneTeamUpdateNotify
	sceneTeamUpdateNotify := new(proto.SceneTeamUpdateNotify)
	sceneTeamUpdateNotify.IsInMp = false
	sceneTeamUpdateNotify.SceneTeamAvatarList = []*proto.SceneTeamAvatar{{
		PlayerUid:  player.PlayerID,
		AvatarGuid: 429496733894967297,
		SceneId:    uint32(player.SceneId),
		EntityId:   player.AvatarEntityId,
		SceneEntityInfo: &proto.SceneEntityInfo{
			EntityType: proto.ProtEntityType_PROT_ENTITY_AVATAR,
			EntityId:   player.AvatarEntityId,
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
			PropList: []*proto.PropPair{{Type: 4001, PropValue: &proto.PropValue{
				Type:  4001,
				Value: &proto.PropValue_Ival{Ival: 1},
				Val:   1,
			}}},
			FightPropList: []*proto.FightPropPair{
				{
					PropType:  1010,
					PropValue: 911.791,
				},
				{
					PropType:  4,
					PropValue: 41.053,
				},
				{
					PropType:  2002,
					PropValue: 57.225,
				},
				{
					PropType:  2001,
					PropValue: 41.053,
				},
				{
					PropType:  2000,
					PropValue: 911.791,
				},
				{
					PropType:  1,
					PropValue: 911.791,
				},
				{
					PropType:  7,
					PropValue: 57.225,
				},
				{
					PropType:  23,
					PropValue: 1.0,
				},
				{
					PropType:  22,
					PropValue: 0.5,
				},
				{
					PropType:  20,
					PropValue: 0.05,
				},
			},
			LifeState:        1,
			AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
			Entity: &proto.SceneEntityInfo_Avatar{
				Avatar: &proto.SceneAvatarInfo{
					Uid:          player.PlayerID,
					AvatarId:     10000007,
					Guid:         429496733894967297,
					PeerId:       1,
					EquipIdList:  []uint32{11509},
					SkillDepotId: 704,
					Weapon: &proto.SceneWeaponInfo{
						EntityId:    player.WeaponEntityId,
						GadgetId:    50011509,
						ItemId:      11509,
						Guid:        429496733894967298,
						Level:       1,
						AbilityInfo: new(proto.AbilitySyncStateInfo),
					},
					SkillLevelMap: map[uint32]uint32{
						10067:  1,
						10068:  1,
						100553: 1,
					},
					WearingFlycloakId: 140001,
					BornTime:          1652555787,
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
		WeaponGuid:        429496733894967298,
		WeaponEntityId:    player.WeaponEntityId,
		IsPlayerCurAvatar: true,
		IsOnScene:         true,
		AvatarAbilityInfo: empty,
		WeaponAbilityInfo: empty,
		// 不填了 看看会不会报错
		AbilityControlBlock: &proto.AbilityControlBlock{
			AbilityEmbryoList: []*proto.AbilityEmbryo{
				{
					AbilityId:               1,
					AbilityNameHash:         4291357363,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               2,
					AbilityNameHash:         1410219662,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               3,
					AbilityNameHash:         1474894886,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               4,
					AbilityNameHash:         3832178184,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               5,
					AbilityNameHash:         2306062007,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               6,
					AbilityNameHash:         3105629177,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               7,
					AbilityNameHash:         3771526669,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               8,
					AbilityNameHash:         100636247,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               9,
					AbilityNameHash:         1564404322,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               10,
					AbilityNameHash:         497711942,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               11,
					AbilityNameHash:         825255509,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               12,
					AbilityNameHash:         1142761247,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               13,
					AbilityNameHash:         518324758,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               14,
					AbilityNameHash:         3276790745,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               15,
					AbilityNameHash:         3429175060,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               16,
					AbilityNameHash:         3429175061,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               17,
					AbilityNameHash:         4253958193,
					AbilityOverrideNameHash: 1178079449,
				},
				{
					AbilityId:               18,
					AbilityNameHash:         209033715,
					AbilityOverrideNameHash: 1178079449,
				},
			},
		},
	}}
	g.SendMsg(api.ApiSceneTeamUpdateNotify, userId, nil, sceneTeamUpdateNotify)
	// PacketSyncTeamEntityNotify
	syncTeamEntityNotify := new(proto.SyncTeamEntityNotify)
	syncTeamEntityNotify.SceneId = uint32(player.SceneId)
	syncTeamEntityNotify.TeamEntityInfoList = make([]*proto.TeamEntityInfo, 0)
	g.SendMsg(api.ApiSyncTeamEntityNotify, userId, nil, syncTeamEntityNotify)
	// PacketSyncScenePlayTeamEntityNotify
	syncScenePlayTeamEntityNotify := new(proto.SyncScenePlayTeamEntityNotify)
	syncScenePlayTeamEntityNotify.SceneId = uint32(player.SceneId)
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
	sceneEntityAppearNotify.AppearType = proto.VisionType_VISION_BORN
	sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{{
		EntityType: proto.ProtEntityType_PROT_ENTITY_AVATAR,
		EntityId:   player.AvatarEntityId,
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
		PropList: []*proto.PropPair{{Type: 4001, PropValue: &proto.PropValue{
			Type:  4001,
			Value: &proto.PropValue_Ival{Ival: 1},
			Val:   1,
		}}},
		FightPropList: []*proto.FightPropPair{
			{
				PropType:  1010,
				PropValue: 911.791,
			},
			{
				PropType:  4,
				PropValue: 41.053,
			},
			{
				PropType:  2002,
				PropValue: 57.225,
			},
			{
				PropType:  2001,
				PropValue: 41.053,
			},
			{
				PropType:  2000,
				PropValue: 911.791,
			},
			{
				PropType:  1,
				PropValue: 911.791,
			},
			{
				PropType:  7,
				PropValue: 57.225,
			},
			{
				PropType:  23,
				PropValue: 1.0,
			},
			{
				PropType:  22,
				PropValue: 0.5,
			},
			{
				PropType:  20,
				PropValue: 0.05,
			},
		},
		LifeState:        1,
		AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
		Entity: &proto.SceneEntityInfo_Avatar{
			Avatar: &proto.SceneAvatarInfo{
				Uid:          player.PlayerID,
				AvatarId:     10000007,
				Guid:         429496733894967297,
				PeerId:       1,
				EquipIdList:  []uint32{11509},
				SkillDepotId: 704,
				Weapon: &proto.SceneWeaponInfo{
					EntityId:    player.WeaponEntityId,
					GadgetId:    50011509,
					ItemId:      11509,
					Guid:        429496733894967298,
					Level:       1,
					AbilityInfo: new(proto.AbilitySyncStateInfo),
				},
				SkillLevelMap: map[uint32]uint32{
					10067:  1,
					10068:  1,
					100553: 1,
				},
				WearingFlycloakId: 140001,
				BornTime:          1652555787,
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
	}}
	g.SendMsg(api.ApiSceneEntityAppearNotify, userId, g.getHeadMsg(11), sceneEntityAppearNotify)
	// PacketWorldPlayerLocationNotify
	worldPlayerLocationNotify := new(proto.WorldPlayerLocationNotify)
	worldPlayerLocationNotify.PlayerWorldLocList = []*proto.PlayerWorldLocationInfo{{
		SceneId: uint32(player.SceneId),
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
	scenePlayerLocationNotify.SceneId = uint32(player.SceneId)
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
	player := g.userManager.GetTargetUser(userId)
	// PacketSceneEntityDisappearNotify
	sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
	sceneEntityDisappearNotify.EntityList = []uint32{player.AvatarEntityId}
	sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_REMOVE
	g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	player.AvatarEntityId = g.worldManager.GetNextWorldEntityID(entityIdTypeConst.AVATAR)
	player.WeaponEntityId = g.worldManager.GetNextWorldEntityID(entityIdTypeConst.WEAPON)
	transPointId := strconv.Itoa(int(req.SceneId)) + "_" + strconv.Itoa(int(req.PointId))
	transPoint, exist := g.gameDataConfig.ScenePointEntries[transPointId]
	if !exist {
		// PacketSceneTransToPointRsp
		sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
		// TODO Retcode.proto
		sceneTransToPointRsp.Retcode = 1 // RET_SVR_ERROR_VALUE
		g.SendMsg(api.ApiSceneTransToPointRsp, userId, nil, sceneTransToPointRsp)
		return
	}
	player.Pos.X = transPoint.PointData.TranPos.X
	player.Pos.Y = transPoint.PointData.TranPos.Y
	player.Pos.Z = transPoint.PointData.TranPos.Z
	oldSceneId := player.SceneId
	player.SceneId = uint16(req.SceneId)
	g.log.Info("player goto scene: %v, pos x: %v, y: %v, z: %v", player.SceneId, player.Pos.X, player.Pos.Y, player.Pos.Z)
	g.userManager.UpdateUser(player)
	// PacketPlayerEnterSceneNotify
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
	playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
	playerEnterSceneNotify.PrevSceneId = uint32(player.SceneId)
	playerEnterSceneNotify.PrevPos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneId = uint32(player.SceneId)
	playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
	if player.SceneId == oldSceneId {
		playerEnterSceneNotify.Type = proto.EnterType_ENTER_GOTO
	} else {
		playerEnterSceneNotify.Type = proto.EnterType_ENTER_JUMP
	}
	playerEnterSceneNotify.TargetUid = player.PlayerID
	playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneNotify.WorldLevel = player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.TransPoint)
	playerEnterSceneNotify.SceneTagIdList = []uint32{102, 107, 109, 113, 117}
	playerEnterSceneNotify.WorldType = 1
	playerEnterSceneNotify.SceneTransaction = strconv.FormatInt(int64(player.SceneId), 10) + "-" + strconv.FormatInt(int64(player.PlayerID), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + "18402"
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
		case proto.CombatTypeArgument_ENTITY_MOVE:
			{
				entityMoveInfo := new(proto.EntityMoveInfo)
				err := pb.Unmarshal(v.CombatData, entityMoveInfo)
				if err != nil {
					g.log.Error("parse combat invocations entity move info error: %v", err)
					continue
				}
				if entityMoveInfo.EntityId == player.AvatarEntityId {
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
	if operation == proto.MarkMapReq_ADD {
		g.log.Debug("user mark type: %v", req.Mark.PointType)
		if req.Mark.PointType == proto.MapMarkPointType_MAP_MARK_POINT_TYPE_NPC {
			// 传送玩家
			posYInt, err := strconv.ParseInt(req.Mark.Name, 10, 64)
			if err != nil {
				g.log.Error("parse pos y error: %v", err)
				return
			}
			x := float64(req.Mark.Pos.X)
			y := float64(posYInt)
			z := float64(req.Mark.Pos.Z)
			player := g.userManager.GetTargetUser(userId)
			player.Pos.X = x
			player.Pos.Y = y
			player.Pos.Z = z
			oldSceneId := player.SceneId
			player.SceneId = uint16(req.Mark.SceneId)
			g.log.Info("player goto scene: %v, pos x: %v, y: %v, z: %v", player.SceneId, x, y, z)
			g.userManager.UpdateUser(player)
			// PacketSceneEntityDisappearNotify
			sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
			sceneEntityDisappearNotify.EntityList = []uint32{player.AvatarEntityId}
			sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_REMOVE
			g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)
			entityIdTypeConst := constant.GetEntityIdTypeConst()
			player.AvatarEntityId = g.worldManager.GetNextWorldEntityID(entityIdTypeConst.AVATAR)
			player.WeaponEntityId = g.worldManager.GetNextWorldEntityID(entityIdTypeConst.WEAPON)
			g.userManager.UpdateUser(player)
			// PacketPlayerEnterSceneNotify
			playerPropertyConst := constant.GetPlayerPropertyConst()
			player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
			playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
			playerEnterSceneNotify.PrevSceneId = uint32(oldSceneId)
			playerEnterSceneNotify.PrevPos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
			playerEnterSceneNotify.SceneId = uint32(player.SceneId)
			playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
			playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
			if player.SceneId == oldSceneId {
				playerEnterSceneNotify.Type = proto.EnterType_ENTER_GOTO
			} else {
				playerEnterSceneNotify.Type = proto.EnterType_ENTER_JUMP
			}
			playerEnterSceneNotify.TargetUid = player.PlayerID
			playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
			playerEnterSceneNotify.WorldLevel = player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
			enterReasonConst := constant.GetEnterReasonConst()
			playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.TransPoint)
			playerEnterSceneNotify.SceneTagIdList = []uint32{102, 107, 109, 113, 117}
			playerEnterSceneNotify.WorldType = 1
			playerEnterSceneNotify.SceneTransaction = strconv.FormatInt(int64(player.SceneId), 10) + "-" + strconv.FormatInt(int64(player.PlayerID), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + "18402"
			g.SendMsg(api.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)
			g.userManager.UpdateUser(player)
		}
	}
}
