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

type msgReplyType struct {
	callback  func(timedOut bool, msg *Msg)
	timeout   int32
	timestamp int32
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
	msg := composeMsgServiceMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedPayload)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) SendAndGetReply(dst string, payload string, timeout int32, callback func(timedOut bool, msg *Msg)) {
	encryptedPayload := encrypt(msgService.aesKey, payload)

	msg := composeMsgServiceMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedPayload)

	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	msgService.edgeNode.replyTable[msg.Id] = reply

	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) Reply(msg *Msg, payload string) {
	encryptedPayload := encrypt(msgService.aesKey, payload)
	replyMsg := composeMsgServiceMsg(msgService.edgeNode.Id(), msg.Src, msgService.id, encryptedPayload)
	replyMsg.Id = msg.Id // use the same id as the sender

	msgService.edgeNode.send(replyMsg)
}
