package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/russross/blackfriday"
	"gopkg.in/fsnotify.v1"
)

const filePath = "README.md"

func main() {
	changeCh := make(chan string)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, <-changeCh)
	})

	s := &http.Server{Addr: ":8080"}
	go s.ListenAndServe()

	if err := exec.Command("open", "http://localhost:8080").Run(); err != nil {
		fmt.Println(err.Error())
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer watcher.Close()

	if err := sendChanges(filePath, changeCh); err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		if err := watcher.Add(filePath); err != nil {
			fmt.Println(err.Error())
			return
		}

		select {
		case <-watcher.Events:
			time.Sleep(10 * time.Millisecond)
			if err := sendChanges(filePath, changeCh); err != nil {
				fmt.Println(err.Error())
				return
			}
		case err := <-watcher.Errors:
			fmt.Println(err.Error())
			return
		}
	}
}

func sendChanges(path string, changeCh chan<- string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	html := header
	html += string(blackfriday.MarkdownCommon(content))
	html += footer

	changeCh <- html
	return nil
}
