package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
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
			files, err := ioutil.ReadDir(".")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var mdFiles []string
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".md") {
					mdFiles = append(mdFiles, f.Name())
				}
			}

			for _, f := range mdFiles {
				content, err := ioutil.ReadFile(f)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				html := blackfriday.MarkdownCommon(content)
				out := fmt.Sprintf("%s\n%s\n%s", header+"<!-- yo -->", html, footer)

				dest := filepath.Join(dir, "index.html")
				if err := ioutil.WriteFile(dest, []byte(out), 0600); err != nil {
					fmt.Println(err.Error())
					return
				}

				changes <- out
			}
			time.Sleep(10 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, <-changes)
	})

	s := &http.Server{Addr: ":8080"}
	go s.ListenAndServe()

	time.Sleep(1 * time.Second)

	if err := exec.Command("open", "http://localhost:8080").Run(); err != nil {
		fmt.Println(err.Error())
		return
	}

	doneCh := make(chan struct{})
	<-doneCh
}
