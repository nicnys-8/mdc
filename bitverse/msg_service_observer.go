package bitverse

type MsgServiceObserver interface {
	OnDeliver(msgService *MsgService, msg *Msg)
}
