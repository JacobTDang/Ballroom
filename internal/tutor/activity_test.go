package tutor

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestActivityFeed_StartedAddsARunningCall(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	calls := f.currentCalls()
	if len(calls) != 1 || calls[0].name != "read_solution_file" || calls[0].status != "running" {
		t.Errorf("currentCalls() = %+v, want one running read_solution_file", calls)
	}
}

func TestActivityFeed_StartedStoresTheArgsPreviewAsDetail(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "highlight_lines", `{"start_line":10,"end_line":20}`)
	if got := f.currentCalls()[0].detail; got != `{"start_line":10,"end_line":20}` {
		t.Errorf("detail = %q, want the args preview stored for the header renderer", got)
	}
}

func TestActivityFeed_FinishedUpdatesTheMatchingCallToDone(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	f.finished("call-1", "312 bytes")
	calls := f.currentCalls()
	if len(calls) != 1 || calls[0].status != "done" || calls[0].detail != "312 bytes" {
		t.Errorf("currentCalls() = %+v, want the call marked done with its result as detail", calls)
	}
}

func TestActivityFeed_FailedUpdatesTheMatchingCallToFailed(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_test_output", "")
	f.failed("call-1", "no test run yet")
	calls := f.currentCalls()
	if len(calls) != 1 || calls[0].status != "failed" || calls[0].detail != "no test run yet" {
		t.Errorf("currentCalls() = %+v, want the call marked failed with the error as detail", calls)
	}
}

func TestActivityFeed_FinishedForUnknownCallIDIsANoOp(t *testing.T) {
	// A callID that was never started (or already dropped by the cap
	// below) must not panic or fabricate a new entry -- eino's own
	// OnEnd/OnError always follow a real OnStart for the same call, but
	// this call may have aged out of the capped list already.
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	f.finished("call-unknown", "some result")
	calls := f.currentCalls()
	if len(calls) != 1 || calls[0].status != "running" {
		t.Errorf("currentCalls() = %+v, want the existing call untouched and no new entry added", calls)
	}
}

func TestActivityFeed_CapsAtFourDroppingTheOldest(t *testing.T) {
	f := &activityFeed{}
	for i := 1; i <= 5; i++ {
		f.started(fmt.Sprintf("call-%d", i), fmt.Sprintf("tool_%d", i), "")
	}
	f.started("call-6", "tool_6", "")
	calls := f.currentCalls()
	if len(calls) != activityToolLines {
		t.Fatalf("len(calls) = %d, want %d (the cap)", len(calls), activityToolLines)
	}
	if calls[0].name != "tool_3" {
		t.Errorf("calls[0].name = %q, want the oldest (tool_1, tool_2) dropped, starting at tool_3", calls[0].name)
	}
	if calls[len(calls)-1].name != "tool_6" {
		t.Errorf("calls[last].name = %q, want the newest call last", calls[len(calls)-1].name)
	}
}

func TestActivityFeed_CurrentCallsReturnsACopyNotTheLiveSlice(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")

	calls := f.currentCalls()
	calls[0].name = "mutated"

	if f.currentCalls()[0].name != "read_solution_file" {
		t.Error("mutating the returned slice affected the feed's internal state -- currentCalls must return a copy")
	}
}

func TestActivityFeed_CurrentCallsMatchesStartedOrder(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	f.started("call-2", "read_problem_statement", "")

	calls := f.currentCalls()
	if len(calls) != 2 || calls[0].name != "read_solution_file" || calls[1].name != "read_problem_statement" {
		t.Errorf("currentCalls() = %+v, want both calls in start order", calls)
	}
}

// --- activityCall.completedAt -- drives the settle fade (see
// activitySettledDotColor): the dot fades from its glow color down to
// the resting base color over activitySettleFadeDuration once a call
// settles, instead of snapping to base instantly, per an explicit
// request for smoother transitions ("after tool calling there is a
// slight fade").

func TestActivityFeed_StartedLeavesCompletedAtZero(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	calls := f.currentCalls()
	if !calls[0].completedAt.IsZero() {
		t.Errorf("completedAt = %v, want zero -- a just-started call hasn't settled yet", calls[0].completedAt)
	}
}

func TestActivityFeed_FinishedSetsCompletedAt(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	before := time.Now()
	f.finished("call-1", "312 bytes")
	after := time.Now()

	calls := f.currentCalls()
	if calls[0].completedAt.Before(before) || calls[0].completedAt.After(after) {
		t.Errorf("completedAt = %v, want it set to roughly now (between %v and %v)", calls[0].completedAt, before, after)
	}
}

func TestActivityFeed_FailedSetsCompletedAt(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_test_output", "")
	f.failed("call-1", "no test run yet")
	if f.currentCalls()[0].completedAt.IsZero() {
		t.Error("completedAt is zero, want it set once a call fails too, not just when it succeeds")
	}
}

// TestActivityDotColor_RunningAtHalfPeriodIsExactlyBaseColor and
// TestActivityDotColor_RunningAtPhaseZeroGlowsBrighterThanBase replace
// the old dim-toward-black pulse's tests -- per live feedback ("I would
// like a glowing dot") the pulse now blends from the base teal (the
// resting point, at the half-period mark) toward a brighter, paler
// highlight at the peak (phase 0), instead of dimming the base color
// down toward black. It never gets darker than the resting base color
// at any point in the cycle.
func TestActivityDotColor_RunningAtHalfPeriodIsExactlyBaseColor(t *testing.T) {
	r, g, b := activityDotColor("running", activityPulsePeriodTicks/2)
	if r != activityPulseBaseR || g != activityPulseBaseG || b != activityPulseBaseB {
		t.Errorf("activityDotColor(running, period/2) = (%d,%d,%d), want exactly the resting base color (%d,%d,%d)", r, g, b, activityPulseBaseR, activityPulseBaseG, activityPulseBaseB)
	}
}

func TestActivityDotColor_RunningAtPhaseZeroGlowsBrighterThanBase(t *testing.T) {
	r, g, b := activityDotColor("running", 0)
	if r < activityPulseBaseR || g < activityPulseBaseG || b < activityPulseBaseB {
		t.Errorf("activityDotColor(running, 0) = (%d,%d,%d), want every channel at least as bright as the base (%d,%d,%d) -- this is the glow peak", r, g, b, activityPulseBaseR, activityPulseBaseG, activityPulseBaseB)
	}
	if r == activityPulseBaseR && g == activityPulseBaseG && b == activityPulseBaseB {
		t.Error("activityDotColor(running, 0) equals the base color exactly, want a visibly brighter glow at the peak")
	}
}

func TestActivityDotColor_RunningNeverDimmerThanBase(t *testing.T) {
	for phase := 0; phase < activityPulsePeriodTicks; phase++ {
		r, g, b := activityDotColor("running", phase)
		if r < activityPulseBaseR || g < activityPulseBaseG || b < activityPulseBaseB {
			t.Errorf("activityDotColor(running, %d) = (%d,%d,%d), want no channel dimmer than the base (%d,%d,%d) anywhere in the cycle", phase, r, g, b, activityPulseBaseR, activityPulseBaseG, activityPulseBaseB)
		}
	}
}

func TestActivityDotColor_RunningIsPeriodic(t *testing.T) {
	a := [3]int{}
	b := [3]int{}
	a[0], a[1], a[2] = activityDotColor("running", 3)
	b[0], b[1], b[2] = activityDotColor("running", 3+activityPulsePeriodTicks)
	if a != b {
		t.Errorf("activityDotColor(running, 3) = %v, activityDotColor(running, 3+period) = %v, want equal (periodic)", a, b)
	}
}

func TestActivityDotColor_DoneAndFailedAreStaticRegardlessOfPhase(t *testing.T) {
	for _, status := range []string{"done", "failed"} {
		for _, phase := range []int{0, 1, activityPulsePeriodTicks / 2, 100} {
			r, g, b := activityDotColor(status, phase)
			if r != activityPulseBaseR || g != activityPulseBaseG || b != activityPulseBaseB {
				t.Errorf("activityDotColor(%s, %d) = (%d,%d,%d), want the static base color, unaffected by phase", status, phase, r, g, b)
			}
		}
	}
}

// --- activitySettledDotColor -- the settle fade shown by pulsedCallLines
// once a call is done/failed (activityDotColor above, used directly by
// pulsedStatusLine, deliberately stays a pure phase->color function with
// no notion of "just settled"; this is a separate function specifically
// for that transition instead of reshaping activityDotColor's signature
// for every caller). Takes elapsed directly rather than a completedAt
// timestamp so it's a pure, deterministic function of its input -- no
// real-clock dependency to fake out in tests.

func TestActivitySettledDotColor_AtZeroElapsedIsTheGlowColor(t *testing.T) {
	r, g, b := activitySettledDotColor(0)
	if r != activityPulseGlowR || g != activityPulseGlowG || b != activityPulseGlowB {
		t.Errorf("activitySettledDotColor(0) = (%d,%d,%d), want exactly the glow color (%d,%d,%d) -- the fade starts here", r, g, b, activityPulseGlowR, activityPulseGlowG, activityPulseGlowB)
	}
}

func TestActivitySettledDotColor_AtFadeDurationIsExactlyBase(t *testing.T) {
	r, g, b := activitySettledDotColor(activitySettleFadeDuration)
	if r != activityPulseBaseR || g != activityPulseBaseG || b != activityPulseBaseB {
		t.Errorf("activitySettledDotColor(fadeDuration) = (%d,%d,%d), want exactly the resting base color (%d,%d,%d)", r, g, b, activityPulseBaseR, activityPulseBaseG, activityPulseBaseB)
	}
}

func TestActivitySettledDotColor_PastFadeDurationStaysAtBase(t *testing.T) {
	r, g, b := activitySettledDotColor(activitySettleFadeDuration * 10)
	if r != activityPulseBaseR || g != activityPulseBaseG || b != activityPulseBaseB {
		t.Errorf("activitySettledDotColor(10x fadeDuration) = (%d,%d,%d), want the resting base color, not clamped/wrapped weirdly", r, g, b)
	}
}

func TestActivitySettledDotColor_MidwayIsBetweenGlowAndBase(t *testing.T) {
	r, g, b := activitySettledDotColor(activitySettleFadeDuration / 2)
	if r <= activityPulseBaseR || r >= activityPulseGlowR {
		t.Errorf("activitySettledDotColor(fadeDuration/2) red = %d, want strictly between the base (%d) and glow (%d)", r, activityPulseBaseR, activityPulseGlowR)
	}
	if g <= activityPulseBaseG || g >= activityPulseGlowG {
		t.Errorf("activitySettledDotColor(fadeDuration/2) green = %d, want strictly between the base (%d) and glow (%d)", g, activityPulseBaseG, activityPulseGlowG)
	}
	if b <= activityPulseBaseB || b >= activityPulseGlowB {
		t.Errorf("activitySettledDotColor(fadeDuration/2) blue = %d, want strictly between the base (%d) and glow (%d)", b, activityPulseBaseB, activityPulseGlowB)
	}
}

func TestActivitySettledDotColor_NegativeElapsedClampsToGlow(t *testing.T) {
	// Defensive: elapsed should never actually be negative in practice,
	// but this must not produce a nonsensical color if it somehow is.
	r, g, b := activitySettledDotColor(-1)
	if r != activityPulseGlowR || g != activityPulseGlowG || b != activityPulseGlowB {
		t.Errorf("activitySettledDotColor(-1) = (%d,%d,%d), want clamped to the glow color", r, g, b)
	}
}

func TestPulsedStatusLine_ContainsThePlainStatusText(t *testing.T) {
	got := pulsedStatusLine(0, 80)
	if !strings.Contains(got, "Thinking") {
		t.Errorf("pulsedStatusLine(0, 80) = %q, want it to contain the plain status text", got)
	}
	if !strings.Contains(got, "\033[38;2;") {
		t.Errorf("pulsedStatusLine(0, 80) = %q, want a truecolor escape for the dot", got)
	}
}

// --- thinkingWaveDots -- the animated ellipsis after "Thinking" ---

func TestThinkingWaveDots_UsesWaveGlyphs(t *testing.T) {
	got := thinkingWaveDots(0)
	found := false
	for _, g := range thinkingWaveGlyphs {
		if strings.Contains(got, g) {
			found = true
		}
	}
	if !found {
		t.Errorf("thinkingWaveDots(0) = %q, want at least one wave glyph from %v", got, thinkingWaveGlyphs)
	}
}

func TestThinkingWaveDots_AnimatesOverPhase(t *testing.T) {
	if thinkingWaveDots(0) == thinkingWaveDots(9) {
		t.Error("thinkingWaveDots identical at phase 0 and 9 -- the dots aren't animating")
	}
}

func TestThinkingWaveDots_WaveTravelsAcrossDots(t *testing.T) {
	// The defining property of a traveling wave: each dot repeats its
	// neighbor's motion a fixed phase later. Dot i at phase p must show
	// the same height dot i+1 shows at phase p + spread.
	for p := 0; p < activityPulsePeriodTicks; p += 5 {
		for i := 0; i < thinkingWaveDotCount-1; i++ {
			a := thinkingWaveLevel(p, i)
			b := thinkingWaveLevel(p+thinkingWaveSpreadTicks, i+1)
			if a != b {
				t.Fatalf("wave not traveling: dot %d at phase %d has level %d, dot %d at phase %d has level %d -- want equal", i, p, a, i+1, p+thinkingWaveSpreadTicks, b)
			}
		}
	}
}

func TestPulsedStatusLine_NarrowWidthDropsWaveNotEscapes(t *testing.T) {
	got := pulsedStatusLine(0, 6)
	if strings.Count(got, "\033[38;2;") > 0 && !strings.Contains(got, "\033[0m") {
		t.Errorf("pulsedStatusLine(0, 6) = %q, want any emitted escape properly closed", got)
	}
}

// --- activityCallHeader / activityOutputLines / pulsedCallLines --
// replaced the old single-line pulsedCallLine, which crammed a
// completed call's raw result onto the same line as its name (a real
// UX complaint: it truncated to almost nothing in a normal-width pane).
// A completed/failed call's output now renders on its own indented
// lines beneath the header, like Claude Code's own tool-result display.

func TestActivityCallHeader_RunningWithArgsShowsThemInline(t *testing.T) {
	c := activityCall{name: "highlight_lines", status: "running", detail: `{"start_line":1}`}
	got := activityCallHeader(c, 80)
	if plain := stripAnsiTest(got); plain != `highlight_lines({"start_line":1})` {
		t.Errorf("activityCallHeader(%+v) stripped = %q, want args shown in parens", c, plain)
	}
	if !strings.Contains(got, mdBoldOn+"highlight_lines"+mdBoldOff) {
		t.Errorf("activityCallHeader(%+v) = %q, want the tool name bold", c, got)
	}
	if !strings.Contains(got, mdDimColor) {
		t.Errorf("activityCallHeader(%+v) = %q, want the args dimmed", c, got)
	}
}

func TestActivityCallHeader_RunningWithNoArgsOmitsParens(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "running", detail: "{}"}
	if got := stripAnsiTest(activityCallHeader(c, 80)); got != "read_solution_file" {
		t.Errorf("activityCallHeader(%+v) stripped = %q, want no parens for empty args", c, got)
	}
}

func TestActivityCallHeader_DoneOmitsResultInline(t *testing.T) {
	// The result now belongs to activityOutputLines, indented beneath
	// this header -- a real bug found live: it used to be crammed onto
	// this same line, truncating to almost nothing in a normal-width
	// pane.
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes of real file content"}
	got := activityCallHeader(c, 80)
	if plain := stripAnsiTest(got); plain != "read_solution_file" {
		t.Errorf("activityCallHeader(%+v) stripped = %q, want just the name, no inline result", c, plain)
	}
	if !strings.Contains(got, mdBoldOn+"read_solution_file"+mdBoldOff) {
		t.Errorf("activityCallHeader(%+v) = %q, want the name bold", c, got)
	}
}

func TestActivityCallHeader_FailedOmitsErrorInlineButFlagsFailure(t *testing.T) {
	c := activityCall{name: "read_test_output", status: "failed", detail: "no test run yet"}
	got := activityCallHeader(c, 80)
	if plain := stripAnsiTest(got); plain != "read_test_output - failed" {
		t.Errorf("activityCallHeader(%+v) stripped = %q, want the name flagged failed, no inline error detail", c, plain)
	}
	red := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityErrorNoteR, activityErrorNoteG, activityErrorNoteB)
	if !strings.Contains(got, red) {
		t.Errorf("activityCallHeader(%+v) = %q, want the failed flag in the error red", c, got)
	}
}

func TestActivityCallHeader_NarrowWidthTruncatesPlainTextNotEscapes(t *testing.T) {
	c := activityCall{name: strings.Repeat("x", 200), status: "running", detail: `{"a":1}`}
	got := activityCallHeader(c, 20)
	if plain := stripAnsiTest(got); len([]rune(plain)) > 20 {
		t.Errorf("stripped header %q is %d runes, want at most the 20-rune budget", plain, len([]rune(plain)))
	}
	// Every escape opener must still have its terminating 'm' -- a
	// sliced escape would leave an unterminated "\033[38;2;..." tail.
	if opens, whole := strings.Count(got, "\033["), len(ansiEscapePattern.FindAllString(got, -1)); opens != whole {
		t.Errorf("activityCallHeader(...) = %q: %d escape openers but %d complete escapes, want every ANSI escape intact after truncation", got, opens, whole)
	}
}

// ansiEscapePattern matches one complete SGR escape sequence, for
// escape-integrity assertions after truncation.
var ansiEscapePattern = regexp.MustCompile("\033\\[[0-9;]*m")

func TestActivityOutputLines_RunningCallHasNoOutputYet(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "running", detail: "{}"}
	if got := activityOutputLines(c, 80); got != nil {
		t.Errorf("activityOutputLines(running) = %v, want nil -- a running call has no result yet", got)
	}
}

func TestActivityOutputLines_EmptyDetailReturnsNil(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: ""}
	if got := activityOutputLines(c, 80); got != nil {
		t.Errorf("activityOutputLines(empty detail) = %v, want nil", got)
	}
}

func TestActivityOutputLines_ShortResultIsOneIndentedLine(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := activityOutputLines(c, 80)
	if len(got) != 1 {
		t.Fatalf("activityOutputLines(...) = %v, want exactly 1 line", got)
	}
	if !strings.HasPrefix(got[0], activityIndent) || !strings.Contains(got[0], "312 bytes") {
		t.Errorf("activityOutputLines(...)[0] = %q, want it indented and containing the result", got[0])
	}
}

func TestActivityOutputLines_FailedDetailIsIndentedToo(t *testing.T) {
	c := activityCall{name: "read_test_output", status: "failed", detail: "no test run yet"}
	got := activityOutputLines(c, 80)
	if len(got) != 1 || !strings.HasPrefix(got[0], activityIndent) || !strings.Contains(got[0], "no test run yet") {
		t.Errorf("activityOutputLines(failed) = %v, want the error indented and shown", got)
	}
}

func TestActivityOutputLines_LongResultWrapsAcrossMultipleIndentedLines(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: strings.Repeat("word ", 40)}
	got := activityOutputLines(c, 30)
	if len(got) <= 1 {
		t.Fatalf("activityOutputLines(...) = %v, want more than 1 wrapped line for a long result", got)
	}
	for _, line := range got {
		if !strings.HasPrefix(line, activityIndent) {
			t.Errorf("line %q missing the indent prefix", line)
		}
	}
}

func TestActivityOutputLines_CapsAtThreeLinesWithEllipsisMarkerOnLast(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: strings.Repeat("word ", 200)}
	got := activityOutputLines(c, 20)
	if len(got) != activityOutputPreviewLines {
		t.Fatalf("activityOutputLines(...) = %d lines, want capped at %d", len(got), activityOutputPreviewLines)
	}
	last := got[len(got)-1]
	// Not HasSuffix -- the line's true suffix is now the color reset
	// escape (activityOutputLines colors the content gray), so the
	// ellipsis marker itself is Contains'd instead.
	if !strings.Contains(last, truncateLineEllipsis) {
		t.Errorf("last line %q, want it to contain %q to signal the result was cut off", last, truncateLineEllipsis)
	}
}

func TestActivityOutputLines_ContentIsColoredGray(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := activityOutputLines(c, 80)
	want := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityOutputHighlightR, activityOutputHighlightG, activityOutputHighlightB)
	if len(got) != 1 || !strings.Contains(got[0], want) {
		t.Errorf("activityOutputLines(...) = %v, want the output colored faded gray", got)
	}
}

func TestActivityErrorNote_ContainsTheErrorColorAndTheText(t *testing.T) {
	got := activityErrorNote("could not reach http://localhost:11434: connection refused")
	color := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityErrorNoteR, activityErrorNoteG, activityErrorNoteB)
	if !strings.Contains(got, color) {
		t.Errorf("activityErrorNote(...) = %q, want the error-red color escape", got)
	}
	if !strings.Contains(got, "could not reach http://localhost:11434: connection refused") {
		t.Errorf("activityErrorNote(...) = %q, want the original text preserved", got)
	}
}

func TestActivityErrorNote_StartsWithADefensiveReset(t *testing.T) {
	got := activityErrorNote("some error")
	if !strings.HasPrefix(got, "\033[0m\033[38;2;") {
		t.Errorf("activityErrorNote(...) = %q, want it prefixed with an explicit reset before its own color, same defensive pattern as coloredDot/activityOutputHighlight", got)
	}
}

func TestActivityOutputLines_ColorEscapeStartsWithADefensiveReset(t *testing.T) {
	// Guards the fix for a real bug found live: a later, supposedly
	// uncolored line was seen visibly inheriting an earlier line's
	// color in the real terminal renderer. Every colored span now opens
	// with an explicit \033[0m before its own color code so it can
	// never inherit stray state from whatever rendered immediately
	// before it.
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := activityOutputLines(c, 80)
	if len(got) != 1 {
		t.Fatalf("activityOutputLines(...) = %v, want exactly 1 line", got)
	}
	if !strings.Contains(got[0], "\033[0m\033[38;2;") {
		t.Errorf("activityOutputLines(...)[0] = %q, want the color escape prefixed with an explicit reset", got[0])
	}
}

func TestPulsedCallLines_RunningCallIsJustTheHeaderLine(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "running", detail: "{}"}
	got := pulsedCallLines(c, 0, 80)
	if len(got) != 1 {
		t.Fatalf("pulsedCallLines(running) = %v, want exactly 1 line (no output yet)", got)
	}
	if !strings.Contains(got[0], "read_solution_file") || !strings.Contains(got[0], "\033[38;2;") {
		t.Errorf("pulsedCallLines(running)[0] = %q, want the colored header", got[0])
	}
}

func TestPulsedCallLines_DoneCallHasHeaderThenIndentedOutput(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := pulsedCallLines(c, 0, 80)
	if len(got) != 2 {
		t.Fatalf("pulsedCallLines(done) = %v, want header + 1 output line", got)
	}
	if !strings.Contains(got[0], "read_solution_file") || strings.Contains(got[0], "312 bytes") {
		t.Errorf("pulsedCallLines(done)[0] = %q, want just the colored header, no inline result", got[0])
	}
	if !strings.HasPrefix(got[1], activityIndent) || !strings.Contains(got[1], "312 bytes") {
		t.Errorf("pulsedCallLines(done)[1] = %q, want the indented result", got[1])
	}
}

// TestPulsedCallLines_JustSettledCallGlowsNotBase confirms pulsedCallLines
// actually wires activitySettledDotColor in: a call that completedAt
// just now should render with the glow color, not the base color a
// naive "done -> always base" implementation would produce.
func TestPulsedCallLines_JustSettledCallGlowsNotBase(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes", completedAt: time.Now()}
	got := pulsedCallLines(c, 0, 80)
	glow := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityPulseGlowR, activityPulseGlowG, activityPulseGlowB)
	if !strings.Contains(got[0], glow) {
		t.Errorf("pulsedCallLines(...)[0] = %q, want the just-settled header to start the fade at the glow color", got[0])
	}
}

// TestPulsedCallLines_LongSettledCallIsBase confirms a call that settled
// well outside the fade window renders with the plain resting base
// color, same as before this feature existed.
func TestPulsedCallLines_LongSettledCallIsBase(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes", completedAt: time.Now().Add(-10 * activitySettleFadeDuration)}
	got := pulsedCallLines(c, 0, 80)
	base := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityPulseBaseR, activityPulseBaseG, activityPulseBaseB)
	if !strings.Contains(got[0], base) {
		t.Errorf("pulsedCallLines(...)[0] = %q, want the long-settled header at the resting base color", got[0])
	}
}

func TestPulsedCallLines_HeaderTruncatesNotTheEscapeSequence(t *testing.T) {
	c := activityCall{name: strings.Repeat("x", 200), status: "running", detail: "{}"}
	got := pulsedCallLines(c, 0, 20) // narrow width
	if !strings.Contains(got[0], "\033[38;2;") || !strings.Contains(got[0], "mo\033[0m") {
		t.Errorf("pulsedCallLines(...)[0] = %q, want the truecolor escape sequence intact", got[0])
	}
}

// --- toolUsageSummary -- the permanent, settled record of which tools
// a completed turn used (plus each one's indented output), appended to
// tutorModel's displayLines (unlike the live activity region, which
// vanishes entirely once the turn ends).

func TestToolUsageSummary_EmptyForNoCalls(t *testing.T) {
	if got := toolUsageSummary(nil, 80); got != "" {
		t.Errorf("toolUsageSummary(nil, 80) = %q, want empty", got)
	}
	if got := toolUsageSummary([]activityCall{}, 80); got != "" {
		t.Errorf("toolUsageSummary([], 80) = %q, want empty", got)
	}
}

// --- settledCallLine -- the compact one-row-per-call form the settled
// summary uses: dot + bold name + dim inline first-result-line, so the
// permanent transcript record stays quiet (the multi-line indented
// preview remains a live-region-only affordance, where watching a
// result stream in is actually useful).

func TestSettledCallLine_DoneIsOneRowWithBoldNameAndDimSummary(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := settledCallLine(c, 80)
	if strings.Contains(got, "\n") {
		t.Fatalf("settledCallLine(...) = %q, want exactly one row", got)
	}
	if !strings.Contains(got, mdBoldOn+"read_solution_file"+mdBoldOff) {
		t.Errorf("settledCallLine(...) = %q, want the tool name bold", got)
	}
	if !strings.Contains(got, mdDimColor) || !strings.Contains(got, "312 bytes") {
		t.Errorf("settledCallLine(...) = %q, want the result summary present and dimmed", got)
	}
	if strings.Contains(got, activityIndent+"\033") {
		t.Errorf("settledCallLine(...) = %q, want no indented output block in the settled form", got)
	}
}

func TestSettledCallLine_EmptyDetailIsJustDotAndName(t *testing.T) {
	c := activityCall{name: "highlight_lines", status: "done", detail: ""}
	got := settledCallLine(c, 80)
	if plain := stripAnsiTest(got); plain != "o highlight_lines" {
		t.Errorf("settledCallLine(...) stripped = %q, want just the dot and name", plain)
	}
}

func TestSettledCallLine_FailedIsRedWithDetailInline(t *testing.T) {
	c := activityCall{name: "read_test_output", status: "failed", detail: "no test run yet"}
	got := settledCallLine(c, 80)
	if strings.Contains(got, "\n") {
		t.Fatalf("settledCallLine(failed) = %q, want exactly one row", got)
	}
	plain := stripAnsiTest(got)
	if !strings.Contains(plain, "read_test_output") || !strings.Contains(plain, "failed") || !strings.Contains(plain, "no test run yet") {
		t.Errorf("settledCallLine(failed) stripped = %q, want name, failed flag, and error inline", plain)
	}
	red := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityErrorNoteR, activityErrorNoteG, activityErrorNoteB)
	if !strings.Contains(got, red) {
		t.Errorf("settledCallLine(failed) = %q, want the failure in the error red", got)
	}
}

func TestSettledCallLine_FailedWithEmptyDetailStillFlagsFailure(t *testing.T) {
	// A failed call must never render indistinguishable from a clean
	// one just because the error text was empty.
	c := activityCall{name: "read_test_output", status: "failed", detail: ""}
	got := settledCallLine(c, 80)
	if plain := stripAnsiTest(got); !strings.Contains(plain, "failed") {
		t.Errorf("settledCallLine(failed, empty detail) stripped = %q, want the failure still flagged", plain)
	}
}

func TestSettledCallLine_MultilineDetailUsesOnlyTheFirstLine(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "first line of result\nsecond line\nthird"}
	got := settledCallLine(c, 80)
	plain := stripAnsiTest(got)
	if !strings.Contains(plain, "first line of result") {
		t.Errorf("settledCallLine(...) stripped = %q, want the first result line inline", plain)
	}
	if strings.Contains(plain, "second line") || strings.Contains(got, "\n") {
		t.Errorf("settledCallLine(...) = %q, want later result lines dropped, not wrapped", got)
	}
}

func TestSettledCallLine_NarrowWidthTruncatesPlainTextNotEscapes(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: strings.Repeat("word ", 40)}
	got := settledCallLine(c, 24)
	if plain := stripAnsiTest(got); len([]rune(plain)) > 24 {
		t.Errorf("stripped settled line %q is %d runes, want at most 24", plain, len([]rune(plain)))
	}
	if opens, whole := strings.Count(got, "\033["), len(ansiEscapePattern.FindAllString(got, -1)); opens != whole {
		t.Errorf("settledCallLine(...) = %q: %d escape openers but %d complete escapes, want none sliced", got, opens, whole)
	}
}

func TestSettledCallLine_TinyWidthDropsTheSummaryEntirely(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := settledCallLine(c, len("o read_solution_file")+3)
	plain := stripAnsiTest(got)
	if strings.Contains(plain, "312") {
		t.Errorf("settledCallLine(tiny) stripped = %q, want the summary dropped whole rather than a useless fragment", plain)
	}
}

func TestToolUsageSummary_OneRowPerCall(t *testing.T) {
	calls := []activityCall{
		{name: "read_solution_file", status: "done", detail: "312 bytes"},
		{name: "read_problem_statement", status: "done", detail: "problem text"},
		{name: "read_test_output", status: "failed", detail: "no test run yet"},
	}
	got := toolUsageSummary(calls, 80)
	if rows := strings.Split(got, "\n"); len(rows) != len(calls) {
		t.Fatalf("toolUsageSummary(...) = %d rows, want exactly one per call (%d):\n%s", len(rows), len(calls), got)
	}
}

func TestToolUsageSummary_InlineSummariesReplaceIndentedOutput(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done", detail: "312 bytes"}}
	got := toolUsageSummary(calls, 80)
	if !strings.Contains(got, "read_solution_file") || !strings.Contains(got, "312 bytes") {
		t.Errorf("toolUsageSummary(...) = %q, want name and result summary on the row", got)
	}
	if strings.Contains(got, "\n") {
		t.Errorf("toolUsageSummary(one call) = %q, want a single row, no indented output block", got)
	}
}

func TestToolUsageSummary_CallsStayInOrder(t *testing.T) {
	calls := []activityCall{
		{name: "read_solution_file", status: "done", detail: "312 bytes"},
		{name: "read_problem_statement", status: "done", detail: "problem text"},
	}
	got := toolUsageSummary(calls, 80)
	if strings.Index(got, "read_solution_file") > strings.Index(got, "read_problem_statement") {
		t.Errorf("toolUsageSummary(...) = %q, want calls in order", got)
	}
}

// TestToolUsageSummary_NameNeverCarriesTheSummaryDim guards the settled
// row's visual contract: the bold tool name stays the normal text
// color -- only the inline result summary after it is dimmed.
func TestToolUsageSummary_NameNeverCarriesTheSummaryDim(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done", detail: "312 bytes"}}
	got := toolUsageSummary(calls, 80)
	dimAt := strings.Index(got, mdDimColor)
	nameAt := strings.Index(got, "read_solution_file")
	if dimAt == -1 || nameAt == -1 || nameAt > dimAt {
		t.Errorf("toolUsageSummary(...) = %q, want the name before any dim span so it never renders dimmed", got)
	}
}

func TestToolUsageSummary_HeaderIsColored(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done"}}
	if got := toolUsageSummary(calls, 80); !strings.Contains(got, "\033[38;2;") {
		t.Errorf("toolUsageSummary(...) = %q, want the dot colored, matching the live activity display", got)
	}
}

func TestActivityFeed_ConcurrentStartedFinishedDoNotRace(t *testing.T) {
	f := &activityFeed{}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("call-%d", n)
			name := fmt.Sprintf("tool_%d", n)
			f.started(id, name, "")
			f.finished(id, "done")
		}(i)
	}
	wg.Wait()
}

// --- truncateLine (moved from scrollbox_test.go alongside truncateLine
// itself when scrollbox.go's hand-rolled ANSI box was deleted) ---

func TestTruncateLine_ShortStringUnchanged(t *testing.T) {
	if got := truncateLine("hello", 10); got != "hello" {
		t.Errorf("truncateLine(%q, 10) = %q, want unchanged", "hello", got)
	}
}

func TestTruncateLine_LongStringTruncatedWithEllipsis(t *testing.T) {
	// ASCII "..." only -- a real bug found live: the Unicode ellipsis
	// (…) and every other symbol this package originally used (⟳ → ✓ ✗)
	// rendered as unrecognizable fallback glyphs (tofu, looking like
	// stray underscores) in a real user's terminal font. Everything this
	// package writes must be plain ASCII plus the one glyph confirmed to
	// render everywhere: o (see formatActivityLine).
	got := truncateLine("this is a much longer string than the limit allows", 10)
	if runes := []rune(got); len(runes) != 10 {
		t.Errorf("truncateLine(...) = %q (len %d), want exactly 10 runes", got, len(runes))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("truncateLine(...) = %q, want it to end with \"...\"", got)
	}
}

func TestTruncateLine_MaxOfZeroOrLessReturnsEmpty(t *testing.T) {
	if got := truncateLine("anything", 0); got != "" {
		t.Errorf("truncateLine(_, 0) = %q, want empty", got)
	}
}

func TestTruncateLine_MaxSmallerThanEllipsisTruncatesTheEllipsisItself(t *testing.T) {
	if got := truncateLine("anything long enough to truncate", 2); got != ".." {
		t.Errorf("truncateLine(_, 2) = %q, want \"..\"", got)
	}
}
