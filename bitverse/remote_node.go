package bitverse

import (
	"encoding/json"
	"fmt"
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
	Id                NodeId
	localNodeId       NodeId
	state             RemoteNodeState
}

func makeRemoteNode(remoteNodeChannel chan *RemoteNode, writer io.Writer, localNodeId NodeId, remoteNodeId NodeId) *RemoteNode {
	remoteNode := new(RemoteNode)
	remoteNode.remoteNodeChannel = remoteNodeChannel
	remoteNode.writer = writer
	remoteNode.Id = remoteNodeId
	remoteNode.localNodeId = localNodeId
	remoteNode.state = Alive

	return remoteNode
}

func (remoteNode *RemoteNode) SendChildrenRequest() {
	msg := ComposeChildrenRequestMsg(remoteNode.localNodeId.String(), remoteNode.Id.String())
	remoteNode.send(msg)
}

/// PRIVATE

func (remoteNode *RemoteNode) send(msg *Msg) {
	enc := json.NewEncoder(remoteNode.writer)
	err := enc.Encode(msg)

	if err != nil {
		remoteNode.state = Dead
		fmt.Println("link: detecting dead link")
		remoteNode.remoteNodeChannel <- remoteNode // notify the node so it can remove it
	}
}
