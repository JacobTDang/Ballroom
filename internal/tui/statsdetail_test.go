package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// statsDetailFixture returns three attempts against the same exercise,
// newest first (as tracker.ListAttemptsFor returns them): one full
// tutor-assisted attempt, one gave-up attempt with no tutor metadata,
// and one bare-bones attempt predating the hints_used/tutor_mode/model
// columns (NULL -- migration 002).
func statsDetailFixture() []tracker.Attempt {
	hints := 2
	mode := "hints-first"
	model := "qwen2.5-coder:14b"
	return []tracker.Attempt{
		{
			ExerciseID: "two-sum-01-python", Category: "arrays-hashing", Language: "python",
			Date: "2026-07-17", Result: tracker.ResultPass, TimeSpentMin: 12.5,
			HintsUsed: &hints, TutorMode: &mode, Model: &model, Notes: "clean on the first pass",
		},
		{
			ExerciseID: "two-sum-01-python", Category: "arrays-hashing", Language: "python",
			Date: "2026-07-10", Result: tracker.ResultGaveUp, TimeSpentMin: 30,
		},
		{
			// Predates migration 002: every nullable column is nil.
			ExerciseID: "two-sum-01-python", Category: "arrays-hashing", Language: "python",
			Date: "2026-06-01", Result: tracker.ResultFail, TimeSpentMin: 8,
		},
	}
}

// --- updateStatsDetail ---

func TestAppModel_StatsDetail_EscAndQGoBackToStatsListNotMain(t *testing.T) {
	for _, key := range []tea.KeyMsg{{Type: tea.KeyEsc}, {Type: tea.KeyRunes, Runes: []rune("q")}} {
		m := appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture()}
		newM, cmd := m.Update(key)
		if cmd != nil {
			t.Errorf("%v: expected no external command", key)
		}
		if got := newM.(appModel).stage; got != stageStats {
			t.Errorf("%v: stage = %v, want stageStats (one level back, not stageMain)", key, got)
		}
	}
}

func TestAppModel_StatsDetail_ArrowsAndJKMoveCursorWithoutExiting(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture()}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageStatsDetail {
		t.Fatalf("stage = %v, want stageStatsDetail (arrows must not exit)", got.stage)
	}
	if got.statsDetailCursor != 1 {
		t.Errorf("statsDetailCursor = %d, want 1", got.statsDetailCursor)
	}

	back, _ := got.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if backGot := back.(appModel); backGot.statsDetailCursor != 0 {
		t.Errorf("statsDetailCursor = %d after k, want 0", backGot.statsDetailCursor)
	}
}

func TestAppModel_StatsDetail_CursorClampsAtBothEnds(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture(), statsDetailCursor: 0}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if got := newM.(appModel).statsDetailCursor; got != 0 {
		t.Errorf("statsDetailCursor = %d, want 0 (clamped at top)", got)
	}

	last := len(statsDetailFixture()) - 1
	m = appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture(), statsDetailCursor: last}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if got := newM.(appModel).statsDetailCursor; got != last {
		t.Errorf("statsDetailCursor = %d, want %d (clamped at bottom)", got, last)
	}
}

func TestAppModel_StatsDetail_PgUpPgDownPageAndClamp(t *testing.T) {
	attempts := make([]tracker.Attempt, 40)
	for i := range attempts {
		attempts[i] = tracker.Attempt{ExerciseID: "two-sum-01-python", Date: "2026-07-17", Result: tracker.ResultPass}
	}
	m := appModel{stage: stageStatsDetail, statsDetailAttempts: attempts, statsDetailCursor: 0}

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	if got := newM.(appModel).statsDetailCursor; got <= 0 {
		t.Fatalf("statsDetailCursor after PgDown = %d, want it to advance", got)
	}
	for i := 0; i < 10; i++ {
		newM, _ = newM.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	}
	if got := newM.(appModel).statsDetailCursor; got != len(attempts)-1 {
		t.Errorf("statsDetailCursor = %d after repeated PgDown, want %d (clamped at the last row)", got, len(attempts)-1)
	}
	for i := 0; i < 10; i++ {
		newM, _ = newM.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	}
	if got := newM.(appModel).statsDetailCursor; got != 0 {
		t.Errorf("statsDetailCursor = %d after repeated PgUp, want 0 (clamped at the top)", got)
	}
}

func TestAppModel_StatsDetail_QuestionMarkOpensHelpAndReturnsToDetail(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture(), statsDetailCursor: 1}
	opened, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	got := opened.(appModel)
	if got.stage != stageHelp {
		t.Fatalf("stage = %v, want stageHelp", got.stage)
	}
	if got.helpOrigin != stageStatsDetail {
		t.Errorf("helpOrigin = %v, want stageStatsDetail", got.helpOrigin)
	}

	closed, _ := got.Update(tea.KeyMsg{Type: tea.KeyEsc})
	backGot := closed.(appModel)
	if backGot.stage != stageStatsDetail || backGot.statsDetailCursor != 1 {
		t.Errorf("after closing help: stage=%v statsDetailCursor=%d, want stageStatsDetail with cursor 1 preserved", backGot.stage, backGot.statsDetailCursor)
	}
}

func TestAppModel_StatsDetail_UnmappedKeyIsNoop(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture(), statsDetailCursor: 1}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageStatsDetail || got.statsDetailCursor != 1 {
		t.Errorf("unmapped key changed state: stage=%v cursor=%d, want unchanged", got.stage, got.statsDetailCursor)
	}
}

func TestAppModel_StatsDetail_EnterIsNoop(t *testing.T) {
	// There's nothing further to drill into below a single exercise's
	// history -- a deliberate no-op, not a dead end that happens to
	// swallow the key silently for some other reason.
	m := appModel{stage: stageStatsDetail, statsDetailAttempts: statsDetailFixture(), statsDetailCursor: 1}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageStatsDetail || got.statsDetailCursor != 1 {
		t.Errorf("enter changed state: stage=%v cursor=%d, want unchanged", got.stage, got.statsDetailCursor)
	}
}

// --- renderStatsDetail ---

func TestRenderStatsDetail_NeverAttemptedExerciseRendersEmptyState(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailExerciseID: "ghost-01", statsDetailAttempts: nil, width: 220, height: 60}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "ghost-01") {
		t.Errorf("expected the exercise id in the header, got:\n%s", view)
	}
	if !strings.Contains(strings.ToLower(view), "no attempts") {
		t.Errorf("expected a friendly empty state, got:\n%s", view)
	}
}

func TestRenderStatsDetail_NullColumnsRenderAsDashes(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailExerciseID: "two-sum-01-python", statsDetailAttempts: statsDetailFixture(), width: 220, height: 60}
	view := stripAnsiTUI(m.View())

	var nullRow string
	for _, l := range strings.Split(view, "\n") {
		if strings.Contains(l, "2026-06-01") {
			nullRow = l
		}
	}
	if nullRow == "" {
		t.Fatalf("couldn't find the pre-migration row (NULL hints/mode/model) in:\n%s", view)
	}
	// hints_used, tutor_mode, and model are all nil on this row -- three
	// dash-rendered columns.
	if c := strings.Count(nullRow, " - "); c < 3 {
		t.Errorf("expected at least 3 dash-rendered NULL columns in the pre-migration row, got %d: %q", c, nullRow)
	}
}

// TestFormatHintsUsed_NilIsDashZeroIsZero is the precise, non-fragile
// version of "NULL columns render as dashes, not zeros" (issue #252):
// nil (unknown -- predates migration 002) must render as the dash, and
// a genuinely-recorded zero must still say "0", never collapse into the
// same dash unknown uses.
func TestFormatHintsUsed_NilIsDashZeroIsZero(t *testing.T) {
	if got := formatHintsUsed(nil); got != statsDash {
		t.Errorf("formatHintsUsed(nil) = %q, want %q", got, statsDash)
	}
	zero := 0
	if got := formatHintsUsed(&zero); got != "0" {
		t.Errorf("formatHintsUsed(&0) = %q, want %q -- a recorded zero must not look like unknown", got, "0")
	}
	two := 2
	if got := formatHintsUsed(&two); got != "2" {
		t.Errorf("formatHintsUsed(&2) = %q, want %q", got, "2")
	}
}

func TestFormatOptional_NilAndEmptyAreDashOtherwiseAsIs(t *testing.T) {
	if got := formatOptional(nil); got != statsDash {
		t.Errorf("formatOptional(nil) = %q, want %q", got, statsDash)
	}
	empty := ""
	if got := formatOptional(&empty); got != statsDash {
		t.Errorf("formatOptional(&\"\") = %q, want %q", got, statsDash)
	}
	val := "hints-first"
	if got := formatOptional(&val); got != val {
		t.Errorf("formatOptional(&%q) = %q, want %q", val, got, val)
	}
}

func TestFormatNotes_BlankIsDashOtherwiseAsIs(t *testing.T) {
	if got := formatNotes(""); got != statsDash {
		t.Errorf("formatNotes(\"\") = %q, want %q", got, statsDash)
	}
	if got := formatNotes("   "); got != statsDash {
		t.Errorf("formatNotes(whitespace) = %q, want %q", got, statsDash)
	}
	if got := formatNotes("clean solution"); got != "clean solution" {
		t.Errorf("formatNotes(text) = %q, want it unchanged", got)
	}
}

func TestRenderStatsDetail_RecordedZeroHintsRendersAsZeroNotDash(t *testing.T) {
	// Every other field in this fixture is deliberately zero-free (no
	// bare "0" can come from anywhere else in the rendered view), so a
	// literal "0" appearing at all can only be the recorded HintsUsed.
	zero := 0
	mode := "syntax-only"
	model := "llama3.1:8b"
	attempts := []tracker.Attempt{{
		ExerciseID: "two-sum-01-python", Category: "arrays-hashing", Language: "python",
		Date: "2026-07-17", Result: tracker.ResultPass, TimeSpentMin: 9.5,
		HintsUsed: &zero, TutorMode: &mode, Model: &model,
	}}
	m := appModel{stage: stageStatsDetail, statsDetailExerciseID: "two-sum-01-python", statsDetailAttempts: attempts, width: 220, height: 60}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "0") {
		t.Errorf("expected the recorded 0 hints to render as \"0\", got:\n%s", view)
	}
}

func TestRenderStatsDetail_ShowsAllRequiredFieldsForARealAttempt(t *testing.T) {
	m := appModel{stage: stageStatsDetail, statsDetailExerciseID: "two-sum-01-python", statsDetailAttempts: statsDetailFixture(), width: 220, height: 60}
	view := stripAnsiTUI(m.View())
	for _, want := range []string{
		"2026-07-17", tracker.ResultPass, "hints-first", "qwen2.5-coder:14b", "clean on the first pass",
		tracker.ResultGaveUp,
	} {
		if !strings.Contains(view, want) {
			t.Errorf("renderStatsDetail missing %q:\n%s", want, view)
		}
	}
	if !strings.Contains(view, "12.5") {
		t.Errorf("expected the time-spent column to show 12.5, got:\n%s", view)
	}
}

func TestRenderStatsDetail_WindowsLongHistoryWithMoreIndicator(t *testing.T) {
	attempts := make([]tracker.Attempt, 30)
	for i := range attempts {
		attempts[i] = tracker.Attempt{ExerciseID: "two-sum-01-python", Date: "2026-07-17", Result: tracker.ResultPass}
	}
	m := appModel{stage: stageStatsDetail, statsDetailExerciseID: "two-sum-01-python", statsDetailAttempts: attempts, width: 220, height: 60}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "more") {
		t.Errorf("expected a windowed history with a '... more' indicator, not every row dumped at once:\n%s", view)
	}
}
