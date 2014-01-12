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
