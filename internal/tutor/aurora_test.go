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

func TestOverlayAurora_GlowsAtEdgesNotCenter(t *testing.T) {
	const w, h = 60, 21
	got := overlayAurora("hi", w, h, 1.0, auroraBrightness)
	lines := strings.Split(got, "\n")
	if len(lines) != h {
		t.Fatalf("overlayAurora produced %d lines, want %d (content padded to full height)", len(lines), h)
	}

	// The top border row sits fully inside the glow: every cell painted.
	if n := strings.Count(lines[0], "\x1b[48;2;"); n != w {
		t.Errorf("top row has %d background escapes, want %d (the border itself glows end to end)", n, w)
	}

	// The middle row is painted only near its left/right edges -- the
	// glow must dissolve before reaching the center, leaving the
	// conversation area on the terminal's own background.
	mid := lines[h/2]
	nMid := strings.Count(mid, "\x1b[48;2;")
	if nMid == 0 {
		t.Fatal("middle row has no background escapes at all -- the side glow is missing")
	}
	if nMid >= w/2 {
		t.Errorf("middle row has %d background escapes, want well under %d -- the glow is flooding the center instead of fading at the borders", nMid, w/2)
	}
	// The exact center cell must be untouched: unpainted center is what
	// keeps text readable and makes this a border glow, not a wash.
	if centerRun := strings.Contains(ansi.Strip(mid), strings.Repeat(" ", 10)); !centerRun {
		t.Errorf("middle row = %q, want a wide plain-space run at the center", mid)
	}
}

// leftGlowWidth counts the painted cells at the start of a rendered
// blank row -- one background escape per padding cell up to the first
// background reset, which is exactly how deep the left glow reaches.
func leftGlowWidth(line string) int {
	cut := strings.Index(line, "\x1b[49m")
	if cut < 0 {
		cut = len(line)
	}
	return strings.Count(line[:cut], "\x1b[48;2;")
}

func TestOverlayAurora_GlowDepthUndulatesAlongTheBorder(t *testing.T) {
	// A constant-depth glow reads as a rigid box frame (real feedback:
	// "rn its just a box"). The inner boundary must vary along the
	// border -- different rows reach different depths at the same
	// instant.
	const w, h = 60, 25
	lines := strings.Split(overlayAurora("", w, h, 2.0, auroraBrightness), "\n")
	minRun, maxRun := w, 0
	// Rows 8..16 are far enough from the top/bottom glow that their
	// left run is governed purely by horizontal distance.
	for r := 8; r <= 16; r++ {
		run := leftGlowWidth(lines[r])
		if run < minRun {
			minRun = run
		}
		if run > maxRun {
			maxRun = run
		}
	}
	if minRun == maxRun {
		t.Errorf("left glow reaches exactly %d cells on every middle row -- the boundary is a straight box edge, not an undulating one", minRun)
	}
}

func TestOverlayAurora_GlowBoundaryMovesOverTime(t *testing.T) {
	// The undulation must travel, not just exist -- the same row's
	// glow reach should differ a moment later.
	const w, h = 60, 25
	a := strings.Split(overlayAurora("", w, h, 0.0, auroraBrightness), "\n")
	b := strings.Split(overlayAurora("", w, h, 1.5, auroraBrightness), "\n")
	for r := 8; r <= 16; r++ {
		if leftGlowWidth(a[r]) != leftGlowWidth(b[r]) {
			return // boundary moved somewhere in the band -- good
		}
	}
	t.Error("no middle row's glow reach changed over 1.5s -- the boundary undulation isn't traveling")
}

func TestOverlayAurora_ResetsBackgroundLeavingTheGlow(t *testing.T) {
	// Where a painted edge cell is followed by an unpainted interior
	// cell, the background must be explicitly reset -- otherwise the
	// glow color would smear across the whole line through
	// background-color inheritance.
	const w, h = 60, 21
	got := overlayAurora("", w, h, 1.0, auroraBrightness)
	mid := strings.Split(got, "\n")[h/2]
	lastPaint := strings.LastIndex(mid, "\x1b[48;2;")
	lastReset := strings.LastIndex(mid, "\x1b[49m")
	firstPaint := strings.Index(mid, "\x1b[48;2;")
	firstReset := strings.Index(mid, "\x1b[49m")
	if firstReset < firstPaint {
		t.Errorf("middle row resets background before ever painting -- reset placement is wrong: %q", mid)
	}
	if lastReset < lastPaint {
		t.Errorf("middle row's last paint at %d is never followed by a reset (last reset at %d) -- the glow will smear across the line", lastPaint, lastReset)
	}
}

func TestOverlayAurora_PreservesGlyphsAndForegroundStyling(t *testing.T) {
	fg := "\x1b[38;2;10;20;30m"
	content := fg + "hi\x1b[0m"
	got := overlayAurora(content, 6, 1, 1.0, auroraBrightness)
	if !strings.Contains(got, fg) {
		t.Errorf("overlayAurora dropped the content's own foreground escape %q", fg)
	}
	stripped := ansi.Strip(strings.Split(got, "\n")[0])
	if strings.TrimRight(stripped, " ") != "hi" {
		t.Errorf("stripped output = %q, want the original glyphs %q (padded with spaces)", stripped, "hi")
	}
}

func TestOverlayAurora_VisiblyAnimatesWithinASecond(t *testing.T) {
	// The motion has to be perceptible, not technically-nonzero -- the
	// first version drifted so slowly it read as a static image (real
	// feedback: "the aurora should be moving"). One second apart, the
	// top border row must render differently.
	a := strings.Split(overlayAurora("hi", 60, 12, 0.0, auroraBrightness), "\n")[0]
	b := strings.Split(overlayAurora("hi", 60, 12, 1.0, auroraBrightness), "\n")[0]
	if a == b {
		t.Error("top border row identical 1 second apart -- the glow isn't visibly moving")
	}
}

func TestOverlayAurora_KeepsContentTallerThanPane(t *testing.T) {
	content := "1\n2\n3\n4\n5"
	got := overlayAurora(content, 4, 3, 1.0, auroraBrightness)
	if n := len(strings.Split(got, "\n")); n != 5 {
		t.Errorf("overlayAurora produced %d lines for 5-line content, want 5 -- a background paints behind content, it never clips it", n)
	}
}

// --- auroraFade -- lifecycle ---

func TestAuroraFade_FullOnceRampCompletes(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnInFlight = true
	m.turnStartedAt = time.Now().Add(-10 * time.Second)
	m.turnSettledAt = time.Now().Add(-10 * time.Second) // stale, must not matter
	if got := m.auroraFade(); got != 1.0 {
		t.Errorf("auroraFade() = %v long after the turn started, want 1.0", got)
	}
}

func TestAuroraFade_RampsInGentlyAtTurnStart(t *testing.T) {
	// The glow must bloom in, not pop to full strength the instant the
	// user hits enter (real feedback: "it should[n't] start so strong
	// ... come in slowly and more naturally").
	m := newTutorLayoutOnly()
	m.turnInFlight = true

	m.turnStartedAt = time.Now()
	if got := m.auroraFade(); got >= 0.1 {
		t.Errorf("auroraFade() = %v at the instant a turn starts, want near 0 -- the glow pops instead of blooming", got)
	}

	m.turnStartedAt = time.Now().Add(-auroraFadeInDuration / 2)
	mid := m.auroraFade()
	if mid <= 0.3 || mid >= 0.7 {
		t.Errorf("auroraFade() = %v at half the ramp, want around 0.5 (smoothstep midpoint)", mid)
	}

	m.turnStartedAt = time.Now().Add(-auroraFadeInDuration)
	if got := m.auroraFade(); got != 1.0 {
		t.Errorf("auroraFade() = %v at exactly the ramp duration, want 1.0", got)
	}
}

func TestAuroraFadeInLead_ResumesTheRampFromAGivenLevel(t *testing.T) {
	// When a new turn starts while the previous glow is still fading
	// out, turnStartedAt is backdated by auroraFadeInLead(level) so the
	// bloom continues smoothly from the current level instead of
	// blinking down to zero. The lead must therefore invert the ramp.
	for _, level := range []float64{0.1, 0.25, 0.5, 0.75, 0.9} {
		m := newTutorLayoutOnly()
		m.turnInFlight = true
		m.turnStartedAt = time.Now().Add(-auroraFadeInLead(level))
		got := m.auroraFade()
		if got < level-0.05 || got > level+0.05 {
			t.Errorf("auroraFade() = %v with turnStartedAt backdated for level %v, want within 0.05", got, level)
		}
	}
}

func TestAuroraFadeOutLead_ResumesTheDecayFromAGivenLevel(t *testing.T) {
	// Mirror of the fade-in lead: if the reply lands mid-bloom, the
	// fade-out must start from the bloom's current level, not jump up
	// to full brightness first.
	for _, level := range []float64{0.1, 0.25, 0.5, 0.75, 0.9} {
		m := newTutorLayoutOnly()
		m.turnSettledAt = time.Now().Add(-auroraFadeOutLead(level))
		got := m.auroraFade()
		if got < level-0.05 || got > level+0.05 {
			t.Errorf("auroraFade() = %v with turnSettledAt backdated for level %v, want within 0.05", got, level)
		}
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
	m.turnSettledAt = time.Now().Add(-auroraFadeOutDuration)
	if got := m.auroraFade(); got != 0.0 {
		t.Errorf("auroraFade() = %v at exactly fadeDuration elapsed, want 0.0 -- fully gone, not asymptotically dim", got)
	}
}

func TestAuroraFade_MidFadeEaseOutBelowLinear(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnSettledAt = time.Now().Add(-auroraFadeOutDuration / 2)
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
