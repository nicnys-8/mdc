package main

const (
	Handshake = iota
	Data
	Discovery
	Bye
)

type MsgId string

type Msg struct {
	Type        int
	Payload     string
	Id          MsgId
	Src         NodeId
	Dst         NodeId
	RouteRecord map[NodeId]NodeId
	LastHop     NodeId
}
