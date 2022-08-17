package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	"game-genshin/constant"
	"game-genshin/model"
	"time"
)

type TickManager struct {
	ticker      *time.Ticker
	tickCount   uint64
	gameManager *GameManager
}

func NewTickManager(gameManager *GameManager) (r *TickManager) {
	r = new(TickManager)
	r.ticker = time.NewTicker(time.Millisecond * 100)
	logger.LOG.Info("game server tick start at: %v", time.Now().UnixMilli())
	r.gameManager = gameManager
	return r
}

func (t *TickManager) OnGameServerTick() {
	t.tickCount++
	now := time.Now().UnixMilli()
	t.onTick100MilliSecond(now)
	if t.tickCount%(10*1) == 0 {
		t.onTickSecond(now)
	}
	if t.tickCount%(10*10) == 0 {
		t.onTick10Second(now)
	}
	if t.tickCount%(10*60) == 0 {
		t.onTickMinute(now)
	}
	if t.tickCount%(10*3600) == 0 {
		t.onTickHour(now)
	}
	if t.tickCount%(10*3600*24) == 0 {
		t.onTickDay(now)
	}
	if t.tickCount%(10*3600*24*7) == 0 {
		t.onTickWeek(now)
	}
}

func (t *TickManager) onTickWeek(now int64) {
	logger.LOG.Info("on tick week, time: %v", now)
}

func (t *TickManager) onTickDay(now int64) {
	logger.LOG.Info("on tick day, time: %v", now)
}

func (t *TickManager) onTickHour(now int64) {
	logger.LOG.Info("on tick hour, time: %v", now)
}

func (t *TickManager) onTickMinute(now int64) {
	logger.LOG.Info("on tick minute, time: %v", now)
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, player := range world.playerMap {
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 201, ChangeCount: 10}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 202, ChangeCount: 10}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 223, ChangeCount: 1000}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 224, ChangeCount: 1000}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 104003, ChangeCount: 1}}, true)

			t.createMonster(world, player)
		}
	}
}

func (t *TickManager) onTick10Second(now int64) {
	logger.LOG.Info("on tick 10 second, time: %v", now)
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, player := range world.playerMap {
			// PacketWorldPlayerRTTNotify
			worldPlayerRTTNotify := new(proto.WorldPlayerRTTNotify)
			worldPlayerRTTNotify.PlayerRttList = []*proto.PlayerRTTInfo{{Uid: player.PlayerID, Rtt: player.ClientRTT}}
			t.gameManager.SendMsg(api.ApiWorldPlayerRTTNotify, player.PlayerID, nil, worldPlayerRTTNotify)
		}
	}
}

func (t *TickManager) onTickSecond(now int64) {
}

func (t *TickManager) onTick100MilliSecond(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, scene := range world.sceneMap {
			scene.AttackHandler(t.gameManager)
		}
	}
}

func (t *TickManager) createMonster(world *World, player *model.Player) {
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	scene := world.GetSceneById(player.SceneId)
	fightPropertyConst := constant.GetFightPropertyConst()
	monsterEntityId := scene.CreateEntity(entityIdTypeConst.MONSTER, map[uint32]float32{uint32(fightPropertyConst.FIGHT_PROP_CUR_HP): 72.91699}, nil)

	playerPropertyConst := constant.GetPlayerPropertyConst()
	pos := &proto.Vector{
		X: 2747,
		Y: 194,
		Z: -1719,
	}
	// PacketSceneEntityAppearNotify
	sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
	sceneEntityAppearNotify.AppearType = proto.VisionType_VISION_TYPE_BORN
	sceneEntityInfo := &proto.SceneEntityInfo{
		EntityType: proto.ProtEntityType_PROT_ENTITY_TYPE_MONSTER,
		EntityId:   monsterEntityId,
		MotionInfo: &proto.MotionInfo{
			Pos:   pos,
			Rot:   &proto.Vector{},
			Speed: &proto.Vector{},
		},
		PropList: []*proto.PropPair{{Type: uint32(playerPropertyConst.PROP_LEVEL), PropValue: &proto.PropValue{
			Type:  uint32(playerPropertyConst.PROP_LEVEL),
			Value: &proto.PropValue_Ival{Ival: int64(1)},
			Val:   int64(1),
		}}},
		FightPropList: []*proto.FightPropPair{
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_HP),
				PropValue: float32(72.91699),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_PHYSICAL_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE),
				PropValue: float32(505.0),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK),
				PropValue: float32(45.679916),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_ICE_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK),
				PropValue: float32(45.679916),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_MAX_HP),
				PropValue: float32(72.91699),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_FIRE_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_ELEC_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_WIND_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_ROCK_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_GRASS_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_WATER_SUB_HURT),
				PropValue: float32(0.1),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_HP),
				PropValue: float32(72.91699),
			},
			{
				PropType:  uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE),
				PropValue: float32(505.0),
			},
		},
		LifeState:        1,
		AnimatorParaList: make([]*proto.AnimatorParameterValueInfoPair, 0),
		Entity: &proto.SceneEntityInfo_Monster{
			Monster: &proto.SceneMonsterInfo{
				MonsterId:       21010101,
				AuthorityPeerId: 1,
				BornType:        proto.MonsterBornType_MONSTER_BORN_TYPE_DEFAULT,
				BlockId:         3001,
				TitleId:         3001,
				SpecialNameId:   40,
			},
		},
		EntityClientData: new(proto.EntityClientData),
		EntityAuthorityInfo: &proto.EntityAuthorityInfo{
			AbilityInfo:         new(proto.AbilitySyncStateInfo),
			RendererChangedInfo: new(proto.EntityRendererChangedInfo),
			AiInfo: &proto.SceneEntityAiInfo{
				IsAiOpen: true,
				BornPos:  pos,
			},
			BornPos: pos,
		},
	}
	sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{sceneEntityInfo}
	t.gameManager.SendMsg(api.ApiSceneEntityAppearNotify, player.PlayerID, t.gameManager.getHeadMsg(11), sceneEntityAppearNotify)
}
