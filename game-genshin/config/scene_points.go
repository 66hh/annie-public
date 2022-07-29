package config

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

type ScenePointEntry struct {
	Name      string     `json:"name"`
	PointData *PointData `json:"pointData"`
}

type ScenePointConfig struct {
	Points map[string]*PointData `json:"points"`
}

type PointData struct {
	Id                int32
	DollarType        string    `json:"$type"`
	TranPos           *Position `json:"tranPos"`
	DungeonIds        []int32   `json:"dungeonIds"`
	DungeonRandomList []int32   `json:"dungeonRandomList"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (g *GameDataConfig) loadScenePoints() {
	g.ScenePointEntries = make(map[string]*ScenePointEntry)
	g.ScenePointIdList = make([]uint32, 0)
	dirPath := g.binPrefix + "Scene/Point"
	fileList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		g.log.Error("open dir error: %v", err)
		return
	}
	for _, file := range fileList {
		fileName := file.Name()
		if !strings.Contains(fileName, "scene") {
			continue
		}
		startIndex := strings.Index(fileName, "scene")
		endIndex := strings.Index(fileName, "_point.json")
		if startIndex == -1 || endIndex == -1 || startIndex+5 > endIndex {
			g.log.Error("file name format error: %v", fileName)
			continue
		}
		sceneId := fileName[startIndex+5 : endIndex]
		fileData, err := ioutil.ReadFile(dirPath + "/" + fileName)
		if err != nil {
			g.log.Error("open file error: %v", err)
			continue
		}
		scenePointConfig := new(ScenePointConfig)
		err = json.Unmarshal(fileData, scenePointConfig)
		if err != nil {
			g.log.Error("parse file error: %v", err)
			continue
		}
		if len(scenePointConfig.Points) == 0 {
			continue
		}
		for k, v := range scenePointConfig.Points {
			sceneIdInt32, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				g.log.Error("parse file error: %v", err)
				continue
			}
			v.Id = int32(sceneIdInt32)
			scenePointEntry := new(ScenePointEntry)
			scenePointEntry.Name = sceneId + "_" + k
			scenePointEntry.PointData = v
			g.ScenePointIdList = append(g.ScenePointIdList, uint32(sceneIdInt32))
			g.ScenePointEntries[scenePointEntry.Name] = scenePointEntry
		}
	}
	g.log.Info("load %v ScenePointEntries", len(g.ScenePointEntries))
}
