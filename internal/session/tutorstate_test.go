package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTutorStateFixture(t *testing.T, workDir string, s tutorState) {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal tutor state fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workDir, tutorStateFile), data, 0o644); err != nil {
		t.Fatalf("write tutor state fixture: %v", err)
	}
}

func TestReadTutorState_ReadsRealValues(t *testing.T) {
	dir := t.TempDir()
	writeTutorStateFixture(t, dir, tutorState{HintsUsed: 4, TutorMode: "hints-first", Model: "llama3.1:8b"})

	got := readTutorState(dir)
	if got.HintsUsed != 4 {
		t.Errorf("HintsUsed = %d, want 4", got.HintsUsed)
	}
	if got.TutorMode != "hints-first" {
		t.Errorf("TutorMode = %q, want %q", got.TutorMode, "hints-first")
	}
	if got.Model != "llama3.1:8b" {
		t.Errorf("Model = %q, want %q", got.Model, "llama3.1:8b")
	}
}

// TestReadTutorState_MissingFileDegradesToZeroValue is the explicit
// contract: a missing dotfile (sandbox mode, or a submit/reference that
// happens before the tutor pane has written anything) must never error
// -- it degrades silently to the zero value.
func TestReadTutorState_MissingFileDegradesToZeroValue(t *testing.T) {
	dir := t.TempDir()

	got := readTutorState(dir)
	if got.HintsUsed != 0 || got.TutorMode != "" || got.Model != "" {
		t.Errorf("readTutorState on a missing file = %+v, want the zero value", got)
	}
}

// TestReadTutorState_MalformedFileDegradesToZeroValue matches
// readLastTestResult's OPPOSITE contract on purpose: unlike a bad test
// result (a real error, since that IS the attempt's own outcome), a
// broken tutor-state dotfile is just metadata about an attempt and must
// never fail it.
func TestReadTutorState_MalformedFileDegradesToZeroValue(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, tutorStateFile), []byte("not json"), 0o644); err != nil {
		t.Fatalf("write malformed fixture: %v", err)
	}

	got := readTutorState(dir)
	if got.HintsUsed != 0 || got.TutorMode != "" || got.Model != "" {
		t.Errorf("readTutorState on a malformed file = %+v, want the zero value", got)
	}
}
