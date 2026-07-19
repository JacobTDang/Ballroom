package tutor

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// styleMarkdown is display-only styling for the chat pane -- the raw
// reply text the model sees in history must never carry these escapes,
// which is why it's applied at the displayLines append, not before.

func TestStyleMarkdown_BoldRendersWithoutMarkers(t *testing.T) {
	got := styleMarkdown("return true if any value appears **at least twice** in the array", 0)
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
	got := styleMarkdown("use `set()` to track seen values", 0)
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
	got := styleMarkdown(in, 0)
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
	got := stripAnsiTest(styleMarkdown(in, 0))
	if !strings.Contains(got, "a ** b") {
		t.Errorf("styleMarkdown mangled ** inside a code fence: %q", got)
	}
	if !strings.Contains(got, "`power`") {
		t.Errorf("styleMarkdown mangled backticks inside a code fence: %q", got)
	}
}

func TestStyleMarkdown_HeadersBold(t *testing.T) {
	got := styleMarkdown("## Approach\nuse a set", 0)
	if strings.Contains(stripAnsiTest(got), "##") {
		t.Errorf("styleMarkdown left raw ## marker in %q", got)
	}
	if !strings.Contains(got, "\x1b[1m") {
		t.Errorf("styleMarkdown output %q has no bold escape for the header", got)
	}
}

func TestStyleMarkdown_PlainTextPassesThroughUntouched(t *testing.T) {
	in := "just a plain sentence with no markdown at all"
	if got := styleMarkdown(in, 0); got != in {
		t.Errorf("styleMarkdown(%q, 0) = %q, want unchanged", in, got)
	}
}

func TestStyleMarkdown_StrayAsteriskNotMangled(t *testing.T) {
	in := "the result is 2 * 3 and 4 * 5"
	if got := styleMarkdown(in, 0); got != in {
		t.Errorf("styleMarkdown(%q, 0) = %q, want unchanged -- lone asterisks aren't bold markers", in, got)
	}
}

func TestStyleMarkdown_KnownLanguageFenceGetsTokenColors(t *testing.T) {
	in := "```python\ndef add(a, b):\n    return a + b\n```"
	got := styleMarkdown(in, 0)
	if strings.Contains(got, "```") {
		t.Errorf("fence markers leaked:\n%s", got)
	}
	stripped := stripAnsiTest(got)
	if !strings.Contains(stripped, "def add(a, b):") || !strings.Contains(stripped, "return a + b") {
		t.Errorf("code content lost:\n%s", stripped)
	}
	// Real token coloring means the keyword and the identifier carry
	// DIFFERENT colors -- a flat single-color block has exactly one
	// distinct foreground sequence.
	colors := distinctForegrounds(got)
	if len(colors) < 2 {
		t.Errorf("highlighted python block has %d distinct foreground colors, want >= 2 (token-level highlighting)", len(colors))
	}
}

func TestStyleMarkdown_UnknownLanguageFenceFallsBackToFlatColor(t *testing.T) {
	in := "```notareallang\nblorp blip 42\n```"
	got := styleMarkdown(in, 0)
	if !strings.Contains(got, mdCodeColor) {
		t.Errorf("unknown-language fence should keep the flat accent color:\n%q", got)
	}
	if !strings.Contains(stripAnsiTest(got), "blorp blip 42") {
		t.Errorf("code content lost:\n%s", stripAnsiTest(got))
	}
}

// TestStyleMarkdown_BulletListStyledMarkerAndHangIndent: a long item
// must wrap HERE with a hang indent, not in refreshViewport's outer
// pass — the outer wrap is indent-blind and would fold continuations
// flush-left.
func TestStyleMarkdown_BulletListStyledMarkerAndHangIndent(t *testing.T) {
	got := styleMarkdown("- alpha beta gamma delta epsilon zeta eta theta iota kappa lambda", 40)
	if !strings.Contains(got, mdCodeColor+"-"+mdColorReset) {
		t.Errorf("bullet marker not styled in the accent color: %q", got)
	}
	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Fatalf("long bullet item did not wrap: %q", got)
	}
	if first := stripAnsiTest(lines[0]); !strings.HasPrefix(first, "- alpha") {
		t.Errorf("first line = %q, want it to start with the marker and item text", first)
	}
	for i, line := range lines[1:] {
		if plain := stripAnsiTest(line); !strings.HasPrefix(plain, "  ") {
			t.Errorf("continuation line %d = %q, want it hang-indented under the item text", i+1, plain)
		}
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w > 40 {
			t.Errorf("list line %d is %d cells wide, want <= 40 so the outer wrap never re-breaks it", i, w)
		}
	}
}

func TestStyleMarkdown_NestedBulletKeepsIndent(t *testing.T) {
	got := styleMarkdown("  - nested item", 0)
	if plain := stripAnsiTest(got); plain != "  - nested item" {
		t.Errorf("stripped = %q, want the nested indent and text preserved", plain)
	}
	if !strings.Contains(got, mdCodeColor+"-"+mdColorReset) {
		t.Errorf("nested bullet marker not styled: %q", got)
	}
}

func TestStyleMarkdown_OrderedListDimNumberAndHangIndent(t *testing.T) {
	got := styleMarkdown("12. alpha beta gamma delta epsilon zeta eta theta iota kappa", 40)
	if !strings.Contains(got, mdDimColor+"12."+mdColorReset) {
		t.Errorf("ordered marker not dimmed: %q", got)
	}
	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Fatalf("long ordered item did not wrap: %q", got)
	}
	if first := stripAnsiTest(lines[0]); !strings.HasPrefix(first, "12. alpha") {
		t.Errorf("first line = %q, want marker then text", first)
	}
	if cont := stripAnsiTest(lines[1]); !strings.HasPrefix(cont, "    ") {
		t.Errorf("continuation = %q, want 4-cell hang indent aligning under the text", cont)
	}
}

func TestStyleMarkdown_BlockquoteBarAndDimText(t *testing.T) {
	got := styleMarkdown("> stay with the invariant", 0)
	plain := stripAnsiTest(got)
	if plain != "│ stay with the invariant" {
		t.Errorf("stripped = %q, want the bar replacing the > marker", plain)
	}
	if !strings.Contains(got, ansiFg(paneRule)+"│") {
		t.Errorf("quote bar missing the structural rule color: %q", got)
	}
	if !strings.Contains(got, mdDimColor) {
		t.Errorf("quote text not dimmed: %q", got)
	}

	if bare := stripAnsiTest(styleMarkdown(">", 0)); !strings.HasPrefix(bare, "│") {
		t.Errorf("bare > line = %q, want just the bar", bare)
	}
}

// TestStyleMarkdown_BlockquoteInlineCodeStaysDimAfterSpan: an inline
// code span inside a quote closes with a default-foreground reset —
// the quote must re-arm its dim color after it or the rest of the
// line "leaks" back to full brightness (same re-arm trick the editor
// cards use on chroma's resets).
func TestStyleMarkdown_BlockquoteInlineCodeStaysDimAfterSpan(t *testing.T) {
	got := styleMarkdown("> use `seen` before the loop ends", 0)
	if !strings.Contains(got, mdCodeColor+"seen"+mdDimColor) {
		t.Errorf("code span inside a quote must close back into dim, got %q", got)
	}
	if !strings.Contains(stripAnsiTest(got), "before the loop ends") {
		t.Errorf("quote text after the span lost: %q", got)
	}
}

// TestStyleMarkdown_LongBlockquoteKeepsBarOnEveryWrappedLine: like
// list items, a quote must wrap itself — the outer pass would fold a
// long quote flush-left and orphan the continuation from its bar
// (seen live in the preview harness at 44 cols).
func TestStyleMarkdown_LongBlockquoteKeepsBarOnEveryWrappedLine(t *testing.T) {
	got := styleMarkdown("> the invariant is that everything left of the write pointer is already final and sorted", 40)
	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Fatalf("long quote did not wrap: %q", got)
	}
	for i, line := range lines {
		if plain := stripAnsiTest(line); !strings.HasPrefix(plain, "│ ") {
			t.Errorf("quote line %d = %q, want the bar on every wrapped line", i, plain)
		}
		if w := lipgloss.Width(line); w > 40 {
			t.Errorf("quote line %d is %d cells wide, want <= 40", i, w)
		}
	}
}

func TestStyleMarkdown_HorizontalRuleRendersDoubleDashes(t *testing.T) {
	if got := stripAnsiTest(styleMarkdown("---", 20)); got != strings.Repeat("═", 20) {
		t.Errorf("hr at width 20 = %q, want 20 double-rule cells", got)
	}
	if got := stripAnsiTest(styleMarkdown("***", 0)); got != strings.Repeat("═", 40) {
		t.Errorf("hr at width 0 = %q, want the 40-cell default", got)
	}
	if got := styleMarkdown("-- not a rule", 20); stripAnsiTest(got) != "-- not a rule" {
		t.Errorf("two dashes = %q, want untouched (three+ makes a rule)", got)
	}
}

func TestStyleMarkdown_LinkUnderlinedWithDimURL(t *testing.T) {
	got := styleMarkdown("see [the docs](https://ex.am/p) for more", 0)
	if strings.Contains(stripAnsiTest(got), "[") {
		t.Errorf("raw link brackets left in %q", got)
	}
	if !strings.Contains(got, "\x1b[4m"+"the docs"+"\x1b[24m") {
		t.Errorf("link text not underlined: %q", got)
	}
	if !strings.Contains(got, mdDimColor+" (https://ex.am/p)"+mdColorReset) {
		t.Errorf("link URL not shown dim in parens: %q", got)
	}
}

func TestStyleMarkdown_ImageSyntaxLeftRaw(t *testing.T) {
	in := "diagram: ![flow](https://ex.am/f.png)"
	if got := styleMarkdown(in, 0); got != in {
		t.Errorf("image syntax = %q, want left raw (the pane can't render it)", got)
	}
}

func TestStyleMarkdown_LinkInsideBackticksStaysLiteral(t *testing.T) {
	got := styleMarkdown("the pattern is `[text](url)` exactly", 0)
	if strings.Contains(got, "\x1b[4m") {
		t.Errorf("link inside a code span must stay literal, got underline in %q", got)
	}
	if !strings.Contains(stripAnsiTest(got), "[text](url)") {
		t.Errorf("code span content lost: %q", got)
	}
}

func TestStyleMarkdown_NewConstructsInertInsideFences(t *testing.T) {
	in := "```\n- item\n> quote\n---\n[a](b)\n```"
	got := styleMarkdown(in, 0)
	if strings.Contains(got, "\x1b[4m") {
		t.Errorf("link styling leaked into a fence: %q", got)
	}
	plain := stripAnsiTest(got)
	for _, want := range []string{"- item", "> quote", "---", "[a](b)"} {
		if !strings.Contains(plain, want) {
			t.Errorf("fence content %q mangled:\n%s", want, plain)
		}
	}
}

// TestStyleMarkdown_PartialStreamedListSafe: styleMarkdown runs on the
// in-flight partial reply every streaming tick — a list cut mid-item
// must render, not panic or drop lines.
func TestStyleMarkdown_PartialStreamedListSafe(t *testing.T) {
	got := styleMarkdown("- item one\n- it", 40)
	plain := stripAnsiTest(got)
	if !strings.Contains(plain, "item one") || !strings.HasSuffix(strings.TrimRight(plain, " "), "- it") {
		t.Errorf("partial list = %q, want both lines rendered (trailing pad spaces aside)", plain)
	}
}

func TestStyleMarkdown_EmphasisAtLineStartNotABullet(t *testing.T) {
	in := "*emphasis* is not a bullet"
	if got := styleMarkdown(in, 0); got != in {
		t.Errorf("styleMarkdown(%q) = %q, want unchanged — no space after the star means no list marker", in, got)
	}
}

// distinctForegrounds collects the unique truecolor foreground
// sequences in s.
func distinctForegrounds(s string) map[string]bool {
	out := map[string]bool{}
	for _, m := range regexp.MustCompile(`\x1b\[38;2;\d+;\d+;\d+m`).FindAllString(s, -1) {
		out[m] = true
	}
	return out
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
