package main

import (
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	// cp "../src/clientproperties"
	"testing"
	"time"
)

//TestDisplayRecentUnseenMessages - tests the code to display recent unseen messages
func TestDisplayRecentUnseenMessages(t *testing.T) {

	mymessages := cp.MyReceivedMessages{Counter: 0}

	messageRequest1 := cp.MessageRequest{
		SenderQuery: "message_request", SenderName: "abc", Message: "Hey There!!"}

	messageRequest2 := cp.MessageRequest{
		SenderQuery: "message_request", SenderName: "def", Message: "Hey There, again!!"}

	mymessages.MyMessages = append(mymessages.MyMessages, messageRequest1)
	mymessages.MyMessages = append(mymessages.MyMessages, messageRequest2)

	startTime := time.Now()
	displayStatus := cp.DisplayRecentUnseenMessages(&mymessages)
	endTime := time.Now()

	if displayStatus != "Displayed" || endTime.Sub(startTime) > 1000000 {
		t.Fatal("Problem in displaying messages...")
	}
}

//TestDisplayRecentNumMessages - tests the code to display few recent messages
func TestDisplayRecentNumMessages(t *testing.T) {

	mymessages := cp.MyReceivedMessages{Counter: 0}

	messageRequest1 := cp.MessageRequest{
		SenderQuery: "message_request", SenderName: "abc", Message: "Hey There!!"}

	messageRequest2 := cp.MessageRequest{
		SenderQuery: "message_request", SenderName: "def", Message: "Hey There, again!!"}

	mymessages.MyMessages = append(mymessages.MyMessages, messageRequest1)
	mymessages.MyMessages = append(mymessages.MyMessages, messageRequest2)

	startTime := time.Now()
	displayStatus := cp.DisplayNumRecentMessages(&mymessages, 3)
	endTime := time.Now()

	if displayStatus != "Displayed" || endTime.Sub(startTime) > 1000000 {
		t.Fatal("Problem in displaying messages...")
	}
}
