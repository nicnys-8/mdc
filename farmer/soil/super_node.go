package main

import (
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	transport        Transport
}

func makeSuperNode(transport Transport, localAddress string, localPort string) (*Node, chan int) {
	node := new(Node)

	node.localPort = localPort
	node.links = make(map[NodeId]*Link)
	node.transport = transport

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	node.id = NodeId(u.String())

	fmt.Println("Node: my node id is " + node.Id())

	done := make(chan int)
	node.msgChannel = make(chan Msg)
	node.linkChannel = make(chan *Link)

	// initialize transport
	node.transport.SetLinkChannel(node.linkChannel)
	node.transport.SetMsgChannel(node.msgChannel)
	node.transport.SetLocalNodeId(node.id)
	node.transport.CreateLocalEndPoint(localAddress, localPort)

	go func() {
		for {
			select {
			case msg := <-node.msgChannel:
				if msg.Dst == node.id && msg.Type == Data {
					fmt.Println("Node: RECEIVING a DATA message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
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

var localFlag = flag.String("local", "", "ip address and port which this node should bound to, e.g. --local localhost:1111")
var testHttpServerFlag = flag.Bool("test-http-server", false, "starts a http test server at port 8080 for debuging")

func main() {
	flag.Parse()

	// TODO: check that the localFlag format is correct, it should be host:port
	temp := strings.Split(*localFlag, ":")
	localAddr := temp[0]
	localPort := temp[1]

	transport := makeWSTransport()
	_, done := makeSuperNode(transport, localAddr, localPort)

	if *testHttpServerFlag {
		fmt.Println("Starting a HTTP test server at port 8080")
		log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
	}

	<-done
}
