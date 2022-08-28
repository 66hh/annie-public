package game

import (
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/constant"
)

type UserItem struct {
	ItemId      uint32
	ChangeCount uint32
}

func (g *GameManager) GetAllItemDataConfig() map[int32]*gdc.ItemData {
	allItemDataConfig := make(map[int32]*gdc.ItemData)
	itemTypeConst := constant.GetItemTypeConst()
	for itemId, itemData := range gdc.CONF.ItemDataMap {
		if itemData.ItemEnumType == itemTypeConst.ITEM_WEAPON {
			// 排除武器
			continue
		}
		if itemData.ItemEnumType == itemTypeConst.ITEM_RELIQUARY {
			// 排除圣遗物
			continue
		}
		if itemId == 100086 ||
			itemId == 100087 ||
			(itemId >= 100100 && itemId <= 101000) ||
			(itemId >= 101106 && itemId <= 101110) ||
			itemId == 101306 ||
			(itemId >= 101500 && itemId <= 104000) ||
			itemId == 105001 ||
			itemId == 105004 ||
			(itemId >= 106000 && itemId <= 107000) ||
			itemId == 107011 ||
			itemId == 108000 ||
			(itemId >= 109000 && itemId <= 110000) ||
			(itemId >= 115000 && itemId <= 130000) ||
			(itemId >= 200200 && itemId <= 200899) ||
			itemId == 220050 ||
			itemId == 220054 {
			// 排除无效道具
			continue
		}
		allItemDataConfig[itemId] = itemData
	}
	return allItemDataConfig
}

func (g *GameManager) AddUserItem(userId uint32, itemList []*UserItem, isHint bool) {
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	for _, userItem := range itemList {
		player.AddItem(userItem.ItemId, userItem.ChangeCount)
	}

	// PacketStoreItemChangeNotify
	storeItemChangeNotify := new(proto.StoreItemChangeNotify)
	storeItemChangeNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	for _, userItem := range itemList {
		pbItem := &proto.Item{
			ItemId: userItem.ItemId,
			Guid:   player.GetItemGuid(userItem.ItemId),
			Detail: &proto.Item_Material{
				Material: &proto.Material{
					Count: player.GetItemCount(userItem.ItemId),
				},
			},
		}
		storeItemChangeNotify.ItemList = append(storeItemChangeNotify.ItemList, pbItem)
	}
	g.SendMsg(proto.ApiStoreItemChangeNotify, userId, nil, storeItemChangeNotify)

	if isHint {
		actionReasonConst := constant.GetActionReasonConst()
		// PacketItemAddHintNotify
		itemAddHintNotify := new(proto.ItemAddHintNotify)
		itemAddHintNotify.Reason = uint32(actionReasonConst.SubfieldDrop)
		for _, userItem := range itemList {
			itemAddHintNotify.ItemList = append(itemAddHintNotify.ItemList, &proto.ItemHint{
				ItemId: userItem.ItemId,
				Count:  userItem.ChangeCount,
				IsNew:  false,
			})
		}
		g.SendMsg(proto.ApiItemAddHintNotify, userId, nil, itemAddHintNotify)
	}

	// PacketPlayerPropNotify
	playerPropNotify := new(proto.PlayerPropNotify)
	playerPropNotify.PropMap = make(map[uint32]*proto.PropValue)
	for _, userItem := range itemList {
		isVirtualItem, prop := player.GetVirtualItemProp(userItem.ItemId)
		if !isVirtualItem {
			continue
		}
		playerPropNotify.PropMap[uint32(prop)] = &proto.PropValue{
			Type: uint32(prop),
			Val:  int64(player.Properties[prop]),
			Value: &proto.PropValue_Ival{
				Ival: int64(player.Properties[prop]),
			},
		}
	}
	if len(playerPropNotify.PropMap) > 0 {
		g.SendMsg(proto.ApiPlayerPropNotify, userId, nil, playerPropNotify)
	}
}

func (g *GameManager) CostUserItem(userId uint32, itemList []*UserItem) {
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	for _, userItem := range itemList {
		player.CostItem(userItem.ItemId, userItem.ChangeCount)
	}

	// PacketStoreItemChangeNotify
	storeItemChangeNotify := new(proto.StoreItemChangeNotify)
	storeItemChangeNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	for _, userItem := range itemList {
		count := player.GetItemCount(userItem.ItemId)
		if count == 0 {
			continue
		}
		pbItem := &proto.Item{
			ItemId: userItem.ItemId,
			Guid:   player.GetItemGuid(userItem.ItemId),
			Detail: &proto.Item_Material{
				Material: &proto.Material{
					Count: count,
				},
			},
		}
		storeItemChangeNotify.ItemList = append(storeItemChangeNotify.ItemList, pbItem)
	}
	if len(storeItemChangeNotify.ItemList) > 0 {
		g.SendMsg(proto.ApiStoreItemChangeNotify, userId, nil, storeItemChangeNotify)
	}

	// PacketStoreItemDelNotify
	storeItemDelNotify := new(proto.StoreItemDelNotify)
	storeItemDelNotify.StoreType = proto.StoreType_STORE_TYPE_PACK
	for _, userItem := range itemList {
		count := player.GetItemCount(userItem.ItemId)
		if count > 0 {
			continue
		}
		storeItemDelNotify.GuidList = append(storeItemDelNotify.GuidList, player.GetItemGuid(userItem.ItemId))
	}
	if len(storeItemDelNotify.GuidList) > 0 {
		g.SendMsg(proto.ApiStoreItemDelNotify, userId, nil, storeItemDelNotify)
	}

	// PacketPlayerPropNotify
	playerPropNotify := new(proto.PlayerPropNotify)
	playerPropNotify.PropMap = make(map[uint32]*proto.PropValue)
	for _, userItem := range itemList {
		isVirtualItem, prop := player.GetVirtualItemProp(userItem.ItemId)
		if !isVirtualItem {
			continue
		}
		playerPropNotify.PropMap[uint32(prop)] = &proto.PropValue{
			Type: uint32(prop),
			Val:  int64(player.Properties[prop]),
			Value: &proto.PropValue_Ival{
				Ival: int64(player.Properties[prop]),
			},
		}
	}
	if len(playerPropNotify.PropMap) > 0 {
		g.SendMsg(proto.ApiPlayerPropNotify, userId, nil, playerPropNotify)
	}
}
