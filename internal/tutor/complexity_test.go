package tutor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckComplexity_SendsSolutionAndClaimReturnsVerdict(t *testing.T) {
	mock := newSequencedOllama(t, "AGREE\nO(n) time, O(n) space -- one pass over nums with a map.")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "solution.py"), []byte("def f(nums):\n    seen = {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := Config{OllamaHost: mock.URL, Model: "test-model", WorkDir: dir, MaxContextBytes: 8000}

	got, err := CheckComplexity(context.Background(), cfg, "O(n) time O(1) space")
	if err != nil {
		t.Fatalf("CheckComplexity: %v", err)
	}
	if !strings.HasPrefix(got, "AGREE") {
		t.Errorf("verdict = %q, want the model's reply", got)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("made %d model calls, want exactly 1", len(reqs))
	}
	var sawSolution, sawClaim bool
	for _, m := range reqs[0].Messages {
		if strings.Contains(m.Content, "seen = {}") {
			sawSolution = true
		}
		if strings.Contains(m.Content, "O(n) time O(1) space") {
			sawClaim = true
		}
	}
	if !sawSolution || !sawClaim {
		t.Errorf("request missing solution (%v) or claim (%v)", sawSolution, sawClaim)
	}
}

func TestCheckComplexity_NoSolutionFileErrors(t *testing.T) {
	cfg := Config{OllamaHost: "http://unused", Model: "m", WorkDir: t.TempDir(), MaxContextBytes: 8000}
	if _, err := CheckComplexity(context.Background(), cfg, "O(1)"); err == nil {
		t.Fatal("want an error with no solution file")
	}
}
