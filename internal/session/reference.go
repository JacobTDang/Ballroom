package session

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Reference requests the exercise's reference solution be revealed and
// waits for it to land in the workspace -- the in-container side of
// `ballroom reference` (bound to M-g). Mirrors Submit's own request/wait
// shape (requestReveal/waitForReady), reusing the same control-directory
// handshake pattern against a different pair of files so the host's
// reveal watcher (orchestrator.WaitAndRevealReference) can run
// independently of the submit watcher and this can be requested before,
// after, or without ever submitting.
func Reference(cfg Config, stdout io.Writer) error {
	if err := requestReferenceReveal(cfg); err != nil {
		return err
	}
	if err := waitForReferenceReady(cfg); err != nil {
		return err
	}
	fmt.Fprintln(stdout, "reference solution revealed at reference/solution.*")
	return nil
}

func requestReferenceReveal(cfg Config) error {
	path := filepath.Join(cfg.ControlDir, "reference.request")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		return fmt.Errorf("session: write reference request: %w", err)
	}
	return nil
}

func waitForReferenceReady(cfg Config) error {
	readyPath := filepath.Join(cfg.ControlDir, "reference.ready")
	deadline := time.Now().Add(cfg.RevealTimeout)

	for time.Now().Before(deadline) {
		if _, err := os.Stat(readyPath); err == nil {
			return nil
		}
		time.Sleep(cfg.PollInterval)
	}
	return fmt.Errorf("session: timed out after %s waiting for the reference solution to be revealed", cfg.RevealTimeout)
}

// ReferenceSolutionPath returns the revealed reference solution's file
// path inside workspaceDir (see orchestrator.WaitAndRevealReference,
// which is what actually copies it into the reference/ subdirectory), or
// a clear error if nothing has been revealed there yet.
func ReferenceSolutionPath(workspaceDir string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(workspaceDir, "reference", "solution.*"))
	if err != nil {
		return "", fmt.Errorf("session: glob reference solution: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("session: no reference solution found in the workspace (has it been revealed yet?)")
	}
	return matches[0], nil
}
