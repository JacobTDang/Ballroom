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
// completedAt is the wall-clock time updateLocked set status to
// "done"/"failed" -- zero while still "running". Used by
// activitySettledDotColor to fade the dot from its glow color down to
// the resting base color over activitySettleFadeDuration instead of
// snapping to base the instant a call settles.
type activityCall struct {
	callID, name, status, detail string // status: "running" | "done" | "failed"
	completedAt                  time.Time
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
// activityToolLines by dropping the oldest entry. (started/finished/
// failed used to also return pre-formatted display lines; every caller
// ignored them — display rendering reads currentCalls instead — so the
// string-formatting side of the feed was deleted outright.)
func (f *activityFeed) started(callID, name, argsPreview string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, activityCall{callID: callID, name: name, status: "running", detail: argsPreview})
	if len(f.calls) > activityToolLines {
		f.calls = f.calls[len(f.calls)-activityToolLines:]
	}
}

// finished marks callID done with resultPreview. A callID that isn't
// found (already dropped by started's cap) is a no-op, not an error —
// eino always pairs a real OnStart with the matching OnEnd/OnError, but
// that OnStart's entry may have aged out of the capped list already.
func (f *activityFeed) finished(callID, resultPreview string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updateLocked(callID, "done", resultPreview)
}

// failed marks callID failed with errDetail. Same no-op-if-unknown
// behavior as finished.
func (f *activityFeed) failed(callID, errDetail string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updateLocked(callID, "failed", errDetail)
}

func (f *activityFeed) updateLocked(callID, status, detail string) {
	for i := range f.calls {
		if f.calls[i].callID == callID {
			f.calls[i].status = status
			f.calls[i].detail = detail
			f.calls[i].completedAt = time.Now()
			return
		}
	}
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

// activityPulseBaseR/G/B is the activity dot's resting color — the same
// teal used elsewhere in this project's palette (docker/tmux.conf,
// internal/catalog/theme.go: #2FA6A6), so the tutor pane's own animation
// reads as the same app rather than a mismatched accent. This is the
// pulse's trough, not a floor to dim below (see activityPulseGlow*).
// Kept as untyped constants, not derived from palette at runtime: the
// blend math uses them in both integer and float contexts, which only
// untyped constants satisfy. TestActivityColorsMatchPalette guarantees
// they stay equal to their palette source.
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

// activitySettleFadeDuration is how long a call's dot takes to fade
// from the glow color down to the resting base color once it settles
// (done/failed), instead of snapping to base instantly — per an
// explicit request for smoother transitions ("after tool calling there
// is a slight fade").
const activitySettleFadeDuration = 250 * time.Millisecond

// activitySettledDotColor returns a settled call's dot color for one
// redraw, elapsed time after it completed: the glow color at elapsed=0,
// fading (in Luv, same perceptually-uniform blend activityDotColor
// itself uses) down to the resting base color by
// activitySettleFadeDuration, and exactly the base color from then on.
// A separate function from activityDotColor (rather than folding this
// into it) deliberately: activityDotColor is also called directly by
// pulsedStatusLine, which has no specific call/completion time to fade
// from at all, so giving it a "just settled" mode would mean every
// caller has to reason about a state that doesn't apply to it.
//
// Takes elapsed directly rather than a completedAt timestamp so this
// stays a pure, deterministic function of its input, with no real-clock
// dependency to fake out in tests — pulsedCallLines is the one place
// that actually computes time.Since(c.completedAt) and calls this.
func activitySettledDotColor(elapsed time.Duration) (r, g, b int) {
	if elapsed <= 0 {
		return activityPulseGlowR, activityPulseGlowG, activityPulseGlowB
	}
	if elapsed >= activitySettleFadeDuration {
		return activityPulseBaseR, activityPulseBaseG, activityPulseBaseB
	}
	fade := elapsed.Seconds() / activitySettleFadeDuration.Seconds()
	blended := activityPulseGlowColor.BlendLuv(activityPulseBaseColor, fade)
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
	return coloredGlyph(r, g, b, activityDotGlyph)
}

// coloredGlyph is coloredDot generalized to any glyph — same defensive
// reset-open-content-close span (see coloredDot's comment for why the
// leading plain reset is load-bearing).
func coloredGlyph(r, g, b int, glyph string) string {
	return fmt.Sprintf("\033[0m\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, glyph)
}

// The thinking indicator's trailing dots are a traveling wave rather
// than a static "...": each dot cycles through glyphs of increasing
// visual size, offset from its neighbor by thinkingWaveSpreadTicks, so
// the ripple visibly rolls left-to-right while a turn is in flight.
// Rides the same free-running pulse phase as the dot color, so no new
// tick plumbing. Plain ASCII, like every other glyph in this file (see
// activityDotGlyph's doc comment for the history) -- an earlier version
// of this slice used "·" (U+00B7) and "˙" (U+02D9) for the taller
// frames, which is exactly the class of character this file has already
// been burned by twice. TestActivityRenderingIsPlainASCII is what
// catches a repeat of that.
var thinkingWaveGlyphs = []string{".", "o", "O"}

const (
	thinkingWaveDotCount    = 3
	thinkingWaveSpreadTicks = 6
)

// thinkingWaveLevel returns dot i's height (an index into
// thinkingWaveGlyphs) at the given pulse phase. Pure function of
// phase - i*spread only — that's exactly what makes the wave *travel*:
// dot i+1 replays dot i's motion thinkingWaveSpreadTicks later.
func thinkingWaveLevel(phase, i int) int {
	x := float64(phase-i*thinkingWaveSpreadTicks) / float64(activityPulsePeriodTicks)
	s := math.Sin(2 * math.Pi * x)
	level := int(math.Round((s + 1) / 2 * float64(len(thinkingWaveGlyphs)-1)))
	return min(max(level, 0), len(thinkingWaveGlyphs)-1)
}

// thinkingWaveDots renders the wave for one pulse frame: each dot gets
// its level's glyph, colored by blending base→glow with height, so a
// crest both rises and brightens.
func thinkingWaveDots(phase int) string {
	var b strings.Builder
	for i := 0; i < thinkingWaveDotCount; i++ {
		lvl := thinkingWaveLevel(phase, i)
		frac := float64(lvl) / float64(len(thinkingWaveGlyphs)-1)
		blended := activityPulseBaseColor.BlendLuv(activityPulseGlowColor, frac)
		r8, g8, b8 := blended.RGB255()
		b.WriteString(coloredGlyph(int(r8), int(g8), int(b8), thinkingWaveGlyphs[lvl]))
	}
	return b.String()
}

// pulsedStatusLine builds the activity region's status line for one
// pulse frame: a color-wrapped dot (always animating, driven by
// model.go's free-running pulseTickCmd for as long as a turn is in
// flight), the plain "Thinking" text, then the traveling wave dots.
// Truncation happens on the plain text *before* any color escape is
// added, and the wave is dropped whole when the width can't fit it —
// width-limiting can never slice a truecolor sequence in half.
func pulsedStatusLine(phase, cols int) string {
	const plain = "Thinking"
	r, g, b := activityDotColor("running", phase)
	line := coloredDot(r, g, b) + " " + truncateLine(plain, max(cols-2, 0))
	if cols-2-len(plain) >= thinkingWaveDotCount {
		line += thinkingWaveDots(phase)
	}
	return line
}

// dimSpan wraps s in the pane's dim metadata color (markdown.go's
// mdDimColor), with the same defensive leading reset every colored
// span in this file carries (see coloredDot's doc comment for the
// renderer-diffing bug that makes the reset load-bearing).
func dimSpan(s string) string {
	return "\033[0m" + mdDimColor + s + "\033[0m"
}

// activityCallHeader renders one call's header line body (everything
// after the dot): the tool name in bold, plus — while still running —
// its args dimmed in parens; a failed call's "- failed" flag in the
// error red. It never includes the result/error detail inline — that's
// shown indented beneath the header instead (see activityOutputLines)
// — a real UX fix: a completed call's raw result/JSON used to be
// crammed onto this same line, truncating to almost nothing in a
// normal-width pane.
//
// Takes cols and truncates the plain text itself, each part against
// what's left of the budget, BEFORE any escape is applied — the same
// rule pulsedStatusLine documents: width-limiting can never slice a
// truecolor sequence in half. (Its callers used to truncate the
// returned string instead, which was only safe while this returned
// plain text.)
func activityCallHeader(c activityCall, cols int) string {
	name := truncateLine(c.name, max(cols, 0))
	rest := max(cols-len([]rune(name)), 0)
	styled := mdBoldOn + name + mdBoldOff
	switch c.status {
	case "done":
		return styled
	case "failed":
		if flag := truncateLine(" - failed", rest); flag != "" {
			styled += activityErrorNote(flag)
		}
		return styled
	default: // "running"
		if c.detail == "" || c.detail == "{}" {
			return styled
		}
		if args := truncateLine("("+c.detail+")", rest); args != "" {
			styled += dimSpan(args)
		}
		return styled
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
//
// A settled (done/failed) call's dot uses activitySettledDotColor
// instead of activityDotColor's own static base-color branch, so it
// fades from the glow color down to base over activitySettleFadeDuration
// right after it completes, instead of snapping to base instantly. An
// unset (zero-value) completedAt -- every test in this file that builds
// an activityCall{status: "done", ...} literal directly, never having
// gone through activityFeed's real updateLocked -- safely resolves to
// the plain base color anyway: time.Since of the zero time overflows
// Duration's own range, and time.Time.Sub's documented behavior is to
// clamp to the maximum representable Duration rather than wrap, which
// is still comfortably past activitySettleFadeDuration.
func pulsedCallLines(c activityCall, phase, cols int) []string {
	var r, g, b int
	if c.status == "running" {
		r, g, b = activityDotColor(c.status, phase)
	} else {
		r, g, b = activitySettledDotColor(time.Since(c.completedAt))
	}
	header := coloredDot(r, g, b) + " " + activityCallHeader(c, max(cols-2, 0))
	return append([]string{header}, activityOutputLines(c, cols)...)
}

// settledCallLine renders one call for the permanent post-turn summary
// as a single compact row: resting-teal dot, bold tool name, then the
// first line of its result dimmed inline (or the failure flag and error
// in red). One row per call deliberately — the multi-line indented
// output preview stays a live-region-only affordance (pulsedCallLines),
// where watching a result arrive is useful; baked into the transcript
// forever it read as noise, so the settled record supersedes the old
// indented-output summary with this quieter form.
//
// All truncation happens on plain text against the remaining budget
// BEFORE styling, same escape-safety rule as activityCallHeader; a
// budget too small for a meaningful summary drops the summary whole
// rather than leaving a useless fragment.
func settledCallLine(c activityCall, cols int) string {
	budget := max(cols-2, 0) // dot + space
	name := truncateLine(c.name, budget)
	line := coloredDot(activityPulseBaseR, activityPulseBaseG, activityPulseBaseB) +
		" " + mdBoldOn + name + mdBoldOff
	rest := budget - len([]rune(name))

	detail := strings.SplitN(c.detail, "\n", 2)[0]
	switch {
	case c.status == "failed":
		// The failure flag survives an empty error detail — a failed
		// call must never render indistinguishable from a clean one.
		flag := " failed"
		if detail != "" {
			flag += ": " + detail
		}
		if t := truncateLine(flag, rest); len([]rune(t)) >= len(" failed") {
			line += activityErrorNote(t)
		}
	case detail == "":
		return line
	default: // done
		const sep = " - "
		if s := truncateLine(sep+detail, rest); len([]rune(s)) > len(sep)+3 {
			line += dimSpan(s)
		}
	}
	return line
}

// toolUsageSummary renders a permanent, settled record of which tools a
// completed turn used — one settledCallLine row per call. Unlike the
// live activity region (which disappears entirely once the turn ends,
// see tutorModel.activityView), this gets appended to displayLines so
// the conversation history keeps showing what the model actually did —
// a real feature request from live use ("leave behind the toolname it
// used"). Empty for a turn that made no tool calls at all, so a normal
// reasoning-only turn doesn't get a spurious blank entry.
//
// (The old form reused pulsedCallLines — full indented output previews,
// plus a completedAt-backdating dance so the settle fade wouldn't bake
// mid-glow into the permanent string. settledCallLine always renders at
// the resting base color, so the backdating went with it.)
func toolUsageSummary(calls []activityCall, cols int) string {
	if len(calls) == 0 {
		return ""
	}
	lines := make([]string, len(calls))
	for i, c := range calls {
		lines[i] = settledCallLine(c, cols)
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
