package constant

type ElementTypeValue struct {
	value           uint16
	energyProperty  uint16
	teamResonanceId uint16
	configName      string
}

type ElementType struct {
	None     *ElementTypeValue
	Fire     *ElementTypeValue
	Water    *ElementTypeValue
	Grass    *ElementTypeValue
	Electric *ElementTypeValue
	Ice      *ElementTypeValue
	Frozen   *ElementTypeValue
	Wind     *ElementTypeValue
	Rock     *ElementTypeValue
	AntiFire *ElementTypeValue
	Default  *ElementTypeValue
}

func GetElementTypeConst() (r *ElementType) {
	r = new(ElementType)
	fightProperty := GetFightPropertyConst()
	r.None = &ElementTypeValue{0, fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY, 0, ""}
	r.Fire = &ElementTypeValue{1, fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY, 10101, "TeamResonance_Fire_Lv2"}
	r.Water = &ElementTypeValue{2, fightProperty.FIGHT_PROP_MAX_WATER_ENERGY, 10201, "TeamResonance_Water_Lv2"}
	r.Grass = &ElementTypeValue{3, fightProperty.FIGHT_PROP_MAX_GRASS_ENERGY, 0, ""}
	r.Electric = &ElementTypeValue{4, fightProperty.FIGHT_PROP_MAX_ELEC_ENERGY, 10401, "TeamResonance_Electric_Lv2"}
	r.Ice = &ElementTypeValue{5, fightProperty.FIGHT_PROP_MAX_ICE_ENERGY, 10601, "TeamResonance_Ice_Lv2"}
	r.Frozen = &ElementTypeValue{6, fightProperty.FIGHT_PROP_MAX_ICE_ENERGY, 0, ""}
	r.Wind = &ElementTypeValue{7, fightProperty.FIGHT_PROP_MAX_WIND_ENERGY, 10301, "TeamResonance_Wind_Lv2"}
	r.Rock = &ElementTypeValue{8, fightProperty.FIGHT_PROP_MAX_ROCK_ENERGY, 10701, "TeamResonance_Rock_Lv2"}
	r.AntiFire = &ElementTypeValue{9, fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY, 0, ""}
	r.Default = &ElementTypeValue{255, fightProperty.FIGHT_PROP_MAX_FIRE_ENERGY, 10801, "TeamResonance_AllDifferent"}
	return r
}
