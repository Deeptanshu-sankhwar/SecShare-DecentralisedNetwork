package main

import (
	"testing"
	// cp "../src/clientproperties"
	// fp "../fileproperties"
	"fmt"
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	"net"
)

//TestSendFileParts testing send file parts
func TestSendFileParts(t *testing.T) {
	filename := "image.jpg"
	allFileParts := cp.GetSplitFile(filename, 2)
	myname := "def"
	ln, err := net.Listen("tcp", ":45000")

	if err != nil {
		fmt.Println("Error ", err, " in listening on address ", ln)
	}

	var list []string
	TestMapPeerIP := make(map[string]string)
	TestMapPeerListenPort := make(map[string]string)
	list = append(list, "abc")
	TestMapPeerIP["abc"] = "127.0.0.1"
	TestMapPeerListenPort["abc"] = "45000"
	activeClient := cp.ClientListen{List: list, PeerIP: TestMapPeerIP, PeerListenPort: TestMapPeerListenPort}
	newfilerequest := cp.FileRequest{Query: "receive_file", MyAddress: "127.0.0.1", MyName: "abc", RequestedFile: "dummy.txt"}
	countSent := cp.SendFileParts(newfilerequest, allFileParts, &activeClient, myname)
	fmt.Println(countSent, len(activeClient.PeerListenPort))

	if countSent != len(activeClient.PeerListenPort) {
		t.Fatal("File parts improperly sent!")
	}

}
