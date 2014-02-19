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
	if msg.Payload == "hello" {
		fmt.Println("got a message: hello")
		fmt.Println("sending: hi dude!")
		msgService.Send(msg.Src, "hi dude!")

		fmt.Println("sending: how are you doing?")
		msgService.SendAndGetReply(msg.Src, "how are you doing?", 10, func(timedOut bool, reply *bitverse.Msg) {
			fmt.Println("got a reply (how are you doing?): " + reply.Payload)
		})
	} else if msg.Payload == "how are you doing?" {
		fmt.Println("got a message: how are you doing?")
		fmt.Println("sending reply (how are you doing): i am fine")
		msgService.Reply(msg, "i am fine")
	} else if msg.Payload == "hi dude!" {
		fmt.Println("got a message: hi dude!")
	} else if msg.Payload == "who are you?" {
		fmt.Println("got a message: who are you?")
		fmt.Println("sending reply (who are you?): i am joker")
		msgService.Reply(msg, "i am joker")
	} else {
		fmt.Println("ERROR got a msg: " + msg.Payload)
	}
}

/// BITVERSE OBSERVER

type BitverseObserver struct {
}

func (bitverseObserver *BitverseObserver) OnSiblingJoined(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " joined")

	msgService := edgeNode.GetMsgService(serviceId)
	msgService.Send(id, "hello")
	fmt.Println("sending: hello")

	fmt.Println("sending: who are you?")
	msgService.SendAndGetReply(id, "who are you?", 10, func(timedOut bool, reply *bitverse.Msg) {
		fmt.Println("got a reply (who are you?): " + reply.Payload)
	})
}

func (bitverseObserver *BitverseObserver) OnSiblingLeft(edgeNode *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " left")
}

func (bitverseObserver *BitverseObserver) OnSiblingHeartbeat(edgeNode *bitverse.EdgeNode, id string) {
	//fmt.Println("-> sibling " + id + " heartbeat")
}

func (bitverseObserver *BitverseObserver) OnChildrenReply(edgeNode *bitverse.EdgeNode, id string, children []string) {
	fmt.Println("-> received children list from " + id)
	for _, childNodeId := range children {
		fmt.Println("-> learned about a sibling: " + childNodeId)
	}
}

func (bitverseObserver *BitverseObserver) OnConnected(edgeNode *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("-> now connected to super node " + remoteNode.Id())

	remoteNode.SendChildrenRequest()

	msgService := edgeNode.GetMsgService(serviceId)
	msgService.SendAndGetReply("does not exists", "", 10, func(timedOut bool, reply *bitverse.Msg) {
		if timedOut {
			fmt.Println("failed to send message to node with id <does not exists>")
		} else {
			fmt.Println("that was a surprise " + reply.Payload)
		}
	})
}

/// MAIN

func main() {
	var done chan int

	edgeNode, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
	//edgeNode.Debug()
	fmt.Println("-> my id is " + edgeNode.Id())

	msgServiceObserver := new(MsgServiceObserver)
	edgeNode.CreateMsgService(secret, serviceId, msgServiceObserver)

	go edgeNode.Connect("localhost:1111")

	<-done
}
