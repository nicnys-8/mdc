package bitverse

type WSTransport struct {
	localPort   string
	wsServer    *WsServer
	wsClient    *WsClient
	localNodeId NodeId
}

func MakeWSTransport() *WSTransport {
	wsTransport := new(WSTransport)
	return wsTransport
}

func (wsTransport *WSTransport) SetLocalNodeId(localNodeId NodeId) {
	wsTransport.localNodeId = localNodeId
}

func (wsTransport *WSTransport) Listen(localAddress string, localPort string, remoteNodeChannel chan *RemoteNode, msgChannel chan Msg) {
	wsServer := makeWsServer(wsTransport.localNodeId, msgChannel, remoteNodeChannel)
	wsTransport.localPort = localPort
	wsTransport.wsServer = wsServer
	wsServer.start(wsTransport.localPort)
}

func (wsTransport *WSTransport) ConnectToNode(remoteAddress string, remoteNodeChannel chan *RemoteNode, msgChannel chan Msg) {
	wsClient := makeWsClient(msgChannel, remoteNodeChannel, wsTransport.localNodeId)
	wsTransport.wsClient = wsClient
	wsClient.connect(remoteAddress)
}
