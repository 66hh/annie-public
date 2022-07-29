//go:build linux
// +build linux

package kcp

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// the read loop for a client session
func (s *UDPSession) readLoop() {
	// default version
	if s.xconn == nil {
		s.defaultReadLoop()
		return
	}

	// x/net version
	var src string
	msgs := make([]ipv4.Message, batchSize)
	for k := range msgs {
		msgs[k].Buffers = [][]byte{make([]byte, mtuLimit)}
	}

	for {
		if count, err := s.xconn.ReadBatch(msgs, 0); err == nil {
			for i := 0; i < count; i++ {
				msg := &msgs[i]

				// make sure the packet is from the same source
				if src == "" { // set source address if nil
					src = msg.Addr.String()
				} else if msg.Addr.String() != src {
					atomic.AddUint64(&DefaultSnmp.InErrs, 1)
					continue
				}

				udpPayload := msg.Buffers[0][:msg.N]

				// source and size has validated
				s.packetInput(udpPayload)
			}
		} else {
			// compatibility issue:
			// for linux kernel<=2.6.32, support for sendmmsg is not available
			// an error of type os.SyscallError will be returned
			if operr, ok := err.(*net.OpError); ok {
				if se, ok := operr.Err.(*os.SyscallError); ok {
					if se.Syscall == "recvmmsg" {
						s.defaultReadLoop()
						return
					}
				}
			}
			s.notifyReadError(errors.WithStack(err))
			return
		}
	}
}

// monitor incoming data for all connections of server
func (l *Listener) monitor() {
	var xconn batchConn
	if _, ok := l.conn.(*net.UDPConn); ok {
		addr, err := net.ResolveUDPAddr("udp", l.conn.LocalAddr().String())
		if err == nil {
			if addr.IP.To4() != nil {
				xconn = ipv4.NewPacketConn(l.conn)
			} else {
				xconn = ipv6.NewPacketConn(l.conn)
			}
		}
	}

	// default version
	if xconn == nil {
		l.defaultMonitor()
		return
	}

	// x/net version
	msgs := make([]ipv4.Message, batchSize)
	for k := range msgs {
		msgs[k].Buffers = [][]byte{make([]byte, mtuLimit)}
	}

	for {
		if count, err := xconn.ReadBatch(msgs, 0); err == nil {
			for i := 0; i < count; i++ {
				msg := &msgs[i]
				udpPayload := msg.Buffers[0][:msg.N]
				if msg.N == 20 {
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
					conn, exist := l.sessions[msg.Addr.String()]
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
							Addr:     msg.Addr.String(),
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
							Addr:     msg.Addr.String(),
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
							Addr:     msg.Addr.String(),
							ConvId:   convId,
							ConnType: ConnEnetFin,
							EnetType: enetType,
						}
						continue
					}
				}
				l.packetInput(udpPayload, msg.Addr)
			}
		} else {
			// compatibility issue:
			// for linux kernel<=2.6.32, support for sendmmsg is not available
			// an error of type os.SyscallError will be returned
			if operr, ok := err.(*net.OpError); ok {
				if se, ok := operr.Err.(*os.SyscallError); ok {
					if se.Syscall == "recvmmsg" {
						l.defaultMonitor()
						return
					}
				}
			}
			l.notifyReadError(errors.WithStack(err))
			return
		}
	}
}
