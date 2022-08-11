package constant

import "flswld.com/common/utils/endec"

type GameConstant struct {
	DEFAULT_ABILITY_STRINGS []string
	DEFAULT_ABILITY_HASHES  []int32
	DEFAULT_ABILITY_NAME    int32
}

func GetGameConstant() (r *GameConstant) {
	r = new(GameConstant)
	r.DEFAULT_ABILITY_STRINGS = []string{
		"Avatar_DefaultAbility_VisionReplaceDieInvincible",
		"Avatar_DefaultAbility_AvartarInShaderChange",
		"Avatar_SprintBS_Invincible",
		"Avatar_Freeze_Duration_Reducer",
		"Avatar_Attack_ReviveEnergy",
		"Avatar_Component_Initializer",
		"Avatar_FallAnthem_Achievement_Listener",
	}
	r.DEFAULT_ABILITY_HASHES = make([]int32, 0)
	for _, v := range r.DEFAULT_ABILITY_STRINGS {
		r.DEFAULT_ABILITY_HASHES = append(r.DEFAULT_ABILITY_HASHES, endec.GenshinAbilityHashCode(v))
	}
	r.DEFAULT_ABILITY_NAME = endec.GenshinAbilityHashCode("Default")
	return r
}
