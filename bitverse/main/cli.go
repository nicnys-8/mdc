package main

import (
	"flag"
	"fmt"
	"log"
	"mdc/bitverse"
	"net/http"
	"strings"
)

var debugFlag = flag.Bool("debug", false, "run the node in debug mode")
var aesSecretFlag = flag.Bool("generate-aes-secret", false, "generate hex encoded aes secret")
var certFlag = flag.String("generate-rsa-keys", "", "generate rsa public and private keys, e.g. --generate-cert mycert")
var localFlag = flag.String("local", "", "ip address and port which this super node should bound to, e.g. --local localhost:1111")
var joinFlag = flag.String("join", "", "ip address and port to a node to join, e.g. --join localhost:2222")
var testHttpServerFlag = flag.Bool("test-http-server", false, "starts a http test server at port 8080 for debuging")

/// MAIN

func main() {
	flag.Parse()
	if *aesSecretFlag {
		aesSecret, err := bitverse.GenerateAesSecret()

		if err != nil {
			panic(err)
		}

		fmt.Println(aesSecret)
	} else if *certFlag != "" {
		fmt.Println("generating files " + *certFlag + " " + *certFlag + ".pub")
		bitverse.GeneratePem(*certFlag)
	} else {
		// set up super node
		var done chan int

		transport := bitverse.MakeWSTransport()

		var superNode *bitverse.SuperNode
		temp := strings.Split(*localFlag, ":")
		localAddr := temp[0]
		localPort := temp[1]

		superNode, done = bitverse.MakeSuperNode(transport, localAddr, localPort)

		if *debugFlag {
			superNode.Debug()
		}

		if *testHttpServerFlag {
			fmt.Println("Starting a HTTP test server at port 8080")
			log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./js/"))))
		}

		<-done
	}
}
