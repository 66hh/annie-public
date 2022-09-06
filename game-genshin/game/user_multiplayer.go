package game

import (
	"flswld.com/common/utils/object"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"game-genshin/constant"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
	"time"
)

func (g *GameManager) PlayerApplyEnterMpReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user apply enter world, user id: %v", userId)
	req := payloadMsg.(*proto.PlayerApplyEnterMpReq)
	targetUid := req.TargetUid

	g.UserApplyEnterWorld(player, targetUid)

	// PacketPlayerApplyEnterMpRsp
	playerApplyEnterMpRsp := new(proto.PlayerApplyEnterMpRsp)
	playerApplyEnterMpRsp.TargetUid = targetUid
	g.SendMsg(proto.ApiPlayerApplyEnterMpRsp, player.PlayerID, player.ClientSeq, playerApplyEnterMpRsp)
}

func (g *GameManager) PlayerApplyEnterMpResultReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user deal world enter apply, user id: %v", userId)
	req := payloadMsg.(*proto.PlayerApplyEnterMpResultReq)
	applyUid := req.ApplyUid
	isAgreed := req.IsAgreed

	g.UserDealEnterWorld(player, applyUid, isAgreed)

	// PacketPlayerApplyEnterMpResultRsp
	playerApplyEnterMpResultRsp := new(proto.PlayerApplyEnterMpResultRsp)
	playerApplyEnterMpResultRsp.ApplyUid = applyUid
	playerApplyEnterMpResultRsp.IsAgreed = isAgreed
	g.SendMsg(proto.ApiPlayerApplyEnterMpResultRsp, player.PlayerID, player.ClientSeq, playerApplyEnterMpResultRsp)
}

func (g *GameManager) PlayerGetForceQuitBanInfoReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user exit world, user id: %v", userId)

	ok := g.UserLeaveWorld(player)

	// PacketPlayerGetForceQuitBanInfoRsp
	playerGetForceQuitBanInfoRsp := new(proto.PlayerGetForceQuitBanInfoRsp)
	if ok {
		playerGetForceQuitBanInfoRsp.Retcode = int32(proto.Retcode_RETCODE_RET_SUCC)
	} else {
		playerGetForceQuitBanInfoRsp.Retcode = int32(proto.Retcode_RETCODE_RET_SVR_ERROR)
	}
	g.SendMsg(proto.ApiPlayerGetForceQuitBanInfoRsp, player.PlayerID, player.ClientSeq, playerGetForceQuitBanInfoRsp)
}

func (g *GameManager) UserApplyEnterWorld(hostPlayer *model.Player, otherUid uint32) {
	otherPlayer := g.userManager.GetOnlineUser(otherUid)
	if otherPlayer == nil {
		// PacketPlayerApplyEnterMpResultNotify
		playerApplyEnterMpResultNotify := new(proto.PlayerApplyEnterMpResultNotify)
		playerApplyEnterMpResultNotify.TargetUid = otherUid
		playerApplyEnterMpResultNotify.TargetNickname = ""
		playerApplyEnterMpResultNotify.IsAgreed = false
		playerApplyEnterMpResultNotify.Reason = proto.PlayerApplyEnterMpResultNotify_REASON_PLAYER_CANNOT_ENTER_MP
		g.SendMsg(proto.ApiPlayerApplyEnterMpResultNotify, hostPlayer.PlayerID, hostPlayer.ClientSeq, playerApplyEnterMpResultNotify)
		return
	}
	world := g.worldManager.GetWorldByID(hostPlayer.WorldId)
	if world.multiplayer {
		return
	}
	applyTime, exist := otherPlayer.CoopApplyMap[hostPlayer.PlayerID]
	if exist && time.Now().UnixNano() < applyTime+int64(10*time.Second) {
		return
	}
	otherPlayer.CoopApplyMap[hostPlayer.PlayerID] = time.Now().UnixNano()

	// PacketPlayerApplyEnterMpNotify
	playerApplyEnterMpNotify := new(proto.PlayerApplyEnterMpNotify)
	playerApplyEnterMpNotify.SrcPlayerInfo = g.PacketOnlinePlayerInfo(hostPlayer)
	g.SendMsg(proto.ApiPlayerApplyEnterMpNotify, otherPlayer.PlayerID, otherPlayer.ClientSeq, playerApplyEnterMpNotify)
}

func (g *GameManager) UserDealEnterWorld(hostPlayer *model.Player, otherUid uint32, agree bool) {
	otherPlayer := g.userManager.GetOnlineUser(otherUid)
	if otherPlayer == nil {
		return
	}
	applyTime, exist := hostPlayer.CoopApplyMap[otherUid]
	if !exist || time.Now().UnixNano() > applyTime+int64(10*time.Second) {
		return
	}
	delete(hostPlayer.CoopApplyMap, otherUid)
	otherPlayerWorld := g.worldManager.GetWorldByID(otherPlayer.WorldId)
	if otherPlayerWorld.multiplayer {
		// PacketPlayerApplyEnterMpResultNotify
		playerApplyEnterMpResultNotify := new(proto.PlayerApplyEnterMpResultNotify)
		playerApplyEnterMpResultNotify.TargetUid = hostPlayer.PlayerID
		playerApplyEnterMpResultNotify.TargetNickname = hostPlayer.NickName
		playerApplyEnterMpResultNotify.IsAgreed = false
		playerApplyEnterMpResultNotify.Reason = proto.PlayerApplyEnterMpResultNotify_REASON_PLAYER_CANNOT_ENTER_MP
		g.SendMsg(proto.ApiPlayerApplyEnterMpResultNotify, otherPlayer.PlayerID, otherPlayer.ClientSeq, playerApplyEnterMpResultNotify)
		return
	}

	// PacketPlayerApplyEnterMpResultNotify
	playerApplyEnterMpResultNotify := new(proto.PlayerApplyEnterMpResultNotify)
	playerApplyEnterMpResultNotify.TargetUid = hostPlayer.PlayerID
	playerApplyEnterMpResultNotify.TargetNickname = hostPlayer.NickName
	playerApplyEnterMpResultNotify.IsAgreed = agree
	playerApplyEnterMpResultNotify.Reason = proto.PlayerApplyEnterMpResultNotify_REASON_PLAYER_JUDGE
	g.SendMsg(proto.ApiPlayerApplyEnterMpResultNotify, otherPlayer.PlayerID, otherPlayer.ClientSeq, playerApplyEnterMpResultNotify)

	if !agree {
		return
	}

	enterReasonConst := constant.GetEnterReasonConst()
	hostWorld := g.worldManager.GetWorldByID(hostPlayer.WorldId)
	if hostWorld.multiplayer == false {
		g.UserWorldRemovePlayer(hostWorld, hostPlayer)
		hostWorld = g.worldManager.CreateWorld(hostPlayer, true)
		g.UserWorldAddPlayer(hostWorld, hostPlayer)
		hostPlayer.BornInScene = false

		// PacketPlayerEnterSceneNotify
		hostPlayerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyMp(
			hostPlayer,
			hostPlayer,
			proto.EnterType_ENTER_TYPE_SELF,
			uint32(enterReasonConst.HostFromSingleToMp),
			hostPlayer.SceneId,
			hostPlayer.Pos,
		)
		g.SendMsg(proto.ApiPlayerEnterSceneNotify, hostPlayer.PlayerID, hostPlayer.ClientSeq, hostPlayerEnterSceneNotify)
	}

	otherWorld := g.worldManager.GetWorldByID(otherPlayer.WorldId)
	g.UserWorldRemovePlayer(otherWorld, otherPlayer)
	_ = object.ObjectDeepCopy(hostPlayer.Pos, otherPlayer.Pos)
	_ = object.ObjectDeepCopy(hostPlayer.Rot, otherPlayer.Rot)
	otherPlayer.Pos.Y += 1
	otherPlayer.SceneId = hostPlayer.SceneId
	g.UserWorldAddPlayer(hostWorld, otherPlayer)
	otherPlayer.BornInScene = false

	// PacketPlayerEnterSceneNotify
	playerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyMp(
		otherPlayer,
		hostPlayer,
		proto.EnterType_ENTER_TYPE_OTHER,
		uint32(enterReasonConst.TeamJoin),
		hostPlayer.SceneId,
		hostPlayer.Pos,
	)
	g.SendMsg(proto.ApiPlayerEnterSceneNotify, otherPlayer.PlayerID, otherPlayer.ClientSeq, playerEnterSceneNotify)
}

func (g *GameManager) UserLeaveWorld(player *model.Player) bool {
	oldWorld := g.worldManager.GetWorldByID(player.WorldId)
	if !oldWorld.multiplayer {
		return false
	}
	// TODO SceneLoadState

	g.UserWorldRemovePlayer(oldWorld, player)
	newWorld := g.worldManager.CreateWorld(player, false)
	g.UserWorldAddPlayer(newWorld, player)
	player.BornInScene = false

	// PacketPlayerEnterSceneNotify
	enterReasonConst := constant.GetEnterReasonConst()
	hostPlayerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyMp(
		player,
		player,
		proto.EnterType_ENTER_TYPE_SELF,
		uint32(enterReasonConst.TeamBack),
		player.SceneId,
		player.Pos,
	)
	g.SendMsg(proto.ApiPlayerEnterSceneNotify, player.PlayerID, player.ClientSeq, hostPlayerEnterSceneNotify)
	return true
}

func (g *GameManager) UserWorldAddPlayer(world *World, player *model.Player) {
	_, exist := world.playerMap[player.PlayerID]
	if exist {
		return
	}
	world.AddPlayer(player, player.SceneId)
	player.WorldId = world.id
	scene := world.GetSceneById(player.SceneId)
	scene.UpdatePlayerTeamEntity(player)
	if len(world.playerMap) > 1 {
		g.UpdateWorldPlayerInfo(world, player)
	}
}

func (g *GameManager) UserWorldRemovePlayer(world *World, player *model.Player) {
	// PacketDelTeamEntityNotify
	delTeamEntityNotify := new(proto.DelTeamEntityNotify)
	delTeamEntityNotify.SceneId = player.SceneId
	delTeamEntityNotify.DelEntityIdList = []uint32{player.TeamConfig.TeamEntityId}
	g.SendMsg(proto.ApiDelTeamEntityNotify, player.PlayerID, player.ClientSeq, delTeamEntityNotify)

	g.RemoveSceneEntityAvatarBroadcastNotify(player)
	world.RemovePlayer(player)

	if len(world.playerMap) > 0 {
		g.UpdateWorldPlayerInfo(world, player)
	}
	if world.owner.PlayerID == player.PlayerID {
		// 房主离线清空所有玩家并销毁世界
		for _, worldPlayer := range world.playerMap {
			newWorld := g.worldManager.CreateWorld(worldPlayer, false)
			g.UserWorldAddPlayer(newWorld, worldPlayer)
			worldPlayer.BornInScene = false

			// PacketPlayerEnterSceneNotify
			enterReasonConst := constant.GetEnterReasonConst()
			hostPlayerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyMp(
				worldPlayer,
				worldPlayer,
				proto.EnterType_ENTER_TYPE_SELF,
				uint32(enterReasonConst.TeamKick),
				worldPlayer.SceneId,
				worldPlayer.Pos,
			)
			g.SendMsg(proto.ApiPlayerEnterSceneNotify, worldPlayer.PlayerID, worldPlayer.ClientSeq, hostPlayerEnterSceneNotify)
		}
		g.worldManager.DestroyWorld(world.id)
	}
}

func (g *GameManager) UpdateWorldPlayerInfo(hostWorld *World, excludePlayer *model.Player) {
	for _, worldPlayer := range hostWorld.playerMap {
		if worldPlayer.PlayerID == excludePlayer.PlayerID {
			continue
		}

		// TODO 更新队伍
		// PacketSceneTeamUpdateNotify
		sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(hostWorld)
		g.SendMsg(proto.ApiSceneTeamUpdateNotify, worldPlayer.PlayerID, worldPlayer.ClientSeq, sceneTeamUpdateNotify)

		// PacketWorldPlayerInfoNotify
		worldPlayerInfoNotify := new(proto.WorldPlayerInfoNotify)
		playerPropertyConst := constant.GetPlayerPropertyConst()
		for _, subWorldPlayer := range hostWorld.playerMap {
			onlinePlayerInfo := new(proto.OnlinePlayerInfo)
			onlinePlayerInfo.Uid = subWorldPlayer.PlayerID
			onlinePlayerInfo.Nickname = subWorldPlayer.NickName
			onlinePlayerInfo.PlayerLevel = subWorldPlayer.Properties[playerPropertyConst.PROP_PLAYER_LEVEL]
			onlinePlayerInfo.MpSettingType = proto.MpSettingType(subWorldPlayer.MpSetting)
			onlinePlayerInfo.NameCardId = subWorldPlayer.NameCard
			onlinePlayerInfo.Signature = subWorldPlayer.Signature
			// 头像
			onlinePlayerInfo.ProfilePicture = &proto.ProfilePicture{AvatarId: subWorldPlayer.HeadImage}
			onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(hostWorld.playerMap))
			worldPlayerInfoNotify.PlayerInfoList = append(worldPlayerInfoNotify.PlayerInfoList, onlinePlayerInfo)
			worldPlayerInfoNotify.PlayerUidList = append(worldPlayerInfoNotify.PlayerUidList, subWorldPlayer.PlayerID)
		}
		g.SendMsg(proto.ApiWorldPlayerInfoNotify, worldPlayer.PlayerID, worldPlayer.ClientSeq, worldPlayerInfoNotify)

		// PacketScenePlayerInfoNotify
		scenePlayerInfoNotify := new(proto.ScenePlayerInfoNotify)
		for _, subWorldPlayer := range hostWorld.playerMap {
			onlinePlayerInfo := new(proto.OnlinePlayerInfo)
			onlinePlayerInfo.Uid = subWorldPlayer.PlayerID
			onlinePlayerInfo.Nickname = subWorldPlayer.NickName
			onlinePlayerInfo.PlayerLevel = subWorldPlayer.Properties[playerPropertyConst.PROP_PLAYER_LEVEL]
			onlinePlayerInfo.MpSettingType = proto.MpSettingType(subWorldPlayer.MpSetting)
			onlinePlayerInfo.NameCardId = subWorldPlayer.NameCard
			onlinePlayerInfo.Signature = subWorldPlayer.Signature
			// 头像
			onlinePlayerInfo.ProfilePicture = &proto.ProfilePicture{AvatarId: subWorldPlayer.HeadImage}
			onlinePlayerInfo.CurPlayerNumInWorld = uint32(len(hostWorld.playerMap))
			scenePlayerInfoNotify.PlayerInfoList = append(scenePlayerInfoNotify.PlayerInfoList, &proto.ScenePlayerInfo{
				Uid:              subWorldPlayer.PlayerID,
				PeerId:           subWorldPlayer.PeerId,
				Name:             subWorldPlayer.NickName,
				SceneId:          subWorldPlayer.SceneId,
				OnlinePlayerInfo: onlinePlayerInfo,
			})
		}
		g.SendMsg(proto.ApiScenePlayerInfoNotify, worldPlayer.PlayerID, worldPlayer.ClientSeq, scenePlayerInfoNotify)

		// PacketWorldPlayerRTTNotify
		worldPlayerRTTNotify := new(proto.WorldPlayerRTTNotify)
		worldPlayerRTTNotify.PlayerRttList = make([]*proto.PlayerRTTInfo, 0)
		for _, subWorldPlayer := range hostWorld.playerMap {
			playerRTTInfo := &proto.PlayerRTTInfo{Uid: subWorldPlayer.PlayerID, Rtt: subWorldPlayer.ClientRTT}
			worldPlayerRTTNotify.PlayerRttList = append(worldPlayerRTTNotify.PlayerRttList, playerRTTInfo)
		}
		g.SendMsg(proto.ApiWorldPlayerRTTNotify, worldPlayer.PlayerID, 0, worldPlayerRTTNotify)

		// PacketSyncTeamEntityNotify
		syncTeamEntityNotify := new(proto.SyncTeamEntityNotify)
		syncTeamEntityNotify.SceneId = worldPlayer.SceneId
		syncTeamEntityNotify.TeamEntityInfoList = make([]*proto.TeamEntityInfo, 0)
		if hostWorld.multiplayer {
			for _, subWorldPlayer := range hostWorld.playerMap {
				if subWorldPlayer.PlayerID == worldPlayer.PlayerID {
					continue
				}
				teamEntityInfo := &proto.TeamEntityInfo{
					TeamEntityId:    subWorldPlayer.TeamConfig.TeamEntityId,
					AuthorityPeerId: subWorldPlayer.PeerId,
					TeamAbilityInfo: new(proto.AbilitySyncStateInfo),
				}
				syncTeamEntityNotify.TeamEntityInfoList = append(syncTeamEntityNotify.TeamEntityInfoList, teamEntityInfo)
			}
		}
		g.SendMsg(proto.ApiSyncTeamEntityNotify, worldPlayer.PlayerID, worldPlayer.ClientSeq, syncTeamEntityNotify)

		// PacketSyncScenePlayTeamEntityNotify
		syncScenePlayTeamEntityNotify := new(proto.SyncScenePlayTeamEntityNotify)
		syncScenePlayTeamEntityNotify.SceneId = worldPlayer.SceneId
		g.SendMsg(proto.ApiSyncScenePlayTeamEntityNotify, worldPlayer.PlayerID, worldPlayer.ClientSeq, syncScenePlayTeamEntityNotify)
	}
}
