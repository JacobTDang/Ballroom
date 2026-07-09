package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
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

	// minPanelWidth is big enough that the ball, the gap, and at least a
	// little text always fit without lipgloss having to word-wrap the
	// joined two-column block — that wrapping doesn't respect the column
	// boundaries and mangles the layout instead of degrading cleanly.
	minPanelWidth  = discoBallWidth + dashboardBorderPadW + dashboardGapWidth + 20
	minPanelHeight = 12

	// bannerScaleSmall is the pixel-scale used for the animated mosaic
	// banner in the right column — small enough to sit beside the disco
	// ball instead of the full-width scale used nowhere else anymore.
	bannerScaleSmall = 1
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

// dashboardBanner picks the full animated mosaic "BALLROOM" banner when
// it fits within availWidth, or falls back to the compact single-line
// wordmark when the terminal is too narrow for it. Without this, a
// terminal narrower than the full banner would force lipgloss to
// word-wrap the already-joined ball+banner block, which doesn't respect
// the column boundaries and spills text outside the panel border.
func dashboardBanner(phase, availWidth int) string {
	full := catalog.MosaicBannerScaled(phase, bannerScaleSmall)
	if lipgloss.Width(full) <= availWidth {
		return full
	}
	return catalog.CompactBanner()
}

// renderDashboardPanel joins the shared ball grid and an animated banner
// (sized to fit) with the given right-column body into the shared
// bordered two-column panel, sized to fill most of the terminal at a
// fixed, content-independent size, top-aligned so the title/menu always
// starts at the same spot regardless of ball height.
func renderDashboardPanel(termW, termH, phase int, rightBody string) string {
	panelW, panelH := panelDimensions(termW, termH)
	innerW := panelW - dashboardBorderPadW
	rightAvail := innerW - discoBallWidth - dashboardGapWidth

	right := dashboardBanner(phase, rightAvail) + "\n\n" + rightBody

	ball := renderDiscoBall(sharedBallGrid, phase)
	content := lipgloss.JoinHorizontal(lipgloss.Top, ball, dashboardGap, right)
	// Width/Height set the box excluding the border, so subtract it here
	// to make the final rendered panel exactly panelW x panelH.
	return dashboardPanelStyle.Width(panelW - 2).Height(panelH - 2).Render(content)
}
