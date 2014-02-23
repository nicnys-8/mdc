package bitverse

import (
	"encoding/json"
	"fmt"
)

type SuperNode struct {
	nodeId            NodeId
	children          map[string]*RemoteNode
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	seqNumberCounter  int
	localAddr         string
	localPort         string
	transport         Transport
	repoTable         map[string]*string // repoid:public key
}

func MakeSuperNode(transport Transport, localAddress string, localPort string) (*SuperNode, chan int) {
	superNode := new(SuperNode)

	superNode.localPort = localPort
	superNode.children = make(map[string]*RemoteNode)
	superNode.transport = transport

	superNode.repoTable = make(map[string]*string)

	superNode.nodeId = generateNodeId()
	debug("supernode: my id is " + superNode.Id())

	superNode.transport.SetLocalNodeId(superNode.nodeId)

	done := make(chan int)
	superNode.msgChannel = make(chan Msg)
	superNode.remoteNodeChannel = make(chan *RemoteNode, 10)

	go superNode.transport.Listen(localAddress, localPort, superNode.remoteNodeChannel, superNode.msgChannel)

	go func() {
		for {
			select {
			case msg := <-superNode.msgChannel:
				debug("supernode: received " + msg.String())
				if msg.Dst == superNode.Id() && msg.Type == Data {
					// ignore, not supported

				} else if msg.Type == Data && msg.ServiceType == Repo && msg.RepoCmd == Claim { // repo claim request
					repoId := msg.RepoId
					pubKeyPem := msg.Signature
					debug("supernode: got a repo claim request for " + repoId + " with key " + pubKeyPem)

					if superNode.repoTable[repoId] == nil {
						// it is free, claim it!
						superNode.repoTable[repoId] = &pubKeyPem // XXX is this safe?
						msg.Status = Ok

					} else {
						// already claimed
						if pubKeyPem == *superNode.repoTable[repoId] {
							// but, it the same owner
							msg.Status = Ok
						} else {
							msg.Status = Error
							msg.Payload = "repo id already claimed"
						}
					}

					childId := msg.Src
					msg.Dst = childId
					msg.Src = superNode.Id()
					superNode.sendToChild(msg)

				} else if msg.Type == Data && msg.ServiceType == Repo && msg.RepoCmd == Store { // repo store request
					repoId := msg.RepoId

					if superNode.repoTable[repoId] == nil {
						msg.Status = Error
						msg.Payload = "no such repo " + repoId

					} else {
						//key := msg.RepoKey
						//value := msg.RepoValue
						//signature := msg.Signature

						msg.Status = Ok
						msg.PayloadType = Nil
						//msg.Payload = "OLD VALUE"
					}

					childId := msg.Src
					msg.Dst = childId
					msg.Src = superNode.Id()
					superNode.sendToChild(msg)

				} else if msg.Type == Heartbeat {
					superNode.forwardToChildren(msg)

				} else if msg.Type == Children {
					superNode.sendChildrenReply(msg.Src)

				} else {
					superNode.sendToChild(msg)
				}
			case remoteNode := <-superNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					delete(superNode.children, remoteNode.Id())

					str := fmt.Sprintf("supernode: removing remote node %s, number of remote nodes are now %d", remoteNode.Id(), len(superNode.children))
					fmt.Println(str)

					msg := composeChildLeft(superNode.nodeId.String(), remoteNode.Id())
					superNode.forwardToChildren(*msg)
				} else {
					superNode.children[remoteNode.Id()] = remoteNode

					str := fmt.Sprintf("supernode: adding remote node %s, number of remote nodes are now %d", remoteNode.Id(), len(superNode.children))
					info(str)

					msg := composeChildJoin(superNode.nodeId.String(), remoteNode.Id())
					superNode.forwardToChildren(*msg)
				}
			}
		}
	}()

	return superNode, done
}

// BITVERSE MANAGEMENT

func (superNode *SuperNode) Id() string {
	return superNode.nodeId.String()
}

// DEBUG

func (superNode *SuperNode) Debug() {
	debugFlag = true
}

/// PRIVATE

func (superNode *SuperNode) sendChildrenReply(nodeId string) {
	debug("supernode: sending children reply to " + nodeId)
	childrenIds := make([]string, len(superNode.children))
	i := 0
	for childNodeId, _ := range superNode.children {
		childrenIds[i] = childNodeId
		i++
	}

	json, _ := json.Marshal(childrenIds)
	reply := composeChildrenReplyMsg(superNode.Id(), nodeId, string(json))

	remoteNode := superNode.children[nodeId]

	if remoteNode != nil {
		remoteNode.deliver(reply)
	}
}

func (superNode *SuperNode) sendToChild(msg Msg) {
	for _, remoteNode := range superNode.children {
		if msg.Src != remoteNode.Id() && msg.Dst == remoteNode.Id() { // do not forward messages to a remote node where it came from
			debug("supernode: forwarding " + msg.String() + " to " + remoteNode.Id())
			msg.Dst = remoteNode.Id()
			remoteNode.deliver(&msg)
		} else {
			debug("failed to forward message to child")
		}
	}
}

func (superNode *SuperNode) forwardToChildren(msg Msg) {
	for _, remoteNode := range superNode.children {
		if msg.Src != remoteNode.Id() { // do not forward messages to a remote node where it came from
			debug("supernode: forwarding " + msg.String() + " to " + remoteNode.Id())

			remoteNode.deliver(&msg)
		}
	}
}
