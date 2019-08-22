package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

//SendFileRequestToPeer To send the request to corresponding peer
func SendFileRequestToPeer(connection net.Conn, fileRequest FileRequest, requestType string) {
	//handle with care, FilePartInfo field for this struct is Nil. Will throw seg fault if accessed

	// encoding as a baseRequest
	someRequest := BaseRequest{RequestType: requestType, FileRequest: fileRequest}
	// encoding the request over json
	encoder1 := json.NewEncoder(connection)
	encoder1.Encode(&someRequest)
	connection.Close() // closing connection after the requet is sent
}

//RequestSomeFile To broadcast the file request to everyone on network
func RequestSomeFile(activeClient *ClientListen, myname string, fileName string) string {

	// creating the object to hold queries related to file
	fileRequest := FileRequest{Query: "ask_for_file",
		MyAddress: activeClient.PeerIP[myname],
		MyName:    myname, RequestedFile: fileName}

	// sending to all peers
	for names := range activeClient.PeerListenPort {
		if names != myname {
			tempSplit := strings.Split(activeClient.PeerIP[names], ":")
			connection, err := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[names])
			// diling for finite number of times
			count := 0
			for err != nil {
				fmt.Println("Error in dialing to: ", names, " dialing again...")
				connection1, err1 := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[names])
				connection = connection1
				err = err1
				count++
				if count > 10 {
					break
				}
			}
			// sending the request
			SendFileRequestToPeer(connection, fileRequest, "ask_for_file")
		}
	}
	// ConnectionKey := en.PerformHandshake(conn, PublicKey)
	requestStatus := "completed"
	return requestStatus

}

//GetRequestedFile To get file from a peer
func GetRequestedFile(activeClient *ClientListen, myname string, fileSenderName string, fileName string) string {

	// creating the object to hold queries related to file
	fileRequest := FileRequest{Query: "receive_from_peer",
		MyAddress: activeClient.PeerIP[myname],
		MyName:    myname, RequestedFile: fileName}
	tempSplit := strings.Split(activeClient.PeerIP[fileSenderName], ":")

	// sending to respective peer
	connection, err := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[fileSenderName])
	count := 0
	for err != nil {
		fmt.Println("Error in dialing to: ", fileSenderName, " dialing again...")
		connection1, err1 := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[fileSenderName])
		connection = connection1
		err = err1
		count++
		if count > 10 {
			break
		}
	}

	SendFileRequestToPeer(connection, fileRequest, "receive_from_peer")
	requestStatus := "completed"
	return requestStatus

}

//FileSenderCredentials getting the details for sending file request
func FileSenderCredentials(broadcast bool) (string, string) {

	// getting details of file Sender
	if broadcast == false {
		var fileSenderName string // is the person who will send the file
		fmt.Print("Whom do you want to receive the file from ? : ")
		fmt.Scanln(&fileSenderName)
		var fileName string
		fmt.Print("What file do you want ? ")
		fmt.Scanln(&fileName) // file we want to receive

		return fileSenderName, fileName
	}

	// if message is to be broadcasted, just get file name
	var fileName string
	fmt.Print("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive

	return "None", fileName
}
