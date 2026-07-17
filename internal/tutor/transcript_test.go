package tutor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func TestAppendTranscriptTurn_WritesAllPathsCreatingDirs(t *testing.T) {
	dir := t.TempDir()
	direct := filepath.Join(dir, "transcript.md")
	nested := filepath.Join(dir, "transcripts", "two-sum-01.md") // parent doesn't exist yet

	if err := appendTranscriptTurn([]string{direct, nested}, "how do I start?", "Try a hash map."); err != nil {
		t.Fatalf("appendTranscriptTurn: %v", err)
	}
	if err := appendTranscriptTurn([]string{direct, nested}, "why a map?", "O(1) lookups."); err != nil {
		t.Fatalf("appendTranscriptTurn second turn: %v", err)
	}

	for _, path := range []string{direct, nested} {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		got := string(b)
		for _, want := range []string{"how do I start?", "Try a hash map.", "why a map?", "O(1) lookups."} {
			if !strings.Contains(got, want) {
				t.Errorf("%s missing %q:\n%s", path, want, got)
			}
		}
		if strings.Index(got, "how do I start?") > strings.Index(got, "why a map?") {
			t.Errorf("%s turns out of order:\n%s", path, got)
		}
	}
}

func TestAppendTranscriptTurn_ReturnsFirstError(t *testing.T) {
	bad := filepath.Join("/dev/null", "cannot", "exist.md")
	if err := appendTranscriptTurn([]string{bad}, "u", "r"); err == nil {
		t.Fatal("want an error for an unwritable path")
	}
}

// TestTutorModel_TurnsAppendToTranscriptPaths drives a real mocked turn
// and checks the exchange lands in every configured transcript file --
// and that a failed turn writes nothing (transcripts mirror history's
// clean-pairs-only semantics).
func TestTutorModel_TurnsAppendToTranscriptPaths(t *testing.T) {
	mock := newSequencedOllama(t, "use a set here")
	dir := t.TempDir()
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly
	cfg.TranscriptPaths = []string{
		filepath.Join(dir, "transcript.md"),
		filepath.Join(dir, "mirror", "two-sum-01.md"),
	}

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	m = submitAndRun(t, m, "what does this syntax mean?")

	for _, path := range cfg.TranscriptPaths {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		got := string(b)
		if !strings.Contains(got, "what does this syntax mean?") || !strings.Contains(got, "use a set here") {
			t.Errorf("%s missing the exchange:\n%s", path, got)
		}
	}
}

func TestTutorModel_TranscriptWriteFailureShowsOneDimNoteAndNeverCrashes(t *testing.T) {
	mock := newSequencedOllama(t, "first reply", "second reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly
	cfg.TranscriptPaths = []string{filepath.Join("/dev/null", "nope", "t.md")}

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	m = submitAndRun(t, m, "first question")
	m = submitAndRun(t, m, "second question")

	notes := 0
	for _, b := range m.displayBlocks {
		if b.kind == blockNote && strings.Contains(b.raw, "transcript") {
			notes++
		}
	}
	if notes != 1 {
		t.Errorf("got %d transcript-failure notes across two failing turns, want exactly 1 (warn once, then stay quiet)", notes)
	}
	// The turns themselves must still have completed normally.
	if len(m.history) != 5 { // system + 2 (user, assistant) pairs
		t.Errorf("history has %d messages, want 5 -- transcript failures must not break turns", len(m.history))
	}
}
