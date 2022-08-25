package net

import (
	"flswld.com/gate-genshin-api/api"
	"flswld.com/gate-genshin-api/api/proto"
	"flswld.com/logger"
	pb "google.golang.org/protobuf/proto"
	"reflect"
	"runtime"
)

type ProtoEnDecode struct {
	kcpMsgInput      chan *KcpMsg
	kcpMsgOutput     chan *KcpMsg
	protoMsgInput    chan *ProtoMsg
	protoMsgOutput   chan *ProtoMsg
	apiIdProtoObjMap map[uint16]reflect.Type
	protoObjApiIdMap map[reflect.Type]uint16
	bypassApiMap     map[uint16]bool
}

func NewProtoEnDecode(kcpMsgInput chan *KcpMsg, kcpMsgOutput chan *KcpMsg, protoMsgInput chan *ProtoMsg, protoMsgOutput chan *ProtoMsg) (r *ProtoEnDecode) {
	r = new(ProtoEnDecode)
	r.kcpMsgInput = kcpMsgInput
	r.kcpMsgOutput = kcpMsgOutput
	r.protoMsgInput = protoMsgInput
	r.protoMsgOutput = protoMsgOutput
	r.apiIdProtoObjMap = make(map[uint16]reflect.Type)
	r.protoObjApiIdMap = make(map[reflect.Type]uint16)
	r.bypassApiMap = make(map[uint16]bool)
	r.initMsgProtoMap()
	return r
}

func (p *ProtoEnDecode) Start() {
	cpuCoreNum := runtime.NumCPU()
	for i := 0; i < cpuCoreNum; i++ {
		go p.protoDecode()
		go p.protoEncode()
	}
}

type ProtoMessage struct {
	apiId   uint16
	message any
}

func (p *ProtoEnDecode) protoDecode() {
	for {
		kcpMsg := <-p.kcpMsgOutput
		protoMsg := new(ProtoMsg)
		protoMsg.ConvId = kcpMsg.ConvId
		protoMsg.ApiId = kcpMsg.ApiId
		// head msg
		if kcpMsg.HeadData != nil && len(kcpMsg.HeadData) != 0 {
			headMsg := new(api.PacketHead)
			err := pb.Unmarshal(kcpMsg.HeadData, headMsg)
			if err != nil {
				logger.LOG.Error("unmarshal head data err: %v", err)
				continue
			}
			protoMsg.HeadMessage = headMsg
		} else {
			protoMsg.HeadMessage = nil
		}
		// payload msg
		isBypass := p.bypassApiMap[kcpMsg.ApiId]
		if !isBypass {
			protoMessageList := make([]*ProtoMessage, 0)
			p.protoDecodePayloadCore(kcpMsg.ApiId, kcpMsg.ProtoData, &protoMessageList)
			if len(protoMessageList) == 0 {
				logger.LOG.Error("decode proto object is nil")
				continue
			}
			if kcpMsg.ApiId == api.ApiUnionCmdNotify {
				for _, protoMessage := range protoMessageList {
					msg := new(ProtoMsg)
					msg.ConvId = kcpMsg.ConvId
					msg.ApiId = protoMessage.apiId
					msg.PayloadMessage = protoMessage.message
					//logger.LOG.Debug("[recv] union proto msg, convId: %v, apiId: %v", msg.ConvId, msg.ApiId)
					if protoMessage.apiId == api.ApiUnionCmdNotify {
						// 聚合消息自身不再往后发送
						continue
					}
					//logger.LOG.Debug("[recv] proto msg, convId: %v, apiId: %v, headMsg: %v", protoMsg.ConvId, protoMsg.ApiId, protoMsg.HeadMessage)
					p.protoMsgOutput <- msg
				}
				// 聚合消息自身不再往后发送
				continue
			} else {
				protoMsg.PayloadMessage = protoMessageList[0].message
			}
		} else {
			protoMsg.PayloadMessage = nil
		}
		//logger.LOG.Debug("[recv] proto msg, convId: %v, apiId: %v, headMsg: %v", protoMsg.ConvId, protoMsg.ApiId, protoMsg.HeadMessage)
		p.protoMsgOutput <- protoMsg
	}
}

func (p *ProtoEnDecode) protoDecodePayloadCore(apiId uint16, protoData []byte, protoMessageList *[]*ProtoMessage) {
	protoObj := p.decodePayloadToProto(apiId, protoData)
	if protoObj == nil {
		logger.LOG.Error("decode proto object is nil")
		return
	}
	if apiId == api.ApiUnionCmdNotify {
		// 处理聚合消息
		unionCmdNotify, ok := protoObj.(*proto.UnionCmdNotify)
		if !ok {
			logger.LOG.Error("parse union cmd error")
			return
		}
		for _, cmd := range unionCmdNotify.GetCmdList() {
			p.protoDecodePayloadCore(uint16(cmd.MessageId), cmd.Body, protoMessageList)
		}
	}
	*protoMessageList = append(*protoMessageList, &ProtoMessage{
		apiId:   apiId,
		message: protoObj,
	})
}

func (p *ProtoEnDecode) protoEncode() {
	for {
		protoMsg := <-p.protoMsgInput
		//logger.LOG.Debug("[send] proto msg, convId: %v, apiId: %v, headMsg: %v", protoMsg.ConvId, protoMsg.ApiId, protoMsg.HeadMessage)
		kcpMsg := new(KcpMsg)
		kcpMsg.ConvId = protoMsg.ConvId
		kcpMsg.ApiId = protoMsg.ApiId
		// head msg
		if protoMsg.HeadMessage != nil {
			headData, err := pb.Marshal(protoMsg.HeadMessage)
			if err != nil {
				logger.LOG.Error("marshal head data err: %v", err)
				continue
			}
			kcpMsg.HeadData = headData
		} else {
			kcpMsg.HeadData = nil
		}
		// payload msg
		if protoMsg.PayloadMessage != nil {
			apiId, protoData := p.encodeProtoToPayload(protoMsg.PayloadMessage)
			if apiId == 0 || protoData == nil {
				logger.LOG.Error("encode proto data is nil")
				continue
			}
			if apiId != protoMsg.ApiId {
				logger.LOG.Error("api id is not match with proto obj")
				continue
			}
			kcpMsg.ProtoData = protoData
		} else {
			kcpMsg.ProtoData = nil
		}
		p.kcpMsgInput <- kcpMsg
	}
}

func (p *ProtoEnDecode) decodePayloadToProto(apiId uint16, protoData []byte) (protoObj any) {
	protoObj = p.getProtoObjByApiId(apiId)
	if protoObj == nil {
		logger.LOG.Error("get new proto object is nil")
		return nil
	}
	err := pb.Unmarshal(protoData, protoObj.(pb.Message))
	if err != nil {
		logger.LOG.Error("unmarshal proto data err: %v", err)
		return nil
	}
	return protoObj
}

func (p *ProtoEnDecode) encodeProtoToPayload(protoObj any) (apiId uint16, protoData []byte) {
	apiId = p.getApiIdByProtoObj(protoObj)
	var err error = nil
	protoData, err = pb.Marshal(protoObj.(pb.Message))
	if err != nil {
		logger.LOG.Error("marshal proto object err: %v", err)
		return 0, nil
	}
	return apiId, protoData
}
