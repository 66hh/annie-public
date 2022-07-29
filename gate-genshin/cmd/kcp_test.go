package main

import (
	"fmt"
	"gate-genshin/kcp"
	"testing"
	"time"
)

func TestKcpServer(t *testing.T) {
	//port := strconv.FormatInt(int64(k.conf.KcpPort), 10)
	listener, err := kcp.ListenWithOptions("0.0.0.0:22102", nil, 0, 0)
	if err != nil {
		panic(err)
	}
	conn, err := listener.AcceptKCP()
	conn.SetACKNoDelay(true)
	conn.SetWriteDelay(false)
	convId := conn.GetConv()
	fmt.Printf("convId: %v\n", convId)
	for {
		recvBuf := make([]byte, 1024*10)
		recvLen, err := conn.Read(recvBuf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("recv len: %v data: %v\n", recvLen, recvBuf[:recvLen])
		bin := []byte{0xAA, 0xBB, 0xCC}
		sendLen, err := conn.Write(bin)
		if err != nil {
			panic(err)
		}
		fmt.Printf("send len: %v data: %v\n", sendLen, bin)
	}
}

func TestKcpClient(t *testing.T) {
	//port := strconv.FormatInt(int64(k.conf.KcpPort), 10)
	conn, err := kcp.DialWithOptions("127.0.0.1:22102", nil, 0, 0)
	if err != nil {
		panic(err)
	}
	raw := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}
	bin := make([]byte, 0)
	for i := 0; i < 200; i++ {
		bin = append(bin, raw...)
	}
	for {
		sendLen, err := conn.Write(bin)
		if err != nil {
			panic(err)
		}
		fmt.Printf("send len: %v data: %v\n", sendLen, bin)
		recvBuf := make([]byte, 1024*10)
		recvLen, err := conn.Read(recvBuf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("recv len: %v data: %v\n", recvLen, recvBuf[:recvLen])
		time.Sleep(time.Second * 10)
	}
}
