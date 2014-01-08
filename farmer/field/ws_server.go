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
