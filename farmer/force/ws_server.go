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
	"net/http"
)

type WsServer struct {
	node *Node
}

func (wsServer *WsServer) WsHandler(ws *websocket.Conn) {
	var err error
	var msg Msg

	for {
		dec := json.NewDecoder(ws)
		err = dec.Decode(&msg)

		if err != nil {
			fmt.Println("WsServer.WsHandler: connection closed")
			break
		}

		if msg.Type == Handshake {
			remoteNodeId := NodeId(msg.Payload)
			link := NewLink(wsServer.node, ws, remoteNodeId)
			wsServer.node.linkChannel <- link

			//fmt.Printf("WsServer.WsHandler: node " + string(wsServer.node.id) + " is now connected to node " + string(remoteNodeId) + "\n")

			// send our node id to the remote node so that it can also create a link
			reply := Msg{Type: Handshake, Payload: string(wsServer.node.id)}
			enc := json.NewEncoder(ws)
			enc.Encode(reply)
		} else {
			wsServer.node.msgChannel <- msg
		}
	}
}

func NewWsServer(node *Node) *WsServer {
	server := new(WsServer)
	server.node = node

	return server
}

func (wsServer *WsServer) start(port string) {
	fmt.Printf("WSServer.bind: starting a new ws server at port " + port + "\n")

	http.Handle("/node", websocket.Handler(wsServer.WsHandler))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("WsServer.start: " + err.Error())
	}
}
