package gitlink

import "testing"

func TestRepostioryGetHostingKind(t *testing.T) {
	github := Repository{Hostname: "github.com"}
	if github.GetHostingKind() != HK_GITHUB {
		t.Errorf("Unexpected HostingKind %d for %v, want %d", github.GetHostingKind(), github, HK_GITHUB)
	}
	gitlab_backup := Repository{Hostname: "gitlab-backup.example.com"}
	if gitlab_backup.GetHostingKind() != HK_GITLAB {
		t.Errorf("Unexpected HostingKind %d for %v, want %d", gitlab_backup.GetHostingKind(), gitlab_backup, HK_GITLAB)
	}
}

func TestGetCommitURL(t *testing.T) {
	x, err := Extract("git@gitlab.example.com:nested/org/repo.git")
	if err != nil {
		t.Fatalf("Could not parse url: %s", err)
	}

	commitUrl := x.GetCommitURL("some-sha-or-other")
	expected := "https://gitlab.example.com/nested/org/repo/-/commit/some-sha-or-other"

	if commitUrl.String() != expected {
		t.Errorf("Expected %#v + some-sha-or-other =>\n%s\ngot\n%s", x, commitUrl, expected)
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		in       string
		out      Repository
		HTTPLink string
	}{
		// Known-good things
		{
			in:       "git@github.com:msiebuhr/foobar.git",
			out:      Repository{Hostname: "github.com", Organisation: "msiebuhr", Repository: "foobar"},
			HTTPLink: "https://github.com/msiebuhr/foobar",
		},
		{
			in:       "https://github.com/msiebuhr/foobar.git",
			out:      Repository{Hostname: "github.com", Organisation: "msiebuhr", Repository: "foobar"},
			HTTPLink: "https://github.com/msiebuhr/foobar",
		},
		{
			in:       "git@gitlab.com:msiebuhr/foobar.git",
			out:      Repository{Hostname: "gitlab.com", Organisation: "msiebuhr", Repository: "foobar"},
			HTTPLink: "https://gitlab.com/msiebuhr/foobar",
		},
		{
			in:       "git@gitlab.com:msiebuhr/suborg/foobar.git",
			out:      Repository{Hostname: "gitlab.com", Organisation: "msiebuhr/suborg", Repository: "foobar"},
			HTTPLink: "https://gitlab.com/msiebuhr/suborg/foobar",
		},

		// Self-hosted stuff
		{
			in:       "git@gitlab.example.com:dept/subdept/repo.git",
			out:      Repository{Hostname: "gitlab.example.com", Organisation: "dept/subdept", Repository: "repo"},
			HTTPLink: "https://gitlab.example.com/dept/subdept/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := Extract(tt.in)

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if out.Hostname != tt.out.Hostname {
				t.Errorf("Unexpected hostname %s, expected %s", out.Hostname, tt.out.Hostname)
			}
			if out.Organisation != tt.out.Organisation {
				t.Errorf("Unexpected organisation %s, expected %s", out.Organisation, tt.out.Organisation)
			}
			if out.Repository != tt.out.Repository {
				t.Errorf("Unexpected repository %s, expected %s", out.Repository, tt.out.Repository)
			}

			httpLink := out.GetHTTPURL().String()
			if httpLink != tt.HTTPLink {
				t.Errorf("Unexpected link %s from %v\nexpected %s", httpLink, out, tt.HTTPLink)
			}
		})
	}
}
