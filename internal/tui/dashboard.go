package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// The two-column dashboard layout (disco ball left, content right, framed
// in a single bordered panel) is shared by the menu and boot screens so
// they read as one consistent app shell rather than two different UIs.

const (
	// dashboardMarginW/H leave a little breathing room from the terminal
	// edges — the panel should read as "almost full screen", not
	// literally edge-to-edge.
	dashboardMarginW = 4
	dashboardMarginH = 2

	// dashboardBorderPadW/H account for dashboardPanelStyle's own
	// border (1 cell each side) and padding (1 row / 3 cols each side).
	dashboardBorderPadW = 8
	dashboardBorderPadH = 4

	dashboardGapWidth = 4

	// minBallHeight keeps the ball from collapsing into an unrecognizable
	// scribble on a very small or not-yet-sized terminal.
	minBallHeight = 10
)

var dashboardGap = func() string {
	s := ""
	for i := 0; i < dashboardGapWidth; i++ {
		s += " "
	}
	return s
}()

var dashboardPanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#9B5FB0")).
	Padding(1, 3)

// ballAreaSize computes how much space is left for the disco ball once
// the terminal margins, the panel's own border/padding, the gap between
// columns, and the right column's width are all reserved.
func ballAreaSize(termW, termH, rightColWidth int) (maxW, maxH int) {
	maxW = termW - dashboardMarginW - dashboardBorderPadW - rightColWidth - dashboardGapWidth
	maxH = termH - dashboardMarginH - dashboardBorderPadH
	return maxW, maxH
}

// dashboardBallSize picks (height, width) for the ball within the given
// budget, preferring to use the full available height — width is always
// exactly 2x height to read as circular (terminal cells are roughly
// twice as tall as wide) — and only shrinking to fit a tight width.
func dashboardBallSize(maxW, maxH int) (h, w int) {
	h = maxH
	if h < minBallHeight {
		h = minBallHeight
	}
	w = h * 2
	if w > maxW {
		w = maxW
		h = w / 2
		if h < minBallHeight {
			h = minBallHeight
		}
	}
	return h, w
}

// renderDashboardPanel joins a pre-built ball grid and right-column
// content into the shared bordered two-column panel, top-aligned so the
// title/menu always starts at the same spot regardless of how much
// taller the ball is — matching the original, smaller layout's look.
func renderDashboardPanel(ballGrid [][]discoBallCell, phase int, right string) string {
	ball := renderDiscoBall(ballGrid, phase)
	return dashboardPanelStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, ball, dashboardGap, right))
}
