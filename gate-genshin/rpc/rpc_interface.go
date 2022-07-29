package rpc

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/gm"
	"gate-genshin/kcp"
	"gate-genshin/net"
	"github.com/pkg/errors"
)

// rpc interface

// 从GS接收消息
func (r *RpcManager) RecvNetMsgFromGameServer(netMsg *api.NetMsg, res *bool) error {
	if netMsg == nil || res == nil {
		return errors.New("param is nil")
	}
	*res = true
	if netMsg.EventId == api.NormalMsg {
		protoMsg := new(net.ProtoMsg)
		convId, exist := r.getConvIdByUserId(netMsg.UserId)
		if exist {
			protoMsg.ConvId = convId
		} else {
			r.log.Error("can not find convId by userId")
			return nil
		}
		protoMsg.ApiId = netMsg.ApiId
		protoMsg.HeadMessage = netMsg.HeadMessage
		protoMsg.PayloadMessage = netMsg.PayloadMessage
		r.protoMsgInput <- protoMsg
		return nil
	} else {
		r.log.Info("recv event from game server, event id: %v", netMsg.EventId)
		return nil
	}
}

// 改变网关开放状态
func (r *RpcManager) ChangeGateOpenState(isOpen *bool, res *bool) error {
	if isOpen == nil || res == nil {
		return errors.New("param is nil")
	}
	*res = true
	r.kcpEventInput <- &net.KcpEvent{
		EventId:      net.KcpGateOpenState,
		EventMessage: *isOpen,
	}
	r.log.Info("change gate open state to: %v", *isOpen)
	return nil
}

// 剔除玩家下线
func (r *RpcManager) KickPlayer(info *gm.KickPlayerInfo, result *bool) error {
	if info == nil || result == nil {
		return errors.New("param is nil")
	}
	convId, exist := r.getConvIdByUserId(info.UserId)
	if !exist {
		*result = false
		return nil
	}
	r.kcpEventInput <- &net.KcpEvent{
		ConvId:       convId,
		EventId:      net.KcpConnForceClose,
		EventMessage: info.Reason,
	}
	*result = true
	return nil
}

// 获取网关在线玩家信息
func (r *RpcManager) GetOnlineUser(uid *uint32, list *gm.OnlineUserList) error {
	if uid == nil || list == nil {
		return errors.New("param is nil")
	}
	if *uid == 0 {
		// 获取全部玩家
		r.convUserIdMapLock.RLock()
		r.convAddrMapLock.RLock()
		for convId, userId := range r.convUserIdMap {
			addr := r.convAddrMap[convId]
			info := &gm.OnlineUserInfo{
				Uid:    userId,
				ConvId: convId,
				Addr:   addr,
			}
			list.UserList = append(list.UserList, info)
		}
		r.convAddrMapLock.RUnlock()
		r.convUserIdMapLock.RUnlock()
	} else {
		// 获取指定uid玩家
		convId, exist := r.getConvIdByUserId(*uid)
		if !exist {
			return nil
		}
		addr, exist := r.getAddrByConvId(convId)
		if !exist {
			return nil
		}
		info := &gm.OnlineUserInfo{
			Uid:    *uid,
			ConvId: convId,
			Addr:   addr,
		}
		list.UserList = append(list.UserList, info)
	}
	return nil
}

// 用户密码改变
func (r *RpcManager) UserPasswordChange(uid *uint32, result *bool) error {
	if uid == nil || result == nil {
		return errors.New("param is nil")
	}
	// dispatch登录态失效
	_, err := r.dao.UpdateAccountFieldByFieldName("uid", *uid, "token", "")
	if err != nil {
		*result = false
		return nil
	}
	// 游戏内登录态失效
	account, err := r.dao.QueryAccountByField("uid", *uid)
	if err != nil {
		*result = false
		return nil
	}
	if account == nil {
		*result = false
		return nil
	}
	convId, exist := r.getConvIdByUserId(uint32(account.PlayerID))
	if !exist {
		*result = true
		return nil
	}
	r.kcpEventInput <- &net.KcpEvent{
		ConvId:       convId,
		EventId:      net.KcpConnForceClose,
		EventMessage: uint32(kcp.EnetAccountPasswordChange),
	}
	*result = true
	return nil
}

// 封号
func (r *RpcManager) ForbidUser(info *gm.ForbidUserInfo, result *bool) error {
	if info == nil || result == nil {
		return errors.New("param is nil")
	}
	// 写入账号封禁信息
	_, err := r.dao.UpdateAccountFieldByFieldName("uid", info.UserId, "forbid", true)
	if err != nil {
		*result = false
		return nil
	}
	_, err = r.dao.UpdateAccountFieldByFieldName("uid", info.UserId, "forbidEndTime", info.ForbidEndTime)
	if err != nil {
		*result = false
		return nil
	}
	// 游戏强制下线
	account, err := r.dao.QueryAccountByField("uid", info.UserId)
	if err != nil {
		*result = false
		return nil
	}
	if account == nil {
		*result = false
		return nil
	}
	convId, exist := r.getConvIdByUserId(uint32(account.PlayerID))
	if !exist {
		*result = true
		return nil
	}
	r.kcpEventInput <- &net.KcpEvent{
		ConvId:       convId,
		EventId:      net.KcpConnForceClose,
		EventMessage: uint32(kcp.EnetServerKillClient),
	}
	*result = true
	return nil
}

// 解封
func (r *RpcManager) UnForbidUser(uid *uint32, result *bool) error {
	if uid == nil || result == nil {
		return errors.New("param is nil")
	}
	// 解除账号封禁
	_, err := r.dao.UpdateAccountFieldByFieldName("uid", *uid, "forbid", false)
	if err != nil {
		*result = false
		return nil
	}
	*result = true
	return nil
}
