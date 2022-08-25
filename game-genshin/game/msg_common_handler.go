package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
)

func (g *GameManager) PlayerSetPauseReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user pause, user id: %v", userId)
	if headMsg != nil {
		logger.LOG.Debug("client sequence id: %v", headMsg.ClientSequenceId)
	}
	req := payloadMsg.(*proto.PlayerSetPauseReq)
	isPaused := req.IsPaused
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	player.Pause = isPaused
}

func (g *GameManager) PathfindingEnterSceneReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user pathfinding enter scene, user id: %v", userId)
	g.SendMsg(api.ApiPathfindingEnterSceneRsp, userId, nil, new(proto.NullMsg))
}

func (g *GameManager) GetScenePointReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get scene point, user id: %v", userId)
	req := payloadMsg.(*proto.GetScenePointReq)
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
	g.SendMsg(api.ApiGetScenePointRsp, userId, nil, getScenePointRsp)
}

func (g *GameManager) GetSceneAreaReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get scene area, user id: %v", userId)
	req := payloadMsg.(*proto.GetSceneAreaReq)
	getSceneAreaRsp := new(proto.GetSceneAreaRsp)
	getSceneAreaRsp.SceneId = req.SceneId
	getSceneAreaRsp.AreaIdList = []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 17, 18, 19, 20, 21, 22, 23, 24, 25, 100, 101, 102, 103, 200, 210, 300, 400, 401, 402, 403}
	getSceneAreaRsp.CityInfoList = make([]*proto.CityInfo, 0)
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 1, Level: 1})
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 2, Level: 1})
	getSceneAreaRsp.CityInfoList = append(getSceneAreaRsp.CityInfoList, &proto.CityInfo{CityId: 3, Level: 1})
	g.SendMsg(api.ApiGetSceneAreaRsp, userId, nil, getSceneAreaRsp)
}

func (g *GameManager) EnterWorldAreaReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user enter world area, user id: %v", userId)
	req := payloadMsg.(*proto.EnterWorldAreaReq)
	enterWorldAreaRsp := new(proto.EnterWorldAreaRsp)
	enterWorldAreaRsp.AreaType = req.AreaType
	enterWorldAreaRsp.AreaId = req.AreaId
	g.SendMsg(api.ApiEnterWorldAreaRsp, userId, nil, enterWorldAreaRsp)
}

func (g *GameManager) PostEnterSceneReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user post enter scene, user id: %v", userId)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	postEnterSceneRsp := new(proto.PostEnterSceneRsp)
	postEnterSceneRsp.EnterSceneToken = player.EnterSceneToken
	g.SendMsg(api.ApiPostEnterSceneRsp, userId, nil, postEnterSceneRsp)
}

func (g *GameManager) TowerAllDataReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user get tower all data, user id: %v", userId)
	towerAllDataRsp := new(proto.TowerAllDataRsp)
	towerAllDataRsp.TowerScheduleId = 29
	towerAllDataRsp.TowerFloorRecordList = []*proto.TowerFloorRecord{{FloorId: 1001}}
	towerAllDataRsp.CurLevelRecord = &proto.TowerCurLevelRecord{IsEmpty: true}
	towerAllDataRsp.NextScheduleChangeTime = 4294967295
	towerAllDataRsp.FloorOpenTimeMap = make(map[uint32]uint32)
	towerAllDataRsp.FloorOpenTimeMap[1024] = 1630486800
	towerAllDataRsp.FloorOpenTimeMap[1025] = 1630486800
	towerAllDataRsp.FloorOpenTimeMap[1026] = 1630486800
	towerAllDataRsp.FloorOpenTimeMap[1027] = 1630486800
	towerAllDataRsp.ScheduleStartTime = 1630486800
	g.SendMsg(api.ApiTowerAllDataRsp, userId, nil, towerAllDataRsp)
}

func (g *GameManager) QueryPathReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	//logger.LOG.Debug("user query path, user id: %v", userId)
	req := payloadMsg.(*proto.QueryPathReq)
	queryPathRsp := new(proto.QueryPathRsp)
	queryPathRsp.Corners = []*proto.Vector{req.DestinationPos[0]}
	queryPathRsp.QueryId = req.QueryId
	queryPathRsp.QueryStatus = proto.QueryPathRsp_PATH_STATUS_TYPE_SUCC
	g.SendMsg(api.ApiQueryPathRsp, userId, nil, queryPathRsp)
}

func (g *GameManager) PingReq(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user ping, user id: %v", userId)
	pingReq := payloadMsg.(*proto.PingReq)
	logger.LOG.Debug("ping req data: %v", pingReq.String())
	player := g.userManager.GetOnlineUser(userId)
	if player != nil {
		player.ClientTime = pingReq.ClientTime
		player.ClientRTT = 10
	}
	pingRsp := new(proto.PingRsp)
	pingRsp.ClientTime = pingReq.ClientTime
	g.SendMsg(api.ApiPingRsp, userId, nil, pingRsp)
}

func (g *GameManager) ClientRttNotify(userId uint32, payloadMsg any) {
	rtt := payloadMsg.(int32)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	logger.LOG.Debug("user rtt notify, user id: %v, rtt: %v", userId, rtt)
	player.ClientRTT = uint32(rtt)
}

func (g *GameManager) EntityAiSyncNotify(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	logger.LOG.Debug("user entity ai sync, user id: %v", userId)
	if payloadMsg == nil {
		return
	}
	req := payloadMsg.(*proto.EntityAiSyncNotify)
	if len(req.LocalAvatarAlertedMonsterList) == 0 {
		return
	}

	// PacketEntityAiSyncNotify
	entityAiSyncNotify := new(proto.EntityAiSyncNotify)
	entityAiSyncNotify.InfoList = make([]*proto.AiSyncInfo, 0)
	for _, monsterId := range req.LocalAvatarAlertedMonsterList {
		entityAiSyncNotify.InfoList = append(entityAiSyncNotify.InfoList, &proto.AiSyncInfo{
			EntityId:        monsterId,
			HasPathToTarget: true,
			IsSelfKilling:   false,
		})
	}
	g.SendMsg(api.ApiEntityAiSyncNotify, userId, nil, entityAiSyncNotify)
}
