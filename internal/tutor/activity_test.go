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
	if len(lines) != 1 || lines[0] != "● read_solution_file" {
		t.Errorf("lines = %v, want [\"● read_solution_file\"]", lines)
	}
}

func TestActivityFeed_StartedWithArgsShowsThemInParens(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "highlight_lines", `{"start_line":10,"end_line":20}`)
	if len(lines) != 1 || lines[0] != `● highlight_lines({"start_line":10,"end_line":20})` {
		t.Errorf("lines = %v, want the args shown in parens", lines)
	}
}

func TestActivityFeed_StartedWithEmptyOrNoArgsOmitsParens(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "read_solution_file", "{}")
	if len(lines) != 1 || lines[0] != "● read_solution_file" {
		t.Errorf("lines = %v, want no parens for empty/no-op args", lines)
	}
}

func TestActivityFeed_FinishedUpdatesTheMatchingCallToDone(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.finished("call-1", "312 bytes")
	if len(lines) != 1 || lines[0] != "● read_solution_file  312 bytes" {
		t.Errorf("lines = %v, want the call marked done with its result", lines)
	}
}

func TestActivityFeed_FinishedWithEmptyResultOmitsTrailingSpace(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "highlight_lines", "")
	lines := f.finished("call-1", "")
	if len(lines) != 1 || lines[0] != "● highlight_lines" {
		t.Errorf("lines = %v, want just the dot and name, no trailing separator", lines)
	}
}

func TestActivityFeed_FailedUpdatesTheMatchingCallToFailed(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_test_output", "")
	lines := f.failed("call-1", "no test run yet")
	if len(lines) != 1 || lines[0] != "● read_test_output - failed: no test run yet" {
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
	if len(lines) != 1 || lines[0] != "● read_solution_file" {
		t.Errorf("lines = %v, want the existing call untouched and no new entry added", lines)
	}
}

func TestActivityFeed_MultipleCallsPreserveStartOrder(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.started("call-2", "read_problem_statement", "")
	if len(lines) != 2 || lines[0] != "● read_solution_file" || lines[1] != "● read_problem_statement" {
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
	if lines[0] != "● tool_3" {
		t.Errorf("lines[0] = %q, want the oldest (tool_1, tool_2) dropped, starting at tool_3", lines[0])
	}
	if lines[len(lines)-1] != "● tool_6" {
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
		want := "● " + activityLineBody(c)
		if got := formatActivityLine(c); got != want {
			t.Errorf("formatActivityLine(%+v) = %q, want %q", c, got, want)
		}
	}
}

func TestActivityDotColor_RunningAtPhaseZeroIsFullBrightness(t *testing.T) {
	r, g, b := activityDotColor("running", 0)
	if r != activityPulseBaseR || g != activityPulseBaseG || b != activityPulseBaseB {
		t.Errorf("activityDotColor(running, 0) = (%d,%d,%d), want the full base color (%d,%d,%d)", r, g, b, activityPulseBaseR, activityPulseBaseG, activityPulseBaseB)
	}
}

func TestActivityDotColor_RunningAtHalfPeriodIsDimmerThanBase(t *testing.T) {
	r, _, _ := activityDotColor("running", activityPulsePeriodTicks/2)
	if r >= activityPulseBaseR {
		t.Errorf("activityDotColor(running, period/2) red = %d, want dimmer than the base %d", r, activityPulseBaseR)
	}
	if r <= 0 {
		t.Errorf("activityDotColor(running, period/2) red = %d, want never fully dark (min brightness floor)", r)
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

func TestPulsedCallLine_ContainsTheLineBodyUntruncatedEscapes(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: "312 bytes"}
	got := pulsedCallLine(c, 0, 80)
	if !strings.Contains(got, "read_solution_file  312 bytes") {
		t.Errorf("pulsedCallLine(...) = %q, want it to contain the formatted body", got)
	}
}

func TestPulsedCallLine_TruncatesBodyNotTheEscapeSequence(t *testing.T) {
	c := activityCall{name: "read_solution_file", status: "done", detail: strings.Repeat("x", 200)}
	got := pulsedCallLine(c, 0, 20) // narrow width
	// The color escape itself must survive intact (never sliced mid-sequence).
	if !strings.Contains(got, "\033[38;2;") || !strings.Contains(got, "m●\033[0m") {
		t.Errorf("pulsedCallLine(...) = %q, want the truecolor escape sequence intact", got)
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
	// render everywhere: ● (see formatActivityLine).
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
