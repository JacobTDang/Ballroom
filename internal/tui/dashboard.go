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

	// minPanelWidth/Height keep the panel usable on a tiny terminal
	// instead of collapsing to nothing.
	minPanelWidth  = 30
	minPanelHeight = 12

	dashboardGapWidth = 4
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

// sharedBallGrid is built once at package init, at a fixed size — the
// ball reads as the same, consistent shape on every screen and every
// resize; only its sparkle colors animate.
var sharedBallGrid = buildDiscoBall(discoBallHeight, discoBallWidth)

// panelDimensions is the panel's total rendered size (border included):
// the terminal size minus a small fixed margin, clamped so it never
// collapses to an unusable size on a tiny terminal. It does NOT depend on
// how much content the panel currently holds — the panel stays the same
// size as checks/build-log lines come and go, instead of resizing itself
// around whatever content happens to be present.
func panelDimensions(termW, termH int) (w, h int) {
	w = termW - dashboardMarginW
	h = termH - dashboardMarginH
	if w < minPanelWidth {
		w = minPanelWidth
	}
	if h < minPanelHeight {
		h = minPanelHeight
	}
	return w, h
}

// renderDashboardPanel joins the shared ball grid and right-column
// content into the shared bordered two-column panel, sized to fill most
// of the terminal at a fixed, content-independent size (lipgloss wraps
// rather than grows the box past it) and top-aligned so the title/menu
// always starts at the same spot regardless of ball height.
func renderDashboardPanel(termW, termH, phase int, right string) string {
	panelW, panelH := panelDimensions(termW, termH)
	ball := renderDiscoBall(sharedBallGrid, phase)
	content := lipgloss.JoinHorizontal(lipgloss.Top, ball, dashboardGap, right)
	// Width/Height set the box excluding the border, so subtract it here
	// to make the final rendered panel exactly panelW x panelH.
	return dashboardPanelStyle.Width(panelW - 2).Height(panelH - 2).Render(content)
}
