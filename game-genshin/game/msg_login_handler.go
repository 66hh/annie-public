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

	// 创建世界
	world := g.worldManager.CreateWorld(player)
	world.AddPlayer(player, player.SceneId)
	player.WorldId = world.id
	player.PeerId = world.GetNextWorldPeerId()
	// 初始化
	player.InitAll(g.gameDataConfig)
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.Properties[playerPropertyConst.PROP_PLAYER_MP_SETTING_TYPE] = uint32(player.MpSetting.Number())
	player.Properties[playerPropertyConst.PROP_IS_MP_MODE_AVAILABLE] = 1
	g.userManager.UpdateUser(player)
	player.TeamConfig.UpdateTeam(world.GetNextWorldEntityId, g.gameDataConfig.AvatarSkillDepotDataMap)

	// PacketPlayerDataNotify
	playerDataNotify := new(proto.PlayerDataNotify)
	playerDataNotify.NickName = player.NickName
	playerDataNotify.ServerTime = uint64(time.Now().UnixMilli())
	playerDataNotify.IsFirstLoginToday = true
	playerDataNotify.RegionId = player.RegionId
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
	storeWeightLimitNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	// TODO 原神背包容量限制 写到配置文件
	storeWeightLimitNotify.WeightLimit = 30000
	storeWeightLimitNotify.WeaponCountLimit = 2000
	storeWeightLimitNotify.ReliquaryCountLimit = 2000
	storeWeightLimitNotify.MaterialCountLimit = 2000
	storeWeightLimitNotify.FurnitureCountLimit = 2000
	g.SendMsg(api.ApiStoreWeightLimitNotify, userId, nil, storeWeightLimitNotify)

	// PacketPlayerStoreNotify
	playerStoreNotify := new(proto.PlayerStoreNotify)
	playerStoreNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	playerStoreNotify.WeightLimit = 30000
	itemDataMapConfig := g.gameDataConfig.ItemDataMap
	itemTypeConst := constant.GetItemTypeConst()
	for _, weapon := range player.WeaponMap {
		pbItem := &proto.Item{
			ItemId: weapon.ItemId,
			Guid:   weapon.Guid,
			Detail: nil,
		}
		if itemDataMapConfig[int32(weapon.ItemId)].ItemEnumType != itemTypeConst.ITEM_WEAPON {
			continue
		}
		affixMap := make(map[uint32]uint32)
		for _, affixId := range weapon.AffixIdList {
			affixMap[affixId] = uint32(weapon.Refinement)
		}
		pbItem.Detail = &proto.Item_Equip{
			Equip: &proto.Equip{
				Detail: &proto.Equip_Weapon{
					Weapon: &proto.Weapon{
						Level:        uint32(weapon.Level),
						Exp:          weapon.Exp,
						PromoteLevel: uint32(weapon.Promote),
						AffixMap:     affixMap,
					},
				},
				IsLocked: weapon.Lock,
			},
		}
		playerStoreNotify.ItemList = append(playerStoreNotify.ItemList, pbItem)
	}
	for _, reliquary := range player.ReliquaryMap {
		pbItem := &proto.Item{
			ItemId: reliquary.ItemId,
			Guid:   reliquary.Guid,
			Detail: nil,
		}
		if itemDataMapConfig[int32(reliquary.ItemId)].ItemEnumType != itemTypeConst.ITEM_RELIQUARY {
			continue
		}
		pbItem.Detail = &proto.Item_Equip{
			Equip: &proto.Equip{
				Detail: &proto.Equip_Reliquary{
					Reliquary: &proto.Reliquary{
						Level:        uint32(reliquary.Level),
						Exp:          reliquary.Exp,
						PromoteLevel: uint32(reliquary.Promote),
						MainPropId:   reliquary.MainPropId,
						// TODO 圣遗物副词条
						AppendPropIdList: nil,
					},
				},
				IsLocked: reliquary.Lock,
			},
		}
		playerStoreNotify.ItemList = append(playerStoreNotify.ItemList, pbItem)
	}
	for _, item := range player.ItemMap {
		pbItem := &proto.Item{
			ItemId: item.ItemId,
			Guid:   item.Guid,
			Detail: nil,
		}
		switch itemDataMapConfig[int32(item.ItemId)].ItemEnumType {
		case itemTypeConst.ITEM_FURNITURE:
			pbItem.Detail = &proto.Item_Furniture{
				Furniture: &proto.Furniture{
					Count: item.Count,
				},
			}
		default:
			pbItem.Detail = &proto.Item_Material{
				Material: &proto.Material{
					Count:      item.Count,
					DeleteInfo: nil,
				},
			}
		}
		playerStoreNotify.ItemList = append(playerStoreNotify.ItemList, pbItem)
	}
	g.SendMsg(api.ApiPlayerStoreNotify, userId, g.getHeadMsg(2), playerStoreNotify)

	// PacketAvatarDataNotify
	avatarDataNotify := new(proto.AvatarDataNotify)
	chooseAvatarId := player.TeamConfig.GetActiveAvatarId()
	avatarDataNotify.CurAvatarTeamId = uint32(player.TeamConfig.GetActiveTeamId())
	avatarDataNotify.ChooseAvatarGuid = player.AvatarMap[chooseAvatarId].Guid
	avatarDataNotify.OwnedFlycloakList = player.FlyCloakList
	// 角色衣装
	avatarDataNotify.OwnedCostumeList = player.CostumeList
	fetterStateConst := constant.GetFetterStateConst()
	for avatarId, avatar := range player.AvatarMap {
		pbAvatar := &proto.AvatarInfo{
			AvatarId: avatar.AvatarId,
			Guid:     avatar.Guid,
			PropMap: map[uint32]*proto.PropValue{
				uint32(playerPropertyConst.PROP_LEVEL): {
					Type:  uint32(playerPropertyConst.PROP_LEVEL),
					Value: &proto.PropValue_Ival{Ival: int64(avatar.Level)},
				},
				uint32(playerPropertyConst.PROP_EXP): {
					Type:  uint32(playerPropertyConst.PROP_EXP),
					Value: &proto.PropValue_Ival{Ival: int64(avatar.Exp)},
				},
				uint32(playerPropertyConst.PROP_BREAK_LEVEL): {
					Type:  uint32(playerPropertyConst.PROP_BREAK_LEVEL),
					Value: &proto.PropValue_Ival{Ival: int64(avatar.Promote)},
				},
				uint32(playerPropertyConst.PROP_SATIATION_VAL): {
					Type:  uint32(playerPropertyConst.PROP_SATIATION_VAL),
					Value: &proto.PropValue_Ival{Ival: 0},
				},
				uint32(playerPropertyConst.PROP_SATIATION_PENALTY_TIME): {
					Type:  uint32(playerPropertyConst.PROP_SATIATION_PENALTY_TIME),
					Value: &proto.PropValue_Ival{Ival: 0},
				},
			},
			LifeState:     1,
			EquipGuidList: avatar.EquipGuidList,
			FightPropMap:  nil,
			SkillDepotId:  avatar.SkillDepotId,
			FetterInfo: &proto.AvatarFetterInfo{
				ExpLevel:  uint32(avatar.FetterLevel),
				ExpNumber: avatar.FetterExp,
				// FetterList 不知道是啥 该角色在配置表里的所有FetterId
				// TODO 资料解锁条目
				FetterList:              nil,
				RewardedFetterLevelList: []uint32{10},
			},
			SkillLevelMap:     nil,
			AvatarType:        1,
			WearingFlycloakId: avatar.FlyCloak,
			BornTime:          uint32(avatar.BornTime),
		}

		player.AvatarMap[avatarId] = avatar
		pbAvatar.FightPropMap = avatar.FightPropMap
		for _, v := range avatar.FetterList {
			pbAvatar.FetterInfo.FetterList = append(pbAvatar.FetterInfo.FetterList, &proto.FetterData{
				FetterId:    v,
				FetterState: uint32(fetterStateConst.FINISH),
			})
		}
		// 解锁全部资料
		for _, v := range g.gameDataConfig.AvatarFetterDataMap[int32(avatar.AvatarId)] {
			pbAvatar.FetterInfo.FetterList = append(pbAvatar.FetterInfo.FetterList, &proto.FetterData{
				FetterId:    uint32(v),
				FetterState: uint32(fetterStateConst.FINISH),
			})
		}
		pbAvatar.SkillLevelMap = make(map[uint32]uint32)
		for k, v := range avatar.SkillLevelMap {
			pbAvatar.SkillLevelMap[k] = v
		}
		avatarDataNotify.AvatarList = append(avatarDataNotify.AvatarList, pbAvatar)
	}
	avatarDataNotify.AvatarTeamMap = make(map[uint32]*proto.AvatarTeam)
	for teamIndex, team := range player.TeamConfig.TeamList {
		var teamAvatarGuidList []uint64 = nil
		for _, avatarId := range team.AvatarIdList {
			if avatarId == 0 {
				break
			}
			teamAvatarGuidList = append(teamAvatarGuidList, player.AvatarMap[avatarId].Guid)
		}
		avatarDataNotify.AvatarTeamMap[uint32(teamIndex)+1] = &proto.AvatarTeam{
			AvatarGuidList: teamAvatarGuidList,
			TeamName:       team.Name,
		}
	}
	g.SendMsg(api.ApiAvatarDataNotify, userId, g.getHeadMsg(2), avatarDataNotify)

	// PacketPlayerEnterSceneNotify
	player.EnterSceneToken = uint32(random.GetRandomInt32(1000, 99999))
	playerEnterSceneNotify := new(proto.PlayerEnterSceneNotify)
	playerEnterSceneNotify.SceneId = player.SceneId
	playerEnterSceneNotify.Pos = &proto.Vector{X: float32(player.Pos.X), Y: float32(player.Pos.Y), Z: float32(player.Pos.Z)}
	playerEnterSceneNotify.SceneBeginTime = uint64(time.Now().UnixMilli())
	playerEnterSceneNotify.Type = proto.EnterType_ENTER_TYPE_SELF
	playerEnterSceneNotify.TargetUid = player.PlayerID
	playerEnterSceneNotify.EnterSceneToken = player.EnterSceneToken
	playerEnterSceneNotify.WorldLevel = player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify.EnterReason = uint32(enterReasonConst.Login)
	// TODO 刚登录进入场景的时候才为true
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
