package forward

import (
	"flswld.com/common/config"
	"flswld.com/gate-genshin-api/gm"
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	"gate-genshin/dao"
	"gate-genshin/kcp"
	"gate-genshin/net"
	"gate-genshin/region"
	"io/ioutil"
	"runtime"
	"sync"
	"time"
)

const (
	ConnWaitToken = iota
	ConnWaitLogin
	ConnAlive
)

type ForwardManager struct {
	dao            *dao.Dao
	protoMsgInput  chan *net.ProtoMsg
	protoMsgOutput chan *net.ProtoMsg
	netMsgInput    chan *proto.NetMsg
	netMsgOutput   chan *proto.NetMsg
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
	// kcpConv -> seq
	convSeqMap      map[uint64]uint32
	convSeqMapLock  sync.RWMutex
	secretKeyBuffer []byte
	kcpEventInput   chan *net.KcpEvent
	kcpEventOutput  chan *net.KcpEvent
	regionCurr      *proto.QueryCurrRegionHttpRsp
	signRsaKey      []byte
	osEncRsaKey     []byte
	cnEncRsaKey     []byte
}

func NewForwardManager(dao *dao.Dao,
	protoMsgInput chan *net.ProtoMsg, protoMsgOutput chan *net.ProtoMsg,
	kcpEventInput chan *net.KcpEvent, kcpEventOutput chan *net.KcpEvent,
	netMsgInput chan *proto.NetMsg, netMsgOutput chan *proto.NetMsg) (r *ForwardManager) {
	r = new(ForwardManager)
	r.dao = dao
	r.protoMsgInput = protoMsgInput
	r.protoMsgOutput = protoMsgOutput
	r.netMsgInput = netMsgInput
	r.netMsgOutput = netMsgOutput
	r.connStateMap = make(map[uint64]uint8)
	r.convUserIdMap = make(map[uint64]uint32)
	r.userIdConvMap = make(map[uint32]uint64)
	r.convAddrMap = make(map[uint64]string)
	r.convSeqMap = make(map[uint64]uint32)
	r.kcpEventInput = kcpEventInput
	r.kcpEventOutput = kcpEventOutput
	return r
}

func (f *ForwardManager) getHeadMsg(convId uint64) (headMsg *proto.PacketHead) {
	headMsg = new(proto.PacketHead)
	seq, exist := f.getSeqByConvId(convId)
	if !exist {
		logger.LOG.Error("can not find seq by convId")
		return nil
	}
	f.setSeqByConvId(convId, seq+1)
	headMsg.ClientSequenceId = seq
	headMsg.SentMs = uint64(time.Now().UnixMilli())
	return headMsg
}

func (f *ForwardManager) kcpEventHandle() {
	for {
		event := <-f.kcpEventOutput
		logger.LOG.Info("rpc manager recv event, ConvId: %v, EventId: %v", event.ConvId, event.EventId)
		switch event.EventId {
		case net.KcpPacketSendNotify:
			// 发包通知
			// 关闭发包监听
			f.kcpEventInput <- &net.KcpEvent{
				ConvId:       event.ConvId,
				EventId:      net.KcpPacketSendListen,
				EventMessage: "Disable",
			}
			// 登录成功 通知GS初始化相关数据
			userId, exist := f.getUserIdByConvId(event.ConvId)
			if !exist {
				logger.LOG.Error("can not find userId by convId")
				continue
			}
			netMsg := new(proto.NetMsg)
			netMsg.UserId = userId
			netMsg.EventId = proto.UserLogin
			netMsg.ApiId = 0
			netMsg.HeadMessage = nil
			netMsg.PayloadMessage = nil
			f.netMsgInput <- netMsg
			logger.LOG.Info("send to gs user login ok, ConvId: %v, UserId: %v", event.ConvId, netMsg.UserId)
		case net.KcpConnCloseNotify:
			// 连接断开通知
			userId, exist := f.getUserIdByConvId(event.ConvId)
			if !exist {
				logger.LOG.Error("can not find userId by convId")
				continue
			}
			if f.getConnState(event.ConvId) == ConnAlive {
				// 通知GS玩家下线
				netMsg := new(proto.NetMsg)
				netMsg.UserId = userId
				netMsg.EventId = proto.UserOffline
				netMsg.ApiId = 0
				netMsg.HeadMessage = nil
				netMsg.PayloadMessage = nil
				f.netMsgInput <- netMsg
				logger.LOG.Info("send to gs user offline, ConvId: %v, UserId: %v", event.ConvId, netMsg.UserId)
			}
			// 删除各种map数据
			f.deleteConnState(event.ConvId)
			f.deleteUserIdByConvId(event.ConvId)
			currConvId, currExist := f.getConvIdByUserId(userId)
			if currExist && currConvId == event.ConvId {
				// 防止误删顶号的新连接数据
				f.deleteConvIdByUserId(userId)
			}
			f.deleteAddrByConvId(event.ConvId)
			f.deleteSeqByConvId(event.ConvId)
		case net.KcpConnEstNotify:
			// 连接建立通知
			addr, ok := event.EventMessage.(string)
			if !ok {
				logger.LOG.Error("event KcpConnEstNotify msg type error")
				continue
			}
			f.setAddrByConvId(event.ConvId, addr)
			f.setSeqByConvId(event.ConvId, 1)
		case net.KcpConnRttNotify:
			// 客户端往返时延通知
			rtt, ok := event.EventMessage.(int32)
			if !ok {
				logger.LOG.Error("event KcpConnRttNotify msg type error")
				continue
			}
			// 通知GS玩家客户端往返时延
			userId, exist := f.getUserIdByConvId(event.ConvId)
			if !exist {
				logger.LOG.Error("can not find userId by convId")
				continue
			}
			netMsg := new(proto.NetMsg)
			netMsg.UserId = userId
			netMsg.EventId = proto.ClientRttNotify
			netMsg.ApiId = 0
			netMsg.HeadMessage = nil
			netMsg.PayloadMessage = nil
			netMsg.ClientRtt = uint32(rtt)
			f.netMsgInput <- netMsg
		case net.KcpConnAddrChangeNotify:
			// 客户端网络地址改变通知
			f.convAddrMapLock.Lock()
			_, exist := f.convAddrMap[event.ConvId]
			if !exist {
				f.convAddrMapLock.Unlock()
				logger.LOG.Error("conn addr change but conn can not be found")
				continue
			}
			addr := event.EventMessage.(string)
			f.convAddrMap[event.ConvId] = addr
			f.convAddrMapLock.Unlock()
		}
	}
}

func (f *ForwardManager) Start() {
	// 读取密钥相关文件
	var err error = nil
	f.secretKeyBuffer, err = ioutil.ReadFile("static/secretKeyBuffer.bin")
	if err != nil {
		logger.LOG.Error("open secretKeyBuffer.bin error")
		return
	}
	f.signRsaKey, f.osEncRsaKey, f.cnEncRsaKey, _ = region.LoadRsaKey()
	// region
	regionCurr, _ := region.InitRegion(config.CONF.Genshin.KcpAddr, config.CONF.Genshin.KcpPort)
	f.regionCurr = regionCurr
	// kcp事件监听
	go f.kcpEventHandle()
	go f.recvNetMsgFromGameServer()
	// 接收客户端消息
	cpuCoreNum := runtime.NumCPU()
	for i := 0; i < cpuCoreNum*10; i++ {
		go f.sendNetMsgToGameServer()
	}
}

// 发送消息到GS
func (f *ForwardManager) sendNetMsgToGameServer() {
	for {
		protoMsg := <-f.protoMsgOutput
		connState := f.getConnState(protoMsg.ConvId)
		// gate本地处理的请求
		if protoMsg.ApiId == proto.ApiPingReq {
			// ping请求
			// 未登录禁止ping
			if connState != ConnAlive {
				continue
			}
			pingReq := protoMsg.PayloadMessage.(*proto.PingReq)
			logger.LOG.Debug("user ping req, data: %v", pingReq.String())
			// TODO 记录客户端最后一次ping时间做超时下线处理
			pingRsp := new(proto.PingRsp)
			pingRsp.ClientTime = pingReq.ClientTime
			// 返回数据到客户端
			resp := new(net.ProtoMsg)
			resp.ConvId = protoMsg.ConvId
			resp.ApiId = proto.ApiPingRsp
			resp.HeadMessage = f.getHeadMsg(protoMsg.ConvId)
			resp.PayloadMessage = pingRsp
			f.protoMsgInput <- resp
			// 通知GS玩家客户端的本地时钟
			userId, exist := f.getUserIdByConvId(protoMsg.ConvId)
			if !exist {
				logger.LOG.Error("can not find userId by convId")
				continue
			}
			netMsg := new(proto.NetMsg)
			netMsg.UserId = userId
			netMsg.EventId = proto.ClientTimeNotify
			netMsg.ApiId = 0
			netMsg.HeadMessage = nil
			netMsg.PayloadMessage = nil
			netMsg.ClientTime = pingReq.ClientTime
			f.netMsgInput <- netMsg
		} else if protoMsg.ApiId == proto.ApiGetPlayerTokenReq {
			// 获取玩家token请求
			if connState != ConnWaitToken {
				continue
			}
			getPlayerTokenReq := protoMsg.PayloadMessage.(*proto.GetPlayerTokenReq)
			getPlayerTokenRsp := f.getPlayerToken(protoMsg.ConvId, getPlayerTokenReq)
			if getPlayerTokenRsp == nil {
				continue
			}
			// 改变解密密钥
			f.kcpEventInput <- &net.KcpEvent{
				ConvId:       protoMsg.ConvId,
				EventId:      net.KcpXorKeyChange,
				EventMessage: "DEC",
			}
			// 返回数据到客户端
			resp := new(net.ProtoMsg)
			resp.ConvId = protoMsg.ConvId
			resp.ApiId = proto.ApiGetPlayerTokenRsp
			resp.HeadMessage = f.getHeadMsg(protoMsg.ConvId)
			resp.PayloadMessage = getPlayerTokenRsp
			f.protoMsgInput <- resp
		} else if protoMsg.ApiId == proto.ApiPlayerLoginReq {
			// 玩家登录请求
			if connState != ConnWaitLogin {
				continue
			}
			playerLoginReq := protoMsg.PayloadMessage.(*proto.PlayerLoginReq)
			playerLoginRsp := f.playerLogin(protoMsg.ConvId, playerLoginReq)
			if playerLoginRsp == nil {
				continue
			}
			// 改变加密密钥
			f.kcpEventInput <- &net.KcpEvent{
				ConvId:       protoMsg.ConvId,
				EventId:      net.KcpXorKeyChange,
				EventMessage: "ENC",
			}
			// 开启发包监听
			f.kcpEventInput <- &net.KcpEvent{
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
				resp.ApiId = proto.ApiPlayerLoginRsp
				resp.HeadMessage = f.getHeadMsg(protoMsg.ConvId)
				resp.PayloadMessage = playerLoginRsp
				f.protoMsgInput <- resp
			}()
		} else {
			// 转发到GS
			// 未登录禁止访问GS
			if connState != ConnAlive {
				continue
			}
			netMsg := new(proto.NetMsg)
			userId, exist := f.getUserIdByConvId(protoMsg.ConvId)
			if exist {
				netMsg.UserId = userId
			} else {
				logger.LOG.Error("can not find userId by convId")
				continue
			}
			netMsg.EventId = proto.NormalMsg
			netMsg.ApiId = protoMsg.ApiId
			netMsg.HeadMessage = protoMsg.HeadMessage
			netMsg.PayloadMessage = protoMsg.PayloadMessage
			f.netMsgInput <- netMsg
		}
	}
}

// 从GS接收消息
func (f *ForwardManager) recvNetMsgFromGameServer() {
	for {
		netMsg := <-f.netMsgOutput
		if netMsg.EventId == proto.NormalMsg {
			protoMsg := new(net.ProtoMsg)
			convId, exist := f.getConvIdByUserId(netMsg.UserId)
			if exist {
				protoMsg.ConvId = convId
			} else {
				logger.LOG.Error("can not find convId by userId")
				continue
			}
			protoMsg.ApiId = netMsg.ApiId
			protoMsg.HeadMessage = f.getHeadMsg(protoMsg.ConvId)
			protoMsg.PayloadMessage = netMsg.PayloadMessage
			f.protoMsgInput <- protoMsg
			continue
		} else {
			logger.LOG.Info("recv event from game server, event id: %v", netMsg.EventId)
			continue
		}
	}
}

func (f *ForwardManager) getConnState(convId uint64) uint8 {
	f.connStateMapLock.RLock()
	connState, connStateExist := f.connStateMap[convId]
	f.connStateMapLock.RUnlock()
	if !connStateExist {
		connState = ConnWaitToken
		f.connStateMapLock.Lock()
		f.connStateMap[convId] = ConnWaitToken
		f.connStateMapLock.Unlock()
	}
	return connState
}

func (f *ForwardManager) setConnState(convId uint64, state uint8) {
	f.connStateMapLock.Lock()
	f.connStateMap[convId] = state
	f.connStateMapLock.Unlock()
}

func (f *ForwardManager) deleteConnState(convId uint64) {
	f.connStateMapLock.Lock()
	delete(f.connStateMap, convId)
	f.connStateMapLock.Unlock()
}

func (f *ForwardManager) getUserIdByConvId(convId uint64) (userId uint32, exist bool) {
	f.convUserIdMapLock.RLock()
	userId, exist = f.convUserIdMap[convId]
	f.convUserIdMapLock.RUnlock()
	return userId, exist
}

func (f *ForwardManager) setUserIdByConvId(convId uint64, userId uint32) {
	f.convUserIdMapLock.Lock()
	f.convUserIdMap[convId] = userId
	f.convUserIdMapLock.Unlock()
}

func (f *ForwardManager) deleteUserIdByConvId(convId uint64) {
	f.convUserIdMapLock.Lock()
	delete(f.convUserIdMap, convId)
	f.convUserIdMapLock.Unlock()
}

func (f *ForwardManager) getConvIdByUserId(userId uint32) (convId uint64, exist bool) {
	f.userIdConvMapLock.RLock()
	convId, exist = f.userIdConvMap[userId]
	f.userIdConvMapLock.RUnlock()
	return convId, exist
}

func (f *ForwardManager) setConvIdByUserId(userId uint32, convId uint64) {
	f.userIdConvMapLock.Lock()
	f.userIdConvMap[userId] = convId
	f.userIdConvMapLock.Unlock()
}

func (f *ForwardManager) deleteConvIdByUserId(userId uint32) {
	f.userIdConvMapLock.Lock()
	delete(f.userIdConvMap, userId)
	f.userIdConvMapLock.Unlock()
}

func (f *ForwardManager) getAddrByConvId(convId uint64) (addr string, exist bool) {
	f.convAddrMapLock.RLock()
	addr, exist = f.convAddrMap[convId]
	f.convAddrMapLock.RUnlock()
	return addr, exist
}

func (f *ForwardManager) setAddrByConvId(convId uint64, addr string) {
	f.convAddrMapLock.Lock()
	f.convAddrMap[convId] = addr
	f.convAddrMapLock.Unlock()
}

func (f *ForwardManager) deleteAddrByConvId(convId uint64) {
	f.convAddrMapLock.Lock()
	delete(f.convAddrMap, convId)
	f.convAddrMapLock.Unlock()
}

func (f *ForwardManager) getSeqByConvId(convId uint64) (seq uint32, exist bool) {
	f.convSeqMapLock.RLock()
	seq, exist = f.convSeqMap[convId]
	f.convSeqMapLock.RUnlock()
	return seq, exist
}

func (f *ForwardManager) setSeqByConvId(convId uint64, seq uint32) {
	f.convSeqMapLock.Lock()
	f.convSeqMap[convId] = seq
	f.convSeqMapLock.Unlock()
}

func (f *ForwardManager) deleteSeqByConvId(convId uint64) {
	f.convSeqMapLock.Lock()
	delete(f.convSeqMap, convId)
	f.convSeqMapLock.Unlock()
}

// 改变网关开放状态
func (f *ForwardManager) ChangeGateOpenState(isOpen bool) bool {
	f.kcpEventInput <- &net.KcpEvent{
		EventId:      net.KcpGateOpenState,
		EventMessage: isOpen,
	}
	logger.LOG.Info("change gate open state to: %v", isOpen)
	return true
}

// 剔除玩家下线
func (f *ForwardManager) KickPlayer(info *gm.KickPlayerInfo) bool {
	if info == nil {
		return false
	}
	convId, exist := f.getConvIdByUserId(info.UserId)
	if !exist {
		return false
	}
	f.kcpEventInput <- &net.KcpEvent{
		ConvId:       convId,
		EventId:      net.KcpConnForceClose,
		EventMessage: info.Reason,
	}
	return true
}

// 获取网关在线玩家信息
func (f *ForwardManager) GetOnlineUser(uid uint32) (list *gm.OnlineUserList) {
	list = &gm.OnlineUserList{
		UserList: make([]*gm.OnlineUserInfo, 0),
	}
	if uid == 0 {
		// 获取全部玩家
		f.convUserIdMapLock.RLock()
		f.convAddrMapLock.RLock()
		for convId, userId := range f.convUserIdMap {
			addr := f.convAddrMap[convId]
			info := &gm.OnlineUserInfo{
				Uid:    userId,
				ConvId: convId,
				Addr:   addr,
			}
			list.UserList = append(list.UserList, info)
		}
		f.convAddrMapLock.RUnlock()
		f.convUserIdMapLock.RUnlock()
	} else {
		// 获取指定uid玩家
		convId, exist := f.getConvIdByUserId(uid)
		if !exist {
			return list
		}
		addr, exist := f.getAddrByConvId(convId)
		if !exist {
			return list
		}
		info := &gm.OnlineUserInfo{
			Uid:    uid,
			ConvId: convId,
			Addr:   addr,
		}
		list.UserList = append(list.UserList, info)
	}
	return list
}

// 用户密码改变
func (f *ForwardManager) UserPasswordChange(uid uint32) bool {
	// dispatch登录态失效
	_, err := f.dao.UpdateAccountFieldByFieldName("uid", uid, "token", "")
	if err != nil {
		return false
	}
	// 游戏内登录态失效
	account, err := f.dao.QueryAccountByField("uid", uid)
	if err != nil {
		return false
	}
	if account == nil {
		return false
	}
	convId, exist := f.getConvIdByUserId(uint32(account.PlayerID))
	if !exist {
		return true
	}
	f.kcpEventInput <- &net.KcpEvent{
		ConvId:       convId,
		EventId:      net.KcpConnForceClose,
		EventMessage: uint32(kcp.EnetAccountPasswordChange),
	}
	return true
}

// 封号
func (f *ForwardManager) ForbidUser(info *gm.ForbidUserInfo) bool {
	if info == nil {
		return false
	}
	// 写入账号封禁信息
	_, err := f.dao.UpdateAccountFieldByFieldName("uid", info.UserId, "forbid", true)
	if err != nil {
		return false
	}
	_, err = f.dao.UpdateAccountFieldByFieldName("uid", info.UserId, "forbidEndTime", info.ForbidEndTime)
	if err != nil {
		return false
	}
	// 游戏强制下线
	account, err := f.dao.QueryAccountByField("uid", info.UserId)
	if err != nil {
		return false
	}
	if account == nil {
		return false
	}
	convId, exist := f.getConvIdByUserId(uint32(account.PlayerID))
	if !exist {
		return true
	}
	f.kcpEventInput <- &net.KcpEvent{
		ConvId:       convId,
		EventId:      net.KcpConnForceClose,
		EventMessage: uint32(kcp.EnetServerKillClient),
	}
	return true
}

// 解封
func (f *ForwardManager) UnForbidUser(uid uint32) bool {
	// 解除账号封禁
	_, err := f.dao.UpdateAccountFieldByFieldName("uid", uid, "forbid", false)
	if err != nil {
		return false
	}
	return true
}
