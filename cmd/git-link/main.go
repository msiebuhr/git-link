package main

import (
	"flag"
	"fmt"
	"log"
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

	if flag.NArg() == 0 {
		fmt.Println(repos.GetHTTPLink())
	} else {
		fmt.Println(repos.GetCommitLink(flag.Arg(0)))
	}
}
