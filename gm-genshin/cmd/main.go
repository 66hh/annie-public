package main

import (
	"flswld.com/common/config"
	"flswld.com/light"
	"flswld.com/logger"
	"gm-genshin/controller"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	filePath := "./application.toml"
	conf := config.NewConfig(filePath)

	log := logger.NewLogger(conf)
	log.Info("gm genshin start")

	httpProvider := light.NewHttpProvider(conf, log)

	// 认证服务
	rpcWaterAuthConsumer := light.NewRpcConsumer(conf, log, "water-auth")

	rpcGenshinGatewayConsumer := light.NewRpcConsumer(conf, log, "genshin-gateway")

	_ = controller.NewController(conf, log, rpcWaterAuthConsumer, rpcGenshinGatewayConsumer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcWaterAuthConsumer.CloseRpcConsumer()
			rpcGenshinGatewayConsumer.CloseRpcConsumer()
			httpProvider.CloseHttpProvider()
			log.Info("gm genshin exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
