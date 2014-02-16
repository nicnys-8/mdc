package bitverse

type BitverseObserver interface {
	OnError(err error)
	OnSiblingJoin(nodeId string)
	OnSiblingExit(nodeId string)
	OnSiblingHeartbeat(node string)
	OnChildrenReply(nodeId string)
	OnConnected(edgeNode *EdgeNode, remoteNode *RemoteNode)
}
