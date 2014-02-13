package bitverse

type Transport interface {
	SetLocalNodeId(localNodeId NodeId)
	Listen(localAddress string, localPort string, remoteNodeChannels chan *RemoteNode, msgChannel chan Msg)
	ConnectToNode(remoteAddress string, remoteNodeChannels chan *RemoteNode, msgChannel chan Msg)
}
