package main

import (
	"flswld.com/common/config"
	_ "flswld.com/gate-genshin-api/api/proto"
	"flswld.com/light"
	"flswld.com/logger"
	"gate-genshin/controller"
	"gate-genshin/dao"
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

	// 用户服务
	rpcUserConsumer := light.NewRpcConsumer("annie-user-app")

	db := dao.NewDao()

	_ = controller.NewController(db, rpcUserConsumer)

	kcpEventInput := make(chan *net.KcpEvent)
	kcpEventOutput := make(chan *net.KcpEvent)
	kcpMsgInput := make(chan *net.KcpMsg, 1000)
	kcpMsgOutput := make(chan *net.KcpMsg, 1000)
	protoMsgInput := make(chan *net.ProtoMsg, 1000)
	protoMsgOutput := make(chan *net.ProtoMsg, 1000)

	connectManager := net.NewKcpConnectManager(kcpEventInput, kcpEventOutput, kcpMsgInput, kcpMsgOutput)
	protoEnDecode := net.NewProtoEnDecode(kcpMsgInput, kcpMsgOutput, protoMsgInput, protoMsgOutput)
	connectManager.Start()
	protoEnDecode.Start()

	gameServiceConsumer := light.NewRpcConsumer("game-genshin-app")
	rpcManager := rpc.NewRpcManager(db, gameServiceConsumer, protoMsgInput, protoMsgOutput, kcpEventInput, kcpEventOutput)
	rpcMsgProvider := light.NewRpcProvider(rpcManager)
	rpcManager.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logger.LOG.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.LOG.Info("gate genshin exit")
			db.CloseDao()
			rpcUserConsumer.CloseRpcConsumer()
			rpcMsgProvider.CloseRpcProvider()
			gameServiceConsumer.CloseRpcConsumer()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
