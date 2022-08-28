package main

import (
	"flswld.com/common/config"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/light"
	"flswld.com/logger"
	"gate-genshin/controller"
	"gate-genshin/dao"
	"gate-genshin/forward"
	"gate-genshin/mq"
	"gate-genshin/net"
	"gate-genshin/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	filePath := "./application.toml"
	config.InitConfig(filePath)

	logger.InitLogger()
	logger.LOG.Info("gate genshin start")

	db := dao.NewDao()

	// 用户服务
	rpcUserConsumer := light.NewRpcConsumer("annie-user-app")

	_ = controller.NewController(db, rpcUserConsumer)

	kcpEventInput := make(chan *net.KcpEvent)
	kcpEventOutput := make(chan *net.KcpEvent)
	protoMsgInput := make(chan *net.ProtoMsg, 10000)
	protoMsgOutput := make(chan *net.ProtoMsg, 10000)
	netMsgInput := make(chan *proto.NetMsg, 10000)
	netMsgOutput := make(chan *proto.NetMsg, 10000)

	connectManager := net.NewKcpConnectManager(protoMsgInput, protoMsgOutput, kcpEventInput, kcpEventOutput)
	connectManager.Start()

	forwardManager := forward.NewForwardManager(db, protoMsgInput, protoMsgOutput, kcpEventInput, kcpEventOutput, netMsgInput, netMsgOutput)
	forwardManager.Start()

	gameServiceConsumer := light.NewRpcConsumer("game-genshin-app")

	rpcManager := rpc.NewRpcManager(forwardManager)
	rpcMsgProvider := light.NewRpcProvider(rpcManager)

	messageQueue := mq.NewMessageQueue(netMsgInput, netMsgOutput)
	messageQueue.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logger.LOG.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.LOG.Info("gate genshin exit")
			messageQueue.Close()
			rpcMsgProvider.CloseRpcProvider()
			gameServiceConsumer.CloseRpcConsumer()
			rpcUserConsumer.CloseRpcConsumer()
			db.CloseDao()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
