package rpc

import (
	"flswld.com/gate-genshin-api/api"
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
	var flag bool
	ok := r.genshinGatewayConsumer.CallFunction("RpcManager", "RecvNetMsgFromGameServer", netMsg, &flag)
	if ok == true && flag == true {
		return
	}
	return
}

// rpc interface
func (r *RpcManager) RecvNetMsgFromGenshinGateway(netMsg *api.NetMsg, res *bool) error {
	r.netMsgInput <- netMsg
	*res = true
	return nil
}
