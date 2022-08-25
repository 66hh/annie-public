package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/constant"
	"strconv"
)

func (g *GameManager) SceneTransToPointReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get scene trans to point, user id: %v", userId)
	req := payloadMsg.(*proto.SceneTransToPointReq)

	transPointId := strconv.Itoa(int(req.SceneId)) + "_" + strconv.Itoa(int(req.PointId))
	transPointConfig, exist := gdc.CONF.ScenePointEntries[transPointId]
	if !exist {
		// PacketSceneTransToPointRsp
		sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
		// TODO Retcode.proto
		sceneTransToPointRsp.Retcode = 1 // RET_SVR_ERROR_VALUE
		g.SendMsg(api.ApiSceneTransToPointRsp, userId, nil, sceneTransToPointRsp)
		return
	}

	// 传送玩家
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}

	oldSceneId := player.SceneId
	newSceneId := req.SceneId

	world := g.worldManager.GetWorldByID(player.WorldId)
	oldScene := world.GetSceneById(oldSceneId)

	//// PacketSceneEntityDisappearNotify
	//sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
	//activeAvatarId := player.TeamConfig.GetActiveAvatarId()
	//playerTeamEntity := oldScene.GetPlayerTeamEntity(player.PlayerID)
	//sceneEntityDisappearNotify.EntityList = []uint32{playerTeamEntity.avatarEntityMap[activeAvatarId]}
	//sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REMOVE
	//g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)

	g.RemoveSceneEntityAvatarBroadcastNotify(player)

	oldScene.RemovePlayer(player)

	newScene := world.GetSceneById(newSceneId)
	newScene.AddPlayer(player)
	newScene.UpdatePlayerTeamEntity(player)
	player.Pos.X = transPointConfig.PointData.TranPos.X
	player.Pos.Y = transPointConfig.PointData.TranPos.Y
	player.Pos.Z = transPointConfig.PointData.TranPos.Z
	player.SceneId = newSceneId
	logger.LOG.Info("player goto scene: %v, pos x: %v, y: %v, z: %v", newSceneId, player.Pos.X, player.Pos.Y, player.Pos.Z)
	//g.userManager.UpdateUser(player)

	player.BornInScene = false

	// PacketPlayerEnterSceneNotify
	var enterType proto.EnterType
	if newSceneId == oldSceneId {
		enterType = proto.EnterType_ENTER_TYPE_GOTO
	} else {
		enterType = proto.EnterType_ENTER_TYPE_JUMP
	}
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyTp(player, enterType, uint32(enterReasonConst.TransPoint), newSceneId, player.Pos)
	g.SendMsg(api.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)
	//g.userManager.UpdateUser(player)

	// PacketSceneTransToPointRsp
	sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
	sceneTransToPointRsp.Retcode = 0
	sceneTransToPointRsp.PointId = req.PointId
	sceneTransToPointRsp.SceneId = req.SceneId
	g.SendMsg(api.ApiSceneTransToPointRsp, userId, nil, sceneTransToPointRsp)
}

func (g *GameManager) MarkMapReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user mark map, user id: %v", userId)
	req := payloadMsg.(*proto.MarkMapReq)
	operation := req.Op
	if operation == proto.MarkMapReq_OPERATION_ADD {
		logger.LOG.Debug("user mark type: %v", req.Mark.PointType)
		if req.Mark.PointType == proto.MapMarkPointType_MAP_MARK_POINT_TYPE_NPC {
			posYInt, err := strconv.ParseInt(req.Mark.Name, 10, 64)
			if err != nil {
				logger.LOG.Error("parse pos y error: %v", err)
				posYInt = 0
			}

			// 传送玩家
			player := g.userManager.GetOnlineUser(userId)
			if player == nil {
				logger.LOG.Error("player is nil, userId: %v", userId)
				return
			}

			oldSceneId := player.SceneId
			newSceneId := req.Mark.SceneId

			world := g.worldManager.GetWorldByID(player.WorldId)
			oldScene := world.GetSceneById(oldSceneId)

			//// PacketSceneEntityDisappearNotify
			//sceneEntityDisappearNotify := new(proto.SceneEntityDisappearNotify)
			//activeAvatarId := player.TeamConfig.GetActiveAvatarId()
			//playerTeamEntity := oldScene.GetPlayerTeamEntity(player.PlayerID)
			//sceneEntityDisappearNotify.EntityList = []uint32{playerTeamEntity.avatarEntityMap[activeAvatarId]}
			//sceneEntityDisappearNotify.DisappearType = proto.VisionType_VISION_TYPE_REMOVE
			//g.SendMsg(api.ApiSceneEntityDisappearNotify, userId, nil, sceneEntityDisappearNotify)

			//g.userManager.UpdateUser(player)

			g.RemoveSceneEntityAvatarBroadcastNotify(player)

			oldScene.RemovePlayer(player)

			newScene := world.GetSceneById(newSceneId)
			newScene.AddPlayer(player)
			newScene.UpdatePlayerTeamEntity(player)
			x := float64(req.Mark.Pos.X)
			y := float64(posYInt)
			z := float64(req.Mark.Pos.Z)
			player.Pos.X = x
			player.Pos.Y = y
			player.Pos.Z = z
			player.SceneId = newSceneId
			logger.LOG.Info("player goto scene: %v, pos x: %v, y: %v, z: %v", newSceneId, x, y, z)
			//g.userManager.UpdateUser(player)

			player.BornInScene = false

			// PacketPlayerEnterSceneNotify
			var enterType proto.EnterType
			if newSceneId == oldSceneId {
				enterType = proto.EnterType_ENTER_TYPE_GOTO
			} else {
				enterType = proto.EnterType_ENTER_TYPE_JUMP
			}
			enterReasonConst := constant.GetEnterReasonConst()
			playerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyTp(player, enterType, uint32(enterReasonConst.TransPoint), newSceneId, player.Pos)
			g.SendMsg(api.ApiPlayerEnterSceneNotify, userId, nil, playerEnterSceneNotify)
			//g.userManager.UpdateUser(player)
		}
	}
}
