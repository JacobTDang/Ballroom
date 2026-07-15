package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

const testPollInterval = 5 * time.Millisecond

// simulateHostReveal mimics what orchestrator.WaitAndReveal does on the host
// side: waits for submit.request, then writes tests.ready. Kept independent
// of the orchestrator package so session's tests don't couple to it.
func simulateHostReveal(t *testing.T, controlDir string) {
	t.Helper()
	requestPath := filepath.Join(controlDir, "submit.request")
	go func() {
		for {
			if _, err := os.Stat(requestPath); err == nil {
				os.WriteFile(filepath.Join(controlDir, "tests.ready"), nil, 0o644)
				return
			}
			time.Sleep(testPollInterval)
		}
	}()
}

func baseConfig(t *testing.T) Config {
	t.Helper()
	dir := t.TempDir()
	return Config{
		ControlDir:    filepath.Join(dir, "control"),
		WorkspaceDir:  filepath.Join(dir, "workspace"),
		ExerciseID:    "two-pointers-01",
		Category:      "pattern",
		Language:      "go",
		StartedAt:     time.Now(),
		DBPath:        filepath.Join(dir, "tracker.db"),
		PollInterval:  testPollInterval,
		RevealTimeout: time.Second,
	}
}

func mkdirs(t *testing.T, cfg Config) {
	t.Helper()
	if err := os.MkdirAll(cfg.ControlDir, 0o755); err != nil {
		t.Fatalf("mkdir control: %v", err)
	}
	if err := os.MkdirAll(cfg.WorkspaceDir, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
}

// designConfig is baseConfig shaped as a design-kind session: no test
// command (validated empty at authoring time), session style in the
// language slot.
func designConfig(t *testing.T) Config {
	cfg := baseConfig(t)
	cfg.Kind = "design"
	cfg.TestCommand = ""
	cfg.ExerciseID = "url-shortener-01-interviewer"
	cfg.Category = "system-design"
	cfg.Language = "interviewer"
	return cfg
}

func TestSubmit_DesignUsesGraderVerdict(t *testing.T) {
	cfg := designConfig(t)
	cfg.Grade = func() (string, string, error) {
		return tracker.ResultFail, "Estimates: missing. Sharding: adequate.", nil
	}
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	var out bytes.Buffer
	attempt, err := Submit(cfg, strings.NewReader("tough one\n"), &out)
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultFail {
		t.Errorf("Result = %q, want the grader's fail verdict", attempt.Result)
	}
	if attempt.Notes != "tough one" {
		t.Errorf("Notes = %q, want the user's own note", attempt.Notes)
	}
	if !strings.Contains(out.String(), "Estimates: missing") {
		t.Errorf("output %q should show the grader's summary", out.String())
	}
	if attempt.GradeSummary != "Estimates: missing. Sharding: adequate." {
		t.Errorf("GradeSummary = %q, want the grader's summary persisted on the attempt", attempt.GradeSummary)
	}
	if strings.Contains(out.String(), "Self-assessment") {
		t.Errorf("output %q ran the self-assessment prompt despite a successful grade", out.String())
	}

	data, err := os.ReadFile(filepath.Join(cfg.WorkspaceDir, lastTestResultFile))
	if err != nil {
		t.Fatalf("read last-test-result file: %v", err)
	}
	var got lastTestResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Result != tracker.ResultFail {
		t.Errorf("file Result = %q, want fail", got.Result)
	}
	if !strings.Contains(got.Output, "Estimates: missing") {
		t.Errorf("file Output = %q, want the grading summary so read_test_output can show it", got.Output)
	}
}

func TestSubmit_DesignGraderErrorFallsBackToSelfAssessment(t *testing.T) {
	cfg := designConfig(t)
	cfg.Grade = func() (string, string, error) {
		return "", "", fmt.Errorf("empty choices from provider")
	}
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	var out bytes.Buffer
	attempt, err := Submit(cfg, strings.NewReader("p\n\n"), &out)
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultPass {
		t.Errorf("Result = %q, want the self-assessed pass after fallback", attempt.Result)
	}
	if attempt.GradeSummary != "" {
		t.Errorf("GradeSummary = %q, want empty for a self-assessed attempt", attempt.GradeSummary)
	}
	if !strings.Contains(out.String(), "empty choices from provider") {
		t.Errorf("output %q should surface the grading failure, not swallow it", out.String())
	}
	if !strings.Contains(out.String(), "Self-assessment") {
		t.Errorf("output %q should have fallen back to the self-assessment prompt", out.String())
	}
}

func TestSubmit_DesignSelfAssessedPass(t *testing.T) {
	cfg := designConfig(t)
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	var out bytes.Buffer
	attempt, err := Submit(cfg, strings.NewReader("p\nnailed the estimates\n"), &out)
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultPass {
		t.Errorf("Result = %q, want pass from the self-assessment", attempt.Result)
	}
	if attempt.Notes != "nailed the estimates" {
		t.Errorf("Notes = %q, want the notes line AFTER the self-assessment line", attempt.Notes)
	}
	if !strings.Contains(out.String(), "rubric") {
		t.Errorf("output %q should tell the user the rubric was revealed", out.String())
	}
}

func TestSubmit_DesignSelfAssessedFail(t *testing.T) {
	cfg := designConfig(t)
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("f\n\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultFail {
		t.Errorf("Result = %q, want fail from the self-assessment", attempt.Result)
	}
}

func TestSubmit_DesignRejectsNonAnswerUntilExplicit(t *testing.T) {
	// No default: a bare Enter or noise must re-prompt, never silently
	// record a result -- pass/fail feeds the "solved" stats.
	cfg := designConfig(t)
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	var out bytes.Buffer
	attempt, err := Submit(cfg, strings.NewReader("\nmaybe\npass\n\n"), &out)
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultPass {
		t.Errorf("Result = %q, want pass once an explicit answer finally arrives", attempt.Result)
	}
	if n := strings.Count(out.String(), "pass or fail"); n < 2 {
		t.Errorf("expected at least 2 re-prompts for the 2 non-answers, saw %d in %q", n, out.String())
	}
}

func TestSubmit_DesignDoesNotRunAnyCommandAndWritesCoherentResultFile(t *testing.T) {
	cfg := designConfig(t)
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	if _, err := Submit(cfg, strings.NewReader("p\n\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("Submit: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cfg.WorkspaceDir, lastTestResultFile))
	if err != nil {
		t.Fatalf("read last-test-result file: %v", err)
	}
	var got lastTestResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Result != tracker.ResultPass {
		t.Errorf("file Result = %q, want pass", got.Result)
	}
	if !strings.Contains(got.Output, "self-assessed") {
		t.Errorf("file Output = %q, want it to say the result was self-assessed so the tutor's read_test_output stays coherent", got.Output)
	}
	if got.TestCommand != "" {
		t.Errorf("file TestCommand = %q, want empty -- nothing was run", got.TestCommand)
	}
}

func TestSubmit_PassResult(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultPass {
		t.Errorf("Result = %q, want %q", attempt.Result, tracker.ResultPass)
	}
}

func TestSubmit_FailResult(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "false"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Result != tracker.ResultFail {
		t.Errorf("Result = %q, want %q", attempt.Result, tracker.ResultFail)
	}
}

func TestSubmit_ComputesTimeSpent(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	cfg.StartedAt = time.Now().Add(-10 * time.Minute)
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.TimeSpentMin < 9.9 || attempt.TimeSpentMin > 10.5 {
		t.Errorf("TimeSpentMin = %v, want ~10", attempt.TimeSpentMin)
	}
}

func TestSubmit_PromptsForAndStoresNotes(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("great problem\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if attempt.Notes != "great problem" {
		t.Errorf("Notes = %q, want %q", attempt.Notes, "great problem")
	}
}

func TestSubmit_TimesOutIfTestsNeverRevealed(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	cfg.RevealTimeout = 30 * time.Millisecond
	mkdirs(t, cfg)
	// deliberately no simulateHostReveal call

	_, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected timeout error when tests are never revealed, got nil")
	}
}

func TestSubmit_WritesLastTestResultFile(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "echo hello-from-test"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	if _, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("Submit: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cfg.WorkspaceDir, lastTestResultFile))
	if err != nil {
		t.Fatalf("read last test result file: %v", err)
	}

	var got lastTestResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal last test result: %v", err)
	}
	if got.Result != tracker.ResultPass {
		t.Errorf("Result = %q, want %q", got.Result, tracker.ResultPass)
	}
	if !strings.Contains(got.Output, "hello-from-test") {
		t.Errorf("Output = %q, want it to contain the command's output", got.Output)
	}
	if got.TestCommand != cfg.TestCommand {
		t.Errorf("TestCommand = %q, want %q", got.TestCommand, cfg.TestCommand)
	}
	if got.RecordedAt.IsZero() {
		t.Error("RecordedAt is zero, want a real timestamp")
	}
}

func TestSubmit_WritesAttemptToTrackerDB(t *testing.T) {
	cfg := baseConfig(t)
	cfg.TestCommand = "true"
	mkdirs(t, cfg)
	simulateHostReveal(t, cfg.ControlDir)

	attempt, err := Submit(cfg, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}

	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("tracker.Open: %v", err)
	}
	defer tr.Close()

	attempts, err := tr.ListAttempts()
	if err != nil {
		t.Fatalf("ListAttempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt logged, got %d", len(attempts))
	}
	if attempts[0].ID != attempt.ID {
		t.Errorf("logged attempt ID = %d, want %d", attempts[0].ID, attempt.ID)
	}
	if attempts[0].ExerciseID != cfg.ExerciseID {
		t.Errorf("logged ExerciseID = %q, want %q", attempts[0].ExerciseID, cfg.ExerciseID)
	}
}
