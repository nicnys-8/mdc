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
	linkChannel  chan *Link
	ws           *websocket.Conn
	remoteNodeId NodeId
	localNodeId  NodeId
	state        LinkState
}

func makeLink(linkChannel chan *Link, ws *websocket.Conn, localNodeId NodeId, remoteNodeId NodeId) *Link {
	link := new(Link)
	link.linkChannel = linkChannel
	link.ws = ws
	link.remoteNodeId = remoteNodeId
	link.localNodeId = localNodeId
	link.state = Alive

	return link
}

func (link *Link) send(msg *Msg) {
	enc := json.NewEncoder(link.ws)
	err := enc.Encode(msg)

	if err != nil {
		fmt.Printf("Link: detecting broken link <" + link.remoteNodeId.String() + ">\n")
		link.state = Dead
		link.linkChannel <- link // notify the node so it can remove it
		fmt.Println("done noitifying node")
	}
}
