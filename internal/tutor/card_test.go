package tutor

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestRenderCodeCard_ExactFrameShape(t *testing.T) {
	lines := renderCodeCard("python", []string{"nums = [1, 2, 2]", "unique = set(nums)"}, 40, false)

	if len(lines) != 5 {
		t.Fatalf("card has %d lines, want 5 (top border, header, 2 code lines, bottom border):\n%s", len(lines), strings.Join(lines, "\n"))
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 40 {
			t.Errorf("card line %d is %d cells wide, want exactly 40: %q", i, w, stripAnsiTest(line))
		}
	}

	plain := make([]string, len(lines))
	for i, line := range lines {
		plain[i] = stripAnsiTest(line)
	}
	if !strings.HasPrefix(plain[0], "╭") || !strings.HasSuffix(plain[0], "╮") {
		t.Errorf("top border = %q, want a rounded ╭...╮ rule", plain[0])
	}
	if !strings.Contains(plain[1], "● ● ●") || !strings.Contains(plain[1], "python") {
		t.Errorf("header = %q, want the three traffic-light dots and the language label", plain[1])
	}
	if !strings.Contains(plain[2], "1") || !strings.Contains(plain[2], "nums = [1, 2, 2]") {
		t.Errorf("first code line = %q, want gutter number 1 and the code", plain[2])
	}
	if !strings.Contains(plain[3], "2") || !strings.Contains(plain[3], "unique = set(nums)") {
		t.Errorf("second code line = %q, want gutter number 2 and the code", plain[3])
	}
	if !strings.HasPrefix(plain[4], "╰") || !strings.HasSuffix(plain[4], "╯") {
		t.Errorf("bottom border = %q, want a rounded ╰...╯ rule", plain[4])
	}
}

func TestRenderCodeCard_UnterminatedRendersBottomless(t *testing.T) {
	lines := renderCodeCard("go", []string{"x := 1"}, 40, true)

	if len(lines) != 3 {
		t.Fatalf("bottomless card has %d lines, want 3 (top, header, 1 code line -- no bottom border):\n%s", len(lines), strings.Join(lines, "\n"))
	}
	last := stripAnsiTest(lines[len(lines)-1])
	if strings.Contains(last, "╰") || strings.Contains(last, "╯") {
		t.Errorf("last line = %q, want no closing border while the fence is still streaming", last)
	}
}

func TestRenderCodeCard_LongLinesHardTruncateInsideTheFrame(t *testing.T) {
	long := strings.Repeat("abcdef ", 20) // ~140 cells, far wider than the card
	lines := renderCodeCard("python", []string{long}, 40, false)

	for i, line := range lines {
		if w := lipgloss.Width(line); w != 40 {
			t.Errorf("card line %d is %d cells wide, want exactly 40 -- long code must truncate, never wrap inside the gutter", i, w)
		}
	}
	code := stripAnsiTest(lines[2])
	if !strings.Contains(code, "…") {
		t.Errorf("truncated code line = %q, want a … marker showing content was cut", code)
	}
}

func TestRenderCodeCard_EmptyLabelSaysCode(t *testing.T) {
	lines := renderCodeCard("", []string{"plain"}, 40, false)
	if !strings.Contains(stripAnsiTest(lines[1]), "code") {
		t.Errorf("header = %q, want the generic \"code\" label for an unlabeled fence", stripAnsiTest(lines[1]))
	}
}

func TestRenderCodeCard_TooNarrowDegradesToFlatCode(t *testing.T) {
	// Below the minimum usable card width there's no room for borders +
	// gutter + content -- fall back to the flat accent-colored lines
	// rather than rendering a broken frame.
	lines := renderCodeCard("python", []string{"x = 1"}, 8, false)
	joined := stripAnsiTest(strings.Join(lines, "\n"))
	if strings.Contains(joined, "╭") {
		t.Errorf("got a framed card at width 8:\n%s\nwant the flat fallback", joined)
	}
	if !strings.Contains(joined, "x = 1") {
		t.Errorf("fallback lost the code itself: %q", joined)
	}
}

func TestStyleMarkdown_FencesBecomeCardsAtRealWidths(t *testing.T) {
	in := "look:\n```python\nx = 1\n```\ndone"
	got := styleMarkdown(in, 60)

	if !strings.Contains(got, "╭") || !strings.Contains(got, "╯") {
		t.Fatalf("no card frame in output:\n%s", stripAnsiTest(got))
	}
	plain := stripAnsiTest(got)
	if !strings.Contains(plain, "look:") || !strings.Contains(plain, "done") {
		t.Errorf("prose around the fence lost:\n%s", plain)
	}
	for _, line := range strings.Split(got, "\n") {
		if w := lipgloss.Width(line); w > 60 {
			t.Errorf("output line %q is %d cells wide, want within 60", stripAnsiTest(line), w)
		}
	}
}

func TestStyleMarkdown_WidthZeroKeepsRuleStyleFences(t *testing.T) {
	in := "```python\nx = 1\n```"
	got := styleMarkdown(in, 0)

	if strings.Contains(got, "╭") {
		t.Errorf("width<=0 must keep the width-independent rule-style fences, got a card:\n%s", stripAnsiTest(got))
	}
	if !strings.Contains(stripAnsiTest(got), "── python ──") {
		t.Errorf("rule-style fence label missing:\n%s", stripAnsiTest(got))
	}
}

func TestStyleMarkdown_UnterminatedFenceStreamsAsBottomlessCard(t *testing.T) {
	in := "here is code:\n```python\nx = 1" // cut mid-stream, no closing fence
	got := styleMarkdown(in, 60)

	if !strings.Contains(got, "╭") {
		t.Fatalf("no card frame for the unterminated fence:\n%s", stripAnsiTest(got))
	}
	if strings.Contains(got, "╰") {
		t.Errorf("unterminated fence must render bottomless (no ╰ yet):\n%s", stripAnsiTest(got))
	}
}

// TestTutorModel_ResizeReRendersCardsAtNewWidth pins the whole point of
// the render-at-width pipeline: a card built for one pane width must be
// rebuilt, not re-wrapped, when the pane resizes.
func TestTutorModel_ResizeReRendersCardsAtNewWidth(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m.displayBlocks = []displayBlock{{kind: blockTutor, raw: "```python\nx = 1\n```"}}
	m.refreshViewport()

	wide := m.viewport.View()
	if !strings.Contains(stripAnsiTest(wide), "╭") {
		t.Fatalf("no card at width 80:\n%s", stripAnsiTest(wide))
	}
	wideTop := cardTopWidth(t, wide)

	newM, _ = m.Update(tea.WindowSizeMsg{Width: 50, Height: 24})
	m = newM.(tutorModel)
	narrow := m.viewport.View()
	narrowTop := cardTopWidth(t, narrow)

	if narrowTop >= wideTop {
		t.Errorf("card top border is %d cells at width 50 vs %d at width 80 -- resize must rebuild the card narrower", narrowTop, wideTop)
	}
	for _, line := range strings.Split(narrow, "\n") {
		if w := lipgloss.Width(line); w > 50 {
			t.Errorf("post-resize line %q is %d cells wide, want within 50", stripAnsiTest(line), w)
		}
	}
}

// cardTopWidth finds the card's top border row in a rendered view and
// returns its visible width.
func cardTopWidth(t *testing.T, view string) int {
	t.Helper()
	for _, line := range strings.Split(view, "\n") {
		plain := stripAnsiTest(line)
		if idx := strings.Index(plain, "╭"); idx >= 0 {
			end := strings.LastIndex(plain, "╮")
			if end <= idx {
				t.Fatalf("card top border has no closing corner: %q", plain)
			}
			return lipgloss.Width(strings.TrimSpace(plain))
		}
	}
	t.Fatalf("no card top border found in view:\n%s", stripAnsiTest(view))
	return 0
}
