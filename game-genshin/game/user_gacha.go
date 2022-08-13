package game

import (
	"flswld.com/common/config"
	"flswld.com/common/utils/random"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type UserInfo struct {
	UserId uint32 `json:"userId"`
	jwt.RegisteredClaims
}

// 获取卡池信息
func (g *GameManager) GetGachaInfoReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get gacha info, userId: %v", userId)
	serverAddr := config.CONF.Genshin.GachaHistoryServer
	getGachaInfoRsp := new(proto.GetGachaInfoRsp)
	getGachaInfoRsp.GachaRandom = 12345
	userInfo := &UserInfo{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userInfo)
	jwtStr, err := token.SignedString([]byte("flswld"))
	if err != nil {
		logger.LOG.Error("generate jwt error: %v", err)
		jwtStr = "default.jwt.token"
	}
	getGachaInfoRsp.GachaInfoList = []*proto.GachaInfo{
		// 温迪
		{
			GachaType:              300,
			ScheduleId:             823,
			BeginTime:              0,
			EndTime:                2051193600,
			GachaSortId:            9998,
			GachaPrefabPath:        "GachaShowPanel_A019",
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A019",
			TitleTextmap:           "UI_GACHA_SHOW_PANEL_A019_TITLE",
			LeftGachaTimes:         2147483647,
			GachaTimesLimit:        2147483647,
			CostItemId:             223,
			CostItemNum:            1,
			TenCostItemId:          223,
			TenCostItemNum:         10,
			GachaRecordUrl:         serverAddr + "/gm/gacha?gachaType=300&jwt=" + jwtStr,
			GachaRecordUrlOversea:  serverAddr + "/gm/gacha?gachaType=300&jwt=" + jwtStr,
			GachaProbUrl:           serverAddr + "/gm/gacha/details?scheduleId=823&jwt=" + jwtStr,
			GachaProbUrlOversea:    serverAddr + "/gm/gacha/details?scheduleId=823&jwt=" + jwtStr,
			GachaUpInfoList: []*proto.GachaUpInfo{
				{
					ItemParentType: 1,
					ItemIdList:     []uint32{1022},
				},
				{
					ItemParentType: 2,
					ItemIdList:     []uint32{1023, 1031, 1014},
				},
			},
			DisplayUp_4ItemList: []uint32{1023},
			DisplayUp_5ItemList: []uint32{1022},
			WishItemId:          0,
			WishProgress:        0,
			WishMaxProgress:     0,
			IsNewWish:           false,
		},
		// 可莉
		{
			GachaType:              400,
			ScheduleId:             833,
			BeginTime:              0,
			EndTime:                2051193600,
			GachaSortId:            9998,
			GachaPrefabPath:        "GachaShowPanel_A018",
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A018",
			TitleTextmap:           "UI_GACHA_SHOW_PANEL_A018_TITLE",
			LeftGachaTimes:         2147483647,
			GachaTimesLimit:        2147483647,
			CostItemId:             223,
			CostItemNum:            1,
			TenCostItemId:          223,
			TenCostItemNum:         10,
			GachaRecordUrl:         serverAddr + "/gm/gacha?gachaType=400&jwt=" + jwtStr,
			GachaRecordUrlOversea:  serverAddr + "/gm/gacha?gachaType=400&jwt=" + jwtStr,
			GachaProbUrl:           serverAddr + "/gm/gacha/details?scheduleId=833&jwt=" + jwtStr,
			GachaProbUrlOversea:    serverAddr + "/gm/gacha/details?scheduleId=833&jwt=" + jwtStr,
			GachaUpInfoList: []*proto.GachaUpInfo{
				{
					ItemParentType: 1,
					ItemIdList:     []uint32{1029},
				},
				{
					ItemParentType: 2,
					ItemIdList:     []uint32{1025, 1034, 1043},
				},
			},
			DisplayUp_4ItemList: []uint32{1025},
			DisplayUp_5ItemList: []uint32{1029},
			WishItemId:          0,
			WishProgress:        0,
			WishMaxProgress:     0,
			IsNewWish:           false,
		},
		// 阿莫斯之弓&风鹰剑
		{
			GachaType:              426,
			ScheduleId:             1103,
			BeginTime:              0,
			EndTime:                2051193600,
			GachaSortId:            9997,
			GachaPrefabPath:        "GachaShowPanel_A020",
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A020",
			TitleTextmap:           "UI_GACHA_SHOW_PANEL_A020_TITLE",
			LeftGachaTimes:         2147483647,
			GachaTimesLimit:        2147483647,
			CostItemId:             223,
			CostItemNum:            1,
			TenCostItemId:          223,
			TenCostItemNum:         10,
			GachaRecordUrl:         serverAddr + "/gm/gacha?gachaType=426&jwt=" + jwtStr,
			GachaRecordUrlOversea:  serverAddr + "/gm/gacha?gachaType=426&jwt=" + jwtStr,
			GachaProbUrl:           serverAddr + "/gm/gacha/details?scheduleId=1103&jwt=" + jwtStr,
			GachaProbUrlOversea:    serverAddr + "/gm/gacha/details?scheduleId=1103&jwt=" + jwtStr,
			GachaUpInfoList: []*proto.GachaUpInfo{
				{
					ItemParentType: 1,
					ItemIdList:     []uint32{15502, 11501},
				},
				{
					ItemParentType: 2,
					ItemIdList:     []uint32{11402, 12402, 15402, 11402, 13407},
				},
			},
			DisplayUp_4ItemList: []uint32{11402},
			DisplayUp_5ItemList: []uint32{15502, 11501},
			WishItemId:          0,
			WishProgress:        0,
			WishMaxProgress:     2,
			IsNewWish:           false,
		},
		// 常驻
		{
			GachaType:              201,
			ScheduleId:             813,
			BeginTime:              0,
			EndTime:                2051193600,
			GachaSortId:            1000,
			GachaPrefabPath:        "GachaShowPanel_A017",
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A017",
			TitleTextmap:           "UI_GACHA_SHOW_PANEL_A017_TITLE",
			LeftGachaTimes:         2147483647,
			GachaTimesLimit:        2147483647,
			CostItemId:             224,
			CostItemNum:            1,
			TenCostItemId:          224,
			TenCostItemNum:         10,
			GachaRecordUrl:         serverAddr + "/gm/gacha?gachaType=201&jwt=" + jwtStr,
			GachaRecordUrlOversea:  serverAddr + "/gm/gacha?gachaType=201&jwt=" + jwtStr,
			GachaProbUrl:           serverAddr + "/gm/gacha/details?scheduleId=813&jwt=" + jwtStr,
			GachaProbUrlOversea:    serverAddr + "/gm/gacha/details?scheduleId=813&jwt=" + jwtStr,
			GachaUpInfoList: []*proto.GachaUpInfo{
				{
					ItemParentType: 1,
					ItemIdList:     []uint32{1003, 1016},
				},
				{
					ItemParentType: 2,
					ItemIdList:     []uint32{1021, 1006, 1015},
				},
			},
			DisplayUp_4ItemList: []uint32{1021},
			DisplayUp_5ItemList: []uint32{1003, 1016},
			WishItemId:          0,
			WishProgress:        0,
			WishMaxProgress:     0,
			IsNewWish:           false,
		},
	}
	g.SendMsg(api.ApiGetGachaInfoRsp, userId, nil, getGachaInfoRsp)
}

func (g *GameManager) DoGachaReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user do gacha, userId: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	req := payloadMsg.(*proto.DoGachaReq)
	gachaScheduleId := req.GachaScheduleId
	gachaTimes := req.GachaTimes

	poolIndex := 0
	gachaType := uint32(0)
	costItemId := uint32(0)
	switch gachaScheduleId {
	case 823:
		// 温迪
		poolIndex = 1
		gachaType = 300
		costItemId = 223
	case 833:
		// 可莉
		poolIndex = 2
		gachaType = 400
		costItemId = 223
	case 1103:
		// 阿莫斯之弓&风鹰剑
		poolIndex = 3
		gachaType = 426
		costItemId = 223
	case 813:
		// 常驻
		poolIndex = 4
		gachaType = 201
		costItemId = 224
	}
	_ = poolIndex

	// PacketDoGachaRsp
	doGachaRsp := new(proto.DoGachaRsp)
	doGachaRsp.GachaType = gachaType
	doGachaRsp.GachaScheduleId = gachaScheduleId
	doGachaRsp.GachaTimes = gachaTimes
	doGachaRsp.NewGachaRandom = 12345
	doGachaRsp.LeftGachaTimes = 2147483647
	doGachaRsp.GachaTimesLimit = 2147483647
	doGachaRsp.CostItemId = costItemId
	doGachaRsp.CostItemNum = 1
	doGachaRsp.TenCostItemId = costItemId
	doGachaRsp.TenCostItemNum = 10

	// 先扣掉粉球或蓝球再进行抽卡
	g.CostUserItem(player.PlayerID, []*UserItem{
		{
			ItemId:      costItemId,
			ChangeCount: gachaTimes,
		},
	})

	doGachaRsp.GachaItemList = make([]*proto.GachaItem, 0)
	for i := uint32(0); i < gachaTimes; i++ {
		//awardItemId := uint32(0)
		//allWeaponDataConfig := g.GetAllAvatarDataConfig()
		//for k, v := range allWeaponDataConfig {
		//	if v.RankLevel < 5 {
		//		delete(allWeaponDataConfig, k)
		//	}
		//}
		//for k, v := range allWeaponDataConfig {
		//	if v.QualityType == "QUALITY_PURPLE" {
		//		delete(allWeaponDataConfig, k)
		//	} else if v.QualityType == "QUALITY_ORANGE" {
		//	}
		//}
		//rn := random.GetRandomInt32(0, int32(len(allWeaponDataConfig)-1))
		//index := int32(0)
		//for itemId, itemData := range allWeaponDataConfig {
		//	if rn == index {
		//		logger.LOG.Debug("itemData: %v", itemData)
		//		awardItemId = uint32(itemId)
		//	}
		//	index++
		//}
		//_ = awardItemId
		//awardItemId %= 1000
		//awardItemId += 1000

		ok, itemId := g.doGachaOnce(userId, gachaType)
		if !ok {
			itemId = 11301
		}

		gachaItem := new(proto.GachaItem)
		gachaItem.GachaItem_ = &proto.ItemParam{
			ItemId: itemId,
			Count:  1,
		}

		//gachaItem.TokenItemList = []*proto.ItemParam{{
		//	// 星尘
		//	ItemId: 222,
		//	Count:  15,
		//}}
		//gachaItem.TransferItems = []*proto.GachaTransferItem{{
		//	Item: &proto.ItemParam{
		//		// 星辉
		//		ItemId: 221,
		//		Count:  5,
		//	},
		//}}

		doGachaRsp.GachaItemList = append(doGachaRsp.GachaItemList, gachaItem)

		//// 添加抽卡获得的道具
		//g.AddUserItem(player.PlayerID, []*UserItem{
		//	// 星尘
		//	{
		//		ItemId:      222,
		//		ChangeCount: 15,
		//	},
		//	// 星辉
		//	{
		//		ItemId:      221,
		//		ChangeCount: 5,
		//	},
		//}, false)
		//g.AddUserWeapon(player.PlayerID, 13303)
	}

	//logger.LOG.Debug("doGachaRsp: %v", doGachaRsp.String())

	g.SendMsg(api.ApiDoGachaRsp, userId, nil, doGachaRsp)
}

var (
	OrangeTimesFixThreshold uint32 = 74   // 触发5星概率修正阈值的抽卡次数
	OrangeTimesFixValue     int32  = 600  // 5星概率修正因子
	PurpleTimesFixThreshold uint32 = 9    // 触发4星概率修正阈值的抽卡次数
	PurpleTimesFixValue     int32  = 5100 // 4星概率修正因子
)

// 单抽一次
func (g *GameManager) doGachaOnce(userId uint32, gachaType uint32) (bool, uint32) {
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return false, 0
	}

	// 找到卡池对应的掉落组
	dropGroupDataConfig := gdc.CONF.DropGroupDataMap[int32(gachaType)]
	if dropGroupDataConfig == nil {
		logger.LOG.Error("drop group not found, drop id: %v", gachaType)
		return false, 0
	}

	// 获取用户的卡池保底信息
	gachaPoolInfo := player.DropInfo.GachaPoolInfo[gachaType]
	if gachaPoolInfo == nil {
		logger.LOG.Error("user gacha pool info not found, gacha type: %v", gachaType)
		return false, 0
	}

	// 保底计数+1
	gachaPoolInfo.OrangeTimes++
	gachaPoolInfo.PurpleTimes++

	// 4星和5星概率修正
	if gachaPoolInfo.OrangeTimes >= OrangeTimesFixThreshold || gachaPoolInfo.PurpleTimes >= PurpleTimesFixThreshold {
		fixDropGroupDataConfig := new(gdc.DropGroupData)
		fixDropGroupDataConfig.DropId = dropGroupDataConfig.DropId
		fixDropGroupDataConfig.WeightAll = dropGroupDataConfig.WeightAll
		// 计算4星和5星权重修正值
		addOrangeWeight := int32(gachaPoolInfo.OrangeTimes-OrangeTimesFixThreshold+1) * OrangeTimesFixValue
		if addOrangeWeight < 0 {
			addOrangeWeight = 0
		}
		addPurpleWeight := int32(gachaPoolInfo.PurpleTimes-PurpleTimesFixThreshold+1) * PurpleTimesFixValue
		if addPurpleWeight < 0 {
			addPurpleWeight = 0
		}
		for _, drop := range dropGroupDataConfig.DropConfig {
			fixDrop := new(gdc.Drop)
			fixDrop.Result = drop.Result
			fixDrop.DropId = drop.DropId
			fixDrop.IsEnd = drop.IsEnd
			// 要求配置表的5/4/3星掉落组id规则固定为卡池类型*10+1/2/3
			orangeDropId := int32(gachaType*10 + 1)
			purpleDropId := int32(gachaType*10 + 2)
			blueDropId := int32(gachaType*10 + 3)
			// 权重修正
			if drop.Result == orangeDropId {
				fixDrop.Weight = drop.Weight + addOrangeWeight
			} else if drop.Result == purpleDropId {
				fixDrop.Weight = drop.Weight + addPurpleWeight
			} else if drop.Result == blueDropId {
				fixDrop.Weight = drop.Weight - addOrangeWeight - addPurpleWeight
			} else {
				logger.LOG.Error("invalid drop group id, does not match any case of orange/purple/blue, result group id: %v", drop.Result)
				fixDrop.Weight = drop.Weight
			}
			fixDropGroupDataConfig.DropConfig = append(fixDropGroupDataConfig.DropConfig, fixDrop)
		}
		dropGroupDataConfig = fixDropGroupDataConfig
	}

	// 掉落
	ok, drop := g.doFullRandDrop(dropGroupDataConfig)
	if !ok {
		return false, 0
	}
	gachaItemId := uint32(drop.Result)
	if gachaItemId < 2000 {
		// 抽到角色
		avatarId := (gachaItemId % 1000) + 10000000
		allAvatarDataConfig := g.GetAllAvatarDataConfig()
		avatarDataConfig := allAvatarDataConfig[int32(avatarId)]
		if avatarDataConfig == nil {
			logger.LOG.Error("avatar data config not found, avatar id: %v", avatarId)
			return false, 0
		}
		// 重置4星和5星保底计数
		if avatarDataConfig.QualityType == "QUALITY_ORANGE" {
			logger.LOG.Debug("[orange avatar], times: %v, gachaItemId: %v", gachaPoolInfo.OrangeTimes, gachaItemId)
			gachaPoolInfo.OrangeTimes = 0
			// 找到UP的5星对应的掉落组id
			upOrangeDropId := int32(gachaType*100 + 12)
			// 触发大保底
			if gachaPoolInfo.MustGetUpOrange {
				logger.LOG.Debug("trigger must get up orange, user id: %v", userId)
				upOrangeDropGroupDataConfig := gdc.CONF.DropGroupDataMap[upOrangeDropId]
				if upOrangeDropGroupDataConfig == nil {
					logger.LOG.Error("drop group not found, drop id: %v", upOrangeDropId)
					return false, 0
				}
				upOrangeOk, upOrangeDrop := g.doFullRandDrop(upOrangeDropGroupDataConfig)
				if !upOrangeOk {
					return false, 0
				}
				gachaPoolInfo.MustGetUpOrange = false
				upOrangeGachaItemId := uint32(upOrangeDrop.Result)
				return upOrangeOk, upOrangeGachaItemId
			}
			if drop.DropId != upOrangeDropId {
				gachaPoolInfo.MustGetUpOrange = true
			}
		} else if avatarDataConfig.QualityType == "QUALITY_PURPLE" {
			logger.LOG.Debug("[purple avatar], times: %v, gachaItemId: %v", gachaPoolInfo.PurpleTimes, gachaItemId)
			gachaPoolInfo.PurpleTimes = 0
		}
	} else {
		// 抽到武器
		allWeaponDataConfig := g.GetAllWeaponDataConfig()
		weaponDataConfig := allWeaponDataConfig[int32(gachaItemId)]
		if weaponDataConfig == nil {
			logger.LOG.Error("weapon item data config not found, item id: %v", gachaItemId)
			return false, 0
		}
		// 重置4星和5星保底计数
		if weaponDataConfig.RankLevel == 5 {
			logger.LOG.Debug("[orange weapon], times: %v, gachaItemId: %v", gachaPoolInfo.OrangeTimes, gachaItemId)
			gachaPoolInfo.OrangeTimes = 0
		} else if weaponDataConfig.RankLevel == 4 {
			logger.LOG.Debug("[purple weapon], times: %v, gachaItemId: %v", gachaPoolInfo.PurpleTimes, gachaItemId)
			gachaPoolInfo.PurpleTimes = 0
		}
	}
	return ok, gachaItemId
}

// 走一次完整流程的掉落组
func (g *GameManager) doFullRandDrop(dropGroupDataConfig *gdc.DropGroupData) (bool, *gdc.Drop) {
	for {
		drop := g.doRandDropOnce(dropGroupDataConfig)
		if drop == nil {
			logger.LOG.Error("weight error, drop group config: %v", dropGroupDataConfig)
			return false, nil
		}
		if drop.IsEnd {
			// 成功抽到物品
			return true, drop
		}
		// 进行下一步掉落流程
		dropGroupDataConfig = gdc.CONF.DropGroupDataMap[drop.Result]
		if dropGroupDataConfig == nil {
			logger.LOG.Error("drop config tab exist error, invalid drop id: %v", drop.Result)
			return false, nil
		}
	}
}

// 进行单次随机掉落
func (g *GameManager) doRandDropOnce(dropGroupDataConfig *gdc.DropGroupData) *gdc.Drop {
	randNum := random.GetRandomInt32(0, dropGroupDataConfig.WeightAll-1)
	sumWeight := int32(0)
	for _, drop := range dropGroupDataConfig.DropConfig {
		sumWeight += drop.Weight
		if sumWeight > randNum {
			return drop
		}
	}
	return nil
}
