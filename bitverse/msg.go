package bitverse

const (
	Handshake = iota
	Data
	Heartbeat
	Bye
)

type MsgId string

type Msg struct {
	Type    int
	Payload string
	Src     string
	Dst     string
	Service string
}

func ComposeDataMsg(src string, dst string, service string, payload string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = payload
	msg.Src = src
	msg.Dst = dst
	msg.Service = service
	return msg
}

func ComposeHeartbeatMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Heartbeat
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	return msg
}

func (msg *Msg) String() string {
	if msg.Type == Heartbeat {
		return "msg[type:heartbeat to:" + msg.Dst + " from:" + msg.Src + "]"
	} else if msg.Type == Data {
		return "msg[type:data to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + " service:" + msg.Service + "]"
	} else {
		return "msg[type:unkown]"
	}
}
