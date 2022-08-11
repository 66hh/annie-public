package constant

type EntityIdType struct {
	AVATAR  uint16
	MONSTER uint16
	NPC     uint16
	GADGET  uint16
	WEAPON  uint16
	TEAM    uint16
	MPLEVEL uint16
}

func GetEntityIdTypeConst() (r *EntityIdType) {
	r = new(EntityIdType)
	r.AVATAR = 0x01
	r.MONSTER = 0x02
	r.NPC = 0x03
	r.GADGET = 0x04
	r.WEAPON = 0x06
	r.TEAM = 0x09
	r.MPLEVEL = 0x0b
	return r
}
