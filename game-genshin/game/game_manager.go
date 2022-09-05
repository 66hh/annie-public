package game

import (
	"flswld.com/common/utils/alg"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"game-genshin/dao"
	"game-genshin/rpc"
	pb "google.golang.org/protobuf/proto"
)

type GameManager struct {
	dao          *dao.Dao
	rpcManager   *rpc.RpcManager
	netMsgInput  chan *proto.NetMsg
	netMsgOutput chan *proto.NetMsg
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

func NewGameManager(dao *dao.Dao, rpcManager *rpc.RpcManager, netMsgInput chan *proto.NetMsg, netMsgOutput chan *proto.NetMsg) (r *GameManager) {
	r = new(GameManager)
	r.dao = dao
	r.rpcManager = rpcManager
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
			case netMsg := <-g.netMsgOutput:
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
func (g *GameManager) SendMsg(apiId uint16, userId uint32, clientSeq uint32, payloadMsg pb.Message) {
	netMsg := new(proto.NetMsg)
	netMsg.UserId = userId
	netMsg.EventId = proto.NormalMsg
	netMsg.ApiId = apiId
	netMsg.ClientSeq = clientSeq
	// 在这里直接序列化成二进制数据 防止发送的消息内包含各种游戏数据指针 而造成并发读写的问题
	payloadMessageData, err := pb.Marshal(payloadMsg)
	if err != nil {
		logger.LOG.Error("parse payload msg to bin error: %v", err)
		return
	}
	netMsg.PayloadMessageData = payloadMessageData
	g.netMsgInput <- netMsg
}

func (g *GameManager) KickPlayer(userId uint32) {
	g.rpcManager.SendKickPlayerToGenshinGateway(userId)
}
