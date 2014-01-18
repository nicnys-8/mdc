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
	"log"
)

type WsClient struct {
	node *Node
	ws   *websocket.Conn
}

func NewWsClient(node *Node) *WsClient {
	wsClient := new(WsClient)
	wsClient.node = node

	return wsClient
}

func (wsClient *WsClient) connect(nodeAddress string) {
	origin := "http://localhost/"
	url := "ws://" + nodeAddress + "/node"

	var err error
	wsClient.ws, err = websocket.Dial(url, "", origin)
	link := wsClient.handshake()

	if err != nil {
		log.Fatal(err)
	}

	wsClient.node.linkChannel <- link

	for {
		msg := wsClient.receive()

		if msg == nil {
			// TODO: remove the link
			return
		}
		wsClient.node.msgChannel <- *msg
	}
}

func (wsClient *WsClient) send(msg *Msg) {
	enc := json.NewEncoder(wsClient.ws)
	err := enc.Encode(msg)
	if err != nil {
		fmt.Println("WSClient: failed to send message")
	}
}

func (wsClient *WsClient) handshake() *Link {
	msg := Msg{Type: Handshake, Payload: string(wsClient.node.id)}

	wsClient.send(&msg)
	reply := wsClient.receive()

	remoteNodeId := NodeId(reply.Payload)
	link := NewLink(wsClient.node, wsClient.ws, remoteNodeId)

	//fmt.Printf("WsClient.handshake: node " + string(wsClient.node.id) + " is now connected to node " + string(remoteNodeId) + "\n")

	return link
}

func (wsClient *WsClient) receive() *Msg { // TODO: return error instead of nil
	dec := json.NewDecoder(wsClient.ws)
	var err error
	var msg Msg

	err = dec.Decode(&msg)
	if err != nil {
		fmt.Println("WSClient: failed to decode message")
		//log.Fatal("decode error:", err)
		return nil
	}

	//fmt.Printf("WSClient.receive: received: " + msg.Payload + "\n")

	return &msg
}
