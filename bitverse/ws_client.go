package bitverse

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"os"
)

type WsClient struct {
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	localNodeId       NodeId
	ws                *websocket.Conn
}

func makeWsClient(msgChannel chan Msg, remoteNodeChannel chan *RemoteNode, localNodeId NodeId) *WsClient {
	wsClient := new(WsClient)
	wsClient.msgChannel = msgChannel
	wsClient.remoteNodeChannel = remoteNodeChannel
	wsClient.localNodeId = localNodeId

	return wsClient
}

func (wsClient *WsClient) connect(ipAddress string) {
	origin := "http://localhost/"
	url := "ws://" + ipAddress + "/node"

	var err error
	wsClient.ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		info("failed to connect to supernode at " + ipAddress + ", connection refused")
		os.Exit(0)
	}

	remoteNode := wsClient.handshake()

	wsClient.remoteNodeChannel <- remoteNode

	for {
		msg := wsClient.receive()

		if msg == nil {
			// TODO: remove the link
			return
		}
		wsClient.msgChannel <- *msg
	}
}

func (wsClient *WsClient) send(msg *Msg) {
	enc := json.NewEncoder(wsClient.ws)
	err := enc.Encode(msg)
	if err != nil {
		info("wsclient: failed to send message")
	}
}

func (wsClient *WsClient) handshake() *RemoteNode {
	msg := ComposeHandshakeMsg(wsClient.localNodeId.String())

	wsClient.send(msg)
	reply := wsClient.receive()

	remoteNodeId := makeNodeIdFromString(reply.Src)
	remoteNode := makeRemoteNode(wsClient.remoteNodeChannel, wsClient.ws, wsClient.localNodeId.String(), remoteNodeId.String())

	return remoteNode
}

func (wsClient *WsClient) receive() *Msg { // TODO: return error instead of nil
	dec := json.NewDecoder(wsClient.ws)
	var err error
	var msg Msg

	err = dec.Decode(&msg)
	if err != nil {
		info("wsclient: failed to decode message")
		return nil
	}

	return &msg
}
