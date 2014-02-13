package bitverse

import (
	"fmt"
)

type SuperNode struct {
	id                NodeId
	wsServer          *WsServer
	children          map[string]*RemoteNode
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	seqNumberCounter  int
	localAddr         string
	localPort         string
	transport         Transport
}

func MakeSuperNode(transport Transport, localAddress string, localPort string) (*SuperNode, chan int) {
	superNode := new(SuperNode)

	superNode.localPort = localPort
	superNode.children = make(map[string]*RemoteNode)
	superNode.transport = transport

	superNode.id = generateNodeId()
	superNode.transport.SetLocalNodeId(superNode.id)

	done := make(chan int)
	superNode.msgChannel = make(chan Msg)
	superNode.remoteNodeChannel = make(chan *RemoteNode, 10)

	go superNode.transport.Listen(localAddress, localPort, superNode.remoteNodeChannel, superNode.msgChannel)

	go func() {
		for {
			select {
			case msg := <-superNode.msgChannel:
				fmt.Println("supernode: received " + msg.String())
				if msg.Dst == superNode.id.String() && msg.Type == Data { // ignore
				} else if msg.Type == Heartbeat {
					superNode.Forward(msg)
				} else { // to someone else, forward
					superNode.Forward(msg)
				}
			case remoteNode := <-superNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					delete(superNode.children, remoteNode.id.String())
					fmt.Printf("supernode: removing remote node %s, number of remote nodes are now %d\n", remoteNode.id.String(), len(superNode.children))
				} else {
					superNode.children[remoteNode.id.String()] = remoteNode
					fmt.Printf("supernode: adding remote node %s, number of remote nodes are now %d\n", remoteNode.id.String(), len(superNode.children))
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
	for _, remoteNode := range superNode.children {
		if msg.Src != remoteNode.id.String() { // do not forward messages to a remote node where it came from
			fmt.Println("supernode: forwarding " + msg.String() + " to " + remoteNode.id.String())
			remoteNode.send(&msg)
		}
	}
}
