package main

import (
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"log"
	"strconv"
	"time"
)

type NodeId string

type Node struct {
	id               NodeId
	wsServer         *WsServer
	coreNodeLink     *Link
	msgChannel       chan Msg
	linkChannel      chan *Link
	seqNumberCounter int
	transport        Transport
}

func makeNode(transport Transport) (*Node, chan int) {
	node := new(Node)
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

	go func() {
		for {
			select {
			case msg := <-node.msgChannel:
				if msg.Dst == node.id && msg.Type == Data {
					fmt.Println("Node: RECEIVING a DATA message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
				} else {
					fmt.Println("Node: ERROR received a message with wrong destination: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
				}
			case link := <-node.linkChannel:
				if link.state == Dead {
					fmt.Println("Node: ERROR we just lost our connection to the core node <" + link.remoteNodeId + ">\n")
					node.coreNodeLink = nil
				} else {
					fmt.Println("Node: adding link to core node <" + link.remoteNodeId + ">\n")
					node.coreNodeLink = link
				}
			}
		}
	}()

	return node, done
}

func (node *Node) Id() NodeId {
	return node.id
}

func (node *Node) Join(remoteAddress string) {
	node.transport.ConnectRemoteEndPoint(remoteAddress)
}

func (node *Node) Send(str string, dst NodeId) {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = str
	msg.Src = node.id
	msg.Dst = dst
	msg.Id = MsgId(string(node.id) + "-" + strconv.Itoa(node.seqNumberCounter))

	node.seqNumberCounter++
	node.coreNodeLink.send(msg)
}

func (node *Node) Test() {
	fmt.Printf("Node.Test: \n")

	go func() {
		counter := 0
		for {
			msgStr := "hi"
			time.Sleep(2 * time.Second)
			if node.coreNodeLink != nil {
				node.Send(msgStr, node.coreNodeLink.remoteNodeId)
			}
			counter++
		}
	}()
}

var joinFlag = flag.String("join", "", "ip address and port to a node to join, e.g. --join localhost:2222")

func main() {
	flag.Parse()

	transport := makeWSTransport()
	node, done := makeNode(transport)
	node.Test()

	// join super node
	remoteAddress := *joinFlag
	if remoteAddress != "" {
		fmt.Println("Node: joining node at " + remoteAddress)
		go node.Join(remoteAddress)
	}

	<-done
}
