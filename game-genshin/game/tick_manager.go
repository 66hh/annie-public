package game

import (
	"flswld.com/common/utils/random"
	"flswld.com/gate-genshin-api/proto"
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
	if t.tickCount%(10*5) == 0 {
		t.onTick5Second(now)
	}
	if t.tickCount%(10*10) == 0 {
		t.onTick10Second(now)
	}
	if t.tickCount%(10*60) == 0 {
		t.onTickMinute(now)
	}
	if t.tickCount%(10*60*10) == 0 {
		t.onTick10Minute(now)
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

func (t *TickManager) onTick10Minute(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, player := range world.playerMap {
			// 蓝球粉球
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 223, ChangeCount: 1}}, true, 0)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 224, ChangeCount: 1}}, true, 0)
		}
	}
}

func (t *TickManager) onTickMinute(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, player := range world.playerMap {
			// 随机物品
			allItemDataConfig := t.gameManager.GetAllItemDataConfig()
			count := random.GetRandomInt32(0, 4)
			i := int32(0)
			itemTypeConst := constant.GetItemTypeConst()
			for itemId := range allItemDataConfig {
				itemDataConfig := allItemDataConfig[itemId]
				// TODO 3.0.0REL版本中 发送某些无效家具 可能会导致客户端背包家具界面卡死
				if itemDataConfig.ItemEnumType == itemTypeConst.ITEM_FURNITURE {
					continue
				}
				num := random.GetRandomInt32(1, 9)
				t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: uint32(itemId), ChangeCount: uint32(num)}}, true, 0)
				i++
				if i > count {
					break
				}
			}
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 102, ChangeCount: 30}}, true, 0)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 201, ChangeCount: 10}}, true, 0)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 202, ChangeCount: 100}}, true, 0)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 203, ChangeCount: 10}}, true, 0)
		}
	}
}

func (t *TickManager) onTick10Second(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		if !world.IsBigWorld() && (world.multiplayer || !world.owner.Pause) {
			// 刷怪
			scene := world.GetSceneById(3)
			monsterEntityCount := 0
			for _, entity := range scene.entityMap {
				if entity.entityType == uint32(proto.ProtEntityType_PROT_ENTITY_TYPE_MONSTER) {
					monsterEntityCount++
				}
			}
			if monsterEntityCount < 30 {
				monsterEntityId := t.createMonster(scene)

				// PacketSceneEntityAppearNotify
				sceneEntityAppearNotify := new(proto.SceneEntityAppearNotify)
				sceneEntityAppearNotify.AppearType = proto.VisionType_VISION_TYPE_BORN
				sceneEntityInfo := t.gameManager.PacketSceneEntityInfoMonster(scene, monsterEntityId)
				sceneEntityAppearNotify.EntityList = []*proto.SceneEntityInfo{sceneEntityInfo}
				for _, scenePlayer := range scene.playerMap {
					t.gameManager.SendMsg(proto.ApiSceneEntityAppearNotify, scenePlayer.PlayerID, 0, sceneEntityAppearNotify)
				}
			}
		}
		for _, player := range world.playerMap {
			if world.multiplayer || !world.owner.Pause {
				// 改面板
				team := player.TeamConfig.GetActiveTeam()
				for _, avatarId := range team.AvatarIdList {
					if avatarId == 0 {
						break
					}
					avatar := player.AvatarMap[avatarId]
					fightPropertyConst := constant.GetFightPropertyConst()
					avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK)] = 1000000
					avatar.FightPropMap[uint32(fightPropertyConst.FIGHT_PROP_CRITICAL)] = 1.0
					t.gameManager.UpdateUserAvatarFightProp(player.PlayerID, avatarId)
				}
			}
		}
	}
}

func (t *TickManager) onTick5Second(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		if world.IsBigWorld() {
			for applyUid, _ := range world.owner.CoopApplyMap {
				t.gameManager.UserDealEnterWorld(world.owner, applyUid, true)
			}
		}
		for _, player := range world.playerMap {
			if world.multiplayer {
				// PacketWorldPlayerLocationNotify
				worldPlayerLocationNotify := new(proto.WorldPlayerLocationNotify)
				for _, worldPlayer := range world.playerMap {
					playerWorldLocationInfo := &proto.PlayerWorldLocationInfo{
						SceneId: worldPlayer.SceneId,
						PlayerLoc: &proto.PlayerLocationInfo{
							Uid: worldPlayer.PlayerID,
							Pos: &proto.Vector{
								X: float32(worldPlayer.Pos.X),
								Y: float32(worldPlayer.Pos.Y),
								Z: float32(worldPlayer.Pos.Z),
							},
							Rot: &proto.Vector{
								X: float32(worldPlayer.Rot.X),
								Y: float32(worldPlayer.Rot.Y),
								Z: float32(worldPlayer.Rot.Z),
							},
						},
					}
					worldPlayerLocationNotify.PlayerWorldLocList = append(worldPlayerLocationNotify.PlayerWorldLocList, playerWorldLocationInfo)
				}
				t.gameManager.SendMsg(proto.ApiWorldPlayerLocationNotify, player.PlayerID, 0, worldPlayerLocationNotify)

				// PacketScenePlayerLocationNotify
				scene := world.GetSceneById(player.SceneId)
				scenePlayerLocationNotify := new(proto.ScenePlayerLocationNotify)
				scenePlayerLocationNotify.SceneId = player.SceneId
				for _, scenePlayer := range scene.playerMap {
					playerLocationInfo := &proto.PlayerLocationInfo{
						Uid: scenePlayer.PlayerID,
						Pos: &proto.Vector{
							X: float32(scenePlayer.Pos.X),
							Y: float32(scenePlayer.Pos.Y),
							Z: float32(scenePlayer.Pos.Z),
						},
						Rot: &proto.Vector{
							X: float32(scenePlayer.Rot.X),
							Y: float32(scenePlayer.Rot.Y),
							Z: float32(scenePlayer.Rot.Z),
						},
					}
					scenePlayerLocationNotify.PlayerLocList = append(scenePlayerLocationNotify.PlayerLocList, playerLocationInfo)
				}
				t.gameManager.SendMsg(proto.ApiScenePlayerLocationNotify, player.PlayerID, 0, scenePlayerLocationNotify)
			}
		}
	}
}

func (t *TickManager) onTickSecond(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, player := range world.playerMap {
			// PacketWorldPlayerRTTNotify
			worldPlayerRTTNotify := new(proto.WorldPlayerRTTNotify)
			worldPlayerRTTNotify.PlayerRttList = make([]*proto.PlayerRTTInfo, 0)
			for _, worldPlayer := range world.playerMap {
				playerRTTInfo := &proto.PlayerRTTInfo{Uid: worldPlayer.PlayerID, Rtt: worldPlayer.ClientRTT}
				worldPlayerRTTNotify.PlayerRttList = append(worldPlayerRTTNotify.PlayerRttList, playerRTTInfo)
			}
			t.gameManager.SendMsg(proto.ApiWorldPlayerRTTNotify, player.PlayerID, 0, worldPlayerRTTNotify)
		}
	}
}

func (t *TickManager) onTick100MilliSecond(now int64) {
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, scene := range world.sceneMap {
			scene.AttackHandler(t.gameManager)
		}
	}
}

func (t *TickManager) createMonster(scene *Scene) uint32 {
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	fightPropertyConst := constant.GetFightPropertyConst()
	pos := &model.Vector{
		X: 2747,
		Y: 194,
		Z: -1719,
	}
	fpm := map[uint32]float32{
		uint32(fightPropertyConst.FIGHT_PROP_CUR_HP):            float32(72.91699),
		uint32(fightPropertyConst.FIGHT_PROP_PHYSICAL_SUB_HURT): float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_CUR_DEFENSE):       float32(505.0),
		uint32(fightPropertyConst.FIGHT_PROP_CUR_ATTACK):        float32(45.679916),
		uint32(fightPropertyConst.FIGHT_PROP_ICE_SUB_HURT):      float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_BASE_ATTACK):       float32(45.679916),
		uint32(fightPropertyConst.FIGHT_PROP_MAX_HP):            float32(72.91699),
		uint32(fightPropertyConst.FIGHT_PROP_FIRE_SUB_HURT):     float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_ELEC_SUB_HURT):     float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_WIND_SUB_HURT):     float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_ROCK_SUB_HURT):     float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_GRASS_SUB_HURT):    float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_WATER_SUB_HURT):    float32(0.1),
		uint32(fightPropertyConst.FIGHT_PROP_BASE_HP):           float32(72.91699),
		uint32(fightPropertyConst.FIGHT_PROP_BASE_DEFENSE):      float32(505.0),
	}
	entityId := scene.CreateEntityMonster(entityIdTypeConst.MONSTER, pos, 1, fpm)
	return entityId
}
