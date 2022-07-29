package constant

type LifeState struct {
	LIFE_NONE   uint16
	LIFE_ALIVE  uint16
	LIFE_DEAD   uint16
	LIFE_REVIVE uint16
}

func GetLifeStateConst() (r *LifeState) {
	r = new(LifeState)
	r.LIFE_NONE = 0
	r.LIFE_ALIVE = 1
	r.LIFE_DEAD = 2
	r.LIFE_REVIVE = 3
	return r
}
