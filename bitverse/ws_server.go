package bitverse

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"net/http"
)

type WsServer struct {
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	localNodeId       NodeId
}

func (wsServer *WsServer) WsHandler(ws *websocket.Conn) {
	var err error
	var msg Msg

	for {
		dec := json.NewDecoder(ws)
		err = dec.Decode(&msg)

		if err != nil {
			fmt.Println("wsserver: connection closed")
			break
		}

		if msg.Type == Handshake {
			remoteNodeId := makeNodeIdFromString(msg.Payload)
			remoteNode := makeRemoteNode(wsServer.remoteNodeChannel, ws, wsServer.localNodeId, remoteNodeId)
			wsServer.remoteNodeChannel <- remoteNode

			// send our node id to the remote node so that it can also create a link
			reply := Msg{Type: Handshake, Payload: string(wsServer.localNodeId.String())}
			enc := json.NewEncoder(ws)
			enc.Encode(reply)
		} else {
			wsServer.msgChannel <- msg
		}
	}
}

func makeWsServer(localNodeId NodeId, msgChannel chan Msg, remoteNodeChannel chan *RemoteNode) *WsServer {
	wsServer := new(WsServer)
	wsServer.msgChannel = msgChannel
	wsServer.remoteNodeChannel = remoteNodeChannel
	wsServer.localNodeId = localNodeId

	return wsServer
}

func (wsServer *WsServer) start(port string) {
	fmt.Printf("wsserver: starting a new server at port " + port + "\n")

	http.Handle("/node", websocket.Handler(wsServer.WsHandler))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("wsserver.start: " + err.Error())
	}
}
