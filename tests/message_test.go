package main

import (
	"fmt"
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	"net"
	"testing"
	// cp "../src/clientproperties"
)

//TestRequestMessage message request testing
func TestRequestMessage(t *testing.T) {

	ln, err := net.Listen("tcp", ":40000")
	if err != nil {
		fmt.Println("Error ", err, " in listening on address ", ln)
	}

	var list []string
	TestMapPeerIP := make(map[string]string)
	TestMapPeerListenPort := make(map[string]string)

	list = append(list, "abc")
	TestMapPeerIP["abc"] = "127.0.0.1"
	TestMapPeerListenPort["abc"] = "40000"
	activeClient := cp.ClientListen{List: list, PeerIP: TestMapPeerIP, PeerListenPort: TestMapPeerListenPort}

	name := "my_name"

	messageStatus := cp.RequestMessage(&activeClient, name, "abc", "hey buddy!")
	if messageStatus != "sent" {
		t.Fatal("Error in sending message ...")
	}

}
