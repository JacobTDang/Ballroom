package tutor

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
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

// activityPulseBaseR/G/B is the activity dot's resting color — the same
// teal used elsewhere in this project's palette (docker/tmux.conf,
// internal/catalog/theme.go: #2FA6A6), so the tutor pane's own animation
// reads as the same app rather than a mismatched accent. This is the
// pulse's trough, not a floor to dim below (see activityPulseGlow*).
const (
	activityPulseBaseR = 0x2F
	activityPulseBaseG = 0xA6
	activityPulseBaseB = 0xA6
)

// activityPulseGlowR/G/B is the pulse's peak color — a pale, brighter
// tint of the same teal. Per live feedback ("I would like a glowing dot
// not just line"), the pulse blends toward THIS at its brightest point
// instead of dimming the base color down toward black: it never reads
// as flickering/off, only ever brightening above a color that's already
// clearly lit.
const (
	activityPulseGlowR = 0xBF
	activityPulseGlowG = 0xFC
	activityPulseGlowB = 0xF7
)

// activityPulseBaseColor/activityPulseGlowColor are the same two colors
// above, pre-converted for go-colorful's blend functions.
var (
	activityPulseBaseColor = colorful.Color{R: activityPulseBaseR / 255.0, G: activityPulseBaseG / 255.0, B: activityPulseBaseB / 255.0}
	activityPulseGlowColor = colorful.Color{R: activityPulseGlowR / 255.0, G: activityPulseGlowG / 255.0, B: activityPulseGlowB / 255.0}
)

// activityPulsePeriodTicks is how many activityPulseInterval ticks make
// up one full pulse cycle (glow -> base -> glow). At the production
// interval this is roughly a 1.4s breathing cadence — slow and calm,
// not frantic. 36 ticks (up from an earlier 12, at a proportionally
// shorter interval — see activityPulseInterval) keeps the same overall
// cadence while giving the animation three times as many discrete
// positions to move through, for a visibly smoother motion.
const activityPulsePeriodTicks = 36

// activityDotColor returns the dot's color for one redraw: while status
// is "running" (the call is actively in flight), a smooth glow pulse
// between activityPulseBaseColor (the resting point, at the half-period
// mark) and activityPulseGlowColor (the peak, at phase 0) — blended in
// Luv, a perceptually-uniform color space, via go-colorful, so the
// transition stays visually smooth end to end instead of passing
// through a muddy off-hue midpoint the way plain per-channel RGB
// interpolation can between two different colors. Once a call has
// settled (done/failed), always the static base color — a real design
// choice, not an oversight: a still-pulsing dot next to a call that
// already finished would read as "still working" when it isn't.
func activityDotColor(status string, phase int) (r, g, b int) {
	if status != "running" {
		return activityPulseBaseR, activityPulseBaseG, activityPulseBaseB
	}
	t := float64(phase%activityPulsePeriodTicks) / float64(activityPulsePeriodTicks)
	glow := (1 + math.Cos(2*math.Pi*t)) / 2 // 1 at phase 0 (peak glow), 0 at half period (resting base)
	blended := activityPulseBaseColor.BlendLuv(activityPulseGlowColor, glow)
	r8, g8, b8 := blended.RGB255()
	return int(r8), int(g8), int(b8)
}

// coloredDot returns activityDotGlyph wrapped in a 24-bit truecolor
// escape for (r,g,b), reset immediately after — tmux.conf already
// enables truecolor passthrough (`terminal-features ",*:RGB"`) for this
// project's session, so this is safe to rely on inside the practice
// container.
//
// The color-open escape is prefixed with an explicit plain reset
// (\033[0m) — defensive, not decorative: a real bug found live had a
// later, supposedly-uncolored line (a tool call's plain-text name)
// visibly inheriting an earlier line's color. Every one of this
// function's own escape spans was independently self-contained when
// checked directly (open, content, close, all present, verified via raw
// -e captures), so the leak wasn't from a missing reset here — it's
// bubbletea's real terminal renderer doing incremental, diffed redraws
// across frames (the pulsing activity region redraws every ~120ms while
// idle content around it does not), which this package has no way to
// fully audit without a live terminal. Rather than depend on "the
// previous span definitely reset cleanly" holding across renderer
// internals this package doesn't control, every colored span now
// explicitly forces a clean slate before applying its own color, so it
// can never inherit stray state from whatever rendered immediately
// before it, regardless of the cause.
func coloredDot(r, g, b int) string {
	return fmt.Sprintf("\033[0m\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, activityDotGlyph)
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

// activityOutputHighlightR/G/B is the tool output preview's color — a
// faded gray, per explicit live feedback: yellow was tried first (per
// an earlier request to "highlight the tool output with yellow for
// now") but read as highlighting the wrong thing — the actual ask was
// to visually de-emphasize the raw tool output as secondary/quieter
// text, not draw the eye to it, while the tool call header above it
// stays the normal (unhighlighted) text color.
const (
	activityOutputHighlightR = 0x80
	activityOutputHighlightG = 0x80
	activityOutputHighlightB = 0x80
)

// activityOutputHighlight wraps s in the faded-gray foreground escape,
// same truecolor mechanism as coloredDot -- including the same
// defensive leading reset (see coloredDot's doc comment for why: a real
// bug found live had this color bleeding into the plain-text header
// line above it).
func activityOutputHighlight(s string) string {
	return fmt.Sprintf("\033[0m\033[38;2;%d;%d;%dm%s\033[0m", activityOutputHighlightR, activityOutputHighlightG, activityOutputHighlightB, s)
}

// activityErrorNoteR/G/B is used for a turn or routing-decision
// failure's real underlying error detail, shown directly in the chat --
// the same red already used for a failed check elsewhere in this
// project (internal/tui/boot.go's checkFailStyle: #F03C3C), so a
// failure reads consistently as a failure across the whole app.
const (
	activityErrorNoteR = 0xF0
	activityErrorNoteG = 0x3C
	activityErrorNoteB = 0x3C
)

// activityErrorNote wraps s in the error-red foreground escape, same
// truecolor mechanism (and defensive leading reset) as coloredDot and
// activityOutputHighlight. This is what a turn/routing failure's real
// error detail is rendered with in tutorModel.Update — a real bug found
// live: that detail used to go to a raw fmt.Fprintf(m.stderr, ...) call
// instead, and since a real interactive session has stderr and stdout
// on the very same tty, that write bypassed bubbletea's renderer
// entirely and visibly corrupted the alt-screen frame (stray text
// landing wherever the cursor happened to be, never cleared by the next
// redraw). Routing it through displayLines instead means it goes
// through the same safe, diffed rendering pipeline as everything else
// on screen.
func activityErrorNote(s string) string {
	return fmt.Sprintf("\033[0m\033[38;2;%d;%d;%dm%s\033[0m", activityErrorNoteR, activityErrorNoteG, activityErrorNoteB, s)
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
// much shorter cadence instead of waiting on the real ~40ms production
// interval — same pattern this package already uses for
// ollamaRequestTimeout. Lowered from an earlier 120ms (paired with
// activityPulsePeriodTicks going from 12 to 36, keeping the same overall
// ~1.4s cadence) specifically for smoother-looking motion: three times
// as many redraws per cycle to step through, per a live request to make
// the animation smoother.
var activityPulseInterval = 40 * time.Millisecond

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
