package bitverse

import (
	"fmt"
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

type Msg struct {
	Type      int
	Payload   string
	Src       string
	Dst       string
	ServiceId string
	Id        string
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

func composeDataMsg(src string, dst string, serviceId string, payload string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = payload
	msg.Src = src
	msg.Dst = dst
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())
	msg.ServiceId = serviceId

	return msg
}

func composeHeartbeatMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Heartbeat
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Id = ""
	msg.ServiceId = ""
	return msg
}

func composeChildrenRequestMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Id = ""
	msg.ServiceId = ""
	return msg
}

func composeChildrenReplyMsg(src string, dst string, json string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = json
	msg.Src = src
	msg.Dst = dst
	msg.Id = ""
	msg.ServiceId = ""
	return msg
}

func composeChildJoin(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildJoined
	msg.Payload = childId
	msg.Src = src
	msg.Dst = ""
	msg.Id = ""
	msg.ServiceId = ""
	return msg
}

func composeChildLeft(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildLeft
	msg.Payload = childId
	msg.Src = src
	msg.Dst = ""
	msg.Id = ""
	msg.ServiceId = ""
	return msg
}

func composeHandshakeMsg(src string) *Msg {
	msg := new(Msg)
	msg.Type = Handshake
	msg.Payload = ""
	msg.Src = src
	msg.Dst = ""
	msg.Id = ""
	msg.ServiceId = ""
	return msg
}

func getSeqNr() int {
	mutex.Lock()
	seqNrCounter++
	mutex.Unlock()
	return seqNrCounter
}
