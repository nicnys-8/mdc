package bitverse

type BitverseObserver interface {
	OnSiblingJoined(edgeNode *EdgeNode, id string)
	OnSiblingLeft(edgeNode *EdgeNode, id string)
	OnSiblingHeartbeat(edgeNode *EdgeNode, id string)
	OnChildrenReply(edgeNode *EdgeNode, id string, children []string)
	OnConnected(edgeNode *EdgeNode, remoteNode *RemoteNode)
}
