package game

import (
	"flswld.com/common/utils/object"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	"game-genshin/constant"
	"game-genshin/model"
	"regexp"
	"time"
	"unicode/utf8"
)

func (g *GameManager) GetPlayerSocialDetailReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get player social detail, user id: %v", userId)
	req := payloadMsg.(*proto.GetPlayerSocialDetailReq)
	targetUid := req.Uid
	getPlayerSocialDetailRsp := new(proto.GetPlayerSocialDetailRsp)
	// TODO 同步阻塞待优化
	targetPlayer := g.userManager.LoadTempOfflineUserSync(targetUid)
	if targetPlayer != nil {
		socialDetail := new(proto.SocialDetail)
		socialDetail.Uid = targetPlayer.PlayerID
		socialDetail.ProfilePicture = &proto.ProfilePicture{AvatarId: targetPlayer.HeadImage}
		socialDetail.Nickname = targetPlayer.NickName
		socialDetail.Signature = targetPlayer.Signature
		playerPropertyConst := constant.GetPlayerPropertyConst()
		socialDetail.Level = targetPlayer.Properties[playerPropertyConst.PROP_PLAYER_LEVEL]
		socialDetail.Birthday = &proto.Birthday{Month: 2, Day: 13}
		socialDetail.WorldLevel = targetPlayer.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL]
		socialDetail.NameCardId = targetPlayer.NameCard
		socialDetail.IsShowAvatar = false
		socialDetail.FinishAchievementNum = 0
		socialDetail.IsFriend = false
		getPlayerSocialDetailRsp.DetailData = socialDetail
	} else {
		getPlayerSocialDetailRsp.Retcode = int32(proto.Retcode_RET_PLAYER_NOT_EXIST)
	}
	g.SendMsg(api.ApiGetPlayerSocialDetailRsp, userId, nil, getPlayerSocialDetailRsp)
}

func (g *GameManager) SetPlayerBirthdayReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user set birthday, user id: %v", userId)
	req := payloadMsg.(*proto.SetPlayerBirthdayReq)
	_ = req
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
}

func (g *GameManager) SetNameCardReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user change name card, user id: %v", userId)
	req := payloadMsg.(*proto.SetNameCardReq)
	nameCardId := req.NameCardId
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	exist := false
	for _, nameCard := range player.NameCardList {
		if nameCard == nameCardId {
			exist = true
		}
	}
	if !exist {
		logger.LOG.Error("name card not exist, userId: %v", userId)
		return
	}
	player.NameCard = nameCardId

	// PacketSetNameCardRsp
	setNameCardRsp := new(proto.SetNameCardRsp)
	setNameCardRsp.NameCardId = nameCardId
	g.SendMsg(api.ApiSetNameCardRsp, userId, nil, setNameCardRsp)
}

func (g *GameManager) SetPlayerSignatureReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user change signature, user id: %v", userId)
	req := payloadMsg.(*proto.SetPlayerSignatureReq)
	signature := req.Signature
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	// PacketSetPlayerSignatureRsp
	setPlayerSignatureRsp := new(proto.SetPlayerSignatureRsp)
	if !object.IsUtf8String(signature) {
		setPlayerSignatureRsp.Retcode = int32(proto.Retcode_RET_SIGNATURE_ILLEGAL)
	} else if utf8.RuneCountInString(signature) > 50 {
		setPlayerSignatureRsp.Retcode = int32(proto.Retcode_RET_SIGNATURE_ILLEGAL)
	} else {
		player.Signature = signature
		setPlayerSignatureRsp.Signature = player.Signature
	}
	g.SendMsg(api.ApiSetPlayerSignatureRsp, userId, nil, setPlayerSignatureRsp)
}

func (g *GameManager) SetPlayerNameReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user change nickname, user id: %v", userId)
	req := payloadMsg.(*proto.SetPlayerNameReq)
	nickName := req.NickName
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	// PacketSetPlayerNameRsp
	setPlayerNameRsp := new(proto.SetPlayerNameRsp)
	if len(nickName) == 0 {
		setPlayerNameRsp.Retcode = int32(proto.Retcode_RET_NICKNAME_IS_EMPTY)
	} else if !object.IsUtf8String(nickName) {
		setPlayerNameRsp.Retcode = int32(proto.Retcode_RET_NICKNAME_UTF8_ERROR)
	} else if utf8.RuneCountInString(nickName) > 14 {
		setPlayerNameRsp.Retcode = int32(proto.Retcode_RET_NICKNAME_TOO_LONG)
	} else if len(regexp.MustCompile(`\d`).FindAllString(nickName, -1)) > 6 {
		setPlayerNameRsp.Retcode = int32(proto.Retcode_RET_NICKNAME_TOO_MANY_DIGITS)
	} else {
		player.NickName = nickName
		setPlayerNameRsp.NickName = player.NickName
	}
	g.SendMsg(api.ApiSetPlayerNameRsp, userId, nil, setPlayerNameRsp)
}

func (g *GameManager) SetPlayerHeadImageReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user change head image, user id: %v", userId)
	req := payloadMsg.(*proto.SetPlayerHeadImageReq)
	avatarId := req.AvatarId
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	_, exist := player.AvatarMap[avatarId]
	if !exist {
		logger.LOG.Error("the head img of the avatar not exist, userId: %v", userId)
		return
	}
	player.HeadImage = avatarId

	// PacketSetPlayerHeadImageRsp
	setPlayerHeadImageRsp := new(proto.SetPlayerHeadImageRsp)
	setPlayerHeadImageRsp.ProfilePicture = &proto.ProfilePicture{AvatarId: player.HeadImage}
	g.SendMsg(api.ApiSetPlayerHeadImageRsp, userId, nil, setPlayerHeadImageRsp)
}

func (g *GameManager) GetAllUnlockNameCardReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get all unlock name card, user id: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	// PacketGetAllUnlockNameCardRsp
	getAllUnlockNameCardRsp := new(proto.GetAllUnlockNameCardRsp)
	getAllUnlockNameCardRsp.NameCardList = player.NameCardList
	g.SendMsg(api.ApiGetAllUnlockNameCardRsp, userId, nil, getAllUnlockNameCardRsp)
}

func (g *GameManager) GetPlayerFriendListReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get friend list, user id: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	// PacketGetPlayerFriendListRsp
	getPlayerFriendListRsp := new(proto.GetPlayerFriendListRsp)
	getPlayerFriendListRsp.FriendList = make([]*proto.FriendBrief, 0)
	for _, uid := range player.FriendList {
		// TODO 同步阻塞待优化
		var onlineState proto.FriendOnlineState
		online := g.userManager.GetUserOnlineState(uid)
		if online {
			onlineState = proto.FriendOnlineState_FRIEND_ONLINE_STATE_ONLINE
		} else {
			onlineState = proto.FriendOnlineState_FRIEND_ONLINE_STATE_FREIEND_DISCONNECT
		}
		friendPlayer := g.userManager.LoadTempOfflineUserSync(uid)
		if friendPlayer == nil {
			logger.LOG.Error("target player is nil, userId: %v", userId)
			continue
		}
		playerPropertyConst := constant.GetPlayerPropertyConst()
		friendBrief := &proto.FriendBrief{
			Uid:               friendPlayer.PlayerID,
			Nickname:          friendPlayer.NickName,
			Level:             friendPlayer.Properties[playerPropertyConst.PROP_PLAYER_LEVEL],
			ProfilePicture:    &proto.ProfilePicture{AvatarId: friendPlayer.HeadImage},
			WorldLevel:        friendPlayer.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL],
			Signature:         friendPlayer.Signature,
			OnlineState:       onlineState,
			IsMpModeAvailable: true,
			LastActiveTime:    player.OfflineTime,
			NameCardId:        friendPlayer.NameCard,
			Param:             (uint32(time.Now().Unix()) - player.OfflineTime) / 3600 / 24,
			IsGameSource:      true,
			PlatformType:      proto.PlatformType_PLATFORM_TYPE_PC,
		}
		getPlayerFriendListRsp.FriendList = append(getPlayerFriendListRsp.FriendList, friendBrief)
	}
	g.SendMsg(api.ApiGetPlayerFriendListRsp, userId, nil, getPlayerFriendListRsp)
}

func (g *GameManager) GetPlayerAskFriendListReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get friend apply list, user id: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	// PacketGetPlayerAskFriendListRsp
	getPlayerAskFriendListRsp := new(proto.GetPlayerAskFriendListRsp)
	getPlayerAskFriendListRsp.AskFriendList = make([]*proto.FriendBrief, 0)
	for _, uid := range player.FriendApplyList {
		// TODO 同步阻塞待优化
		var onlineState proto.FriendOnlineState
		online := g.userManager.GetUserOnlineState(uid)
		if online {
			onlineState = proto.FriendOnlineState_FRIEND_ONLINE_STATE_ONLINE
		} else {
			onlineState = proto.FriendOnlineState_FRIEND_ONLINE_STATE_FREIEND_DISCONNECT
		}
		friendPlayer := g.userManager.LoadTempOfflineUserSync(uid)
		if friendPlayer == nil {
			logger.LOG.Error("target player is nil, userId: %v", userId)
			continue
		}
		playerPropertyConst := constant.GetPlayerPropertyConst()
		friendBrief := &proto.FriendBrief{
			Uid:               friendPlayer.PlayerID,
			Nickname:          friendPlayer.NickName,
			Level:             friendPlayer.Properties[playerPropertyConst.PROP_PLAYER_LEVEL],
			ProfilePicture:    &proto.ProfilePicture{AvatarId: friendPlayer.HeadImage},
			WorldLevel:        friendPlayer.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL],
			Signature:         friendPlayer.Signature,
			OnlineState:       onlineState,
			IsMpModeAvailable: true,
			LastActiveTime:    player.OfflineTime,
			NameCardId:        friendPlayer.NameCard,
			Param:             (uint32(time.Now().Unix()) - player.OfflineTime) / 3600 / 24,
			IsGameSource:      true,
			PlatformType:      proto.PlatformType_PLATFORM_TYPE_PC,
		}
		getPlayerAskFriendListRsp.AskFriendList = append(getPlayerAskFriendListRsp.AskFriendList, friendBrief)
	}
	g.SendMsg(api.ApiGetPlayerAskFriendListRsp, userId, nil, getPlayerAskFriendListRsp)
}

func (g *GameManager) AskAddFriendReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user apply add friend, user id: %v", userId)
	req := payloadMsg.(*proto.AskAddFriendReq)
	targetUid := req.TargetUid
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	// TODO 同步阻塞待优化
	targetPlayerOnline := g.userManager.GetUserOnlineState(targetUid)
	targetPlayer := g.userManager.LoadTempOfflineUserSync(targetUid)
	if targetPlayer == nil {
		logger.LOG.Error("apply add friend target player is nil, userId: %v", userId)
		return
	}
	exist := false
	for _, uid := range targetPlayer.FriendApplyList {
		if uid == player.PlayerID {
			exist = true
		}
	}
	for _, uid := range targetPlayer.FriendList {
		if uid == player.PlayerID {
			exist = true
		}
	}
	if exist {
		logger.LOG.Error("friend or apply already exist, user id: %v", userId)
		return
	}
	targetPlayer.FriendApplyList = append(targetPlayer.FriendApplyList, player.PlayerID)

	if targetPlayerOnline {
		playerPropertyConst := constant.GetPlayerPropertyConst()
		// PacketAskAddFriendNotify
		askAddFriendNotify := new(proto.AskAddFriendNotify)
		askAddFriendNotify.TargetUid = player.PlayerID
		askAddFriendNotify.TargetFriendBrief = &proto.FriendBrief{
			Uid:               player.PlayerID,
			Nickname:          player.NickName,
			Level:             player.Properties[playerPropertyConst.PROP_PLAYER_LEVEL],
			ProfilePicture:    &proto.ProfilePicture{AvatarId: player.HeadImage},
			WorldLevel:        player.Properties[playerPropertyConst.PROP_PLAYER_WORLD_LEVEL],
			Signature:         player.Signature,
			OnlineState:       proto.FriendOnlineState_FRIEND_ONLINE_STATE_ONLINE,
			IsMpModeAvailable: true,
			LastActiveTime:    player.OfflineTime,
			NameCardId:        player.NameCard,
			Param:             (uint32(time.Now().Unix()) - player.OfflineTime) / 3600 / 24,
			IsGameSource:      true,
			PlatformType:      proto.PlatformType_PLATFORM_TYPE_PC,
		}
		g.SendMsg(api.ApiAskAddFriendNotify, targetPlayer.PlayerID, nil, askAddFriendNotify)
	}

	// PacketAskAddFriendRsp
	askAddFriendRsp := new(proto.AskAddFriendRsp)
	askAddFriendRsp.TargetUid = targetUid
	g.SendMsg(api.ApiAskAddFriendRsp, userId, nil, askAddFriendRsp)
}

func (g *GameManager) DealAddFriendReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user deal friend apply, user id: %v", userId)
	req := payloadMsg.(*proto.DealAddFriendReq)
	targetUid := req.TargetUid
	result := req.DealAddFriendResult
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	if result == proto.DealAddFriendResultType_DEAL_ADD_FRIEND_RESULT_TYPE_ACCEPT {
		player.FriendList = append(player.FriendList, targetUid)
		// TODO 同步阻塞待优化
		targetPlayer := g.userManager.LoadTempOfflineUserSync(targetUid)
		if targetPlayer == nil {
			logger.LOG.Error("agree friend apply target player is nil, userId: %v", userId)
			return
		}
		targetPlayer.FriendList = append(targetPlayer.FriendList, player.PlayerID)
	}

	newFriendApplyList := make([]uint32, 0)
	for _, uid := range player.FriendApplyList {
		if uid == targetUid {
			continue
		}
		newFriendApplyList = append(newFriendApplyList, uid)
	}
	player.FriendApplyList = newFriendApplyList

	// PacketDealAddFriendRsp
	dealAddFriendRsp := new(proto.DealAddFriendRsp)
	dealAddFriendRsp.TargetUid = targetUid
	dealAddFriendRsp.DealAddFriendResult = result
	g.SendMsg(api.ApiDealAddFriendRsp, userId, nil, dealAddFriendRsp)
}

func (g *GameManager) GetOnlinePlayerListReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get online player list, user id: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	count := 0
	onlinePlayerList := make([]*model.Player, 0)
	g.userManager.playerMapLock.RLock()
	for _, onlinePlayer := range g.userManager.playerMap {
		if onlinePlayer.PlayerID == player.PlayerID {
			continue
		}
		if onlinePlayer.Online == false {
			continue
		}
		onlinePlayerList = append(onlinePlayerList, onlinePlayer)
		count++
		if count >= 50 {
			break
		}
	}
	g.userManager.playerMapLock.RUnlock()

	// PacketGetOnlinePlayerListRsp
	getOnlinePlayerListRsp := new(proto.GetOnlinePlayerListRsp)
	getOnlinePlayerListRsp.PlayerInfoList = make([]*proto.OnlinePlayerInfo, 0)
	for _, onlinePlayer := range onlinePlayerList {
		onlinePlayerInfo := g.PacketOnlinePlayerInfo(onlinePlayer)
		getOnlinePlayerListRsp.PlayerInfoList = append(getOnlinePlayerListRsp.PlayerInfoList, onlinePlayerInfo)
	}
	g.SendMsg(api.ApiGetOnlinePlayerListRsp, userId, nil, getOnlinePlayerListRsp)
}

func (g *GameManager) PacketOnlinePlayerInfo(player *model.Player) *proto.OnlinePlayerInfo {
	playerPropertyConst := constant.GetPlayerPropertyConst()
	onlinePlayerInfo := &proto.OnlinePlayerInfo{
		Uid:                 player.PlayerID,
		Nickname:            player.NickName,
		PlayerLevel:         player.Properties[playerPropertyConst.PROP_PLAYER_LEVEL],
		MpSettingType:       player.MpSetting,
		NameCardId:          player.NameCard,
		Signature:           player.Signature,
		ProfilePicture:      &proto.ProfilePicture{AvatarId: player.HeadImage},
		CurPlayerNumInWorld: 1,
	}
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world != nil && world.playerMap != nil {
		onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(world.playerMap))
	}
	return onlinePlayerInfo
}
