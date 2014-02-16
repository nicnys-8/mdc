package bitverse

import (
	"encoding/json"
	"fmt"
	"time"
)

var HEARTBEAT_RATE time.Duration = 10

type EdgeNode struct {
	nodeId            NodeId
	wsServer          *WsServer
	superNode         *RemoteNode
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	seqNumberCounter  int
	transport         Transport
	msgServices       map[string]*MsgService
	bitverseObserver  BitverseObserver
}

func MakeEdgeNode(transport Transport, bitverseObserver BitverseObserver) (*EdgeNode, chan int) {
	edgeNode := new(EdgeNode)
	edgeNode.transport = transport
	edgeNode.bitverseObserver = bitverseObserver

	edgeNode.nodeId = generateNodeId()
	fmt.Println("edgenode: my id is " + edgeNode.Id())

	edgeNode.transport.SetLocalNodeId(edgeNode.nodeId)

	done := make(chan int)
	edgeNode.msgChannel = make(chan Msg)
	edgeNode.remoteNodeChannel = make(chan *RemoteNode, 10)
	edgeNode.msgServices = make(map[string]*MsgService)

	go func() {
		for {
			select {
			case msg := <-edgeNode.msgChannel:
				fmt.Println("edgenode: received " + msg.String())
				if msg.Dst == edgeNode.Id() && msg.Type == Data {
					msgService := edgeNode.msgServices[msg.ServiceId]
					if msgService == nil {
						fmt.Println("edgenode: failed to deliver message, no such service with id <" + msg.ServiceId + "> created")
					} else {
						observer := msgService.observer
						if observer == nil {
							fmt.Println("edgenode: failed to deliver message, no observer registered")
						} else {
							var err error
							msg.Payload, err = decrypt(msgService.aesKey, msg.Payload)
							if err != nil {
								fmt.Println("edgenode: failed to decrypt payload, ignoring msg")
							} else {
								observer.OnDeliver(msgService, &msg)

							}
						}
					}
				} else if msg.Type == Heartbeat {
					//fmt.Println("edgenode: got HEARBEAT message from <" + msg.Src + ">")
				} else if msg.Type == Children {
					if bitverseObserver != nil {
						var children []string
						if err := json.Unmarshal([]byte(msg.Payload), &children); err != nil {
							panic(err)
						}
						bitverseObserver.OnChildrenReply(edgeNode, msg.Src, children)
					}
				} else { // ignore
				}
			case remoteNode := <-edgeNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					fmt.Println("edgenode: ERROR we just lost our connection to the super node <" + remoteNode.Id() + ">")
					edgeNode.superNode = nil
				} else {
					fmt.Println("edgenode: adding link to super node <" + remoteNode.Id() + ">")
					edgeNode.superNode = remoteNode
					if bitverseObserver != nil {
						bitverseObserver.OnConnected(edgeNode, edgeNode.superNode)
					}
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

func (edgeNode *EdgeNode) Id() string {
	return edgeNode.nodeId.String()
}

func (edgeNode *EdgeNode) Connect(remoteAddress string) {
	edgeNode.transport.ConnectToNode(remoteAddress, edgeNode.remoteNodeChannel, edgeNode.msgChannel)
}

func (edgeNode *EdgeNode) CreateMsgService(secret string, name string, observer MsgServiceObserver) *MsgService {
	if edgeNode.msgServices[name] == nil {
		msgService := composeMsgService(secret, name, observer, edgeNode)
		edgeNode.msgServices[name] = msgService
		return msgService
	} else {
		return edgeNode.msgServices[name]
	}
}

func (edgeNode *EdgeNode) GetMsgService(name string) *MsgService {
	return edgeNode.msgServices[name]
}

func (edgeNode *EdgeNode) SendHeartbeat() {
	msg := ComposeHeartbeatMsg(edgeNode.Id(), edgeNode.superNode.Id())
	edgeNode.superNode.send(msg)
}

func (edgeNode *EdgeNode) Checkout(id string, rev int) (dict *Dictionary) {
	return nil
}

func (edgeNode *EdgeNode) Checkin(dictionary *Dictionary) (rev int) {
	return 0
}

/// PRIVATE

func (edgeNode *EdgeNode) send(dst string, payload string, service string) {
	msg := ComposeDataMsg(edgeNode.Id(), dst, service, payload)
	edgeNode.superNode.send(msg)
}
