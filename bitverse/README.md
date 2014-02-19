## About
The bitverse project is aiming at developing a P2P Messaging and Storage framework based on the [Chord DHT algorithm](http://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf). So far, only the messaging part is implemented (and is only possible to have one super node). More to come in the future.

![](https://raw.github.com/ltu-cloudberry/mdc/master/bitverse/images/bitverse.png)

A bitverse network consists of two different types of nodes, so called **edge nodes** and **super nodes**. Edge nodes are typically connected to a super node, which makes sure messages are delivered even if an edge-node is located behind a firewall. Instead of addressing servers or devices using IP addresses, each node is globally addressable in the bitverse network using a hash key (e.g. aCfy3BQP09RdVIw78SRb1gL_df8=). Each node (independent of other nodes) is responsible for calculating the hash key (node id). The node id currently based on Base64 encoded SHA-1 hash key of a self-generated UUID string.

Edge nodes always send messages to the super node. Messages send to other nodes directly connected to the super node will be delivered directly by the super node. However, if an edge node is connected to a foreign super node somewhere, the super node will use the DHT to lookup the address (IP address) and deliver the message to that super node. This means that the distance to any other node in the bitverse network is maximum 3 hops away independent of the size of the bitverse network.

Instead of using port numbers as in TCP or UDP to identify connections, the bitverse uses a globally unique service identifier, which is also a UUID string. To be able to send and receive messages, developer has to register and create a **service** object, which very similar to a Unix socket. A service object only resides in the edge nodes that create such a service. When a service object is created, developer has to provide an AES encryption key. Only edge nodes having access to that encryption key can send and receive messages to that particular service. The super nodes will only forward messages towards their destination and are not able to decrypt the content enclosed in the messages.     

Currently, the only language supported is Go, but a JavaScript library will soon be released. 

## Installation
Just type `go get ltu-cloudberry/mdc/bitverse` to install the bitverse Go project. You can also find a bitverse binary file in the bin directory if you just want to run a standalone super node. 
To setup a supernode, call `bitverse --super --local localhost:1111`, where the `--local` flag the specifies host and port the super node should bind to. You may also pass the `--debug` flag if you want to enable debugging,.

## Example
To be able to create an edge node, a *BitverseObserver* compliant object must first be implemented. The edge node object will call functions in the bitverse observer object when it becomes connected to a super node, or when other nodes (siblings) joins or leaves the super node (it will also be possible to retreive a list of edge nodes of any super node in the bitverse network).   

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

To setup an edge node, we need to create a WebSocket transport and pass a our reference to our bitverse observer.

```go
var done chan int

edgeNode, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
fmt.Println("-> my id is " + edgeNode.Id())

go edgeNode.Connect("localhost:1111")

<-done
```

To create a messaging service (a storage service will be supported in the future) we need to create a MsgServiceObserver. The OnDeliver function in the MsgServiceObserver object will be called when messages are received by our service object.

```go
func (msgServiceObserver *MsgServiceObserver) OnDeliver(msgService *bitverse.MsgService, msg *bitverse.Msg) {
...	
}
```

To create a service we need to call the *CreateMsgService* function and pass along a service id and an encryption key. It is up to developer to securely pass and store the encryption key.

```go
var serviceId = "6107911a-7554-4ea7-80fc-25ec5e2462a7" // uuid
var secret = "x very very very very secret key"        // aes encryption key, 16, 24, or 32 bytes

msgServiceObserver := new(MsgServiceObserver)
edgeNode.CreateMsgService(secret, serviceId, msgServiceObserver)
```

Messages can easily be send using the *Send* function on the messaging service object.  

```go
var serviceId = "6107911a-7554-4ea7-80fc-25ec5e2462a7" // uuid
msgService := edgeNode.GetMsgService(serviceId)
msgService.Send(id, "hello")

```

If we except a reply back from the other node, we can call the *SendAndGetReply* and provide a closure which will be called when the reply is received. We also need to provide a timeout in seconds in case the other node failed to reply or the node simply does not exists.

```go
msgService.SendAndGetReply("does not exists", "", 10, func(timedOut bool, reply *bitverse.Msg) {
	if timedOut {
		fmt.Println("failed to send message")
	} else {
		fmt.Println("got a message" + reply.Payload)
	}
})

For a full example, see 
```

## Documentation
See http://godoc.org/github.com/ltu-cloudberry/mdc/bitverse
