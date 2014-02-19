package bitverse

import (
	"crypto/sha1"
	//"encoding/base64"
	"fmt"
	"github.com/nu7hatch/gouuid"
)

type NodeId struct {
	hashkey string
}

func makeNodeIdFromString(id string) NodeId {
	nodeId := NodeId{}
	nodeId.hashkey = id

	return nodeId
}

func generateNodeId() NodeId {
	nodeId := NodeId{}

	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(u.String()))
	//nodeId.hashkey = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	nodeId.hashkey = fmt.Sprintf("%x", hasher.Sum(nil))

	return nodeId
}

func (nodeId *NodeId) Hashkey() string {
	return nodeId.hashkey
}

func (nodeId *NodeId) String() string {
	return nodeId.hashkey
}
