package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
)

func (g *GameManager) CombatInvocationsNotify(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	//logger.LOG.Debug("user combat invocations, user id: %v", userId)
	req := payloadMsg.(*proto.CombatInvocationsNotify)
	//logger.LOG.Debug("req CombatInvocationsNotify: %v", req)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, user id: %v", userId)
		return
	}
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world == nil {
		return
	}
	scene := world.GetSceneById(player.SceneId)
	invokeHandler := NewInvokeHandler[proto.CombatInvokeEntry]()
	for _, entry := range req.InvokeList {
		switch entry.ArgumentType {
		case proto.CombatTypeArgument_COMBAT_TYPE_ARGUMENT_ENTITY_MOVE:
			entityMoveInfo := new(proto.EntityMoveInfo)
			err := pb.Unmarshal(entry.CombatData, entityMoveInfo)
			if err != nil {
				logger.LOG.Error("parse combat invocations entity move info error: %v", err)
				continue
			}

			motionInfo := entityMoveInfo.MotionInfo

			sceneEntity := scene.GetEntity(entityMoveInfo.EntityId)
			if sceneEntity != nil {
				if motionInfo.Pos != nil && motionInfo.Rot != nil {
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
				}
				sceneEntity.moveState = uint16(motionInfo.State)
				sceneEntity.lastMoveSceneTimeMs = entityMoveInfo.SceneTime
				sceneEntity.lastMoveReliableSeq = entityMoveInfo.ReliableSeq
			}
			if entityMoveInfo.EntityId == player.TeamConfig.GetActiveAvatarEntity().AvatarEntityId {
				// 玩家在移动
				if motionInfo.Pos != nil && motionInfo.Rot != nil {
					player.Pos.X = float64(motionInfo.Pos.X)
					player.Pos.Y = float64(motionInfo.Pos.Y)
					player.Pos.Z = float64(motionInfo.Pos.Z)
					player.Rot.X = float64(motionInfo.Rot.X)
					player.Rot.Y = float64(motionInfo.Rot.Y)
					player.Rot.Z = float64(motionInfo.Rot.Z)
				}
			}
			invokeHandler.addEntry(entry.ForwardType, entry)
		case proto.CombatTypeArgument_COMBAT_TYPE_ARGUMENT_EVT_BEING_HIT:
			scene.AddAttack(&Attack{
				combatInvokeEntry: entry,
				uid:               player.PlayerID,
			})
		}
	}

	// PacketCombatInvocationsNotify
	if invokeHandler.AllLen() > 0 {
		combatInvocationsNotify := new(proto.CombatInvocationsNotify)
		combatInvocationsNotify.InvokeList = invokeHandler.entryListForwardAll
		for _, v := range scene.playerMap {
			g.SendMsg(api.ApiCombatInvocationsNotify, v.PlayerID, nil, combatInvocationsNotify)
		}
	}
	if invokeHandler.AllExceptCurLen() > 0 {
		combatInvocationsNotify := new(proto.CombatInvocationsNotify)
		combatInvocationsNotify.InvokeList = invokeHandler.entryListForwardAllExceptCur
		for _, v := range scene.playerMap {
			if player.PlayerID == v.PlayerID {
				continue
			}
			g.SendMsg(api.ApiCombatInvocationsNotify, v.PlayerID, nil, combatInvocationsNotify)
		}
	}
	if invokeHandler.HostLen() > 0 {
		combatInvocationsNotify := new(proto.CombatInvocationsNotify)
		combatInvocationsNotify.InvokeList = invokeHandler.entryListForwardHost
		g.SendMsg(api.ApiCombatInvocationsNotify, world.owner.PlayerID, nil, combatInvocationsNotify)
	}
}

func (g *GameManager) AbilityInvocationsNotify(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	//logger.LOG.Debug("user ability invocations, user id: %v", userId)
	req := payloadMsg.(*proto.AbilityInvocationsNotify)
	//logger.LOG.Debug("req AbilityInvocationsNotify: %v", req)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, user id: %v", userId)
		return
	}
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world == nil {
		return
	}
	scene := world.GetSceneById(player.SceneId)
	invokeHandler := NewInvokeHandler[proto.AbilityInvokeEntry]()
	for _, entry := range req.Invokes {
		invokeHandler.addEntry(entry.ForwardType, entry)
	}

	// PacketAbilityInvocationsNotify
	if invokeHandler.AllLen() > 0 {
		abilityInvocationsNotify := new(proto.AbilityInvocationsNotify)
		abilityInvocationsNotify.Invokes = invokeHandler.entryListForwardAll
		for _, v := range scene.playerMap {
			g.SendMsg(api.ApiAbilityInvocationsNotify, v.PlayerID, nil, abilityInvocationsNotify)
		}
	}
	if invokeHandler.AllExceptCurLen() > 0 {
		abilityInvocationsNotify := new(proto.AbilityInvocationsNotify)
		abilityInvocationsNotify.Invokes = invokeHandler.entryListForwardAllExceptCur
		for _, v := range scene.playerMap {
			if player.PlayerID == v.PlayerID {
				continue
			}
			g.SendMsg(api.ApiAbilityInvocationsNotify, v.PlayerID, nil, abilityInvocationsNotify)
		}
	}
	if invokeHandler.HostLen() > 0 {
		abilityInvocationsNotify := new(proto.AbilityInvocationsNotify)
		abilityInvocationsNotify.Invokes = invokeHandler.entryListForwardHost
		g.SendMsg(api.ApiAbilityInvocationsNotify, world.owner.PlayerID, nil, abilityInvocationsNotify)
	}
}

func (g *GameManager) ClientAbilityInitFinishNotify(userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	//logger.LOG.Debug("user client ability ok, user id: %v", userId)
	req := payloadMsg.(*proto.ClientAbilityInitFinishNotify)
	//logger.LOG.Debug("req ClientAbilityInitFinishNotify: %v", req)
	player := g.userManager.GetOnlineUser(userId)
	if player == nil {
		logger.LOG.Error("player is nil, user id: %v", userId)
		return
	}
	world := g.worldManager.GetWorldByID(player.WorldId)
	if world == nil {
		return
	}
	scene := world.GetSceneById(player.SceneId)
	invokeHandler := NewInvokeHandler[proto.AbilityInvokeEntry]()
	for _, entry := range req.Invokes {
		invokeHandler.addEntry(entry.ForwardType, entry)
	}

	// PacketClientAbilityInitFinishNotify
	if invokeHandler.AllLen() > 0 {
		clientAbilityInitFinishNotify := new(proto.ClientAbilityInitFinishNotify)
		clientAbilityInitFinishNotify.Invokes = invokeHandler.entryListForwardAll
		for _, v := range scene.playerMap {
			g.SendMsg(api.ApiClientAbilityInitFinishNotify, v.PlayerID, nil, clientAbilityInitFinishNotify)
		}
	}
	if invokeHandler.AllExceptCurLen() > 0 {
		clientAbilityInitFinishNotify := new(proto.ClientAbilityInitFinishNotify)
		clientAbilityInitFinishNotify.Invokes = invokeHandler.entryListForwardAllExceptCur
		for _, v := range scene.playerMap {
			if player.PlayerID == v.PlayerID {
				continue
			}
			g.SendMsg(api.ApiClientAbilityInitFinishNotify, v.PlayerID, nil, clientAbilityInitFinishNotify)
		}
	}
	if invokeHandler.HostLen() > 0 {
		clientAbilityInitFinishNotify := new(proto.ClientAbilityInitFinishNotify)
		clientAbilityInitFinishNotify.Invokes = invokeHandler.entryListForwardHost
		g.SendMsg(api.ApiClientAbilityInitFinishNotify, world.owner.PlayerID, nil, clientAbilityInitFinishNotify)
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