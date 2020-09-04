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
		fmt.Println(repos.GetHTTPLink())
		return
	}

	// Loop through arguments to find out what they represent
	for i := 0; i < flag.NArg(); i += 1 {
		arg := flag.Arg(i)
		// File exists?
		if _, err := os.Stat(arg); err == nil {
			fmt.Println(repos.GetFileLink("master", arg))
			continue
		}

		// Parseable as a commit?
		if out, err := exec.Command("git", "rev-parse", arg).Output(); err == nil {
			fmt.Println(repos.GetCommitLink(string(out)))
		}

	}
}
