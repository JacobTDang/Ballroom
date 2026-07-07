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
