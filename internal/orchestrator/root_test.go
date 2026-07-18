package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

// rootCacheEnvVar mirrors internal/config's unexported const of the same
// name (see that package's root.go) — this package no longer owns the
// cached-root logic itself (issue #255: dockerBuildRoot now just
// delegates to config.ResolveRoot), but its tests still need to point
// the cache at an isolated temp path so they never read or write the
// developer's real machine-wide cache.
const rootCacheEnvVar = "PRACTICE_ROOT_CACHE_FILE"

// checkoutDir creates a temp dir containing docker/Dockerfile, so it
// passes config.ResolveRoot's checkout check — a stand-in for a real
// ballroom repo clone.
func checkoutDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "docker"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "docker", "Dockerfile"), []byte("FROM scratch\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return dir
}

// withRootCache points the cache at a fresh, not-yet-existing path under
// a temp dir for the duration of the test, so tests never read or write
// the developer's real machine-wide cache.
func withRootCache(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "root-cache", "root")
	t.Setenv(rootCacheEnvVar, path)
	return path
}

// The exhaustive looks-like-a-checkout / cache-fallback / error-message
// behavior is now tested once, in internal/config (root_test.go), since
// that's where the logic actually lives. These two just confirm
// dockerBuildRoot still delegates to it correctly.

func TestDockerBuildRoot_ReturnsCfgRootWhenItLooksLikeCheckout(t *testing.T) {
	withRootCache(t)
	dir := checkoutDir(t)

	got, err := dockerBuildRoot(dir)
	if err != nil {
		t.Fatalf("dockerBuildRoot: %v", err)
	}
	if got != dir {
		t.Errorf("got %q, want %q", got, dir)
	}
}

func TestDockerBuildRoot_ErrorsWhenNeitherCfgRootNorCacheIsACheckout(t *testing.T) {
	withRootCache(t) // empty cache, nothing saved yet
	notACheckout := t.TempDir()

	_, err := dockerBuildRoot(notACheckout)
	if err == nil {
		t.Fatal("expected an error when neither cfgRoot nor the cache resolves to a real checkout")
	}
}
