package bitverse

const (
	Handshake = iota
	Data
	Heartbeat
	Children
	Bye
)

type MsgId string

type Msg struct {
	Type      int
	Payload   string
	Src       string
	Dst       string
	ServiceId string
}

func ComposeDataMsg(src string, dst string, serviceId string, payload string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = payload
	msg.Src = src
	msg.Dst = dst
	msg.ServiceId = serviceId
	return msg
}

func ComposeHeartbeatMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Heartbeat
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.ServiceId = ""
	return msg
}

func ComposeChildrenRequestMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.ServiceId = ""
	return msg
}

func ComposeChildrenReplyMsg(src string, dst string, json string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = json
	msg.Src = src
	msg.Dst = dst
	msg.ServiceId = ""
	return msg
}

func ComposeHandshakeMsg(src string) *Msg {
	msg := new(Msg)
	msg.Type = Handshake
	msg.Payload = ""
	msg.Src = src
	msg.Dst = ""
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
	} else if msg.Type == Data {
		return "msg[type:data to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + " service:" + msg.ServiceId + "]"
	} else {
		return "msg[type:unkown]"
	}
}
