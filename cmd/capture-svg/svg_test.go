package main

import (
	"strings"
	"testing"
)

func TestApplySGR_Truecolor(t *testing.T) {
	s := applySGR(newStyle(), "38;2;47;166;166")
	if s.fg != "#2fa6a6" {
		t.Errorf("fg = %s, want #2fa6a6", s.fg)
	}
	s = applySGR(newStyle(), "48;2;20;21;28")
	if s.bg != "#14151c" {
		t.Errorf("bg = %s, want #14151c", s.bg)
	}
}

func TestApplySGR_256Color(t *testing.T) {
	// 234 is in the grayscale ramp: 8 + (234-232)*10 = 28.
	if s := applySGR(newStyle(), "48;5;234"); s.bg != "#1c1c1c" {
		t.Errorf("bg = %s, want #1c1c1c", s.bg)
	}
	// 196 sits in the color cube.
	if s := applySGR(newStyle(), "38;5;196"); s.fg != "#ff0000" {
		t.Errorf("fg = %s, want #ff0000", s.fg)
	}
}

func TestApplySGR_BasicAndBright(t *testing.T) {
	if s := applySGR(newStyle(), "31"); s.fg != base16[1] {
		t.Errorf("fg = %s, want %s", s.fg, base16[1])
	}
	if s := applySGR(newStyle(), "91"); s.fg != base16[9] {
		t.Errorf("bright fg = %s, want %s", s.fg, base16[9])
	}
}

func TestApplySGR_Attributes(t *testing.T) {
	s := applySGR(newStyle(), "1;2;3;4;7")
	if !s.bold || !s.dim || !s.italic || !s.underline || !s.reverse {
		t.Errorf("attributes not all set: %+v", s)
	}
	s = applySGR(s, "22;23;24;27")
	if s.bold || s.dim || s.italic || s.underline || s.reverse {
		t.Errorf("attributes not all cleared: %+v", s)
	}
}

func TestApplySGR_ResetRestoresDefaults(t *testing.T) {
	s := applySGR(newStyle(), "1;38;2;255;0;0")
	s = applySGR(s, "0")
	if s.fg != defaultFg || s.bold {
		t.Errorf("reset left %+v", s)
	}
}

// TestApplySGR_MalformedDegradesRatherThanErroring: a capture is a
// recording of whatever the terminal emitted. Refusing to render one
// because of a single unrecognized parameter would be worse than
// rendering it plainly.
func TestApplySGR_MalformedDegradesRatherThanErroring(t *testing.T) {
	s := applySGR(newStyle(), "38;2;not;a;color")
	if s.fg != defaultFg {
		t.Errorf("fg = %s, want the default left intact", s.fg)
	}
	if s := applySGR(newStyle(), "999;abc"); s.fg != defaultFg {
		t.Errorf("unknown codes changed the style: %+v", s)
	}
}

func TestParse_StripsEscapesAndPadsToWidth(t *testing.T) {
	rows := parse("\x1b[38;2;47;166;166mhi\x1b[0m", 6)
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	if len(rows[0]) != 6 {
		t.Errorf("row has %d cells, want padding to 6", len(rows[0]))
	}
	if rows[0][0].r != 'h' || rows[0][1].r != 'i' {
		t.Errorf("text not parsed: %q%q", rows[0][0].r, rows[0][1].r)
	}
	if rows[0][0].st.fg != "#2fa6a6" {
		t.Errorf("color not applied: %s", rows[0][0].st.fg)
	}
}

func TestParse_MultiByteGlyphsStayIntact(t *testing.T) {
	rows := parse("╔═╗", 3)
	got := string([]rune{rows[0][0].r, rows[0][1].r, rows[0][2].r})
	if got != "╔═╗" {
		t.Errorf("got %q, want the box characters intact", got)
	}
}

// TestRender_SpacesArePositionOnly is the rule that keeps the output
// aligned: a run including spaces would stretch its textLength over
// them, and any viewer whose font metrics differ slightly then spreads
// that error across the run as visibly uneven spacing.
func TestRender_SpacesArePositionOnly(t *testing.T) {
	out := render(parse("ab   cd", 7), 7)
	if strings.Contains(out, ">ab   cd<") {
		t.Error("spaces were rendered inside a text run")
	}
	if !strings.Contains(out, ">ab<") || !strings.Contains(out, ">cd<") {
		t.Errorf("expected two separate runs:\n%s", out)
	}
}

func TestRender_EmitsBackgroundRectsOnlyForNonDefault(t *testing.T) {
	plain := render(parse("hello", 5), 5)
	// The page rect is always present; no run rects should join it.
	if strings.Count(plain, "<rect") != 1 {
		t.Errorf("plain text emitted extra background rects:\n%s", plain)
	}
	styled := render(parse("\x1b[48;2;30;32;41mhi\x1b[0m", 2), 2)
	if strings.Count(styled, "<rect") < 2 {
		t.Errorf("styled background produced no rect:\n%s", styled)
	}
}

func TestRender_EscapesXML(t *testing.T) {
	out := render(parse("a<b>&c", 6), 6)
	if !strings.Contains(out, "a&lt;b&gt;&amp;c") {
		t.Errorf("XML not escaped:\n%s", out)
	}
}

func TestRender_MatchesCommittedStructure(t *testing.T) {
	out := render(parse("\x1b[38;2;155;95;176m│\x1b[0m", 1), 1)
	for _, want := range []string{
		`<svg xmlns="http://www.w3.org/2000/svg"`,
		`lengthAdjust="spacingAndGlyphs"`,
		`font-family="SF Mono, Menlo, Consolas, monospace"`,
		`fill="#9b5fb0"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q — it must reproduce the committed SVG shape:\n%s", want, out)
		}
	}
}

func TestRender_AttributesMapToSVG(t *testing.T) {
	out := render(parse("\x1b[1mB\x1b[0m", 1), 1)
	if !strings.Contains(out, `font-weight="600"`) {
		t.Error("bold did not map to font-weight")
	}
	out = render(parse("\x1b[2mD\x1b[0m", 1), 1)
	if !strings.Contains(out, `opacity="0.62"`) {
		t.Error("dim did not map to opacity")
	}
}

func TestRender_ReverseSwapsColors(t *testing.T) {
	out := render(parse("\x1b[38;2;255;0;0;48;2;0;0;255;7mX\x1b[0m", 1), 1)
	if !strings.Contains(out, `fill="#0000ff"`) {
		t.Errorf("reverse did not swap foreground to the background color:\n%s", out)
	}
}
