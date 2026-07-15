package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobTDang/Ballroom/internal/exercise"
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

	// Exercises that ship a problem.md get a plain-text render written
	// alongside it -- problem.txt is what the editor pane actually opens
	// (docker/entrypoint.sh prefers it), showing the statement as clean
	// structured text instead of raw markdown markers. problem.md itself
	// stays in the workspace untouched: it remains the authoring format
	// and what the tutor's read_problem_statement tool reads. Written to
	// the disposable workspace only, never the source repo.
	md, err := os.ReadFile(filepath.Join(dir, "problem.md"))
	if err == nil {
		text := exercise.RenderProblemText(string(md))
		if err := os.WriteFile(filepath.Join(dir, "problem.txt"), []byte(text), 0o644); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("orchestrator: write problem.txt: %w", err)
		}
	} else if !os.IsNotExist(err) {
		cleanup()
		return "", nil, fmt.Errorf("orchestrator: read problem.md: %w", err)
	}
	return dir, cleanup, nil
}
