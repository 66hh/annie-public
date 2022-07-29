package constant

type EnterReason struct {
	None                    uint16
	Login                   uint16
	DungeonReplay           uint16
	DungeonReviveOnWaypoint uint16
	DungeonEnter            uint16
	DungeonQuit             uint16
	Gm                      uint16
	QuestRollback           uint16
	Revival                 uint16
	PersonalScene           uint16
	TransPoint              uint16
	ClientTransmit          uint16
	ForceDragBack           uint16
	TeamKick                uint16
	TeamJoin                uint16
	TeamBack                uint16
	Muip                    uint16
	DungeonInviteAccept     uint16
	Lua                     uint16
	ActivityLoadTerrain     uint16
	HostFromSingleToMp      uint16
	MpPlay                  uint16
	AnchorPoint             uint16
	LuaSkipUi               uint16
	ReloadTerrain           uint16
	DraftTransfer           uint16
	EnterHome               uint16
	ExitHome                uint16
	ChangeHomeModule        uint16
	Gallery                 uint16
	HomeSceneJump           uint16
	HideAndSeek             uint16
}

func GetEnterReasonConst() (r *EnterReason) {
	r = new(EnterReason)
	r.None = 0
	r.Login = 1
	r.DungeonReplay = 11
	r.DungeonReviveOnWaypoint = 12
	r.DungeonEnter = 13
	r.DungeonQuit = 14
	r.Gm = 21
	r.QuestRollback = 31
	r.Revival = 32
	r.PersonalScene = 41
	r.TransPoint = 42
	r.ClientTransmit = 43
	r.ForceDragBack = 44
	r.TeamKick = 51
	r.TeamJoin = 52
	r.TeamBack = 53
	r.Muip = 54
	r.DungeonInviteAccept = 55
	r.Lua = 56
	r.ActivityLoadTerrain = 57
	r.HostFromSingleToMp = 58
	r.MpPlay = 59
	r.AnchorPoint = 60
	r.LuaSkipUi = 61
	r.ReloadTerrain = 62
	r.DraftTransfer = 63
	r.EnterHome = 64
	r.ExitHome = 65
	r.ChangeHomeModule = 66
	r.Gallery = 67
	r.HomeSceneJump = 68
	r.HideAndSeek = 69
	return r
}
