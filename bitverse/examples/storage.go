package main

import (
	"fmt"
	"mdc/bitverse"
)

var repoId = "test"                             // uuid
var secret = "x very very very very secret key" // aes encryption key, 16, 24, or 32 bytes
// SERVICE OBSERVER

/// BITVERSE OBSERVER

type BitverseObserver struct {
}

func (bitverseObserver *BitverseObserver) OnSiblingJoined(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " joined")
}

func (bitverseObserver *BitverseObserver) OnSiblingLeft(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " left")
}

func (bitverseObserver *BitverseObserver) OnSiblingHeartbeat(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " heartbeat")
}

func (bitverseObserver *BitverseObserver) OnChildrenReply(node *bitverse.EdgeNode, id string, children []string) {
	fmt.Println("-> received children list from " + id)
}

func (bitverseObserver *BitverseObserver) OnConnected(node *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("-> now connected to super node " + remoteNode.Id())

	storageService := node.CreateStorageService(secret, repoId)

	storageService.Store("myKey", "myValue", 10, func(success bool, reply *bitverse.Msg) {
		if success {
			fmt.Println("managed to store key in bitverse network")
		} else {
			fmt.Println("failed to store key in bitverse network")
		}
	})
}

/// MAIN

func main() {
	var done chan int

	node, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
	//node.Debug()
	fmt.Println("-> my id is " + node.Id())

	node.CreateStorageService(secret, repoId)

	go node.Connect("localhost:1111")

	<-done
}
