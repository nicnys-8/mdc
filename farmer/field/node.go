package main

import (
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type NodeId string

const Broadcast NodeId = NodeId("BROADCAST")

type Node struct {
	id               NodeId
	wsServer         *WsServer
	links            map[NodeId]*Link
	msgChannel       chan Msg
	linkChannel      chan *Link
	seqNumbers       map[NodeId]int // contains the higest sequence number received from a node
	seqNumberCounter int
}

func NewNode(port string, nodeAddress string) (*Node, chan int) {
	node := new(Node)

	node.links = make(map[NodeId]*Link)
	node.seqNumbers = make(map[NodeId]int)
	node.seqNumberCounter = 0

	//fmt.Printf("\n\bXXXX node.seqNumberCounter=%d\n\n", node.seqNumberCounter)

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	node.id = NodeId(u.String())

	fmt.Println("NewNode: my node id is " + node.Id())

	done := make(chan int)
	node.msgChannel = make(chan Msg)
	node.linkChannel = make(chan *Link)

	// start websocket server
	node.wsServer = NewWsServer(node)
	go node.wsServer.start(port)

	// join another node
	if nodeAddress != "" {
		fmt.Println("NewNode: joining node at " + nodeAddress)
		go node.Join(nodeAddress)
	}

	go func() {
		for {
			select {
			case msg := <-node.msgChannel:
				// check if this is a new message, or if we have already got it before
				seqNumberForNode, found := node.seqNumbers[node.id]
				seqNr := msg.SeqNr

				//fmt.Printf("\n\bGOT A MSG seqNumberForNode=%d\n\n", seqNumberForNode)
				//fmt.Printf("\n\bGOT A MSG seqNumberForMsg=%d\n\n", seqNr)

				if !found {
					seqNumberForNode = seqNr - 1 // this is a msg from a new node, assume we got the previous msg
					fmt.Printf("node <"+string(node.id)+">  not found setting seq nr to %d\n", msg.SeqNr-1)
				}

				if seqNr == seqNumberForNode+1 {
					// this a new message
					node.seqNumbers[node.id] = seqNumberForNode + 1 // increase the counter

					if msg.Dst == node.id {
						// yes, we finnaly got a message
						fmt.Println("NewNode: got a new message: " + msg.Payload)
					} else {
						// forward this message
						fmt.Println("NewNode: forward message: " + msg.Payload)
						node.Forward(msg)
					}

				} else {
					// we already go this message, just drop it
					fmt.Printf("NewNode: dropping message")
				}

			case link := <-node.linkChannel:
				if link.state == Dead {
					delete(node.links, link.remoteNodeId)
					fmt.Println("NewNode: removing link: <" + link.remoteNodeId + ">\n")
				} else {
					fmt.Println("NewNode: got a new link: <" + link.remoteNodeId + ">\n")
					node.links[link.remoteNodeId] = link
				}
			}
		}
	}()

	return node, done

}

func (node *Node) Id() NodeId {
	return node.id
}

func (node *Node) Join(nodeAddress string) {
	wsClient := NewWsClient(node)
	wsClient.connect(nodeAddress)
}

func (node *Node) SendStrToNode(str string, dst NodeId) {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = str
	msg.Src = node.id
	msg.Dst = dst
	msg.SeqNr = node.seqNumberCounter

	node.seqNumberCounter++

	// send the message on all our links
	for _, link := range node.links {
		//fmt.Printf("\nsending msg msg.SeqNr=<%d>\n", msg.SeqNr)
		link.send(msg)
	}
}

func (node *Node) Forward(msg Msg) {
	// send the message on all our links
	for _, link := range node.links {
		//fmt.Printf("\nsending msg msg.SeqNr=<%d>\n", msg.SeqNr)
		link.send(&msg)
	}
}

func (node *Node) Test() {
	fmt.Printf("Node.Test: \n")

	go func() {
		counter := 0
		for {
			msgStr := "hejsan " + strconv.Itoa(counter) + " from " + string(node.Id())
			fmt.Printf("Sending MSG <" + msgStr + ">\n")
			time.Sleep(2 * time.Second)
			for remoteNodeId, _ := range node.links {
				node.SendStrToNode(msgStr, remoteNodeId)
			}
			counter++
		}
	}()
}

var nodeAddressFlag = flag.String("join", "", "ip address and port to a node to join, e.g. --join localhost:3121")
var nodePortFlag = flag.String("port", "12345", "port to bind this node to, e.g --port 12345")
var testHttpServerFlag = flag.Bool("test-http-server", false, "starts a http test server at port 8080 for debuging")

func main() {
	flag.Parse()

	node, done := NewNode(*nodePortFlag, *nodeAddressFlag)
	node.Test()

	if *testHttpServerFlag {
		fmt.Println("Starting a HTTP test server at port 8080")
		log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
	}

	<-done
}
