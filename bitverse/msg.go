package bitverse

import (
	"fmt"
	"sync"
)

// Message type definition
const (
	Handshake = iota
	Data
	Heartbeat
	Children
	ChildJoined
	ChildLeft
	Bye
)

// Service type definition
const (
	Messaging = iota
	Storage
	Control
)

// storage service cmd
const (
	Store = iota
	Lookup
	CreateRepo
)

var mutex sync.Mutex
var seqNrCounter int = 0

type Msg struct {
	Type              int
	Payload           string
	Src               string
	Dst               string
	Id                string
	ServiceType       int
	MsgChannelId      string // used by messaging service
	RepositoryId      string // used by storage service
	StorageServiceCmd int    // used by storage service
	Key               string // used by storage service
	Value             string // used by storage service
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
		return "msg[type:data to:" + msg.Dst + " from:" + msg.Src + " payload:" + msg.Payload + " msgchannelid:" + msg.MsgChannelId + "]"
	} else {
		return "msg[type:unkown]"
	}
}

/// PRIVATE

func composeMsgServiceMsg(src string, dst string, serviceId string, payload string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = payload
	msg.Src = src
	msg.Dst = dst
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())
	msg.ServiceType = Messaging
	msg.RepositoryId = ""
	msg.MsgChannelId = serviceId
	return msg
}

func composeStorageServiceStoreMsg(src string, superNodeId string, repoId string, key string, value string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = ""
	msg.Src = src
	msg.Dst = superNodeId
	msg.Key = key
	msg.Value = value
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())
	msg.ServiceType = Storage
	msg.RepositoryId = repoId
	msg.StorageServiceCmd = Store
	msg.MsgChannelId = "internal"
	return msg
}

func composeStorageServiceLookupMsg(src string, superNodeId string, repoId string, key string) *Msg {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = ""
	msg.Src = src
	msg.Dst = superNodeId
	msg.Key = key
	msg.Value = ""
	msg.Id = msg.Src + ":" + fmt.Sprintf("%d", getSeqNr())
	msg.ServiceType = Storage
	msg.RepositoryId = repoId
	msg.StorageServiceCmd = Store
	msg.MsgChannelId = "internal"
	return msg
}

func composeHeartbeatMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Heartbeat
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Id = ""
	msg.ServiceType = Control
	msg.RepositoryId = ""
	msg.MsgChannelId = ""
	return msg
}

func composeChildrenRequestMsg(src string, dst string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = ""
	msg.Src = src
	msg.Dst = dst
	msg.Id = ""
	msg.ServiceType = Control
	msg.RepositoryId = ""
	msg.MsgChannelId = ""
	return msg
}

func composeChildrenReplyMsg(src string, dst string, json string) *Msg {
	msg := new(Msg)
	msg.Type = Children
	msg.Payload = json
	msg.Src = src
	msg.Dst = dst
	msg.Id = ""
	msg.ServiceType = Control
	msg.RepositoryId = ""
	msg.MsgChannelId = ""
	return msg
}

func composeChildJoin(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildJoined
	msg.Payload = childId
	msg.Src = src
	msg.Dst = ""
	msg.Id = ""
	msg.ServiceType = Control
	msg.RepositoryId = ""
	msg.MsgChannelId = ""
	return msg
}

func composeChildLeft(src string, childId string) *Msg {
	msg := new(Msg)
	msg.Type = ChildLeft
	msg.Payload = childId
	msg.Src = src
	msg.Dst = ""
	msg.Id = ""
	msg.ServiceType = Control
	msg.RepositoryId = ""
	msg.MsgChannelId = ""

	return msg
}

func composeHandshakeMsg(src string) *Msg {
	msg := new(Msg)
	msg.Type = Handshake
	msg.Payload = ""
	msg.Src = src
	msg.Dst = ""
	msg.Id = ""
	msg.ServiceType = Control
	msg.RepositoryId = ""
	msg.MsgChannelId = ""
	return msg
}

func getSeqNr() int {
	mutex.Lock()
	seqNrCounter++
	mutex.Unlock()
	return seqNrCounter
}
