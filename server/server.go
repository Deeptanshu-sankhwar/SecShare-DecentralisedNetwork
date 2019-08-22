package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	// "strings"
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	en "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/encryptionproperties"
	sp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/serverproperties"
	"net"
	"sync"
	// cp "../src/clientproperties"
	// en "../src/encryptionproperties"
	// sp "../src/serverproperties"
)

/* Add relevant Print statements where confused , and comment print statements while pushing */

var mutex = &sync.Mutex{} // Lock and unlock (Mutex)

var clients []cp.Client

var cli cp.ClientListen

var jobs []cp.ClientJob

func pingAll(clients []cp.Client, cli cp.ClientListen) {
	for i := 0; i < len(clients); i++ {
		encode := json.NewEncoder(clients[i].ConnectionServer) //sending to each online peer
		encode.Encode(cli)
	}
}

func performJobs() { // storing each job in a queue in the server and executing it one by one

	for {
		if len(jobs) != 0 {
			mutex.Lock()
			// fmt.Println("number of jobs currently are ", len(jobs))
			getJob := jobs[0]
			jobs = jobs[1:]
			mutex.Unlock()
			// fmt.Println("number of jobs currently are ", len(jobs))
			handler(getJob.Conn, getJob.Name, getJob.Query, getJob.ClientListenPort)
		}
	}
}

func handler(c net.Conn, name string, query string, ClientListenPort string) { // handling each connection

	if query == "login" {

		remoteAddress := c.RemoteAddr().String()
		newClient := cp.Client{Address: remoteAddress, Name: name, ConnectionServer: c} //making struct
		cli.PeerIP[name] = remoteAddress
		cli.PeerListenPort[name] = ClientListenPort //creating the map
		clients = append(clients, newClient)        //append
		cli.List = append(cli.List, name)
		go pingAll(clients, cli)

	} else if query == "quit" {

		cli = sp.QueryDeal(&clients, cli, name)
	} else if query == "" {

		var name string
		remoteAddress := c.RemoteAddr().String()
		for i := 0; i < len(clients); i++ {
			if clients[i].Address == remoteAddress {
				name = clients[i].Name
			}
		}

		cli = sp.QueryDeal(&clients, cli, name)
	}
	fmt.Print("no of active clients are : ", len(cli.List))
	fmt.Print("Active clients are -> ", cli.List, "\n")
	fmt.Print("Active clients IPs are -> ", cli.PeerIP, "\n")

}

func maintainConnection(conn net.Conn, PeerKeys map[net.Conn]*rsa.PublicKey,
	pub *rsa.PublicKey, pri *rsa.PrivateKey) { //maintaining the connection between client and server

	//performing handshake
	peerKey := &rsa.PublicKey{}

	decoder := json.NewDecoder(conn)
	decoder.Decode(&peerKey)
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)
	PeerKeys[conn] = peerKey
	// fmt.Print(peerKey.N)
	for {
		clientQuery := cp.ClientQuery{}

		decoder := json.NewDecoder(conn)
		decoder.Decode(&clientQuery)

		Name := string(en.DecryptWithPrivateKey(clientQuery.Name, pri))
		Query := string(en.DecryptWithPrivateKey(clientQuery.Query, pri))
		ClientListenPort := string(en.DecryptWithPrivateKey(clientQuery.ClientListenPort, pri))
		// fmt.Println("name and query are ", Name, Query)
		job := cp.ClientJob{Name: Name, Query: Query, Conn: conn, ClientListenPort: ClientListenPort}
		// fmt.Println("current job is ", job.Query)

		mutex.Lock()
		if job.Query != "" {
			jobs = append(jobs, job)
			fmt.Println("appended job is ", job)
		}
		mutex.Unlock()

		if Query == "quit" || Query == "" {
			break
		}
	}
	conn = nil
}

func main() {

	ln, _ := net.Listen("tcp", ":8081") // making a server
	fmt.Println(": SERVER STARTED ON PORT 8081  : ")
	// fmt.Print(ln.LocalAddr().String())
	PrivateKey, PublicKey := en.GenerateKeyPair()

	PeerKeys := make(map[net.Conn]*rsa.PublicKey)

	cli = cp.ClientListen{List: []string{}, PeerIP: make(map[string]string),
		PeerListenPort: make(map[string]string)}
	go performJobs()

	for {
		conn, _ := ln.Accept()
		go maintainConnection(conn, PeerKeys, PublicKey, PrivateKey) //accept a new connection and maintain it using the function above
	}

}
