package game

import (
	"flswld.com/common/utils/alg"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"game-genshin/constant"
	"game-genshin/model"
	pb "google.golang.org/protobuf/proto"
)

type WorldManager struct {
	worldMap  map[uint32]*World
	snowflake *alg.SnowflakeWorker
	bigWorld  *World
}

func NewWorldManager(snowflake *alg.SnowflakeWorker) (r *WorldManager) {
	r = new(WorldManager)
	r.worldMap = make(map[uint32]*World)
	r.snowflake = snowflake
	return r
}

func (w *WorldManager) GetWorldByID(worldId uint32) *World {
	return w.worldMap[worldId]
}

func (w *WorldManager) GetWorldMap() map[uint32]*World {
	return w.worldMap
}

func (w *WorldManager) CreateWorld(owner *model.Player, multiplayer bool) *World {
	worldId := uint32(w.snowflake.GenId())
	world := &World{
		id:              worldId,
		owner:           owner,
		playerMap:       make(map[uint32]*model.Player),
		sceneMap:        make(map[uint32]*Scene),
		entityIdCounter: 0,
		worldLevel:      0,
		multiplayer:     multiplayer,
		mpLevelEntityId: 0,
		chatMsgList:     make([]*proto.ChatInfo, 0),
	}
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	world.mpLevelEntityId = world.GetNextWorldEntityId(entityIdTypeConst.MPLEVEL)
	w.worldMap[worldId] = world
	return world
}

func (w *WorldManager) DestroyWorld(worldId uint32) {
	world := w.GetWorldByID(worldId)
	for _, player := range world.playerMap {
		world.RemovePlayer(player)
		player.WorldId = 0
	}
	delete(w.worldMap, worldId)
}

func (w *WorldManager) InitBigWorld(owner *model.Player) {
	w.bigWorld = w.GetWorldByID(owner.WorldId)
	w.bigWorld.multiplayer = true
}

type World struct {
	id              uint32
	owner           *model.Player
	playerMap       map[uint32]*model.Player
	sceneMap        map[uint32]*Scene
	entityIdCounter uint32
	worldLevel      uint8
	multiplayer     bool
	mpLevelEntityId uint32
	chatMsgList     []*proto.ChatInfo
}

func (w *World) GetNextWorldEntityId(entityType uint16) uint32 {
	w.entityIdCounter++
	ret := (uint32(entityType) << 24) + w.entityIdCounter
	return ret
}

func (w *World) AddPlayer(player *model.Player, sceneId uint32) {
	player.PeerId = uint32(len(w.playerMap) + 1)
	w.playerMap[player.PlayerID] = player
	scene := w.GetSceneById(sceneId)
	scene.AddPlayer(player)
}

func (w *World) RemovePlayer(player *model.Player) {
	scene := w.sceneMap[player.SceneId]
	scene.RemovePlayer(player)
	delete(w.playerMap, player.PlayerID)
}

func (w *World) CreateScene(sceneId uint32) *Scene {
	scene := &Scene{
		id:                  sceneId,
		world:               w,
		playerMap:           make(map[uint32]*model.Player),
		entityMap:           make(map[uint32]*Entity),
		playerTeamEntityMap: make(map[uint32]*PlayerTeamEntity),
		time:                18 * 60,
		attackQueue:         alg.NewQueue(),
	}
	w.sceneMap[sceneId] = scene
	return scene
}

func (w *World) GetSceneById(sceneId uint32) *Scene {
	scene, exist := w.sceneMap[sceneId]
	if !exist {
		scene = w.CreateScene(sceneId)
	}
	return scene
}

func (w *World) AddChat(chatInfo *proto.ChatInfo) {
	w.chatMsgList = append(w.chatMsgList, chatInfo)
}

func (w *World) GetChatList() []*proto.ChatInfo {
	return w.chatMsgList
}

func (w *World) IsBigWorld() bool {
	return w.owner.PlayerID == 1
}

type Scene struct {
	id                  uint32
	world               *World
	playerMap           map[uint32]*model.Player
	entityMap           map[uint32]*Entity
	playerTeamEntityMap map[uint32]*PlayerTeamEntity
	time                uint32
	attackQueue         *alg.Queue
}

type Entity struct {
	id                  uint32
	scene               *Scene
	pos                 *model.Vector
	rot                 *model.Vector
	moveState           uint16
	lastMoveSceneTimeMs uint32
	lastMoveReliableSeq uint32
	fightProp           map[uint32]float32
	entityType          uint32
	uid                 uint32
	avatarId            uint32
	level               uint8
}

type PlayerTeamEntity struct {
	teamEntityId    uint32
	avatarEntityMap map[uint32]uint32
	weaponEntityMap map[uint64]uint32
}

type Attack struct {
	combatInvokeEntry *proto.CombatInvokeEntry
	uid               uint32
}

func (s *Scene) ChangeTime(time uint32) {
	s.time = time % 1440
}

func (s *Scene) GetPlayerTeamEntity(userId uint32) *PlayerTeamEntity {
	return s.playerTeamEntityMap[userId]
}

func (s *Scene) CreatePlayerTeamEntity(player *model.Player) {
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	playerTeamEntity := &PlayerTeamEntity{
		teamEntityId:    s.world.GetNextWorldEntityId(entityIdTypeConst.TEAM),
		avatarEntityMap: make(map[uint32]uint32),
		weaponEntityMap: make(map[uint64]uint32),
	}
	s.playerTeamEntityMap[player.PlayerID] = playerTeamEntity
}

func (s *Scene) UpdatePlayerTeamEntity(player *model.Player) {
	team := player.TeamConfig.GetActiveTeam()
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	playerTeamEntity := s.playerTeamEntityMap[player.PlayerID]
	for _, avatarId := range team.AvatarIdList {
		if avatarId == 0 {
			break
		}
		avatar := player.AvatarMap[avatarId]
		s.DestroyEntity(playerTeamEntity.avatarEntityMap[avatarId])
		playerTeamEntity.avatarEntityMap[avatarId] = s.CreateEntityAvatar(entityIdTypeConst.AVATAR, player, avatarId)
		s.DestroyEntity(playerTeamEntity.weaponEntityMap[avatar.EquipWeapon.WeaponId])
		playerTeamEntity.weaponEntityMap[avatar.EquipWeapon.WeaponId] = s.CreateEntityWeapon(entityIdTypeConst.WEAPON)
	}
}

func (s *Scene) AddPlayer(player *model.Player) {
	s.playerMap[player.PlayerID] = player
	s.CreatePlayerTeamEntity(player)
	s.UpdatePlayerTeamEntity(player)
}

func (s *Scene) RemovePlayer(player *model.Player) {
	playerTeamEntity := s.GetPlayerTeamEntity(player.PlayerID)
	for _, avatarEntityId := range playerTeamEntity.avatarEntityMap {
		s.DestroyEntity(avatarEntityId)
	}
	for _, weaponEntityId := range playerTeamEntity.weaponEntityMap {
		s.DestroyEntity(weaponEntityId)
	}
	delete(s.playerTeamEntityMap, player.PlayerID)
	delete(s.playerMap, player.PlayerID)
}

func (s *Scene) CreateEntityAvatar(entityType uint16, player *model.Player, avatarId uint32) uint32 {
	entityId := s.world.GetNextWorldEntityId(entityType)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		pos:                 player.Pos,
		rot:                 player.Rot,
		moveState:           uint16(proto.MotionState_MOTION_STATE_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp:           player.AvatarMap[avatarId].FightPropMap,
		entityType:          uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_AVATAR),
		uid:                 player.PlayerID,
		avatarId:            avatarId,
		level:               player.AvatarMap[avatarId].Level,
	}
	s.entityMap[entity.id] = entity
	return entity.id
}

func (s *Scene) CreateEntityWeapon(entityType uint16) uint32 {
	entityId := s.world.GetNextWorldEntityId(entityType)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		pos:                 new(model.Vector),
		rot:                 new(model.Vector),
		moveState:           uint16(proto.MotionState_MOTION_STATE_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp:           nil,
		entityType:          uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_WEAPON),
		uid:                 0,
		avatarId:            0,
		level:               0,
	}
	s.entityMap[entity.id] = entity
	return entity.id
}

func (s *Scene) CreateEntityMonster(entityType uint16, pos *model.Vector, level uint8, fightProp map[uint32]float32) uint32 {
	entityId := s.world.GetNextWorldEntityId(entityType)
	entity := &Entity{
		id:                  entityId,
		scene:               s,
		pos:                 pos,
		rot:                 new(model.Vector),
		moveState:           uint16(proto.MotionState_MOTION_STATE_NONE),
		lastMoveSceneTimeMs: 0,
		lastMoveReliableSeq: 0,
		fightProp:           fightProp,
		entityType:          uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_MONSTER),
		uid:                 0,
		avatarId:            0,
		level:               level,
	}
	s.entityMap[entity.id] = entity
	return entity.id
}

func (s *Scene) DestroyEntity(entityId uint32) {
	delete(s.entityMap, entityId)
}

func (s *Scene) GetEntity(entityId uint32) *Entity {
	return s.entityMap[entityId]
}

func (s *Scene) AddAttack(attack *Attack) {
	s.attackQueue.EnQueue(attack)
}

func (s *Scene) AttackHandler(gameManager *GameManager) {
	combatInvokeEntryListAll := make([]*proto.CombatInvokeEntry, 0)
	combatInvokeEntryListOther := make(map[uint32][]*proto.CombatInvokeEntry)
	combatInvokeEntryListHost := make([]*proto.CombatInvokeEntry, 0)

	for s.attackQueue.Len() != 0 {
		value := s.attackQueue.DeQueue()
		attack, ok := value.(*Attack)
		if !ok {
			logger.LOG.Error("error attack type, attack value: %v", value)
			continue
		}
		if attack.combatInvokeEntry == nil {
			logger.LOG.Error("error attack data, attack value: %v", value)
			continue
		}

		hitInfo := new(proto.EvtBeingHitInfo)
		err := pb.Unmarshal(attack.combatInvokeEntry.CombatData, hitInfo)
		if err != nil {
			logger.LOG.Error("parse combat invocations entity hit info error: %v", err)
			continue
		}

		attackResult := hitInfo.AttackResult
		//logger.LOG.Debug("run attack handler, attackResult: %v", attackResult)
		target := s.entityMap[attackResult.DefenseId]
		if target == nil {
			logger.LOG.Error("could not found target, defense id: %v", attackResult.DefenseId)
			continue
		}
		attackResult.Damage *= 100
		damage := attackResult.Damage
		attackerId := attackResult.AttackerId
		_ = attackerId
		fightPropertyConst := constant.GetFightPropertyConst()
		currHp := float32(0)
		if target.fightProp != nil {
			currHp = target.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)]
			currHp -= damage
			if currHp < 0 {
				currHp = 0
			}
			target.fightProp[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)] = currHp
		}

		// PacketEntityFightPropUpdateNotify
		entityFightPropUpdateNotify := new(proto.EntityFightPropUpdateNotify)
		entityFightPropUpdateNotify.EntityId = target.id
		entityFightPropUpdateNotify.FightPropMap = make(map[uint32]float32)
		entityFightPropUpdateNotify.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_HP)] = currHp
		for _, player := range s.playerMap {
			gameManager.SendMsg(proto.ApiEntityFightPropUpdateNotify, player.PlayerID, player.ClientSeq, entityFightPropUpdateNotify)
		}

		combatData, err := pb.Marshal(hitInfo)
		if err != nil {
			logger.LOG.Error("create combat invocations entity hit info error: %v", err)
		}
		attack.combatInvokeEntry.CombatData = combatData
		switch attack.combatInvokeEntry.ForwardType {
		case proto.ForwardType_FORWARD_TYPE_TO_ALL:
			combatInvokeEntryListAll = append(combatInvokeEntryListAll, attack.combatInvokeEntry)
		case proto.ForwardType_FORWARD_TYPE_TO_ALL_EXCEPT_CUR:
			fallthrough
		case proto.ForwardType_FORWARD_TYPE_TO_ALL_EXIST_EXCEPT_CUR:
			if combatInvokeEntryListOther[attack.uid] == nil {
				combatInvokeEntryListOther[attack.uid] = make([]*proto.CombatInvokeEntry, 0)
			}
			combatInvokeEntryListOther[attack.uid] = append(combatInvokeEntryListOther[attack.uid], attack.combatInvokeEntry)
		case proto.ForwardType_FORWARD_TYPE_TO_HOST:
			combatInvokeEntryListHost = append(combatInvokeEntryListHost, attack.combatInvokeEntry)
		default:
		}
	}

	// PacketCombatInvocationsNotify
	if len(combatInvokeEntryListAll) > 0 {
		combatInvocationsNotifyAll := new(proto.CombatInvocationsNotify)
		combatInvocationsNotifyAll.InvokeList = combatInvokeEntryListAll
		for _, player := range s.playerMap {
			gameManager.SendMsg(proto.ApiCombatInvocationsNotify, player.PlayerID, player.ClientSeq, combatInvocationsNotifyAll)
		}
	}
	if len(combatInvokeEntryListOther) > 0 {
		for uid, list := range combatInvokeEntryListOther {
			combatInvocationsNotifyOther := new(proto.CombatInvocationsNotify)
			combatInvocationsNotifyOther.InvokeList = list
			for _, player := range s.playerMap {
				if player.PlayerID == uid {
					continue
				}
				gameManager.SendMsg(proto.ApiCombatInvocationsNotify, player.PlayerID, player.ClientSeq, combatInvocationsNotifyOther)
			}
		}
	}
	if len(combatInvokeEntryListHost) > 0 {
		combatInvocationsNotifyHost := new(proto.CombatInvocationsNotify)
		combatInvocationsNotifyHost.InvokeList = combatInvokeEntryListHost
		gameManager.SendMsg(proto.ApiCombatInvocationsNotify, s.world.owner.PlayerID, s.world.owner.ClientSeq, combatInvocationsNotifyHost)
	}
}
