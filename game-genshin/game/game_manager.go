package game

import (
	"flswld.com/common/utils/alg"
	"flswld.com/gate-genshin-api/api"
	"game-genshin/dao"
	"time"
)

type GameManager struct {
	dao          *dao.Dao
	netMsgInput  chan *api.NetMsg
	netMsgOutput chan *api.NetMsg
	snowflake    *alg.SnowflakeWorker
	// 本地事件队列管理器
	localEventManager *LocalEventManager
	// 接口路由管理器
	routeManager *RouteManager
	// 用户管理器
	userManager *UserManager
	// 世界管理器
	worldManager *WorldManager
	// 游戏服务器tick
	tickManager *TickManager
}

func NewGameManager(dao *dao.Dao, netMsgInput chan *api.NetMsg, netMsgOutput chan *api.NetMsg) (r *GameManager) {
	r = new(GameManager)
	r.dao = dao
	r.netMsgInput = netMsgInput
	r.netMsgOutput = netMsgOutput
	r.snowflake = alg.NewSnowflakeWorker(1)
	r.localEventManager = NewLocalEventManager(r)
	r.routeManager = NewRouteManager(r)
	r.userManager = NewUserManager(dao, r.localEventManager.localEventChan)
	r.worldManager = NewWorldManager(r.snowflake)
	r.tickManager = NewTickManager(r)
	return r
}

func (g *GameManager) Start() {
	g.routeManager.InitRoute()
	g.userManager.StartAutoSaveUser()
	go func() {
		for {
			select {
			case netMsg := <-g.netMsgInput:
				// 接收客户端消息
				g.routeManager.RouteHandle(netMsg)
			case <-g.tickManager.ticker.C:
				// 游戏服务器定时帧
				g.tickManager.OnGameServerTick()
			case localEvent := <-g.localEventManager.localEventChan:
				// 处理本地事件
				g.localEventManager.LocalEventHandle(localEvent)
			}
		}
	}()
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
