package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type AvatarConfigAbility struct {
	AbilityName string `json:"abilityName"`
}

type AvatarConfig struct {
	Abilities []*AvatarConfigAbility `json:"abilities"`
}

type AbilityEmbryoEntry struct {
	Name      string
	Abilities []string
}

func (g *GameDataConfig) loadAbilityEmbryos() {
	dirPath := g.binPrefix + "Avatar"
	fileList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		g.log.Error("open dir error: %v", err)
		return
	}
	embryoList := make([]*AbilityEmbryoEntry, 0)
	for _, file := range fileList {
		fileName := file.Name()
		if !strings.Contains(fileName, "ConfigAvatar_") {
			continue
		}
		startIndex := strings.Index(fileName, "ConfigAvatar_")
		endIndex := strings.Index(fileName, ".json")
		if startIndex == -1 || endIndex == -1 || startIndex+13 > endIndex {
			g.log.Error("file name format error: %v", fileName)
			continue
		}
		avatarName := fileName[startIndex+13 : endIndex]
		fileData, err := ioutil.ReadFile(dirPath + "/" + fileName)
		if err != nil {
			g.log.Error("open file error: %v", err)
			continue
		}
		avatarConfig := new(AvatarConfig)
		err = json.Unmarshal(fileData, avatarConfig)
		if err != nil {
			g.log.Error("parse file error: %v", err)
			continue
		}
		if len(avatarConfig.Abilities) == 0 {
			continue
		}
		abilityEmbryoEntry := new(AbilityEmbryoEntry)
		abilityEmbryoEntry.Name = avatarName
		for _, v := range avatarConfig.Abilities {
			abilityEmbryoEntry.Abilities = append(abilityEmbryoEntry.Abilities, v.AbilityName)
		}
		embryoList = append(embryoList, abilityEmbryoEntry)
	}
	if len(embryoList) == 0 {
		g.log.Error("no embryo load")
	}
	g.AbilityEmbryos = make(map[string]*AbilityEmbryoEntry)
	for _, v := range embryoList {
		g.AbilityEmbryos[v.Name] = v
	}
	g.log.Info("load %v AbilityEmbryos", len(g.AbilityEmbryos))
}
