package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"sync"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent/react"
	template "github.com/cloudwego/eino/utils/callbacks"
)

// thinkingTickInterval matches internal/tui/tick.go's tickInterval, for
// a consistent spin rate with the rest of the app.
const thinkingTickInterval = 150 * time.Millisecond

// dotBaseColor is the tool-call progress dot's accent color — the same
// warm gold internal/tui uses for its hint/highlight style (#E8A93C),
// reused here so it reads as this app's established "in progress" color
// rather than introducing a new one.
var dotBaseColor = color.RGBA{R: 232, G: 169, B: 60, A: 255}

// dotPulsePeriod is how many ticks one full dim->bright->dim pulse takes
// (10 ticks * 150ms = 1.5s per cycle).
const dotPulsePeriod = 10

// toolCallEntry is one tool call's display state: its label (name plus
// an optional short argument summary) and whether it has finished. The
// dot next to an in-progress entry pulses; once done it holds steady at
// full brightness — the pulse stopping is the "complete" signal, no
// separate glyph needed.
type toolCallEntry struct {
	label string
	done  bool
}

// thinkingCallKey correlates one onStart call to its matching onEnd call
// via the context eino threads between them (see compose/utils.go's
// runWithCallbacks: the context onStart returns is what's passed to the
// tool call and then to onEnd) — works correctly even when several tool
// calls are in flight concurrently, since each has its own derived
// context carrying its own entry pointer.
type thinkingCallKey struct{}

// discoBallRenderer abstracts how the ball itself actually gets drawn —
// ansiBallRenderer (discoball.go, the half-block mosaic, always
// correct, works in any terminal) or kittyBallRenderer (kittyimage.go,
// the real Kitty graphics protocol, genuine full-quality image data
// instead of an approximation, only usable when kittyAvailable()).
// newThinkingDisplay picks one at construction time; the tool-call dot
// list underneath is always plain ANSI regardless — there's no reason
// to risk that on a still-best-effort image path.
type discoBallRenderer interface {
	// init runs once, before the first showFrame call.
	init(w io.Writer)
	// showFrame draws frame (not yet reduced mod the frame count) at
	// the writer's current cursor position, and must leave the cursor
	// exactly rows() lines below where it started — the same contract
	// the ANSI renderer has always had (it prints exactly rows() \n-
	// terminated lines), which redrawLocked's cursor-up math depends on.
	showFrame(w io.Writer, frame int)
	// rows is how many terminal lines this renderer's ball occupies —
	// fixed for its lifetime.
	rows() int
	// close releases any renderer-owned resources, once, when the
	// display finishes.
	close(w io.Writer)
}

// thinkingDisplay renders a small spinning disco ball in the tutor pane
// for the whole duration of one turn — from the moment it's constructed
// (right before agent.Generate) until finish is called — with each tool
// call's name appended as a line below the ball, prefixed by a dot that
// pulses while it's running and holds steady once it's done. Replaces
// the old toolCallDisplay (per-tool spinner/checkmark lines), which
// rendered as literal underscores in the real tutor pane: the Braille
// spinner and checkmark glyphs weren't in that terminal's font. Colored
// background blocks sidestep that entirely — there's no glyph to look
// up, just terminal cell background painting.
//
// Tool calls within a single turn can run concurrently — eino's
// ToolsNode executes them via goroutines when a model requests several
// at once (verified by reading compose/tool_node.go) — so every
// mutation goes through mu and redrawLocked.
type thinkingDisplay struct {
	mu    sync.Mutex
	w     io.Writer
	frame int
	calls []*toolCallEntry
	ball  discoBallRenderer
	drawn int // lines currently on screen, for the cursor-up redraw math

	stop       chan struct{}
	finishOnce sync.Once
}

// newThinkingDisplay starts the ball animating immediately — unlike the
// old toolCallDisplay, which only activated on a tool call, this runs
// for every turn (including ones where the model never calls a tool at
// all), matching "an animation for when the tutor is thinking" literally.
//
// boxInScrollRegion forces the ANSI renderer even when Kitty is
// available — a real bug found live, not theorized: with the anchored
// input box active (a DECSTBM-confined scroll region), the ball visibly
// intersected with conversation text during ordinary turns, not just on
// resize. The Kitty graphics protocol was designed assuming a normal,
// full-screen-scrolling terminal; whether a placed image correctly
// scrolls along with a DECSTBM-confined region the way text does, or
// stays visually fixed while text scrolls past/through it, isn't
// something verifiable from this environment (it needs a real Kitty
// terminal). The ANSI renderer is fully verified working correctly
// inside the box across many real multi-turn tests, so this removes
// Kitty as a variable for that combination specifically rather than
// keep guessing at undocumented protocol behavior. Kitty rendering is
// unaffected when no box is active (already proven working live,
// unrelated to this bug).
func newThinkingDisplay(w io.Writer, boxInScrollRegion bool) *thinkingDisplay {
	var ball discoBallRenderer = ansiBallRenderer{}
	if kittyAvailable() && !boxInScrollRegion {
		ball = kittyBallRenderer{}
	}
	d := &thinkingDisplay{w: w, ball: ball, stop: make(chan struct{})}
	d.ball.init(w)
	d.mu.Lock()
	d.redrawLocked()
	d.mu.Unlock()
	go d.tick()
	return d
}

// callbackHandler builds the eino callbacks.Handler that drives this
// display. Only a tool handler is registered (no chat-model handler) —
// verified safe by reading eino's dispatch code (internal/callbacks/inject.go
// filters handlers via Needed() before ever calling OnStart/OnEnd, and
// handlerTemplate.Needed nil-checks each sub-handler before delegating),
// so the react agent's own chat-model callbacks simply never reach this
// handler rather than risking a nil-pointer dispatch.
func (d *thinkingDisplay) callbackHandler() callbacks.Handler {
	return react.BuildAgentCallback(nil, &template.ToolCallbackHandler{
		OnStart: d.onStart,
		OnEnd:   d.onEnd,
	})
}

func (d *thinkingDisplay) onStart(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
	label := info.Name
	if summary := summarizeArgs(info.Name, input.ArgumentsInJSON); summary != "" {
		label = info.Name + " " + summary
	}
	entry := &toolCallEntry{label: label}

	d.mu.Lock()
	d.calls = append(d.calls, entry)
	d.redrawLocked()
	d.mu.Unlock()

	return context.WithValue(ctx, thinkingCallKey{}, entry)
}

func (d *thinkingDisplay) onEnd(ctx context.Context, _ *callbacks.RunInfo, _ *tool.CallbackOutput) context.Context {
	entry, _ := ctx.Value(thinkingCallKey{}).(*toolCallEntry)

	d.mu.Lock()
	if entry != nil {
		entry.done = true
	}
	d.redrawLocked()
	d.mu.Unlock()

	return ctx
}

func (d *thinkingDisplay) tick() {
	ticker := time.NewTicker(thinkingTickInterval)
	defer ticker.Stop()
	for {
		select {
		case <-d.stop:
			return
		case <-ticker.C:
			d.mu.Lock()
			d.frame++
			d.redrawLocked()
			d.mu.Unlock()
		}
	}
}

// finish stops the animation — call via defer right after starting a
// display for a turn, so the ticker goroutine never outlives that turn.
// Erases the ball and tool-call list once the turn is done, so the
// reply prints where they were instead of below a static recap left on
// screen — an earlier version deliberately left the block in place, but
// the user found that made a long practice session's tool-call block
// pile up behind every single reply, which reads as clutter, not a
// useful recap. Also releases the renderer's own resources (for Kitty:
// deletes its uploaded image data) before erasing, so a Kitty placement
// is torn down through its own protocol rather than just painted over.
// Both happen once, together, inside finishOnce so a second finish()
// call (tutor.go doesn't do this, but the type stays correct on its own
// terms) is a safe no-op.
func (d *thinkingDisplay) finish() {
	d.finishOnce.Do(func() {
		close(d.stop)
		d.mu.Lock()
		d.ball.close(d.w)
		d.eraseLocked()
		d.mu.Unlock()
	})
}

// eraseLocked blanks every row redrawLocked last drew and leaves the
// cursor back at the top of that now-empty space, ready for the turn's
// reply to print there. Must be called with mu held.
func (d *thinkingDisplay) eraseLocked() {
	if d.drawn == 0 {
		return
	}
	fmt.Fprintf(d.w, "\033[%dA", d.drawn)
	for i := 0; i < d.drawn; i++ {
		io.WriteString(d.w, "\033[2K\r\n")
	}
	fmt.Fprintf(d.w, "\033[%dA", d.drawn)
}

// redrawLocked repaints the whole block in place: one line per tool
// call (each prefixed with its progress dot) followed by the current
// ball frame below it. Must be called with mu held. Every caller (the
// constructor, onStart, onEnd, tick) goes through this same locked
// path, so concurrent tool calls just serialize into a consistent
// repaint rather than a torn/interleaved write.
//
// Moves the cursor up by the exact row count drawn last time, then
// reprints everything — the same technique the ANSI renderer always
// used, scroll-safe as long as the row count is tracked and used
// consistently (content and cursor scroll together at the bottom
// margin, so relative movement stays correct regardless of where
// scrolling has put the block on the physical screen). A DECSC/DECRC
// (save/restore cursor) version was tried instead — save once, restore
// to that exact point every redraw — but broke under real scrolling in
// a live tutor pane (confirmed from a screenshot: the block re-printed
// at a new position each redraw instead of overwriting, since DECRC's
// restored position doesn't track a scroll that happened after the
// save). The standalone checker that seemed to validate DECSC/DECRC
// ran in a mostly-empty terminal window and never exercised scrolling at
// all, which is why the bug didn't show up until a real, already-full
// tutor pane.
func (d *thinkingDisplay) redrawLocked() {
	if d.drawn > 0 {
		fmt.Fprintf(d.w, "\033[%dA", d.drawn)
	}
	for _, entry := range d.calls {
		fmt.Fprintf(d.w, "\033[2K\r  %s %s\n", renderDot(entry, d.frame), entry.label)
	}
	d.ball.showFrame(d.w, d.frame)
	d.drawn = len(d.calls) + d.ball.rows()
}

// renderDot renders one tool call's progress indicator: a single
// truecolor-background space, pulsing in brightness while entry is
// still running, held at full brightness once it's done. Same
// colored-space technique as the disco ball itself (renderDiscoBallFrame)
// — no glyph lookup, so nothing here can hit the font-coverage gap that
// broke the previous spinner/checkmark indicator.
func renderDot(entry *toolCallEntry, frame int) string {
	brightness := 1.0
	if !entry.done {
		brightness = pulseBrightness(frame)
	}
	c := dotBaseColor
	return fmt.Sprintf("\033[48;2;%d;%d;%dm \033[0m",
		uint8(float64(c.R)*brightness),
		uint8(float64(c.G)*brightness),
		uint8(float64(c.B)*brightness),
	)
}

// pulseBrightness is a triangle wave over [0.25, 1.0] with period
// dotPulsePeriod ticks — never fully dark, so the dot stays visible even
// at its dimmest point.
func pulseBrightness(frame int) float64 {
	t := frame % dotPulsePeriod
	half := dotPulsePeriod / 2
	var tri float64
	if t < half {
		tri = float64(t) / float64(half)
	} else {
		tri = float64(dotPulsePeriod-t) / float64(half)
	}
	return 0.25 + 0.75*tri
}

// summarizeArgs returns a short human-readable summary of a tool call's
// arguments for display, or "" for tools with nothing worth showing
// (read_solution_file, read_problem_statement, read_test_output, and
// read_cursor_position all take no model-supplied arguments). This is
// purely cosmetic — a malformed argument string just falls back to no
// summary rather than surfacing an error, since the tool call's actual
// execution (with its own error handling via
// utils.WrapToolWithErrorHandler) is unaffected either way.
func summarizeArgs(toolName, argumentsJSON string) string {
	if toolName != "highlight_lines" {
		return ""
	}
	var in highlightLinesInput
	if err := json.Unmarshal([]byte(argumentsJSON), &in); err != nil {
		return ""
	}
	return fmt.Sprintf("(lines %d-%d)", in.Start, in.End)
}
