package main

import (
	"flswld.com/common/config"
	"flswld.com/gate-genshin-api/api"
	_ "flswld.com/gate-genshin-api/api/proto"
	"flswld.com/light"
	"flswld.com/logger"
	"game-genshin/dao"
	"game-genshin/game"
	"game-genshin/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	filePath := "./application.toml"
	conf := config.NewConfig(filePath)

	log := logger.NewLogger(conf)
	log.Info("game-genshin start")

	db := dao.NewDao(conf, log)

	netMsgInput := make(chan *api.NetMsg, 1000)
	netMsgOutput := make(chan *api.NetMsg, 1000)

	gameManager := game.NewGameManager(log, conf, db, netMsgInput, netMsgOutput)
	gameManager.Start()

	genshinGatewayConsumer := light.NewRpcConsumer(conf, log, "genshin-gateway")
	rpcManager := rpc.NewRpcManager(genshinGatewayConsumer, netMsgInput, netMsgOutput)
	gameServiceProvider := light.NewRpcProvider(conf, log, rpcManager)
	rpcManager.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("game-genshin exit")
			db.CloseDao()
			gameServiceProvider.CloseRpcProvider()
			genshinGatewayConsumer.CloseRpcConsumer()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
