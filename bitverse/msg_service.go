package bitverse

type MsgService struct {
	id               string
	observer         MsgServiceObserver
	edgeNode         *EdgeNode
	aesEncryptionKey string
}

type msgReplyType struct {
	callback  func(err error, data interface{})
	timeout   int32
	timestamp int32
}

func composeMsgService(aesEncryptionKey string, id string, observe MsgServiceObserver, edgeNode *EdgeNode) *MsgService {
	service := new(MsgService)
	service.id = id
	service.observer = observe
	service.edgeNode = edgeNode
	service.aesEncryptionKey = aesEncryptionKey
	return service
}

func (msgService *MsgService) Send(dst string, data string) {
	encryptedData := encryptAes(msgService.aesEncryptionKey, data)
	msg := composeMsgServiceMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedData)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) SendAndGetReply(dst string, data string, timeout int32, callback func(err error, data interface{})) {
	encryptedData := encryptAes(msgService.aesEncryptionKey, data)
	msg := composeMsgServiceMsg(msgService.edgeNode.Id(), dst, msgService.id, encryptedData)
	msgService.edgeNode.registerReplyCallback(msg.Id, timeout, callback)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) Reply(msg *Msg, data string) {
	encryptedData := encryptAes(msgService.aesEncryptionKey, data)
	replyMsg := composeMsgServiceMsg(msgService.edgeNode.Id(), msg.Src, msgService.id, encryptedData)
	replyMsg.Id = msg.Id // use the same id as the sender

	msgService.edgeNode.send(replyMsg)
}

/// PRIVATE

func (msgService *MsgService) sendMsg(msg *Msg) {
	msg.Payload = encryptAes(msgService.aesEncryptionKey, msg.Payload)
	msgService.edgeNode.send(msg)
}

func (msgService *MsgService) sendMsgAndGetReply(msg *Msg, timeout int32, callback func(err error, data interface{})) {
	if msg == nil {
		panic("msg is nil")
	}

	if msgService == nil {
		panic("msg service is nil")
	}

	msg.Payload = encryptAes(msgService.aesEncryptionKey, msg.Payload)
	msgService.edgeNode.registerReplyCallback(msg.Id, timeout, callback)
	msgService.edgeNode.send(msg)
}
