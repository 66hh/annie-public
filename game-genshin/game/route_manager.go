package game

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/logger"
)

type HandlerFunc func(userId uint32, headMsg *api.PacketHead, payloadMsg any)

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

func (r *RouteManager) doRoute(apiId uint16, userId uint32, headMsg *api.PacketHead, payloadMsg any) {
	handlerFunc, ok := r.handlerFuncRouteMap[apiId]
	if !ok {
		logger.LOG.Error("no route for msg, apiId: %v", apiId)
		return
	}
	handlerFunc(userId, headMsg, payloadMsg)
}

func (r *RouteManager) InitRoute() {
	r.registerRouter(api.ApiPlayerSetPauseReq, r.gameManager.PlayerSetPauseReq)                         // 玩家暂停请求
	r.registerRouter(api.ApiSetPlayerBornDataReq, r.gameManager.SetPlayerBornDataReq)                   // 玩家设置初始信息请求
	r.registerRouter(api.ApiGetPlayerSocialDetailReq, r.gameManager.GetPlayerSocialDetailReq)           // 获取玩家社区信息请求
	r.registerRouter(api.ApiEnterSceneReadyReq, r.gameManager.EnterSceneReadyReq)                       // 进入场景准备就绪请求
	r.registerRouter(api.ApiPathfindingEnterSceneReq, r.gameManager.PathfindingEnterSceneReq)           // 寻路进入场景请求
	r.registerRouter(api.ApiGetScenePointReq, r.gameManager.GetScenePointReq)                           // 获取场景信息请求
	r.registerRouter(api.ApiGetSceneAreaReq, r.gameManager.GetSceneAreaReq)                             // 获取场景区域请求
	r.registerRouter(api.ApiSceneInitFinishReq, r.gameManager.SceneInitFinishReq)                       // 场景初始化完成请求
	r.registerRouter(api.ApiEnterSceneDoneReq, r.gameManager.EnterSceneDoneReq)                         // 进入场景完成请求
	r.registerRouter(api.ApiEnterWorldAreaReq, r.gameManager.EnterWorldAreaReq)                         // 进入世界区域请求
	r.registerRouter(api.ApiPostEnterSceneReq, r.gameManager.PostEnterSceneReq)                         // 提交进入场景请求
	r.registerRouter(api.ApiTowerAllDataReq, r.gameManager.TowerAllDataReq)                             // 深渊数据请求
	r.registerRouter(api.ApiSceneTransToPointReq, r.gameManager.SceneTransToPointReq)                   // 场景传送点请求
	r.registerRouter(api.ApiCombatInvocationsNotify, r.gameManager.CombatInvocationsNotify)             // 战斗调用通知
	r.registerRouter(api.ApiMarkMapReq, r.gameManager.MarkMapReq)                                       // 标记地图请求
	r.registerRouter(api.ApiChangeAvatarReq, r.gameManager.ChangeAvatarReq)                             // 更换角色请求
	r.registerRouter(api.ApiSetUpAvatarTeamReq, r.gameManager.SetUpAvatarTeamReq)                       // 配置队伍请求
	r.registerRouter(api.ApiChooseCurAvatarTeamReq, r.gameManager.ChooseCurAvatarTeamReq)               // 切换队伍请求
	r.registerRouter(api.ApiGetGachaInfoReq, r.gameManager.GetGachaInfoReq)                             // 卡池获取请求
	r.registerRouter(api.ApiDoGachaReq, r.gameManager.DoGachaReq)                                       // 抽卡请求
	r.registerRouter(api.ApiQueryPathReq, r.gameManager.QueryPathReq)                                   // 寻路请求
	r.registerRouter(api.ApiPingReq, r.gameManager.PingReq)                                             // ping请求
	r.registerRouter(api.ApiAbilityInvocationsNotify, r.gameManager.AbilityInvocationsNotify)           // 技能使用通知
	r.registerRouter(api.ApiClientAbilityInitFinishNotify, r.gameManager.ClientAbilityInitFinishNotify) // 客户端技能初始化完成通知
	r.registerRouter(api.ApiEntityAiSyncNotify, r.gameManager.EntityAiSyncNotify)                       // 实体AI怪物同步通知
	r.registerRouter(api.ApiWearEquipReq, r.gameManager.WearEquipReq)                                   // 装备穿戴请求
	r.registerRouter(api.ApiChangeGameTimeReq, r.gameManager.ChangeGameTimeReq)                         // 改变游戏场景时间请求
	r.registerRouter(api.ApiSetPlayerBirthdayReq, r.gameManager.SetPlayerBirthdayReq)                   // 设置生日请求
	r.registerRouter(api.ApiSetNameCardReq, r.gameManager.SetNameCardReq)                               // 修改名片请求
	r.registerRouter(api.ApiSetPlayerSignatureReq, r.gameManager.SetPlayerSignatureReq)                 // 修改签名请求
	r.registerRouter(api.ApiSetPlayerNameReq, r.gameManager.SetPlayerNameReq)                           // 修改昵称请求
	r.registerRouter(api.ApiSetPlayerHeadImageReq, r.gameManager.SetPlayerHeadImageReq)                 // 修改头像请求
	r.registerRouter(api.ApiGetAllUnlockNameCardReq, r.gameManager.GetAllUnlockNameCardReq)             // 获取全部已解锁名片请求
	r.registerRouter(api.ApiGetPlayerFriendListReq, r.gameManager.GetPlayerFriendListReq)               // 好友列表请求
	r.registerRouter(api.ApiGetPlayerAskFriendListReq, r.gameManager.GetPlayerAskFriendListReq)         // 好友申请列表请求
	r.registerRouter(api.ApiAskAddFriendReq, r.gameManager.AskAddFriendReq)                             // 加好友请求
	r.registerRouter(api.ApiDealAddFriendReq, r.gameManager.DealAddFriendReq)                           // 处理好友申请请求
	r.registerRouter(api.ApiGetOnlinePlayerListReq, r.gameManager.GetOnlinePlayerListReq)               // 在线玩家列表请求
	r.registerRouter(api.ApiPlayerForceExitReq, r.gameManager.PlayerForceExitReq)                       // 退出游戏请求
	r.registerRouter(api.ApiPlayerApplyEnterMpReq, r.gameManager.PlayerApplyEnterMpReq)                 // 世界敲门请求
	r.registerRouter(api.ApiPlayerApplyEnterMpResultReq, r.gameManager.PlayerApplyEnterMpResultReq)     // 世界敲门处理请求
	r.registerRouter(api.ApiPlayerGetForceQuitBanInfoReq, r.gameManager.PlayerGetForceQuitBanInfoReq)   // 退出世界请求
}

func (r *RouteManager) RouteHandle(netMsg *api.NetMsg) {
	switch netMsg.EventId {
	case api.NormalMsg:
		r.doRoute(netMsg.ApiId, netMsg.UserId, netMsg.HeadMessage, netMsg.PayloadMessage)
	case api.UserLogin:
		r.gameManager.OnLogin(netMsg.UserId)
	case api.UserOffline:
		r.gameManager.OnUserOffline(netMsg.UserId)
	case api.ClientRttNotify:
		r.gameManager.ClientRttNotify(netMsg.UserId, netMsg.PayloadMessage)
	}
}
