package tutor

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func TestPaneModeColor_MapsEveryModeFamily(t *testing.T) {
	cases := map[string]string{
		exercise.TutorModeHintsFirst:            trafficGold,
		exercise.TutorModeFullAssist:            panePink,
		exercise.TutorModeSyntaxOnly:            paneTeal,
		exercise.TutorModeDesignCoach:           paneTeal,
		exercise.TutorModeStoryCoach:            paneTeal,
		exercise.TutorModeInterviewer:           trafficRed,
		exercise.TutorModeBehavioralInterviewer: trafficRed,
		"no-such-mode":                          paneRule,
	}
	for mode, want := range cases {
		if got := paneModeColor(mode); got != want {
			t.Errorf("paneModeColor(%q) = %q, want %q", mode, got, want)
		}
	}
}

// TestStatusLeftText_HintsFirstShowsTheHintBudget: hints-first's whole
// contract revolves around "first ask vs repeat ask", but the count
// lived only in model state the user couldn't see. The status bar is
// the always-visible place for it. (Ported from the old header.)
func TestStatusLeftText_HintsFirstShowsTheHintBudget(t *testing.T) {
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "test-model", Mode: exercise.TutorModeHintsFirst}

	if got := m.statusLeftText(); !strings.Contains(got, "hints: 0") {
		t.Errorf("statusLeftText() = %q, want the zero hint count shown before any request", got)
	}
	m.helpRequestCount = 3
	if got := m.statusLeftText(); !strings.Contains(got, "hints: 3") {
		t.Errorf("statusLeftText() = %q, want the live hint count", got)
	}
}

func TestStatusLeftText_OtherModesShowNoHintCount(t *testing.T) {
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "test-model", Mode: exercise.TutorModeFullAssist}
	m.helpRequestCount = 2
	if got := m.statusLeftText(); strings.Contains(got, "hints") {
		t.Errorf("statusLeftText() = %q, want no hint count outside hints-first", got)
	}
}

func TestStatusLeftText_RoutingNamesBothModels(t *testing.T) {
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "worker-model", OrchestratorModel: "orchestrator-model", Mode: exercise.TutorModeSyntaxOnly}
	m.routingEnabled = true
	got := m.statusLeftText()
	if !strings.Contains(got, "worker-model") || !strings.Contains(got, "orchestrator-model") {
		t.Errorf("statusLeftText() = %q, want it to name both models", got)
	}
}

// statusBarModel builds a layout-only model with a fixed identity and
// width, the shared setup for the width-behavior tests below.
func statusBarModel(t *testing.T, width int) tutorModel {
	t.Helper()
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "m", Mode: exercise.TutorModeSyntaxOnly}
	m.workerEndpoint = "http://long-endpoint:11434"
	newM, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 24})
	return newM.(tutorModel)
}

// TestStatusBarView_ExactlyOneRowAtExactWidth pins the bar's structural
// contract: exactly one row, exactly the pane's width, at any width —
// a wrapped or overflowing bar would corrupt the fixed layout
// arithmetic the same way a wrapped header used to.
func TestStatusBarView_ExactlyOneRowAtExactWidth(t *testing.T) {
	for _, width := range []int{12, 20, 40, 100} {
		bar := statusBarModel(t, width).statusBarView()
		if strings.Contains(bar, "\n") {
			t.Errorf("width %d: statusBarView contains a newline, want exactly one row", width)
		}
		if got := lipgloss.Width(bar); got != width {
			t.Errorf("width %d: statusBarView renders %d cells wide, want exactly %d", width, got, width)
		}
	}
}

func TestStatusBarView_ModePillIsUppercaseWithModeColor(t *testing.T) {
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "test-model", Mode: exercise.TutorModeHintsFirst}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	bar := m.statusBarView()
	if !strings.Contains(stripAnsiTest(bar), "H I N T S - F I R S T") {
		t.Errorf("statusBarView = %q, want the uppercase letterspaced mode pill", stripAnsiTest(bar))
	}
	if !strings.Contains(bar, ansiBg(trafficGold)) {
		t.Error("statusBarView missing the hints-first pill's gold background escape")
	}
}

func TestStatusBarView_WideShowsScrollEndpointAndExit(t *testing.T) {
	bar := stripAnsiTest(statusBarModel(t, 120).statusBarView())
	for _, want := range []string{"http://long-endpoint:11434", "scroll", "ctrl+d exit"} {
		if !strings.Contains(bar, want) {
			t.Errorf("wide statusBarView = %q, want it to contain %q", bar, want)
		}
	}
}

// TestStatusBarView_NarrowDropsEndpointThenScroll: as width shrinks the
// right side gives way piecewise — endpoint first, then the scroll
// percentage — and the exit hint survives the longest, mirroring how
// the old header dropped its endpoint half. The two widths below (60,
// 45) are wider than they'd need to be pre-letterspacing: the mode
// pill's own width roughly doubled (e.g. "SYNTAX-ONLY" -> "S Y N T A X -
// O N L Y"), which shifts every one of these gives-way thresholds out
// by the same amount -- picked empirically against statusBarView's
// actual output rather than hand-computed, the same way the pill's own
// width isn't hand-computed elsewhere in this file.
func TestStatusBarView_NarrowDropsEndpointThenScroll(t *testing.T) {
	bar60 := stripAnsiTest(statusBarModel(t, 60).statusBarView())
	if strings.Contains(bar60, "long-endpoint") {
		t.Errorf("width 60: statusBarView = %q, want the endpoint dropped first", bar60)
	}
	if !strings.Contains(bar60, "scroll") || !strings.Contains(bar60, "ctrl+d exit") {
		t.Errorf("width 60: statusBarView = %q, want scroll %% and the exit hint kept", bar60)
	}

	bar45 := stripAnsiTest(statusBarModel(t, 45).statusBarView())
	if strings.Contains(bar45, "scroll") {
		t.Errorf("width 45: statusBarView = %q, want the scroll %% dropped after the endpoint", bar45)
	}
	if !strings.Contains(bar45, "ctrl+d exit") {
		t.Errorf("width 45: statusBarView = %q, want the exit hint kept", bar45)
	}
}

func TestTutorModel_View_LastLineIsStatusBar(t *testing.T) {
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "test-model", Mode: exercise.TutorModeHintsFirst}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	lines := strings.Split(m.View(), "\n")
	last := stripAnsiTest(lines[len(lines)-1])
	if !strings.Contains(last, "H I N T S - F I R S T") || !strings.Contains(last, "ctrl+d exit") {
		t.Errorf("View's last line = %q, want the status bar pinned there", last)
	}
	first := stripAnsiTest(lines[0])
	if strings.Contains(first, "test-model") {
		t.Errorf("View's first line = %q, want no header row above the viewport anymore", first)
	}
}

// TestTutorModel_View_ExactVerticalBudget tightens the old
// sum-within-terminal check to equality: with no activity region the
// viewport, textarea (plus its border frame), and status bar must
// account for every terminal row exactly — a gap would show as a
// blank band, an excess would scroll the pane.
func TestTutorModel_View_ExactVerticalBudget(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := newM.(tutorModel)

	frame := textareaBoxStyle.GetVerticalFrameSize()
	total := got.viewport.Height + got.textarea.Height() + frame + statusBarHeight
	if total != 30 {
		t.Errorf("viewport(%d) + textarea(%d) + frame(%d) + statusbar(%d) = %d, want exactly the terminal height 30",
			got.viewport.Height, got.textarea.Height(), frame, statusBarHeight, total)
	}
}
