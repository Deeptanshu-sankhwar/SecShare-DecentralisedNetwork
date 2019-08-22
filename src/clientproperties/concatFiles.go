package clientproperties

import (
	"fmt"
	"github.com/andlabs/ui"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var popwin = map[string]*ui.Window{}

// creating waitgroup to wait for all goroutines to finish
var wgConcat sync.WaitGroup
var curr string

// func to write the filepart details to all Files byte slice
func concatFiles(i int, allFiles []byte, filePartContent FilePartContents) {
	defer wgConcat.Done()

	filePartContents := filePartContent.Contents

	for j := i * len(filePartContents); j < i*len(filePartContents)+len(filePartContents); j++ {
		allFiles[j] = filePartContents[j-i*len(filePartContents)]
	}

}

func concatenateFileParts(file MyReceivedFiles) {

	// getting total size of all parts
	var byteSizeLength int
	fileName := file.MyFileName
	fileParts := file.MyFile

	for i := 0; i < len(fileParts); i++ {
		byteSizeLength += len(fileParts[i].Contents)
	}

	// creating new byte slice for creating new file
	allFiles := make([]byte, byteSizeLength)

	// writing the received parts to allFiles slice
	for i := int(0); i < file.FilePartInfo.TotalParts; i++ {
		wgConcat.Add(1)
		go concatFiles(i, allFiles, fileParts[i])
	}

	wgConcat.Wait()

	// writing the received file
	currentfilename := "Received_" + fileName
	ioutil.WriteFile(currentfilename, allFiles, os.ModeAppend)

	// Test File existence.
	_, err := os.Stat(currentfilename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file doesn't exist")
		}
	}

	// Change permissions in Linux.
	err = os.Chmod(currentfilename, 0777)
	if err != nil {
		fmt.Println(err)
	}
	curr = fileName
	go ui.Main(setupPop)
	time.Sleep(4000 * time.Millisecond)
	popwin[curr].Destroy()

}

func setupPop() {
	var fname string
	fname = curr
	popwin[fname] = ui.NewWindow("Receiving_"+fname, 300, 150, true)

	tab := ui.NewTab()
	popwin[fname].SetChild(tab)
	popwin[fname].SetMargined(true)

	tab.Append("Received ", makePopup())
	tab.SetMargined(0, true)

	popwin[fname].Show()
}
func makePopup() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	text := ui.NewEntry()
	text.SetReadOnly(true)
	text.SetText("Received" + curr)
	hbox.Append(text, false)
	return vbox
}
