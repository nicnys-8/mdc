package bitverse

type ServiceObserver interface {
	OnError(err error)
	OnDeliver(msg *Msg)
	OnSiblingJoin(nodeId string)
	OnSiblingExit(nodeId string)
	OnSiblingHeartbeat(nodeId string)
}
