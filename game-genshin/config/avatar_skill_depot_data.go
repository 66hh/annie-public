package config

import (
	"encoding/json"
	"flswld.com/common/utils/endec"
	"flswld.com/logger"
	"game-genshin/constant"
	"io/ioutil"
)

type InherentProudSkillOpens struct {
	ProudSkillGroupId      int32 `json:"proudSkillGroupId"`
	NeedAvatarPromoteLevel int32 `json:"needAvatarPromoteLevel"`
}

type AvatarSkillDepotData struct {
	Id              int32 `json:"id"`
	EnergySkill     int32 `json:"energySkill"`
	AttackModeSkill int32 `json:"attackModeSkill"`

	Skills                  []int32                    `json:"skills"`
	SubSkills               []int32                    `json:"subSkills"`
	ExtraAbilities          []string                   `json:"extraAbilities"`
	Talents                 []int32                    `json:"talents"`
	InherentProudSkillOpens []*InherentProudSkillOpens `json:"inherentProudSkillOpens"`
	TalentStarName          string                     `json:"talentStarName"`
	SkillDepotAbilityGroup  string                     `json:"skillDepotAbilityGroup"`

	// 计算属性
	EnergySkillData *AvatarSkillData           `json:"-"`
	ElementType     *constant.ElementTypeValue `json:"-"`
	Abilities       []int32                    `json:"-"`
}

func (g *GameDataConfig) loadAvatarSkillDepotData() {
	g.AvatarSkillDepotDataMap = make(map[int32]*AvatarSkillDepotData)
	fileNameList := []string{"AvatarSkillDepotExcelConfigData.json"}
	for _, fileName := range fileNameList {
		fileData, err := ioutil.ReadFile(g.excelBinPrefix + fileName)
		if err != nil {
			logger.LOG.Error("open file error: %v", err)
			continue
		}
		list := make([]map[string]any, 0)
		err = json.Unmarshal(fileData, &list)
		if err != nil {
			logger.LOG.Error("parse file error: %v", err)
			continue
		}
		for _, v := range list {
			i, err := json.Marshal(v)
			if err != nil {
				logger.LOG.Error("parse file error: %v", err)
				continue
			}
			avatarSkillDepotData := new(AvatarSkillDepotData)
			err = json.Unmarshal(i, avatarSkillDepotData)
			if err != nil {
				logger.LOG.Error("parse file error: %v", err)
				continue
			}
			g.AvatarSkillDepotDataMap[avatarSkillDepotData.Id] = avatarSkillDepotData
		}
	}
	logger.LOG.Info("load %v AvatarSkillDepotData", len(g.AvatarSkillDepotDataMap))
	elementTypeConst := constant.GetElementTypeConst()
	for _, v := range g.AvatarSkillDepotDataMap {
		// set energy skill data
		v.EnergySkillData = g.AvatarSkillDataMap[v.EnergySkill]
		if v.EnergySkillData != nil {
			v.ElementType = v.EnergySkillData.CostElemTypeX
		} else {
			v.ElementType = elementTypeConst.None
		}
		// set embryo abilities if player skill depot
		if v.SkillDepotAbilityGroup != "" {
			config := g.GameDepot.PlayerAbilities[v.SkillDepotAbilityGroup]
			if config != nil {
				for _, targetAbility := range config.TargetAbilities {
					v.Abilities = append(v.Abilities, endec.GenshinAbilityHashCode(targetAbility.AbilityName))
				}
			}
		}
	}
}
