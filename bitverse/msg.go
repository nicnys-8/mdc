package bitverse

import (
//"encoding/base64"
//"fmt"
)

const (
	Handshake = iota
	Data
	Heartbeat
	Children
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
	msg.Service = ""
	return msg
}

func ComposeChildrenRequestMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Service = ""
	return msg
}

func (msg *Msg) String() string {
	if msg.Type == Heartbeat {
		return "msg[type:heartbeat to:" + msg.Dst + " from:" + msg.Src + "]"
	} else if msg.Type == Children {
		return "msg[type:children to:" + msg.Dst + " from:" + msg.Src + "]"
	} else if msg.Type == Data {
		return "msg[type:data to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + " service:" + msg.Service + "]"
	} else {
		return "msg[type:unkown]"
	}
}

/// PRIVATE

func (msg *Msg) encodePayload(payload string) (encodedPayload string) {
	//encodedPayload = base64.StdEncoding.EncodeToString([]byte(payload))
	//fmt.Println(encodedPayload)
	return
}

func (msg *Msg) decodePayload(encodedPayload string) (payload string) {
	/*var err err
	payload, err = base64.StdEncoding.DecodeString(str)

	if err != nil {
		log.Fatal("msg: failed to decode payload", err)
	} */
	return
}
