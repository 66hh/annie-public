package api

const (
	NormalMsg = iota
	UserLogin
	UserOffline
	ClientRttNotify
)

type NetMsg struct {
	UserId         uint32
	EventId        uint16
	ApiId          uint16
	HeadMessage    *PacketHead
	PayloadMessage any
}
