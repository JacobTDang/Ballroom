package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestWrapRightColumn_NoLineExceedsWidth is the direct contract: after
// wrapping, no right-column line is wider than its column, so none can
// be soft-wrapped by the panel onto the disco ball.
func TestWrapRightColumn_NoLineExceedsWidth(t *testing.T) {
	body := strings.Join([]string{
		"today: Diameter of Binary Tree",
		"Choose your worker and orchestrator models — local (Ollama) or API (OpenRouter)",
		"recent: contains-duplicate-01-python · contains-duplicate-01-python · valid-anagram-01-go",
	}, "\n")
	const avail = 40
	for _, l := range strings.Split(wrapRightColumn(body, avail), "\n") {
		if w := lipgloss.Width(l); w > avail {
			t.Errorf("line %q is %d wide, want <= %d", l, w, avail)
		}
	}
}

func TestWrapRightColumn_ShortLinesUnchanged(t *testing.T) {
	body := "1. Practice\n\n2. Daily"
	if got := wrapRightColumn(body, 40); got != body {
		t.Errorf("short body was altered:\n%q", got)
	}
}

func TestWrapRightColumn_PreservesEveryWord(t *testing.T) {
	// Wrapping must not drop text -- it flows to the next line, it
	// doesn't truncate.
	line := "Choose your worker and orchestrator models — local (Ollama) or API (OpenRouter)"
	wrapped := wrapRightColumn(line, 30)
	for _, word := range strings.Fields(line) {
		if !strings.Contains(wrapped, word) {
			t.Errorf("wrapping dropped %q:\n%s", word, wrapped)
		}
	}
}

func TestWrapRightColumn_ZeroWidthIsNoop(t *testing.T) {
	body := "anything at all"
	if got := wrapRightColumn(body, 0); got != body {
		t.Errorf("zero avail should pass through, got %q", got)
	}
}

// TestRenderDashboardPanel_RightTextNeverEntersBallColumns is the
// regression at the render level: the reported bug was right-column
// text (the Settings subtitle) wrapping onto the disco ball's columns.
// The ball occupies the left of the panel; no row may carry a
// distinctive right-column word in that left region.
func TestRenderDashboardPanel_RightTextNeverEntersBallColumns(t *testing.T) {
	// The exact content from the report: a long menu subtitle.
	body := "5. Settings\nChoose your worker and orchestrator models — local (Ollama) or API (OpenRouter)"

	for _, termW := range []int{96, 110, 132, 150} {
		out := renderDashboardPanel(termW, 42, 0, body, layoutTop, "")
		panelW, panelH := panelDimensions(termW, 42)
		_, ballW := ballDimensions(panelW-dashboardBorderPadW, panelH-dashboardBorderPadH)
		for _, row := range strings.Split(out, "\n") {
			plain := stripAnsiTUI(row)
			// "OpenRouter" is a right-column-only word; it must never
			// appear starting within the ball's left columns.
			if idx := strings.Index(plain, "OpenRouter"); idx >= 0 && idx < ballW {
				t.Errorf("width %d: right-column text at column %d is inside the ball (width %d):\n%s",
					termW, idx, ballW, plain)
			}
		}
	}
}
