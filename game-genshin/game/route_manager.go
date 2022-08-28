package game

import (
	"flswld.com/gate-genshin-api/proto"
	"flswld.com/logger"
	pb "google.golang.org/protobuf/proto"
)

type HandlerFunc func(userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message)

type RouteManager struct {
	gameManager *GameManager
	// k:apiId v:HandlerFunc
	handlerFuncRouteMap map[uint16]HandlerFunc
}

func NewRouteManager(gameManager *GameManager) (r *RouteManager) {
	r = new(RouteManager)
	r.gameManager = gameManager
	r.handlerFuncRouteMap = make(map[uint16]HandlerFunc)
	return r
}

func (r *RouteManager) registerRouter(apiId uint16, handlerFunc HandlerFunc) {
	r.handlerFuncRouteMap[apiId] = handlerFunc
}

func (r *RouteManager) doRoute(apiId uint16, userId uint32, headMsg *proto.PacketHead, payloadMsg pb.Message) {
	handlerFunc, ok := r.handlerFuncRouteMap[apiId]
	if !ok {
		logger.LOG.Error("no route for msg, apiId: %v", apiId)
		return
	}
	handlerFunc(userId, headMsg, payloadMsg)
}

func (r *RouteManager) InitRoute() {
	r.registerRouter(proto.ApiPlayerSetPauseReq, r.gameManager.PlayerSetPauseReq)                         // 玩家暂停请求
	r.registerRouter(proto.ApiSetPlayerBornDataReq, r.gameManager.SetPlayerBornDataReq)                   // 玩家设置初始信息请求
	r.registerRouter(proto.ApiGetPlayerSocialDetailReq, r.gameManager.GetPlayerSocialDetailReq)           // 获取玩家社区信息请求
	r.registerRouter(proto.ApiEnterSceneReadyReq, r.gameManager.EnterSceneReadyReq)                       // 进入场景准备就绪请求
	r.registerRouter(proto.ApiPathfindingEnterSceneReq, r.gameManager.PathfindingEnterSceneReq)           // 寻路进入场景请求
	r.registerRouter(proto.ApiGetScenePointReq, r.gameManager.GetScenePointReq)                           // 获取场景信息请求
	r.registerRouter(proto.ApiGetSceneAreaReq, r.gameManager.GetSceneAreaReq)                             // 获取场景区域请求
	r.registerRouter(proto.ApiSceneInitFinishReq, r.gameManager.SceneInitFinishReq)                       // 场景初始化完成请求
	r.registerRouter(proto.ApiEnterSceneDoneReq, r.gameManager.EnterSceneDoneReq)                         // 进入场景完成请求
	r.registerRouter(proto.ApiEnterWorldAreaReq, r.gameManager.EnterWorldAreaReq)                         // 进入世界区域请求
	r.registerRouter(proto.ApiPostEnterSceneReq, r.gameManager.PostEnterSceneReq)                         // 提交进入场景请求
	r.registerRouter(proto.ApiTowerAllDataReq, r.gameManager.TowerAllDataReq)                             // 深渊数据请求
	r.registerRouter(proto.ApiSceneTransToPointReq, r.gameManager.SceneTransToPointReq)                   // 场景传送点请求
	r.registerRouter(proto.ApiCombatInvocationsNotify, r.gameManager.CombatInvocationsNotify)             // 战斗调用通知
	r.registerRouter(proto.ApiMarkMapReq, r.gameManager.MarkMapReq)                                       // 标记地图请求
	r.registerRouter(proto.ApiChangeAvatarReq, r.gameManager.ChangeAvatarReq)                             // 更换角色请求
	r.registerRouter(proto.ApiSetUpAvatarTeamReq, r.gameManager.SetUpAvatarTeamReq)                       // 配置队伍请求
	r.registerRouter(proto.ApiChooseCurAvatarTeamReq, r.gameManager.ChooseCurAvatarTeamReq)               // 切换队伍请求
	r.registerRouter(proto.ApiGetGachaInfoReq, r.gameManager.GetGachaInfoReq)                             // 卡池获取请求
	r.registerRouter(proto.ApiDoGachaReq, r.gameManager.DoGachaReq)                                       // 抽卡请求
	r.registerRouter(proto.ApiQueryPathReq, r.gameManager.QueryPathReq)                                   // 寻路请求
	r.registerRouter(proto.ApiAbilityInvocationsNotify, r.gameManager.AbilityInvocationsNotify)           // 技能使用通知
	r.registerRouter(proto.ApiClientAbilityInitFinishNotify, r.gameManager.ClientAbilityInitFinishNotify) // 客户端技能初始化完成通知
	r.registerRouter(proto.ApiEntityAiSyncNotify, r.gameManager.EntityAiSyncNotify)                       // 实体AI怪物同步通知
	r.registerRouter(proto.ApiWearEquipReq, r.gameManager.WearEquipReq)                                   // 装备穿戴请求
	r.registerRouter(proto.ApiChangeGameTimeReq, r.gameManager.ChangeGameTimeReq)                         // 改变游戏场景时间请求
	r.registerRouter(proto.ApiSetPlayerBirthdayReq, r.gameManager.SetPlayerBirthdayReq)                   // 设置生日请求
	r.registerRouter(proto.ApiSetNameCardReq, r.gameManager.SetNameCardReq)                               // 修改名片请求
	r.registerRouter(proto.ApiSetPlayerSignatureReq, r.gameManager.SetPlayerSignatureReq)                 // 修改签名请求
	r.registerRouter(proto.ApiSetPlayerNameReq, r.gameManager.SetPlayerNameReq)                           // 修改昵称请求
	r.registerRouter(proto.ApiSetPlayerHeadImageReq, r.gameManager.SetPlayerHeadImageReq)                 // 修改头像请求
	r.registerRouter(proto.ApiGetAllUnlockNameCardReq, r.gameManager.GetAllUnlockNameCardReq)             // 获取全部已解锁名片请求
	r.registerRouter(proto.ApiGetPlayerFriendListReq, r.gameManager.GetPlayerFriendListReq)               // 好友列表请求
	r.registerRouter(proto.ApiGetPlayerAskFriendListReq, r.gameManager.GetPlayerAskFriendListReq)         // 好友申请列表请求
	r.registerRouter(proto.ApiAskAddFriendReq, r.gameManager.AskAddFriendReq)                             // 加好友请求
	r.registerRouter(proto.ApiDealAddFriendReq, r.gameManager.DealAddFriendReq)                           // 处理好友申请请求
	r.registerRouter(proto.ApiGetOnlinePlayerListReq, r.gameManager.GetOnlinePlayerListReq)               // 在线玩家列表请求
	r.registerRouter(proto.ApiPlayerForceExitReq, r.gameManager.PlayerForceExitReq)                       // 退出游戏请求
	r.registerRouter(proto.ApiPlayerApplyEnterMpReq, r.gameManager.PlayerApplyEnterMpReq)                 // 世界敲门请求
	r.registerRouter(proto.ApiPlayerApplyEnterMpResultReq, r.gameManager.PlayerApplyEnterMpResultReq)     // 世界敲门处理请求
	r.registerRouter(proto.ApiPlayerGetForceQuitBanInfoReq, r.gameManager.PlayerGetForceQuitBanInfoReq)   // 退出世界请求
}

func (r *RouteManager) RouteHandle(netMsg *proto.NetMsg) {
	switch netMsg.EventId {
	case proto.NormalMsg:
		r.doRoute(netMsg.ApiId, netMsg.UserId, netMsg.HeadMessage, netMsg.PayloadMessage)
	case proto.UserLogin:
		r.gameManager.OnLogin(netMsg.UserId)
	case proto.UserOffline:
		r.gameManager.OnUserOffline(netMsg.UserId)
	case proto.ClientRttNotify:
		r.gameManager.ClientRttNotify(netMsg.UserId, netMsg.ClientRtt)
	case proto.ClientTimeNotify:
		r.gameManager.ClientTimeNotify(netMsg.UserId, netMsg.ClientTime)
	}
}
