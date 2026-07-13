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
		"ballroom", "home", "practice <id>", "sandbox", "submit", "tutor", "return", "help", "config",
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

// fakeCheckToolCallingCLI substitutes checkToolCallingFn in tests so
// setModelCmd's warning path never makes a real network call — same
// indirection pattern internal/tui/app_test.go uses for the identical
// reason.
func fakeCheckToolCallingCLI(supported bool, err error) func() {
	orig := checkToolCallingFn
	checkToolCallingFn = func(string, string, string) (bool, error) { return supported, err }
	return func() { checkToolCallingFn = orig }
}

func TestConfigCmd_RequiresASubcommand(t *testing.T) {
	err := configCmd(nil)
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("configCmd(nil) error = %v, want a usage error", err)
	}
}

func TestConfigCmd_UnknownSubcommandReturnsError(t *testing.T) {
	err := configCmd([]string{"bogus"})
	if err == nil || !strings.Contains(err.Error(), "bogus") {
		t.Errorf("configCmd([\"bogus\"]) error = %v, want it to mention the unknown subcommand", err)
	}
}

func TestConfigCmd_SetModelRequiresATag(t *testing.T) {
	err := configCmd([]string{"set-model"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("configCmd([\"set-model\"]) error = %v, want a usage error", err)
	}
}

func TestConfigCmd_SetKeyRequiresAKey(t *testing.T) {
	err := configCmd([]string{"set-key"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("configCmd([\"set-key\"]) error = %v, want a usage error", err)
	}
}

func TestConfigCmd_SetOrchestratorModelRequiresATag(t *testing.T) {
	err := configCmd([]string{"set-orchestrator-model"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("configCmd([\"set-orchestrator-model\"]) error = %v, want a usage error", err)
	}
}

func TestConfigCmd_DispatchesSetOrchestratorModel(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := configCmd([]string{"set-orchestrator-model", "nemotron"}); err != nil {
		t.Fatalf("configCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OrchestratorModel != "nemotron" {
		t.Errorf("OrchestratorModel = %q, want %q", cfg.OrchestratorModel, "nemotron")
	}
}

func TestSetModelCmd_PersistsToSettings(t *testing.T) {
	defer fakeCheckToolCallingCLI(true, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setModelCmd("llama3.1:8b"); err != nil {
		t.Fatalf("setModelCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.TutorModel != "llama3.1:8b" {
		t.Errorf("TutorModel = %q, want %q", cfg.TutorModel, "llama3.1:8b")
	}
}

func TestSetModelCmd_PreservesExistingOpenRouterAPIKey(t *testing.T) {
	defer fakeCheckToolCallingCLI(true, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setKeyCmd("sk-preserve-me"); err != nil {
		t.Fatalf("setKeyCmd: %v", err)
	}
	if err := setModelCmd("llama3.1:8b"); err != nil {
		t.Fatalf("setModelCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OpenRouterAPIKey != "sk-preserve-me" {
		t.Errorf("OpenRouterAPIKey = %q, want it preserved across a set-model call", cfg.OpenRouterAPIKey)
	}
}

func TestSetModelCmd_OpenRouterModelWithNoKeyWarnsWithoutCallingCheck(t *testing.T) {
	checkCalled := false
	orig := checkToolCallingFn
	checkToolCallingFn = func(string, string, string) (bool, error) {
		checkCalled = true
		return true, nil
	}
	defer func() { checkToolCallingFn = orig }()
	t.Setenv("PRACTICE_ROOT", t.TempDir())
	t.Setenv("OPENROUTER_API_KEY", "")

	out := captureStdout(t, func() {
		if err := setModelCmd("openrouter:anthropic/claude-3.5-sonnet"); err != nil {
			t.Fatalf("setModelCmd: %v", err)
		}
	})

	if checkCalled {
		t.Error("expected the tool-calling check to be skipped when no OpenRouter key is configured")
	}
	if !strings.Contains(out, "no OpenRouter API key configured") {
		t.Errorf("output %q missing the missing-key warning", out)
	}
}

func TestSetModelCmd_WarnsWhenToolCallingUnsupported(t *testing.T) {
	defer fakeCheckToolCallingCLI(false, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	out := captureStdout(t, func() {
		if err := setModelCmd("some-model"); err != nil {
			t.Fatalf("setModelCmd: %v", err)
		}
	})

	if !strings.Contains(out, "may not support real tool calling") {
		t.Errorf("output %q missing the unsupported-tool-calling warning", out)
	}
}

func TestSetKeyCmd_PersistsToSettings(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())
	t.Setenv("OPENROUTER_API_KEY", "")

	if err := setKeyCmd("sk-abc123"); err != nil {
		t.Fatalf("setKeyCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OpenRouterAPIKey != "sk-abc123" {
		t.Errorf("OpenRouterAPIKey = %q, want %q", cfg.OpenRouterAPIKey, "sk-abc123")
	}
}

func TestSetModelCmd_PreservesExistingOrchestratorModel(t *testing.T) {
	defer fakeCheckToolCallingCLI(true, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setOrchestratorModelCmd("nemotron"); err != nil {
		t.Fatalf("setOrchestratorModelCmd: %v", err)
	}
	if err := setModelCmd("llama3.1:8b"); err != nil {
		t.Fatalf("setModelCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OrchestratorModel != "nemotron" {
		t.Errorf("OrchestratorModel = %q, want it preserved across a set-model call", cfg.OrchestratorModel)
	}
}

func TestSetKeyCmd_PreservesExistingOrchestratorModel(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())
	t.Setenv("OPENROUTER_API_KEY", "")

	if err := setOrchestratorModelCmd("nemotron"); err != nil {
		t.Fatalf("setOrchestratorModelCmd: %v", err)
	}
	if err := setKeyCmd("sk-abc123"); err != nil {
		t.Fatalf("setKeyCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OrchestratorModel != "nemotron" {
		t.Errorf("OrchestratorModel = %q, want it preserved across a set-key call", cfg.OrchestratorModel)
	}
}

func TestSetKeyCmd_PreservesExistingTutorModel(t *testing.T) {
	defer fakeCheckToolCallingCLI(true, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setModelCmd("llama3.1:8b"); err != nil {
		t.Fatalf("setModelCmd: %v", err)
	}
	if err := setKeyCmd("sk-abc123"); err != nil {
		t.Fatalf("setKeyCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.TutorModel != "llama3.1:8b" {
		t.Errorf("TutorModel = %q, want it preserved across a set-key call", cfg.TutorModel)
	}
}

func TestSetOrchestratorModelCmd_PersistsToSettings(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setOrchestratorModelCmd("nemotron"); err != nil {
		t.Fatalf("setOrchestratorModelCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OrchestratorModel != "nemotron" {
		t.Errorf("OrchestratorModel = %q, want %q", cfg.OrchestratorModel, "nemotron")
	}
}

func TestSetOrchestratorModelCmd_NonePreservedFieldsButClearsOrchestratorModel(t *testing.T) {
	defer fakeCheckToolCallingCLI(true, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setModelCmd("llama3.1:8b"); err != nil {
		t.Fatalf("setModelCmd: %v", err)
	}
	if err := setOrchestratorModelCmd("nemotron"); err != nil {
		t.Fatalf("setOrchestratorModelCmd: %v", err)
	}
	if err := setOrchestratorModelCmd("none"); err != nil {
		t.Fatalf("setOrchestratorModelCmd(none): %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.OrchestratorModel != "" {
		t.Errorf("OrchestratorModel = %q, want empty after set-orchestrator-model none", cfg.OrchestratorModel)
	}
	if cfg.TutorModel != "llama3.1:8b" {
		t.Errorf("TutorModel = %q, want it preserved across set-orchestrator-model none", cfg.TutorModel)
	}
}

func TestSetOrchestratorModelCmd_PreservesExistingTutorModelAndKey(t *testing.T) {
	defer fakeCheckToolCallingCLI(true, nil)()
	t.Setenv("PRACTICE_ROOT", t.TempDir())

	if err := setModelCmd("llama3.1:8b"); err != nil {
		t.Fatalf("setModelCmd: %v", err)
	}
	if err := setKeyCmd("sk-preserve-me"); err != nil {
		t.Fatalf("setKeyCmd: %v", err)
	}
	if err := setOrchestratorModelCmd("nemotron"); err != nil {
		t.Fatalf("setOrchestratorModelCmd: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if cfg.TutorModel != "llama3.1:8b" {
		t.Errorf("TutorModel = %q, want it preserved across a set-orchestrator-model call", cfg.TutorModel)
	}
	if cfg.OpenRouterAPIKey != "sk-preserve-me" {
		t.Errorf("OpenRouterAPIKey = %q, want it preserved across a set-orchestrator-model call", cfg.OpenRouterAPIKey)
	}
}

// captureStdout runs fn with os.Stdout redirected to a pipe and returns
// what was written — setModelCmd/setKeyCmd print their confirmation/
// warning messages directly to stdout rather than returning them.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	fn()
	os.Stdout = orig
	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return buf.String()
}
