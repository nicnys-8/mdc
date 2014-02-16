package main

import (
	"fmt"
	"mdc/bitverse"
)

var serviceId = "6107911a-7554-4ea7-80fc-25ec5e2462a7"
var secret = "x very very very very secret key" // 16, 24, or 32 bytes

/// SERVICE OBSERVER

type MsgServiceObserver struct {
}

func (msgServiceObserver *MsgServiceObserver) OnDeliver(msgService *bitverse.MsgService, msg *bitverse.Msg) {
	fmt.Println("got a message <" + msg.Payload + ">" + " from " + msg.Src)

	// send a reply back
	if msg.Payload == "hello" {
		msgService.Send(msg.Src, "hi dude")
	}

}

/// BITVERSE OBSERVER

type BitverseObserver struct {
}

func (bitverseObserver *BitverseObserver) OnSiblingJoin(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("sibling " + id + " joined")
}

func (bitverseObserver *BitverseObserver) OnSiblingExit(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("sibling " + id + " exit")
}

func (bitverseObserver *BitverseObserver) OnSiblingHeartbeat(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("sibling " + id + " heartbeat")
}

func (bitverseObserver *BitverseObserver) OnChildrenReply(edgeNode *bitverse.EdgeNode, id string, children []string) {
	fmt.Println("received children list from " + id)
	for _, childNodeId := range children {
		fmt.Println("child: " + childNodeId)

		msgService := edgeNode.GetMsgService(serviceId)
		msgService.Send(childNodeId, "hello")
	}
}

func (bitverseObserver *BitverseObserver) OnConnected(edgeNode *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("now connected to super node " + remoteNode.Id())

	remoteNode.SendChildrenRequest()
}

/// MAIN

func main() {
	var done chan int

	edgeNode, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))

	msgServiceObserver := new(MsgServiceObserver)
	edgeNode.CreateMsgService(secret, serviceId, msgServiceObserver)

	go edgeNode.Connect("localhost:1111")

	<-done
}
