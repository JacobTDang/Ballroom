package tutor

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudwego/eino/schema"
)

// --- thinkingBorderColorAt -- the traveling gradient itself ---

func TestThinkingBorderColorAt_WrapsAroundFullPerimeter(t *testing.T) {
	const perimeter = 100
	c0 := thinkingBorderColorAt(0, perimeter, 7)
	cP := thinkingBorderColorAt(perimeter, perimeter, 7)
	r0, g0, b0 := c0.RGB255()
	rP, gP, bP := cP.RGB255()
	if r0 != rP || g0 != gP || b0 != bP {
		t.Errorf("colorAt(0) = (%d,%d,%d), colorAt(perimeter) = (%d,%d,%d), want identical -- one full lap must land back on the same color", r0, g0, b0, rP, gP, bP)
	}
}

func TestThinkingBorderColorAt_AnimatesOverPhase(t *testing.T) {
	const perimeter = 100
	c1 := thinkingBorderColorAt(5, perimeter, 0)
	c2 := thinkingBorderColorAt(5, perimeter, 25)
	r1, g1, b1 := c1.RGB255()
	r2, g2, b2 := c2.RGB255()
	if r1 == r2 && g1 == g2 && b1 == b2 {
		t.Errorf("colorAt(pos=5) is (%d,%d,%d) at both phase 0 and phase 25 -- the gradient isn't moving", r1, g1, b1)
	}
}

func TestThinkingBorderColorAt_SmoothBetweenAdjacentCells(t *testing.T) {
	// With 7 stops spread over 210 cells, adjacent cells are 1/30th of a
	// stop-to-stop blend apart -- any per-channel jump bigger than this
	// threshold means the blend degenerated into discrete color bands
	// (a genuine band between palette neighbors jumps by ~190). The
	// threshold is loose enough to accommodate Luv blending's honest
	// nonlinearity in RGB terms: perceptually-uniform steps bunch the
	// RGB delta toward one end of a wide hue gap like teal->red.
	const perimeter = 210
	const maxChannelDelta = 40
	for pos := 0; pos < perimeter; pos++ {
		r1, g1, b1 := thinkingBorderColorAt(pos, perimeter, 0).RGB255()
		r2, g2, b2 := thinkingBorderColorAt(pos+1, perimeter, 0).RGB255()
		dr, dg, db := absDiffU8(r1, r2), absDiffU8(g1, g2), absDiffU8(b1, b2)
		if dr > maxChannelDelta || dg > maxChannelDelta || db > maxChannelDelta {
			t.Fatalf("colorAt(%d)->colorAt(%d) jumps by (%d,%d,%d), want every channel <= %d -- hard color band, not a smooth blend", pos, pos+1, dr, dg, db, maxChannelDelta)
		}
	}
}

func absDiffU8(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

// --- thinkingBorderFade -- when the border shows and how it dies ---

func TestThinkingBorderFade_FullWhileTurnInFlight(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnInFlight = true
	// A stale turnSettledAt from a previous turn must not matter: a new
	// in-flight turn is always full brightness.
	m.turnSettledAt = time.Now().Add(-10 * time.Second)
	if got := m.thinkingBorderFade(); got != 1.0 {
		t.Errorf("thinkingBorderFade() = %v while turnInFlight, want 1.0", got)
	}
}

func TestThinkingBorderFade_ZeroBeforeAnyTurnCompleted(t *testing.T) {
	m := newTutorLayoutOnly()
	if got := m.thinkingBorderFade(); got != 0.0 {
		t.Errorf("thinkingBorderFade() = %v on a fresh model, want 0.0 -- no border before the first turn ever runs", got)
	}
}

func TestThinkingBorderFade_ZeroAfterFadeDuration(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnSettledAt = time.Now().Add(-thinkingBorderFadeDuration)
	if got := m.thinkingBorderFade(); got != 0.0 {
		t.Errorf("thinkingBorderFade() = %v at exactly fadeDuration elapsed, want 0.0 -- the border must be fully gone, not asymptotically dim", got)
	}
}

func TestThinkingBorderFade_MidFadeEaseOutBelowLinear(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnSettledAt = time.Now().Add(-thinkingBorderFadeDuration / 2)
	got := m.thinkingBorderFade()
	if got <= 0 || got >= 1 {
		t.Fatalf("thinkingBorderFade() = %v at half the fade duration, want strictly between 0 and 1", got)
	}
	// The ease-out curve (t*t) drops fast early and tapers at the end,
	// so at the halfway point it must sit below the 0.5 a linear ramp
	// would give. Allow a little slack for the wall-clock time that
	// passes between setting turnSettledAt and reading the fade.
	if got >= 0.5 {
		t.Errorf("thinkingBorderFade() = %v at half duration, want < 0.5 (ease-out, not linear)", got)
	}
}

// --- renderThinkingBorder -- the frame itself ---

func TestRenderThinkingBorder_FadeZeroRingIsPlainBlank(t *testing.T) {
	got := renderThinkingBorder("hi", 10, 5, 0, 0)
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("renderThinkingBorder produced %d lines, want 5", len(lines))
	}
	if strings.Contains(got, "\x1b") {
		t.Error("fade=0 output contains an ANSI escape -- the resting ring must be completely unstyled, not dimmed-to-black (this is the no-residual-sidebar guarantee)")
	}
	blank := strings.Repeat(" ", 10)
	if lines[0] != blank {
		t.Errorf("top ring row = %q, want 10 plain spaces", lines[0])
	}
	if lines[4] != blank {
		t.Errorf("bottom ring row = %q, want 10 plain spaces", lines[4])
	}
	for i := 1; i <= 3; i++ {
		if !strings.HasPrefix(lines[i], " ") || !strings.HasSuffix(lines[i], " ") {
			t.Errorf("ring row %d = %q, want a plain space at each end", i, lines[i])
		}
	}
	if !strings.Contains(lines[1], "hi") {
		t.Errorf("content row = %q, want the content preserved inside the ring", lines[1])
	}
}

func TestRenderThinkingBorder_TooSmallReturnsContentUnchanged(t *testing.T) {
	content := "tiny"
	if got := renderThinkingBorder(content, 3, 5, 0, 1.0); got != content {
		t.Errorf("renderThinkingBorder(w=3) = %q, want the content passed through untouched", got)
	}
	if got := renderThinkingBorder(content, 10, 3, 0, 1.0); got != content {
		t.Errorf("renderThinkingBorder(h=3) = %q, want the content passed through untouched", got)
	}
}

func TestRenderThinkingBorder_ActiveBorderColorsRingAndPreservesContent(t *testing.T) {
	got := renderThinkingBorder("hi", 10, 5, 0, 1.0)
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("renderThinkingBorder produced %d lines, want 5", len(lines))
	}
	if !strings.Contains(lines[0], "╭") || !strings.Contains(lines[0], "╮") {
		t.Errorf("top row = %q, want rounded corners at both ends", lines[0])
	}
	if !strings.Contains(lines[0], "\x1b[38;2;") {
		t.Errorf("top row = %q, want per-cell truecolor escapes", lines[0])
	}
	if n := strings.Count(lines[1], "│"); n != 2 {
		t.Errorf("content row = %q, want exactly 2 side glyphs, got %d", lines[1], n)
	}
	if !strings.Contains(lines[1], "hi") {
		t.Errorf("content row = %q, want the content preserved inside the border", lines[1])
	}
	if !strings.Contains(lines[4], "╰") || !strings.Contains(lines[4], "╯") {
		t.Errorf("bottom row = %q, want rounded corners at both ends", lines[4])
	}
}

func TestRenderThinkingBorder_PadsShortContentToFullHeight(t *testing.T) {
	got := renderThinkingBorder("x", 10, 6, 0, 1.0)
	if n := len(strings.Split(got, "\n")); n != 6 {
		t.Errorf("renderThinkingBorder produced %d lines, want 6 -- short content must be padded so the bottom edge stays anchored", n)
	}
}

// --- model wiring ---

func TestTurnCompleteMsg_SetsTurnSettledAt(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnInFlight = true
	before := time.Now()

	newM, _ := m.Update(turnCompleteMsg{reply: schema.AssistantMessage("done", nil), userMessage: "x"})
	got := newM.(tutorModel)

	if got.turnSettledAt.IsZero() {
		t.Fatal("turnSettledAt still zero after turnCompleteMsg -- the fade has no origin point")
	}
	if got.turnSettledAt.Before(before) {
		t.Errorf("turnSettledAt = %v, want at or after %v", got.turnSettledAt, before)
	}
}

func TestView_BorderAppearsWhileThinkingAndIsAbsentAtRest(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	m = newM.(tutorModel)

	restLines := strings.Split(m.View(), "\n")
	if len(restLines) != 12 {
		t.Fatalf("View() at rest = %d lines, want 12 (the blank ring keeps the frame footprint stable)", len(restLines))
	}
	if strings.Contains(restLines[0], "\x1b") || strings.TrimSpace(restLines[0]) != "" {
		t.Errorf("View() top row at rest = %q, want completely blank -- no visible border of any kind when idle", restLines[0])
	}

	m.turnInFlight = true
	thinkingLines := strings.Split(m.View(), "\n")
	if !strings.Contains(thinkingLines[0], "╭") || !strings.Contains(thinkingLines[0], "\x1b[38;2;") {
		t.Errorf("View() top row while thinking = %q, want a colored border corner", thinkingLines[0])
	}
}
