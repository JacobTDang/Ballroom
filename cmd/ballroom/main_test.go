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
		"ballroom", "home", "practice <id>", "sandbox", "submit", "return", "help",
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

// clearSessionEnv unsets the session-scoped env vars orchestrator.RunExercise
// sets via `docker run -e`, so tests start from a known "on the host"
// baseline regardless of what's in the ambient environment.
func clearSessionEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PRACTICE_WORKSPACE_DIR", "")
	t.Setenv("PRACTICE_CONTROL_DIR", "")
	t.Setenv("PRACTICE_STARTED_AT", "")
}

func TestIsSessionContext_TrueWhenAllSessionVarsSet(t *testing.T) {
	clearSessionEnv(t)
	t.Setenv("PRACTICE_WORKSPACE_DIR", "/workspace")
	t.Setenv("PRACTICE_CONTROL_DIR", "/control")
	t.Setenv("PRACTICE_STARTED_AT", "2026-07-08T00:00:00Z")

	if !isSessionContext() {
		t.Error("isSessionContext() = false, want true when all session env vars are set")
	}
}

func TestIsSessionContext_FalseOnHost(t *testing.T) {
	clearSessionEnv(t)

	if isSessionContext() {
		t.Error("isSessionContext() = true, want false when no session env vars are set")
	}
}

func TestIsSessionContext_FalseWhenOnlyPartiallySet(t *testing.T) {
	clearSessionEnv(t)
	t.Setenv("PRACTICE_WORKSPACE_DIR", "/workspace")

	if isSessionContext() {
		t.Error("isSessionContext() = true, want false when only some session env vars are set")
	}
}

func TestReturnCmd_ErrorsOutsideSession(t *testing.T) {
	clearSessionEnv(t)

	err := returnCmd()
	if err == nil || !strings.Contains(err.Error(), "not running inside an active practice session") {
		t.Errorf("returnCmd() outside a session = %v, want an error about not being in a session", err)
	}
}
