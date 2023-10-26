package gitlink

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git struct {
	Workdir string
}

func (g *Git) GetCurrentCommitSHA() (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("unable to get current commit: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
