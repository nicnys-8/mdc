package bitverse

import (
	//"fmt"
	"time"
)

type RepoMsgServiceObserver struct {
}

func (repoMsgServiceObserver *RepoMsgServiceObserver) OnDeliver(msgService *MsgService, msg *Msg) {
	// ignore, we wil only use SendAndGetReply
}

type RepoService struct {
	repoId     string
	edgeNode   *EdgeNode
	aesKey     string
	msgService *MsgService
}

func composeRepoService(secret string, repoId string, edgeNode *EdgeNode) *RepoService {
	service := new(RepoService)
	service.repoId = repoId
	service.edgeNode = edgeNode
	service.aesKey = secret

	return service
}

func (repoService *RepoService) Store(key string, value string, timeout int32, callback func(timedOut bool, oldValue interface{})) {
	msg := composeStorageServiceStoreMsg(repoService.edgeNode.Id(), repoService.edgeNode.superNode.Id(), repoService.repoId, key, value)

	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	repoService.edgeNode.replyTable[msg.Id] = reply

	repoService.edgeNode.send(msg)
}

func (repoService *RepoService) Lookup(key string, timeout int32, callback func(timedOut bool, value interface{})) {
	msg := composeStorageServiceLookupMsg(repoService.edgeNode.Id(), repoService.edgeNode.superNode.Id(), repoService.repoId, key)

	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	repoService.edgeNode.replyTable[msg.Id] = reply

	repoService.edgeNode.send(msg)
}
