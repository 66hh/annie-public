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
	r.registerRouter(api.ApiPlayerSetPauseReq, r.gameManager.PlayerSetPauseReq)               // 玩家暂停请求
	r.registerRouter(api.ApiSetPlayerBornDataReq, r.gameManager.SetPlayerBornDataReq)         // 玩家设置初始信息请求
	r.registerRouter(api.ApiGetPlayerSocialDetailReq, r.gameManager.GetPlayerSocialDetailReq) // 获取玩家社区信息请求
	r.registerRouter(api.ApiEnterSceneReadyReq, r.gameManager.EnterSceneReadyReq)             // 进入场景准备就绪请求
	r.registerRouter(api.ApiPathfindingEnterSceneReq, r.gameManager.PathfindingEnterSceneReq) // 寻路进入场景请求
	r.registerRouter(api.ApiGetScenePointReq, r.gameManager.GetScenePointReq)                 // 获取场景信息请求
	r.registerRouter(api.ApiGetSceneAreaReq, r.gameManager.GetSceneAreaReq)                   // 获取场景区域请求
	r.registerRouter(api.ApiSceneInitFinishReq, r.gameManager.SceneInitFinishReq)             // 场景初始化完成请求
	r.registerRouter(api.ApiEnterSceneDoneReq, r.gameManager.EnterSceneDoneReq)               // 进入场景完成请求
	r.registerRouter(api.ApiEnterWorldAreaReq, r.gameManager.EnterWorldAreaReq)               // 进入世界区域请求
	r.registerRouter(api.ApiPostEnterSceneReq, r.gameManager.PostEnterSceneReq)               // 提交进入场景请求
	r.registerRouter(api.ApiTowerAllDataReq, r.gameManager.TowerAllDataReq)                   // 深渊数据请求
	r.registerRouter(api.ApiSceneTransToPointReq, r.gameManager.SceneTransToPointReq)         // 场景传送点请求
	r.registerRouter(api.ApiCombatInvocationsNotify, r.gameManager.CombatInvocationsNotify)   // 战斗调用通知
	r.registerRouter(api.ApiMarkMapReq, r.gameManager.MarkMapReq)                             // 标记地图请求
	r.registerRouter(api.ApiChangeAvatarReq, r.gameManager.ChangeAvatarReq)                   // 更换角色请求
	r.registerRouter(api.ApiSetUpAvatarTeamReq, r.gameManager.SetUpAvatarTeamReq)             // 配置队伍请求
	r.registerRouter(api.ApiChooseCurAvatarTeamReq, r.gameManager.ChooseCurAvatarTeamReq)     // 切换队伍请求
	r.registerRouter(api.ApiGetGachaInfoReq, r.gameManager.GetGachaInfoReq)                   // 卡池获取请求
	r.registerRouter(api.ApiDoGachaReq, r.gameManager.DoGachaReq)                             // 抽卡请求
}

func (r *RouteManager) RouteHandle(netMsg *api.NetMsg) {
	switch netMsg.EventId {
	case api.NormalMsg:
		r.doRoute(netMsg.ApiId, netMsg.UserId, netMsg.HeadMessage, netMsg.PayloadMessage)
	case api.UserLogin:
		r.gameManager.OnLogin(netMsg.UserId)
	case api.UserOffline:
		r.gameManager.OnUserOffline(netMsg.UserId)
	}
}
