package main

type Transport interface {
	SetMsgChannel(msgChannel chan Msg)
	SetLinkChannel(linkChannel chan *Link)
	SetLocalNodeId(localNodeId NodeId)
	CreateLocalEndPoint(localAddress string, localPort string)
	ConnectRemoteEndPoint(address string)
}
