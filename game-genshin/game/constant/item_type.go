package constant

type ItemType struct {
	ITEM_NONE      uint16
	ITEM_VIRTUAL   uint16
	ITEM_MATERIAL  uint16
	ITEM_RELIQUARY uint16
	ITEM_WEAPON    uint16
	ITEM_DISPLAY   uint16
	ITEM_FURNITURE uint16
	STRING_MAP     map[string]uint16
}

func GetItemTypeConst() (r *ItemType) {
	r = new(ItemType)

	r.ITEM_NONE = 0
	r.ITEM_VIRTUAL = 1
	r.ITEM_MATERIAL = 2
	r.ITEM_RELIQUARY = 3
	r.ITEM_WEAPON = 4
	r.ITEM_DISPLAY = 5
	r.ITEM_FURNITURE = 6

	r.STRING_MAP = make(map[string]uint16)

	r.STRING_MAP["ITEM_NONE"] = 0
	r.STRING_MAP["ITEM_VIRTUAL"] = 1
	r.STRING_MAP["ITEM_MATERIAL"] = 2
	r.STRING_MAP["ITEM_RELIQUARY"] = 3
	r.STRING_MAP["ITEM_WEAPON"] = 4
	r.STRING_MAP["ITEM_DISPLAY"] = 5
	r.STRING_MAP["ITEM_FURNITURE"] = 6

	return r
}
