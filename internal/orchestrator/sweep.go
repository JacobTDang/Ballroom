package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// workspaceDirPrefix and controlDirPrefix are the os.MkdirTemp patterns
// PrepareWorkspace (workspace.go) and RunExercise (run.go) create a
// session's two disposable dirs with -- named here too so
// sweepStaleSessionDirs matches exactly what actually gets created,
// instead of carrying a second, driftable copy of the same string.
const (
	workspaceDirPrefix = "practice-workspace-"
	controlDirPrefix   = "practice-control-"
)

// staleSessionDirAge is how old an orphaned workspace/control dir (see
// sweepStaleSessionDirs) has to be before a startup sweep removes it --
// generous enough to never touch a session that's merely still running,
// while still reclaiming space from one left behind by, e.g., a closed
// terminal window predating issue #231's signal handling.
const staleSessionDirAge = 24 * time.Hour

// SweepStaleSessionDirs removes workspaceDirPrefix/controlDirPrefix
// directories under os.TempDir() whose modtime is older than
// staleSessionDirAge -- a safety net for sessions that never got the
// chance to clean up after themselves (a killed process, or a session
// run under a ballroom binary older than issue #231's signal handling).
// Meant to be called once, cheaply, at startup (see cmd/ballroom's
// main), and is entirely best-effort: it never touches anything outside
// the two prefixes above, and never reports an error back to its caller
// -- a sweep that fails to run isn't worth blocking the app over.
func SweepStaleSessionDirs() {
	sweepStaleSessionDirs(os.TempDir(), staleSessionDirAge, time.Now())
}

// sweepStaleSessionDirs is SweepStaleSessionDirs' testable core: dir,
// maxAge, and now are parameters instead of os.TempDir(),
// staleSessionDirAge, and time.Now() directly, so a test can point it at
// a t.TempDir() and control the cutoff without waiting real hours.
// Silently returns if dir can't be read -- same best-effort contract as
// the exported wrapper.
func sweepStaleSessionDirs(dir string, maxAge time.Duration, now time.Time) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	cutoff := now.Add(-maxAge)
	for _, entry := range entries {
		// os.MkdirTemp only ever creates directories for these two
		// prefixes -- a file that happens to share one is never ours to
		// remove, however stale it looks.
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, workspaceDirPrefix) && !strings.HasPrefix(name, controlDirPrefix) {
			continue
		}
		info, err := entry.Info()
		if err != nil || info.ModTime().After(cutoff) {
			continue
		}
		os.RemoveAll(filepath.Join(dir, name))
	}
}
