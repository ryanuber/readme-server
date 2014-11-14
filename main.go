package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/russross/blackfriday"
)

func main() {
	changes := make(chan string)

	dir, err := ioutil.TempDir("", "ms")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		for {
			content, err := ioutil.ReadFile("README.md")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			html := blackfriday.MarkdownCommon(content)
			out := fmt.Sprintf("%s\n%s\n%s", header, html, footer)

			dest := filepath.Join(dir, "index.html")
			if err := ioutil.WriteFile(dest, []byte(out), 0600); err != nil {
				fmt.Println(err.Error())
				return
			}

			changes <- out
			time.Sleep(5 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, <-changes)
	})

	s := &http.Server{Addr: ":8080"}
	go s.ListenAndServe()

	if err := exec.Command("open", "http://localhost:8080").Run(); err != nil {
		fmt.Println(err.Error())
		return
	}

	doneCh := make(chan struct{})
	<-doneCh
}
