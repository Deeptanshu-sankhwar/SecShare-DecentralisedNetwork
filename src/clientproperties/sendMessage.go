package clientproperties

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

//sendMessageRequestToPeer encodes the baserequest
func sendMessageToPeer(connection net.Conn, messageRequest MessageRequest) {
	baseRequest := BaseRequest{RequestType: "receive_message", MessageRequest: messageRequest}
	encoder1 := json.NewEncoder(connection)
	encoder1.Encode(&baseRequest)
	connection.Close()
}

// RequestMessage takes message from client and dials to receiver
func RequestMessage(activeClient *ClientListen, name string, messageReceiverName string,
	message string) string {

	messageRequest := MessageRequest{
		SenderQuery: "message_request", SenderName: name,
		SenderAddress: activeClient.PeerIP[name], Message: message}
	tempSplit := strings.Split(activeClient.PeerIP[messageReceiverName], ":")
	connection, err := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[messageReceiverName])

	count := 0
	for err != nil {

		connection1, err1 := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[messageReceiverName])
		connection = connection1
		err = err1
		count++
		if count > 100 {
			fmt.Println("Error in dialing to: ", messageReceiverName, " recheck credentials")
			messageStatus := "not sent"
			return messageStatus
			break
		}
	}

	sendMessageToPeer(connection, messageRequest)
	messageStatus := "sent"
	return messageStatus
}

//MessageReceiverCredentials for receiving message credientials
func MessageReceiverCredentials() (string, string) {

	// getting credentials of the person to send message to
	var messageReceiverName string
	var message string // what message to send?
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Message (Person's name) : ")
	fmt.Scanln(&messageReceiverName)
	fmt.Print("Message to send : ")
	// fmt.Scanln(&message)
	message, err := in.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return messageReceiverName, message
}
