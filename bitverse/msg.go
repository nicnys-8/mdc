package bitverse

import (
	"fmt"
	"sync"
)

// message type definition
const (
	Handshake = iota
	Data
	Heartbeat
	Children
	ChildJoined
	ChildLeft
	Bye
)

// service type definition
const (
	Messaging = iota
	Repo
	Control
)

// repo cmd:s
const (
	Store = iota
	Lookup
	Claim
)

// status
const (
	Ok = iota
	Error
)

// payload type
const (
	String = iota
	Nil
)

var mutex sync.Mutex
var seqNrCounter int = 0

type Msg struct {
	Type           int    // message type
	Payload        string // payload
	PayloadType    int    // payload format, e.g nil or string
	Src            string // source
	Dst            string // desintation
	Id             string // unique id as calculated by sender
	ServiceType    int    // service type
	Signature      string // rsa signature
	MsgServiceName string // used by messaging service
	RepoId         string // used by repo service
	RepoCmd        int    // used by repo service
	RepoKey        string // used by repo service
	RepoValue      string // used by repo service
	Status         int    // status, e.g. Ok or Error
	msgService     *MsgService
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
		return "msg[type:data to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + " msgchannelid:" + msg.MsgServiceName + "]"
	} else {
		return "msg[type:unkown]"
	}
}

/// PRIVATE

/// Messaging service messages

func composeMsgServiceMsg(src string, dst string, serviceId string, payload string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = payload
	msg.PayloadType = String
	msg.Src = src
	msg.Dst = dst
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())

	msg.ServiceType = Messaging
	msg.MsgServiceName = serviceId

	msg.Status = Ok

	return msg
}

/// Repo service messages

func composeRepoClaimMsg(src string, superNodeId string, repoId string, publicKey string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Src = src
	msg.Dst = superNodeId
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())
	msg.Signature = publicKey

	msg.MsgServiceName = repoId
	msg.ServiceType = Repo

	msg.RepoId = repoId
	msg.RepoCmd = Claim

	msg.Status = Ok

	return msg
}

func composeRepoStoreMsg(src string, superNodeId string, repoId string, key string, value string, signature string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Src = src
	msg.Dst = superNodeId
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())

	msg.MsgServiceName = repoId
	msg.ServiceType = Repo

	msg.RepoId = repoId
	msg.RepoCmd = Store
	msg.RepoKey = key
	msg.RepoValue = value

	msg.Signature = signature

	msg.Status = Ok

	return msg
}

func composeRepoLookupMsg(src string, superNodeId string, repoId string, key string, signature string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Src = src
	msg.Dst = superNodeId
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())

	msg.MsgServiceName = repoId
	msg.ServiceType = Repo

	msg.RepoId = repoId
	msg.RepoCmd = Lookup
	msg.RepoKey = key

	msg.Signature = signature

	msg.Status = Ok

	return msg
}

func (msg *Msg) Reply(data string) {
	msg.msgService.reply(msg, data)
}

/// PRIVATE

// Bitverse control messages

func composeHeartbeatMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Heartbeat
	msg.Src = src
	msg.Dst = dst
	msg.ServiceType = Control
	return msg
}

func composeChildrenRequestMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Src = src
	msg.Dst = dst
	msg.ServiceType = Control
	return msg
}

func composeChildrenReplyMsg(src string, dst string, json string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = json
	msg.Src = src
	msg.Dst = dst
	msg.ServiceType = Control
	return msg
}

func composeChildJoin(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildJoined
	msg.Payload = childId
	msg.Src = src
	msg.ServiceType = Control
	return msg
}

func composeChildLeft(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildLeft
	msg.Payload = childId
	msg.Src = src
	msg.ServiceType = Control
	return msg
}

func composeHandshakeMsg(src string) *Msg {
	msg := new(Msg)
	msg.Type = Handshake
	msg.Src = src
	msg.ServiceType = Control
	return msg
}

func getSeqNr() int {
	mutex.Lock()
	seqNrCounter++
	mutex.Unlock()
	return seqNrCounter
}
