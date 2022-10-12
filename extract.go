package gitlink

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type HostingKind uint8

const (
	HK_UNDECIDED HostingKind = iota
	HK_UNKNOWN
	HK_GITHUB
	HK_GITLAB
)

type Repository struct {
	Hostname     string
	Organisation string
	Repository   string

	// What kind of host are we dealing with
	hostKind HostingKind

	// Some local state
	BranchName string
}

func (r *Repository) GetHostingKind() HostingKind {
	if r.hostKind != HK_UNDECIDED {
		return r.hostKind
	}

	// Detect hostKind
	if strings.Contains(r.Hostname, "gitlab") {
		r.hostKind = HK_GITLAB
	} else if strings.Contains(r.Hostname, "github") {
		r.hostKind = HK_GITHUB
	} else {
		r.hostKind = HK_UNKNOWN
	}

	return r.hostKind
}

// HTTP Link returns an educated guess at where the repository can be found
func (r Repository) GetHTTPLink() string {
	return fmt.Sprintf("https://%s/%s/%s", r.Hostname, r.Organisation, r.Repository)
}

// HTTP Commit link
func (r Repository) GetCommitLink(sha string) string {
	// TODO: Try resolving commitish to a real SHA?
	if r.GetHostingKind() == HK_GITLAB {
		return fmt.Sprintf(
			"https://%s/%s/%s/-/commit/%s",
			r.Hostname, r.Organisation, r.Repository, sha)
	}

	if r.GetHostingKind() == HK_GITHUB {
		return fmt.Sprintf(
			"https://%s/%s/%s/commit/%s",
			r.Hostname, r.Organisation, r.Repository, sha)
	}

	return "UNKNOWN"
}

// Get a link to a file
func (r Repository) GetFileLink(filename string) string {
	if strings.HasPrefix(filename, "./") {
		filename = filename[2:]
	}

	// Get current SHA
	git := Git{}
	sha, err := git.GetCurrentCommitSHA()
	if err != nil {
		panic(err)
	}
	sha = strings.TrimSpace(sha)
	// TODO: Try resolving commitish to a real SHA?
	if r.GetHostingKind() == HK_GITLAB {
		return fmt.Sprintf(
			"https://%s/%s/%s/-/blob/%s/%s",
			r.Hostname, r.Organisation, r.Repository, sha, filename)
	}

	if r.GetHostingKind() == HK_GITHUB {
		return fmt.Sprintf(
			"https://%s/%s/%s/blob/%s/%s",
			r.Hostname, r.Organisation, r.Repository, sha, filename)
	}

	return "UNKNOWN"
}

func Extract(gitlink string) (Repository, error) {
	if strings.HasPrefix(gitlink, "https://") {
		u, err := url.Parse(gitlink)

		if err != nil {
			return Repository{}, err
		}

		// Strip leading / from path
		path := strings.TrimPrefix(u.Path, "/")

		// Find last /
		lastSlash := strings.LastIndex(path, "/")

		repository := path[lastSlash+1 : len(path)]
		repository = strings.TrimSuffix(repository, ".git")

		return Repository{
			Hostname:     u.Hostname(),
			Organisation: path[0:lastSlash],
			Repository:   repository,
		}, nil

	}
	// A SSH remote link, perhaps?
	if strings.HasPrefix(gitlink, "git@") {
		hostAndReposSplitIndex := strings.LastIndex(gitlink, ":")
		lastSlash := strings.LastIndex(gitlink, "/")

		repository := gitlink[lastSlash+1 : len(gitlink)]
		repository = strings.TrimSuffix(repository, ".git")

		return Repository{
			Hostname:     gitlink[4:hostAndReposSplitIndex], // Grab 'git@[xxx]:org/repos'
			Organisation: gitlink[hostAndReposSplitIndex+1 : lastSlash],
			Repository:   repository,
		}, nil
	}

	return Repository{}, errors.New("Unknown repository link style")
}
