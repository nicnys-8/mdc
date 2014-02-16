package main

import (
	"fmt"
	"mdc/bitverse"
)

/// SERVICE OBSERVER

type MyServiceObserver struct {
}

func (myService *MyServiceObserver) OnError(err error) {
}

func (myService *MyServiceObserver) OnDeliver(msg *bitverse.Msg) {
	fmt.Println("got a message " + msg.String())
}

/// BITVERSE OBSERVER

type MyBitverseObserver struct {
}

func (myBitverseObserver *MyBitverseObserver) OnError(err error) {
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingJoin(nodeId string) {
	fmt.Println("sibling " + nodeId + " joined")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingExit(nodeId string) {
	fmt.Println("sibling " + nodeId + " exit")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingHeartbeat(nodeId string) {
	fmt.Println("sibling " + nodeId + " heartbeat")
}

func (myBitverseObserver *MyBitverseObserver) OnChildrenReply(nodeId string) {
	fmt.Println("received children list from " + nodeId)
}

func (myBitverseObserver *MyBitverseObserver) OnConnected(edgeNode *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("now connected to super node " + remoteNode.Id.String())
	remoteNode.SendChildrenRequest()

	serviceObserver := new(MyServiceObserver)

	myService := edgeNode.GetService("myservice", serviceObserver)
	myService.Nop()
}

/// MAIN

func main() {
	var done chan int

	bitverseObserver := new(MyBitverseObserver)

	transport := bitverse.MakeWSTransport()
	edgeNode, done := bitverse.MakeEdgeNode(transport, bitverseObserver)

	go edgeNode.Connect("localhost:1111")

	<-done
}
