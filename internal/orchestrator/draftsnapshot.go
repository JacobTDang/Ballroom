package orchestrator

import (
	"context"
	"time"

	"github.com/JacobTDang/Ballroom/internal/draft"
)

// SnapshotLoop periodically snapshots workspaceDir's solution.* files
// into exerciseID's draft directory under dataDir (see
// internal/draft.Snapshot) every interval, until ctx is cancelled.
//
// Mirrors WaitAndReveal's goroutine convention: started from
// RunExercise as `go func() { ... }()`, cancelled when the container
// exits. Unlike WaitAndReveal, cancellation is the only way this loop
// ever ends -- there's no equivalent of "submit happened, stop early"
// -- so returning nil on ctx.Done() is the normal, expected outcome,
// not something callers need to filter out. Snapshot errors along the
// way are non-fatal (a transient read/write hiccup on a cheap poll
// shouldn't tear down the session) but are remembered and returned
// once the loop stops, so a caller can still surface them.
func SnapshotLoop(ctx context.Context, dataDir, exerciseID, workspaceDir string, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return lastErr
		case <-ticker.C:
			if _, err := draft.Snapshot(dataDir, exerciseID, workspaceDir); err != nil {
				lastErr = err
			}
		}
	}
}
