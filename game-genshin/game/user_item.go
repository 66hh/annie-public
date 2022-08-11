package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	"game-genshin/constant"
)

type UserItem struct {
	ItemId      uint32
	ChangeCount uint32
}

func (g *GameManager) AddUserItem(userId uint32, itemList []*UserItem, isHint bool) {
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		logger.LOG.Error("player not found, user id: %v", userId)
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
	g.SendMsg(api.ApiStoreItemChangeNotify, userId, nil, storeItemChangeNotify)

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
		g.SendMsg(api.ApiItemAddHintNotify, userId, nil, itemAddHintNotify)
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
		g.SendMsg(api.ApiPlayerPropNotify, userId, g.getHeadMsg(0), playerPropNotify)
	}
}

func (g *GameManager) CostUserItem(userId uint32, itemList []*UserItem) {
	player := g.userManager.GetTargetUser(userId)
	if player == nil {
		logger.LOG.Error("player not found, user id: %v", userId)
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
		g.SendMsg(api.ApiStoreItemChangeNotify, userId, nil, storeItemChangeNotify)
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
		g.SendMsg(api.ApiStoreItemDelNotify, userId, nil, storeItemDelNotify)
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
		g.SendMsg(api.ApiPlayerPropNotify, userId, g.getHeadMsg(0), playerPropNotify)
	}
}
