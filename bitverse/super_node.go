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
					superNode.forwardToChildren(msg)
				} else if msg.Type == Children {
					superNode.sendChildrenReply(msg.Src)
				} else { // to someone else, forward
					superNode.forwardToChild(msg)
				}
			case remoteNode := <-superNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					delete(superNode.children, remoteNode.Id.String())
					fmt.Printf("supernode: removing remote node %s, number of remote nodes are now %d\n", remoteNode.Id.String(), len(superNode.children))
				} else {
					superNode.children[remoteNode.Id.String()] = remoteNode
					fmt.Printf("supernode: adding remote node %s, number of remote nodes are now %d\n", remoteNode.Id.String(), len(superNode.children))
				}
			}
		}
	}()

	return superNode, done
}

func (superNode *SuperNode) Id() NodeId {
	return superNode.id
}

/// PRIVATE

func (superNode *SuperNode) sendChildrenReply(nodeId string) {
	fmt.Println("supernode: sending children reply to " + nodeId)
	//reply := ""

	for childNodeId, _ := range superNode.children {
		fmt.Println("child node id = " + childNodeId)
	}
}

func (superNode *SuperNode) forwardToChild(msg Msg) {
	// send the message on all our links
	for _, remoteNode := range superNode.children {
		if msg.Src != remoteNode.Id.String() && msg.Dst == remoteNode.Id.String() { // do not forward messages to a remote node where it came from
			fmt.Println("supernode: forwarding " + msg.String() + " to " + remoteNode.Id.String())
			remoteNode.send(&msg)
		}
	}
}

func (superNode *SuperNode) forwardToChildren(msg Msg) {
	// send the message on all our links
	for _, remoteNode := range superNode.children {
		if msg.Src != remoteNode.Id.String() { // do not forward messages to a remote node where it came from
			fmt.Println("supernode: forwarding " + msg.String() + " to " + remoteNode.Id.String())
			remoteNode.send(&msg)
		}
	}
}
