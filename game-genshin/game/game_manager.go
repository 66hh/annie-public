package game

import (
	"flswld.com/common/config"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/logger"
	gameDataConfig "game-genshin/config"
	"game-genshin/dao"
	"time"
)

type GameManager struct {
	log          *logger.Logger
	conf         *config.Config
	dao          *dao.Dao
	netMsgInput  chan *api.NetMsg
	netMsgOutput chan *api.NetMsg
	// 配置表
	gameDataConfig *gameDataConfig.GameDataConfig
	// 接口路由管理器
	routeManager *RouteManager
	// 用户管理器
	userManager *UserManager
	// 世界管理器
	worldManager *WorldManager
}

func NewGameManager(log *logger.Logger, conf *config.Config, dao *dao.Dao, netMsgInput chan *api.NetMsg, netMsgOutput chan *api.NetMsg) (r *GameManager) {
	r = new(GameManager)
	r.log = log
	r.conf = conf
	r.dao = dao
	r.netMsgInput = netMsgInput
	r.netMsgOutput = netMsgOutput
	r.gameDataConfig = gameDataConfig.NewGameDataConfig(log, conf)
	r.routeManager = NewRouteManager(log, r)
	r.userManager = NewUserManager(log, dao)
	r.worldManager = NewWorldManager()
	return r
}

func (g *GameManager) Start() {
	g.gameDataConfig.LoadAll()
	g.userManager.StartAutoSaveUser()
	g.routeManager.InitRoute()
	g.routeManager.StartRouteHandle(g.netMsgInput)
}

// 发送消息给客户端
func (g *GameManager) SendMsg(apiId uint16, userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	netMsg := new(api.NetMsg)
	netMsg.UserId = userId
	netMsg.EventId = api.NormalMsg
	netMsg.ApiId = apiId
	netMsg.HeadMessage = headMsg
	netMsg.PayloadMessage = payloadMsg
	g.netMsgOutput <- netMsg
}

func (g *GameManager) getHeadMsg(seq uint32) (headMsg *api.PacketHead) {
	headMsg = new(api.PacketHead)
	headMsg.ClientSequenceId = seq
	headMsg.Timestamp = uint64(time.Now().UnixMilli())
	return headMsg
}
