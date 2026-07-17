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
	ControlDir   string
	WorkspaceDir string
	TestCommand  string
	ExerciseID   string
	Category     string
	Language     string // programming language, or session style for design kind
	Kind         string // exercise.KindCoding (default "" means coding) or exercise.KindDesign
	// Grade, when set on a design submit, produces the pass/fail verdict
	// and a per-dimension summary by grading the design against the
	// revealed rubric (see tutor.GradeDesign, wired by cmd/ballroom).
	// Injected as a function rather than imported so this package stays
	// decoupled from the tutor's model plumbing and tests can fake it.
	// Nil, or any error from it, falls back to explicit self-assessment.
	Grade func() (verdict, summary string, err error)
	// CheckComplexity, when set, powers the post-pass complexity quiz on
	// coding submits: it receives the user's claimed time/space
	// complexity and returns the model's verdict text (see
	// tutor.CheckComplexity, wired by cmd/ballroom). Same
	// injected-function decoupling as Grade. Nil disables the quiz; an
	// error degrades to a notice and never blocks recording the attempt.
	CheckComplexity func(claim string) (string, error)
	// Recap, when set, writes the post-session recap: it receives the
	// attempt's result and raw grading output and returns a short
	// summary (see tutor.SessionRecap, wired by cmd/ballroom), which is
	// printed and appended to the attempt's notes tagged "[recap]".
	// Same injected-function decoupling as Grade/CheckComplexity; nil
	// disables it, an error degrades to a notice.
	Recap         func(result, output string) (string, error)
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

	var result, output, gradeSummary string
	if cfg.Kind == exercise.KindDesign {
		fmt.Fprintln(stdout, "\nrubric.md has been revealed in your workspace.")
		result, output, gradeSummary = gradeOrSelfAssess(cfg, scanner, stdout)
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

	// The complexity quiz: only after a green coding run (a failing
	// submit has more urgent things to think about, and design answers
	// have no complexity), and only when the checker is wired. Entirely
	// optional and entirely non-blocking -- empty answer skips, a model
	// error becomes a notice, the attempt records regardless.
	if cfg.Kind != exercise.KindDesign && result == tracker.ResultPass && cfg.CheckComplexity != nil {
		fmt.Fprint(stdout, "\nTime/space complexity of your solution? (enter to skip): ")
		claim := ""
		if scanner.Scan() {
			claim = strings.TrimSpace(scanner.Text())
		}
		if claim != "" {
			verdict, err := cfg.CheckComplexity(claim)
			if err != nil {
				fmt.Fprintf(stdout, "complexity check unavailable: %v\n", err)
			} else {
				fmt.Fprintf(stdout, "\ncomplexity check:\n%s\n", verdict)
			}
		}
	}

	notes := promptNotes(scanner, stdout)

	// The model-written recap lands in the same notes column, tagged so
	// it never masquerades as something the user wrote. After the notes
	// prompt on purpose: the user's own reflection comes uncontaminated
	// by the model's.
	if cfg.Recap != nil {
		recap, err := cfg.Recap(result, output)
		if err != nil {
			fmt.Fprintf(stdout, "recap unavailable: %v\n", err)
		} else {
			fmt.Fprintf(stdout, "\nsession recap:\n%s\n", recap)
			if notes != "" {
				notes += "\n\n"
			}
			notes += "[recap] " + recap
		}
	}

	attempt := tracker.Attempt{
		ExerciseID:   cfg.ExerciseID,
		Category:     cfg.Category,
		Language:     cfg.Language,
		Date:         time.Now().Format("2006-01-02"),
		TimeSpentMin: time.Since(cfg.StartedAt).Minutes(),
		Result:       result,
		Notes:        notes,
		GradeSummary: gradeSummary,
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

// gradeOrSelfAssess produces a design submit's result: the injected
// grader's verdict when it succeeds (summary shown to the user and
// persisted for read_test_output), or the explicit self-assessment
// prompt when no grader is wired or grading fails. Grading failures are
// printed, never swallowed -- a free-tier model hiccup must be visible,
// and must degrade to the human answering, not to a silent guess.
func gradeOrSelfAssess(cfg Config, scanner *bufio.Scanner, stdout io.Writer) (result, output, gradeSummary string) {
	if cfg.Grade != nil {
		verdict, summary, err := cfg.Grade()
		if err == nil {
			fmt.Fprintf(stdout, "\ntutor grade:\n%s\n\n", summary)
			return verdict, "(design session: model-graded)\n\n" + summary, summary
		}
		fmt.Fprintf(stdout, "\ngrading failed (%v); falling back to self-assessment\n", err)
	} else {
		fmt.Fprintln(stdout, "open it in the editor (M-1) or ask the tutor for a graded assessment (M-2) before you assess yourself.")
	}
	return promptSelfAssessment(scanner, stdout), "(design session: self-assessed)", ""
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
