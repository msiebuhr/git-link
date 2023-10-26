package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/msiebuhr/git-link" // gitlink
)

var remoteName = flag.String("remote", "origin", "What remote to genenrate links for")
var doOpen = flag.Bool("open", false, "Immediately open in browser")
var isTerminal bool
var outputMarkdown bool

func init() {
	isTerminalDefault := false
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		isTerminalDefault = (fileInfo.Mode()&os.ModeCharDevice != 0)
	}
	flag.BoolVar(&isTerminal, "term", isTerminalDefault, "Is comment running in terminal")
	flag.BoolVar(&outputMarkdown, "md", false, "Output links in markdown-syntax")

}

// Formats a URL as a Hyperlink in a terminal
func format_url(raw_name string, u *url.URL) string {
	if outputMarkdown {
		return fmt.Sprintf("[%s/%s](%s)", u.Host, u.Path, u.String())
	} else if isTerminal {
		return fmt.Sprintf("\x1b]8;;%s\x1b\\%s/%s\x1b]8;;\x1b\\", u.String(), u.Host, u.Path)
	} else {
		return u.String()
	}
}

func main() {
	flag.Parse()

	// Get  remote url by `git config --get remote.$REMOTE_NAME.url`
	out, err := exec.Command("git", "config", "--get", fmt.Sprintf("remote.%s.url", *remoteName)).Output()
	if err != nil {
		log.Fatal(err)
	}
	u := strings.Trim(string(out), "\n")

	repos, err := gitlink.Extract(u)
	if err != nil {
		log.Fatal(err)
	}

	// No arguments? Then just dump a link to browse origin
	if flag.NArg() == 0 {
		link := repos.GetHTTPURL()
		fmt.Println(format_url(link.String(), link))

		if *doOpen {
			cmd := exec.Command("xdg-open", link.String())
			cmd.Start()
		}
		return
	}

	// Loop through arguments to find out what they represent
	for i := 0; i < flag.NArg(); i += 1 {
		arg := flag.Arg(i)
		link := &url.URL{}

		if _, err := os.Stat(arg); err == nil {
			// File exists?
			link = repos.GetFileURL(arg)
		} else if out, err := exec.Command("git", "rev-parse", arg).Output(); err == nil {
			link = repos.GetCommitURL(strings.TrimSpace(string(out)))
		}

		if link != nil {
			fmt.Println(format_url(arg, link))
			if *doOpen {
				cmd := exec.Command("xdg-open", link.String())
				cmd.Start()
			}
		}
	}
}
