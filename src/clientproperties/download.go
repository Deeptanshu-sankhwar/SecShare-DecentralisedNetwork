package clientproperties

import (
	"bytes"
	"fmt"
	"github.com/andlabs/ui"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WriteCounter - Count of writer
type WriteCounter struct {
	Total uint64
}

var counter WriteCounter

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	// wc.PrintProgress()
	return n, nil
}

var helpMap = map[string]*WriteCounter{}
var lenMap = map[string]int{}
var mainwin = map[string]*ui.Window{}
var prog = map[string]*ui.ProgressBar{}
var name string

//DummyAsync simulates Async
func DummyAsync(wg *sync.WaitGroup, client *http.Client, start int64, end int, i int, size int, url string, f *os.File) error {
	helpMap[name] = &WriteCounter{}
	err := AsyncDownloader(wg, client, start, end, i, size, url, f)
	return err
}

//AsyncDownloader downloader function
func AsyncDownloader(wg *sync.WaitGroup, client *http.Client, start int64, end int, i int, size int, url string, f *os.File) error {

	if end > size {
		end = size
	}
	var fname string
	fname = name
	startString := strconv.FormatInt(start, 10)

	endString := strconv.Itoa(end)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", "bytes="+startString+"-"+endString)
	res, err2 := client.Do(req)
	if err2 != nil {
		return err2
	}
	f.Seek(start, 0)
	var buf bytes.Buffer
	io.Copy(&buf, io.TeeReader(res.Body, helpMap[fname]))
	// io.Copy(&buf, res.Body)
	var buffer []byte
	buffer = buf.Bytes()
	f.Write(buffer)

	wg.Done()
	return nil
}

//set is used to update progressbar
func set(ip *ui.ProgressBar) {
	var fname string
	fname = name
	lenth := lenMap[fname]
	for int(helpMap[fname].Total)+1 < lenth {

		g := int(helpMap[fname].Total) * 100
		val := int(g / lenth)
		time.Sleep(200 * time.Millisecond)
		if val > 90 {
			break
		}
		ip.SetValue(val)
	}

}

//Download ..
func Download(args string) {

	s := strings.Split(args, " ")
	url := s[0]
	// fmt.Println(s)
	tempSplit := strings.Split(url, "/")
	flen := len(tempSplit)
	name = tempSplit[flen-1]
	var fname string
	fname = name
	path := s[1] + "/" + name
	helpMap[fname] = &WriteCounter{}
	client := &http.Client{}
	resp, _ := http.Head(url)
	length := resp.Header.Get("content-length")
	lenh, _ := strconv.Atoi(length)
	lenMap[fname] = lenh
	dummy := make([]byte, lenh)
	ioutil.WriteFile(fmt.Sprintf(path), dummy, 0644)
	f, _ := os.OpenFile(path, os.O_RDWR, 0666)
	var start int64
	start = 0
	end := 0
	partLength := int(lenh / 4)
	end = partLength
	var wg sync.WaitGroup
	wg.Add(4)
	// for testing use http://file-examples.com/wp-content/uploads/2017/11/file_example_MP3_1MG.mp3
	for i := 0; i < 4; i++ {

		go AsyncDownloader(&wg, client, start, end, i, lenh, url, f)
		start = int64(end) + 1
		end = end + partLength

	}
	go ui.Main(setupUI)
	wg.Wait()
	mainwin[fname].Destroy() // waiting for the goroutines to finish
	ui.Quit()
}

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	var fname string
	fname = name
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	prog[fname] = ui.NewProgressBar()
	hbox.Append(prog[fname], false)
	go set(prog[fname])

	return vbox
}

func setupUI() {
	var fname string
	fname = name
	mainwin[fname] = ui.NewWindow("Downloading"+fname, 300, 150, true)
	mainwin[fname].OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin[fname].Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin[fname].SetChild(tab)
	mainwin[fname].SetMargined(true)

	tab.Append("Downloading", makeBasicControlsPage())
	tab.SetMargined(0, true)
	// go set(tab)
	mainwin[fname].Show()
}
