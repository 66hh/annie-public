package constant

type EquipType struct {
	EQUIP_NONE     uint16
	EQUIP_BRACER   uint16
	EQUIP_NECKLACE uint16
	EQUIP_SHOES    uint16
	EQUIP_RING     uint16
	EQUIP_DRESS    uint16
	EQUIP_WEAPON   uint16
	STRING_MAP     map[string]uint16
}

func GetEquipTypeConst() (r *EquipType) {
	r = new(EquipType)

	r.EQUIP_NONE = 0
	r.EQUIP_BRACER = 1
	r.EQUIP_NECKLACE = 2
	r.EQUIP_SHOES = 3
	r.EQUIP_RING = 4
	r.EQUIP_DRESS = 5
	r.EQUIP_WEAPON = 6

	r.STRING_MAP = make(map[string]uint16)

	r.STRING_MAP["EQUIP_NONE"] = 0
	r.STRING_MAP["EQUIP_BRACER"] = 1
	r.STRING_MAP["EQUIP_NECKLACE"] = 2
	r.STRING_MAP["EQUIP_SHOES"] = 3
	r.STRING_MAP["EQUIP_RING"] = 4
	r.STRING_MAP["EQUIP_DRESS"] = 5
	r.STRING_MAP["EQUIP_WEAPON"] = 6

	return r
}
