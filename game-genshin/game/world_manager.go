package game

import (
	"flswld.com/common/utils/alg"
	"game-genshin/constant"
	"game-genshin/model"
)

type WorldManager struct {
	worldMap  map[uint32]*World
	snowflake *alg.SnowflakeWorker
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

func (w *WorldManager) CreateWorld(owner *model.Player) *World {
	worldId := uint32(w.snowflake.GenId())
	world := &World{
		id:              worldId,
		owner:           owner,
		playerMap:       make(map[uint32]*model.Player),
		sceneMap:        make(map[uint32]*Scene),
		entityIdCounter: 0,
		peerIdCounter:   0,
		worldLevel:      0,
		multiplayer:     false,
		mpLevelEntityId: 0,
	}
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	world.mpLevelEntityId = world.GetNextWorldEntityId(entityIdTypeConst.MPLEVEL)
	w.worldMap[worldId] = world
	return world
}

type World struct {
	id              uint32
	owner           *model.Player
	playerMap       map[uint32]*model.Player
	sceneMap        map[uint32]*Scene
	entityIdCounter uint32
	peerIdCounter   uint32
	worldLevel      uint8
	multiplayer     bool
	mpLevelEntityId uint32
}

func (w *World) GetNextWorldEntityId(entityType uint16) uint32 {
	w.entityIdCounter++
	ret := (uint32(entityType) << 24) + w.entityIdCounter
	return ret
}

func (w *World) GetNextWorldPeerId() uint32 {
	w.peerIdCounter++
	return w.peerIdCounter
}

func (w *World) AddPlayer(player *model.Player, sceneId uint32) {
	w.playerMap[player.PlayerID] = player
	scene := w.GetSceneById(sceneId)
	scene.AddPlayer(player)
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	player.TeamConfig.TeamEntityId = w.GetNextWorldEntityId(entityIdTypeConst.TEAM)
}

func (w *World) RemovePlayer(player *model.Player) {
	scene := w.sceneMap[player.SceneId]
	scene.RemovePlayer(player)
	delete(w.playerMap, player.PlayerID)
}

func (w *World) CreateScene(sceneId uint32) *Scene {
	scene := &Scene{
		id:        sceneId,
		world:     w,
		playerMap: make(map[uint32]*model.Player),
		entityMap: make(map[uint32]*Entity),
		time:      0,
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

type Scene struct {
	id        uint32
	world     *World
	playerMap map[uint32]*model.Player
	entityMap map[uint32]*Entity
	time      int64
}

type Entity struct {
	id    uint32
	scene *Scene
}

func (s *Scene) AddPlayer(player *model.Player) {
	s.playerMap[player.PlayerID] = player
}

func (s *Scene) RemovePlayer(player *model.Player) {
	delete(s.playerMap, player.PlayerID)
}

func (s *Scene) CreateEntity(entityType uint16) *Entity {
	entity := &Entity{
		id:    s.world.GetNextWorldEntityId(entityType),
		scene: s,
	}
	s.entityMap[0] = entity
	return entity
}
