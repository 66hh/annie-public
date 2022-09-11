package game

import (
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
)

func (g *GameManager) CombatInvocationsNotify(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	//logger.LOG.Debug("user combat invocations, user id: %v", userId)
	req := payloadMsg.(*proto.CombatInvocationsNotify)
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world == nil {
		return
	}
	scene := world.GetSceneById(player.SceneId)
	invokeHandler := NewInvokeHandler[proto.CombatInvokeEntry]()
	for _, entry := range req.InvokeList {
		//logger.LOG.Debug("AT: %v, FT: %v, UID: %v", entry.ArgumentType, entry.ForwardType, player.PlayerID)
		switch entry.ArgumentType {
		case proto.CombatTypeArgument_COMBAT_TYPE_ARGUMENT_EVT_BEING_HIT:
			scene.AddAttack(&Attack{
				combatInvokeEntry: entry,
				uid:               player.PlayerID,
			})
		case proto.CombatTypeArgument_COMBAT_TYPE_ARGUMENT_ENTITY_MOVE:
			entityMoveInfo := new(proto.EntityMoveInfo)
			err := pb.Unmarshal(entry.CombatData, entityMoveInfo)
			if err != nil {
				logger.LOG.Error("parse combat invocations entity move info error: %v", err)
				continue
			}

			motionInfo := entityMoveInfo.MotionInfo

			if motionInfo.Pos == nil || motionInfo.Rot == nil {
				continue
			}

			sceneEntity := scene.GetEntity(entityMoveInfo.EntityId)
			if sceneEntity != nil {
				sceneEntity.pos = &model.Vector{
					X: float64(motionInfo.Pos.X),
					Y: float64(motionInfo.Pos.Y),
					Z: float64(motionInfo.Pos.Z),
				}
				sceneEntity.rot = &model.Vector{
					X: float64(motionInfo.Rot.X),
					Y: float64(motionInfo.Rot.Y),
					Z: float64(motionInfo.Rot.Z),
				}
				sceneEntity.moveState = uint16(motionInfo.State)
				sceneEntity.lastMoveSceneTimeMs = entityMoveInfo.SceneTime
				sceneEntity.lastMoveReliableSeq = entityMoveInfo.ReliableSeq
				//logger.LOG.Debug("entity move, id: %v, pos: %v, uid: %v", sceneEntity.id, sceneEntity.pos, player.PlayerID)
			}
			activeAvatarId := player.TeamConfig.GetActiveAvatarId()
			playerTeamEntity := scene.GetPlayerTeamEntity(player.PlayerID)
			if playerTeamEntity != nil && entityMoveInfo.EntityId == playerTeamEntity.avatarEntityMap[activeAvatarId] {
				// 玩家在移动
				team := player.TeamConfig.GetActiveTeam()
				for _, avatarId := range team.AvatarIdList {
					if avatarId == activeAvatarId {
						continue
					}
					entityId := playerTeamEntity.avatarEntityMap[avatarId]
					entity := scene.GetEntity(entityId)
					if entity != nil {
						entity.pos.X = float64(motionInfo.Pos.X)
						entity.pos.Y = float64(motionInfo.Pos.Y)
						entity.pos.Z = float64(motionInfo.Pos.Z)
						entity.rot.X = float64(motionInfo.Rot.X)
						entity.rot.Y = float64(motionInfo.Rot.Y)
						entity.rot.Z = float64(motionInfo.Rot.Z)
					}
				}
				player.Pos.X = float64(motionInfo.Pos.X)
				player.Pos.Y = float64(motionInfo.Pos.Y)
				player.Pos.Z = float64(motionInfo.Pos.Z)
				player.Rot.X = float64(motionInfo.Rot.X)
				player.Rot.Y = float64(motionInfo.Rot.Y)
				player.Rot.Z = float64(motionInfo.Rot.Z)
			}
			invokeHandler.addEntry(entry.ForwardType, entry)
		default:
			invokeHandler.addEntry(entry.ForwardType, entry)
		}
	}

	// PacketCombatInvocationsNotify
	if invokeHandler.AllLen() > 0 {
		combatInvocationsNotify := new(proto.CombatInvocationsNotify)
		combatInvocationsNotify.InvokeList = invokeHandler.entryListForwardAll
		for _, v := range scene.playerMap {
			g.SendMsg(proto.ApiCombatInvocationsNotify, v.PlayerID, v.ClientSeq, combatInvocationsNotify)
		}
	}
	if invokeHandler.AllExceptCurLen() > 0 {
		combatInvocationsNotify := new(proto.CombatInvocationsNotify)
		combatInvocationsNotify.InvokeList = invokeHandler.entryListForwardAllExceptCur
		for _, v := range scene.playerMap {
			if player.PlayerID == v.PlayerID {
				continue
			}
			g.SendMsg(proto.ApiCombatInvocationsNotify, v.PlayerID, v.ClientSeq, combatInvocationsNotify)
		}
	}
	if invokeHandler.HostLen() > 0 {
		combatInvocationsNotify := new(proto.CombatInvocationsNotify)
		combatInvocationsNotify.InvokeList = invokeHandler.entryListForwardHost
		g.SendMsg(proto.ApiCombatInvocationsNotify, world.owner.PlayerID, world.owner.ClientSeq, combatInvocationsNotify)
	}
}

func (g *GameManager) AbilityInvocationsNotify(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	//logger.LOG.Debug("user ability invocations, user id: %v", userId)
	req := payloadMsg.(*proto.AbilityInvocationsNotify)
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world == nil {
		return
	}
	scene := world.GetSceneById(player.SceneId)
	invokeHandler := NewInvokeHandler[proto.AbilityInvokeEntry]()
	for _, entry := range req.Invokes {
		//logger.LOG.Debug("AT: %v, FT: %v, UID: %v", entry.ArgumentType, entry.ForwardType, player.PlayerID)
		invokeHandler.addEntry(entry.ForwardType, entry)
	}

	// PacketAbilityInvocationsNotify
	if invokeHandler.AllLen() > 0 {
		abilityInvocationsNotify := new(proto.AbilityInvocationsNotify)
		abilityInvocationsNotify.Invokes = invokeHandler.entryListForwardAll
		for _, v := range scene.playerMap {
			g.SendMsg(proto.ApiAbilityInvocationsNotify, v.PlayerID, v.ClientSeq, abilityInvocationsNotify)
		}
	}
	if invokeHandler.AllExceptCurLen() > 0 {
		abilityInvocationsNotify := new(proto.AbilityInvocationsNotify)
		abilityInvocationsNotify.Invokes = invokeHandler.entryListForwardAllExceptCur
		for _, v := range scene.playerMap {
			if player.PlayerID == v.PlayerID {
				continue
			}
			g.SendMsg(proto.ApiAbilityInvocationsNotify, v.PlayerID, v.ClientSeq, abilityInvocationsNotify)
		}
	}
	if invokeHandler.HostLen() > 0 {
		abilityInvocationsNotify := new(proto.AbilityInvocationsNotify)
		abilityInvocationsNotify.Invokes = invokeHandler.entryListForwardHost
		g.SendMsg(proto.ApiAbilityInvocationsNotify, world.owner.PlayerID, world.owner.ClientSeq, abilityInvocationsNotify)
	}
}

func (g *GameManager) ClientAbilityInitFinishNotify(userId uint32, player *model.Player, clientSeq uint32, payloadMsg pb.Message) {
	//logger.LOG.Debug("user client ability ok, user id: %v", userId)
	req := payloadMsg.(*proto.ClientAbilityInitFinishNotify)
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world == nil {
		return
	}
	scene := world.GetSceneById(player.SceneId)
	invokeHandler := NewInvokeHandler[proto.AbilityInvokeEntry]()
	for _, entry := range req.Invokes {
		//logger.LOG.Debug("AT: %v, FT: %v, UID: %v", entry.ArgumentType, entry.ForwardType, player.PlayerID)
		invokeHandler.addEntry(entry.ForwardType, entry)
	}

	// PacketClientAbilityInitFinishNotify
	if invokeHandler.AllLen() > 0 {
		clientAbilityInitFinishNotify := new(proto.ClientAbilityInitFinishNotify)
		clientAbilityInitFinishNotify.Invokes = invokeHandler.entryListForwardAll
		for _, v := range scene.playerMap {
			g.SendMsg(proto.ApiClientAbilityInitFinishNotify, v.PlayerID, v.ClientSeq, clientAbilityInitFinishNotify)
		}
	}
	if invokeHandler.AllExceptCurLen() > 0 {
		clientAbilityInitFinishNotify := new(proto.ClientAbilityInitFinishNotify)
		clientAbilityInitFinishNotify.Invokes = invokeHandler.entryListForwardAllExceptCur
		for _, v := range scene.playerMap {
			if player.PlayerID == v.PlayerID {
				continue
			}
			g.SendMsg(proto.ApiClientAbilityInitFinishNotify, v.PlayerID, v.ClientSeq, clientAbilityInitFinishNotify)
		}
	}
	if invokeHandler.HostLen() > 0 {
		clientAbilityInitFinishNotify := new(proto.ClientAbilityInitFinishNotify)
		clientAbilityInitFinishNotify.Invokes = invokeHandler.entryListForwardHost
		g.SendMsg(proto.ApiClientAbilityInitFinishNotify, world.owner.PlayerID, world.owner.ClientSeq, clientAbilityInitFinishNotify)
	}
}

type InvokeType interface {
	proto.AbilityInvokeEntry | proto.CombatInvokeEntry
}

type InvokeHandler[T InvokeType] struct {
	entryListForwardAll          []*T
	entryListForwardAllExceptCur []*T
	entryListForwardHost         []*T
}

func NewInvokeHandler[T InvokeType]() (r *InvokeHandler[T]) {
	r = new(InvokeHandler[T])
	r.InitInvokeHandler()
	return r
}

func (i *InvokeHandler[T]) InitInvokeHandler() {
	i.entryListForwardAll = make([]*T, 0)
	i.entryListForwardAllExceptCur = make([]*T, 0)
	i.entryListForwardHost = make([]*T, 0)
}

func (i *InvokeHandler[T]) addEntry(forward proto.ForwardType, entry *T) {
	switch forward {
	case proto.ForwardType_FORWARD_TYPE_TO_ALL:
		i.entryListForwardAll = append(i.entryListForwardAll, entry)
	case proto.ForwardType_FORWARD_TYPE_TO_ALL_EXCEPT_CUR:
		fallthrough
	case proto.ForwardType_FORWARD_TYPE_TO_ALL_EXIST_EXCEPT_CUR:
		i.entryListForwardAllExceptCur = append(i.entryListForwardAllExceptCur, entry)
	case proto.ForwardType_FORWARD_TYPE_TO_HOST:
		i.entryListForwardHost = append(i.entryListForwardHost, entry)
	default:
		if forward != proto.ForwardType_FORWARD_TYPE_ONLY_SERVER {
			logger.LOG.Error("forward: %v, entry: %v", forward, entry)
		}
	}
}

func (i *InvokeHandler[T]) AllLen() int {
	return len(i.entryListForwardAll)
}

func (i *InvokeHandler[T]) AllExceptCurLen() int {
	return len(i.entryListForwardAllExceptCur)
}

func (i *InvokeHandler[T]) HostLen() int {
	return len(i.entryListForwardHost)
}
