package rpc

import (
	"flswld.com/common/config"
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/light"
	"flswld.com/logger"
	"gate-genshin/dao"
	"gate-genshin/net"
	"gate-genshin/region"
	"io/ioutil"
	"sync"
	"time"
)

const (
	ConnWaitToken = iota
	ConnWaitLogin
	ConnAlive
)

type RpcManager struct {
	conf                *config.Config
	log                 *logger.Logger
	dao                 *dao.Dao
	gameServiceConsumer *light.Consumer
	protoMsgInput       chan *net.ProtoMsg
	protoMsgOutput      chan *net.ProtoMsg
	// 玩家登录相关
	connStateMap     map[uint64]uint8
	connStateMapLock sync.RWMutex
	// kcpConv -> userID
	convUserIdMap     map[uint64]uint32
	convUserIdMapLock sync.RWMutex
	// userID -> kcpConv
	userIdConvMap     map[uint32]uint64
	userIdConvMapLock sync.RWMutex
	// kcpConv -> ipAddr
	convAddrMap     map[uint64]string
	convAddrMapLock sync.RWMutex
	secretKeyBuffer []byte
	kcpEventInput   chan *net.KcpEvent
	kcpEventOutput  chan *net.KcpEvent
	regionCurr      *proto.QueryCurrRegionHttpRsp
	signRsaKey      []byte
	encRsaKey       []byte
}

func NewRpcManager(conf *config.Config, log *logger.Logger, dao *dao.Dao, gameServiceConsumer *light.Consumer, protoMsgInput chan *net.ProtoMsg, protoMsgOutput chan *net.ProtoMsg, kcpEventInput chan *net.KcpEvent, kcpEventOutput chan *net.KcpEvent) (r *RpcManager) {
	r = new(RpcManager)
	r.conf = conf
	r.log = log
	r.dao = dao
	r.gameServiceConsumer = gameServiceConsumer
	r.protoMsgInput = protoMsgInput
	r.protoMsgOutput = protoMsgOutput
	r.connStateMap = make(map[uint64]uint8)
	r.convUserIdMap = make(map[uint64]uint32)
	r.userIdConvMap = make(map[uint32]uint64)
	r.convAddrMap = make(map[uint64]string)
	r.kcpEventInput = kcpEventInput
	r.kcpEventOutput = kcpEventOutput
	return r
}

func (r *RpcManager) getHeadMsg(seq uint32) (headMsg *api.PacketHead) {
	headMsg = new(api.PacketHead)
	headMsg.ClientSequenceId = seq
	headMsg.Timestamp = uint64(time.Now().UnixMilli())
	return headMsg
}

func (r *RpcManager) kcpEventHandle() {
	for {
		event := <-r.kcpEventOutput
		r.log.Info("rpc manager recv event, ConvId: %v, EventId: %v", event.ConvId, event.EventId)
		switch event.EventId {
		case net.KcpPacketSendNotify:
			// 关闭发包监听
			r.kcpEventInput <- &net.KcpEvent{
				ConvId:       event.ConvId,
				EventId:      net.KcpPacketSendListen,
				EventMessage: "Disable",
			}
			// 登录成功 通知GS初始化相关数据
			userId, exist := r.getUserIdByConvId(event.ConvId)
			if !exist {
				r.log.Error("can not find userId by convId")
				continue
			}
			netMsg := new(api.NetMsg)
			netMsg.UserId = userId
			netMsg.EventId = api.UserLogin
			netMsg.ApiId = 0
			netMsg.HeadMessage = nil
			netMsg.PayloadMessage = nil
			r.sendNetMsgToGameServer(netMsg)
			r.log.Info("send to gs user login ok, ConvId: %v, UserId: %v", event.ConvId, netMsg.UserId)
		case net.KcpConnCloseNotify:
			// 连接断开
			userId, exist := r.getUserIdByConvId(event.ConvId)
			if !exist {
				r.log.Error("can not find userId by convId")
				continue
			}
			if r.getConnState(event.ConvId) == ConnAlive {
				// 通知GS玩家下线
				netMsg := new(api.NetMsg)
				netMsg.UserId = userId
				netMsg.EventId = api.UserOffline
				netMsg.ApiId = 0
				netMsg.HeadMessage = nil
				netMsg.PayloadMessage = nil
				r.sendNetMsgToGameServer(netMsg)
				r.log.Info("send to gs user offline, ConvId: %v, UserId: %v", event.ConvId, netMsg.UserId)
			}
			// 删除各种map数据
			r.deleteConnState(event.ConvId)
			r.deleteUserIdByConvId(event.ConvId)
			currConvId, currExist := r.getConvIdByUserId(userId)
			if currExist && currConvId == event.ConvId {
				// 防止误删顶号的新连接数据
				r.deleteConvIdByUserId(userId)
			}
			r.deleteAddrByConvId(event.ConvId)
		case net.KcpConnEstNotify:
			addr, ok := event.EventMessage.(string)
			if !ok {
				r.log.Error("event KcpConnEstNotify msg type error")
				continue
			}
			r.setAddrByConvId(event.ConvId, addr)
		}
	}
}

func (r *RpcManager) getConnState(convId uint64) uint8 {
	r.connStateMapLock.RLock()
	connState, connStateExist := r.connStateMap[convId]
	r.connStateMapLock.RUnlock()
	if !connStateExist {
		connState = ConnWaitToken
		r.connStateMapLock.Lock()
		r.connStateMap[convId] = ConnWaitToken
		r.connStateMapLock.Unlock()
	}
	return connState
}

func (r *RpcManager) setConnState(convId uint64, state uint8) {
	r.connStateMapLock.Lock()
	r.connStateMap[convId] = state
	r.connStateMapLock.Unlock()
}

func (r *RpcManager) deleteConnState(convId uint64) {
	r.connStateMapLock.Lock()
	delete(r.connStateMap, convId)
	r.connStateMapLock.Unlock()
}

func (r *RpcManager) getUserIdByConvId(convId uint64) (userId uint32, exist bool) {
	r.convUserIdMapLock.RLock()
	userId, exist = r.convUserIdMap[convId]
	r.convUserIdMapLock.RUnlock()
	return userId, exist
}

func (r *RpcManager) setUserIdByConvId(convId uint64, userId uint32) {
	r.convUserIdMapLock.Lock()
	r.convUserIdMap[convId] = userId
	r.convUserIdMapLock.Unlock()
}

func (r *RpcManager) deleteUserIdByConvId(convId uint64) {
	r.convUserIdMapLock.Lock()
	delete(r.convUserIdMap, convId)
	r.convUserIdMapLock.Unlock()
}

func (r *RpcManager) getConvIdByUserId(userId uint32) (convId uint64, exist bool) {
	r.userIdConvMapLock.RLock()
	convId, exist = r.userIdConvMap[userId]
	r.userIdConvMapLock.RUnlock()
	return convId, exist
}

func (r *RpcManager) setConvIdByUserId(userId uint32, convId uint64) {
	r.userIdConvMapLock.Lock()
	r.userIdConvMap[userId] = convId
	r.userIdConvMapLock.Unlock()
}

func (r *RpcManager) deleteConvIdByUserId(userId uint32) {
	r.userIdConvMapLock.Lock()
	delete(r.userIdConvMap, userId)
	r.userIdConvMapLock.Unlock()
}

func (r *RpcManager) getAddrByConvId(convId uint64) (addr string, exist bool) {
	r.convAddrMapLock.RLock()
	addr, exist = r.convAddrMap[convId]
	r.convAddrMapLock.RUnlock()
	return addr, exist
}

func (r *RpcManager) setAddrByConvId(convId uint64, addr string) {
	r.convAddrMapLock.Lock()
	r.convAddrMap[convId] = addr
	r.convAddrMapLock.Unlock()
}

func (r *RpcManager) deleteAddrByConvId(convId uint64) {
	r.convAddrMapLock.Lock()
	delete(r.convAddrMap, convId)
	r.convAddrMapLock.Unlock()
}

func (r *RpcManager) Start() {
	// 读取密钥相关文件
	var err error = nil
	r.secretKeyBuffer, err = ioutil.ReadFile("static/secretKeyBuffer.bin")
	if err != nil {
		r.log.Error("open secretKeyBuffer.bin error")
		return
	}
	r.signRsaKey, r.encRsaKey, _ = region.LoadRsaKey(r.log)
	// region
	regionCurr, _ := region.InitRegion(r.log, r.conf.Genshin.KcpAddr, r.conf.Genshin.KcpPort)
	r.regionCurr = regionCurr
	// kcp事件监听
	go r.kcpEventHandle()
	// 接收客户端消息
	for i := 0; i < 100; i++ {
		go func(index int) {
			for {
				protoMsg := <-r.protoMsgOutput
				connState := r.getConnState(protoMsg.ConvId)
				// gate本地处理的请求
				if protoMsg.ApiId == api.ApiPingReq {
					// ping请求
					// 未登录禁止ping
					if connState != ConnAlive {
						continue
					}
					pingReq := protoMsg.PayloadMessage.(*proto.PingReq)
					// TODO 记录客户端最后一次ping时间做超时下线处理
					pingRsp := new(proto.PingRsp)
					pingRsp.ClientTime = uint32(pingReq.ClientTime)
					// 返回数据到客户端
					resp := new(net.ProtoMsg)
					resp.ConvId = protoMsg.ConvId
					resp.ApiId = api.ApiPingRsp
					resp.HeadMessage = r.getHeadMsg(protoMsg.HeadMessage.ClientSequenceId)
					resp.PayloadMessage = pingRsp
					r.protoMsgInput <- resp
				} else if protoMsg.ApiId == api.ApiGetPlayerTokenReq {
					// 获取玩家token请求
					if connState != ConnWaitToken {
						continue
					}
					getPlayerTokenReq := protoMsg.PayloadMessage.(*proto.GetPlayerTokenReq)
					getPlayerTokenRsp := r.getPlayerToken(protoMsg.ConvId, getPlayerTokenReq)
					if getPlayerTokenRsp == nil {
						continue
					}
					// 改变解密密钥
					r.kcpEventInput <- &net.KcpEvent{
						ConvId:       protoMsg.ConvId,
						EventId:      net.KcpXorKeyChange,
						EventMessage: "DEC",
					}
					// 返回数据到客户端
					resp := new(net.ProtoMsg)
					resp.ConvId = protoMsg.ConvId
					resp.ApiId = api.ApiGetPlayerTokenRsp
					resp.HeadMessage = r.getHeadMsg(11)
					resp.PayloadMessage = getPlayerTokenRsp
					r.protoMsgInput <- resp
				} else if protoMsg.ApiId == api.ApiPlayerLoginReq {
					// 玩家登录请求
					if connState != ConnWaitLogin {
						continue
					}
					playerLoginReq := protoMsg.PayloadMessage.(*proto.PlayerLoginReq)
					playerLoginRsp := r.playerLogin(protoMsg.ConvId, playerLoginReq)
					if playerLoginRsp == nil {
						continue
					}
					// 改变加密密钥
					r.kcpEventInput <- &net.KcpEvent{
						ConvId:       protoMsg.ConvId,
						EventId:      net.KcpXorKeyChange,
						EventMessage: "ENC",
					}
					// 开启发包监听
					r.kcpEventInput <- &net.KcpEvent{
						ConvId:       protoMsg.ConvId,
						EventId:      net.KcpPacketSendListen,
						EventMessage: "Enable",
					}
					go func() {
						// 保证kcp事件已成功生效
						time.Sleep(time.Millisecond * 50)
						// 返回数据到客户端
						resp := new(net.ProtoMsg)
						resp.ConvId = protoMsg.ConvId
						resp.ApiId = api.ApiPlayerLoginRsp
						resp.HeadMessage = r.getHeadMsg(1)
						resp.PayloadMessage = playerLoginRsp
						r.protoMsgInput <- resp
					}()
				} else {
					// 转发到GS
					// 未登录禁止访问GS
					if connState != ConnAlive {
						continue
					}
					netMsg := new(api.NetMsg)
					userId, exist := r.getUserIdByConvId(protoMsg.ConvId)
					if exist {
						netMsg.UserId = userId
					} else {
						r.log.Error("can not find userId by convId")
						continue
					}
					netMsg.EventId = api.NormalMsg
					netMsg.ApiId = protoMsg.ApiId
					netMsg.HeadMessage = protoMsg.HeadMessage
					netMsg.PayloadMessage = protoMsg.PayloadMessage
					r.sendNetMsgToGameServer(netMsg)
				}
			}
		}(i)
	}
}

func (r *RpcManager) sendNetMsgToGameServer(netMsg *api.NetMsg) {
	var flag bool
	ok := r.gameServiceConsumer.CallFunction("RpcManager", "RecvNetMsgFromGenshinGateway", netMsg, &flag)
	if ok == true && flag == true {
		return
	}
	return
}
