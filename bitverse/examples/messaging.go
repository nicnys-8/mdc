package main

import (
	"fmt"
	"mdc/bitverse"
)

var serviceId = "6107911a-7554-4ea7-80fc-25ec5e2462a7" // uuid
var secret = "x very very very very secret key"        // aes encryption key, 16, 24, or 32 bytes

/// SERVICE OBSERVER

type MsgServiceObserver struct {
}

func (msgServiceObserver *MsgServiceObserver) OnDeliver(msgService *bitverse.MsgService, msg *bitverse.Msg) {
	fmt.Println("got a message <" + msg.Payload + ">" + " from " + msg.Src)

	// send a reply back
	if msg.Payload == "hello" {
		msgService.Send(msg.Src, "hi dude!")
	}

	if msg.Payload == "how are you?" {
		fmt.Println("XXXXXXXX sending reply back xxxx")
		msgService.Reply(msg, "i am fine")
	}
}

/// BITVERSE OBSERVER

type BitverseObserver struct {
}

func (bitverseObserver *BitverseObserver) OnSiblingJoined(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("sibling " + id + " joined")
}

func (bitverseObserver *BitverseObserver) OnSiblingLeft(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("sibling " + id + " left")
}

func (bitverseObserver *BitverseObserver) OnSiblingHeartbeat(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("sibling " + id + " heartbeat")
}

func (bitverseObserver *BitverseObserver) OnChildrenReply(edgeNode *bitverse.EdgeNode, id string, children []string) {
	fmt.Println("received children list from " + id)
	for _, childNodeId := range children {
		fmt.Println("learned about a child: " + childNodeId)

		msgService := edgeNode.GetMsgService(serviceId)
		msgService.Send(childNodeId, "hello")

		msgService.SendAndGetReply(childNodeId, "how are you?", 10, func(timedOut bool, reply *bitverse.Msg) {
			fmt.Println("got a reply: " + reply.Payload)
		})
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
	//edgeNode.Debug()
	fmt.Println("my id is " + edgeNode.Id())

	msgServiceObserver := new(MsgServiceObserver)
	edgeNode.CreateMsgService(secret, serviceId, msgServiceObserver)

	go edgeNode.Connect("localhost:1111")

	<-done
}
