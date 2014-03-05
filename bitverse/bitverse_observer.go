package bitverse

type BitverseObserver interface {
	OnSiblingJoined(node *EdgeNode, id string)
	OnSiblingLeft(node *EdgeNode, id string)
	OnSiblingHeartbeat(node *EdgeNode, id string)
	OnChildrenReply(node *EdgeNode, id string, children []string)
	OnConnected(node *EdgeNode, superNode *RemoteNode)
}
