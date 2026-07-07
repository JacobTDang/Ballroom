package orchestrator

import (
	"fmt"
	"os"
)

// PrepareWorkspace copies repoPath into a fresh temp directory and returns
// it along with a cleanup func that removes it.
//
// The returned directory — not repoPath — is what gets mounted as
// /workspace. repoPath must never be mounted directly: it's the permanent,
// reusable source of truth for the exercise, and anything written into a
// running session (edits, or hidden tests revealed on submit) must land
// somewhere disposable instead of leaking back into it.
func PrepareWorkspace(repoPath string) (workspaceDir string, cleanup func(), err error) {
	dir, err := os.MkdirTemp("", "practice-workspace-")
	if err != nil {
		return "", nil, fmt.Errorf("orchestrator: create workspace dir: %w", err)
	}
	cleanup = func() { os.RemoveAll(dir) }

	if err := copyTree(repoPath, dir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("orchestrator: copy repo into workspace: %w", err)
	}
	return dir, cleanup, nil
}
