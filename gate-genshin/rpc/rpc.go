package rpc

import "gate-genshin/forward"

type RpcManager struct {
	forwardManager *forward.ForwardManager
}

func NewRpcManager(forwardManager *forward.ForwardManager) (r *RpcManager) {
	r = new(RpcManager)
	r.forwardManager = forwardManager
	return r
}
