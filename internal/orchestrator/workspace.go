package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
//
// videoURL, when non-empty, is appended to the rendered problem.txt as
// a dim one-line footer -- the user's chosen placement ("problem footer
// + submit"): always visible, spoiler-control left to them.
//
// draftDir, when non-empty, has any solution.* files it contains (see
// internal/draft.Snapshot, which is what populates it) overlaid on top
// of the fresh copy of repoPath -- resuming a previous session's saved
// progress instead of the exercise's pristine starter code. An empty
// draftDir, or one that doesn't exist yet / holds no solution.* files
// (an exercise with no draft saved), is a complete no-op: today's
// behavior, unchanged.
func PrepareWorkspace(repoPath, videoURL, draftDir string) (workspaceDir string, cleanup func(), err error) {
	dir, err := os.MkdirTemp("", workspaceDirPrefix)
	if err != nil {
		return "", nil, fmt.Errorf("orchestrator: create workspace dir: %w", err)
	}
	cleanup = func() { os.RemoveAll(dir) }

	if err := copyTree(repoPath, dir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("orchestrator: copy repo into workspace: %w", err)
	}

	if err := overlayDraft(draftDir, dir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("orchestrator: overlay draft: %w", err)
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
		if videoURL != "" {
			text = strings.TrimRight(text, "\n") + "\n\n─────\nsolution video (spoilers!): " + videoURL + "\n"
		}
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

// overlayDraft copies every solution.* file directly inside draftDir
// on top of workspaceDir, overwriting the starter files copyTree just
// placed there with the user's saved in-progress edits. A no-op when
// draftDir is empty -- guarded explicitly rather than left to
// filepath.Glob, since joining "" onto a pattern would silently glob
// the process's own working directory instead of doing nothing. A
// draftDir that doesn't exist on disk yet (no draft saved for this
// exercise) is likewise a no-op: Glob against a missing directory
// matches nothing and returns no error.
func overlayDraft(draftDir, workspaceDir string) error {
	if draftDir == "" {
		return nil
	}

	matches, err := filepath.Glob(filepath.Join(draftDir, "solution.*"))
	if err != nil {
		return fmt.Errorf("glob draft solution files: %w", err)
	}
	for _, src := range matches {
		info, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("stat draft file %s: %w", src, err)
		}
		if info.IsDir() {
			continue
		}
		dst := filepath.Join(workspaceDir, filepath.Base(src))
		if err := copyFile(src, dst, info.Mode()); err != nil {
			return fmt.Errorf("overlay draft file %s: %w", src, err)
		}
	}
	return nil
}
