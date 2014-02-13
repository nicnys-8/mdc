package main

import (
	"flag"
	"fmt"
	"log"
	"mdc/bitverse"
	"net/http"
	"strings"
	"time"
)

var superFlag = flag.Bool("super", false, "run the node as a super node")
var localFlag = flag.String("local", "", "ip address and port which this super node should bound to, e.g. --local localhost:1111")
var joinFlag = flag.String("join", "", "ip address and port to a node to join, e.g. --join localhost:2222")
var testHttpServerFlag = flag.Bool("test-http-server", false, "starts a http test server at port 8080 for debuging")

func main() {
	flag.Parse()

	var done chan int

	transport := bitverse.MakeWSTransport()

	if *superFlag {
		temp := strings.Split(*localFlag, ":")
		localAddr := temp[0]
		localPort := temp[1]

		_, done = bitverse.MakeSuperNode(transport, localAddr, localPort)

		if *testHttpServerFlag {
			fmt.Println("Starting a HTTP test server at port 8080")
			log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
		}
	} else {
		var edgeNode *bitverse.EdgeNode
		edgeNode, done = bitverse.MakeEdgeNode(transport)
		edgeNode.Test()

		ticker := time.NewTicker(time.Millisecond * 2000)
		go func() {
			for t := range ticker.C {
				fmt.Println("Sending announcement", t)
				edgeNode.Announce()
			}
		}()

		// join super node
		remoteAddress := *joinFlag
		if remoteAddress != "" {
			fmt.Println("EdgeNode: joining node at " + remoteAddress)
			go edgeNode.Join(remoteAddress)
		}
	}

	<-done
}