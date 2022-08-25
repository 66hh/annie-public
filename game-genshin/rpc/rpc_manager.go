package rpc

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/gate-genshin-api/gm"
	"flswld.com/light"
)

type RpcManager struct {
	genshinGatewayConsumer *light.Consumer
	netMsgInput            chan *api.NetMsg
	netMsgOutput           chan *api.NetMsg
}

func NewRpcManager(genshinGatewayConsumer *light.Consumer, netMsgInput chan *api.NetMsg, netMsgOutput chan *api.NetMsg) (r *RpcManager) {
	r = new(RpcManager)
	r.genshinGatewayConsumer = genshinGatewayConsumer
	r.netMsgInput = netMsgInput
	r.netMsgOutput = netMsgOutput
	return r
}

func (r *RpcManager) Start() {
	for i := 0; i < 1; i++ {
		go func() {
			for {
				netMsg := <-r.netMsgOutput
				r.sendNetMsgToGenshinGateway(netMsg)
			}
		}()
	}
}

func (r *RpcManager) sendNetMsgToGenshinGateway(netMsg *api.NetMsg) {
	if _, ok := (netMsg.PayloadMessage).(*proto.NullMsg); ok {
		// 沙比gob没法序列化空结构体
		netMsg.PayloadMessage = nil
	}
	var result bool
	ok := r.genshinGatewayConsumer.CallFunction("RpcManager", "RecvNetMsgFromGameServer", netMsg, &result)
	if ok == true && result == true {
		return
	}
	return
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

// rpc interface
func (r *RpcManager) RecvNetMsgFromGenshinGateway(netMsg *api.NetMsg, result *bool) error {
	r.netMsgInput <- netMsg
	*result = true
	return nil
}
