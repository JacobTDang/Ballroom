package tutor

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
)

func TestThinkingDisplay_AnimatesImmediatelyEvenWithNoToolCalls(t *testing.T) {
	// The ball should spin for the whole turn, including turns where the
	// model never calls a tool at all — this is the main behavior change
	// from the old toolCallDisplay, which only ever activated on a tool
	// call.
	t.Setenv("KITTY_WINDOW_ID", "") // force the ANSI renderer — this test checks its specific escape codes
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	d.finish()

	// The ball's whole visual signal is carried in the truecolor
	// background escapes on opaque cells — after stripping ANSI, every
	// cell (opaque or background) is just plain spaces, so this checks
	// the raw buffer instead of the stripped text.
	if !strings.Contains(buf.String(), "\033[48;2;") {
		t.Fatal("expected a truecolor background escape from the ball even with no tool calls, got none")
	}
}

func TestThinkingDisplay_ToolNamesAccumulateInOutput(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)

	ctx := context.Background()
	d.onStart(ctx, &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackInput{})
	d.onStart(ctx, &callbacks.RunInfo{Name: "read_problem_statement"}, &tool.CallbackInput{})
	d.finish()

	got := stripAnsi(buf.String())
	if !strings.Contains(got, "read_solution_file") {
		t.Errorf("expected read_solution_file in the output, got:\n%s", got)
	}
	if !strings.Contains(got, "read_problem_statement") {
		t.Errorf("expected read_problem_statement in the output, got:\n%s", got)
	}
}

func TestThinkingDisplay_ShowsArgumentSummaryForHighlightLines(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)

	input := &tool.CallbackInput{ArgumentsInJSON: `{"file":"solution.go","start":4,"end":6,"note":"off by one"}`}
	d.onStart(context.Background(), &callbacks.RunInfo{Name: "highlight_lines"}, input)
	d.finish()

	got := stripAnsi(buf.String())
	if !strings.Contains(got, "highlight_lines (lines 4-6)") {
		t.Errorf("expected an argument summary in the output, got:\n%s", got)
	}
}

func TestThinkingDisplay_OnEndMarksTheRightEntryDoneRegardlessOfOrder(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)

	ctx1 := d.onStart(context.Background(), &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackInput{})
	ctx2 := d.onStart(context.Background(), &callbacks.RunInfo{Name: "read_problem_statement"}, &tool.CallbackInput{})

	// End the SECOND call first — must not mark the first (still
	// running) entry done too.
	d.onEnd(ctx2, &callbacks.RunInfo{Name: "read_problem_statement"}, &tool.CallbackOutput{})

	d.mu.Lock()
	if len(d.calls) != 2 {
		d.mu.Unlock()
		t.Fatalf("expected 2 entries, got %d", len(d.calls))
	}
	if d.calls[0].done {
		t.Error("expected the first (still running) entry to not be marked done")
	}
	if !d.calls[1].done {
		t.Error("expected the second (already ended) entry to be marked done")
	}
	d.mu.Unlock()

	d.onEnd(ctx1, &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackOutput{})
	d.mu.Lock()
	if !d.calls[0].done {
		t.Error("expected the first entry to be marked done after its own OnEnd")
	}
	d.mu.Unlock()

	d.finish()
}

func TestThinkingDisplay_ConcurrentToolCallsRaceSafe(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)

	names := []string{"read_solution_file", "read_problem_statement", "read_cursor_position"}
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			ctx := d.onStart(context.Background(), &callbacks.RunInfo{Name: name}, &tool.CallbackInput{})
			d.onEnd(ctx, &callbacks.RunInfo{Name: name}, &tool.CallbackOutput{})
		}(name)
	}
	wg.Wait()
	d.finish()

	got := stripAnsi(buf.String())
	for _, name := range names {
		if !strings.Contains(got, name) {
			t.Errorf("expected %q in the output, got:\n%s", name, got)
		}
	}
}

func TestThinkingDisplay_FinishIsIdempotent(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	d.finish()
	d.finish() // must not panic (closing d.stop twice)
}

func TestPulseBrightness_VariesAcrossFrames(t *testing.T) {
	dim := pulseBrightness(0)
	bright := pulseBrightness(dotPulsePeriod / 2)
	if dim >= bright {
		t.Errorf("pulseBrightness(0) = %v, pulseBrightness(period/2) = %v, want dim < bright", dim, bright)
	}
	if dim < 0.25 || dim > 1.0 {
		t.Errorf("pulseBrightness(0) = %v, want within [0.25, 1.0]", dim)
	}
}

func TestRenderDot_RunningEntryPulsesAcrossFrames(t *testing.T) {
	entry := &toolCallEntry{label: "read_solution_file"}
	dim := renderDot(entry, 0)
	bright := renderDot(entry, dotPulsePeriod/2)
	if dim == bright {
		t.Error("expected a running entry's dot color to differ across frames (pulsing), got the same color")
	}
}

func TestRenderDot_DoneEntryHoldsSteadyAcrossFrames(t *testing.T) {
	entry := &toolCallEntry{label: "read_solution_file", done: true}
	a := renderDot(entry, 0)
	b := renderDot(entry, dotPulsePeriod/2)
	if a != b {
		t.Errorf("expected a done entry's dot color to stay constant across frames, got %q then %q", a, b)
	}
}

func TestNewThinkingDisplay_UsesAnsiRendererWithoutKitty(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	defer d.finish()

	if _, ok := d.ball.(ansiBallRenderer); !ok {
		t.Errorf("ball = %T, want ansiBallRenderer", d.ball)
	}
}

func TestNewThinkingDisplay_UsesKittyRendererWhenAvailable(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "1")
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	defer d.finish()

	if _, ok := d.ball.(kittyBallRenderer); !ok {
		t.Errorf("ball = %T, want kittyBallRenderer", d.ball)
	}
}

// TestNewThinkingDisplay_ForcesAnsiWhenBoxActiveEvenWithKittyAvailable is
// a regression test for a real bug found live: with the anchored input
// box active (a DECSTBM-confined scroll region), the Kitty-rendered ball
// visibly intersected with conversation text during ordinary turns —
// whether a placed Kitty image correctly scrolls along with a
// DECSTBM-confined region the way text does isn't verifiable outside a
// real Kitty terminal, so boxInScrollRegion=true forces the
// fully-verified ANSI renderer instead, removing Kitty as a variable
// for that specific combination.
func TestNewThinkingDisplay_ForcesAnsiWhenBoxActiveEvenWithKittyAvailable(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "1")
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, true)
	defer d.finish()

	if _, ok := d.ball.(ansiBallRenderer); !ok {
		t.Errorf("ball = %T, want ansiBallRenderer even though Kitty is available, since a box is active", d.ball)
	}
}

func TestThinkingDisplay_UsesRelativeCursorUpMathNotSaveRestore(t *testing.T) {
	// DECSC/DECRC (save/restore cursor) was tried and reverted — it
	// broke under real scrolling in a live tutor pane (confirmed via
	// screenshot: the block re-printed at a new position each redraw
	// instead of overwriting). Relative cursor-up-by-N is the
	// scroll-safe technique the ANSI renderer always used; this asserts
	// the display never regresses back to DECSC/DECRC.
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)

	ctx := d.onStart(context.Background(), &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackInput{})
	d.onEnd(ctx, &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackOutput{})
	d.finish()

	got := buf.String()
	if strings.Contains(got, "\0337") || strings.Contains(got, "\0338") {
		t.Errorf("expected no DECSC/DECRC sequences, got some in:\n%s", got)
	}
	if !strings.Contains(got, "\033["+fmt.Sprint(discoBallRenderedRows)+"A") {
		t.Errorf("expected a cursor-up-by-%d sequence (the constructor's own row count) somewhere, got:\n%s", discoBallRenderedRows, got)
	}
}

func TestThinkingDisplay_ToolCallsPrintBeforeTheBall(t *testing.T) {
	// The constructor's very first redraw (before any tool call starts)
	// draws the ball alone, since d.calls is still empty then — so the
	// ball's *first-ever* appearance in the whole accumulated buffer
	// naturally precedes the tool name regardless of per-redraw
	// ordering. What actually matters is the ordering *within* a single
	// redraw pass that has both — isolate the last (most complete)
	// redraw segment and check ordering only within that.
	//
	// Captured before finish(): finish() now erases the block (see its
	// doc comment), which appends its own cursor-up sequence after the
	// last real content redraw — this test is about redraw ordering,
	// not erasing, so it looks at the buffer as it stood right after the
	// last real draw.
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	d.onStart(context.Background(), &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackInput{})

	raw := buf.String()
	d.finish()
	ups := regexp.MustCompile(`\x1b\[\d+A`).FindAllStringIndex(raw, -1)
	if len(ups) == 0 {
		t.Fatal("expected at least one cursor-up redraw sequence")
	}
	lastRedraw := raw[ups[len(ups)-1][1]:] // everything after the final cursor-up

	nameIdx := strings.Index(lastRedraw, "read_solution_file")
	if nameIdx == -1 {
		t.Fatal("expected read_solution_file in the final redraw")
	}
	ballIdx := strings.Index(lastRedraw, "\033[38;2;") // ANSI ball's half-block foreground escape
	if ballIdx == -1 {
		t.Skip("ball produced no half-block foreground cells in this frame (possible for a mostly-background frame) — not a reliable signal here")
	}
	if ballIdx < nameIdx {
		t.Error("expected the tool-call name to print before the ball (ball goes below the list), but the ball appeared first")
	}
}

func TestThinkingDisplay_FinishErasesTheBallAndToolCallList(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	d.onStart(context.Background(), &callbacks.RunInfo{Name: "read_solution_file"}, &tool.CallbackInput{})

	beforeFinish := buf.Len()
	d.finish()
	erase := buf.String()[beforeFinish:]

	// One clear per drawn row (tool-call line + ball rows), all wrapped
	// in a cursor-up-by-d.drawn before and after -- verifies finish()
	// actually blanks the block rather than leaving it as a static
	// recap.
	wantUp := fmt.Sprintf("\033[%dA", d.drawn)
	if strings.Count(erase, wantUp) != 2 {
		t.Errorf("expected exactly 2 occurrences of %q (up before and after clearing) in finish()'s output, got %d in:\n%q", wantUp, strings.Count(erase, wantUp), erase)
	}
	if got := strings.Count(erase, "\033[2K"); got != d.drawn {
		t.Errorf("expected %d row clears (\\033[2K) in finish()'s output, got %d in:\n%q", d.drawn, got, erase)
	}
}

func TestThinkingDisplay_FinishClosesTheBallRenderer(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "1")
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	d.finish()

	if !strings.Contains(buf.String(), "a=d") {
		t.Error("expected finish() to close the Kitty renderer (delete its uploaded frames), found no a=d commands")
	}
}

func TestThinkingDisplay_CallbackHandlerFiresOnStartAndOnEnd(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	var buf bytes.Buffer
	d := newThinkingDisplay(&buf, false)
	h := d.callbackHandler()

	info := &callbacks.RunInfo{Name: "read_solution_file", Component: "Tool"}
	ctx := h.OnStart(context.Background(), info, &tool.CallbackInput{})
	h.OnEnd(ctx, info, &tool.CallbackOutput{})
	d.finish()

	got := stripAnsi(buf.String())
	if !strings.Contains(got, "read_solution_file") {
		t.Fatalf("expected the tool name via the built callbacks.Handler, got:\n%s", got)
	}
	d.mu.Lock()
	if len(d.calls) != 1 || !d.calls[0].done {
		t.Errorf("expected the entry to be marked done via the built callbacks.Handler's OnEnd, calls=%+v", d.calls)
	}
	d.mu.Unlock()
}
