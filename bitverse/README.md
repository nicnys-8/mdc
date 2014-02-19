## About
The bitverse project aiming at developing a P2P Messaging and Storage framework based on the [Chord DHT algorithm](http://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf). So far, only the messaging is implemented. 

A bitverse network consists of two different nodes, so called edge-nodes and super-nodes.  

![](https://raw.github.com/ltu-cloudberry/mdc/master/bitverse/images/bitverse.png)

## Installation
Just type `go get ltu-cloudberry/mdc/bitverse` to install the bitverse Go project. You can also find a bitverse binary file in the bin directory if you just want to run a super node. 
To setup a supernode, call `bitverse --super --local localhost:1111`, where the `--local` flag the specifies host and port the super node should bind to. If you want to enable debugging, you may also pass the `--debug` flag.

## Example
```go
type MyBitverseObserver struct {
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingJoined(edgeNode *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " joined")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingLeft(edgeNode *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " left")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingHeartbeat(edgeNode *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " heartbeat")
}

func (myBitverseObserver *MyBitverseObserver) OnChildrenReply(edgeNode *bitverse.EdgeNode, nodeId string, children []string) {
	fmt.Println("received children list from " + nodeId)
}

func (myBitverseObserver *MyBitverseObserver) OnConnected(edgeNode *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("now connected to super node " + remoteNode.Id())
}
```

```go
var done chan int

edgeNode, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
fmt.Println("-> my id is " + edgeNode.Id())

go edgeNode.Connect("localhost:1111")

<-done
```

## Documentation
See http://godoc.org/github.com/ltu-cloudberry/mdc/bitverse
