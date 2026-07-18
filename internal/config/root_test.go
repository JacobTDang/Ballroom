package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// checkoutDir creates a temp dir containing docker/Dockerfile, so it
// passes looksLikeCheckout — a stand-in for a real ballroom repo clone.
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

// withRootCache points rootCacheFile at a fresh, not-yet-existing path
// under a temp dir for the duration of the test, so tests never read or
// write the developer's real machine-wide cache.
func withRootCache(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "root-cache", "root")
	t.Setenv(rootCacheEnvVar, path)
	return path
}

func TestLooksLikeCheckout_TrueWhenDockerfileExists(t *testing.T) {
	if !looksLikeCheckout(checkoutDir(t)) {
		t.Error("expected a dir with docker/Dockerfile to look like a checkout")
	}
}

func TestLooksLikeCheckout_FalseWhenDockerfileMissing(t *testing.T) {
	if looksLikeCheckout(t.TempDir()) {
		t.Error("expected a plain empty dir not to look like a checkout")
	}
}

func TestCacheRootAndLoadCachedRoot_RoundTrips(t *testing.T) {
	path := withRootCache(t)
	cacheRoot(path, "/some/checkout/path")

	got, ok := loadCachedRoot(path)
	if !ok {
		t.Fatal("expected loadCachedRoot to find the cached value")
	}
	if got != "/some/checkout/path" {
		t.Errorf("got %q, want /some/checkout/path", got)
	}
}

func TestLoadCachedRoot_MissingFileReturnsFalse(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist")
	_, ok := loadCachedRoot(path)
	if ok {
		t.Error("expected loadCachedRoot to report false for a missing file")
	}
}

func TestCacheRoot_CreatesParentDirIfMissing(t *testing.T) {
	path := withRootCache(t)
	cacheRoot(path, "/x")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected cache file to exist at %s: %v", path, err)
	}
}

func TestResolveRoot_ReturnsCandidateWhenItLooksLikeCheckout(t *testing.T) {
	withRootCache(t)
	dir := checkoutDir(t)

	got, err := ResolveRoot(dir)
	if err != nil {
		t.Fatalf("ResolveRoot: %v", err)
	}
	if got != dir {
		t.Errorf("got %q, want %q", got, dir)
	}
}

func TestResolveRoot_ValidCandidateIsCachedForFutureFallback(t *testing.T) {
	cachePath := withRootCache(t)
	dir := checkoutDir(t)

	if _, err := ResolveRoot(dir); err != nil {
		t.Fatalf("ResolveRoot: %v", err)
	}

	cached, ok := loadCachedRoot(cachePath)
	if !ok || cached != dir {
		t.Errorf("expected %q to be cached, got %q (ok=%v)", dir, cached, ok)
	}
}

func TestResolveRoot_FallsBackToCachedRootWhenCandidateIsNotACheckout(t *testing.T) {
	cachePath := withRootCache(t)
	realCheckout := checkoutDir(t)
	cacheRoot(cachePath, realCheckout)

	notACheckout := t.TempDir()
	got, err := ResolveRoot(notACheckout)
	if err != nil {
		t.Fatalf("ResolveRoot: %v", err)
	}
	if got != realCheckout {
		t.Errorf("got %q, want cached checkout %q", got, realCheckout)
	}
}

func TestResolveRoot_ErrorsWhenNeitherCandidateNorCacheIsACheckout(t *testing.T) {
	withRootCache(t) // empty cache, nothing saved yet
	notACheckout := t.TempDir()

	_, err := ResolveRoot(notACheckout)
	if err == nil {
		t.Fatal("expected an error when neither the candidate nor the cache resolves to a real checkout")
	}
}

func TestResolveRoot_ErrorNamesTheProblemClearly(t *testing.T) {
	withRootCache(t)
	notACheckout := t.TempDir()

	_, err := ResolveRoot(notACheckout)
	if err == nil {
		t.Fatal("expected an error")
	}
	// The message should point at what's missing and how to fix it, not
	// just fail silently/cryptically -- someone hitting this is almost
	// certainly running an installed binary from outside the checkout
	// for the first time.
	want := "docker/Dockerfile"
	if got := err.Error(); !strings.Contains(got, want) {
		t.Errorf("error = %q, want it to mention %q", got, want)
	}
}

func TestResolveRoot_IgnoresStaleCachedRootThatNoLongerLooksLikeACheckout(t *testing.T) {
	cachePath := withRootCache(t)
	staleDir := t.TempDir() // once a checkout, now just an empty dir (e.g. repo moved/deleted)
	cacheRoot(cachePath, staleDir)

	_, err := ResolveRoot(t.TempDir())
	if err == nil {
		t.Fatal("expected an error since the cached root no longer looks like a real checkout")
	}
}
