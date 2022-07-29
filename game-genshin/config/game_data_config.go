package config

import (
	"flswld.com/common/config"
	"flswld.com/logger"
	"os"
)

type GameDataConfig struct {
	log            *logger.Logger
	conf           *config.Config
	binPrefix      string
	excelBinPrefix string
	// 配置表
	// BinOutput
	// 技能列表
	AbilityEmbryos    map[string]*AbilityEmbryoEntry
	OpenConfigEntries map[string]*OpenConfigEntry
	// ExcelBinOutput
	FetterDataMap       map[uint32]*FetterData
	AvatarFetterDataMap map[uint32][]uint32
	// 资源
	// 场景传送点
	ScenePointEntries map[string]*ScenePointEntry
	ScenePointIdList  []uint32
}

func NewGameDataConfig(log *logger.Logger, conf *config.Config) (r *GameDataConfig) {
	r = new(GameDataConfig)
	r.log = log
	r.conf = conf
	r.binPrefix = ""
	r.excelBinPrefix = ""
	return r
}

func (g *GameDataConfig) load() {
	// 技能列表
	g.loadAbilityEmbryos()
	g.loadOpenConfig()
	// 资源
	g.loadFetterData()
	// 场景传送点
	g.loadScenePoints()
}

func (g *GameDataConfig) LoadAll() {
	resourcePath := g.conf.Genshin.ResourcePath
	dirInfo, err := os.Stat(resourcePath)
	if err != nil || !dirInfo.IsDir() {
		g.log.Error("open game data config dir error: %v", err)
		return
	}
	g.binPrefix = resourcePath + "/BinOutput"
	g.excelBinPrefix = resourcePath + "/ExcelBinOutput"
	dirInfo, err = os.Stat(g.binPrefix)
	if err != nil || !dirInfo.IsDir() {
		g.log.Error("open game data bin output config dir error: %v", err)
		return
	}
	dirInfo, err = os.Stat(g.excelBinPrefix)
	if err != nil || !dirInfo.IsDir() {
		g.log.Error("open game data excel bin output config dir error: %v", err)
		return
	}
	g.binPrefix += "/"
	g.excelBinPrefix += "/"
	g.load()
}
