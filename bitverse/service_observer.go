package bitverse

type ServiceObserver interface {
	OnError(err error)
	OnDeliver(msg *Msg)
}
