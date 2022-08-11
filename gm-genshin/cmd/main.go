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
	config.InitConfig(filePath)

	logger.InitLogger()
	logger.LOG.Info("gm genshin start")

	httpProvider := light.NewHttpProvider()

	// 认证服务
	rpcWaterAuthConsumer := light.NewRpcConsumer("water-auth")

	rpcGenshinGatewayConsumer := light.NewRpcConsumer("genshin-gateway")

	_ = controller.NewController(rpcWaterAuthConsumer, rpcGenshinGatewayConsumer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logger.LOG.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcWaterAuthConsumer.CloseRpcConsumer()
			rpcGenshinGatewayConsumer.CloseRpcConsumer()
			httpProvider.CloseHttpProvider()
			logger.LOG.Info("gm genshin exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
