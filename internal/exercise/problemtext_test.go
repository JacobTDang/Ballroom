package exercise

import (
	"strings"
	"testing"
)

// RenderProblemText converts an exercise's problem.md into the clean
// structured plain text the editor pane actually displays -- readable
// as-is, no markdown markers, no terminal needed to interpret anything.

func TestRenderProblemText_TitleGetsUnderlined(t *testing.T) {
	got := RenderProblemText("# Contains Duplicate\n\nbody text")
	if strings.Contains(got, "#") {
		t.Errorf("rendered text still contains # markers:\n%s", got)
	}
	lines := strings.Split(got, "\n")
	if lines[0] != "Contains Duplicate" {
		t.Fatalf("first line = %q, want the bare title", lines[0])
	}
	if lines[1] != strings.Repeat("═", len("Contains Duplicate")) {
		t.Errorf("second line = %q, want a ═ underline matching the title length", lines[1])
	}
}

func TestRenderProblemText_SectionHeaderGetsLighterUnderline(t *testing.T) {
	got := RenderProblemText("## Examples\n\nstuff")
	if strings.Contains(got, "#") {
		t.Errorf("rendered text still contains # markers:\n%s", got)
	}
	lines := strings.Split(got, "\n")
	if lines[0] != "Examples" || lines[1] != strings.Repeat("─", len("Examples")) {
		t.Errorf("section header rendered as %q / %q, want bare text over a ─ underline", lines[0], lines[1])
	}
}

func TestRenderProblemText_BoldMarkersStripped(t *testing.T) {
	got := RenderProblemText("appears **at least twice** in the array")
	if strings.Contains(got, "*") {
		t.Errorf("rendered text still contains * markers: %q", got)
	}
	if !strings.Contains(got, "at least twice") {
		t.Errorf("bold text lost: %q", got)
	}
}

func TestRenderProblemText_InlineCodeBackticksStripped(t *testing.T) {
	got := RenderProblemText("return `true` if `nums` has a duplicate")
	if strings.Contains(got, "`") {
		t.Errorf("rendered text still contains backticks: %q", got)
	}
	if !strings.Contains(got, "return true if nums has a duplicate") {
		t.Errorf("inline code text mangled: %q", got)
	}
}

func TestRenderProblemText_FencedBlockIndentedWithMarkersDropped(t *testing.T) {
	in := "## Examples\n\n```\nInput: nums = [1,2,3,1]\nOutput: true\n```\n"
	got := RenderProblemText(in)
	if strings.Contains(got, "```") {
		t.Errorf("rendered text still contains fence markers:\n%s", got)
	}
	if !strings.Contains(got, "  Input: nums = [1,2,3,1]") || !strings.Contains(got, "  Output: true") {
		t.Errorf("fenced content not indented as a block:\n%s", got)
	}
}

func TestRenderProblemText_NoTransformsInsideFences(t *testing.T) {
	in := "```\nx = a ** b  # `power`\n```"
	got := RenderProblemText(in)
	if !strings.Contains(got, "a ** b") || !strings.Contains(got, "`power`") {
		t.Errorf("code inside a fence was mangled:\n%s", got)
	}
}

func TestRenderProblemText_BoldSpanningAHardWrappedLineBreakStripped(t *testing.T) {
	// Real case from contains-duplicate's problem.md: authors hard-wrap
	// prose, so a bold span can open on one line and close on the next.
	in := "return true if any value appears **at\nleast twice** in the array"
	got := RenderProblemText(in)
	if strings.Contains(got, "*") {
		t.Errorf("rendered text still contains * markers across a wrapped bold span:\n%s", got)
	}
	if !strings.Contains(got, "at\nleast twice") {
		t.Errorf("wrapped bold text lost or re-flowed: %q", got)
	}
}

func TestRenderProblemText_InlineCodeSpanningAHardWrappedLineBreakStripped(t *testing.T) {
	// Real case from greedy-06's problem.md: `target =\n[x, y, z]`.
	in := "an integer array `target =\n[x, y, z]` that describes the triplet"
	got := RenderProblemText(in)
	if strings.Contains(got, "`") {
		t.Errorf("rendered text still contains backticks across a wrapped code span:\n%s", got)
	}
	if !strings.Contains(got, "target =\n[x, y, z]") {
		t.Errorf("wrapped code text lost or re-flowed: %q", got)
	}
}

func TestRenderProblemText_HeaderContainingInlineCodeStripped(t *testing.T) {
	// Real case from off-by-one-01's problem.md: # Off-by-one: `MaxOf`
	got := RenderProblemText("# Off-by-one: `MaxOf`\n\nbody")
	if strings.Contains(got, "`") {
		t.Errorf("rendered header still contains backticks:\n%s", got)
	}
	lines := strings.Split(got, "\n")
	if lines[0] != "Off-by-one: MaxOf" {
		t.Fatalf("header line = %q, want the backticks stripped", lines[0])
	}
	if lines[1] != strings.Repeat("═", len("Off-by-one: MaxOf")) {
		t.Errorf("underline = %q, want its length to match the STRIPPED header text", lines[1])
	}
}

func TestRenderProblemText_StrayAsterisksAndPlainTextUntouched(t *testing.T) {
	in := "the result is 2 * 3 and 4 * 5"
	if got := RenderProblemText(in); got != in {
		t.Errorf("RenderProblemText(%q) = %q, want unchanged", in, got)
	}
}
