package bitverse

import (
	"crypto/rsa"
	//"fmt"
)

type RepoMsgServiceObserver struct {
}

func (repoMsgServiceObserver *RepoMsgServiceObserver) OnDeliver(msgService *MsgService, msg *Msg) {
	// ignore, we wil only use SendAndGetReply
}

type RepoService struct {
	repoId           string
	edgeNode         *EdgeNode
	aesEncryptionKey string
	msgService       *MsgService
	prv              *rsa.PrivateKey
	pub              *rsa.PublicKey
}

func composeRepoService(aesEncryptionKey string, prv *rsa.PrivateKey, pub *rsa.PublicKey, repoId string, edgeNode *EdgeNode, msgService *MsgService) *RepoService {
	service := new(RepoService)
	service.repoId = repoId
	service.prv = prv
	service.pub = pub
	service.edgeNode = edgeNode
	service.aesEncryptionKey = aesEncryptionKey
	service.msgService = msgService

	return service
}

func (repoService *RepoService) Store(key string, value string, timeout int32, callback func(err error, oldValue interface{})) {
	encryptedKey := encryptAes(repoService.aesEncryptionKey, key)
	encryptedValue := encryptAes(repoService.aesEncryptionKey, value)
	signature, err := sign(repoService.prv, encryptedKey+encryptedValue)
	if err != nil {
		panic(err)
	}

	msg := composeRepoStoreMsg(repoService.edgeNode.Id(), repoService.edgeNode.superNode.Id(), repoService.repoId, encryptedKey, encryptedValue, signature)
	repoService.msgService.sendMsgAndGetReply(msg, timeout, callback)
}

func (repoService *RepoService) Lookup(key string, timeout int32, callback func(err error, value interface{})) {
	//msg := composeRepoLookupMsg(repoService.edgeNode.Id(), repoService.edgeNode.superNode.Id(), repoService.repoId, key)
	encryptedKey := encryptAes(repoService.aesEncryptionKey, key)
	encryptedValue := encryptAes(repoService.aesEncryptionKey, value)
	signature, err := sign(repoService.prv, encryptedKey+encryptedValue)
	if err != nil {
		panic(err)
	}

	msg := composeRepoStoreMsg(repoService.edgeNode.Id(), repoService.edgeNode.superNode.Id(), repoService.repoId, encryptedKey, encryptedValue, signature)
	repoService.msgService.sendMsgAndGetReply(msg, timeout, callback)

}
