package net

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"reflect"
)

func (p *ProtoEnDecode) initMsgProtoMap() {
	// apiId -> protoObj
	p.apiIdProtoObjMap[api.ApiGetPlayerTokenReq] = reflect.TypeOf(&proto.GetPlayerTokenReq{})               // 获取玩家token请求
	p.apiIdProtoObjMap[api.ApiPlayerLoginReq] = reflect.TypeOf(&proto.PlayerLoginReq{})                     // 玩家登录请求
	p.apiIdProtoObjMap[api.ApiPingReq] = reflect.TypeOf(&proto.PingReq{})                                   // ping请求
	p.apiIdProtoObjMap[api.ApiPlayerSetPauseReq] = reflect.TypeOf(&proto.PlayerSetPauseReq{})               // 玩家暂停请求
	p.apiIdProtoObjMap[api.ApiSetPlayerBornDataReq] = reflect.TypeOf(&proto.SetPlayerBornDataReq{})         // 玩家设置初始信息请求
	p.apiIdProtoObjMap[api.ApiGetPlayerSocialDetailReq] = reflect.TypeOf(&proto.GetPlayerSocialDetailReq{}) // 获取玩家社区信息请求
	p.apiIdProtoObjMap[api.ApiEnterSceneReadyReq] = reflect.TypeOf(&proto.EnterSceneReadyReq{})             // 进入场景准备就绪请求
	p.apiIdProtoObjMap[api.ApiGetScenePointReq] = reflect.TypeOf(&proto.GetScenePointReq{})                 // 获取场景信息请求
	p.apiIdProtoObjMap[api.ApiGetSceneAreaReq] = reflect.TypeOf(&proto.GetSceneAreaReq{})                   // 获取场景区域请求
	p.apiIdProtoObjMap[api.ApiEnterWorldAreaReq] = reflect.TypeOf(&proto.EnterWorldAreaReq{})               // 进入世界区域请求
	p.apiIdProtoObjMap[api.ApiUnionCmdNotify] = reflect.TypeOf(&proto.UnionCmdNotify{})                     // 聚合消息
	p.apiIdProtoObjMap[api.ApiSceneTransToPointReq] = reflect.TypeOf(&proto.SceneTransToPointReq{})         // 场景传送点请求
	p.apiIdProtoObjMap[api.ApiCombatInvocationsNotify] = reflect.TypeOf(&proto.CombatInvocationsNotify{})   // 战斗调用通知
	p.apiIdProtoObjMap[api.ApiMarkMapReq] = reflect.TypeOf(&proto.MarkMapReq{})                             // 标记地图请求
	p.apiIdProtoObjMap[api.ApiChangeAvatarReq] = reflect.TypeOf(&proto.ChangeAvatarReq{})                   // 更换角色请求
	p.apiIdProtoObjMap[api.ApiSetUpAvatarTeamReq] = reflect.TypeOf(&proto.SetUpAvatarTeamReq{})             // 配置队伍请求
	p.apiIdProtoObjMap[api.ApiChooseCurAvatarTeamReq] = reflect.TypeOf(&proto.ChooseCurAvatarTeamReq{})     // 切换队伍请求
	// protoObj -> apiId
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetPlayerTokenRsp{})] = api.ApiGetPlayerTokenRsp                           // 获取玩家token响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerLoginRsp{})] = api.ApiPlayerLoginRsp                                 // 玩家登录响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PingRsp{})] = api.ApiPingRsp                                               // ping响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerDataNotify{})] = api.ApiPlayerDataNotify                             // 玩家信息通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.StoreWeightLimitNotify{})] = api.ApiStoreWeightLimitNotify                 // 通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerStoreNotify{})] = api.ApiPlayerStoreNotify                           // 通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AvatarDataNotify{})] = api.ApiAvatarDataNotify                             // 角色信息通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerEnterSceneNotify{})] = api.ApiPlayerEnterSceneNotify                 // 玩家进入场景通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.OpenStateUpdateNotify{})] = api.ApiOpenStateUpdateNotify                   // 通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetPlayerSocialDetailRsp{})] = api.ApiGetPlayerSocialDetailRsp             // 获取玩家社区信息响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.EnterScenePeerNotify{})] = api.ApiEnterScenePeerNotify                     // 进入场景对方通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.EnterSceneReadyRsp{})] = api.ApiEnterSceneReadyRsp                         // 进入场景准备就绪响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetScenePointRsp{})] = api.ApiGetScenePointRsp                             // 获取场景信息响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetSceneAreaRsp{})] = api.ApiGetSceneAreaRsp                               // 获取场景区域响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ServerTimeNotify{})] = api.ApiServerTimeNotify                             // 服务器时间通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.WorldPlayerInfoNotify{})] = api.ApiWorldPlayerInfoNotify                   // 世界玩家信息通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.WorldDataNotify{})] = api.ApiWorldDataNotify                               // 世界数据通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerWorldSceneInfoListNotify{})] = api.ApiPlayerWorldSceneInfoListNotify // 场景解锁信息通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.HostPlayerNotify{})] = api.ApiHostPlayerNotify                             // 主机玩家通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneTimeNotify{})] = api.ApiSceneTimeNotify                               // 场景时间通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerGameTimeNotify{})] = api.ApiPlayerGameTimeNotify                     // 玩家游戏内时间通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerEnterSceneInfoNotify{})] = api.ApiPlayerEnterSceneInfoNotify         // 玩家进入场景信息通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneAreaWeatherNotify{})] = api.ApiSceneAreaWeatherNotify                 // 场景区域天气通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ScenePlayerInfoNotify{})] = api.ApiScenePlayerInfoNotify                   // 场景玩家信息通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneTeamUpdateNotify{})] = api.ApiSceneTeamUpdateNotify                   // 场景队伍更新通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SyncTeamEntityNotify{})] = api.ApiSyncTeamEntityNotify                     // 同步队伍实体通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SyncScenePlayTeamEntityNotify{})] = api.ApiSyncScenePlayTeamEntityNotify   // 同步场景玩家队伍实体通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneInitFinishRsp{})] = api.ApiSceneInitFinishRsp                         // 场景初始化完成响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.EnterSceneDoneRsp{})] = api.ApiEnterSceneDoneRsp                           // 进入场景完成响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerTimeNotify{})] = api.ApiPlayerTimeNotify                             // 玩家对时通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneEntityAppearNotify{})] = api.ApiSceneEntityAppearNotify               // 场景实体出现通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.WorldPlayerLocationNotify{})] = api.ApiWorldPlayerLocationNotify           // 世界玩家位置通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ScenePlayerLocationNotify{})] = api.ApiScenePlayerLocationNotify           // 场景玩家位置通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.WorldPlayerRTTNotify{})] = api.ApiWorldPlayerRTTNotify                     // 世界玩家RTT时延
	p.protoObjApiIdMap[reflect.TypeOf(&proto.EnterWorldAreaRsp{})] = api.ApiEnterWorldAreaRsp                           // 进入世界区域响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PostEnterSceneRsp{})] = api.ApiPostEnterSceneRsp                           // 提交进入场景响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.TowerAllDataRsp{})] = api.ApiTowerAllDataRsp                               // 深渊数据响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneTransToPointRsp{})] = api.ApiSceneTransToPointRsp                     // 场景传送点响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SceneEntityDisappearNotify{})] = api.ApiSceneEntityDisappearNotify         // 场景实体消失通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ChangeAvatarRsp{})] = api.ApiChangeAvatarRsp                               // 更换角色响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SetUpAvatarTeamRsp{})] = api.ApiSetUpAvatarTeamRsp                         // 配置队伍响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AvatarTeamUpdateNotify{})] = api.ApiAvatarTeamUpdateNotify                 // 角色队伍更新通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ChooseCurAvatarTeamRsp{})] = api.ApiChooseCurAvatarTeamRsp                 // 切换队伍响应
	// bypass 尚未得知协议的客户端上行消息
	p.bypassApiMap[api.ApiPathfindingEnterSceneReq] = true // 寻路进入场景请求
	p.bypassApiMap[api.ApiSceneInitFinishReq] = true       // 场景初始化完成请求
	p.bypassApiMap[api.ApiEnterSceneDoneReq] = true        // 进入场景完成请求
	p.bypassApiMap[api.ApiPostEnterSceneReq] = true        // 提交进入场景请求
	p.bypassApiMap[api.ApiTowerAllDataReq] = true          // 深渊数据请求
}

func (p *ProtoEnDecode) getProtoObjByApiId(apiId uint16) (protoObj any) {
	protoObjTypePointer, ok := p.apiIdProtoObjMap[apiId]
	if !ok {
		p.log.Error("unknown api id: %v", apiId)
		protoObj = nil
		return protoObj
	}
	protoObjInst := reflect.New(protoObjTypePointer.Elem())
	protoObj = protoObjInst.Interface()
	return protoObj
}

func (p *ProtoEnDecode) getApiIdByProtoObj(protoObj any) (apiId uint16) {
	var ok = false
	apiId, ok = p.protoObjApiIdMap[reflect.TypeOf(protoObj)]
	if !ok {
		p.log.Error("unknown proto object: %v", protoObj)
		apiId = 0
	}
	return apiId
}
