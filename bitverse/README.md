## About
The bitverse is a sub-project to the Cloudberry project currently running at [LTU](http://www.ltu.se) and is aiming at developing a P2P Messaging and Storage framework based on the [Chord DHT algorithm](http://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf). So far, only the messaging is implemented (and is only possible to have one super node). More to come in the future.

![](https://raw.github.com/ltu-cloudberry/mdc/master/bitverse/images/bitverse.png)

A bitverse network consists of two different types of nodes, so called **edge nodes** and **super nodes**. Edge nodes are typically connected to a super node, which makes sure messages are delivered even if an edge-node is located behind a firewall. Instead of addressing servers or devices using IP addresses, each node is globally addressable in the bitverse network using a hash key (e.g. *7d7dbf33abf34bdb7ef47231a7507372e2c908d6*). Each node (independent of other nodes) is responsible for calculating their own hash keys (a.k.a node id's). The node id currently a SHA-1 hash key of a self-generated UUID string.

Edge nodes typicially send messages to the super node. Messages send to other nodes connected to the same super node will be delivered directly by the super node. However, if an edge node is connected to a foreign super node somewhere, the super node will use the DHT to lookup the address (IP address) of the foreign super node and deliver the message to that super node. This means that the distance to any other node in the bitverse network is maximum 3 hops away independent of the size of the bitverse network.

Instead of using port numbers as in TCP or UDP to identify connections, the bitverse uses a service identifier, which can be an arbitrary string. To be able to send and receive messages, developer has to register and create a **service** object (which is very similar to a Unix socket) on their edge node objects. A service object only resides in the edge nodes where it is created. When a service object is created, developer has to provide an AES encryption key. Only edge nodes having access to that encryption key can send and receive messages to that particular service. That is, the super nodes are only responsible for forwarding messages towards their destination and are not able to decrypt the content enclosed in the messages.     

Currently, the only language supported is Go, but a JavaScript library will soon be released. 

## Installation
Just type `go get ltu-cloudberry/mdc/bitverse` to install the bitverse Go project. You can also find a bitverse binary file in the [bin](https://github.com/ltu-cloudberry/mdc/tree/master/bitverse/bin) directory if you just want to run a standalone super node. 

To setup a supernode, call `bitverse --local localhost:1111`, where the `--local` flag the specifies host and port where the super node should bind to. You may also pass the `--debug` flag if you want to enable debugging (more print traces).

## Example
To be able to create an edge node, a *BitverseObserver* compliant object must first be implemented. The edge node object will call functions in the bitverse observer object when it becomes connected to a super node, or when other nodes (siblings) joins or leaves the super node (it is possible to retreive a list of edge nodes on any other foreign super node in the bitverse network).   

```go
type MyBitverseObserver struct {
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingJoined(node *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " joined")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingLeft(node *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " left")
}

func (myBitverseObserver *MyBitverseObserver) OnSiblingHeartbeat(node *bitverse.EdgeNode, nodeId string) {
	fmt.Println("sibling " + nodeId + " heartbeat")
}

func (myBitverseObserver *MyBitverseObserver) OnChildrenReply(node *bitverse.EdgeNode, nodeId string, children []string) {
	fmt.Println("received children list from " + nodeId)
}

func (myBitverseObserver *MyBitverseObserver) OnConnected(node *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("now connected to super node " + remoteNode.Id())
}
```

To setup an edge node, we need to create a WebSocket transport and pass a reference to our bitverse observer.

```go
var done chan int

node, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
fmt.Println("-> my id is " + node.Id())

go node.Connect("localhost:1111")

<-done
```

### Messaging

To create a messaging service (a storage service will be supported in the future) we need to create a MsgServiceObserver. The OnDeliver function in the MsgServiceObserver object will be called when messages are received by our service object.

```go
func (msgServiceObserver *MsgServiceObserver) OnDeliver(msgService *bitverse.MsgService, msg *bitverse.Msg) {
...	
}
```

To create a service we need to call the *bitverse.CreateMsgService(...)* function and pass along a service id (does not have to be globally unique) and an AES encryption key. It is up to developer to securely pass and store the encryption key. It needs to be a 32 bit hex formated string, and can either be generated by calling `bitverse --generate-aes-secret` or by calling *bitverse.GenerateAesSecret()*.

```go
var serviceId = "my service"
var secret = "817130b45245941f4f8fd0ad77ccb5bf115faecf83d1956f9f01c666b9f35f6e"

msgServiceObserver := new(MsgServiceObserver)
edgeNode.CreateMsgService(secret, serviceId, msgServiceObserver)
```

Messages can easily be send using the *bitverse.Send(...)* function on the messaging service object.  

```go
var serviceId = "my service"
msgService := edgeNode.GetMsgService(serviceId)
msgService.Send("6a133a1b41f987210559ceb4ed9b1dbf58aec876", "hello")

```

If we except a reply back from the other node, we can call the *bitverse.SendAndGetReply(...)* and provide a closure which will be called when the reply is received. We also need to provide a timeout in seconds in case the other node failed to reply, or that node simply does not exists.

```go
msgService.SendAndGetReply("6a133a1b41f987210559ceb4ed9b1dbf58aec876", "hello", 10, func(success bool, reply *bitverse.Msg) {
		if success {
			fmt.Println("that was a surprise " + reply.Payload)
		} else {
			// will most likely timeout unless node 6a133a1b41f987210559ceb4ed9b1dbf58aec876 is online
			fmt.Println("failed to send message to node with id 6a133a1b41f987210559ceb4ed9b1dbf58aec876")
		}
	})
```

For a full example, see https://raw.github.com/ltu-cloudberry/mdc/master/bitverse/examples/messaging.go. Setup a super node at localhost:1111 (`bitverse --local localhost:1111`) and call *go run messaging.go*. 

## Documentation
See http://godoc.org/github.com/ltu-cloudberry/mdc/bitverse