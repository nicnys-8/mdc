package bitverse

import (
	//"fmt"
	"time"
)

type MsgService struct {
	id       string
	observer MsgServiceObserver
	edgeNode *EdgeNode
	aesKey   string
}

type msgServiceReplyType struct {
	msgReplyCallback func(timedOut bool, msg *Msg)
	timeout          int32
	timestamp        int32
}

func composeMsgService(secret string, id string, observe MsgServiceObserver, edgeNode *EdgeNode) *MsgService {
	service := new(MsgService)
	service.id = id
	service.observer = observe
	service.edgeNode = edgeNode
	service.aesKey = secret
	return service
}

func (msgService *MsgService) Send(dst string, payload string) {
	encryptedPayload := encrypt(msgService.aesKey, payload)
	msg := composeDataMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedPayload)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) SendAndGetReply(dst string, payload string, timeout int32, msgReplyCallback func(timedOut bool, msg *Msg)) {
	encryptedPayload := encrypt(msgService.aesKey, payload)
	msg := composeDataMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedPayload)

	msgServiceReply := new(msgServiceReplyType)
	msgServiceReply.timeout = timeout
	msgServiceReply.msgReplyCallback = msgReplyCallback
	msgServiceReply.timestamp = int32(time.Now().Unix())
	msgService.edgeNode.msgServiceReplies[msg.Id] = msgServiceReply

	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) Reply(msg *Msg, payload string) {
	encryptedPayload := encrypt(msgService.aesKey, payload)
	replyMsg := composeDataMsg(msgService.edgeNode.Id(), msg.Src, msgService.id, encryptedPayload)
	replyMsg.Id = msg.Id // use the same id as the sender
	msgService.edgeNode.send(replyMsg)
}
