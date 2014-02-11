package main

const (
	Handshake = iota
	Data
	Announce
	Bye
)

type MsgId string

type Msg struct {
	Type    int
	Payload string
	Src     string
	Dst     string
}
