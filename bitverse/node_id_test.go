package bitverse

import (
	"testing"
)

func TestNodeId(t *testing.T) {
	nodeId1 := generateNodeId()
	fmt.Printf("nodeId1=" + nodeId1)
	//t.Fatalf("unexpected err. %s", err)
}
