package bitverse

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"log"
)

type WsClient struct {
	msgChannel  chan Msg
	linkChannel chan *Link
	localNodeId NodeId
	ws          *websocket.Conn
}

func makeWsClient(msgChannel chan Msg, linkChannel chan *Link, localNodeId NodeId) *WsClient {
	wsClient := new(WsClient)
	wsClient.msgChannel = msgChannel
	wsClient.linkChannel = linkChannel
	wsClient.localNodeId = localNodeId

	return wsClient
}

func (wsClient *WsClient) connect(ipAddress string) {
	origin := "http://localhost/"
	url := "ws://" + ipAddress + "/node"

	var err error
	wsClient.ws, err = websocket.Dial(url, "", origin)
	link := wsClient.handshake()

	if err != nil {
		log.Fatal(err)
	}

	wsClient.linkChannel <- link

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
		fmt.Println("WSClient: failed to send message")
	}
}

func (wsClient *WsClient) handshake() *Link {
	msg := Msg{Type: Handshake, Payload: wsClient.localNodeId.String()}

	wsClient.send(&msg)
	reply := wsClient.receive()

	remoteNodeId := makeNodeIdFromString(reply.Payload)
	link := makeLink(wsClient.linkChannel, wsClient.ws, wsClient.localNodeId, remoteNodeId)

	return link
}

func (wsClient *WsClient) receive() *Msg { // TODO: return error instead of nil
	dec := json.NewDecoder(wsClient.ws)
	var err error
	var msg Msg

	err = dec.Decode(&msg)
	if err != nil {
		fmt.Println("WSClient: failed to decode message")
		return nil
	}

	return &msg
}
