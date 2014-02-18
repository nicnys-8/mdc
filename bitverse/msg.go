package bitverse

import (
	"sync"
)

const (
	Handshake = iota
	Data
	Heartbeat
	Children
	ChildJoined
	ChildLeft
	Bye
)

var mutex sync.Mutex
var seqNrCounter int = 0

type MsgId string

type Msg struct {
	Type      int
	Payload   string
	Src       string
	Dst       string
	ServiceId string
	Id        int
}

func ComposeDataMsg(src string, dst string, serviceId string, payload string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = payload
	msg.Src = src
	msg.Dst = dst
	msg.Id = getSeqNr()
	msg.ServiceId = serviceId

	return msg
}

func ComposeHeartbeatMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Heartbeat
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Id = 0
	msg.ServiceId = ""
	return msg
}

func ComposeChildrenRequestMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Id = 0
	msg.ServiceId = ""
	return msg
}

func ComposeChildrenReplyMsg(src string, dst string, json string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = json
	msg.Src = src
	msg.Dst = dst
	msg.Id = 0
	msg.ServiceId = ""
	return msg
}

func ComposeChildJoin(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildJoined
	msg.Payload = childId
	msg.Src = src
	msg.Dst = ""
	msg.Id = 0
	msg.ServiceId = ""
	return msg
}

func ComposeChildLeft(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildLeft
	msg.Payload = childId
	msg.Src = src
	msg.Dst = ""
	msg.Id = 0
	msg.ServiceId = ""
	return msg
}

func ComposeHandshakeMsg(src string) *Msg {
	msg := new(Msg)
	msg.Type = Handshake
	msg.Payload = ""
	msg.Src = src
	msg.Dst = ""
	msg.Id = 0
	msg.ServiceId = ""
	return msg
}

func (msg *Msg) String() string {
	if msg.Type == Heartbeat {
		return "msg[type:heartbeat to:" + msg.Dst + " from:" + msg.Src + "]"
	} else if msg.Type == Handshake {
		return "msg[type:handshake to:" + msg.Dst + " from:" + msg.Src + "]"
	} else if msg.Type == Children {
		return "msg[type:children to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + "]"
	} else if msg.Type == ChildJoined {
		return "msg[type:childjoined to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + "]"
	} else if msg.Type == ChildJoined {
		return "msg[type:childleft to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + "]"
	} else if msg.Type == Data {
		return "msg[type:data to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + " service:" + msg.ServiceId + "]"
	} else {
		return "msg[type:unkown]"
	}
}

/// PRIVATE

func getSeqNr() int {
	mutex.Lock()
	seqNrCounter++
	mutex.Unlock()
	return seqNrCounter
}
