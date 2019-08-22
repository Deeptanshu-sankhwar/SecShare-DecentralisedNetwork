package serverproperties

import (
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	// cp "../clientproperties"
)

// RemoveFromClient - removes the client details who quits
func RemoveFromClient(clients []cp.Client, name string) []cp.Client {
	tempClients := []cp.Client{}
	for i := 0; i < len(clients); i++ {
		if clients[i].Name != name {
			tempClients = append(tempClients, clients[i])
		}
	}
	return tempClients
}

// QueryDeal upadtes client list and PeerIP table when quit query is passed
func QueryDeal(clients *[]cp.Client, cli cp.ClientListen, name string) cp.ClientListen {

	delete(cli.PeerIP, name)
	var j int
	for i := 0; i < len(cli.List); i++ {
		if cli.List[i] == name {
			j = i
			break
		}
	}

	cli.List = append(cli.List[:j], cli.List[j+1:]...)

	*clients = RemoveFromClient(*clients, name)
	return cli
}
