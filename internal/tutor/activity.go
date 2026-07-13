package tutor

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// activityToolLines caps how many concurrent tool-call lines the
// activity region shows at once, dropping the oldest when a turn makes
// more calls than this — the same "most recent N" trade-off a real
// terminal's limited height always requires, independent of how the
// region itself is laid out (see model.go's recomputeLayout, which sizes
// the region to fit exactly len(activeCalls)+1 rows, capped by this).
const activityToolLines = 4

// activityArgsPreviewMax/activityResultPreviewMax cap how much of a raw
// tool-call argument/result string appears on one activity line — keeps
// the preview a readable "at a glance" length even in a wide terminal,
// rather than filling the whole line with e.g. a long file's raw
// content.
const activityArgsPreviewMax = 40
const activityResultPreviewMax = 60

// activityCall is one tool invocation's current display state.
type activityCall struct {
	callID, name, status, detail string // status: "running" | "done" | "failed"
}

// activityFeed tracks the tool calls happening during one Generate call
// (a turn, or a comprehension check) and formats them into the lines
// tutorModel's activity region displays. A fresh feed is used per call
// (see model.go's buildActivityChannelOption/startTurn) — it is not a
// session-wide log.
type activityFeed struct {
	mu    sync.Mutex
	calls []activityCall
}

// started records a new running call, capping the list at
// activityToolLines by dropping the oldest entry. Returns the current
// formatted lines under the same lock, so a caller's redraw is never
// built from a state that's already stale by the time it reads it.
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
// (not pre-formatted into strings) — used by model.go's
// buildActivityChannelOption (pushed onto the activity channel on every
// callback) and tutorModel.activityView (pulsedCallLine needs each
// call's status to decide whether its dot pulses or sits static, per
// activityDotColor), not just rendered text.
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
// used by activityFeed's own started/finished/failed for the channel
// snapshot pushed to the bubbletea Update loop; pulsedCallLine (below)
// renders the same body but color-wraps the leading dot for the pulse
// animation. Every line leads with a single ● — a real bug
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
// pulse frame: a color-wrapped dot (always animating, driven by
// model.go's free-running pulseTickCmd for as long as a turn is in
// flight) followed by the plain "Thinking..." text. Truncation happens
// on the plain text *before* the color escape is added, so width-limiting
// can never slice a truecolor sequence in half.
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

// activityPulseInterval is a var (not const) so tests can substitute a
// much shorter cadence instead of waiting on the real ~120ms production
// interval — same pattern this package already uses for
// ollamaRequestTimeout.
var activityPulseInterval = 120 * time.Millisecond

// truncateLineEllipsis is deliberately plain ASCII, not the Unicode
// ellipsis (…) — a real bug found live: that character (and every other
// symbol this package originally used: ⟳ → ✓ ✗) rendered as an
// unrecognizable fallback glyph (tofu, reading like a stray underscore)
// in a real user's terminal font. Everything this package writes is now
// plain ASCII plus the one glyph confirmed to render everywhere: ● (see
// formatActivityLine).
const truncateLineEllipsis = "..."

// truncateLine caps s at max runes, replacing the tail with
// truncateLineEllipsis when it's cut — used for shortening tool-call
// argument/result previews (buildActivityChannelOption in model.go) and
// truncating activity lines to the current terminal width
// (pulsedStatusLine/pulsedCallLine above). max <= 0 returns empty rather
// than panicking on the slice below; max too small to fit the ellipsis
// itself just returns as much of the ellipsis as fits.
func truncateLine(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= len(truncateLineEllipsis) {
		return truncateLineEllipsis[:max]
	}
	return string(runes[:max-len(truncateLineEllipsis)]) + truncateLineEllipsis
}
