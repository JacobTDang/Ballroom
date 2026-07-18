package orchestrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	submitRequestFile = "submit.request"
	testsReadyFile    = "tests.ready"

	referenceRequestFile = "reference.request"
	referenceReadyFile   = "reference.ready"
)

// WaitAndReveal blocks until controlDir/submit.request appears (written by
// the in-container `practice submit` command), then copies testsSrc into
// workspace and writes controlDir/tests.ready to signal completion.
//
// Hidden tests are never mounted or present in workspace before this runs —
// this is the actual hiding guarantee, not an honor-system convention.
func WaitAndReveal(ctx context.Context, controlDir, testsSrc, workspace string, pollInterval time.Duration) error {
	requestPath := filepath.Join(controlDir, submitRequestFile)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		if _, err := os.Stat(requestPath); err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("orchestrator: wait for submit request: %w", ctx.Err())
		case <-ticker.C:
		}
	}

	if err := copyTree(testsSrc, workspace); err != nil {
		return fmt.Errorf("orchestrator: reveal tests: %w", err)
	}

	readyPath := filepath.Join(controlDir, testsReadyFile)
	if err := os.WriteFile(readyPath, nil, 0o644); err != nil {
		return fmt.Errorf("orchestrator: write ready marker: %w", err)
	}
	return nil
}

// WaitAndRevealReference blocks until controlDir/reference.request
// appears (written by the in-container `ballroom reference` command, or
// by a passing `ballroom submit` -- see internal/session), then copies
// referenceSrc (the exercise's .reference/ directory) plus testsSrc (the
// same hidden tests WaitAndReveal reveals on submit) into
// workspace/reference/ -- a SUBDIRECTORY, never the workspace root, so
// the reference solution can never land on top of or otherwise clobber
// the user's own solution file -- and writes controlDir/reference.ready
// to signal completion. The hidden tests ride along so the reference
// reads next to what it must satisfy, not just the bare answer.
//
// Runs entirely independently of WaitAndReveal: reference and submit are
// unrelated triggers that can arrive in either order, or not at all, so
// a real session spawns both watchers side by side (see RunExercise) and
// this never blocks on or interferes with submit handling.
//
// referenceSrc not existing (a malformed exercise, or a design exercise,
// which never has one) is a real, immediate error rather than silently
// retrying forever -- the copy can never succeed no matter how many more
// times reference.request is polled, so failing fast here is what lets
// the in-container caller's own bounded wait (see session.Reference)
// report a clear failure instead of just timing out in silence.
func WaitAndRevealReference(ctx context.Context, controlDir, referenceSrc, testsSrc, workspace string, pollInterval time.Duration) error {
	requestPath := filepath.Join(controlDir, referenceRequestFile)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		if _, err := os.Stat(requestPath); err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("orchestrator: wait for reference request: %w", ctx.Err())
		case <-ticker.C:
		}
	}

	if info, err := os.Stat(referenceSrc); err != nil || !info.IsDir() {
		return fmt.Errorf("orchestrator: reveal reference: no reference solution available for this exercise (missing %s)", referenceSrc)
	}

	dst := filepath.Join(workspace, "reference")
	if err := copyTree(referenceSrc, dst); err != nil {
		return fmt.Errorf("orchestrator: reveal reference: %w", err)
	}
	if err := copyTree(testsSrc, dst); err != nil {
		return fmt.Errorf("orchestrator: reveal reference tests: %w", err)
	}

	readyPath := filepath.Join(controlDir, referenceReadyFile)
	if err := os.WriteFile(readyPath, nil, 0o644); err != nil {
		return fmt.Errorf("orchestrator: write reference ready marker: %w", err)
	}
	return nil
}

// copyTree recursively copies the contents of src into dst, creating dst
// (and any subdirectories) as needed. dst is not required to exist yet.
func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
