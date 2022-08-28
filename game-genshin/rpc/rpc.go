package rpc

import (
	"flswld.com/gate-genshin-api/gm"
	"flswld.com/light"
)

type RpcManager struct {
	genshinGatewayConsumer *light.Consumer
}

func NewRpcManager(genshinGatewayConsumer *light.Consumer) (r *RpcManager) {
	r = new(RpcManager)
	r.genshinGatewayConsumer = genshinGatewayConsumer
	return r
}

func (r *RpcManager) SendKickPlayerToGenshinGateway(userId uint32) {
	info := new(gm.KickPlayerInfo)
	info.UserId = userId
	// 客户端提示信息为服务器断开连接
	info.Reason = uint32(5)
	var result bool
	ok := r.genshinGatewayConsumer.CallFunction("RpcManager", "KickPlayer", &info, &result)
	if ok == true && result == true {
		return
	}
	return
}
