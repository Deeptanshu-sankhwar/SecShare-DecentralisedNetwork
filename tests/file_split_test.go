package main

import (
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	"testing"
)

//TestFileSplit test file split
func TestFileSplit(t *testing.T) {

	filename := "image.jpg"
	allFileParts := cp.GetSplitFile(filename, 2)

	if len(allFileParts[0].FilePartContents) == 0 {
		t.Fatal("File not properly written to allFileParts slice")
	}
}
