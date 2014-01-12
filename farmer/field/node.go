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
	seqNumberCounter int
	localAddr        string
	localPort        string
}

func NewNode(localAddr string, localPort string, nodeAddress string) (*Node, chan int) {
	node := new(Node)

	node.localAddr = localAddr
	node.localPort = localPort
	node.links = make(map[NodeId]*Link)

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
				node.Announce()
			case msg := <-node.msgChannel:
				fmt.Println("Node: got a message id=<" + msg.Id + ">")

				if msg.Dst == node.id && msg.Type == Data {
					fmt.Println("Node: RECEIVING a DATA message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
				} else if msg.Dst == Broadcast && msg.Type == Discovery {
					fmt.Println("Node: RECEIVING a DISCOVERY message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))

					// TODO, learn path to other nodes
					//nodeId := msg.Src
					//nodeAddress = msg.Payload

					fmt.Println("Node: FORWARDING message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))

					node.Forward(msg)
				} else {
					fmt.Println("Node: FORWARDING message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
					node.Forward(msg)
				}

				fmt.Println("")
			case link := <-node.linkChannel:
				if link.state == Dead {
					delete(node.links, link.remoteNodeId)
					fmt.Println("Node: REMOVING link: <" + link.remoteNodeId + ">\n")
				} else {
					fmt.Println("Node: ADDING link: <" + link.remoteNodeId + ">\n")
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
	msg.Id = MsgId(string(node.id) + "-" + strconv.Itoa(node.seqNumberCounter))

	node.seqNumberCounter++

	// send the message on all our links
	for _, link := range node.links {
		link.send(msg)
	}
}

func (node *Node) Announce() {
	msg := new(Msg)
	msg.Type = Discovery
	msg.Payload = node.localAddr + ":" + node.localPort
	msg.Src = node.id
	msg.Dst = Broadcast
	msg.Id = MsgId(string(node.id) + "-" + strconv.Itoa(node.seqNumberCounter))

	node.seqNumberCounter++

	for _, link := range node.links {
		link.send(msg)
	}
}

func (node *Node) Forward(msg Msg) {
	// send the message on all our links
	for _, link := range node.links {
		if msg.LastHop != link.remoteNodeId { // do not forward messages to a link where it came from
			msg.LastHop = node.id
			link.send(&msg)
		}
	}
}

func (node *Node) Test() {
	fmt.Printf("Node.Test: \n")

	go func() {
		counter := 0
		for {
			msgStr := "hi"
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

	// TODO: check that the localFlag format is correct, it should be host:port
	temp := strings.Split(*localFlag, ":")
	localAddr := temp[0]
	localPort := temp[1]

	node, done := NewNode(localAddr, localPort, *joinFlag)
	node.Test()

	if *testHttpServerFlag {
		fmt.Println("Starting a HTTP test server at port 8080")
		log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
	}

	<-done
}
