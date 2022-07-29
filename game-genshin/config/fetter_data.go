package config

import (
	"encoding/json"
	"io/ioutil"
)

type FetterData struct {
	AvatarId uint32 `json:"AvatarId"`
	FetterId uint32 `json:"FetterId"`
}

func (g *GameDataConfig) loadFetterData() {
	g.FetterDataMap = make(map[uint32]*FetterData)
	fileNameList := []string{"FetterInfoExcelConfigData.json", "FettersExcelConfigData.json", "FetterStoryExcelConfigData.json", "PhotographExpressionExcelConfigData.json", "PhotographPosenameExcelConfigData.json"}
	for _, fileName := range fileNameList {
		fileData, err := ioutil.ReadFile(g.excelBinPrefix + fileName)
		if err != nil {
			g.log.Error("open file error: %v", err)
			continue
		}
		list := make([]map[string]any, 0)
		err = json.Unmarshal(fileData, &list)
		if err != nil {
			g.log.Error("parse file error: %v", err)
			continue
		}
		for _, v := range list {
			i, err := json.Marshal(v)
			if err != nil {
				g.log.Error("parse file error: %v", err)
				continue
			}
			fetterData := new(FetterData)
			err = json.Unmarshal(i, fetterData)
			if err != nil {
				g.log.Error("parse file error: %v", err)
				continue
			}
			g.FetterDataMap[fetterData.FetterId] = fetterData
		}
	}
	g.log.Info("load %v FetterData", len(g.FetterDataMap))
	g.AvatarFetterDataMap = make(map[uint32][]uint32)
	for _, v := range g.FetterDataMap {
		avatarFetterIdList, exist := g.AvatarFetterDataMap[v.AvatarId]
		if !exist {
			avatarFetterIdList = make([]uint32, 0)
		}
		avatarFetterIdList = append(avatarFetterIdList, v.FetterId)
		g.AvatarFetterDataMap[v.AvatarId] = avatarFetterIdList
	}
}
