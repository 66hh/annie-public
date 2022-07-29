package net

import "flswld.com/gate-genshin-api/api"

type KcpMsg struct {
	ConvId    uint64
	ApiId     uint16
	HeadData  []byte
	ProtoData []byte
}

type ProtoMsg struct {
	ConvId         uint64
	ApiId          uint16
	HeadMessage    *api.PacketHead
	PayloadMessage any
}
