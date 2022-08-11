package constant

import "flswld.com/common/utils/endec"

type ElementTypeValue struct {
	Value           uint16
	CurrEnergyProp  uint16
	MaxEnergyProp   uint16
	TeamResonanceId uint16
	ConfigName      string
	ConfigHash      int32
}

type ElementType struct {
	None       *ElementTypeValue
	Fire       *ElementTypeValue
	Water      *ElementTypeValue
	Grass      *ElementTypeValue
	Electric   *ElementTypeValue
	Ice        *ElementTypeValue
	Frozen     *ElementTypeValue
	Wind       *ElementTypeValue
	Rock       *ElementTypeValue
	AntiFire   *ElementTypeValue
	Default    *ElementTypeValue
	STRING_MAP map[string]*ElementTypeValue
	VALUE_MAP  map[uint16]*ElementTypeValue
}

func GetElementTypeConst() (r *ElementType) {
	r = new(ElementType)
	fightProperty := GetFightPropertyConst()

	r.None = &ElementTypeValue{
		0,
		fightProperty.FIGHT_PROP_CUR_FIRE_ENERGY,
		fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY,
		0,
		"",
		endec.GenshinAbilityHashCode(""),
	}
	r.Fire = &ElementTypeValue{
		1,
		fightProperty.FIGHT_PROP_CUR_FIRE_ENERGY,
		fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY,
		10101,
		"TeamResonance_Fire_Lv2",
		endec.GenshinAbilityHashCode("TeamResonance_Fire_Lv2"),
	}
	r.Water = &ElementTypeValue{
		2,
		fightProperty.FIGHT_PROP_CUR_WATER_ENERGY,
		fightProperty.FIGHT_PROP_MAX_WATER_ENERGY,
		10201,
		"TeamResonance_Water_Lv2",
		endec.GenshinAbilityHashCode("TeamResonance_Water_Lv2"),
	}
	r.Grass = &ElementTypeValue{
		3,
		fightProperty.FIGHT_PROP_CUR_GRASS_ENERGY,
		fightProperty.FIGHT_PROP_MAX_GRASS_ENERGY,
		0,
		"",
		endec.GenshinAbilityHashCode(""),
	}
	r.Electric = &ElementTypeValue{
		4,
		fightProperty.FIGHT_PROP_CUR_ELEC_ENERGY,
		fightProperty.FIGHT_PROP_MAX_ELEC_ENERGY,
		10401,
		"TeamResonance_Electric_Lv2",
		endec.GenshinAbilityHashCode("TeamResonance_Electric_Lv2"),
	}
	r.Ice = &ElementTypeValue{
		5,
		fightProperty.FIGHT_PROP_CUR_ICE_ENERGY,
		fightProperty.FIGHT_PROP_MAX_ICE_ENERGY,
		10601,
		"TeamResonance_Ice_Lv2",
		endec.GenshinAbilityHashCode("TeamResonance_Ice_Lv2"),
	}
	r.Frozen = &ElementTypeValue{
		6,
		fightProperty.FIGHT_PROP_CUR_ICE_ENERGY,
		fightProperty.FIGHT_PROP_MAX_ICE_ENERGY,
		0,
		"",
		endec.GenshinAbilityHashCode(""),
	}
	r.Wind = &ElementTypeValue{
		7,
		fightProperty.FIGHT_PROP_CUR_WIND_ENERGY,
		fightProperty.FIGHT_PROP_MAX_WIND_ENERGY,
		10301,
		"TeamResonance_Wind_Lv2",
		endec.GenshinAbilityHashCode("TeamResonance_Wind_Lv2"),
	}
	r.Rock = &ElementTypeValue{
		8,
		fightProperty.FIGHT_PROP_CUR_ROCK_ENERGY,
		fightProperty.FIGHT_PROP_MAX_ROCK_ENERGY,
		10701,
		"TeamResonance_Rock_Lv2",
		endec.GenshinAbilityHashCode("TeamResonance_Rock_Lv2"),
	}
	r.AntiFire = &ElementTypeValue{
		9,
		fightProperty.FIGHT_PROP_CUR_FIRE_ENERGY,
		fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY,
		0,
		"",
		endec.GenshinAbilityHashCode(""),
	}
	r.Default = &ElementTypeValue{
		255,
		fightProperty.FIGHT_PROP_CUR_FIRE_ENERGY,
		fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY,
		10801,
		"TeamResonance_AllDifferent",
		endec.GenshinAbilityHashCode("TeamResonance_AllDifferent"),
	}

	r.STRING_MAP = make(map[string]*ElementTypeValue)

	r.STRING_MAP["None"] = r.None
	r.STRING_MAP["Fire"] = r.Fire
	r.STRING_MAP["Water"] = r.Water
	r.STRING_MAP["Grass"] = r.Grass
	r.STRING_MAP["Electric"] = r.Electric
	r.STRING_MAP["Ice"] = r.Ice
	r.STRING_MAP["Frozen"] = r.Frozen
	r.STRING_MAP["Wind"] = r.Wind
	r.STRING_MAP["Rock"] = r.Rock
	r.STRING_MAP["AntiFire"] = r.AntiFire
	r.STRING_MAP["Default"] = r.Default

	r.VALUE_MAP = make(map[uint16]*ElementTypeValue)

	r.VALUE_MAP[0] = r.None
	r.VALUE_MAP[1] = r.Fire
	r.VALUE_MAP[2] = r.Water
	r.VALUE_MAP[3] = r.Grass
	r.VALUE_MAP[4] = r.Electric
	r.VALUE_MAP[5] = r.Ice
	r.VALUE_MAP[6] = r.Frozen
	r.VALUE_MAP[7] = r.Wind
	r.VALUE_MAP[8] = r.Rock
	r.VALUE_MAP[9] = r.AntiFire
	r.VALUE_MAP[255] = r.Default

	return r
}
