package main

import (
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	en "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/encryptionproperties"
	// cp "../src/clientproperties"
	// en "../src/encryptionproperties"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// function to fetch peers which are currently active
func gettingPeersFromServer(c net.Conn, peers *[]string, msg *cp.ClientListen) {
	for {
		// json decoder to decode the struct obtained over the connection
		d := json.NewDecoder(c)
		d.Decode(msg)
	}
}

func main() {

	argsWithoutProg := os.Args[1:] //taking the args from the user

	var activeClient cp.ClientListen  // to store information about peers
	var directoryFiles cp.ClientFiles // to store information about files in the directory
	fileDirectory := "../files"

	// map to keep account of the files received
	myfiles := make(map[string]cp.MyReceivedFiles)
	// struct - containing counter, which gives the index from which we have to read new messages
	// and an array of struct type - MessageRequest to store the attributes of messages
	mymessages := cp.MyReceivedMessages{Counter: 0}

	// to store current peers
	peers := []string{}

	// credentials of the client logging in
	var name string
	var query string
	var listenPort string // which port does he prefer to listen upon
	var flag = false

	// generating keys for connection with server
	_, PublicKey := en.GenerateKeyPair()
	addr := argsWithoutProg[2]
	conn, err := net.Dial("tcp", addr)

	// limiting the dial count
	dialCount := 0
	for err != nil {
		// fmt.Println("error in connecting to server, dialing again")
		conn1, err1 := net.Dial("tcp", "192.168.111.252:8081")
		conn = conn1
		err = err1
		dialCount++
		if dialCount > 200 {
			fmt.Println("Apparently server's port is not open...")
			os.Exit(1)
		}
	}

	//assigning names and ports.
	name = argsWithoutProg[0]
	listenPort = argsWithoutProg[1]
	// performing handshake with server
	ServerKey := en.PerformHandshake(conn, PublicKey)

	// queries currently supported
	fmt.Println("The follwing queries are supported ->")
	fmt.Println("for quitting - quit")
	fmt.Println("for broadcasting a file request - ask_for_file")
	fmt.Println("for receiving file - receive_file")
	fmt.Println("for sending message - send_message")
	fmt.Println("for displaying recent messages - display_recent_messages")
	fmt.Println("for displaying current peers - display_peers")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		query = "quit"
		// encrypting credentials to notify server, when a client wants to quit
		nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
		mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
		// sending the information to server
		queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)
		fmt.Print(query)
		cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
		os.Exit(1)

	}()

	for {
		// getting others clients who are currenty active
		go gettingPeersFromServer(conn, &peers, &activeClient)

		// flag == false, signifies the client has to login
		if flag == false {

			// fmt.Print("Enter your credentials : ")
			// fmt.Scanln(&name)
			flag = true // has logged in
			query = "login"

			// encrypting the details with the PublicKey
			nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
			queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)

			ln, err := net.Listen("tcp", ":"+listenPort)
			if err == nil {
				fmt.Println("[SUCCESS] Successfully logged in")
			}
			for err != nil {
				fmt.Println("Cant listen on this port, choose another : ")
				os.Exit(2)

			}

			// adds files in directory to a clientFiles
			error1 := filepath.Walk(fileDirectory, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				directoryFiles.FilesInDir = append(directoryFiles.FilesInDir, info.Name())
				return nil
			})

			if error1 != nil {
				panic(error1)
			}

			mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
			// sending credentials to server
			cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
			// activating my own port for listening
			go cp.ListenOnSelfPort(ln, name, &activeClient, myfiles, &mymessages, &directoryFiles)
			continue

		} else {
			// accepting further queries, after login is done
			fmt.Print("[QUERY]>>")
			fmt.Scanln(&query)

			if query == "quit" {
				// encrypting credentials to notify server, when a client wants to quit
				nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
				queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)
				mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
				// sending the information to server
				cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
				os.Exit(2)

			} else if query == "ask_for_file" {
				// broadcasting request for receiving some file
				_, fileName := cp.FileSenderCredentials(true)
				requestStatus := cp.RequestSomeFile(&activeClient, name, fileName)
				if requestStatus == "completed" {
					fmt.Println("Request broadcasted properly")
				} else {
					fmt.Println("Request not broadcasted properly")
				}

				time.Sleep(2 * time.Second)
				// display the status of file existence from all clients
				cp.DisplayRecentUnseenMessages(&mymessages)

			} else if query == "receive_file" {
				// Receiving file from specific peer
				// getting credentials of the person from whom to receive the file
				fileSenderName, fileName := cp.FileSenderCredentials(false)
				// sending the request to receive some file
				// status of the request, whether or not the request is sent properly
				requestStatus := cp.GetRequestedFile(&activeClient, name, fileSenderName, fileName)
				if requestStatus == "completed" {
					fmt.Println("Request sent")
				} else {
					fmt.Println("Request not sent properly")
				}
			} else if query == "send_message" {
				// query to send messages to others
				// activating the message send mode
				messaging := true
				fmt.Println("Currently in messaging mode..")

				for messaging == true {
					// getting credentials of the person, to whom I want to send some message
					messageReceiverName, message := cp.MessageReceiverCredentials()
					// status of the message, whether or not sent properly
					messageStatus := cp.RequestMessage(&activeClient, name, messageReceiverName, message)
					if messageStatus == "sent" {
						fmt.Println("Message sent")
					} else {
						fmt.Println("Message not sent properly")
					}

					// whether I want to exit the message mode
					var queryMessage string
					fmt.Print("Do you want to send more messages? If Yes type Y, else N: ")
					fmt.Scan(&queryMessage)
					if queryMessage == "N" {
						fmt.Println("Exiting messaging mode...")
						break
					}
				}

			} else if query == "display_recent_messages" {
				// to display recent messages, which haven't been seen yet
				fmt.Println("Display recent unseen messages - (type) 1")
				fmt.Println("Display recent Num messages - (type) 2")
				var queryMessage string
				fmt.Scanln(&queryMessage)

				// Display recently unseen messages
				if queryMessage == "1" {
					_ = cp.DisplayRecentUnseenMessages(&mymessages)
				} else {
					// display N recent messages
					var num int
					fmt.Print("Number of recent messages you want to see : ")
					fmt.Scan(&num)
					_ = cp.DisplayNumRecentMessages(&mymessages, num)
				}

			} else if query == "down" {
				// to download files, support within file concurrency and can donwload muliple files simultaneously
				fmt.Print("<URL> <filepath>") // url string
				var url string
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan() // use `for scanner.Scan()` to keep reading
				url = scanner.Text()
				go cp.Download(url)

			} else if query == "display_peers" {
				fmt.Println("Active peers : ", activeClient)
			}
		}
	}
}
