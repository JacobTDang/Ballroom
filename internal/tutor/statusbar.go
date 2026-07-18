package tutor

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// The bottom status bar replaces the old two-row header: the session's
// identity (mode pill, model, hints) on the left, transient state
// (scroll position, endpoint, exit hint) on the right, on one
// background-filled row pinned under the input box. Bottom rather than
// top because the eye rests at the input line between turns — the
// pane's identity sits next to where the user already looks, and the
// transcript gains the freed rows.
//
// Everything here styles with the raw ansiFg/ansiBg escapes rather
// than lipgloss color styles: lipgloss routes colors through termenv
// profile detection, which can strip them entirely under `go test`
// (no TTY), and this bar's content is pinned by string-assertion
// tests (statusbar_test.go) — the same reason markdown.go hand-rolls
// its escapes.

// statusBarHeight is the bar's row count; recomputeLayout subtracts it
// from the viewport's budget the way headerHeight used to be.
const statusBarHeight = 1

// paneModeColor is the mode pill's background: one color per mode
// family, so the session's contract is readable at a glance — gold for
// the hint-budget drill, pink for full assistance, teal for the
// guided/syntax modes, red for timed interviewer mocks. Unknown modes
// get the structural paneRule rather than a loud accent.
func paneModeColor(mode string) string {
	switch mode {
	case exercise.TutorModeHintsFirst:
		return trafficGold
	case exercise.TutorModeFullAssist:
		return panePink
	case exercise.TutorModeSyntaxOnly, exercise.TutorModeDesignCoach, exercise.TutorModeStoryCoach:
		return paneTeal
	case exercise.TutorModeInterviewer, exercise.TutorModeBehavioralInterviewer:
		return trafficRed
	default:
		return paneRule
	}
}

// statusLeftText is the bar's left-side content after the mode pill:
// which model(s) this session runs on, plus the hint count in
// hints-first (the count the mode's prompt machinery already tracks —
// surfaced because "first ask vs repeat ask" is that mode's whole
// drill). Split from statusBarView so tests can assert content without
// caring about width arithmetic, same contract the old
// headerStatusText had.
func (m tutorModel) statusLeftText() string {
	model := m.cfg.Model
	if m.routingEnabled {
		model = m.cfg.Model + " +" + m.cfg.OrchestratorModel
	}
	s := " " + model
	if m.cfg.Mode == exercise.TutorModeHintsFirst {
		s += fmt.Sprintf(" · hints: %d", m.helpRequestCount)
	}
	return s
}

// statusEndpointText is where requests go — both endpoints when
// routing splits them, same rule the old header used.
func (m tutorModel) statusEndpointText() string {
	endpoint := m.workerEndpoint
	if m.routingEnabled && m.orchestratorEndpoint != endpoint {
		endpoint = endpoint + " / " + m.orchestratorEndpoint
	}
	return endpoint
}

// statusBarView renders the bar: exactly one row, exactly m.width
// cells, at every width — a wrapped or overflowing bar would corrupt
// the fixed layout arithmetic exactly the way a wrapped header used
// to. As width shrinks the right side gives way piecewise (endpoint
// first, then the scroll percentage, then everything), and as a last
// resort the left half truncates with an ellipsis.
func (m tutorModel) statusBarView() string {
	bg := ansiBg(paneStatusBg)
	// The pill: mode name on its mode color, dark text for contrast
	// (the pane's own card background doubles as the darkest ink in
	// the palette). The trailing bg re-arms the row background the
	// pill's own background replaced.
	pill := ansiBg(paneModeColor(m.cfg.Mode)) + ansiFg(cardBg) + mdBoldOn +
		" " + strings.ToUpper(m.cfg.Mode) + " " + mdBoldOff + mdColorReset + bg
	left := pill + mdDimColor + m.statusLeftText() + mdColorReset

	scroll := fmt.Sprintf("scroll %d%%", int(m.viewport.ScrollPercent()*100+0.5))
	const exit = "ctrl+d exit "
	// Right-side variants, widest first; the first that fits wins.
	// The exit hint outlives everything else: it's the one piece a
	// stuck user actually needs.
	rights := []string{
		mdDimColor + scroll + " · " + m.statusEndpointText() + " · " + mdColorReset + exit,
		mdDimColor + scroll + " · " + mdColorReset + exit,
		exit,
		"",
	}

	if m.width <= 0 {
		// No real size yet (startup, before the first WindowSizeMsg):
		// render unpadded rather than guessing a width.
		return bg + left + " " + exit + "\x1b[39;49m"
	}

	for _, right := range rights {
		gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
		if gap >= 1 || (right == "" && gap >= 0) {
			return bg + left + strings.Repeat(" ", gap) + right + "\x1b[39;49m"
		}
	}

	// Even the bare left half overflows: truncate it to the pane.
	trunc := ansi.Truncate(left, m.width, "…")
	pad := max(m.width-lipgloss.Width(trunc), 0)
	return bg + trunc + strings.Repeat(" ", pad) + "\x1b[39;49m"
}
