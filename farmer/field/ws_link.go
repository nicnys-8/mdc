package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
)

type LinkState int

const (
	Alive = iota
	Dead
)

type Link struct {
	node         *Node
	ws           *websocket.Conn
	remoteNodeId NodeId
	localNodeId  NodeId
	state        LinkState
}

func NewLink(node *Node, ws *websocket.Conn, remoteNodeId NodeId) *Link {
	link := new(Link)
	link.node = node
	link.ws = ws
	link.remoteNodeId = remoteNodeId
	link.localNodeId = node.id
	link.state = Alive

	return link
}

func (link *Link) send(msg *Msg) {
	enc := json.NewEncoder(link.ws)
	err := enc.Encode(msg)

	if err != nil {
		fmt.Printf("Link: detecting broken link <" + string(link.remoteNodeId) + ">\n")
		link.state = Dead
		link.node.linkChannel <- link // notify the node so it can remove it
	}
}
