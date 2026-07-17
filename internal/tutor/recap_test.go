package tutor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSessionRecap_SendsTranscriptResultAndOutput(t *testing.T) {
	mock := newSequencedOllama(t, "You solved two-sum after one nudge about complements; the hidden tests passed first try.")
	dir := t.TempDir()
	transcript := "## you\n\nhow do I start?\n\n## tutor\n\nThink about complements.\n"
	if err := os.WriteFile(filepath.Join(dir, "transcript.md"), []byte(transcript), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := Config{OllamaHost: mock.URL, Model: "test-model", WorkDir: dir, MaxContextBytes: 8000}

	got, err := SessionRecap(context.Background(), cfg, "pass", "=== RUN TestTwoSum\n--- PASS")
	if err != nil {
		t.Fatalf("SessionRecap: %v", err)
	}
	if !strings.Contains(got, "two-sum") {
		t.Errorf("recap = %q, want the model reply", got)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("made %d calls, want 1", len(reqs))
	}
	joined := ""
	for _, m := range reqs[0].Messages {
		joined += m.Content + "\n"
	}
	for _, want := range []string{"Think about complements", "pass", "--- PASS"} {
		if !strings.Contains(joined, want) {
			t.Errorf("request missing %q", want)
		}
	}
}

func TestSessionRecap_NoTranscriptStillRecaps(t *testing.T) {
	mock := newSequencedOllama(t, "You passed without asking the tutor anything.")
	cfg := Config{OllamaHost: mock.URL, Model: "test-model", WorkDir: t.TempDir(), MaxContextBytes: 8000}

	if _, err := SessionRecap(context.Background(), cfg, "pass", "ok"); err != nil {
		t.Fatalf("SessionRecap without a transcript: %v", err)
	}
	joined := ""
	for _, m := range mock.allRequests()[0].Messages {
		joined += m.Content
	}
	if !strings.Contains(joined, "didn't talk to the tutor") {
		t.Error("request should carry the no-conversation note instead of an empty transcript")
	}
}
