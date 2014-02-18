package main

import (
	"flag"
	"fmt"
	"log"
	"mdc/bitverse"
	"net/http"
	"strings"
)

var superFlag = flag.Bool("super", false, "run the node as a super node")
var debugFlag = flag.Bool("debug", false, "run the node in debug mode")
var localFlag = flag.String("local", "", "ip address and port which this super node should bound to, e.g. --local localhost:1111")
var joinFlag = flag.String("join", "", "ip address and port to a node to join, e.g. --join localhost:2222")
var testHttpServerFlag = flag.Bool("test-http-server", false, "starts a http test server at port 8080 for debuging")

/// BITVERSE OBSERVER

type MyBitverseObserver struct {
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingJoined(edgeNode *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " joined")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingLeft(edgeNode *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " left")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingHeartbeat(edgeNode *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " heartbeat")
}

func (myBitverseObserver *MyBitverseObserver) OnChildrenReply(edgeNode *bitverse.EdgeNode, nodeId string, children []string) {
	fmt.Println("received children list from " + nodeId)
}

func (myBitverseObserver *MyBitverseObserver) OnConnected(edgeNode *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("now connected to super node " + remoteNode.Id())
}

/// MAIN

func main() {
	flag.Parse()

	bitverseObserver := new(MyBitverseObserver)
	var done chan int

	transport := bitverse.MakeWSTransport()

	if *superFlag {
		var superNode *bitverse.SuperNode
		temp := strings.Split(*localFlag, ":")
		localAddr := temp[0]
		localPort := temp[1]

		superNode, done = bitverse.MakeSuperNode(transport, localAddr, localPort)

		if *debugFlag {
			superNode.Debug()
		}

		if *testHttpServerFlag {
			fmt.Println("Starting a HTTP test server at port 8080")
			log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
		}
	} else {
		var edgeNode *bitverse.EdgeNode
		edgeNode, done = bitverse.MakeEdgeNode(transport, bitverseObserver)

		// join super node
		remoteAddress := *joinFlag
		if remoteAddress != "" {
			fmt.Println("EdgeNode: joining node at " + remoteAddress)
			go edgeNode.Connect(remoteAddress)
		}
	}

	<-done
}
