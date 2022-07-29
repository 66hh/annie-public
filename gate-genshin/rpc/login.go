package rpc

import (
	"flswld.com/gate-genshin-api/api/proto"
	"gate-genshin/kcp"
	"gate-genshin/net"
	"strconv"
	"strings"
	"time"
)

func (r *RpcManager) getPlayerToken(convId uint64, req *proto.GetPlayerTokenReq) (rsp *proto.GetPlayerTokenRsp) {
	uidStr := req.AccountUid
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		r.log.Error("parse uid error: %v", err)
		return nil
	}
	account, err := r.dao.QueryAccountByField("uid", uid)
	if err != nil {
		r.log.Error("query account error: %v", err)
		return nil
	}
	if account == nil {
		r.log.Error("account is nil")
		return nil
	}
	if account.ComboToken != req.AccountToken {
		r.log.Error("token error")
		return nil
	}
	// comboToken验证成功
	if account.Forbid {
		if account.ForbidEndTime > uint64(time.Now().Unix()) {
			// 封号通知
			rsp = new(proto.GetPlayerTokenRsp)
			rsp.Uid = uint32(account.PlayerID)
			rsp.IsProficientPlayer = true
			rsp.Retcode = 21
			rsp.Msg = "FORBID_CHEATING_PLUGINS"
			//rsp.BlackUidEndTime = 2051193600 // 2035-01-01 00:00:00
			rsp.BlackUidEndTime = uint32(account.ForbidEndTime)
			rsp.RegPlatform = 3
			rsp.CountryCode = "US"
			addr, exist := r.getAddrByConvId(convId)
			if !exist {
				r.log.Error("can not find addr by convId")
				return nil
			}
			split := strings.Split(addr, ":")
			rsp.ClientIpStr = split[0]
			return rsp
		} else {
			account.Forbid = false
			_, err := r.dao.UpdateAccountFieldByFieldName("uid", account.Uid, "forbid", false)
			if err != nil {
				r.log.Error("update db error: %v", err)
				return nil
			}
		}
	}
	oldConvId, oldExist := r.getConvIdByUserId(uint32(account.PlayerID))
	if oldExist {
		// 顶号
		r.kcpEventInput <- &net.KcpEvent{
			ConvId:       oldConvId,
			EventId:      net.KcpConnForceClose,
			EventMessage: uint32(kcp.EnetServerRelogin),
		}
	}
	r.setUserIdByConvId(convId, uint32(account.PlayerID))
	r.setConvIdByUserId(uint32(account.PlayerID), convId)
	r.setConnState(convId, ConnWaitLogin)
	// 返回响应
	rsp = new(proto.GetPlayerTokenRsp)
	rsp.Uid = uint32(account.PlayerID)
	rsp.Token = account.ComboToken
	rsp.AccountType = 1
	// TODO 要确定一下新注册的号这个值该返回什么
	rsp.IsProficientPlayer = true
	rsp.SecretKeySeed = 11468049314633205968
	rsp.SecurityCmdBuffer = r.secretKeyBuffer
	rsp.PlatformType = 3
	rsp.ChannelId = 1
	rsp.CountryCode = "US"
	rsp.ClientVersionRandomKey = "c25-314dd05b0b5f"
	rsp.RegPlatform = 3
	addr, exist := r.getAddrByConvId(convId)
	if !exist {
		r.log.Error("can not find addr by convId")
		return nil
	}
	split := strings.Split(addr, ":")
	rsp.ClientIpStr = split[0]
	return rsp
}

func (r *RpcManager) playerLogin(convId uint64, req *proto.PlayerLoginReq) (rsp *proto.PlayerLoginRsp) {
	userId, exist := r.getUserIdByConvId(convId)
	if !exist {
		r.log.Error("can not find userId by convId")
		return nil
	}
	account, err := r.dao.QueryAccountByField("playerID", userId)
	if err != nil {
		r.log.Error("query account error: %v", err)
		return nil
	}
	if account == nil {
		r.log.Error("account is nil")
		return nil
	}
	if account.ComboToken != req.Token {
		r.log.Error("token error")
		return nil
	}
	// comboToken验证成功
	r.setConnState(convId, ConnAlive)
	// 返回响应
	rsp = new(proto.PlayerLoginRsp)
	rsp.IsUseAbilityHash = true
	rsp.AbilityHashCode = 1844674
	rsp.GameBiz = "hk4e_global"
	rsp.ClientDataVersion = uint32(r.regionCurr.RegionInfo.ClientDataVersion)
	rsp.ClientSilenceDataVersion = uint32(r.regionCurr.RegionInfo.ClientSilenceDataVersion)
	rsp.ClientMd5 = r.regionCurr.RegionInfo.ClientDataMd5
	rsp.ClientSilenceMd5 = r.regionCurr.RegionInfo.ClientSilenceDataMd5
	rsp.ResVersionConfig = r.regionCurr.RegionInfo.ResVersionConfig
	rsp.ClientVersionSuffix = r.regionCurr.RegionInfo.ClientVersionSuffix
	rsp.ClientSilenceVersionSuffix = r.regionCurr.RegionInfo.ClientSilenceVersionSuffix
	rsp.IsScOpen = false
	rsp.RegisterCps = "mihoyo"
	rsp.CountryCode = "US"
	return rsp
}
