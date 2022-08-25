package net

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	"reflect"
)

func (p *ProtoEnDecode) initMsgProtoMap() {
	// apiId -> protoObj
	p.apiIdProtoObjMap[api.ApiGetPlayerTokenReq] = reflect.TypeOf(&proto.GetPlayerTokenReq{})                         // 获取玩家token请求
	p.apiIdProtoObjMap[api.ApiPlayerLoginReq] = reflect.TypeOf(&proto.PlayerLoginReq{})                               // 玩家登录请求
	p.apiIdProtoObjMap[api.ApiPingReq] = reflect.TypeOf(&proto.PingReq{})                                             // ping请求
	p.apiIdProtoObjMap[api.ApiPlayerSetPauseReq] = reflect.TypeOf(&proto.PlayerSetPauseReq{})                         // 玩家暂停请求
	p.apiIdProtoObjMap[api.ApiSetPlayerBornDataReq] = reflect.TypeOf(&proto.SetPlayerBornDataReq{})                   // 注册请求
	p.apiIdProtoObjMap[api.ApiGetPlayerSocialDetailReq] = reflect.TypeOf(&proto.GetPlayerSocialDetailReq{})           // 获取玩家社区信息请求
	p.apiIdProtoObjMap[api.ApiEnterSceneReadyReq] = reflect.TypeOf(&proto.EnterSceneReadyReq{})                       // 进入场景准备就绪请求
	p.apiIdProtoObjMap[api.ApiGetScenePointReq] = reflect.TypeOf(&proto.GetScenePointReq{})                           // 获取场景信息请求
	p.apiIdProtoObjMap[api.ApiGetSceneAreaReq] = reflect.TypeOf(&proto.GetSceneAreaReq{})                             // 获取场景区域请求
	p.apiIdProtoObjMap[api.ApiEnterWorldAreaReq] = reflect.TypeOf(&proto.EnterWorldAreaReq{})                         // 进入世界区域请求
	p.apiIdProtoObjMap[api.ApiUnionCmdNotify] = reflect.TypeOf(&proto.UnionCmdNotify{})                               // 聚合消息
	p.apiIdProtoObjMap[api.ApiSceneTransToPointReq] = reflect.TypeOf(&proto.SceneTransToPointReq{})                   // 场景传送点请求
	p.apiIdProtoObjMap[api.ApiMarkMapReq] = reflect.TypeOf(&proto.MarkMapReq{})                                       // 标记地图请求
	p.apiIdProtoObjMap[api.ApiChangeAvatarReq] = reflect.TypeOf(&proto.ChangeAvatarReq{})                             // 更换角色请求
	p.apiIdProtoObjMap[api.ApiSetUpAvatarTeamReq] = reflect.TypeOf(&proto.SetUpAvatarTeamReq{})                       // 配置队伍请求
	p.apiIdProtoObjMap[api.ApiChooseCurAvatarTeamReq] = reflect.TypeOf(&proto.ChooseCurAvatarTeamReq{})               // 切换队伍请求
	p.apiIdProtoObjMap[api.ApiDoGachaReq] = reflect.TypeOf(&proto.DoGachaReq{})                                       // 抽卡请求
	p.apiIdProtoObjMap[api.ApiQueryPathReq] = reflect.TypeOf(&proto.QueryPathReq{})                                   // 寻路请求
	p.apiIdProtoObjMap[api.ApiCombatInvocationsNotify] = reflect.TypeOf(&proto.CombatInvocationsNotify{})             // 战斗调用通知
	p.apiIdProtoObjMap[api.ApiAbilityInvocationsNotify] = reflect.TypeOf(&proto.AbilityInvocationsNotify{})           // 技能使用通知
	p.apiIdProtoObjMap[api.ApiClientAbilityInitFinishNotify] = reflect.TypeOf(&proto.ClientAbilityInitFinishNotify{}) // 客户端技能初始化完成通知
	p.apiIdProtoObjMap[api.ApiEntityAiSyncNotify] = reflect.TypeOf(&proto.EntityAiSyncNotify{})                       // 实体AI怪物同步通知
	p.apiIdProtoObjMap[api.ApiWearEquipReq] = reflect.TypeOf(&proto.WearEquipReq{})                                   // 装备穿戴请求
	p.apiIdProtoObjMap[api.ApiChangeGameTimeReq] = reflect.TypeOf(&proto.ChangeGameTimeReq{})                         // 改变游戏场景时间请求
	p.apiIdProtoObjMap[api.ApiSetPlayerBirthdayReq] = reflect.TypeOf(&proto.SetPlayerBirthdayReq{})                   // 设置生日请求
	p.apiIdProtoObjMap[api.ApiSetNameCardReq] = reflect.TypeOf(&proto.SetNameCardReq{})                               // 修改名片请求
	p.apiIdProtoObjMap[api.ApiSetPlayerSignatureReq] = reflect.TypeOf(&proto.SetPlayerSignatureReq{})                 // 修改签名请求
	p.apiIdProtoObjMap[api.ApiSetPlayerNameReq] = reflect.TypeOf(&proto.SetPlayerNameReq{})                           // 修改昵称请求
	p.apiIdProtoObjMap[api.ApiSetPlayerHeadImageReq] = reflect.TypeOf(&proto.SetPlayerHeadImageReq{})                 // 修改头像请求
	p.apiIdProtoObjMap[api.ApiAskAddFriendReq] = reflect.TypeOf(&proto.AskAddFriendReq{})                             // 加好友请求
	p.apiIdProtoObjMap[api.ApiDealAddFriendReq] = reflect.TypeOf(&proto.DealAddFriendReq{})                           // 处理好友申请请求
	p.apiIdProtoObjMap[api.ApiGetOnlinePlayerListReq] = reflect.TypeOf(&proto.GetOnlinePlayerListReq{})               // 在线玩家列表请求
	p.apiIdProtoObjMap[api.ApiPathfindingEnterSceneReq] = reflect.TypeOf(&proto.NullMsg{})                            // 寻路进入场景请求
	p.apiIdProtoObjMap[api.ApiSceneInitFinishReq] = reflect.TypeOf(&proto.NullMsg{})                                  // 场景初始化完成请求
	p.apiIdProtoObjMap[api.ApiEnterSceneDoneReq] = reflect.TypeOf(&proto.NullMsg{})                                   // 进入场景完成请求
	p.apiIdProtoObjMap[api.ApiPostEnterSceneReq] = reflect.TypeOf(&proto.NullMsg{})                                   // 提交进入场景请求
	p.apiIdProtoObjMap[api.ApiTowerAllDataReq] = reflect.TypeOf(&proto.NullMsg{})                                     // 深渊数据请求
	p.apiIdProtoObjMap[api.ApiGetGachaInfoReq] = reflect.TypeOf(&proto.NullMsg{})                                     // 卡池获取请求
	p.apiIdProtoObjMap[api.ApiGetAllUnlockNameCardReq] = reflect.TypeOf(&proto.NullMsg{})                             // 获取全部已解锁名片请求
	p.apiIdProtoObjMap[api.ApiGetPlayerFriendListReq] = reflect.TypeOf(&proto.NullMsg{})                              // 好友列表请求
	p.apiIdProtoObjMap[api.ApiGetPlayerAskFriendListReq] = reflect.TypeOf(&proto.NullMsg{})                           // 好友申请列表请求
	p.apiIdProtoObjMap[api.ApiPlayerForceExitReq] = reflect.TypeOf(&proto.NullMsg{})                                  // 退出游戏请求
	p.apiIdProtoObjMap[api.ApiPlayerApplyEnterMpReq] = reflect.TypeOf(&proto.PlayerApplyEnterMpReq{})                 // 世界敲门请求
	p.apiIdProtoObjMap[api.ApiPlayerApplyEnterMpResultReq] = reflect.TypeOf(&proto.PlayerApplyEnterMpResultReq{})     // 世界敲门处理请求
	p.apiIdProtoObjMap[api.ApiPlayerGetForceQuitBanInfoReq] = reflect.TypeOf(&proto.NullMsg{})                        // 退出世界请求
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
	p.protoObjApiIdMap[reflect.TypeOf(&proto.StoreItemChangeNotify{})] = api.ApiStoreItemChangeNotify                   // 背包道具变动通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ItemAddHintNotify{})] = api.ApiItemAddHintNotify                           // 道具增加提示通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.StoreItemDelNotify{})] = api.ApiStoreItemDelNotify                         // 背包道具删除通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerPropNotify{})] = api.ApiPlayerPropNotify                             // 玩家属性通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetGachaInfoRsp{})] = api.ApiGetGachaInfoRsp                               // 卡池获取响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.DoGachaRsp{})] = api.ApiDoGachaRsp                                         // 抽卡响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.EntityFightPropUpdateNotify{})] = api.ApiEntityFightPropUpdateNotify       // 实体战斗属性更新通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.QueryPathRsp{})] = api.ApiQueryPathRsp                                     // 寻路响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.EntityAiSyncNotify{})] = api.ApiEntityAiSyncNotify                         // 实体AI怪物同步通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AvatarFightPropNotify{})] = api.ApiAvatarFightPropNotify                   // 角色战斗属性通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AvatarEquipChangeNotify{})] = api.ApiAvatarEquipChangeNotify               // 角色装备改变通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AvatarAddNotify{})] = api.ApiAvatarAddNotify                               // 角色新增通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.WearEquipRsp{})] = api.ApiWearEquipRsp                                     // 装备穿戴响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ChangeGameTimeRsp{})] = api.ApiChangeGameTimeRsp                           // 改变游戏场景时间响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SetPlayerBirthdayRsp{})] = api.ApiSetPlayerBirthdayRsp                     // 设置生日响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SetNameCardRsp{})] = api.ApiSetNameCardRsp                                 // 修改名片响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SetPlayerSignatureRsp{})] = api.ApiSetPlayerSignatureRsp                   // 修改签名响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SetPlayerNameRsp{})] = api.ApiSetPlayerNameRsp                             // 修改昵称响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.SetPlayerHeadImageRsp{})] = api.ApiSetPlayerHeadImageRsp                   // 修改头像响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetAllUnlockNameCardRsp{})] = api.ApiGetAllUnlockNameCardRsp               // 获取全部已解锁名片响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.UnlockNameCardNotify{})] = api.ApiUnlockNameCardNotify                     // 名片解锁通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetPlayerFriendListRsp{})] = api.ApiGetPlayerFriendListRsp                 // 好友列表响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetPlayerAskFriendListRsp{})] = api.ApiGetPlayerAskFriendListRsp           // 好友申请列表响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AskAddFriendRsp{})] = api.ApiAskAddFriendRsp                               // 加好友响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AskAddFriendNotify{})] = api.ApiAskAddFriendNotify                         // 加好友通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.DealAddFriendRsp{})] = api.ApiDealAddFriendRsp                             // 处理好友申请响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.GetOnlinePlayerListRsp{})] = api.ApiGetOnlinePlayerListRsp                 // 在线玩家列表响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.NullMsg{})] = api.ApiSceneForceUnlockNotify                                // 场景强制解锁通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.NullMsg{})] = api.ApiSetPlayerBornDataRsp                                  // 注册响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.NullMsg{})] = api.ApiDoSetPlayerBornDataNotify                             // 注册通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.NullMsg{})] = api.ApiPathfindingEnterSceneRsp                              // 寻路进入场景响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.NullMsg{})] = api.ApiPlayerForceExitRsp                                    // 退出游戏响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.CombatInvocationsNotify{})] = api.ApiCombatInvocationsNotify               // 战斗调用通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.AbilityInvocationsNotify{})] = api.ApiAbilityInvocationsNotify             // 技能使用通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.ClientAbilityInitFinishNotify{})] = api.ApiClientAbilityInitFinishNotify   // 客户端技能初始化完成通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.DelTeamEntityNotify{})] = api.ApiDelTeamEntityNotify                       // 删除队伍实体通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerApplyEnterMpRsp{})] = api.ApiPlayerApplyEnterMpRsp                   // 世界敲门响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerApplyEnterMpNotify{})] = api.ApiPlayerApplyEnterMpNotify             // 世界敲门通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerApplyEnterMpResultRsp{})] = api.ApiPlayerApplyEnterMpResultRsp       // 世界敲门处理响应
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerApplyEnterMpResultNotify{})] = api.ApiPlayerApplyEnterMpResultNotify // 世界敲门处理通知
	p.protoObjApiIdMap[reflect.TypeOf(&proto.PlayerGetForceQuitBanInfoRsp{})] = api.ApiPlayerGetForceQuitBanInfoRsp     // 退出世界响应
	// 消息体为空沙比gob无法序列化的消息
	p.bypassApiMap[api.ApiGetOnlinePlayerListReq] = true // 未知
	// 尚未得知的客户端上行消息
	p.bypassApiMap[api.ApiClientAbilityChangeNotify] = true       // 未知
	p.bypassApiMap[api.ApiEntityConfigHashNotify] = true          // 未知
	p.bypassApiMap[api.ApiMonsterAIConfigHashNotify] = true       // 未知
	p.bypassApiMap[api.ApiEvtAiSyncCombatThreatInfoNotify] = true // 未知
	p.bypassApiMap[api.ApiEvtAiSyncSkillCdNotify] = true          // 未知
	p.bypassApiMap[api.ApiGetRegionSearchReq] = true              // 未知
	p.bypassApiMap[api.ApiObstacleModifyNotify] = true            // 未知
}

func (p *ProtoEnDecode) getProtoObjByApiId(apiId uint16) (protoObj any) {
	protoObjTypePointer, ok := p.apiIdProtoObjMap[apiId]
	if !ok {
		logger.LOG.Error("unknown api id: %v", apiId)
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
		logger.LOG.Error("unknown proto object: %v", protoObj)
		apiId = 0
	}
	return apiId
}
