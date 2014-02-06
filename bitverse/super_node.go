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

type SuperNode struct {
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

func makeSuperNode(transport Transport, localAddress string, localPort string) (*SuperNode, chan int) {
	snode := new(SuperNode)

	snode.localPort = localPort
	snode.links = make(map[NodeId]*Link)
	snode.transport = transport

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	snode.id = NodeId(u.String())

	fmt.Println("SuperNode: my node id is " + snode.Id())

	done := make(chan int)
	snode.msgChannel = make(chan Msg)
	snode.linkChannel = make(chan *Link)

	// initialize transport
	snode.transport.SetLinkChannel(snode.linkChannel)
	snode.transport.SetMsgChannel(snode.msgChannel)
	snode.transport.SetLocalNodeId(snode.id)
	snode.transport.CreateLocalEndPoint(localAddress, localPort)

	go func() {
		for {
			select {
			case msg := <-snode.msgChannel:
				if msg.Dst == snode.id && msg.Type == Data {
					fmt.Println("SuperNode: RECEIVING a DATA message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
				} else {
					fmt.Println("SuperNode: FORWARDING message: <" + msg.Payload + "> from " + string(msg.Src) + " to " + string(msg.Dst))
					snode.Forward(msg)
				}

				fmt.Println("")
			case link := <-snode.linkChannel:
				if link.state == Dead {
					delete(snode.links, link.remoteNodeId)
					fmt.Println("SuperNode: REMOVING link: <" + link.remoteNodeId + ">\n")
				} else {
					fmt.Println("SuperNode: ADDING link: <" + link.remoteNodeId + ">\n")
					snode.links[link.remoteNodeId] = link
				}
			}
		}
	}()

	return snode, done
}

func (snode *SuperNode) Id() NodeId {
	return snode.id
}

func (snode *SuperNode) SendStrToNode(str string, dst NodeId) {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = str
	msg.Src = snode.id
	msg.Dst = dst
	msg.Id = MsgId(string(snode.id) + "-" + strconv.Itoa(snode.seqNumberCounter))

	snode.seqNumberCounter++

	// send the message on all our links
	for _, link := range snode.links {
		link.send(msg)
	}
}

func (snode *SuperNode) Announce() {
	msg := new(Msg)
	msg.Type = Discovery
	msg.Payload = snode.localAddr + ":" + snode.localPort
	msg.Src = snode.id
	msg.Dst = Broadcast
	msg.Id = MsgId(string(snode.id) + "-" + strconv.Itoa(snode.seqNumberCounter))

	snode.seqNumberCounter++

	for _, link := range snode.links {
		link.send(msg)
	}
}

func (snode *SuperNode) Forward(msg Msg) {
	// send the message on all our links
	for _, link := range snode.links {
		if msg.LastHop != link.remoteNodeId { // do not forward messages to a link where it came from
			msg.LastHop = snode.id
			link.send(&msg)
		}
	}
}

var localFlag = flag.String("local", "", "ip address and port which this super node should bound to, e.g. --local localhost:1111")
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
