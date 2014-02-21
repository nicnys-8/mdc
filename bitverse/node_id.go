package bitverse

type NodeId struct {
	hashkey string
}

func makeNodeIdFromString(id string) NodeId {
	nodeId := NodeId{}
	nodeId.hashkey = HashkeyFromString(id)

	return nodeId
}

func generateNodeId() NodeId {
	nodeId := NodeId{}
	nodeId.hashkey = UniqueHashkey()

	return nodeId
}

func (nodeId *NodeId) Hashkey() string {
	return nodeId.hashkey
}

func (nodeId *NodeId) String() string {
	return nodeId.hashkey
}
