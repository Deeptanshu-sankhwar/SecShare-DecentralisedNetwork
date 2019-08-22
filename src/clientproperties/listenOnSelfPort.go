package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// creating locks for messages and files array
var mutexFiles = &sync.Mutex{}    // Lock and unlock (mutexFiles)
var mutexMessages = &sync.Mutex{} // Lock and unlock (mutexFiles)

// SendPart - To send each file part concurrently
func SendPart(names string, activeClient *ClientListen, newfilerequest FileRequest,
	countSent int, allfileparts []FilePartInfo, wgSplit *sync.WaitGroup) {
	count := 0
	tempSplit := strings.Split(activeClient.PeerIP[names], ":")
	connection, err := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[names])
	for err != nil {
		// fmt.Println("Error in dialing to: ", names, " dialing again...")
		connection1, err1 := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[names])
		connection = connection1
		err = err1
		count++
		if count > 100 {
			// fmt.Println("Error in sending current file part - ", err)
			break
		}
	}

	// sending the file part corresponding to this peer
	// fmt.Println("Connection established to send a file part to connection - ", connection)
	baseRequest := BaseRequest{RequestType: "received_some_file", FileRequest: newfilerequest,
		FilePartInfo: allfileparts[countSent]}
	encoder := json.NewEncoder(connection)
	encoder.Encode(&baseRequest)

	defer wgSplit.Done()
}

//SendFileParts send various file parts to peers
func SendFileParts(newfilerequest FileRequest, allfileparts []FilePartInfo,
	activeClient *ClientListen, myname string) int {

	// waitgroup to wait for all goroutines to finish
	var wgSplit sync.WaitGroup

	countSent := 0
	for names := range activeClient.PeerListenPort {
		if names != myname {
			wgSplit.Add(1) // incrementing the count of waitgroup, for another goroutine is run
			go SendPart(names, activeClient, newfilerequest, countSent, allfileparts, &wgSplit)
			countSent++
		}
	}
	wgSplit.Wait()
	return countSent
}

// Handle request to send some file to a peer
func handleNewFileSendRequest(newfilerequest FileRequest, myname string, activeClient *ClientListen) {
	// getting all splits of file
	allfileparts := GetSplitFile(newfilerequest.RequestedFile, len(activeClient.List))
	_ = SendFileParts(newfilerequest, allfileparts, activeClient, myname)

}

// Handle some received file part for myself
func handleReceivedFile(newrequest BaseRequest, myfiles map[string]MyReceivedFiles) {

	var TotalFileParts int
	var filePartNum int

	requestedFileName := newrequest.FilePartInfo.FileName

	TotalFileParts = newrequest.FilePartInfo.TotalParts
	filePartNum = newrequest.FilePartInfo.PartNumber

	// If already exists in myfies, append the current part to corresponding file struct
	if _, ok := myfiles[requestedFileName]; ok {

		mutexFiles.Lock()
		var x = myfiles[requestedFileName]
		x.PartsReceived++
		myfiles[requestedFileName] = x
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
		mutexFiles.Unlock()
		// if all parts have been received, concatenate it and create new file
		if myfiles[requestedFileName].PartsReceived == TotalFileParts {
			concatenateFileParts(myfiles[requestedFileName])
		}

	} else {
		// creating new received file object for my own file
		myfiles[requestedFileName] = MyReceivedFiles{0, newrequest.FilePartInfo.FileName,
			make([]FilePartContents, newrequest.FilePartInfo.TotalParts),
			newrequest.FilePartInfo}

		// locking
		mutexFiles.Lock()
		var x = myfiles[requestedFileName]
		x.PartsReceived++
		myfiles[requestedFileName] = x
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
		mutexFiles.Unlock()

		// if all parts have been received, concatenate it and create new file
		if myfiles[requestedFileName].PartsReceived == TotalFileParts {
			concatenateFileParts(myfiles[requestedFileName])
		}
	}
}

// Handle a request
func handleReceivedRequest(connection net.Conn, activeClient *ClientListen, myname string,
	myfiles map[string]MyReceivedFiles, newrequest BaseRequest) {

	// If the file receievd is for me
	if newrequest.FileRequest.MyName == myname {
		handleReceivedFile(newrequest, myfiles)

	} else {
		// if file is received to be forwarded to some other peer

		count := 0
		tempSplit := strings.Split(activeClient.PeerIP[newrequest.FileRequest.MyName], ":")
		connection, err := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[newrequest.FileRequest.MyName])
		for err != nil {
			// fmt.Println("Error in dialing to: ", newrequest.FileRequest.MyName, " dialing again...")
			connection1, err1 := net.Dial("tcp", tempSplit[0]+":"+activeClient.PeerListenPort[newrequest.FileRequest.MyName])
			connection = connection1
			err = err1
			count++
			if count > 100 {
				// fmt.Println("Error in forwarding file part - ", err)
				break
			}
		}

		// sending the request further
		newSendRequest := newrequest
		newconn := json.NewEncoder(connection)
		newconn.Encode(&newSendRequest)
		bytesForwarded := len(newrequest.FilePartInfo.FilePartContents)
		t := strconv.Itoa(bytesForwarded)
		fmt.Println("\nforwarded "+t+" bytes to", newrequest.FileRequest.MyAddress)
		fmt.Print("[QUERY]>>")

	}
}

//CheckFileExistence checking existence of file to send it to some peer
func CheckFileExistence(request FileRequest, directoryFiles *ClientFiles) bool {
	// if the file exists
	for _, file := range directoryFiles.FilesInDir {

		if strings.Compare(file, request.RequestedFile) == 0 {
			return true
		}
	}
	return false
}

func handleConnection(connection net.Conn, activeClient *ClientListen, myname string,
	myfiles map[string]MyReceivedFiles, mymessages *MyReceivedMessages, directoryFiles *ClientFiles) {

	var newrequest BaseRequest
	newconn := json.NewDecoder(connection)
	newconn.Decode(&newrequest)

	// If peer is asking for a file, send the existence status of the file in my directory
	if newrequest.RequestType == "ask_for_file" {

		exists := CheckFileExistence(newrequest.FileRequest, directoryFiles)
		if exists {
			_ = RequestMessage(activeClient, myname, newrequest.FileRequest.MyName, "have the file you requested")
		} else {
			_ = RequestMessage(activeClient, myname, newrequest.FileRequest.MyName, "doesn't have the file you requested")
		}

		// if peer is asking me to send some file
	} else if newrequest.RequestType == "receive_from_peer" {

		handleNewFileSendRequest(newrequest.FileRequest, myname, activeClient)

		// if received some file part
	} else if newrequest.RequestType == "received_some_file" {

		handleReceivedRequest(connection, activeClient, myname, myfiles, newrequest)

		// If recereceive_fileievd some message
	} else if newrequest.RequestType == "receive_message" {
		mutexMessages.Lock()
		mymessages.MyMessages = append(mymessages.MyMessages, newrequest.MessageRequest)
		mutexMessages.Unlock()
	}

}

// ListenOnSelfPort listens for clients on network
func ListenOnSelfPort(ln net.Listener, myname string, activeClient *ClientListen, myfiles map[string]MyReceivedFiles,
	mymessages *MyReceivedMessages, directoryFiles *ClientFiles) {
	for {
		connection, err := ln.Accept()

		if err != nil {
			panic(err)
		}
		// Hanling the received connection
		go handleConnection(connection, activeClient, myname, myfiles, mymessages, directoryFiles)
	}
}
