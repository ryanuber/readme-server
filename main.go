package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"gopkg.in/fsnotify.v1"
)

const filePath = "README.md"

func main() {
	var port int

	flags := flag.NewFlagSet("packer-build-manager", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.Usage = usage

	flags.IntVar(&port, "port", 5678, "port number")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Printf("%s", err)
		return
	}

	os.Exit(run(port))
}

func run(port int) int {
	changeCh := make(chan string)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, <-changeCh)
	})

	s := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	go s.ListenAndServe()

	url := fmt.Sprintf("http://localhost:%d", port)
	if err := exec.Command("open", url).Run(); err != nil {
		log.Printf("%s", err)
		return 1
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("%s", err)
		return 1
	}
	defer watcher.Close()

	if err := sendChanges(filePath, changeCh); err != nil {
		log.Printf("%s", err)
		return 1
	}

	for {
		if err := watcher.Add(filePath); err != nil {
			log.Printf("%s", err)
			return 1
		}

		select {
		case event := <-watcher.Events:
			time.Sleep(10 * time.Millisecond)

			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
			case event.Op&fsnotify.Create == fsnotify.Create:
			case event.Op&fsnotify.Rename == fsnotify.Rename:
			default:
				continue
			}

			log.Printf("Detected file change")
			if err := sendChanges(filePath, changeCh); err != nil {
				log.Printf("%s", err)
				return 1
			}
		case err := <-watcher.Errors:
			log.Printf("%s", err)
			return 1
		}
	}

	return 0
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

func usage() {
	helpText := `
usage: readme [options]

Starts an HTTP server to display live updates to your README file

Options:
  -port=<number>  The port number to start the server on.
`
	log.Printf(strings.TrimSpace(helpText))
}
