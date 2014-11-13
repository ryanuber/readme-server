package main

import (
	//"net/http"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

func main() {
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

	dir, err := ioutil.TempDir("", "ms")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, f := range mdFiles {
		content, err := ioutil.ReadFile(f)
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
	}

	indexPath := filepath.Join(dir, "index.html")
	if err := exec.Command("open", indexPath).Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
