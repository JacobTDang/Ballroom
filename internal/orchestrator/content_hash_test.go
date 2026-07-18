package orchestrator

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// writeTree creates files under root from a map of slash-separated
// relative path -> content, creating parent directories as needed. A
// stand-in for a real ballroom checkout's docker/cmd/internal trees,
// scoped down to just the files a given test cares about.
func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("MkdirAll %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("WriteFile %s: %v", full, err)
		}
	}
}

// baseTree is a minimal stand-in for the parts of a ballroom checkout
// contentHash actually cares about: something under each of docker/,
// cmd/, internal/, plus go.mod and go.sum.
func baseTree() map[string]string {
	return map[string]string{
		"docker/Dockerfile":    "FROM scratch\n",
		"docker/tmux.conf":     "set -g mouse on\n",
		"cmd/ballroom/main.go": "package main\n\nfunc main() {}\n",
		"internal/foo/foo.go":  "package foo\n",
		"go.mod":               "module example.com/x\n",
		"go.sum":               "",
	}
}

func TestContentHash_StableAcrossRepeatedCalls(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, baseTree())

	h1, err := contentHash(dir)
	if err != nil {
		t.Fatalf("contentHash: %v", err)
	}
	h2, err := contentHash(dir)
	if err != nil {
		t.Fatalf("contentHash: %v", err)
	}
	if h1 != h2 {
		t.Errorf("contentHash not stable across repeated calls: %q != %q", h1, h2)
	}
}

func TestContentHash_IndependentOfWalkOrder(t *testing.T) {
	files := baseTree()

	dirA := t.TempDir()
	writeTree(t, dirA, files)

	// Write the exact same files into a second, distinct directory in the
	// reverse order, so contentHash can't be relying on write order or
	// filesystem readdir order — it must sort explicitly.
	keys := make([]string, 0, len(files))
	for rel := range files {
		keys = append(keys, rel)
	}
	sort.Strings(keys)

	dirB := t.TempDir()
	for i := len(keys) - 1; i >= 0; i-- {
		rel := keys[i]
		writeTree(t, dirB, map[string]string{rel: files[rel]})
	}

	hA, err := contentHash(dirA)
	if err != nil {
		t.Fatalf("contentHash(dirA): %v", err)
	}
	hB, err := contentHash(dirB)
	if err != nil {
		t.Fatalf("contentHash(dirB): %v", err)
	}
	if hA != hB {
		t.Errorf("contentHash depends on write/walk order: %q != %q", hA, hB)
	}
}

func TestContentHash_ChangesWhenNonTestFileChanges(t *testing.T) {
	dirA := t.TempDir()
	writeTree(t, dirA, baseTree())

	changed := baseTree()
	changed["internal/foo/foo.go"] = "package foo\n\nfunc Bar() {}\n"
	dirB := t.TempDir()
	writeTree(t, dirB, changed)

	hA, err := contentHash(dirA)
	if err != nil {
		t.Fatalf("contentHash(dirA): %v", err)
	}
	hB, err := contentHash(dirB)
	if err != nil {
		t.Fatalf("contentHash(dirB): %v", err)
	}
	if hA == hB {
		t.Errorf("expected hash to change when a non-test file's content changes, got the same hash %q for both", hA)
	}
}

func TestContentHash_UnchangedWhenOnlyTestFileChanges(t *testing.T) {
	base := baseTree()
	base["internal/foo/foo_test.go"] = "package foo\n\nfunc TestFoo(t *testing.T) {}\n"
	dirA := t.TempDir()
	writeTree(t, dirA, base)

	changed := baseTree()
	changed["internal/foo/foo_test.go"] = "package foo\n\nfunc TestFoo(t *testing.T) { /* totally different body */ }\n"
	dirB := t.TempDir()
	writeTree(t, dirB, changed)

	hA, err := contentHash(dirA)
	if err != nil {
		t.Fatalf("contentHash(dirA): %v", err)
	}
	hB, err := contentHash(dirB)
	if err != nil {
		t.Fatalf("contentHash(dirB): %v", err)
	}
	if hA != hB {
		t.Errorf("expected _test.go changes to be excluded from the hash, got different hashes: %q vs %q", hA, hB)
	}
}

func TestContentHash_UnchangedWhenOnlyTestdataFileChanges(t *testing.T) {
	base := baseTree()
	base["internal/foo/testdata/fixture.txt"] = "v1\n"
	dirA := t.TempDir()
	writeTree(t, dirA, base)

	changed := baseTree()
	changed["internal/foo/testdata/fixture.txt"] = "v2, completely different contents\n"
	dirB := t.TempDir()
	writeTree(t, dirB, changed)

	hA, err := contentHash(dirA)
	if err != nil {
		t.Fatalf("contentHash(dirA): %v", err)
	}
	hB, err := contentHash(dirB)
	if err != nil {
		t.Fatalf("contentHash(dirB): %v", err)
	}
	if hA != hB {
		t.Errorf("expected testdata/ changes to be excluded from the hash, got different hashes: %q vs %q", hA, hB)
	}
}

func TestContentHash_MissingDirectoryDoesNotError(t *testing.T) {
	dir := t.TempDir() // no docker/, cmd/, internal/, go.mod, or go.sum at all

	got, err := contentHash(dir)
	if err != nil {
		t.Fatalf("contentHash: %v", err)
	}
	if got == "" {
		t.Error("expected a non-empty hash even when every content root is missing")
	}
}

func TestContentHash_PartiallyMissingDirectoryDoesNotError(t *testing.T) {
	dir := t.TempDir()
	// Only go.mod present -- docker/, cmd/, internal/, go.sum all absent.
	writeTree(t, dir, map[string]string{"go.mod": "module example.com/x\n"})

	if _, err := contentHash(dir); err != nil {
		t.Fatalf("contentHash: %v", err)
	}
}
