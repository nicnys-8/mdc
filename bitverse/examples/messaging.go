package main

import (
	"fmt"
	"mdc/bitverse"
)

/// SERVICE OBSERVER

type MsgServiceObserver struct {
}

func (msgServiceObserver *MsgServiceObserver) OnDeliver(msgService *bitverse.MsgService, msg *bitverse.Msg) {
	if msg.Payload == "hello" {
		fmt.Println("got a message: hello")
		fmt.Println("sending: hi dude!")
		msgService.Send(msg.Src, "hi dude!")

		fmt.Println("sending: how are you doing?")
		msgService.SendAndGetReply(msg.Src, "how are you doing?", 10, func(success bool, reply *string) {
			fmt.Println("got a reply (how are you doing?): " + *reply)
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

func (bitverseObserver *BitverseObserver) OnSiblingJoined(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " joined")

	msgService := node.GetMsgService(serviceId)
	msgService.Send(id, "hello")
	fmt.Println("sending: hello")

	fmt.Println("sending: who are you?")
	msgService.SendAndGetReply(id, "who are you?", 10, func(success bool, reply *string) {
		fmt.Println("got a reply (who are you?): " + *reply)
	})
}

func (bitverseObserver *BitverseObserver) OnSiblingLeft(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " left")
}

func (bitverseObserver *BitverseObserver) OnSiblingHeartbeat(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " heartbeat")
}

func (bitverseObserver *BitverseObserver) OnChildrenReply(node *bitverse.EdgeNode, id string, children []string) {
	fmt.Println("-> received children list from " + id)
	for _, childNodeId := range children {
		fmt.Println("-> learned about a sibling: " + childNodeId)
	}
}

func (bitverseObserver *BitverseObserver) OnConnected(node *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("-> now connected to super node " + remoteNode.Id())

	remoteNode.SendChildrenRequest()

	msgService := node.GetMsgService(serviceId)
	msgService.SendAndGetReply("6a133a1b41f987210559ceb4ed9b1dbf58aec876", "hello", 10, func(success bool, reply *string) {
		if success {
			fmt.Println("that was a surprise " + *reply)
		} else {
			// we will most likely timeout unless node 6a133a1b41f987210559ceb4ed9b1dbf58aec876 is online
			fmt.Println("failed to send message to node with id <does not exists>")
		}
	})
}

// uuid
var serviceId = "6107911a-7554-4ea7-80fc-25ec5e2462a7"

// aes encryption key should be 32 bytes encoded as hex
// can be genetated by calling ./bitverse --generate-aes-secret from unix shell or
// by calling bitverse.GenerateAesSecret()
var secret = "5da71277f031a9dff561f0a72bb72651e260dab0735b767f2f7a62dec9e99760"

/// MAIN

func main() {
	var done chan int

	secret2, _ := bitverse.GenerateAesSecret()
	fmt.Println(secret2)

	node, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
	//node.Debug()
	fmt.Println("-> my id is " + node.Id())

	msgServiceObserver := new(MsgServiceObserver)
	node.CreateMsgService(secret, serviceId, msgServiceObserver)

	go node.Connect("localhost:1111")

	<-done
}
