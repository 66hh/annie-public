package proto

import pb "google.golang.org/protobuf/proto"

const (
	NormalMsg = iota
	UserLogin
	UserOffline
	ClientRttNotify
	ClientTimeNotify
)

type NetMsg struct {
	UserId             uint32      `msgpack:"UserId"`
	EventId            uint16      `msgpack:"EventId"`
	ApiId              uint16      `msgpack:"ApiId"`
	HeadMessage        *PacketHead `msgpack:"-"`
	HeadMessageData    []byte      `msgpack:"HeadMessageData"`
	PayloadMessage     pb.Message  `msgpack:"-"`
	PayloadMessageData []byte      `msgpack:"PayloadMessageData"`
	ClientRtt          uint32      `msgpack:"ClientRtt"`
	ClientTime         uint32      `msgpack:"ClientTime"`
}
