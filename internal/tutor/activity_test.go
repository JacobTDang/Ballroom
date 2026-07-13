package tutor

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestActivityFeed_StartedAddsARunningLine(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "read_solution_file", "")
	if len(lines) != 1 || lines[0] != "o read_solution_file" {
		t.Errorf("lines = %v, want [\"o read_solution_file\"]", lines)
	}
}

func TestActivityFeed_StartedWithArgsShowsThemInParens(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "highlight_lines", `{"start_line":10,"end_line":20}`)
	if len(lines) != 1 || lines[0] != `o highlight_lines({"start_line":10,"end_line":20})` {
		t.Errorf("lines = %v, want the args shown in parens", lines)
	}
}

func TestActivityFeed_StartedWithEmptyOrNoArgsOmitsParens(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "read_solution_file", "{}")
	if len(lines) != 1 || lines[0] != "o read_solution_file" {
		t.Errorf("lines = %v, want no parens for empty/no-op args", lines)
	}
}

func TestActivityFeed_FinishedUpdatesTheMatchingCallToDone(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.finished("call-1", "312 bytes")
	if len(lines) != 1 || lines[0] != "o read_solution_file  312 bytes" {
		t.Errorf("lines = %v, want the call marked done with its result", lines)
	}
}

func TestActivityFeed_FinishedWithEmptyResultOmitsTrailingSpace(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "highlight_lines", "")
	lines := f.finished("call-1", "")
	if len(lines) != 1 || lines[0] != "o highlight_lines" {
		t.Errorf("lines = %v, want just the dot and name, no trailing separator", lines)
	}
}

func TestActivityFeed_FailedUpdatesTheMatchingCallToFailed(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_test_output", "")
	lines := f.failed("call-1", "no test run yet")
	if len(lines) != 1 || lines[0] != "o read_test_output - failed: no test run yet" {
		t.Errorf("lines = %v, want the call marked failed with the error", lines)
	}
}

func TestActivityFeed_FinishedForUnknownCallIDIsANoOp(t *testing.T) {
	// A callID that was never started (or already dropped by the cap
	// below) must not panic or fabricate a new entry -- eino's own
	// OnEnd/OnError always follow a real OnStart for the same call, but
	// this call may have aged out of the capped list already.
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.finished("call-unknown", "some result")
	if len(lines) != 1 || lines[0] != "o read_solution_file" {
		t.Errorf("lines = %v, want the existing call untouched and no new entry added", lines)
	}
}

func TestActivityFeed_MultipleCallsPreserveStartOrder(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.started("call-2", "read_problem_statement", "")
	if len(lines) != 2 || lines[0] != "o read_solution_file" || lines[1] != "o read_problem_statement" {
		t.Errorf("lines = %v, want both calls in start order", lines)
	}
}

func TestActivityFeed_CapsAtFourDroppingTheOldest(t *testing.T) {
	f := &activityFeed{}
	for i := 1; i <= 5; i++ {
		f.started(fmt.Sprintf("call-%d", i), fmt.Sprintf("tool_%d", i), "")
	}
	lines := f.started("call-6", "tool_6", "")
	if len(lines) != activityToolLines {
		t.Fatalf("len(lines) = %d, want %d (the cap)", len(lines), activityToolLines)
	}
	if lines[0] != "o tool_3" {
		t.Errorf("lines[0] = %q, want the oldest (tool_1, tool_2) dropped, starting at tool_3", lines[0])
	}
	if lines[len(lines)-1] != "o tool_6" {
		t.Errorf("lines[last] = %q, want the newest call last", lines[len(lines)-1])
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

func TestFormatActivityLine_IsADotFollowedByActivityLineBody(t *testing.T) {
	cases := []activityCall{
		{name: "read_solution_file", status: "running"},
		{name: "highlight_lines", status: "running", detail: `{"start_line":1}`},
		{name: "read_solution_file", status: "done", detail: "312 bytes"},
		{name: "read_test_output", status: "failed", detail: "no test run yet"},
	}
	for _, c := range cases {
		want := "o " + activityLineBody(c)
		if got := formatActivityLine(c); got != want {
			t.Errorf("formatActivityLine(%+v) = %q, want %q", c, got, want)
		}
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

func TestPulsedStatusLine_ContainsThePlainStatusText(t *testing.T) {
	got := pulsedStatusLine(0, 80)
	if !strings.Contains(got, "Thinking...") {
		t.Errorf("pulsedStatusLine(0, 80) = %q, want it to contain the plain status text", got)
	}
	if !strings.Contains(got, "\033[38;2;") {
		t.Errorf("pulsedStatusLine(0, 80) = %q, want a truecolor escape for the dot", got)
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
	if got := activityCallHeader(c); got != `highlight_lines({"start_line":1})` {
		t.Errorf("activityCallHeader(%+v) = %q, want args shown in parens", c, got)
	}
}

func TestActivityCallHeader_RunningWithNoArgsOmitsParens(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "running", detail: "{}"}
	if got := activityCallHeader(c); got != "read_solution_file" {
		t.Errorf("activityCallHeader(%+v) = %q, want no parens for empty args", c, got)
	}
}

func TestActivityCallHeader_DoneOmitsResultInline(t *testing.T) {
	// The result now belongs to activityOutputLines, indented beneath
	// this header -- a real bug found live: it used to be crammed onto
	// this same line, truncating to almost nothing in a normal-width
	// pane.
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes of real file content"}
	if got := activityCallHeader(c); got != "read_solution_file" {
		t.Errorf("activityCallHeader(%+v) = %q, want just the name, no inline result", c, got)
	}
}

func TestActivityCallHeader_FailedOmitsErrorInlineButFlagsFailure(t *testing.T) {
	c := activityCall{name: "read_test_output", status: "failed", detail: "no test run yet"}
	if got := activityCallHeader(c); got != "read_test_output - failed" {
		t.Errorf("activityCallHeader(%+v) = %q, want the name flagged failed, no inline error detail", c, got)
	}
}

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

func TestToolUsageSummary_OneDoneCallShowsNameAndIndentedOutput(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done", detail: "312 bytes"}}
	got := toolUsageSummary(calls, 80)
	if !strings.Contains(got, "read_solution_file") {
		t.Errorf("toolUsageSummary(...) = %q, want the tool name", got)
	}
	// Not activityIndent+"312 bytes" as one literal substring -- the
	// output is now color-highlighted, so an escape sequence sits
	// between the indent and the content itself.
	if !strings.Contains(got, activityIndent) || !strings.Contains(got, "312 bytes") {
		t.Errorf("toolUsageSummary(...) = %q, want the result indented beneath the name", got)
	}
}

func TestToolUsageSummary_FailedCallShowsFailedSuffixAndIndentedError(t *testing.T) {
	calls := []activityCall{{name: "read_test_output", status: "failed", detail: "no test run yet"}}
	got := toolUsageSummary(calls, 80)
	if !strings.Contains(got, "read_test_output - failed") {
		t.Errorf("toolUsageSummary(...) = %q, want the name flagged failed", got)
	}
	if !strings.Contains(got, activityIndent) || !strings.Contains(got, "no test run yet") {
		t.Errorf("toolUsageSummary(...) = %q, want the error indented beneath the name", got)
	}
}

func TestToolUsageSummary_MultipleCallsEachGetTheirOwnHeaderAndOutput(t *testing.T) {
	calls := []activityCall{
		{name: "read_solution_file", status: "done", detail: "312 bytes"},
		{name: "read_problem_statement", status: "done", detail: "problem text"},
	}
	got := toolUsageSummary(calls, 80)
	for _, want := range []string{"read_solution_file", "312 bytes", "read_problem_statement", "problem text"} {
		if !strings.Contains(got, want) {
			t.Errorf("toolUsageSummary(...) = %q, want it to contain %q", got, want)
		}
	}
	// The first call's header must precede the second call's, matching
	// call order.
	if strings.Index(got, "read_solution_file") > strings.Index(got, "read_problem_statement") {
		t.Errorf("toolUsageSummary(...) = %q, want calls in order", got)
	}
}

func TestToolUsageSummary_OutputContentIsColoredGray(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done", detail: "312 bytes"}}
	got := toolUsageSummary(calls, 80)
	want := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityOutputHighlightR, activityOutputHighlightG, activityOutputHighlightB)
	if !strings.Contains(got, want) {
		t.Errorf("toolUsageSummary(...) = %q, want the output colored faded gray", got)
	}
}

// TestToolUsageSummary_HeaderLineNeverCarriesTheOutputColor guards a
// real design correction from live use: the tool call's own name
// should stay the normal (unhighlighted) text color -- only the raw
// output beneath it gets colored.
func TestToolUsageSummary_HeaderLineNeverCarriesTheOutputColor(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done", detail: "312 bytes"}}
	got := toolUsageSummary(calls, 80)
	lines := strings.Split(got, "\n")
	headerLine := lines[0]
	outputColor := fmt.Sprintf("\033[38;2;%d;%d;%dm", activityOutputHighlightR, activityOutputHighlightG, activityOutputHighlightB)
	if strings.Contains(headerLine, outputColor) {
		t.Errorf("header line %q, want it to never contain the output's gray color escape", headerLine)
	}
}

func TestToolUsageSummary_HeaderIsColored(t *testing.T) {
	calls := []activityCall{{name: "read_solution_file", status: "done"}}
	if got := toolUsageSummary(calls, 80); !strings.Contains(got, "\033[38;2;") {
		t.Errorf("toolUsageSummary(...) = %q, want the header colored, matching the live activity display", got)
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
