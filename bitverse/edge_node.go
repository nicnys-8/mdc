package bitverse

import (
	"fmt"
	"time"
)

type EdgeNode struct {
	id               NodeId
	wsServer         *WsServer
	superNodeLink    *Link
	msgChannel       chan Msg
	linkChannel      chan *Link
	seqNumberCounter int
	transport        Transport
}

func MakeEdgeNode(transport Transport) (*EdgeNode, chan int) {
	edgeNode := new(EdgeNode)
	edgeNode.transport = transport
	edgeNode.id = generateNodeId()

	done := make(chan int)
	edgeNode.msgChannel = make(chan Msg)
	edgeNode.linkChannel = make(chan *Link)

	// initialize transport
	edgeNode.transport.SetLinkChannel(edgeNode.linkChannel)
	edgeNode.transport.SetMsgChannel(edgeNode.msgChannel)
	edgeNode.transport.SetLocalNodeId(edgeNode.id)

	go func() {
		for {
			select {
			case msg := <-edgeNode.msgChannel:
				if msg.Dst == edgeNode.id.String() && msg.Type == Data {
					fmt.Println("EdgeNode: got DATA message <" + msg.Payload + "> from " + msg.Src + " to " + msg.Dst)
				} else if msg.Type == Announce {
					fmt.Println("SuperNode: got ANNOUNCE message from <" + msg.Src + ">")
				} else {
					fmt.Println("EdgeNode: ERROR got a very strange message <" + msg.Payload + "> from " + msg.Src + " to " + msg.Dst)
				}
			case link := <-edgeNode.linkChannel:
				if link.state == Dead {
					fmt.Println("EdgeNode: ERROR we just lost our connection to the super node <" + link.remoteNodeId.String() + ">\n")
					edgeNode.superNodeLink = nil
				} else {
					fmt.Println("EdgeNode: adding link to super node <" + link.remoteNodeId.String() + ">\n")
					edgeNode.superNodeLink = link
				}
			}
		}
	}()

	return edgeNode, done
}

func (edgeNode *EdgeNode) Id() NodeId {
	return edgeNode.id
}

func (edgeNode *EdgeNode) Join(remoteAddress string) {
	edgeNode.transport.ConnectRemoteEndPoint(remoteAddress, false)
}

func (edgeNode *EdgeNode) Send(str string, dst NodeId) {
	msg := new(Msg)
	msg.Type = Data
	msg.Payload = str
	msg.Src = edgeNode.id.String()
	msg.Dst = dst.String()

	edgeNode.superNodeLink.send(msg)
}

func (edgeNode *EdgeNode) Announce() {
	msg := new(Msg)
	msg.Type = Announce
	msg.Src = edgeNode.id.String()

	edgeNode.superNodeLink.send(msg)
}

func (edgeNode *EdgeNode) Test() {
	fmt.Printf("EdgeNode.Test: \n")

	go func() {
		counter := 0
		for {
			msgStr := "hi"
			time.Sleep(2 * time.Second)
			if edgeNode.superNodeLink != nil {
				edgeNode.Send(msgStr, edgeNode.superNodeLink.remoteNodeId)
			}
			counter++
		}
	}()
}
