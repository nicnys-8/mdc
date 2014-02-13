package bitverse

import (
	"fmt"
	"time"
)

var HEARTBEAT_RATE time.Duration = 2

type EdgeNode struct {
	id                NodeId
	wsServer          *WsServer
	superNodeLink     *RemoteNode
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	seqNumberCounter  int
	transport         Transport
	services          map[string]*Service
}

func MakeEdgeNode(transport Transport) (*EdgeNode, chan int) {
	edgeNode := new(EdgeNode)
	edgeNode.transport = transport

	edgeNode.id = generateNodeId()
	edgeNode.transport.SetLocalNodeId(edgeNode.id)

	done := make(chan int)
	edgeNode.msgChannel = make(chan Msg)
	edgeNode.remoteNodeChannel = make(chan *RemoteNode, 10)
	edgeNode.services = make(map[string]*Service)

	go func() {
		for {
			select {
			case msg := <-edgeNode.msgChannel:
				fmt.Println("edgenode: received " + msg.String())
				if msg.Dst == edgeNode.id.String() && msg.Type == Data {
					service := edgeNode.services[msg.Service]
					if service == nil {
						fmt.Println("edgenode: failed to deliver message, no such service " + msg.Service)
					} else {
						observer := service.observer
						if observer == nil {
							fmt.Println("edgenode: failed to deliver message, no observer registered")
						} else {
							observer.OnDeliver(&msg)
						}
					}
				} else if msg.Type == Heartbeat {
					//fmt.Println("SuperNode: got HEARBEAT message from <" + msg.Src + ">")
				} else { // ignore
				}
			case remoteNode := <-edgeNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					fmt.Println("edgenode: ERROR we just lost our connection to the super node <" + remoteNode.id.String() + ">")
					edgeNode.superNodeLink = nil
				} else {
					fmt.Println("edgenode: adding link to super node <" + remoteNode.id.String() + ">")
					edgeNode.superNodeLink = remoteNode
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Millisecond * HEARTBEAT_RATE * 1000)
	go func() {
		for t := range ticker.C {
			fmt.Println("edgenode: sending heartbeat", t)
			edgeNode.SendHeartbeat()
		}
	}()

	return edgeNode, done
}

func (edgeNode *EdgeNode) Id() NodeId {
	return edgeNode.id
}

func (edgeNode *EdgeNode) Join(remoteAddress string) {
	edgeNode.transport.ConnectToNode(remoteAddress, edgeNode.remoteNodeChannel, edgeNode.msgChannel)
}

func (edgeNode *EdgeNode) GetService(name string, observer ServiceObserver) *Service {
	if edgeNode.services[name] == nil {
		service := composeService(name, observer, edgeNode)
		edgeNode.services[name] = service
		return service
	} else {
		return edgeNode.services[name]
	}
}

func (edgeNode *EdgeNode) SendHeartbeat() {
	msg := ComposeHeartbeatMsg(edgeNode.id.String(), edgeNode.superNodeLink.id.String())
	edgeNode.superNodeLink.send(msg)
}

func (edgeNode *EdgeNode) Checkout(id string, rev int) (dict *Dictionary) {
	return nil
}

func (edgeNode *EdgeNode) Checkin(dictionary *Dictionary) (rev int) {
	return 0
}

/// PRIVATE

func (edgeNode *EdgeNode) send(dst string, payload string, service string) {
	msg := ComposeDataMsg(edgeNode.id.String(), dst, payload, service)
	edgeNode.superNodeLink.send(msg)
}
