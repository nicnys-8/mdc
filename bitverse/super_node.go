package bitverse

import (
	"fmt"
)

type SuperNode struct {
	id               NodeId
	wsServer         *WsServer
	links            map[string]*Link
	msgChannel       chan Msg
	linkChannel      chan *Link
	seqNumberCounter int
	localAddr        string
	localPort        string
	transport        Transport
}

func MakeSuperNode(transport Transport, localAddress string, localPort string) (*SuperNode, chan int) {
	superNode := new(SuperNode)

	superNode.localPort = localPort
	superNode.links = make(map[string]*Link)
	superNode.transport = transport

	superNode.id = generateNodeId()

	done := make(chan int)
	superNode.msgChannel = make(chan Msg)
	superNode.linkChannel = make(chan *Link, 10)

	// initialize transport
	superNode.transport.SetLinkChannel(superNode.linkChannel)
	superNode.transport.SetMsgChannel(superNode.msgChannel)
	superNode.transport.SetLocalNodeId(superNode.id)
	superNode.transport.CreateLocalEndPoint(localAddress, localPort)

	go func() {
		for {
			select {
			case msg := <-superNode.msgChannel:
				if msg.Dst == superNode.id.String() && msg.Type == Data {
					fmt.Println("SuperNode: got DATA message <" + msg.Payload + "> from " + msg.Src + " to " + msg.Dst)

				} else if msg.Type == Announce {
					fmt.Println("SuperNode: got ANNOUNCE message from <" + msg.Src + ">")
					superNode.Forward(msg)
				} else {
					fmt.Println("SuperNode: forwarding message <" + msg.Payload + "> from " + msg.Src + " to " + msg.Dst)
					superNode.Forward(msg)
				}

				fmt.Println("")
			case link := <-superNode.linkChannel:
				if link.state == Dead {
					fmt.Println("SuperNode: REMOVING link: <" + link.remoteNodeId.String() + ">\n")
					delete(superNode.links, link.remoteNodeId.String())
					//fmt.Println(superNode.links)
				} else {
					fmt.Println("SuperNode: ADDING link: <" + link.remoteNodeId.String() + ">\n")
					superNode.links[link.remoteNodeId.String()] = link
					//fmt.Println(superNode.links)
				}
			}
		}
	}()

	return superNode, done
}

func (superNode *SuperNode) Id() NodeId {
	return superNode.id
}

func (superNode *SuperNode) Forward(msg Msg) {
	// send the message on all our links
	for _, link := range superNode.links {
		if msg.Src != link.remoteNodeId.String() { // do not forward messages to a link where it came from
			link.send(&msg)
		}
	}
}
