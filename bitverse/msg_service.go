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
	callback  func(success bool, data interface{})
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

func (msgService *MsgService) Send(dst string, data string) {
	encryptedData := encryptAes(msgService.aesKey, data)
	msg := composeMsgServiceMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedData)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) SendAndGetReply(dst string, data string, timeout int32, callback func(success bool, data interface{})) {
	encryptedData := encryptAes(msgService.aesKey, data)

	msg := composeMsgServiceMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedData)

	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	msgService.edgeNode.replyTable[msg.Id] = reply

	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) Reply(msg *Msg, data string) {
	encryptedData := encryptAes(msgService.aesKey, data)
	replyMsg := composeMsgServiceMsg(msgService.edgeNode.Id(), msg.Src, msgService.id, encryptedData)
	replyMsg.Id = msg.Id // use the same id as the sender

	msgService.edgeNode.send(replyMsg)
}
