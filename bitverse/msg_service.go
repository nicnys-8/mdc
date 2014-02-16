package bitverse

type MsgService struct {
	id       string
	observer MsgServiceObserver
	edgeNode *EdgeNode
	aesKey   string
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
	msgService.edgeNode.send(dst, encryptedPayload, msgService.id)
}
