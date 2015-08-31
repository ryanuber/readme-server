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
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"gopkg.in/fsnotify.v1"
)

func main() {
	var dontOpen bool
	var port int

	flags := flag.NewFlagSet("readme", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.Usage = usage

	flags.IntVar(&port, "port", 5678, "port number")
	flags.BoolVar(&dontOpen, "dont-open", false, "dont open")

	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Printf("%s", err)
		return
	}

	// Check if we were provided the markdown path
	file := ""
	if args := flags.Args(); len(args) > 0 {
		file = args[0]
	}

	os.Exit(run(port, dontOpen, file))
}

func run(port int, dontOpen bool, filePath string) int {
	changeCh := make(chan string)
	go serve(port, changeCh)

	if !dontOpen {
		url := fmt.Sprintf("http://localhost:%d", port)
		if err := exec.Command("open", url).Run(); err != nil {
			log.Printf("%s", err)
			return 1
		}
	}

	if filePath == "" {
		var err error
		filePath, err = findReadme()
		if err != nil {
			log.Printf("%s", err)
			return 1
		}
	}

	if err := sendChanges(filePath, changeCh); err != nil {
		log.Printf("%s", err)
		return 1
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("%s", err)
		return 1
	}
	defer watcher.Close()

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

func serve(port int, changeCh <-chan string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/" {
			path := filepath.Join(".", req.URL.String())
			if _, err := os.Stat(path); err == nil {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					log.Printf("%s", err)
				} else {
					io.WriteString(w, string(content))
				}
			}
			return
		}
		io.WriteString(w, <-changeCh)
	})

	s := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	if err := s.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("error serving: %s", err))
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

func findReadme() (string, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "README") {
			return file.Name(), nil
		}
	}

	return "", fmt.Errorf("No README files found")
}

func usage() {
	helpText := `
usage: readme [options] [file]

Starts an HTTP server to display live updates to your README file.
By default, the current working directory is searched for a file
with the prefix 'README'. The file can be specified explicilty by
passing an additional argument with the path to the desired file.

Options:
  -port=<number>  The port number to start the server on.
  -dont-open      Do not automatically open the page in a browser
`
	os.Stderr.WriteString(strings.TrimSpace(helpText) + "\n")
	os.Exit(1)
}
