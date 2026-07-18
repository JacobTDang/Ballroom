package session

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// Reference requests the exercise's reference solution be revealed,
// waits for it to land in the workspace, and records the attempt as
// given up -- the in-container side of `ballroom reference` (bound to
// M-g). Asking to see the answer before solving it IS the outcome
// (issue #238): recording it here, before the caller ever opens the
// file, keeps the stats honest instead of leaving the row for a submit
// that may never come. A reveal that times out records nothing --
// there was no real "give up" moment, just a broken handshake.
//
// Mirrors Submit's own request/wait/log shape (requestReveal/
// waitForReady, then building and logging a tracker.Attempt), reusing
// the same control-directory handshake pattern against a different
// pair of files so the host's reveal watcher
// (orchestrator.WaitAndRevealReference) can run independently of the
// submit watcher and this can be requested before, after, or without
// ever submitting.
func Reference(cfg Config, stdout io.Writer) (tracker.Attempt, error) {
	if err := requestReferenceReveal(cfg); err != nil {
		return tracker.Attempt{}, err
	}
	if err := waitForReferenceReady(cfg); err != nil {
		return tracker.Attempt{}, err
	}

	// The tutor pane's assistance counters for this session -- same
	// graceful degradation as Submit's own use of this (see
	// readTutorState's doc comment).
	ts := readTutorState(cfg.WorkspaceDir)
	attempt := tracker.Attempt{
		ExerciseID:   cfg.ExerciseID,
		Category:     cfg.Category,
		Language:     cfg.Language,
		Date:         time.Now().Format("2006-01-02"),
		TimeSpentMin: elapsedMinutes(cfg),
		Result:       tracker.ResultGaveUp,
		HintsUsed:    &ts.HintsUsed,
		TutorMode:    &ts.TutorMode,
		Model:        &ts.Model,
	}

	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		return tracker.Attempt{}, fmt.Errorf("session: open tracker: %w", err)
	}
	defer tr.Close()

	id, err := tr.LogAttempt(attempt)
	if err != nil {
		return tracker.Attempt{}, fmt.Errorf("session: log attempt: %w", err)
	}
	attempt.ID = id

	fmt.Fprintf(stdout, "reference solution revealed at reference/solution.* -- recorded as given up (attempt #%d)\n", id)
	return attempt, nil
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
