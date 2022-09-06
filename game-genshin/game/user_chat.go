package game

import (
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
	"time"
)

func (g *GameManager) PullRecentChatReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user pull recent chat, user id: %v", userId)
	req := payloadMsg.(*proto.PullRecentChatReq)
	// 经研究发现 原神现网环境 客户端仅拉取最新的5条未读聊天消息 所以人太多的话小姐姐不回你消息是有原因的
	// 因此 阿米你这样做真的合适吗 不过现在代码到了我手上我想怎么写就怎么写 我才不会重蹈覆辙
	_ = req.PullNum

	retMsgList := make([]*proto.ChatInfo, 0)
	for _, msgList := range player.ChatMsgMap {
		for _, chatMsg := range msgList {
			// 反手就是一个遍历
			if chatMsg.IsRead {
				continue
			}
			retMsgList = append(retMsgList, g.ConvChatMsgToChatInfo(chatMsg))
		}
	}

	world := g.worldManager.GetWorldByID(player.WorldId)
	if world.multiplayer {
		chatList := world.GetChatList()
		count := len(chatList)
		if count > 10 {
			count = 10
		}
		for i := len(chatList) - count; i < len(chatList); i++ {
			// PacketPlayerChatNotify
			playerChatNotify := new(proto.PlayerChatNotify)
			playerChatNotify.ChannelId = 0
			playerChatNotify.ChatInfo = chatList[i]
			g.SendMsg(proto.ApiPlayerChatNotify, player.PlayerID, 0, playerChatNotify)
		}
	}

	// PacketPullRecentChatRsp
	pullRecentChatRsp := new(proto.PullRecentChatRsp)
	pullRecentChatRsp.ChatInfo = retMsgList
	g.SendMsg(proto.ApiPullRecentChatRsp, player.PlayerID, 0, pullRecentChatRsp)
}

func (g *GameManager) PullPrivateChatReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user pull private chat, user id: %v", userId)
	req := payloadMsg.(*proto.PullPrivateChatReq)
	targetUid := req.TargetUid
	pullNum := req.PullNum
	fromSequence := req.FromSequence

	msgList, exist := player.ChatMsgMap[targetUid]
	if !exist {
		return
	}
	if pullNum+fromSequence > uint32(len(msgList)) {
		pullNum = uint32(len(msgList)) - fromSequence
	}
	recentMsgList := msgList[fromSequence : fromSequence+pullNum]
	retMsgList := make([]*proto.ChatInfo, 0)
	for _, chatMsg := range recentMsgList {
		retMsgList = append(retMsgList, g.ConvChatMsgToChatInfo(chatMsg))
	}

	// PacketPullPrivateChatRsp
	pullPrivateChatRsp := new(proto.PullPrivateChatRsp)
	pullPrivateChatRsp.ChatInfo = retMsgList
	g.SendMsg(proto.ApiPullPrivateChatRsp, player.PlayerID, 0, pullPrivateChatRsp)
}

func (g *GameManager) PrivateChatReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user send private chat, user id: %v", userId)
	req := payloadMsg.(*proto.PrivateChatReq)
	targetUid := req.TargetUid
	content := req.Content

	// TODO 同步阻塞待优化
	targetPlayer := g.userManager.LoadTempOfflineUserSync(targetUid)
	if targetPlayer == nil {
		return
	}
	chatInfo := &proto.ChatInfo{
		Time:     uint32(time.Now().Unix()),
		Sequence: 101,
		ToUid:    targetPlayer.PlayerID,
		Uid:      player.PlayerID,
		IsRead:   false,
		Content:  nil,
	}
	switch content.(type) {
	case *proto.PrivateChatReq_Text:
		text := content.(*proto.PrivateChatReq_Text).Text
		if len(text) == 0 {
			return
		}
		chatInfo.Content = &proto.ChatInfo_Text{
			Text: text,
		}
	case *proto.PrivateChatReq_Icon:
		icon := content.(*proto.PrivateChatReq_Icon).Icon
		chatInfo.Content = &proto.ChatInfo_Icon{
			Icon: icon,
		}
	default:
		return
	}

	// 消息加入自己的队列
	msgList, exist := player.ChatMsgMap[targetPlayer.PlayerID]
	if !exist {
		msgList = make([]*model.ChatMsg, 0)
	}
	msgList = append(msgList, g.ConvChatInfoToChatMsg(chatInfo))
	player.ChatMsgMap[targetPlayer.PlayerID] = msgList
	// 消息加入目标玩家的队列
	msgList, exist = targetPlayer.ChatMsgMap[player.PlayerID]
	if !exist {
		msgList = make([]*model.ChatMsg, 0)
	}
	msgList = append(msgList, g.ConvChatInfoToChatMsg(chatInfo))
	targetPlayer.ChatMsgMap[player.PlayerID] = msgList

	if targetPlayer.Online {
		// PacketPrivateChatNotify
		privateChatNotify := new(proto.PrivateChatNotify)
		privateChatNotify.ChatInfo = chatInfo
		g.SendMsg(proto.ApiPrivateChatNotify, targetPlayer.PlayerID, 0, privateChatNotify)
	}

	// PacketPrivateChatNotify
	privateChatNotify := new(proto.PrivateChatNotify)
	privateChatNotify.ChatInfo = chatInfo
	g.SendMsg(proto.ApiPrivateChatNotify, player.PlayerID, 0, privateChatNotify)

	// PacketPrivateChatRsp
	privateChatRsp := new(proto.PrivateChatRsp)
	g.SendMsg(proto.ApiPrivateChatRsp, player.PlayerID, 0, privateChatRsp)
}

func (g *GameManager) ReadPrivateChatReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user read private chat, user id: %v", userId)
	req := payloadMsg.(*proto.ReadPrivateChatReq)
	targetUid := req.TargetUid

	msgList, exist := player.ChatMsgMap[targetUid]
	if !exist {
		return
	}
	for index, chatMsg := range msgList {
		chatMsg.IsRead = true
		msgList[index] = chatMsg
	}
	player.ChatMsgMap[targetUid] = msgList

	// PacketReadPrivateChatRsp
	readPrivateChatRsp := new(proto.ReadPrivateChatRsp)
	g.SendMsg(proto.ApiReadPrivateChatRsp, player.PlayerID, 0, readPrivateChatRsp)
}

func (g *GameManager) PlayerChatReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user multiplayer chat, user id: %v", userId)
	req := payloadMsg.(*proto.PlayerChatReq)
	channelId := req.ChannelId
	chatInfo := req.ChatInfo

	sendChatInfo := &proto.ChatInfo{
		Time:    uint32(time.Now().Unix()),
		Uid:     player.PlayerID,
		Content: nil,
	}
	switch chatInfo.Content.(type) {
	case *proto.ChatInfo_Text:
		text := chatInfo.Content.(*proto.ChatInfo_Text).Text
		if len(text) == 0 {
			return
		}
		sendChatInfo.Content = &proto.ChatInfo_Text{
			Text: text,
		}
	case *proto.ChatInfo_Icon:
		icon := chatInfo.Content.(*proto.ChatInfo_Icon).Icon
		sendChatInfo.Content = &proto.ChatInfo_Icon{
			Icon: icon,
		}
	default:
		return
	}

	world := g.worldManager.GetWorldByID(player.WorldId)
	world.AddChat(sendChatInfo)

	// PacketPlayerChatNotify
	playerChatNotify := new(proto.PlayerChatNotify)
	playerChatNotify.ChannelId = channelId
	playerChatNotify.ChatInfo = sendChatInfo
	for _, worldPlayer := range world.playerMap {
		g.SendMsg(proto.ApiPlayerChatNotify, worldPlayer.PlayerID, 0, playerChatNotify)
	}

	// PacketPlayerChatRsp
	playerChatRsp := new(proto.PlayerChatRsp)
	g.SendMsg(proto.ApiPlayerChatRsp, player.PlayerID, 0, playerChatRsp)
}

func (g *GameManager) ConvChatInfoToChatMsg(chatInfo *proto.ChatInfo) (chatMsg *model.ChatMsg) {
	chatMsg = &model.ChatMsg{
		Time:    chatInfo.Time,
		ToUid:   chatInfo.ToUid,
		Uid:     chatInfo.Uid,
		IsRead:  chatInfo.IsRead,
		MsgType: 0,
		Text:    "",
		Icon:    0,
	}
	switch chatInfo.Content.(type) {
	case *proto.ChatInfo_Text:
		chatMsg.MsgType = model.ChatMsgTypeText
		chatMsg.Text = chatInfo.Content.(*proto.ChatInfo_Text).Text
	case *proto.ChatInfo_Icon:
		chatMsg.MsgType = model.ChatMsgTypeIcon
		chatMsg.Icon = chatInfo.Content.(*proto.ChatInfo_Icon).Icon
	default:
	}
	return chatMsg
}

func (g *GameManager) ConvChatMsgToChatInfo(chatMsg *model.ChatMsg) (chatInfo *proto.ChatInfo) {
	chatInfo = &proto.ChatInfo{
		Time:     chatMsg.Time,
		Sequence: 0,
		ToUid:    chatMsg.ToUid,
		Uid:      chatMsg.Uid,
		IsRead:   chatMsg.IsRead,
		Content:  nil,
	}
	switch chatMsg.MsgType {
	case model.ChatMsgTypeText:
		chatInfo.Content = &proto.ChatInfo_Text{
			Text: chatMsg.Text,
		}
	case model.ChatMsgTypeIcon:
		chatInfo.Content = &proto.ChatInfo_Icon{
			Icon: chatMsg.Icon,
		}
	default:
	}
	return chatInfo
}
