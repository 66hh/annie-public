package game

import (
	"flswld.com/common/utils/reflection"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/constant"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
	"time"
)

func (g *GameManager) OnLogin(userId uint32) {
	logger.LOG.Info("user login, user id: %v", userId)
	player, asyncWait := g.userManager.OnlineUser(userId)
	if !asyncWait {
		g.OnLoginOk(userId, player)
	}
}

func (g *GameManager) OnLoginOk(userId uint32, player *model.Player) {
	if player == nil {
		g.SendMsg(proto.ApiDoSetPlayerBornDataNotify, userId, nil, new(proto.NullMsg))
		return
	}
	{
		// TODO 3.0.0REL版本目前存在当前队伍活跃角色非主角登录进不去场景的情况
		activeAvatarId := player.TeamConfig.GetActiveAvatarId()
		if activeAvatarId != player.MainCharAvatarId {
			activeTeam := player.TeamConfig.GetActiveTeam()
			mainCharIndex := player.TeamConfig.CurrAvatarIndex
			for index, avatarId := range activeTeam.AvatarIdList {
				if avatarId == player.MainCharAvatarId {
					mainCharIndex = uint8(index)
				}
			}
			activeTeam.AvatarIdList[mainCharIndex] = player.MainCharAvatarId
			player.TeamConfig.CurrAvatarIndex = mainCharIndex
		}
	}
	// 创建世界
	player.Online = true
	world := g.worldManager.CreateWorld(player, false)
	world.AddPlayer(player, player.SceneId)
	player.WorldId = world.id
	// 初始化
	player.InitAll()
	playerPropertyConst := constant.GetPlayerPropertyConst()
	player.Properties[playerPropertyConst.PROP_PLAYER_MP_SETTING_TYPE] = uint32(player.MpSetting.Number())
	player.Properties[playerPropertyConst.PROP_IS_MP_MODE_AVAILABLE] = 1
	//g.userManager.UpdateUser(player)
	player.TeamConfig.UpdateTeam()
	scene := world.GetSceneById(player.SceneId)
	scene.UpdatePlayerTeamEntity(player)

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
	g.SendMsg(proto.ApiPlayerDataNotify, userId, nil, playerDataNotify)

	// PacketStoreWeightLimitNotify
	storeWeightLimitNotify := new(proto.StoreWeightLimitNotify)
	storeWeightLimitNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	// TODO 原神背包容量限制 写到配置文件
	storeWeightLimitNotify.WeightLimit = 30000
	storeWeightLimitNotify.WeaponCountLimit = 2000
	storeWeightLimitNotify.ReliquaryCountLimit = 2000
	storeWeightLimitNotify.MaterialCountLimit = 2000
	storeWeightLimitNotify.FurnitureCountLimit = 2000
	g.SendMsg(proto.ApiStoreWeightLimitNotify, userId, nil, storeWeightLimitNotify)

	// PacketPlayerStoreNotify
	playerStoreNotify := new(proto.PlayerStoreNotify)
	playerStoreNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	playerStoreNotify.WeightLimit = 30000
	itemDataMapConfig := gdc.CONF.ItemDataMap
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
		itemDataConfig := itemDataMapConfig[int32(item.ItemId)]
		if itemDataConfig != nil && itemDataConfig.ItemEnumType == itemTypeConst.ITEM_FURNITURE {
			pbItem.Detail = &proto.Item_Furniture{
				Furniture: &proto.Furniture{
					Count: item.Count,
				},
			}
		} else {
			pbItem.Detail = &proto.Item_Material{
				Material: &proto.Material{
					Count:      item.Count,
					DeleteInfo: nil,
				},
			}
		}
		playerStoreNotify.ItemList = append(playerStoreNotify.ItemList, pbItem)
	}
	g.SendMsg(proto.ApiPlayerStoreNotify, userId, nil, playerStoreNotify)

	// PacketAvatarDataNotify
	avatarDataNotify := new(proto.AvatarDataNotify)
	chooseAvatarId := player.TeamConfig.GetActiveAvatarId()
	avatarDataNotify.CurAvatarTeamId = uint32(player.TeamConfig.GetActiveTeamId())
	avatarDataNotify.ChooseAvatarGuid = player.AvatarMap[chooseAvatarId].Guid
	avatarDataNotify.OwnedFlycloakList = player.FlyCloakList
	// 角色衣装
	avatarDataNotify.OwnedCostumeList = player.CostumeList
	for _, avatar := range player.AvatarMap {
		pbAvatar := g.PacketAvatarInfo(avatar)
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
	g.SendMsg(proto.ApiAvatarDataNotify, userId, nil, avatarDataNotify)

	player.BornInScene = false

	// PacketPlayerEnterSceneNotify
	playerEnterSceneNotify := g.PacketPlayerEnterSceneNotify(player)
	g.SendMsg(proto.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)

	//g.userManager.UpdateUser(player)

	// PacketOpenStateUpdateNotify
	openStateUpdateNotify := new(proto.OpenStateUpdateNotify)
	openStateConst := constant.GetOpenStateConst()
	openStateConstMap := reflection.ConvStructToMap(openStateConst)
	openStateUpdateNotify.OpenStateMap = make(map[uint32]uint32)
	for _, v := range openStateConstMap {
		openStateUpdateNotify.OpenStateMap[uint32(v.(uint16))] = 1
	}
	g.SendMsg(proto.ApiOpenStateUpdateNotify, userId, nil, openStateUpdateNotify)
}

func (g *GameManager) SetPlayerBornDataReq(userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message) {
	logger.LOG.Debug("user set born data, user id: %v", userId)
	if headMsg != nil {
		logger.LOG.Debug("client sequence id: %v", headMsg.ClientSequenceId)
	}
	if payloadMsg == nil {
		return
	}
	req := payloadMsg.(*proto.SetPlayerBornDataReq)
	logger.LOG.Debug("avatar id: %v, nickname: %v", req.AvatarId, req.NickName)

	exist, asyncWait := g.userManager.CheckUserExistOnReg(userId, req)
	if !asyncWait {
		g.PlayerReg(exist, req, userId)
	}
}

func (g *GameManager) PlayerReg(exist bool, req *proto.SetPlayerBornDataReq, userId uint32) {
	if exist {
		logger.LOG.Error("recv set born data req, but user is already exist, userId: %v", userId)
		return
	}

	mainCharAvatarId := req.GetAvatarId()
	if mainCharAvatarId != 10000005 && mainCharAvatarId != 10000007 {
		logger.LOG.Error("invalid main char avatar id: %v", mainCharAvatarId)
		return
	}

	player := new(model.Player)
	player.PlayerID = userId
	player.NickName = req.NickName
	player.Signature = "惟愿时光记忆，一路繁花千树。"
	player.MainCharAvatarId = mainCharAvatarId
	player.HeadImage = mainCharAvatarId
	player.NameCard = 210001
	player.NameCardList = make([]uint32, 0)
	player.NameCardList = append(player.NameCardList, 210001, 210042)

	player.FriendList = make([]uint32, 0)
	player.FriendApplyList = make([]uint32, 0)

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
	player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL] = 0
	player.Properties[playerPropertyConst.PROP_IS_SPRING_AUTO_USE] = 1
	player.Properties[playerPropertyConst.PROP_SPRING_AUTO_USE_PERCENT] = 100
	player.Properties[playerPropertyConst.PROP_IS_FLYABLE] = 1
	player.Properties[playerPropertyConst.PROP_IS_TRANSFERABLE] = 1
	player.Properties[playerPropertyConst.PROP_MAX_STAMINA] = 24000
	player.Properties[playerPropertyConst.PROP_CUR_PERSIST_STAMINA] = 24000
	player.Properties[playerPropertyConst.PROP_PLAYER_RESIN] = 160

	player.FlyCloakList = make([]uint32, 0)
	player.FlyCloakList = append(player.FlyCloakList, 140001)

	player.CostumeList = make([]uint32, 0)

	player.Pos = &model.Vector{X: 2747, Y: 194, Z: -1719}
	player.Rot = &model.Vector{X: 0, Y: 307, Z: 0}

	player.MpSetting = proto.MpSettingType_MP_SETTING_TYPE_ENTER_AFTER_APPLY

	player.ItemMap = make(map[uint32]*model.Item)
	player.WeaponMap = make(map[uint64]*model.Weapon)
	player.ReliquaryMap = make(map[uint64]*model.Reliquary)
	player.AvatarMap = make(map[uint32]*model.Avatar)

	player.GameObjectGuidMap = make(map[uint64]model.GameObject)

	// 选哥哥的福报
	if mainCharAvatarId == 10000005 {
		// 添加所有角色
		allAvatarDataConfig := g.GetAllAvatarDataConfig()
		for avatarId, avatarDataConfig := range allAvatarDataConfig {
			player.AddAvatar(uint32(avatarId))
			// 添加初始武器
			weaponId := uint64(g.snowflake.GenId())
			player.AddWeapon(uint32(avatarDataConfig.InitialWeapon), weaponId)
			// 角色装上初始武器
			player.WearWeapon(uint32(avatarId), weaponId)
		}
		// 添加所有武器
		allWeaponDataConfig := g.GetAllWeaponDataConfig()
		for itemId := range allWeaponDataConfig {
			weaponId := uint64(g.snowflake.GenId())
			player.AddWeapon(uint32(itemId), weaponId)
		}
		// 添加所有道具
		allItemDataConfig := g.GetAllItemDataConfig()
		for itemId := range allItemDataConfig {
			player.AddItem(uint32(itemId), 1)
		}
	}

	// 添加选定的主角
	player.AddAvatar(mainCharAvatarId)
	// 添加初始武器
	avatarDataConfig := gdc.CONF.AvatarDataMap[int32(mainCharAvatarId)]
	weaponId := uint64(g.snowflake.GenId())
	player.AddWeapon(uint32(avatarDataConfig.InitialWeapon), weaponId)
	// 角色装上初始武器
	player.WearWeapon(mainCharAvatarId, weaponId)

	player.TeamConfig = model.NewTeamInfo()
	player.TeamConfig.AddAvatarToTeam(mainCharAvatarId, 0)

	player.DropInfo = model.NewDropInfo()

	g.userManager.AddUser(player)

	g.SendMsg(proto.ApiSetPlayerBornDataRsp, userId, nil, new(proto.NullMsg))
	g.OnLogin(userId)
}

func (g *GameManager) PlayerForceExitReq(userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message) {
	// 告诉网关断开玩家的连接
	g.SendMsg(proto.ApiPlayerForceExitRsp, userId, nil, new(proto.NullMsg))
	go func() {
		time.Sleep(time.Second)
		g.KickPlayer(userId)
	}()
}

func (g *GameManager) OnUserOffline(userId uint32) {
	logger.LOG.Info("user offline, user id: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world != nil {
		g.UserWorldRemovePlayer(world, player)
	}
	player.OfflineTime = uint32(time.Now().Unix())
	player.Online = false
	g.userManager.OfflineUser(player)
}
