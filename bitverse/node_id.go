package main

import (
	"github.com/nu7hatch/gouuid"
	"log"
	"strconv"
)

const (
	UNICAST = iota
	BROADCAST
)

type NodeId struct {
	addressType int
	unqiueId    string
}

func makeNodeIdFromString(str string) NodeId {
	nodeId := NodeId{}
	nodeId.unqiueId = str
	nodeId.addressType = UNICAST

	return nodeId
}

func makeBroadcastAddress() NodeId {
	nodeId := NodeId{}
	nodeId.unqiueId = strconv.Itoa(UNICAST)
	nodeId.addressType = BROADCAST

	return nodeId
}

func generateNodeId() NodeId {
	nodeId := NodeId{}

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	nodeId.unqiueId = u.String()

	return nodeId
}

func (nodeId *NodeId) String() string {
	return nodeId.unqiueId
}
