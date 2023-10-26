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
func (r Repository) GetHTTPURL() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   r.Hostname,
		Path:   r.Repository,
	}
}

func (r Repository) GetHTTPLink() string {
	return fmt.Sprintf("https://%s/%s/%s", r.Hostname, r.Organisation, r.Repository)
}

// HTTP Commit link
func (r Repository) GetCommitURL(sha string) *url.URL {
	switch r.GetHostingKind() {
	case HK_GITLAB:
		return &url.URL{
			Scheme: "https",
			Host:   r.Hostname,
			Path: fmt.Sprintf(
				"%s/%s/-/commit/%s",
				r.Organisation, r.Repository, sha,
			),
		}
	case HK_GITHUB:
		return &url.URL{
			Scheme: "https",
			Host:   r.Hostname,
			Path: fmt.Sprintf(
				"%s/%s/commit/%s",
				r.Organisation, r.Repository, sha,
			),
		}
	default:
		return nil
	}
}

func (r Repository) GetCommitLink(sha string) string {
	u := r.GetCommitURL(sha)
	if u == nil {
		return "UNKNOWN"
	}
	return u.String()
}

// Get a link to a file
func (r Repository) GetFileURL(filename string) *url.URL {
	if strings.HasPrefix(filename, "./") {
		filename = filename[2:]
	}

	// Get current SHA
	git := Git{}
	sha, err := git.GetCurrentCommitSHA()
	if err != nil {
		panic(err)
	}
	// TODO: Try resolving commitish to a real SHA?
	u := &url.URL{Scheme: "https", Host: r.Hostname}
	switch r.GetHostingKind() {
	case HK_GITLAB:
		u.Path = fmt.Sprintf(
			"%s/%s/-/blob/%s/%s",
			r.Organisation, r.Repository, sha, filename,
		)
	case HK_GITHUB:
		u.Path = fmt.Sprintf(
			"%s/%s/blob/%s/%s",
			r.Organisation, r.Repository, sha, filename,
		)
	default:
		return nil
	}
	return u
}

func (r Repository) GetFileLink(filename string) string {
	u := r.GetFileURL(filename)
	if u == nil {
		return "UNKNOWN"
	}
	return u.String()
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
