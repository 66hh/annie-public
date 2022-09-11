package game

import (
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
)

func (g *GameManager) PlayerSetPauseReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user pause, user id: %v", userId)
	req := payloadMsg.(*proto.PlayerSetPauseReq)
	isPaused := req.IsPaused

	player.Pause = isPaused

	// PacketPlayerSetPauseRsp
	playerSetPauseRsp := new(proto.PlayerSetPauseRsp)
	g.SendMsg(proto.ApiPlayerSetPauseRsp, userId, player.ClientSeq, playerSetPauseRsp)
}

func (g *GameManager) TowerAllDataReq(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user get tower all data, user id: %v", userId)

	// PacketTowerAllDataRsp
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
	g.SendMsg(proto.ApiTowerAllDataRsp, userId, player.ClientSeq, towerAllDataRsp)
}

func (g *GameManager) EntityAiSyncNotify(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	logger.LOG.Debug("user entity ai sync, user id: %v", userId)
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
	g.SendMsg(proto.ApiEntityAiSyncNotify, userId, player.ClientSeq, entityAiSyncNotify)
}

func (g *GameManager) ClientTimeNotify(userId uint32, clientTime uint32) {
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	logger.LOG.Debug("client time notify, user id: %v, time: %v", userId, clientTime)
	player.ClientTime = clientTime
}

func (g *GameManager) ClientRttNotify(userId uint32, clientRtt uint32) {
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, userId: %v", userId)
		return
	}
	logger.LOG.Debug("client rtt notify, user id: %v, rtt: %v", userId, clientRtt)
	player.ClientRTT = clientRtt
}
