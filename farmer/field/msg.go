package main

const (
	Handshake = iota
	Data
	Discovery
	Bye
)

type Msg struct {
	Type    int
	Payload string
	SeqNr   int
	Src     NodeId
	Dst     NodeId
}
