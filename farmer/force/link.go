/*  The MIT License (MIT)

Copyright (c) 2014 Luleå University of Technology, Sweden

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE. */

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
