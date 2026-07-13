package tutor

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	agentopt "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	template "github.com/cloudwego/eino/utils/callbacks"
)

// activityThinkingStatus is the activity region's status line while a
// turn is in flight. Phase 1: static text. A later pass replaces this
// with a moving gradient-shimmer redraw (see the plan for this feature)
// — nothing else about the activity region changes for that, since the
// status line is just one more argument to showActivity. Uses only ●
// (a single dot, matching Claude Code's own tool-call indicator) and
// plain ASCII — see formatActivityLine's doc comment for why every other
// symbol this package originally used got removed.
const activityThinkingStatus = "● Thinking..."

// activityArgsPreviewMax/activityResultPreviewMax cap how much of a raw
// tool-call argument/result string appears on one activity line — both
// are already subject to showActivity's own width truncation, but
// truncating here first keeps the preview a readable "at a glance"
// length even in a wide terminal, rather than filling the whole line
// with e.g. a long file's raw content.
const activityArgsPreviewMax = 40
const activityResultPreviewMax = 60

// activityCall is one tool invocation's current display state.
type activityCall struct {
	callID, name, status, detail string // status: "running" | "done" | "failed"
}

// activityFeed tracks the tool calls happening during one Generate call
// (a turn, or a comprehension check) and formats them into the lines
// showActivity displays. A fresh feed is used per call (see
// newActivityOption) — it is not a session-wide log.
type activityFeed struct {
	mu    sync.Mutex
	calls []activityCall
}

// started records a new running call, capping the list at
// activityToolLines by dropping the oldest entry — the same "most
// recent N" trade-off the activity region's fixed row count already
// requires. Returns the current formatted lines under the same lock, so
// a caller's redraw is never built from a state that's already stale by
// the time it reads it.
func (f *activityFeed) started(callID, name, argsPreview string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, activityCall{callID: callID, name: name, status: "running", detail: argsPreview})
	if len(f.calls) > activityToolLines {
		f.calls = f.calls[len(f.calls)-activityToolLines:]
	}
	return f.linesLocked()
}

// finished marks callID done with resultPreview. A callID that isn't
// found (already dropped by started's cap) is a no-op, not an error —
// eino always pairs a real OnStart with the matching OnEnd/OnError, but
// that OnStart's entry may have aged out of the capped list already.
func (f *activityFeed) finished(callID, resultPreview string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updateLocked(callID, "done", resultPreview)
	return f.linesLocked()
}

// failed marks callID failed with errDetail. Same no-op-if-unknown
// behavior as finished.
func (f *activityFeed) failed(callID, errDetail string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updateLocked(callID, "failed", errDetail)
	return f.linesLocked()
}

func (f *activityFeed) updateLocked(callID, status, detail string) {
	for i := range f.calls {
		if f.calls[i].callID == callID {
			f.calls[i].status = status
			f.calls[i].detail = detail
			return
		}
	}
}

func (f *activityFeed) linesLocked() []string {
	lines := make([]string, len(f.calls))
	for i, c := range f.calls {
		lines[i] = formatActivityLine(c)
	}
	return lines
}

// currentCalls returns a copy of the feed's current calls, structured
// (not pre-formatted into strings) — used by activityPulse's ticker,
// which needs each call's status to decide whether its dot pulses or
// sits static (see activityDotColor), not just its rendered text.
func (f *activityFeed) currentCalls() []activityCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]activityCall, len(f.calls))
	copy(out, f.calls)
	return out
}

// activityLineBody renders one call's current state, everything after
// the leading "● " dot (see formatActivityLine/pulsedCallLine, its two
// callers — one plain, one color-wraps the dot). "{}" (an empty JSON
// object -- what eino sends for a no-argument tool) is treated the same
// as no args at all, since showing "({})" on every no-arg tool call
// (most of this package's tools) would be noise, not information.
func activityLineBody(c activityCall) string {
	switch c.status {
	case "done":
		if c.detail == "" {
			return c.name
		}
		return fmt.Sprintf("%s  %s", c.name, c.detail)
	case "failed":
		return fmt.Sprintf("%s - failed: %s", c.name, c.detail)
	default: // "running"
		if c.detail == "" || c.detail == "{}" {
			return c.name
		}
		return fmt.Sprintf("%s(%s)", c.name, c.detail)
	}
}

// formatActivityLine renders one call's current state as plain text —
// used for the immediate, event-driven redraw (see OnStart/OnEnd below);
// activityPulse's ticker uses pulsedCallLine instead, which color-wraps
// the same leading dot. Every line leads with a single ● — a real bug
// found live: the previous version used a different glyph per state
// (⟳ → ✓ ✗) plus a Unicode ellipsis (…), and in a real user's terminal
// font every one of those rendered as an unrecognizable fallback glyph
// (tofu, reading like stray underscores). ● alone is confirmed to render
// correctly everywhere it's been tested (it's an extremely old,
// near-universally-supported code point).
func formatActivityLine(c activityCall) string {
	return "● " + activityLineBody(c)
}

// activityPulseBaseR/G/B is the activity dot's base color — the same
// teal used elsewhere in this project's palette (docker/tmux.conf,
// internal/catalog/theme.go: #2FA6A6), so the tutor pane's own animation
// reads as the same app rather than a mismatched accent.
const (
	activityPulseBaseR = 0x2F
	activityPulseBaseG = 0xA6
	activityPulseBaseB = 0xA6
)

// activityPulseMinBrightness floors how dim a pulsing dot ever gets — a
// full fade to black would read as the dot disappearing/flickering
// rather than "breathing"; capping the low end keeps it always at least
// faintly visible while still giving a clear fade.
const activityPulseMinBrightness = 0.35

// activityPulsePeriodTicks is how many activityPulseInterval ticks make
// up one full fade cycle (dim -> bright -> dim). At the production
// interval (120ms) this is roughly a 1.5s breathing cadence — slow and
// calm, not frantic.
const activityPulsePeriodTicks = 12

// activityDotColor returns the dot's color for one redraw: a smooth
// brightness pulse while status is "running" (the call is actively in
// flight — this is the "fade in and out while it's running" effect),
// or the static full-brightness base color once a call has settled
// (done/failed) — a real design choice, not an oversight: a
// still-pulsing dot next to a call that already finished would read as
// "still working" when it isn't.
func activityDotColor(status string, phase int) (r, g, b int) {
	if status != "running" {
		return activityPulseBaseR, activityPulseBaseG, activityPulseBaseB
	}
	t := float64(phase%activityPulsePeriodTicks) / float64(activityPulsePeriodTicks)
	brightness := activityPulseMinBrightness + (1-activityPulseMinBrightness)*(1+math.Cos(2*math.Pi*t))/2
	return int(float64(activityPulseBaseR) * brightness),
		int(float64(activityPulseBaseG) * brightness),
		int(float64(activityPulseBaseB) * brightness)
}

// coloredDot returns a single ● wrapped in a 24-bit truecolor escape for
// (r,g,b), reset immediately after — tmux.conf already enables truecolor
// passthrough (`terminal-features ",*:RGB"`) for this project's session,
// so this is safe to rely on inside the practice container.
func coloredDot(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm●\033[0m", r, g, b)
}

// pulsedStatusLine builds the activity region's status line for one
// pulse frame: a color-wrapped dot (always animating like a running
// call, for as long as the pulse ticker is alive — see startActivityPulse)
// followed by the plain "Thinking..." text. Truncation happens on the
// plain text *before* the color escape is added, so width-limiting can
// never slice a truecolor sequence in half.
func pulsedStatusLine(phase, cols int) string {
	const plain = "Thinking..."
	r, g, b := activityDotColor("running", phase)
	return coloredDot(r, g, b) + " " + truncateLine(plain, max(cols-2, 0))
}

// pulsedCallLine is pulsedStatusLine's counterpart for one tool-call
// line, color-wrapping activityLineBody's dot per activityDotColor(c.status, phase).
func pulsedCallLine(c activityCall, phase, cols int) string {
	r, g, b := activityDotColor(c.status, phase)
	return coloredDot(r, g, b) + " " + truncateLine(activityLineBody(c), max(cols-2, 0))
}

// buildActivityCallbackOption builds the agentopt.AgentOption that wires
// feed into box's activity display via real eino tool-call callbacks
// (react.BuildAgentCallback / utils/callbacks.ToolCallbackHandler —
// OnStart/OnEnd/OnError fire live, while Generate is still running, not
// after the fact) — the immediate, event-driven redraw path (see
// showActivity's doc comment; startActivityPulse is the other, ambient
// ticker-driven path against the same feed).
//
// compose.GetToolCallID(ctx) is what correlates a specific call's OnEnd
// back to its own OnStart entry (not just by tool name) — eino sets this
// on ctx before invoking the tool, and that same ctx (returned by
// OnStart, per eino's own callback contract) is what OnEnd receives, so
// two concurrent calls to the same tool are tracked as separate entries
// rather than one clobbering the other's status.
func buildActivityCallbackOption(feed *activityFeed, box *inputBox) agentopt.AgentOption {
	toolHandler := &template.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			argsPreview := ""
			if input != nil {
				argsPreview = truncateLine(input.ArgumentsInJSON, activityArgsPreviewMax)
			}
			lines := feed.started(compose.GetToolCallID(ctx), info.Name, argsPreview)
			box.showActivity(activityThinkingStatus, lines)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			resultPreview := ""
			if output != nil {
				resultPreview = truncateLine(output.Response, activityResultPreviewMax)
			}
			lines := feed.finished(compose.GetToolCallID(ctx), resultPreview)
			box.showActivity(activityThinkingStatus, lines)
			return ctx
		},
		OnError: func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
			lines := feed.failed(compose.GetToolCallID(ctx), truncateLine(err.Error(), activityResultPreviewMax))
			box.showActivity(activityThinkingStatus, lines)
			return ctx
		},
	}
	handler := react.BuildAgentCallback(nil, toolHandler)
	return agentopt.WithComposeOptions(compose.WithCallbacks(handler))
}

// activityPulseInterval is a var (not const) so tests can substitute a
// much shorter cadence instead of waiting on the real ~120ms production
// interval — same pattern this package already uses for
// ollamaRequestTimeout.
var activityPulseInterval = 120 * time.Millisecond

// activityPulse is the ticker-driven "fade in and out while it's running"
// animation — a real feature request, not decorative filler: it gives a
// continuous, at-a-glance signal that the tutor is still actively
// working, for however long a turn's Generate call takes, independent of
// whether a tool-call event happens to have just fired.
type activityPulse struct {
	stop chan struct{}
	done chan struct{}
}

// startActivityPulse starts redrawing box's activity region on a fixed
// cadence, reading feed's current state fresh every tick so it always
// paints the real, latest tool-call list (never a stale snapshot) — just
// with an animated dot instead of the static one showActivity's
// event-driven redraws use.
func startActivityPulse(box *inputBox, feed *activityFeed) *activityPulse {
	p := &activityPulse{stop: make(chan struct{}), done: make(chan struct{})}
	go func() {
		defer close(p.done)
		ticker := time.NewTicker(activityPulseInterval)
		defer ticker.Stop()
		phase := 0
		for {
			select {
			case <-p.stop:
				return
			case <-ticker.C:
				phase++
				box.showActivityPulse(feed.currentCalls(), phase)
			}
		}
	}()
	return p
}

// close stops the pulse and blocks until its goroutine has actually
// exited (not just requested to stop) — load-bearing: the caller always
// calls box.clearActivity() right after this returns (see Run()), and
// without this synchronous wait a still-running tick could fire between
// that stop signal and the clear, redrawing stale content on top of a
// region the caller just believed was blank.
func (p *activityPulse) close() {
	close(p.stop)
	<-p.done
}

// startActivitySession is the single entry point Run()/runComprehensionCheck
// use: builds a fresh activityFeed (never shared/reused across calls —
// each turn or comprehension check starts its tool-call window empty,
// same as each turn's own conversation starting fresh context), wires
// both the event-driven callback option and the ambient pulse ticker
// against it, and returns the option to pass into generateWithLeakRetry
// plus a stop func the caller must call once that Generate call returns
// (before box.clearActivity() — see close's doc comment for why the
// order matters).
//
// box == nil (no real terminal, e.g. cmd/tutor-eval) returns a harmless
// no-op option and a no-op stop func — there's no reserved region to
// draw into, matching how every other box-dependent feature in this
// package already degrades.
func startActivitySession(box *inputBox) (agentopt.AgentOption, func()) {
	if box == nil {
		return agentopt.WithComposeOptions(), func() {}
	}
	feed := &activityFeed{}
	opt := buildActivityCallbackOption(feed, box)
	pulse := startActivityPulse(box, feed)
	return opt, pulse.close
}
