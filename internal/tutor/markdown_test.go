package tutor

import (
	"strings"
	"testing"
)

// styleMarkdown is display-only styling for the chat pane -- the raw
// reply text the model sees in history must never carry these escapes,
// which is why it's applied at the displayLines append, not before.

func TestStyleMarkdown_BoldRendersWithoutMarkers(t *testing.T) {
	got := styleMarkdown("return true if any value appears **at least twice** in the array")
	if strings.Contains(got, "**") {
		t.Errorf("styleMarkdown left raw ** markers in %q", got)
	}
	if !strings.Contains(got, "\x1b[1m") {
		t.Errorf("styleMarkdown output %q has no bold escape", got)
	}
	if !strings.Contains(stripAnsiTest(got), "at least twice") {
		t.Errorf("styleMarkdown lost the bold text itself: %q", got)
	}
}

func TestStyleMarkdown_InlineCodeStyledWithoutBackticks(t *testing.T) {
	got := styleMarkdown("use `set()` to track seen values")
	if strings.Contains(got, "`") {
		t.Errorf("styleMarkdown left raw backticks in %q", got)
	}
	if !strings.Contains(got, "\x1b[38;2;") {
		t.Errorf("styleMarkdown output %q has no color escape for the code span", got)
	}
	if !strings.Contains(stripAnsiTest(got), "set()") {
		t.Errorf("styleMarkdown lost the code text itself: %q", got)
	}
}

func TestStyleMarkdown_FencedBlockStyledAndFencesReplaced(t *testing.T) {
	in := "Here you go:\n```python\ndef f():\n    return 1\n```\nDone."
	got := styleMarkdown(in)
	if strings.Contains(got, "```") {
		t.Errorf("styleMarkdown left raw fence markers in %q", got)
	}
	stripped := stripAnsiTest(got)
	for _, want := range []string{"def f():", "return 1", "Here you go:", "Done.", "python"} {
		if !strings.Contains(stripped, want) {
			t.Errorf("styleMarkdown output missing %q:\n%s", want, stripped)
		}
	}
}

func TestStyleMarkdown_NoInlineStylingInsideFences(t *testing.T) {
	in := "```\nx = a ** b  # `power` operator\n```"
	got := stripAnsiTest(styleMarkdown(in))
	if !strings.Contains(got, "a ** b") {
		t.Errorf("styleMarkdown mangled ** inside a code fence: %q", got)
	}
	if !strings.Contains(got, "`power`") {
		t.Errorf("styleMarkdown mangled backticks inside a code fence: %q", got)
	}
}

func TestStyleMarkdown_HeadersBold(t *testing.T) {
	got := styleMarkdown("## Approach\nuse a set")
	if strings.Contains(stripAnsiTest(got), "##") {
		t.Errorf("styleMarkdown left raw ## marker in %q", got)
	}
	if !strings.Contains(got, "\x1b[1m") {
		t.Errorf("styleMarkdown output %q has no bold escape for the header", got)
	}
}

func TestStyleMarkdown_PlainTextPassesThroughUntouched(t *testing.T) {
	in := "just a plain sentence with no markdown at all"
	if got := styleMarkdown(in); got != in {
		t.Errorf("styleMarkdown(%q) = %q, want unchanged", in, got)
	}
}

func TestStyleMarkdown_StrayAsteriskNotMangled(t *testing.T) {
	in := "the result is 2 * 3 and 4 * 5"
	if got := styleMarkdown(in); got != in {
		t.Errorf("styleMarkdown(%q) = %q, want unchanged -- lone asterisks aren't bold markers", in, got)
	}
}

// stripAnsiTest removes ANSI escapes for content assertions.
func stripAnsiTest(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		if r == '\x1b' {
			inEsc = true
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
