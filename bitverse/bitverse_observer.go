package bitverse

type BitverseObserver interface {
	OnSiblingJoin(edgeNode *EdgeNode, id string)
	OnSiblingExit(edgeNode *EdgeNode, id string)
	OnSiblingHeartbeat(edgeNode *EdgeNode, id string)
	OnChildrenReply(edgeNode *EdgeNode, id string, children []string)
	OnConnected(edgeNode *EdgeNode, remoteNode *RemoteNode)
}
