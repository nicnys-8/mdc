package bitverse

import (
	"encoding/json"
	"io"
)

type RemoteNodeState int

const (
	Alive = iota
	Dead
)

type RemoteNode struct {
	remoteNodeChannel chan *RemoteNode
	writer            io.Writer
	id                string
	remoteId          string
	state             RemoteNodeState
}

func makeRemoteNode(remoteNodeChannel chan *RemoteNode, writer io.Writer, remoteId string, id string) *RemoteNode {
	remoteNode := new(RemoteNode)
	remoteNode.remoteNodeChannel = remoteNodeChannel
	remoteNode.writer = writer
	remoteNode.id = id
	remoteNode.remoteId = remoteId
	remoteNode.state = Alive

	return remoteNode
}

func (remoteNode *RemoteNode) SendChildrenRequest() {
	msg := ComposeChildrenRequestMsg(remoteNode.remoteId, remoteNode.id)
	remoteNode.send(msg)
}

func (remoteNode *RemoteNode) Id() string {
	return remoteNode.id
}

func (remoteNode *RemoteNode) RemoteId() string {
	return remoteNode.remoteId
}

/// PRIVATE

func (remoteNode *RemoteNode) send(msg *Msg) {
	enc := json.NewEncoder(remoteNode.writer)
	err := enc.Encode(msg)

	if err != nil {
		remoteNode.state = Dead
		debug("link: detecting dead link")
		remoteNode.remoteNodeChannel <- remoteNode // notify the node so it can remove it
	}
}
