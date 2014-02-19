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

type MsgReplyCallback func(timedOut bool, msg *Msg)

type MsgServiceReply struct {
	msgReplyCallback MsgReplyCallback
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
	msg := ComposeDataMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedPayload)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) SendAndGetReply(dst string, payload string, timeout int32, msgReplyCallback MsgReplyCallback) {
	encryptedPayload := encrypt(msgService.aesKey, payload)
	msg := ComposeDataMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedPayload)

	msgServiceReply := new(MsgServiceReply)
	msgServiceReply.timeout = timeout
	msgServiceReply.msgReplyCallback = msgReplyCallback
	msgServiceReply.timestamp = int32(time.Now().Unix())
	msgService.edgeNode.msgServiceReplies[msg.Id] = msgServiceReply

	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) Reply(msg *Msg, payload string) {
	encryptedPayload := encrypt(msgService.aesKey, payload)
	replyMsg := ComposeDataMsg(msgService.edgeNode.Id(), msg.Src, msgService.id, encryptedPayload)
	replyMsg.Id = msg.Id // use the same id as the sender
	msgService.edgeNode.send(replyMsg)
}
