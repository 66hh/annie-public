package game

import (
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/constant"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
	"strconv"
)

func (g *GameManager) SceneTransToPointReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user get scene trans to point, user id: %v", userId)
	req := payloadMsg.(*proto.SceneTransToPointReq)

	transPointId := strconv.Itoa(int(req.SceneId)) + "_" + strconv.Itoa(int(req.PointId))
	transPointConfig, exist := gdc.CONF.ScenePointEntries[transPointId]
	if !exist {
		// PacketSceneTransToPointRsp
		sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
		sceneTransToPointRsp.Retcode = int32(proto.Retcode_RETCODE_RET_SVR_ERROR)
		g.SendMsg(proto.ApiSceneTransToPointRsp, userId, player.ClientSeq, sceneTransToPointRsp)
		return
	}

	// 传送玩家
	newSceneId := req.SceneId
	oldSceneId := player.SceneId
	oldPos := &model.Vector{
		X: player.Pos.X,
		Y: player.Pos.Y,
		Z: player.Pos.Z,
	}
	jumpScene := false
	if newSceneId != oldSceneId {
		jumpScene = true
	}
	g.RemoveSceneEntityAvatarBroadcastNotify(player)
	world := g.worldManager.GetWorldByID(player.WorldId)
	oldScene := world.GetSceneById(oldSceneId)
	if jumpScene {
		// PacketDelTeamEntityNotify
		delTeamEntityNotify := g.PacketDelTeamEntityNotify(oldScene, player)
		g.SendMsg(proto.ApiDelTeamEntityNotify, player.PlayerID, player.ClientSeq, delTeamEntityNotify)

		oldScene.RemovePlayer(player)
		newScene := world.GetSceneById(newSceneId)
		newScene.AddPlayer(player)
	} else {
		oldScene.UpdatePlayerTeamEntity(player)
	}
	player.Pos.X = transPointConfig.PointData.TranPos.X
	player.Pos.Y = transPointConfig.PointData.TranPos.Y
	player.Pos.Z = transPointConfig.PointData.TranPos.Z
	player.SceneId = newSceneId
	player.SceneLoadState = model.SceneNone

	// PacketPlayerEnterSceneNotify
	var enterType proto.EnterType
	if jumpScene {
		logger.LOG.Debug("player jump scene, scene: %v, pos: %v", player.SceneId, player.Pos)
		enterType = proto.EnterType_ENTER_TYPE_JUMP
	} else {
		logger.LOG.Debug("player goto scene, scene: %v, pos: %v", player.SceneId, player.Pos)
		enterType = proto.EnterType_ENTER_TYPE_GOTO
	}
	enterReasonConst := constant.GetEnterReasonConst()
	playerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyTp(player, enterType, uint32(enterReasonConst.TransPoint), oldSceneId, oldPos)
	g.SendMsg(proto.ApiPlayerEnterSceneNotify, userId, player.ClientSeq, playerEnterSceneNotify)

	// PacketSceneTransToPointRsp
	sceneTransToPointRsp := new(proto.SceneTransToPointRsp)
	sceneTransToPointRsp.Retcode = 0
	sceneTransToPointRsp.PointId = req.PointId
	sceneTransToPointRsp.SceneId = req.SceneId
	g.SendMsg(proto.ApiSceneTransToPointRsp, userId, player.ClientSeq, sceneTransToPointRsp)
}

func (g *GameManager) MarkMapReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
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
			newSceneId := req.Mark.SceneId
			oldSceneId := player.SceneId
			oldPos := &model.Vector{
				X: player.Pos.X,
				Y: player.Pos.Y,
				Z: player.Pos.Z,
			}
			jumpScene := false
			if newSceneId != oldSceneId {
				jumpScene = true
			}
			g.RemoveSceneEntityAvatarBroadcastNotify(player)
			world := g.worldManager.GetWorldByID(player.WorldId)
			oldScene := world.GetSceneById(oldSceneId)
			if jumpScene {
				// PacketDelTeamEntityNotify
				delTeamEntityNotify := g.PacketDelTeamEntityNotify(oldScene, player)
				g.SendMsg(proto.ApiDelTeamEntityNotify, player.PlayerID, player.ClientSeq, delTeamEntityNotify)

				oldScene.RemovePlayer(player)
				newScene := world.GetSceneById(newSceneId)
				newScene.AddPlayer(player)
			} else {
				oldScene.UpdatePlayerTeamEntity(player)
			}
			player.Pos.X = float64(req.Mark.Pos.X)
			player.Pos.Y = float64(posYInt)
			player.Pos.Z = float64(req.Mark.Pos.Z)
			player.SceneId = newSceneId
			player.SceneLoadState = model.SceneNone

			// PacketPlayerEnterSceneNotify
			var enterType proto.EnterType
			if jumpScene {
				logger.LOG.Debug("player jump scene, scene: %v, pos: %v", player.SceneId, player.Pos)
				enterType = proto.EnterType_ENTER_TYPE_JUMP
			} else {
				logger.LOG.Debug("player goto scene, scene: %v, pos: %v", player.SceneId, player.Pos)
				enterType = proto.EnterType_ENTER_TYPE_GOTO
			}
			enterReasonConst := constant.GetEnterReasonConst()
			playerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyTp(player, enterType, uint32(enterReasonConst.TransPoint), oldSceneId, oldPos)
			g.SendMsg(proto.ApiPlayerEnterSceneNotify, userId, player.ClientSeq, playerEnterSceneNotify)
		}
	}
}

func (g *GameManager) PathfindingEnterSceneReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user pathfinding enter scene, user id: %v", userId)
	g.SendMsg(proto.ApiPathfindingEnterSceneRsp, userId, player.ClientSeq, new(proto.PathfindingEnterSceneRsp))
}

func (g *GameManager) QueryPathReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	//logger.LOG.Debug("user query path, user id: %v", userId)
	req := payloadMsg.(*proto.QueryPathReq)

	// PacketQueryPathRsp
	queryPathRsp := new(proto.QueryPathRsp)
	queryPathRsp.Corners = []*proto.Vector{req.DestinationPos[0]}
	queryPathRsp.QueryId = req.QueryId
	queryPathRsp.QueryStatus = proto.QueryPathRsp_PATH_STATUS_TYPE_SUCC
	g.SendMsg(proto.ApiQueryPathRsp, userId, player.ClientSeq, queryPathRsp)
}

func (g *GameManager) GetScenePointReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user get scene point, user id: %v", userId)
	req := payloadMsg.(*proto.GetScenePointReq)

	// PacketGetScenePointRsp
	getScenePointRsp := new(proto.GetScenePointRsp)
	getScenePointRsp.SceneId = req.SceneId
	getScenePointRsp.UnlockedPointList = make([]uint32, 0)
	for i := uint32(1); i < 1000; i++ {
		getScenePointRsp.UnlockedPointList = append(getScenePointRsp.UnlockedPointList, i)
	}
	getScenePointRsp.UnlockAreaList = make([]uint32, 0)
	for i := uint32(1); i < 9; i++ {
		getScenePointRsp.UnlockAreaList = append(getScenePointRsp.UnlockAreaList, i)
	}
	g.SendMsg(proto.ApiGetScenePointRsp, userId, player.ClientSeq, getScenePointRsp)
}

func (g *GameManager) GetSceneAreaReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user get scene area, user id: %v", userId)
	req := payloadMsg.(*proto.GetSceneAreaReq)

	// PacketGetSceneAreaRsp
	getSceneAreaRsp := new(proto.GetSceneAreaRsp)
	getSceneAreaRsp.SceneId = req.SceneId
	getSceneAreaRsp.AreaIdList = []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 17, 18, 19, 20, 21, 22, 23, 24, 25, 100, 101, 102, 103, 200, 210, 300, 400, 401, 402, 403}
	getSceneAreaRsp.CityInfoList = make([]*proto.CityInfo, 0)
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 1, Level: 1})
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 2, Level: 1})
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 3, Level: 1})
	g.SendMsg(proto.ApiGetSceneAreaRsp, userId, player.ClientSeq, getSceneAreaRsp)
}
