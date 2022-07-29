package net

import (
	"bytes"
	"encoding/binary"
	"flswld.com/common/config"
	"flswld.com/common/utils/random"
	"flswld.com/logger"
	"gate-genshin/kcp"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

type KcpXorKey struct {
	encKey []byte
	decKey []byte
}

const (
	KcpXorKeyChange = iota
	KcpPacketRecvListen
	KcpPacketSendListen
	KcpConnForceClose
	KcpAllConnForceClose
	KcpGateOpenState
	KcpPacketRecvNotify
	KcpPacketSendNotify
	KcpConnCloseNotify
	KcpConnEstNotify
)

type KcpEvent struct {
	ConvId       uint64
	EventId      int
	EventMessage any
}

type KcpConnectManager struct {
	conf                  *config.Config
	log                   *logger.Logger
	openState             bool
	connMap               map[uint64]*kcp.UDPSession
	connMapLock           sync.RWMutex
	kcpEventInput         chan *KcpEvent
	kcpEventOutput        chan *KcpEvent
	kcpMsgInput           chan *KcpMsg
	kcpMsgOutput          chan *KcpMsg
	kcpRawSendChanMap     map[uint64]chan *KcpMsg
	kcpRawSendChanMapLock sync.RWMutex
	// 收包发包监听标志
	kcpRecvListenMap     map[uint64]bool
	kcpRecvListenMapLock sync.RWMutex
	kcpSendListenMap     map[uint64]bool
	kcpSendListenMapLock sync.RWMutex
	// key
	dispatchKey   []byte
	secretKey     []byte
	kcpKeyMap     map[uint64]*KcpXorKey
	kcpKeyMapLock sync.RWMutex
}

func NewKcpConnectManager(conf *config.Config, log *logger.Logger, kcpEventInput chan *KcpEvent, kcpEventOutput chan *KcpEvent, kcpMsgInput chan *KcpMsg, kcpMsgOutput chan *KcpMsg) (r *KcpConnectManager) {
	r = new(KcpConnectManager)
	r.conf = conf
	r.log = log
	r.openState = true
	r.connMap = make(map[uint64]*kcp.UDPSession)
	r.kcpEventInput = kcpEventInput
	r.kcpEventOutput = kcpEventOutput
	r.kcpMsgInput = kcpMsgInput
	r.kcpMsgOutput = kcpMsgOutput
	r.kcpRawSendChanMap = make(map[uint64]chan *KcpMsg)
	r.kcpRecvListenMap = make(map[uint64]bool)
	r.kcpSendListenMap = make(map[uint64]bool)
	r.kcpKeyMap = make(map[uint64]*KcpXorKey)
	return r
}

func (k *KcpConnectManager) Start() {
	go func() {
		// key
		var err error = nil
		k.dispatchKey, err = ioutil.ReadFile("static/dispatchKey.bin")
		if err != nil {
			k.log.Error("open dispatchKey.bin error")
			return
		}
		k.secretKey, err = ioutil.ReadFile("static/secretKey.bin")
		if err != nil {
			k.log.Error("open secretKey.bin error")
			return
		}
		// kcp
		port := strconv.FormatInt(int64(k.conf.Genshin.KcpPort), 10)
		listener, err := kcp.ListenWithOptions("0.0.0.0:"+port, nil, 0, 0)
		if err != nil {
			k.log.Error("listen kcp err: %v", err)
			return
		} else {
			go k.enetHandle(listener)
			go k.chanSendHandle()
			go k.eventHandle()
			for {
				conn, err := listener.AcceptKCP()
				if err != nil {
					k.log.Error("accept kcp err: %v", err)
					return
				}
				if k.openState == false {
					_ = conn.Close()
					continue
				}
				conn.SetACKNoDelay(true)
				conn.SetWriteDelay(false)
				convId := conn.GetConv()
				k.log.Debug("client connect, convId: %v", convId)
				// 连接建立成功通知
				k.kcpEventOutput <- &KcpEvent{
					ConvId:       convId,
					EventId:      KcpConnEstNotify,
					EventMessage: conn.RemoteAddr().String(),
				}
				k.connMapLock.Lock()
				k.connMap[convId] = conn
				k.connMapLock.Unlock()
				k.kcpKeyMapLock.Lock()
				k.kcpKeyMap[convId] = &KcpXorKey{
					encKey: k.dispatchKey,
					decKey: k.dispatchKey,
				}
				k.kcpKeyMapLock.Unlock()
				go k.recvHandle(convId, conn)
				kcpRawSendChan := make(chan *KcpMsg, 1000)
				k.kcpRawSendChanMapLock.Lock()
				k.kcpRawSendChanMap[convId] = kcpRawSendChan
				k.kcpRawSendChanMapLock.Unlock()
				go k.sendHandle(convId, conn, kcpRawSendChan)
			}
		}
	}()
}

func (k *KcpConnectManager) enetHandle(listener *kcp.Listener) {
	for {
		enetNotify := <-listener.EnetNotify
		k.log.Info("[Enet Notify], addr: %v, conv: %v, conn: %v, enet: %v", enetNotify.Addr, enetNotify.ConvId, enetNotify.ConnType, enetNotify.EnetType)
		switch enetNotify.ConnType {
		case kcp.ConnEnetSyn:
			if enetNotify.EnetType == kcp.EnetClientConnectKey {
				convData := random.GetRandomByte(8)
				convDataBuffer := bytes.NewBuffer(convData)
				var conv uint64
				_ = binary.Read(convDataBuffer, binary.BigEndian, &conv)
				listener.SendEnetNotifyToClient(&kcp.Enet{
					Addr:     enetNotify.Addr,
					ConvId:   conv,
					ConnType: kcp.ConnEnetEst,
					EnetType: enetNotify.EnetType,
				})
			}
		case kcp.ConnEnetEst:
		case kcp.ConnEnetFin:
			k.closeKcpConn(enetNotify.ConvId, enetNotify.EnetType)
		default:
		}
	}
}

func (k *KcpConnectManager) chanSendHandle() {
	// 分发到每个连接具体的发送协程
	for {
		kcpMsg := <-k.kcpMsgInput
		k.kcpRawSendChanMapLock.RLock()
		kcpRawSendChan := k.kcpRawSendChanMap[kcpMsg.ConvId]
		k.kcpRawSendChanMapLock.RUnlock()
		if kcpRawSendChan != nil {
			select {
			case kcpRawSendChan <- kcpMsg:
			default:
				k.log.Error("kcpRawSendChan is full, convId: %v", kcpMsg.ConvId)
			}
		} else {
			k.log.Error("kcpRawSendChan is nil, convId: %v", kcpMsg.ConvId)
		}
	}
}

func (k *KcpConnectManager) recvHandle(convId uint64, conn *kcp.UDPSession) {
	// 接收
	for {
		recvBuf := make([]byte, 384000)
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		recvLen, err := conn.Read(recvBuf)
		if err != nil {
			k.log.Error("exit recv loop, conn read err: %v, convId: %v", err, convId)
			k.closeKcpConn(convId, kcp.EnetServerKick)
			break
		}
		recvBuf = recvBuf[:recvLen]
		k.kcpRecvListenMapLock.RLock()
		flag := k.kcpRecvListenMap[convId]
		k.kcpRecvListenMapLock.RUnlock()
		if flag {
			// 收包通知
			k.kcpEventOutput <- &KcpEvent{
				ConvId:       convId,
				EventId:      KcpPacketRecvNotify,
				EventMessage: recvBuf,
			}
		}
		kcpMsgList := make([]*KcpMsg, 0)
		k.decodeBinToPayload(recvBuf, convId, &kcpMsgList)
		for _, v := range kcpMsgList {
			k.kcpMsgOutput <- v
		}
	}
}

func (k *KcpConnectManager) sendHandle(convId uint64, conn *kcp.UDPSession, kcpRawSendChan chan *KcpMsg) {
	// 发送
	for {
		kcpMsg, ok := <-kcpRawSendChan
		if !ok {
			k.log.Error("exit send loop, send chan close, convId: %v", convId)
			k.closeKcpConn(convId, kcp.EnetServerKick)
			break
		}
		bin := k.encodePayloadToBin(kcpMsg)
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		_, err := conn.Write(bin)
		if err != nil {
			k.log.Error("exit send loop, conn write err: %v, convId: %v", err, convId)
			k.closeKcpConn(convId, kcp.EnetServerKick)
			break
		}
		k.kcpSendListenMapLock.RLock()
		flag := k.kcpSendListenMap[convId]
		k.kcpSendListenMapLock.RUnlock()
		if flag {
			// 发包通知
			k.kcpEventOutput <- &KcpEvent{
				ConvId:       convId,
				EventId:      KcpPacketSendNotify,
				EventMessage: bin,
			}
		}
	}
}

func (k *KcpConnectManager) closeKcpConn(convId uint64, enetType uint32) {
	k.connMapLock.RLock()
	conn, exist := k.connMap[convId]
	k.connMapLock.RUnlock()
	if !exist {
		return
	}
	// 获取待关闭的发送管道
	k.kcpRawSendChanMapLock.RLock()
	kcpRawSendChan := k.kcpRawSendChanMap[convId]
	k.kcpRawSendChanMapLock.RUnlock()
	// 清理数据
	k.connMapLock.Lock()
	delete(k.connMap, convId)
	k.connMapLock.Unlock()
	k.kcpRawSendChanMapLock.Lock()
	delete(k.kcpRawSendChanMap, convId)
	k.kcpRawSendChanMapLock.Unlock()
	k.kcpRecvListenMapLock.Lock()
	delete(k.kcpRecvListenMap, convId)
	k.kcpRecvListenMapLock.Unlock()
	k.kcpSendListenMapLock.Lock()
	delete(k.kcpSendListenMap, convId)
	k.kcpSendListenMapLock.Unlock()
	k.kcpKeyMapLock.Lock()
	delete(k.kcpKeyMap, convId)
	k.kcpKeyMapLock.Unlock()
	// 关闭连接
	conn.SendEnetNotify(&kcp.Enet{
		ConnType: kcp.ConnEnetFin,
		EnetType: enetType,
	})
	_ = conn.Close()
	// 关闭发送管道
	close(kcpRawSendChan)
	// 连接关闭通知
	k.kcpEventOutput <- &KcpEvent{
		ConvId:  convId,
		EventId: KcpConnCloseNotify,
	}
}

func (k *KcpConnectManager) closeAllKcpConn() {
	closeConnList := make([]*kcp.UDPSession, 0)
	k.connMapLock.RLock()
	for _, v := range k.connMap {
		closeConnList = append(closeConnList, v)
	}
	k.connMapLock.RUnlock()
	for _, v := range closeConnList {
		k.closeKcpConn(v.GetConv(), kcp.EnetServerShutdown)
	}
}

func (k *KcpConnectManager) eventHandle() {
	// 事件处理
	for {
		event := <-k.kcpEventInput
		k.log.Info("kcp manager recv event, ConvId: %v, EventId: %v, EventMessage: %v", event.ConvId, event.EventId, event.EventMessage)
		switch event.EventId {
		case KcpXorKeyChange:
			// XOR密钥切换
			k.connMapLock.RLock()
			_, exist := k.connMap[event.ConvId]
			k.connMapLock.RUnlock()
			if !exist {
				k.log.Error("conn not exist, convId: %v", event.ConvId)
				continue
			}
			flag, ok := event.EventMessage.(string)
			if !ok {
				k.log.Error("event KcpXorKeyChange msg type error")
				continue
			}
			if flag == "ENC" {
				k.kcpKeyMapLock.Lock()
				k.kcpKeyMap[event.ConvId].encKey = k.secretKey
				k.kcpKeyMapLock.Unlock()
			} else if flag == "DEC" {
				k.kcpKeyMapLock.Lock()
				k.kcpKeyMap[event.ConvId].decKey = k.secretKey
				k.kcpKeyMapLock.Unlock()
			}
		case KcpPacketRecvListen:
			// 收包监听
			k.connMapLock.RLock()
			_, exist := k.connMap[event.ConvId]
			k.connMapLock.RUnlock()
			if !exist {
				k.log.Error("conn not exist, convId: %v", event.ConvId)
				continue
			}
			flag, ok := event.EventMessage.(string)
			if !ok {
				k.log.Error("event KcpXorKeyChange msg type error")
				continue
			}
			if flag == "Enable" {
				k.kcpRecvListenMapLock.Lock()
				k.kcpRecvListenMap[event.ConvId] = true
				k.kcpRecvListenMapLock.Unlock()
			} else if flag == "Disable" {
				k.kcpRecvListenMapLock.Lock()
				k.kcpRecvListenMap[event.ConvId] = false
				k.kcpRecvListenMapLock.Unlock()
			}
		case KcpPacketSendListen:
			// 发包监听
			k.connMapLock.RLock()
			_, exist := k.connMap[event.ConvId]
			k.connMapLock.RUnlock()
			if !exist {
				k.log.Error("conn not exist, convId: %v", event.ConvId)
				continue
			}
			flag, ok := event.EventMessage.(string)
			if !ok {
				k.log.Error("event KcpXorKeyChange msg type error")
				continue
			}
			if flag == "Enable" {
				k.kcpSendListenMapLock.Lock()
				k.kcpSendListenMap[event.ConvId] = true
				k.kcpSendListenMapLock.Unlock()
			} else if flag == "Disable" {
				k.kcpSendListenMapLock.Lock()
				k.kcpSendListenMap[event.ConvId] = false
				k.kcpSendListenMapLock.Unlock()
			}
		case KcpConnForceClose:
			// 强制关闭某个连接
			k.connMapLock.RLock()
			_, exist := k.connMap[event.ConvId]
			k.connMapLock.RUnlock()
			if !exist {
				k.log.Error("conn not exist, convId: %v", event.ConvId)
				continue
			}
			reason, ok := event.EventMessage.(uint32)
			if !ok {
				k.log.Error("event KcpConnForceClose msg type error")
				continue
			}
			k.closeKcpConn(event.ConvId, reason)
			k.log.Info("conn has been force close, convId: %v", event.ConvId)
		case KcpAllConnForceClose:
			// 强制关闭所有连接
			k.closeAllKcpConn()
			k.log.Info("all conn has been force close")
		case KcpGateOpenState:
			// 改变网关开放状态
			openState, ok := event.EventMessage.(bool)
			if !ok {
				k.log.Error("event KcpGateOpenState msg type error")
				continue
			}
			k.openState = openState
			if openState == false {
				k.closeAllKcpConn()
			}
		}
	}
}
