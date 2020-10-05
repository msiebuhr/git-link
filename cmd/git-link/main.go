package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/msiebuhr/git-link" // gitlink
)

var remoteName = flag.String("remote", "origin", "What remote to genenrate links for")
var doOpen = flag.Bool("open", false, "Immediately open in browser")

func main() {
	flag.Parse()

	// Get  remote url by `git config --get remote.$REMOTE_NAME.url`
	out, err := exec.Command("git", "config", "--get", fmt.Sprintf("remote.%s.url", *remoteName)).Output()
	if err != nil {
		log.Fatal(err)
	}
	url := strings.Trim(string(out), "\n")

	repos, err := gitlink.Extract(url)
	if err != nil {
		log.Fatal(err)
	}

	// No arguments? Then just dump a link to browse origin
	if flag.NArg() == 0 {
		link := repos.GetHTTPLink()
		fmt.Println(link)

		if *doOpen {
			cmd := exec.Command("xdg-open", link)
			cmd.Start()
		}
		return
	}

	// Loop through arguments to find out what they represent
	for i := 0; i < flag.NArg(); i += 1 {
		arg := flag.Arg(i)
		link := ""

		if _, err := os.Stat(arg); err == nil {
			// File exists?
			link = repos.GetFileLink("master", arg)
			continue
		} else
		// Parseable as a commit?
		if out, err := exec.Command("git", "rev-parse", arg).Output(); err == nil {
			fmt.Println(repos.GetCommitLink(string(out)))
		}

		if link != "" {
			fmt.Println(link)
			if *doOpen {
				cmd := exec.Command("xdg-open", link)
				cmd.Start()
			}
		}
	}
}
