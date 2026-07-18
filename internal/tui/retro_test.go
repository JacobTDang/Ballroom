package tui

import (
	"strings"
	"testing"
)

// TestHeading_UppercaseAndLetterspaced pins the retro heading style.
// The precedent is already in the banner's tagline ("I N T E R V I E W
// P R E P"), so this generalizes what the app was doing in one place.
func TestHeading_UppercaseAndLetterspaced(t *testing.T) {
	if got, want := heading("Stats"), "S T A T S"; got != want {
		t.Errorf("heading(%q) = %q, want %q", "Stats", got, want)
	}
}

func TestHeading_PreservesWordGaps(t *testing.T) {
	// Two words stay distinguishable: letters get one space between
	// them, words get three, or "NEXT UP" reads as one long word.
	got := heading("Next up")
	if !strings.Contains(got, "N E X T") || !strings.Contains(got, "U P") {
		t.Errorf("heading(%q) = %q, want both words letterspaced", "Next up", got)
	}
	if strings.Contains(got, "T U") {
		t.Errorf("heading(%q) = %q, word boundary is not distinguishable", "Next up", got)
	}
}

func TestHeading_HandlesMultiByte(t *testing.T) {
	// Letterspacing operates on runes, not bytes -- a naive byte loop
	// would split a multi-byte rune and emit garbage.
	got := heading("café")
	if !strings.Contains(got, "É") && !strings.Contains(got, "É") {
		t.Errorf("heading(%q) = %q, want the accented rune intact and uppercased", "café", got)
	}
}

func TestHeading_EmptyIsEmpty(t *testing.T) {
	if got := heading(""); got != "" {
		t.Errorf("heading(\"\") = %q, want empty", got)
	}
}

// TestProgressBar_BracketedAndBlockShaded pins the retro bar shape.
func TestProgressBar_BracketedAndBlockShaded(t *testing.T) {
	got := stripAnsiTUI(progressBar(5, 10, 10))
	if !strings.HasPrefix(got, "[") || !strings.HasSuffix(got, "]") {
		t.Errorf("progressBar = %q, want bracketed framing", got)
	}
	if !strings.Contains(got, "█") {
		t.Errorf("progressBar = %q, want block-shaded fill", got)
	}
}

// TestProgressBar_KeepsTheFloorToOneRule is the behavior that must
// survive the restyle: one solved problem out of 149 still shows a
// visible cell rather than rounding to an empty bar.
func TestProgressBar_KeepsTheFloorToOneRule(t *testing.T) {
	got := stripAnsiTUI(progressBar(1, 149, 14))
	if !strings.Contains(got, "█") {
		t.Errorf("progressBar(1, 149, 14) = %q, want at least one filled cell", got)
	}
}

func TestProgressBar_ZeroAndFull(t *testing.T) {
	zero := stripAnsiTUI(progressBar(0, 10, 8))
	if strings.Contains(zero, "█") {
		t.Errorf("progressBar(0,...) = %q, want no filled cells", zero)
	}
	full := stripAnsiTUI(progressBar(10, 10, 8))
	if strings.Contains(full, "░") {
		t.Errorf("progressBar(10,10,...) = %q, want no empty cells", full)
	}
}

func TestProgressBar_ZeroTotalDoesNotDivideByZero(t *testing.T) {
	if got := stripAnsiTUI(progressBar(3, 0, 6)); !strings.HasPrefix(got, "[") {
		t.Errorf("progressBar with zero total = %q, want a rendered empty bar", got)
	}
}

// TestPanelUsesDoubleBorder: the double line is the single most
// recognizable retro-terminal cue, and the panel frames every screen.
func TestPanelUsesDoubleBorder(t *testing.T) {
	m := appModel{stage: stageMain, width: 120, height: 40}
	out := m.View()
	for _, corner := range []string{"╔", "╗", "╚", "╝"} {
		if !strings.Contains(out, corner) {
			t.Errorf("View missing double-border corner %q", corner)
		}
	}
}
