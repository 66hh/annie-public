package kcp

import (
	"bytes"
	"encoding/binary"
	"sync/atomic"

	"github.com/pkg/errors"
)

func (s *UDPSession) defaultReadLoop() {
	buf := make([]byte, mtuLimit)
	var src string
	for {
		if n, addr, err := s.conn.ReadFrom(buf); err == nil {
			udpPayload := buf[:n]

			// make sure the packet is from the same source
			if src == "" { // set source address
				src = addr.String()
			} else if addr.String() != src {
				atomic.AddUint64(&DefaultSnmp.InErrs, 1)
				continue
			}

			s.packetInput(udpPayload)
		} else {
			s.notifyReadError(errors.WithStack(err))
			return
		}
	}
}

func (l *Listener) defaultMonitor() {
	buf := make([]byte, mtuLimit)
	for {
		if n, from, err := l.conn.ReadFrom(buf); err == nil {
			udpPayload := buf[:n]
			if n == 20 {
				// 原神KCP的Enet协议
				// 提取Enet协议头部和尾部幻数
				udpPayloadEnetHead := udpPayload[:4]
				udpPayloadEnetTail := udpPayload[len(udpPayload)-4:]
				// 提取Enet协议类型
				enetTypeData := udpPayload[12:16]
				enetTypeDataBuffer := bytes.NewBuffer(enetTypeData)
				var enetType uint32
				_ = binary.Read(enetTypeDataBuffer, binary.BigEndian, &enetType)
				l.sessionLock.RLock()
				conn, exist := l.sessions[from.String()]
				l.sessionLock.RUnlock()
				var convId uint64 = 0
				if exist {
					convId = conn.GetConv()
				}
				equalHead := bytes.Compare(udpPayloadEnetHead, MagicEnetSynHead)
				equalTail := bytes.Compare(udpPayloadEnetTail, MagicEnetSynTail)
				if equalHead == 0 && equalTail == 0 {
					// 客户端前置握手获取conv
					l.EnetNotify <- &Enet{
						Addr:     from.String(),
						ConvId:   convId,
						ConnType: ConnEnetSyn,
						EnetType: enetType,
					}
					continue
				}
				equalHead = bytes.Compare(udpPayloadEnetHead, MagicEnetEstHead)
				equalTail = bytes.Compare(udpPayloadEnetTail, MagicEnetEstTail)
				if equalHead == 0 && equalTail == 0 {
					// 连接建立
					l.EnetNotify <- &Enet{
						Addr:     from.String(),
						ConvId:   convId,
						ConnType: ConnEnetEst,
						EnetType: enetType,
					}
					continue
				}
				equalHead = bytes.Compare(udpPayloadEnetHead, MagicEnetFinHead)
				equalTail = bytes.Compare(udpPayloadEnetTail, MagicEnetFinTail)
				if equalHead == 0 && equalTail == 0 {
					// 连接断开
					l.EnetNotify <- &Enet{
						Addr:     from.String(),
						ConvId:   convId,
						ConnType: ConnEnetFin,
						EnetType: enetType,
					}
					continue
				}
			}
			l.packetInput(udpPayload, from)
		} else {
			l.notifyReadError(errors.WithStack(err))
			return
		}
	}
}
