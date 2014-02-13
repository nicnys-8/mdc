package main

import (
	"fmt"
	"mdc/bitverse"
)

type MyServiceObserver struct {
}

func (myService *MyServiceObserver) OnError(err error) {
}

func (myService *MyServiceObserver) OnDeliver(msg *bitverse.Msg) {
	fmt.Println("got a message " + msg.String())
}

func (myService *MyServiceObserver) OnSiblingJoin(nodeId string) {
	fmt.Println("sibling " + nodeId + " joined")
}

func (myService *MyServiceObserver) OnSiblingExit(nodeId string) {
	fmt.Println("sibling " + nodeId + " exit")
}

func (myService *MyServiceObserver) OnSiblingHeartbeat(nodeId string) {
	fmt.Println("sibling " + nodeId + " heartbeat")
}

func main() {
	var done chan int
	serviceObserver := new(MyServiceObserver)

	transport := bitverse.MakeWSTransport()
	edgeNode, done := bitverse.MakeEdgeNode(transport)
	go edgeNode.Join("localhost:1111")

	myService := edgeNode.GetService("myservice", serviceObserver)
	myService.Nop()

	<-done
}
