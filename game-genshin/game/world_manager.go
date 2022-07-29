package game

type WorldManager struct {
	entityID uint32
}

func NewWorldManager() (r *WorldManager) {
	r = new(WorldManager)
	return r
}

func (w *WorldManager) GetNextWorldEntityID(entityType uint16) uint32 {
	w.entityID++
	ret := (uint32(entityType) << 24) + w.entityID
	return ret
}
