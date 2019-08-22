package main

import (
	// cp "../src/clientproperties"
	// sp "../src/serverproperties"
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	sp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/serverproperties"
	"testing"
)

//TestRemoveFromClient test if client is deleted or not
func TestRemoveFromClient(t *testing.T) {

	clientList := []cp.Client{}
	cli := cp.Client{Address: "001", Name: "dummy"}
	clientList = append(clientList, cli)
	clientList = sp.RemoveFromClient(clientList, "dummy")
	if len(clientList) > 0 {
		t.Fatal("RemoveFromClient not working correctly")
	}
}

// TestQueryDeal check client query is working or not
func TestQueryDeal(t *testing.T) {
	name := "user"
	TestMap := make(map[string]string)
	TestMap[name] = "active"
	var list []string
	list = append(list, "user")
	cli := cp.ClientListen{List: list, PeerIP: TestMap}
	clt := cp.Client{Name: name, Address: ":8087"}
	var clients = []cp.Client{clt}
	Client := sp.QueryDeal(&clients, cli, name)

	if len(clients) > 0 || len(Client.List) > 0 || cli.PeerIP[name] != "" {
		t.Fatal("QueryDeal in the serverproperties not working correctly")
	}

}
