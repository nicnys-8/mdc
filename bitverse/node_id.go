package bitverse

import (
	"github.com/nu7hatch/gouuid"
	"log"
)

type NodeId struct {
	unqiueId string
}

func makeNodeIdFromString(str string) NodeId {
	nodeId := NodeId{}
	nodeId.unqiueId = str

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
