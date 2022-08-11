package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
)

func (g *GameManager) GetGachaInfoReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	getGachaInfoRsp := new(proto.GetGachaInfoRsp)
	getGachaInfoRsp.GachaRandom = 12345
	sUrl := "66b1dc20bcb8e25fa6918fcc0f3d54575714583dea8431a3d42bfd5839c10aba"
	getGachaInfoRsp.GachaInfoList = []*proto.GachaInfo{
		// A
		{
			CostItemNum:            1,
			LeftGachaTimes:         2147483647,
			ScheduleId:             833,
			GachaTimesLimit:        2147483647,
			EndTime:                1924992000,
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A018",
			TenCostItemId:          223,
			GachaRecordUrl:         "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=400",
			CostItemId:             223,
			GachaPrefabPath:        "GachaShowPanel_A018",
			TenCostItemNum:         10,
			GachaType:              400,
			GachaProbUrl:           "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=833",
			GachaSortId:            9998,
			GachaProbUrlOversea:    "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=833",
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
			DisplayUp_4ItemList:   []uint32{1025},
			TitleTextmap:          "UI_GACHA_SHOW_PANEL_A018_TITLE",
			DisplayUp_5ItemList:   []uint32{1029},
			GachaRecordUrlOversea: "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=400",
			// 没填的
			BeginTime:       0,
			WishItemId:      0,
			WishProgress:    0,
			WishMaxProgress: 0,
			IsNewWish:       false,
		},
		// B
		{
			CostItemNum:            1,
			LeftGachaTimes:         2147483647,
			ScheduleId:             1103,
			GachaTimesLimit:        2147483647,
			EndTime:                1924992000,
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A020",
			TenCostItemId:          223,
			GachaRecordUrl:         "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=426",
			CostItemId:             223,
			GachaPrefabPath:        "GachaShowPanel_A020",
			TenCostItemNum:         10,
			GachaType:              426,
			GachaProbUrl:           "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=1103",
			GachaSortId:            9997,
			WishMaxProgress:        2,
			GachaProbUrlOversea:    "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=1103",
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
			DisplayUp_4ItemList:   []uint32{11402},
			TitleTextmap:          "UI_GACHA_SHOW_PANEL_A020_TITLE",
			DisplayUp_5ItemList:   []uint32{15502, 11501},
			GachaRecordUrlOversea: "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=426",
			// 没填的
			BeginTime:    0,
			WishItemId:   0,
			WishProgress: 0,
			IsNewWish:    false,
		},
		// C
		{
			CostItemNum:            1,
			LeftGachaTimes:         2147483647,
			ScheduleId:             813,
			GachaTimesLimit:        2147483647,
			EndTime:                1924992000,
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A017",
			TenCostItemId:          224,
			GachaRecordUrl:         "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=201",
			CostItemId:             224,
			GachaPrefabPath:        "GachaShowPanel_A017",
			TenCostItemNum:         10,
			GachaType:              201,
			GachaProbUrl:           "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=813",
			GachaSortId:            1000,
			GachaProbUrlOversea:    "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=813",
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
			DisplayUp_4ItemList:   []uint32{1021},
			TitleTextmap:          "UI_GACHA_SHOW_PANEL_A017_TITLE",
			DisplayUp_5ItemList:   []uint32{1003, 1016},
			GachaRecordUrlOversea: "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=201",
			// 没填的
			BeginTime:       0,
			WishItemId:      0,
			WishProgress:    0,
			WishMaxProgress: 0,
			IsNewWish:       false,
		},
		// D
		{
			CostItemNum:            1,
			LeftGachaTimes:         2147483647,
			ScheduleId:             823,
			GachaTimesLimit:        2147483647,
			EndTime:                1924992000,
			GachaPreviewPrefabPath: "UI_Tab_GachaShowPanel_A019",
			TenCostItemId:          223,
			GachaRecordUrl:         "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=300",
			CostItemId:             223,
			GachaPrefabPath:        "GachaShowPanel_A019",
			TenCostItemNum:         10,
			GachaType:              300,
			GachaProbUrl:           "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=823",
			GachaSortId:            9998,
			GachaProbUrlOversea:    "https://172.16.2.155:443/gacha/details?s=" + sUrl + "&scheduleId=823",
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
			DisplayUp_4ItemList:   []uint32{1023},
			TitleTextmap:          "UI_GACHA_SHOW_PANEL_A019_TITLE",
			DisplayUp_5ItemList:   []uint32{1022},
			GachaRecordUrlOversea: "https://172.16.2.155:443/gacha?s=" + sUrl + "&gachaType=300",
			// 没填的
			BeginTime:       0,
			WishItemId:      0,
			WishProgress:    0,
			WishMaxProgress: 0,
			IsNewWish:       false,
		},
	}
	g.SendMsg(api.ApiGetGachaInfoRsp, userId, nil, getGachaInfoRsp)
}

func (g *GameManager) DoGachaReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	player := g.userManager.GetTargetUser(userId)
	req := payloadMsg.(*proto.DoGachaReq)
	gachaScheduleId := req.GachaScheduleId
	gachaTimes := req.GachaTimes

	gachaType := uint32(0)
	costItemId := uint32(0)
	switch gachaScheduleId {
	case 833:
		gachaType = 400
		costItemId = 223
	case 1103:
		gachaType = 426
		costItemId = 223
	case 813:
		gachaType = 201
		costItemId = 224
	case 823:
		gachaType = 300
		costItemId = 223
	}

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

	doGachaRsp.GachaItemList = make([]*proto.GachaItem, 0)
	for i := uint32(0); i < gachaTimes; i++ {
		gachaItem := new(proto.GachaItem)
		gachaItem.GachaItem_ = &proto.ItemParam{
			ItemId: 13303,
			Count:  1,
		}
		gachaItem.TokenItemList = []*proto.ItemParam{{
			ItemId: 222,
			Count:  15,
		}}
		gachaItem.TransferItems = []*proto.GachaTransferItem{{
			Item: &proto.ItemParam{
				ItemId: 221,
				Count:  5,
			},
		}}
		doGachaRsp.GachaItemList = append(doGachaRsp.GachaItemList, gachaItem)
		g.AddUserItem(player.PlayerID, []*UserItem{
			{
				ItemId:      222,
				ChangeCount: 15,
			},
			{
				ItemId:      221,
				ChangeCount: 5,
			},
		}, false)
		g.AddUserWeapon(player.PlayerID, 13303)
	}

	g.SendMsg(api.ApiDoGachaRsp, userId, nil, doGachaRsp)
}
