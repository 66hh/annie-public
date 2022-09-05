package proto

import (
	"flswld.com/logger"
	pb "google.golang.org/protobuf/proto"
	"reflect"
)

type ApiProtoMap struct {
	apiIdProtoObjMap map[uint16]reflect.Type
	protoObjApiIdMap map[reflect.Type]uint16
	apiDeDupMap      map[uint16]bool
}

func NewApiProtoMap() (r *ApiProtoMap) {
	r = new(ApiProtoMap)
	r.apiIdProtoObjMap = make(map[uint16]reflect.Type)
	r.protoObjApiIdMap = make(map[reflect.Type]uint16)
	r.apiDeDupMap = make(map[uint16]bool)
	r.registerAllMessage()
	return r
}

func (a *ApiProtoMap) registerAllMessage() {
	// 已接入的消息
	a.registerMessage(ApiGetPlayerTokenReq, &GetPlayerTokenReq{})                           // 获取玩家token请求
	a.registerMessage(ApiPlayerLoginReq, &PlayerLoginReq{})                                 // 玩家登录请求
	a.registerMessage(ApiPingReq, &PingReq{})                                               // ping请求
	a.registerMessage(ApiPlayerSetPauseReq, &PlayerSetPauseReq{})                           // 玩家暂停请求
	a.registerMessage(ApiSetPlayerBornDataReq, &SetPlayerBornDataReq{})                     // 注册请求
	a.registerMessage(ApiGetPlayerSocialDetailReq, &GetPlayerSocialDetailReq{})             // 获取玩家社区信息请求
	a.registerMessage(ApiEnterSceneReadyReq, &EnterSceneReadyReq{})                         // 进入场景准备就绪请求
	a.registerMessage(ApiGetScenePointReq, &GetScenePointReq{})                             // 获取场景信息请求
	a.registerMessage(ApiGetSceneAreaReq, &GetSceneAreaReq{})                               // 获取场景区域请求
	a.registerMessage(ApiEnterWorldAreaReq, &EnterWorldAreaReq{})                           // 进入世界区域请求
	a.registerMessage(ApiUnionCmdNotify, &UnionCmdNotify{})                                 // 聚合消息
	a.registerMessage(ApiSceneTransToPointReq, &SceneTransToPointReq{})                     // 场景传送点请求
	a.registerMessage(ApiMarkMapReq, &MarkMapReq{})                                         // 标记地图请求
	a.registerMessage(ApiChangeAvatarReq, &ChangeAvatarReq{})                               // 更换角色请求
	a.registerMessage(ApiSetUpAvatarTeamReq, &SetUpAvatarTeamReq{})                         // 配置队伍请求
	a.registerMessage(ApiChooseCurAvatarTeamReq, &ChooseCurAvatarTeamReq{})                 // 切换队伍请求
	a.registerMessage(ApiDoGachaReq, &DoGachaReq{})                                         // 抽卡请求
	a.registerMessage(ApiQueryPathReq, &QueryPathReq{})                                     // 寻路请求
	a.registerMessage(ApiCombatInvocationsNotify, &CombatInvocationsNotify{})               // 战斗调用通知
	a.registerMessage(ApiAbilityInvocationsNotify, &AbilityInvocationsNotify{})             // 技能使用通知
	a.registerMessage(ApiClientAbilityInitFinishNotify, &ClientAbilityInitFinishNotify{})   // 客户端技能初始化完成通知
	a.registerMessage(ApiEntityAiSyncNotify, &EntityAiSyncNotify{})                         // 实体AI怪物同步通知
	a.registerMessage(ApiWearEquipReq, &WearEquipReq{})                                     // 装备穿戴请求
	a.registerMessage(ApiChangeGameTimeReq, &ChangeGameTimeReq{})                           // 改变游戏场景时间请求
	a.registerMessage(ApiSetPlayerBirthdayReq, &SetPlayerBirthdayReq{})                     // 设置生日请求
	a.registerMessage(ApiSetNameCardReq, &SetNameCardReq{})                                 // 修改名片请求
	a.registerMessage(ApiSetPlayerSignatureReq, &SetPlayerSignatureReq{})                   // 修改签名请求
	a.registerMessage(ApiSetPlayerNameReq, &SetPlayerNameReq{})                             // 修改昵称请求
	a.registerMessage(ApiSetPlayerHeadImageReq, &SetPlayerHeadImageReq{})                   // 修改头像请求
	a.registerMessage(ApiAskAddFriendReq, &AskAddFriendReq{})                               // 加好友请求
	a.registerMessage(ApiDealAddFriendReq, &DealAddFriendReq{})                             // 处理好友申请请求
	a.registerMessage(ApiGetOnlinePlayerListReq, &GetOnlinePlayerListReq{})                 // 在线玩家列表请求
	a.registerMessage(ApiPathfindingEnterSceneReq, &NullMsg{})                              // 寻路进入场景请求
	a.registerMessage(ApiSceneInitFinishReq, &NullMsg{})                                    // 场景初始化完成请求
	a.registerMessage(ApiEnterSceneDoneReq, &NullMsg{})                                     // 进入场景完成请求
	a.registerMessage(ApiPostEnterSceneReq, &NullMsg{})                                     // 提交进入场景请求
	a.registerMessage(ApiTowerAllDataReq, &NullMsg{})                                       // 深渊数据请求
	a.registerMessage(ApiGetGachaInfoReq, &NullMsg{})                                       // 卡池获取请求
	a.registerMessage(ApiGetAllUnlockNameCardReq, &NullMsg{})                               // 获取全部已解锁名片请求
	a.registerMessage(ApiGetPlayerFriendListReq, &NullMsg{})                                // 好友列表请求
	a.registerMessage(ApiGetPlayerAskFriendListReq, &NullMsg{})                             // 好友申请列表请求
	a.registerMessage(ApiPlayerForceExitReq, &NullMsg{})                                    // 退出游戏请求
	a.registerMessage(ApiPlayerApplyEnterMpReq, &PlayerApplyEnterMpReq{})                   // 世界敲门请求
	a.registerMessage(ApiPlayerApplyEnterMpResultReq, &PlayerApplyEnterMpResultReq{})       // 世界敲门处理请求
	a.registerMessage(ApiPlayerGetForceQuitBanInfoReq, &NullMsg{})                          // 退出世界请求
	a.registerMessage(ApiGetPlayerTokenRsp, &GetPlayerTokenRsp{})                           // 获取玩家token响应
	a.registerMessage(ApiPlayerLoginRsp, &PlayerLoginRsp{})                                 // 玩家登录响应
	a.registerMessage(ApiPingRsp, &PingRsp{})                                               // ping响应
	a.registerMessage(ApiPlayerSetPauseRsp, &PlayerSetPauseRsp{})                           // 玩家暂停响应
	a.registerMessage(ApiPlayerDataNotify, &PlayerDataNotify{})                             // 玩家信息通知
	a.registerMessage(ApiStoreWeightLimitNotify, &StoreWeightLimitNotify{})                 // 通知
	a.registerMessage(ApiPlayerStoreNotify, &PlayerStoreNotify{})                           // 通知
	a.registerMessage(ApiAvatarDataNotify, &AvatarDataNotify{})                             // 角色信息通知
	a.registerMessage(ApiPlayerEnterSceneNotify, &PlayerEnterSceneNotify{})                 // 玩家进入场景通知
	a.registerMessage(ApiOpenStateUpdateNotify, &OpenStateUpdateNotify{})                   // 通知
	a.registerMessage(ApiGetPlayerSocialDetailRsp, &GetPlayerSocialDetailRsp{})             // 获取玩家社区信息响应
	a.registerMessage(ApiEnterScenePeerNotify, &EnterScenePeerNotify{})                     // 进入场景对方通知
	a.registerMessage(ApiEnterSceneReadyRsp, &EnterSceneReadyRsp{})                         // 进入场景准备就绪响应
	a.registerMessage(ApiGetScenePointRsp, &GetScenePointRsp{})                             // 获取场景信息响应
	a.registerMessage(ApiGetSceneAreaRsp, &GetSceneAreaRsp{})                               // 获取场景区域响应
	a.registerMessage(ApiServerTimeNotify, &ServerTimeNotify{})                             // 服务器时间通知
	a.registerMessage(ApiWorldPlayerInfoNotify, &WorldPlayerInfoNotify{})                   // 世界玩家信息通知
	a.registerMessage(ApiWorldDataNotify, &WorldDataNotify{})                               // 世界数据通知
	a.registerMessage(ApiPlayerWorldSceneInfoListNotify, &PlayerWorldSceneInfoListNotify{}) // 场景解锁信息通知
	a.registerMessage(ApiHostPlayerNotify, &HostPlayerNotify{})                             // 主机玩家通知
	a.registerMessage(ApiSceneTimeNotify, &SceneTimeNotify{})                               // 场景时间通知
	a.registerMessage(ApiPlayerGameTimeNotify, &PlayerGameTimeNotify{})                     // 玩家游戏内时间通知
	a.registerMessage(ApiPlayerEnterSceneInfoNotify, &PlayerEnterSceneInfoNotify{})         // 玩家进入场景信息通知
	a.registerMessage(ApiSceneAreaWeatherNotify, &SceneAreaWeatherNotify{})                 // 场景区域天气通知
	a.registerMessage(ApiScenePlayerInfoNotify, &ScenePlayerInfoNotify{})                   // 场景玩家信息通知
	a.registerMessage(ApiSceneTeamUpdateNotify, &SceneTeamUpdateNotify{})                   // 场景队伍更新通知
	a.registerMessage(ApiSyncTeamEntityNotify, &SyncTeamEntityNotify{})                     // 同步队伍实体通知
	a.registerMessage(ApiSyncScenePlayTeamEntityNotify, &SyncScenePlayTeamEntityNotify{})   // 同步场景玩家队伍实体通知
	a.registerMessage(ApiSceneInitFinishRsp, &SceneInitFinishRsp{})                         // 场景初始化完成响应
	a.registerMessage(ApiEnterSceneDoneRsp, &EnterSceneDoneRsp{})                           // 进入场景完成响应
	a.registerMessage(ApiPlayerTimeNotify, &PlayerTimeNotify{})                             // 玩家对时通知
	a.registerMessage(ApiSceneEntityAppearNotify, &SceneEntityAppearNotify{})               // 场景实体出现通知
	a.registerMessage(ApiWorldPlayerLocationNotify, &WorldPlayerLocationNotify{})           // 世界玩家位置通知
	a.registerMessage(ApiScenePlayerLocationNotify, &ScenePlayerLocationNotify{})           // 场景玩家位置通知
	a.registerMessage(ApiWorldPlayerRTTNotify, &WorldPlayerRTTNotify{})                     // 世界玩家RTT时延
	a.registerMessage(ApiEnterWorldAreaRsp, &EnterWorldAreaRsp{})                           // 进入世界区域响应
	a.registerMessage(ApiPostEnterSceneRsp, &PostEnterSceneRsp{})                           // 提交进入场景响应
	a.registerMessage(ApiTowerAllDataRsp, &TowerAllDataRsp{})                               // 深渊数据响应
	a.registerMessage(ApiSceneTransToPointRsp, &SceneTransToPointRsp{})                     // 场景传送点响应
	a.registerMessage(ApiSceneEntityDisappearNotify, &SceneEntityDisappearNotify{})         // 场景实体消失通知
	a.registerMessage(ApiChangeAvatarRsp, &ChangeAvatarRsp{})                               // 更换角色响应
	a.registerMessage(ApiSetUpAvatarTeamRsp, &SetUpAvatarTeamRsp{})                         // 配置队伍响应
	a.registerMessage(ApiAvatarTeamUpdateNotify, &AvatarTeamUpdateNotify{})                 // 角色队伍更新通知
	a.registerMessage(ApiChooseCurAvatarTeamRsp, &ChooseCurAvatarTeamRsp{})                 // 切换队伍响应
	a.registerMessage(ApiStoreItemChangeNotify, &StoreItemChangeNotify{})                   // 背包道具变动通知
	a.registerMessage(ApiItemAddHintNotify, &ItemAddHintNotify{})                           // 道具增加提示通知
	a.registerMessage(ApiStoreItemDelNotify, &StoreItemDelNotify{})                         // 背包道具删除通知
	a.registerMessage(ApiPlayerPropNotify, &PlayerPropNotify{})                             // 玩家属性通知
	a.registerMessage(ApiGetGachaInfoRsp, &GetGachaInfoRsp{})                               // 卡池获取响应
	a.registerMessage(ApiDoGachaRsp, &DoGachaRsp{})                                         // 抽卡响应
	a.registerMessage(ApiEntityFightPropUpdateNotify, &EntityFightPropUpdateNotify{})       // 实体战斗属性更新通知
	a.registerMessage(ApiQueryPathRsp, &QueryPathRsp{})                                     // 寻路响应
	a.registerMessage(ApiAvatarFightPropNotify, &AvatarFightPropNotify{})                   // 角色战斗属性通知
	a.registerMessage(ApiAvatarEquipChangeNotify, &AvatarEquipChangeNotify{})               // 角色装备改变通知
	a.registerMessage(ApiAvatarAddNotify, &AvatarAddNotify{})                               // 角色新增通知
	a.registerMessage(ApiWearEquipRsp, &WearEquipRsp{})                                     // 装备穿戴响应
	a.registerMessage(ApiChangeGameTimeRsp, &ChangeGameTimeRsp{})                           // 改变游戏场景时间响应
	a.registerMessage(ApiSetPlayerBirthdayRsp, &SetPlayerBirthdayRsp{})                     // 设置生日响应
	a.registerMessage(ApiSetNameCardRsp, &SetNameCardRsp{})                                 // 修改名片响应
	a.registerMessage(ApiSetPlayerSignatureRsp, &SetPlayerSignatureRsp{})                   // 修改签名响应
	a.registerMessage(ApiSetPlayerNameRsp, &SetPlayerNameRsp{})                             // 修改昵称响应
	a.registerMessage(ApiSetPlayerHeadImageRsp, &SetPlayerHeadImageRsp{})                   // 修改头像响应
	a.registerMessage(ApiGetAllUnlockNameCardRsp, &GetAllUnlockNameCardRsp{})               // 获取全部已解锁名片响应
	a.registerMessage(ApiUnlockNameCardNotify, &UnlockNameCardNotify{})                     // 名片解锁通知
	a.registerMessage(ApiGetPlayerFriendListRsp, &GetPlayerFriendListRsp{})                 // 好友列表响应
	a.registerMessage(ApiGetPlayerAskFriendListRsp, &GetPlayerAskFriendListRsp{})           // 好友申请列表响应
	a.registerMessage(ApiAskAddFriendRsp, &AskAddFriendRsp{})                               // 加好友响应
	a.registerMessage(ApiAskAddFriendNotify, &AskAddFriendNotify{})                         // 加好友通知
	a.registerMessage(ApiDealAddFriendRsp, &DealAddFriendRsp{})                             // 处理好友申请响应
	a.registerMessage(ApiGetOnlinePlayerListRsp, &GetOnlinePlayerListRsp{})                 // 在线玩家列表响应
	a.registerMessage(ApiSceneForceUnlockNotify, &NullMsg{})                                // 场景强制解锁通知
	a.registerMessage(ApiSetPlayerBornDataRsp, &NullMsg{})                                  // 注册响应
	a.registerMessage(ApiDoSetPlayerBornDataNotify, &NullMsg{})                             // 注册通知
	a.registerMessage(ApiPathfindingEnterSceneRsp, &NullMsg{})                              // 寻路进入场景响应
	a.registerMessage(ApiPlayerForceExitRsp, &NullMsg{})                                    // 退出游戏响应
	a.registerMessage(ApiDelTeamEntityNotify, &DelTeamEntityNotify{})                       // 删除队伍实体通知
	a.registerMessage(ApiPlayerApplyEnterMpRsp, &PlayerApplyEnterMpRsp{})                   // 世界敲门响应
	a.registerMessage(ApiPlayerApplyEnterMpNotify, &PlayerApplyEnterMpNotify{})             // 世界敲门通知
	a.registerMessage(ApiPlayerApplyEnterMpResultRsp, &PlayerApplyEnterMpResultRsp{})       // 世界敲门处理响应
	a.registerMessage(ApiPlayerApplyEnterMpResultNotify, &PlayerApplyEnterMpResultNotify{}) // 世界敲门处理通知
	a.registerMessage(ApiPlayerGetForceQuitBanInfoRsp, &PlayerGetForceQuitBanInfoRsp{})     // 退出世界响应
	a.registerMessage(ApiGetShopmallDataReq, &GetShopmallDataReq{})                         // 商店信息请求
	a.registerMessage(ApiGetShopmallDataRsp, &GetShopmallDataRsp{})                         // 商店信息响应
	a.registerMessage(ApiGetShopReq, &GetShopReq{})                                         // 商店详情请求
	a.registerMessage(ApiGetShopRsp, &GetShopRsp{})                                         // 商店详情响应
	a.registerMessage(ApiBuyGoodsReq, &BuyGoodsReq{})                                       // 商店货物购买请求
	a.registerMessage(ApiBuyGoodsRsp, &BuyGoodsRsp{})                                       // 商店货物购买响应
	a.registerMessage(ApiMcoinExchangeHcoinReq, &McoinExchangeHcoinReq{})                   // 结晶换原石请求
	a.registerMessage(ApiMcoinExchangeHcoinRsp, &McoinExchangeHcoinRsp{})                   // 结晶换原石响应
	a.registerMessage(ApiAvatarChangeCostumeReq, &AvatarChangeCostumeReq{})                 // 角色换装请求
	a.registerMessage(ApiAvatarChangeCostumeRsp, &AvatarChangeCostumeRsp{})                 // 角色换装响应
	a.registerMessage(ApiAvatarChangeCostumeNotify, &AvatarChangeCostumeNotify{})           // 角色换装通知
	a.registerMessage(ApiAvatarWearFlycloakReq, &AvatarWearFlycloakReq{})                   // 角色换风之翼请求
	a.registerMessage(ApiAvatarWearFlycloakRsp, &AvatarWearFlycloakRsp{})                   // 角色换风之翼响应
	a.registerMessage(ApiAvatarFlycloakChangeNotify, &AvatarFlycloakChangeNotify{})         // 角色换风之翼通知
	// 尚未得知的客户端上行消息
	a.registerMessage(ApiClientAbilityChangeNotify, &ClientAbilityChangeNotify{})             // 未知
	a.registerMessage(ApiEvtAiSyncSkillCdNotify, &EvtAiSyncSkillCdNotify{})                   // 未知
	a.registerMessage(ApiEvtAiSyncCombatThreatInfoNotify, &EvtAiSyncCombatThreatInfoNotify{}) // 未知
	a.registerMessage(ApiEntityConfigHashNotify, &NullMsg{})                                  // 未知
	a.registerMessage(ApiMonsterAIConfigHashNotify, &NullMsg{})                               // 未知
	a.registerMessage(ApiGetRegionSearchReq, &NullMsg{})                                      // 未知
	a.registerMessage(ApiObstacleModifyNotify, &NullMsg{})                                    // 未知
	// TODO
	a.registerMessage(ApiEvtDoSkillSuccNotify, &EvtDoSkillSuccNotify{})
	a.registerMessage(ApiEvtCreateGadgetNotify, &EvtCreateGadgetNotify{})
	a.registerMessage(ApiEvtDestroyGadgetNotify, &EvtDestroyGadgetNotify{})
	// 空消息
	a.registerMessage(65535, &NullMsg{})
}

func (a *ApiProtoMap) registerMessage(apiId uint16, protoObj pb.Message) {
	_, exist := a.apiDeDupMap[apiId]
	if exist {
		logger.LOG.Error("reg dup msg, api id: %v", apiId)
		return
	} else {
		a.apiDeDupMap[apiId] = true
	}
	// apiId -> protoObj
	a.apiIdProtoObjMap[apiId] = reflect.TypeOf(protoObj)
	// protoObj -> apiId
	a.protoObjApiIdMap[reflect.TypeOf(protoObj)] = apiId
}

func (a *ApiProtoMap) GetProtoObjByApiId(apiId uint16) (protoObj pb.Message) {
	protoObjTypePointer, ok := a.apiIdProtoObjMap[apiId]
	if !ok {
		logger.LOG.Error("unknown api id: %v", apiId)
		protoObj = nil
		return protoObj
	}
	protoObjInst := reflect.New(protoObjTypePointer.Elem())
	protoObj = protoObjInst.Interface().(pb.Message)
	return protoObj
}

func (a *ApiProtoMap) GetApiIdByProtoObj(protoObj pb.Message) (apiId uint16) {
	var ok = false
	apiId, ok = a.protoObjApiIdMap[reflect.TypeOf(protoObj)]
	if !ok {
		logger.LOG.Error("unknown proto object: %v", protoObj)
		apiId = 0
	}
	return apiId
}
