package tutor

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// activityToolLines caps how many concurrent tool-call lines the
// activity region shows at once, dropping the oldest when a turn makes
// more calls than this — the same "most recent N" trade-off a real
// terminal's limited height always requires, independent of how the
// region itself is laid out (see model.go's recomputeLayout, which sizes
// the region to fit exactly len(activeCalls)+1 rows, capped by this).
const activityToolLines = 4

// activityArgsPreviewMax caps how much of a raw tool-call argument
// string appears inline in a running call's header line (see
// activityCallHeader) — keeps it a readable "at a glance" length even
// in a wide terminal.
const activityArgsPreviewMax = 40

// activityResultPreviewMax caps how much of a raw tool-call result/error
// string is kept at all (see buildActivityChannelOption in model.go),
// before activityOutputLines wraps and further caps it down to
// activityOutputPreviewLines indented lines — generous enough to
// actually fill that 3-line window with real content instead of cutting
// off after what used to fit on one inline line.
const activityResultPreviewMax = 240

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

// activityDotGlyph is the activity indicator's dot — plain ASCII, not a
// Unicode symbol. This project already burned itself once on this exact
// mistake: an earlier version used a different glyph per state
// (⟳ → ✓ ✗) plus a Unicode ellipsis (…), all of which rendered as
// unrecognizable fallback glyphs (tofu) in a real user's terminal font,
// and was replaced with "●" (U+25CF) as "confirmed to render everywhere
// it's been tested." That confirmation turned out to be wrong too — a
// real user later reported ● itself rendering as a bare underscore in
// their terminal/font. There is no Unicode code point this project can
// actually promise will render as a real dot everywhere; plain ASCII
// "o" is the only thing guaranteed to render identically in every
// terminal, encoding, and font, full stop.
const activityDotGlyph = "o"

// activityLineBody renders one call's current state, everything after
// the leading dot (see formatActivityLine/pulsedCallLine, its two
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
// animation. Every line leads with activityDotGlyph.
func formatActivityLine(c activityCall) string {
	return activityDotGlyph + " " + activityLineBody(c)
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
// faintly visible while still giving a clear fade. Raised from an
// earlier, more extreme 0.35: at that floor the dot dimmed to a barely
// perceptible smudge against a black background — a real complaint from
// live use ("I would like a glowing dot not just line") — closer to
// looking flat/off than a glow. 0.6 keeps the dot clearly, consistently
// lit while still leaving real breathing motion (0.6 -> 1.0 -> 0.6).
const activityPulseMinBrightness = 0.6

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

// coloredDot returns activityDotGlyph wrapped in a 24-bit truecolor
// escape for (r,g,b), reset immediately after — tmux.conf already
// enables truecolor passthrough (`terminal-features ",*:RGB"`) for this
// project's session, so this is safe to rely on inside the practice
// container.
func coloredDot(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, activityDotGlyph)
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

// activityCallHeader renders one call's header line body (everything
// after the dot) — the tool name, plus its args in parens while still
// running. Unlike activityLineBody, it never includes the result/error
// detail inline — that's shown indented beneath the header instead (see
// activityOutputLines) — a real UX fix: a completed call's raw
// result/JSON used to be crammed onto this same line, truncating to
// almost nothing in a normal-width pane.
func activityCallHeader(c activityCall) string {
	switch c.status {
	case "done":
		return c.name
	case "failed":
		return c.name + " - failed"
	default: // "running"
		if c.detail == "" || c.detail == "{}" {
			return c.name
		}
		return fmt.Sprintf("%s(%s)", c.name, c.detail)
	}
}

// activityIndent nests a completed/failed call's output preview
// visually beneath the header line that produced it — like Claude
// Code's own tool-result display.
const activityIndent = "  "

// activityOutputPreviewLines caps how many indented output lines a
// completed/failed call's result/error preview shows — a fixed,
// bounded window so one verbose tool call can't push every other active
// call (and the input box below the whole region) off-screen.
const activityOutputPreviewLines = 3

// activityOutputHighlightR/G/B is the tool output preview's highlight
// color — a deliberately loud yellow, per an explicit request from live
// use ("could we like highlight the tool output with yellow for now")
// to make a completed/failed call's result visually unmistakable as
// belonging to the header above it, independent of the indentation
// alone. "for now" in the request itself flags this as provisional.
const (
	activityOutputHighlightR = 0xE6
	activityOutputHighlightG = 0xC3
	activityOutputHighlightB = 0x00
)

// activityOutputHighlight wraps s in the yellow foreground escape,
// same truecolor mechanism as coloredDot.
func activityOutputHighlight(s string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", activityOutputHighlightR, activityOutputHighlightG, activityOutputHighlightB, s)
}

// activityOutputLines returns c's result (done) or error (failed)
// detail, word-wrapped to fit within cols, highlighted in yellow, and
// indented, capped at activityOutputPreviewLines lines — nil for a
// running call (no output yet) or an empty detail. When wrapping
// produces more lines than the cap, the last shown line is re-cut with
// a trailing ellipsis so a long result signals "there's more" rather
// than silently stopping. The indent itself stays outside the color
// escape (activityIndent + activityOutputHighlight(line), not the
// reverse) so callers can still reliably strings.HasPrefix a returned
// line on activityIndent.
func activityOutputLines(c activityCall, cols int) []string {
	if c.detail == "" || (c.status != "done" && c.status != "failed") {
		return nil
	}
	w := cols - len(activityIndent)
	if w <= 0 {
		return nil
	}
	wrapped := strings.Split(lipgloss.NewStyle().Width(w).Render(c.detail), "\n")
	if len(wrapped) > activityOutputPreviewLines {
		wrapped = wrapped[:activityOutputPreviewLines]
		last := activityOutputPreviewLines - 1
		runes := []rune(wrapped[last])
		if cut := w - len(truncateLineEllipsis); cut > 0 && len(runes) > cut {
			runes = runes[:cut]
		}
		wrapped[last] = string(runes) + truncateLineEllipsis
	}
	out := make([]string, len(wrapped))
	for i, line := range wrapped {
		out[i] = activityIndent + activityOutputHighlight(line)
	}
	return out
}

// pulsedCallLines is pulsedStatusLine's counterpart for one tool call: a
// color-wrapped header line (dot + activityCallHeader) followed by the
// call's indented output preview, once it has one (activityOutputLines
// — nil for a still-running call, so it's just the header line).
func pulsedCallLines(c activityCall, phase, cols int) []string {
	r, g, b := activityDotColor(c.status, phase)
	header := coloredDot(r, g, b) + " " + truncateLine(activityCallHeader(c), max(cols-2, 0))
	return append([]string{header}, activityOutputLines(c, cols)...)
}

// toolUsageSummary renders a permanent, settled record of which tools a
// completed turn used, plus each one's indented output preview — reuses
// pulsedCallLines exactly as the live activity region does (phase is
// irrelevant here: activityDotColor only varies by phase for a
// "running" call, and every settled call is "done"/"failed" by the time
// this renders, so it always gets the same static color regardless of
// what phase is passed). Unlike the live activity region (which
// disappears entirely once the turn ends, see tutorModel.activityView),
// this gets appended to displayLines so the conversation history keeps
// showing what the model actually did and what it got back, not just
// its final reply — a real feature request from live use ("leave behind
// the toolname it used" / "the tool output should be indented"). Empty
// for a turn that made no tool calls at all, so a normal reasoning-only
// turn doesn't get a spurious blank entry.
func toolUsageSummary(calls []activityCall, cols int) string {
	if len(calls) == 0 {
		return ""
	}
	var lines []string
	for _, c := range calls {
		lines = append(lines, pulsedCallLines(c, 0, cols)...)
	}
	return strings.Join(lines, "\n")
}

// activityPulseInterval is a var (not const) so tests can substitute a
// much shorter cadence instead of waiting on the real ~120ms production
// interval — same pattern this package already uses for
// ollamaRequestTimeout.
var activityPulseInterval = 120 * time.Millisecond

// truncateLineEllipsis is deliberately plain ASCII, not the Unicode
// ellipsis (…) — a real bug found live: that character (and every other
// symbol this package has tried, see activityDotGlyph's doc comment for
// the full history) rendered as an unrecognizable fallback glyph (tofu,
// reading like a stray underscore) in a real user's terminal font.
// Everything this package writes is plain ASCII, full stop.
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
