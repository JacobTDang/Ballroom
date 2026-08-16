package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/mock"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// The Capital One mock sitting (issue #280): four questions drawn from
// the capital-one slot pools, run back-to-back as ordinary exercise
// sessions against one shared 70-minute wall-clock deadline. The TUI
// owns the start screen and the summary; the sequential container loop
// runs between two program lifetimes, like every other docker handoff.

// mockChoice is the between-questions decision: continue into the next
// container, skip the upcoming question, or end the sitting.
type mockChoice int

const (
	mockContinue mockChoice = iota
	mockSkip
	mockAbort
)

// runMockSitting runs a drawn 4-question plan sequentially against one
// wall-clock deadline. Every collaborator with a side effect (clock,
// docker launch, tracker read, stdin prompt) is injected so the loop is
// unit-testable; runMockReal wires the real ones.
//
// The first launch/tracker error aborts the rest of the sitting and is
// RETURNED, not written to stderr: the very next thing the caller does
// is start a fresh alt-screen program, which would wipe a stderr line
// before anyone could read it (the issue #230 lesson) — the summary
// screen renders it instead.
func runMockSitting(
	cfg config.Config,
	plan [4]exercise.Exercise,
	now func() time.Time,
	runExercise func(config.Config, exercise.Exercise, string) error,
	attemptsFor func(exerciseID string) ([]tracker.Attempt, error),
	prompt func(question int, remainingMin int) mockChoice,
) (mock.Sitting, error) {
	start := now()
	deadline := start.Add(mock.TotalMinutes * time.Minute)
	sitting := mock.Sitting{StartedAt: start.Format(time.RFC3339)}
	for i, ex := range plan {
		sitting.ProblemIDs[i] = ex.ProblemID
		sitting.Outcomes[i] = mock.OutcomeSkipped
	}

	var firstErr error
	aborted, timedOut := false, false
	fail := func(err error) {
		if firstErr == nil {
			firstErr = err
		}
		aborted = true
	}
	for i, ex := range plan {
		if aborted {
			continue // record the rest as skipped, don't break mid-bookkeeping
		}
		remaining := mock.RemainingMinutes(deadline, now())
		if remaining == 0 {
			timedOut = true
			continue
		}
		switch prompt(i+1, remaining) {
		case mockAbort:
			aborted = true
			continue
		case mockSkip:
			continue
		}

		before, err := attemptsFor(ex.ID)
		if err != nil {
			fail(fmt.Errorf("mock: read attempts: %w", err))
			continue
		}
		baseline := mock.MaxAttemptID(before)

		questionStart := now()
		ex.TimeLimitMin = mock.RemainingMinutes(deadline, questionStart)
		if ex.TimeLimitMin == 0 {
			// The interstitial approved this launch moments ago with
			// budget left; never hand the container a zero limit.
			ex.TimeLimitMin = 1
		}
		if runErr := runExercise(cfg, ex, ""); runErr != nil {
			fail(fmt.Errorf("mock: question %d: %w", i+1, runErr))
			continue
		}
		sitting.MinutesPerSlot[i] = now().Sub(questionStart).Minutes()

		after, err := attemptsFor(ex.ID)
		if err != nil {
			fail(fmt.Errorf("mock: read attempts: %w", err))
			continue
		}
		sitting.Outcomes[i] = mock.OutcomeSince(after, baseline)
	}

	sitting.MinutesTotal = now().Sub(start).Minutes()
	// Completed means the sitting ran its course: every question was
	// offered. A voluntary skip still completes; an abort or a deadline
	// that cut questions off does not.
	sitting.Completed = !aborted && !timedOut
	return sitting, firstErr
}

// runMockReal is the production wiring: real clock, real docker launch,
// real tracker, real stdin prompt between questions (the terminal
// belongs to the plain console between docker handoffs, same as the
// moment RunExercise itself returns).
func runMockReal(cfg config.Config, plan [4]exercise.Exercise) (mock.Sitting, error) {
	attempts := func(id string) ([]tracker.Attempt, error) {
		tr, err := tracker.Open(cfg.DBPath)
		if err != nil {
			return nil, err
		}
		defer tr.Close()
		return tr.ListAttemptsFor(id)
	}
	return runMockSitting(cfg, plan, time.Now, orchestrator.RunExercise, attempts, promptMockChoice)
}

// promptMockChoice is the between-question interstitial on the plain
// terminal: Enter continues, s skips the question, q ends the sitting.
func promptMockChoice(question, remainingMin int) mockChoice {
	fmt.Printf("\nMock question %d of 4 — %d min left. [Enter] start · [s] skip · [q] end sitting: ", question, remainingMin)
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	switch strings.TrimSpace(strings.ToLower(line)) {
	case "s":
		return mockSkip
	case "q":
		return mockAbort
	}
	return mockContinue
}

// enterMockStart loads the catalog, draws the sitting, and loads the
// history line — any failure lands in m.err and keeps the user on the
// main menu rather than a half-initialized start screen.
func enterMockStart(m appModel) appModel {
	statuses, err := catalogListFn(m.cfg)
	if err != nil {
		m.err = err
		return m
	}
	plan, err := mock.Draw(statuses)
	if err != nil {
		m.err = err
		return m
	}
	sittings, err := mock.ListSittings(m.cfg.DataDir)
	if err != nil {
		m.err = err
		return m
	}
	m.err = nil
	m.mockPlan = plan
	m.mockHistory = mock.HistoryLine(sittings)
	m.stage = stageMockStart
	return m
}

func (m appModel) updateMockStart(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.outcome = outcomeRunMock
		return m, tea.Quit
	case "esc", "q":
		m.stage = stageMain
	case "?":
		return m.openHelp()
	case "ctrl+c":
		m.quit = true
		return m, tea.Quit
	}
	return m, nil
}

func (m appModel) renderMockStart() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Capital One mock — CodeSignal GCA format"))
	b.WriteString("\n\n")
	for i, ex := range m.mockPlan {
		b.WriteString(fmt.Sprintf("  Q%d  %-6s  %-34s %s\n", i+1, mock.Slots[i], truncateTitle(ex.Title, 34), ex.ProblemID))
	}
	b.WriteString("\n" + fmt.Sprintf("%d minutes, forward-only — skip costs the question, the clock never pauses.", mock.TotalMinutes))
	if m.mockHistory != "" {
		b.WriteString("\n" + menuSubtitleStyle.Render(m.mockHistory))
	}
	b.WriteString("\n\n[enter] start · [esc] back")
	return b.String()
}

func (m appModel) updateMockSummary(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc", "q":
		m.mockSitting = nil
		m.stage = stageMain
	case "ctrl+c":
		m.quit = true
		return m, tea.Quit
	}
	return m, nil
}

// outcomeGlyph maps a slot outcome to its summary-row marker.
func outcomeGlyph(outcome string) string {
	switch outcome {
	case mock.OutcomeSolved:
		return "✓ solved"
	case mock.OutcomeAttempted:
		return "~ attempted"
	default:
		return "– skipped"
	}
}

func (m appModel) renderMockSummary() string {
	s := m.mockSitting
	if s == nil {
		return m.renderMain()
	}
	var b strings.Builder
	b.WriteString(hintStyle.Render("Mock sitting — results"))
	b.WriteString("\n\n")
	for i, pid := range s.ProblemIDs {
		minutes := ""
		if s.Outcomes[i] != mock.OutcomeSkipped {
			minutes = fmt.Sprintf(" · %.1f min", s.MinutesPerSlot[i])
		}
		b.WriteString(fmt.Sprintf("  Q%d  %-6s  %-14s  %s%s\n", i+1, mock.Slots[i], pid, outcomeGlyph(s.Outcomes[i]), minutes))
	}
	b.WriteString(fmt.Sprintf("\nSolved %d/4 · %.0f min total", s.Solved(), s.MinutesTotal))
	if !s.Completed {
		b.WriteString(" · ended early")
	}
	if m.err != nil {
		b.WriteString("\n\n")
		b.WriteString(renderFriendlyError("the sitting hit an error", m.err))
	}
	b.WriteString("\n\n[enter] menu")
	return b.String()
}
