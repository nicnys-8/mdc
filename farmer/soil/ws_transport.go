package main

type WSTransport struct {
	Transport
	msgChannel  chan Msg
	linkChannel chan *Link
	localPort   string
	wsServer    *WsServer
	wsClient    *WsClient
	localNodeId NodeId
}

func makeWSTransport() *WSTransport {
	wsTransport := new(WSTransport)
	return wsTransport
}

func (wsTransport *WSTransport) SetLinkChannel(linkChannel chan *Link) {
	wsTransport.linkChannel = linkChannel
}

func (wsTransport *WSTransport) SetMsgChannel(msgChannel chan Msg) {
	wsTransport.msgChannel = msgChannel
}

func (wsTransport *WSTransport) SetLocalNodeId(localNodeId NodeId) {
	wsTransport.localNodeId = localNodeId
}

func (wsTransport *WSTransport) CreateLocalEndPoint(localAddress string, localPort string) {
	wsServer := makeWsServer(wsTransport.localNodeId, wsTransport.msgChannel, wsTransport.linkChannel)
	wsTransport.localPort = localPort
	wsTransport.wsServer = wsServer
	go wsServer.start(wsTransport.localPort)
}

func (wsTransport *WSTransport) ConnectRemoteEndPoint(ipAddress string) {
	wsClient := makeWsClient(wsTransport.msgChannel, wsTransport.linkChannel, wsTransport.localNodeId)
	wsTransport.wsClient = wsClient
	wsClient.connect(ipAddress)
}
