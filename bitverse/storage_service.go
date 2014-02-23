package bitverse

import (
	//"fmt"
	"time"
)

type StorageMsgServiceObserver struct {
}

func (storageMsgServiceObserver *StorageMsgServiceObserver) OnDeliver(msgService *MsgService, msg *Msg) {
	// ignore, we wil only use SendAndGetReply
}

type StorageService struct {
	repoId     string
	edgeNode   *EdgeNode
	aesKey     string
	msgService *MsgService
}

func composeStorageService(secret string, repoId string, edgeNode *EdgeNode) *StorageService {
	service := new(StorageService)
	service.repoId = repoId
	service.edgeNode = edgeNode
	service.aesKey = secret

	return service
}

func (msgService *StorageService) Store(key string, value string, timeout int32, callback func(timedOut bool, oldValue *string)) {
	msg := composeStorageServiceStoreMsg(msgService.edgeNode.Id(), msgService.edgeNode.superNode.Id(), msgService.repoId, key, value)

	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	msgService.edgeNode.replyTable[msg.Id] = reply

	msgService.edgeNode.send(msg)
}

func (msgService *StorageService) Lookup(key string, timeout int32, callback func(timedOut bool, value *string)) {
	msg := composeStorageServiceLookupMsg(msgService.edgeNode.Id(), msgService.edgeNode.superNode.Id(), msgService.repoId, key)

	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	msgService.edgeNode.replyTable[msg.Id] = reply

	msgService.edgeNode.send(msg)
}
