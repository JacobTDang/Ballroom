package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/config"
)

// captureUsage runs printUsage through an os.Pipe and returns what it wrote.
func captureUsage(t *testing.T) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	printUsage(w)
	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return buf.String()
}

func TestPrintUsage_MentionsEverySubcommand(t *testing.T) {
	out := captureUsage(t)
	for _, want := range []string{
		"ballroom", "home", "practice <id>", "sandbox", "submit", "help",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("usage output missing %q:\n%s", want, out)
		}
	}
}

func TestRunExercise_UnknownIDReturnsClearError(t *testing.T) {
	cfg := config.Config{ExercisesDir: t.TempDir()}

	err := runExercise(cfg, "does-not-exist")
	if err == nil {
		t.Fatal("expected an error for an unknown exercise id, got nil")
	}
	if !strings.Contains(err.Error(), "unknown exercise") {
		t.Errorf("error = %v, want it to mention \"unknown exercise\"", err)
	}
}

func TestPracticeCmd_RequiresAnID(t *testing.T) {
	err := practiceCmd(nil)
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("practiceCmd(nil) error = %v, want a usage error", err)
	}
}

func TestPracticeCmd_UnknownIDReturnsClearError(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	err := practiceCmd([]string{"does-not-exist"})
	if err == nil || !strings.Contains(err.Error(), "unknown exercise") {
		t.Errorf("error = %v, want it to mention \"unknown exercise\"", err)
	}
}
