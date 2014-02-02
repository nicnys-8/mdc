package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"net/http"
)

type WsServer struct {
	msgChannel  chan Msg
	linkChannel chan *Link
	localNodeId NodeId
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
			fmt.Println("XXXXXX remoteNodeId=" + remoteNodeId)
			fmt.Println("msg.Payload=" + msg.Payload)
			link := makeLink(wsServer.linkChannel, ws, wsServer.localNodeId, remoteNodeId)
			wsServer.linkChannel <- link

			// send our node id to the remote node so that it can also create a link
			reply := Msg{Type: Handshake, Payload: string(wsServer.localNodeId)}
			enc := json.NewEncoder(ws)
			enc.Encode(reply)
		} else {
			wsServer.msgChannel <- msg
		}
	}
}

func makeWsServer(localNodeId NodeId, msgChannel chan Msg, linkChannel chan *Link) *WsServer {
	wsServer := new(WsServer)
	wsServer.msgChannel = msgChannel
	wsServer.linkChannel = linkChannel
	wsServer.localNodeId = localNodeId

	return wsServer
}

func (wsServer *WsServer) start(port string) {
	fmt.Printf("WSServer.bind: starting a new ws server at port " + port + "\n")

	http.Handle("/node", websocket.Handler(wsServer.WsHandler))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("WsServer.start: " + err.Error())
	}
}
