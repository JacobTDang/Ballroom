package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// mkdirWithModTime creates a directory at path and backdates its modtime
// to modTime -- how these tests simulate a stale (or fresh) leftover
// session dir without actually waiting real hours.
func mkdirWithModTime(t *testing.T, path string, modTime time.Time) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("chtimes %s: %v", path, err)
	}
}

func TestSweepStaleSessionDirs_RemovesStaleWorkspaceAndControlDirs(t *testing.T) {
	base := t.TempDir()
	now := time.Now()

	stale := []string{
		workspaceDirPrefix + "abc123",
		controlDirPrefix + "abc123",
	}
	for _, name := range stale {
		mkdirWithModTime(t, filepath.Join(base, name), now.Add(-48*time.Hour))
	}

	sweepStaleSessionDirs(base, 24*time.Hour, now)

	for _, name := range stale {
		if _, err := os.Stat(filepath.Join(base, name)); !os.IsNotExist(err) {
			t.Errorf("expected stale dir %s to be removed, stat err = %v", name, err)
		}
	}
}

func TestSweepStaleSessionDirs_SparesFreshDirs(t *testing.T) {
	base := t.TempDir()
	now := time.Now()

	fresh := []string{
		workspaceDirPrefix + "fresh",
		controlDirPrefix + "fresh",
	}
	for _, name := range fresh {
		mkdirWithModTime(t, filepath.Join(base, name), now.Add(-1*time.Minute))
	}

	sweepStaleSessionDirs(base, 24*time.Hour, now)

	for _, name := range fresh {
		if _, err := os.Stat(filepath.Join(base, name)); err != nil {
			t.Errorf("expected fresh dir %s to survive the sweep, stat err = %v", name, err)
		}
	}
}

func TestSweepStaleSessionDirs_SparesUnrelatedNames(t *testing.T) {
	base := t.TempDir()
	now := time.Now()

	unrelated := filepath.Join(base, "some-other-app-tmp-dir")
	mkdirWithModTime(t, unrelated, now.Add(-48*time.Hour))

	sweepStaleSessionDirs(base, 24*time.Hour, now)

	if _, err := os.Stat(unrelated); err != nil {
		t.Errorf("expected an unrelated dir name to survive the sweep untouched, stat err = %v", err)
	}
}

func TestSweepStaleSessionDirs_SparesUnrelatedFilesEvenIfStale(t *testing.T) {
	base := t.TempDir()
	now := time.Now()

	// A stray file (not a directory) that happens to start with the same
	// prefix must never be touched -- os.MkdirTemp only ever creates
	// directories for these two prefixes, so a file is never one of ours.
	stalePath := filepath.Join(base, workspaceDirPrefix+"not-a-dir")
	if err := os.WriteFile(stalePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}
	if err := os.Chtimes(stalePath, now.Add(-48*time.Hour), now.Add(-48*time.Hour)); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	sweepStaleSessionDirs(base, 24*time.Hour, now)

	if _, err := os.Stat(stalePath); err != nil {
		t.Errorf("expected a non-directory entry to survive the sweep, stat err = %v", err)
	}
}

func TestSweepStaleSessionDirs_MissingBaseDirIsANoOp(t *testing.T) {
	// Must never error/panic the app -- see SweepStaleSessionDirs' doc
	// comment. os.TempDir() always exists in practice, but nothing here
	// should assume that.
	sweepStaleSessionDirs(filepath.Join(t.TempDir(), "does-not-exist"), 24*time.Hour, time.Now())
}
