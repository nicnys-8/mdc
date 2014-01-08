package main

import (
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	localAddr        string
	localPort        string
}

func NewNode(localAddr string, localPort string, nodeAddress string) (*Node, chan int) {
	node := new(Node)

	node.localAddr = localAddr
	node.localPort = localPort
	node.links = make(map[NodeId]*Link)
	node.seqNumbers = make(map[NodeId]int)
	node.seqNumberCounter = 0 // this one is shared my multiple go routines and may need to be protected

	//fmt.Printf("\n\bXXXX node.seqNumberCounter=%d\n\n", node.seqNumberCounter)

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	node.id = NodeId(u.String())

	fmt.Println("Node: my node id is " + node.Id())

	done := make(chan int)
	node.msgChannel = make(chan Msg)
	node.linkChannel = make(chan *Link)

	// start websocket server
	node.wsServer = NewWsServer(node)
	go node.wsServer.start(localPort)

	// join another node
	if nodeAddress != "" {
		fmt.Println("Node: joining node at " + nodeAddress)
		go node.Join(nodeAddress)
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("Node: sending announce message")
				node.Announce()
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

					if msg.Dst == node.id || msg.Dst == Broadcast {
						// yes, we finnaly got a message
						fmt.Println("Node: got a new message: " + msg.Payload)
					} else {
						// forward this message
						fmt.Println("Node: forward message: " + msg.Payload)
						node.Forward(msg)
					}

				} else {
					// we already go this message, just drop it
					fmt.Println("Node: dropping message " + msg.Payload)
				}

			case link := <-node.linkChannel:
				if link.state == Dead {
					delete(node.links, link.remoteNodeId)
					fmt.Println("Node: removing link: <" + link.remoteNodeId + ">\n")
				} else {
					fmt.Println("Node: got a new link: <" + link.remoteNodeId + ">\n")
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

func (node *Node) Announce() {
	msg := new(Msg)
	msg.Type = Discovery
	msg.Payload = node.localAddr + ":" + node.localPort
	msg.Src = node.id
	msg.Dst = Broadcast
	msg.SeqNr = node.seqNumberCounter

	node.seqNumberCounter++

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

var joinFlag = flag.String("join", "", "ip address and port to a node to join, e.g. --join localhost:2222")
var localFlag = flag.String("local", "", "ip address and port where this node should bound, e.g. --local localhost:1111")
var testHttpServerFlag = flag.Bool("test-http-server", false, "starts a http test server at port 8080 for debuging")

func main() {
	flag.Parse()

	temp := strings.Split(*localFlag, ":")
	localAddr := temp[0]
	localPort := temp[1]
	//localPortStr = strings.Trim(localPortStr, " ")
	//localPort, _ := strconv.Atoi(localPortStr)

	//fmt.Println("localAddr=" + localAddr)
	//fmt.Printf("localPort=%d\n", localPort)

	node, done := NewNode(localAddr, localPort, *joinFlag)
	node.Test()

	if *testHttpServerFlag {
		fmt.Println("Starting a HTTP test server at port 8080")
		log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
	}

	<-done
}
