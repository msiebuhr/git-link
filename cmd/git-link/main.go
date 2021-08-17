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
var isTerminal bool

func init() {
	isTerminalDefault := false
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		isTerminalDefault = (fileInfo.Mode()&os.ModeCharDevice != 0)
	}
	flag.BoolVar(&isTerminal, "term", isTerminalDefault, "Is comment running in terminal")
}

// Formats a URL as a Hyperlink in a terminal
func format_url(url string) string {
	if isTerminal {
		return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, url)
	} else {
		return url
	}
}

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
		fmt.Println(format_url(link))

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
		} else if out, err := exec.Command("git", "rev-parse", arg).Output(); err == nil {
			link = repos.GetCommitLink(strings.TrimSpace(string(out)))
		}

		if link != "" {
			fmt.Println(format_url(link))
			if *doOpen {
				cmd := exec.Command("xdg-open", link)
				cmd.Start()
			}
		}
	}
}
