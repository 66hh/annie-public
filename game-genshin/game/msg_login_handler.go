package game

import (
	"flswld.com/common/utils/random"
	"flswld.com/common/utils/reflection"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"game-genshin/game/constant"
	"strconv"
	"time"
)

func (g *GameManager) OnLoginOk(userId uint32) {
	g.log.Info("user login, user id: %v", userId)
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		g.SendMsg(api.ApiDoSetPlayerBornDataNotify, userId, nil, nil)
		return
	}
	// 登陆成功
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.Properties[playerPropertyConst.PROP_PLAYER_MP_SETTING_TYPE] = uint32(player.MpSetting.Number())
	player.Properties[playerPropertyConst.PROP_IS_MP_MODE_AVAILABLE] = 1
	g.userManager.UpdateUser(player)
	// PacketPlayerDataNotify
	playerDataNotify := new(proto.PlayerDataNotify)
	playerDataNotify.NickName = player.NickName
	playerDataNotify.ServerTime = uint64(time.Now().UnixMilli())
	playerDataNotify.IsFirstLoginToday = true
	playerDataNotify.RegionId = uint32(player.RegionId)
	playerDataNotify.PropMap = make(map[uint32]*proto.PropValue)
	for k, v := range player.Properties {
		propValue := new(proto.PropValue)
		propValue.Type = uint32(k)
		propValue.Value = &proto.PropValue_Ival{Ival: int64(v)}
		propValue.Val = int64(v)
		playerDataNotify.PropMap[uint32(k)] = propValue
	}
	g.SendMsg(api.ApiPlayerDataNotify, userId, g.getHeadMsg(2), playerDataNotify)
	// PacketStoreWeightLimitNotify
	storeWeightLimitNotify := new(proto.StoreWeightLimitNotify)
	storeWeightLimitNotify.StoreType = proto.StoreType_STORE_PACK
	// TODO 原神背包容量限制 写到配置文件
	storeWeightLimitNotify.WeightLimit = 30000
	storeWeightLimitNotify.WeaponCountLimit = 2000
	storeWeightLimitNotify.ReliquaryCountLimit = 2000
	storeWeightLimitNotify.MaterialCountLimit = 2000
	storeWeightLimitNotify.FurnitureCountLimit = 2000
	g.SendMsg(api.ApiStoreWeightLimitNotify, userId, nil, storeWeightLimitNotify)
	// PacketPlayerStoreNotify
	playerStoreNotify := new(proto.PlayerStoreNotify)
	playerStoreNotify.StoreType = proto.StoreType_STORE_PACK
	playerStoreNotify.WeightLimit = 30000
	// TODO 建立玩家背包道具数据结构
	playerStoreNotify.ItemList = append(playerStoreNotify.ItemList, &proto.Item{
		ItemId: 11509,
		Guid:   429496733894967298,
		Detail: &proto.Item_Equip{
			Equip: &proto.Equip{
				Detail: &proto.Equip_Weapon{
					Weapon: &proto.Weapon{
						Level:        1,
						Exp:          0,
						PromoteLevel: 0,
						AffixMap:     nil,
					},
				},
				IsLocked: false,
			},
		},
	})
	g.SendMsg(api.ApiPlayerStoreNotify, userId, g.getHeadMsg(2), playerStoreNotify)
	// PacketAvatarDataNotify
	avatarDataNotify := new(proto.AvatarDataNotify)
	// TODO 建立玩家队伍管理器
	avatarDataNotify.CurAvatarTeamId = 1
	avatarDataNotify.ChooseAvatarGuid = 429496733894967297
	avatarDataNotify.OwnedFlycloakList = player.FlyCloakList
	// TODO 暂时不知道这是什么鬼东西
	avatarDataNotify.OwnedCostumeList = make([]uint32, 0)
	// TODO 建立玩家角色管理器
	avatarDataNotify.AvatarList = []*proto.AvatarInfo{{
		AvatarId: 10000007,
		Guid:     429496733894967297,
		PropMap: map[uint32]*proto.PropValue{
			1001: {Type: 1001, Value: &proto.PropValue_Ival{Ival: 0}},
			1002: {Type: 1002, Value: &proto.PropValue_Ival{Ival: 0}},
			1003: {Type: 1003, Value: &proto.PropValue_Ival{Ival: 0}},
			1004: {Type: 1004, Value: &proto.PropValue_Ival{Ival: 0}},
			4001: {Type: 4001, Val: 1, Value: &proto.PropValue_Ival{Ival: 1}},
		},
		LifeState:     1,
		EquipGuidList: []uint64{429496733894967298},
		FightPropMap: map[uint32]float32{
			0:    0.0,
			1:    911.791,
			4:    41.053,
			7:    57.225,
			20:   0.05,
			22:   0.5,
			23:   1.0,
			1010: 911.791,
			2000: 911.791,
			2001: 41.053,
			2002: 57.225,
		},
		SkillDepotId: 704,
		FetterInfo: &proto.AvatarFetterInfo{
			ExpLevel: 1,
			// FetterList 不知道是啥 该角色在配置表里的所有FetterId
		},
		SkillLevelMap: map[uint32]uint32{
			10067:  1,
			10068:  1,
			100553: 1,
		},
		AvatarType:        1,
		WearingFlycloakId: 140001,
		BornTime:          1652555787,
	}}
	avatarDataNotify.AvatarList[0].FetterInfo.FetterList = make([]*proto.FetterData, 0)
	for _, v := range g.gameDataConfig.AvatarFetterDataMap[10000007] {
		avatarDataNotify.AvatarList[0].FetterInfo.FetterList = append(avatarDataNotify.AvatarList[0].FetterInfo.FetterList, &proto.FetterData{
			FetterId:    v,
			FetterState: uint32(constant.GetFetterStateConst().FINISH),
		})
	}
	avatarDataNotify.AvatarTeamMap = make(map[uint32]*proto.AvatarTeam)
	avatarDataNotify.AvatarTeamMap[1] = &proto.AvatarTeam{
		AvatarGuidList: []uint64{429496733894967297},
		TeamName:       "",
	}
	avatarDataNotify.AvatarTeamMap[2] = &proto.AvatarTeam{}
	avatarDataNotify.AvatarTeamMap[3] = &proto.AvatarTeam{}
	avatarDataNotify.AvatarTeamMap[4] = &proto.AvatarTeam{}
	g.SendMsg(api.ApiAvatarDataNotify, userId, g.getHeadMsg(2), avatarDataNotify)
	// PacketPlayerEnterSceneNotify
	player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
	playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
	playerEnterSceneNotify.SceneId = uint32(player.SceneId)
	playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
	playerEnterSceneNotify.Type = proto.EnterType_ENTER_SELF
	playerEnterSceneNotify.TargetUid = player.PlayerID
	playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneNotify.WorldLevel = player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.Login)
	// TODO 这个要留意一下JAVA那边的实现
	playerEnterSceneNotify.IsFirstLoginEnterScene = true
	playerEnterSceneNotify.WorldType = 1
	playerEnterSceneNotify.SceneTransaction = "3-" + strconv.FormatInt(int64(player.PlayerID), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + "18402"
	g.SendMsg(api.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)
	g.userManager.UpdateUser(player)
	// PacketOpenStateUpdateNotify
	openStateUpdateNotify := new(proto.OpenStateUpdateNotify)
	openStateConst := constant.GetOpenStateConst()
	openStateConstMap := reflection.ConvStructToMap(openStateConst)
	openStateUpdateNotify.OpenStateMap = make(map[uint32]uint32)
	for _, v := range openStateConstMap {
		openStateUpdateNotify.OpenStateMap[uint32(v.(uint16))] = 1
	}
	g.SendMsg(api.ApiOpenStateUpdateNotify, userId, nil, openStateUpdateNotify)
}

func (g *GameManager) OnUserOffline(userId uint32) {
	g.log.Info("user offline, user id: %v", userId)
	player := g.userManager.GetTargetUser(userId)
	g.userManager.UpdateUser(player)
}
