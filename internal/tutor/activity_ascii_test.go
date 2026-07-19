package tutor

import (
	"testing"
	"time"
	"unicode"
)

// TestActivityRenderingIsPlainASCII pins the hard constraint this file's
// own header comments describe at length: every rendering function in
// activity.go must only ever emit runes in the ASCII range. This project
// has been burned by "safe-looking" Unicode glyphs twice already -- a
// per-state symbol set (⟳ ✓ ✗) plus a Unicode ellipsis (…) rendered as
// tofu in a real user's terminal font, and "●" (U+25CF), picked as the
// fix, later turned out to render as a bare underscore for a different
// user. Plain ASCII is the only thing this project can still promise
// renders identically everywhere, so this test exercises every
// string-producing function in this file across a spread of inputs
// (status, phase, width, call content) and fails the instant any of them
// emits a rune outside the ASCII range. The escape-code plumbing
// (\x1b[38;2;r;g;bm etc.) needs no special-casing or stripping here --
// every byte of an ANSI escape is itself plain ASCII.
func TestActivityRenderingIsPlainASCII(t *testing.T) {
	assertASCII := func(t *testing.T, label, s string) {
		t.Helper()
		for _, r := range s {
			if r > unicode.MaxASCII {
				t.Errorf("%s emitted non-ASCII rune %q (%U) in %q", label, r, r, s)
			}
		}
	}

	// thinkingWaveDots/pulsedStatusLine over a full pulse cycle -- this
	// is exactly the check that would have caught thinkingWaveGlyphs's
	// "·"/"˙": their frames only surface partway through the cycle, so
	// checking phase 0 alone would have missed them.
	for phase := 0; phase < activityPulsePeriodTicks; phase++ {
		assertASCII(t, "thinkingWaveDots", thinkingWaveDots(phase))
		for _, w := range []int{0, 1, 4, 8, 40, 120} {
			assertASCII(t, "pulsedStatusLine", pulsedStatusLine(phase, w))
		}
	}

	calls := []activityCall{
		{callID: "1", name: "run_tests", status: "running", detail: `{"file": "solution.go"}`},
		{callID: "2", name: "read_file", status: "running", detail: ""},
		{callID: "3", name: "grep", status: "done", detail: "3 matches found in solution.go"},
		{callID: "4", name: "run_tests", status: "failed", detail: "exit status 1: assertion failed"},
		{callID: "5", name: "no_detail", status: "done", detail: ""},
		{callID: "6", name: "just_settled", status: "done", detail: "wrapped up fine", completedAt: time.Now()},
		{callID: "7", name: "just_failed", status: "failed", detail: "boom", completedAt: time.Now()},
	}

	for _, w := range []int{0, 1, 4, 8, 20, 40, 80, 120} {
		for _, c := range calls {
			assertASCII(t, "activityCallHeader", activityCallHeader(c, w))
			for _, line := range activityOutputLines(c, w) {
				assertASCII(t, "activityOutputLines", line)
			}
			for _, line := range pulsedCallLines(c, 0, w) {
				assertASCII(t, "pulsedCallLines", line)
			}
			assertASCII(t, "settledCallLine", settledCallLine(c, w))
		}
		assertASCII(t, "toolUsageSummary", toolUsageSummary(calls, w))
	}

	for _, rgb := range [][3]int{{0, 0, 0}, {255, 255, 255}, {47, 166, 166}, {191, 252, 247}, {240, 60, 60}} {
		assertASCII(t, "coloredDot", coloredDot(rgb[0], rgb[1], rgb[2]))
		assertASCII(t, "coloredGlyph", coloredGlyph(rgb[0], rgb[1], rgb[2], activityDotGlyph))
	}
	assertASCII(t, "dimSpan", dimSpan("plain text"))
	assertASCII(t, "activityOutputHighlight", activityOutputHighlight("plain text"))
	assertASCII(t, "activityErrorNote", activityErrorNote("plain text"))
	assertASCII(t, "truncateLine (short)", truncateLine("short", 20))
	assertASCII(t, "truncateLine (cut)", truncateLine("a very long string that needs truncating", 10))

	// The raw glyph inventory itself, not just function output built
	// from it -- belt and suspenders with the functional checks above.
	assertASCII(t, "activityDotGlyph", activityDotGlyph)
	assertASCII(t, "truncateLineEllipsis", truncateLineEllipsis)
	for _, g := range thinkingWaveGlyphs {
		assertASCII(t, "thinkingWaveGlyphs entry", g)
	}
}
