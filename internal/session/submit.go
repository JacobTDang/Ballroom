// Package session implements the in-container side of `practice submit`:
// request the hidden tests be revealed, run test_command once they land,
// and log the result. Runs inside the practice container, talking to the
// host-side orchestrator only through files under a shared control dir
// (bind-mounted from the host) — no network, no Docker socket.
package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// lastTestResultFile is the well-known dotfile written into a workspace
// after every submit, holding the raw test output the tutor's
// read_test_output tool reads from (internal/tutor). Living directly in
// the workspace mirrors how the tutor already finds solution.*/problem.md
// there — no DB schema change, and it's naturally scoped to just this
// exercise run.
const lastTestResultFile = ".ballroom-last-test-result.json"

// lastTestResult is the JSON shape written to lastTestResultFile.
type lastTestResult struct {
	Result      string    `json:"result"`
	Output      string    `json:"output"`
	TestCommand string    `json:"test_command"`
	RecordedAt  time.Time `json:"recorded_at"`
}

// Config describes one submit invocation. All paths are as seen from
// inside the container.
type Config struct {
	ControlDir    string
	WorkspaceDir  string
	TestCommand   string
	ExerciseID    string
	Category      string
	Language      string // programming language, or session style for design kind
	Kind          string // exercise.KindCoding (default "" means coding) or exercise.KindDesign
	StartedAt     time.Time
	DBPath        string
	PollInterval  time.Duration
	RevealTimeout time.Duration
}

// Submit requests the hidden content be revealed, waits for it, grades
// the attempt, logs it, and returns it. For a coding exercise the
// reveal is the hidden tests and grading is running test_command; for a
// design exercise the reveal is the grading rubric (tests/<id>/rubric.md
// lands in the workspace through the very same handshake) and grading
// is an explicit self-assessment -- there is nothing to run.
func Submit(cfg Config, stdin io.Reader, stdout io.Writer) (tracker.Attempt, error) {
	if err := requestReveal(cfg); err != nil {
		return tracker.Attempt{}, err
	}
	if err := waitForReady(cfg); err != nil {
		return tracker.Attempt{}, err
	}

	// One scanner shared by every prompt below: a bufio.Scanner reads
	// ahead of what it returns, so a second scanner on the same stdin
	// would start past input the first one already buffered -- with two
	// prompts in the design path, the notes prompt would silently lose
	// its line.
	scanner := bufio.NewScanner(stdin)

	var result, output string
	if cfg.Kind == exercise.KindDesign {
		fmt.Fprintln(stdout, "\nrubric.md has been revealed in your workspace -- open it in the editor (M-1) or ask the tutor for a graded assessment (M-2) before you assess yourself.")
		result = promptSelfAssessment(scanner, stdout)
		output = "(design session: self-assessed)"
	} else {
		result, output = runTestCommand(cfg)
		fmt.Fprintf(stdout, "\nresult: %s\n%s\n", result, output)
	}

	if err := writeLastTestResult(cfg, result, output); err != nil {
		fmt.Fprintf(stdout, "warning: could not save test output for the tutor: %v\n", err)
		// Not fatal — same graceful-degradation philosophy as the rest of
		// the tutor-adjacent code. A submission should still get logged
		// to the tracker DB even if this write fails.
	}

	notes := promptNotes(scanner, stdout)

	attempt := tracker.Attempt{
		ExerciseID:   cfg.ExerciseID,
		Category:     cfg.Category,
		Language:     cfg.Language,
		Date:         time.Now().Format("2006-01-02"),
		TimeSpentMin: time.Since(cfg.StartedAt).Minutes(),
		Result:       result,
		Notes:        notes,
	}

	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		return tracker.Attempt{}, fmt.Errorf("session: open tracker: %w", err)
	}
	defer tr.Close()

	id, err := tr.LogAttempt(attempt)
	if err != nil {
		return tracker.Attempt{}, fmt.Errorf("session: log attempt: %w", err)
	}
	attempt.ID = id
	return attempt, nil
}

func requestReveal(cfg Config) error {
	path := filepath.Join(cfg.ControlDir, "submit.request")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		return fmt.Errorf("session: write submit request: %w", err)
	}
	return nil
}

func waitForReady(cfg Config) error {
	readyPath := filepath.Join(cfg.ControlDir, "tests.ready")
	deadline := time.Now().Add(cfg.RevealTimeout)

	for time.Now().Before(deadline) {
		if _, err := os.Stat(readyPath); err == nil {
			return nil
		}
		time.Sleep(cfg.PollInterval)
	}
	return fmt.Errorf("session: timed out after %s waiting for hidden tests to be revealed", cfg.RevealTimeout)
}

func runTestCommand(cfg Config) (result string, output string) {
	cmd := exec.Command("sh", "-c", cfg.TestCommand)
	cmd.Dir = cfg.WorkspaceDir
	out, err := cmd.CombinedOutput()

	result = tracker.ResultPass
	if err != nil {
		result = tracker.ResultFail
	}
	return result, string(out)
}

// writeLastTestResult persists result/output to lastTestResultFile in
// the workspace so the tutor's read_test_output tool has something to
// read after this submit — see internal/tutor for the reader.
func writeLastTestResult(cfg Config, result, output string) error {
	data, err := json.Marshal(lastTestResult{
		Result:      result,
		Output:      output,
		TestCommand: cfg.TestCommand,
		RecordedAt:  time.Now(),
	})
	if err != nil {
		return fmt.Errorf("session: marshal last test result: %w", err)
	}

	path := filepath.Join(cfg.WorkspaceDir, lastTestResultFile)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("session: write last test result: %w", err)
	}
	return nil
}

// promptSelfAssessment asks for an explicit pass/fail verdict on a
// design session and loops until it gets one -- no default, because the
// answer feeds the same tracker rows and "solved" stats that real test
// runs do, and a silently-recorded result would corrupt them. EOF
// before an answer records a fail: the honest floor, never a free pass.
func promptSelfAssessment(scanner *bufio.Scanner, stdout io.Writer) string {
	for {
		fmt.Fprint(stdout, "Self-assessment -- did your design meet the rubric? Type pass or fail (p/f): ")
		if !scanner.Scan() {
			fmt.Fprintln(stdout, "\nno answer before EOF; recording fail")
			return tracker.ResultFail
		}
		switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
		case "p", "pass":
			return tracker.ResultPass
		case "f", "fail":
			return tracker.ResultFail
		}
	}
}

func promptNotes(scanner *bufio.Scanner, stdout io.Writer) string {
	fmt.Fprint(stdout, "Notes (optional): ")
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}
