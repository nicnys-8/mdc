package main

import (
	"fmt"
	"mdc/bitverse"
)

var repoId = "test"
var secret = "5da71277f031a9dff561f0a72bb72651e260dab0735b767f2f7a62dec9e99760"

// SERVICE OBSERVER

/// BITVERSE OBSERVER

type BitverseObserver struct {
}

func (bitverseObserver *BitverseObserver) OnSiblingJoined(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " joined")
}

func (bitverseObserver *BitverseObserver) OnSiblingLeft(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " left")
}

func (bitverseObserver *BitverseObserver) OnSiblingHeartbeat(node *bitverse.EdgeNode, id string) {
	fmt.Println("-> sibling " + id + " heartbeat")
}

func (bitverseObserver *BitverseObserver) OnChildrenReply(node *bitverse.EdgeNode, id string, children []string) {
	fmt.Println("-> received children list from " + id)
}

func (bitverseObserver *BitverseObserver) OnConnected(node *bitverse.EdgeNode, remoteNode *bitverse.RemoteNode) {
	fmt.Println("-> now connected to super node " + remoteNode.Id())

	prv, pub, err := bitverse.ImportPem("cert")
	if err != nil {
		fmt.Println(err)
	}

	node.ClaimRepository(repoId, secret, prv, pub, 10, func(err error, repo interface{}) {
		if err != nil {
			fmt.Println("failed to claim repo")
		} else {
			fmt.Println("sucessfully claimed repo <test>")
			testRepo := repo.(*bitverse.RepoService)

			testRepo.Store("myKey", "myValue", 10, func(err error, oldValue interface{}) {
				if err != nil {
					fmt.Println("failed to store key in bitverse network")
				} else {
					switch oldValue.(type) {
					case string:
						fmt.Printf("replacing key-value pair in the bitverse network, old value was " + oldValue.(string))
					case nil:
						fmt.Printf("storing new key-value pair in the bitverse network")
					}
				}
			})

		}
	})
}

/// MAIN

func main() {
	var done chan int

	node, done := bitverse.MakeEdgeNode(bitverse.MakeWSTransport(), new(BitverseObserver))
	//node.Debug()
	fmt.Println("-> my id is " + node.Id())

	go node.Connect("localhost:1111")

	<-done
}
