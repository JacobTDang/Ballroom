package tutor

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudwego/eino/schema"
)

// --- auroraColorAt -- the drifting mesh-gradient field ---

func TestAuroraColorAt_SmoothInSpace(t *testing.T) {
	// Neighboring sample points a terminal-cell apart (1/100th of the
	// pane) must blend, not band -- a hard edge between adjacent cells
	// would read as blocky artifacts, the opposite of the soft mesh
	// look this exists for.
	const maxChannelDelta = 12
	for i := 0; i < 100; i++ {
		u := float64(i) / 100.0
		r1, g1, b1 := auroraColorAt(u, 0.5, 3.0)
		r2, g2, b2 := auroraColorAt(u+0.01, 0.5, 3.0)
		if d := maxAbsChannelDelta(r1, g1, b1, r2, g2, b2); d > maxChannelDelta {
			t.Fatalf("auroraColorAt(%v)->%v jumps by %d in one cell, want <= %d", u, u+0.01, d, maxChannelDelta)
		}
	}
}

func TestAuroraColorAt_AnimatesOverTime(t *testing.T) {
	r1, g1, b1 := auroraColorAt(0.3, 0.3, 0.0)
	r2, g2, b2 := auroraColorAt(0.3, 0.3, 5.0)
	if r1 == r2 && g1 == g2 && b1 == b2 {
		t.Errorf("auroraColorAt(0.3,0.3) = (%v,%v,%v) at both t=0 and t=5 -- the field isn't drifting", r1, g1, b1)
	}
}

func maxAbsChannelDelta(r1, g1, b1, r2, g2, b2 float64) int {
	d := 0
	for _, pair := range [][2]float64{{r1, r2}, {g1, g2}, {b1, b2}} {
		diff := int((pair[0] - pair[1]) * 255)
		if diff < 0 {
			diff = -diff
		}
		if diff > d {
			d = diff
		}
	}
	return d
}

// --- overlayAurora -- compositing the field behind styled content ---

func TestOverlayAurora_ZeroBrightnessReturnsContentUnchanged(t *testing.T) {
	content := "hello\nworld"
	if got := overlayAurora(content, 10, 4, 1.0, 0); got != content {
		t.Errorf("overlayAurora(brightness=0) = %q, want the content byte-identical -- no background may remain at rest", got)
	}
}

func TestOverlayAurora_PaintsBackgroundOnEveryCell(t *testing.T) {
	got := overlayAurora("hi", 6, 3, 1.0, 0.35)
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("overlayAurora produced %d lines, want 3 (content padded to full height)", len(lines))
	}
	for i, line := range lines {
		if n := strings.Count(line, "\x1b[48;2;"); n != 6 {
			t.Errorf("line %d has %d background escapes, want 6 (one per cell, full width)", i, n)
		}
		if !strings.HasSuffix(line, "\x1b[49m") {
			t.Errorf("line %d = %q, want it to end by restoring the default background so nothing bleeds past the pane", i, line)
		}
	}
}

func TestOverlayAurora_PreservesGlyphsAndForegroundStyling(t *testing.T) {
	fg := "\x1b[38;2;10;20;30m"
	content := fg + "hi\x1b[0m"
	got := overlayAurora(content, 6, 1, 1.0, 0.35)
	if !strings.Contains(got, fg) {
		t.Errorf("overlayAurora dropped the content's own foreground escape %q", fg)
	}
	stripped := ansi.Strip(strings.Split(got, "\n")[0])
	if strings.TrimRight(stripped, " ") != "hi" {
		t.Errorf("stripped output = %q, want the original glyphs %q (padded with spaces)", stripped, "hi")
	}
}

func TestOverlayAurora_AnimatesOverTime(t *testing.T) {
	a := overlayAurora("hi", 8, 2, 0.0, 0.35)
	b := overlayAurora("hi", 8, 2, 5.0, 0.35)
	if a == b {
		t.Error("overlayAurora output identical at t=0 and t=5 -- the background isn't animating")
	}
}

func TestOverlayAurora_KeepsContentTallerThanPane(t *testing.T) {
	content := "1\n2\n3\n4\n5"
	got := overlayAurora(content, 4, 3, 1.0, 0.35)
	if n := len(strings.Split(got, "\n")); n != 5 {
		t.Errorf("overlayAurora produced %d lines for 5-line content, want 5 -- a background paints behind content, it never clips it", n)
	}
}

// --- auroraFade -- lifecycle (was the border's fade, same semantics) ---

func TestAuroraFade_FullWhileTurnInFlight(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnInFlight = true
	m.turnSettledAt = time.Now().Add(-10 * time.Second)
	if got := m.auroraFade(); got != 1.0 {
		t.Errorf("auroraFade() = %v while turnInFlight, want 1.0", got)
	}
}

func TestAuroraFade_ZeroBeforeAnyTurnCompleted(t *testing.T) {
	m := newTutorLayoutOnly()
	if got := m.auroraFade(); got != 0.0 {
		t.Errorf("auroraFade() = %v on a fresh model, want 0.0", got)
	}
}

func TestAuroraFade_ZeroAfterFadeDuration(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnSettledAt = time.Now().Add(-auroraFadeDuration)
	if got := m.auroraFade(); got != 0.0 {
		t.Errorf("auroraFade() = %v at exactly fadeDuration elapsed, want 0.0 -- fully gone, not asymptotically dim", got)
	}
}

func TestAuroraFade_MidFadeEaseOutBelowLinear(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnSettledAt = time.Now().Add(-auroraFadeDuration / 2)
	got := m.auroraFade()
	if got <= 0 || got >= 0.5 {
		t.Errorf("auroraFade() = %v at half duration, want in (0, 0.5) -- ease-out sits below the linear ramp", got)
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

func TestView_AuroraBehindContentWhileThinkingAndAbsentAtRest(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	m = newM.(tutorModel)

	if rest := m.View(); strings.Contains(rest, "\x1b[48;2;") {
		t.Error("View() at rest contains background escapes -- the aurora must be completely absent when idle")
	}

	m.turnInFlight = true
	if thinking := m.View(); !strings.Contains(thinking, "\x1b[48;2;") {
		t.Error("View() while thinking has no background escapes -- the aurora isn't rendering")
	}
}
