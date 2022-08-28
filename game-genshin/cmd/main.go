package main

import (
	"flswld.com/common/config"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/light"
	"flswld.com/logger"
	gdc "game-genshin/config"
	"game-genshin/dao"
	"game-genshin/game"
	"game-genshin/mq"
	"game-genshin/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	filePath := "./application.toml"
	config.InitConfig(filePath)

	logger.InitLogger()
	logger.LOG.Info("game-genshin start")

	gdc.InitGameDataConfig()

	db := dao.NewDao()

	netMsgInput := make(chan *proto.NetMsg, 10000)
	netMsgOutput := make(chan *proto.NetMsg, 10000)

	genshinGatewayConsumer := light.NewRpcConsumer("genshin-gateway")
	rpcManager := rpc.NewRpcManager(genshinGatewayConsumer)
	gameServiceProvider := light.NewRpcProvider(rpcManager)

	messageQueue := mq.NewMessageQueue(netMsgInput, netMsgOutput)
	messageQueue.Start()

	gameManager := game.NewGameManager(db, rpcManager, netMsgInput, netMsgOutput)
	gameManager.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logger.LOG.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.LOG.Info("game-genshin exit")
			db.CloseDao()
			gameServiceProvider.CloseRpcProvider()
			genshinGatewayConsumer.CloseRpcConsumer()
			messageQueue.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
