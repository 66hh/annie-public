package config

import (
	"encoding/json"
	"io/ioutil"
)

type GameDepot struct {
	PlayerAbilities map[string]*AvatarConfig
}

func (g *GameDataConfig) loadGameDepot() {
	g.GameDepot = new(GameDepot)
	playerElementsFilePath := g.binPrefix + "AbilityGroup/AbilityGroup_Other_PlayerElementAbility.json"
	playerElementsFile, err := ioutil.ReadFile(playerElementsFilePath)
	if err != nil {
		g.log.Error("open file error: %v", err)
		return
	}
	playerAbilities := make(map[string]*AvatarConfig)
	err = json.Unmarshal(playerElementsFile, &playerAbilities)
	if err != nil {
		g.log.Error("parse file error: %v", err)
		return
	}
	g.GameDepot.PlayerAbilities = playerAbilities
	g.log.Info("load %v PlayerAbilities", len(g.GameDepot.PlayerAbilities))
}
