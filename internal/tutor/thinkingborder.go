package tutor

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
)

// thinkingBorderColors is the gradient the border cycles through while a
// turn is in flight — the same 7 brand colors as the homepage banner's
// mosaic (internal/catalog/theme.go's bannerGradient), duplicated here
// rather than imported: internal/tutor runs inside the practice-session
// container and deliberately has no dependency on the host-side catalog
// package, and 7 color literals aren't worth creating one. Order
// matters — this is the sequence the chase travels through, not just a
// palette.
var thinkingBorderColors = []colorful.Color{
	borderRGB(0xF0, 0x3C, 0x3C), // red
	borderRGB(0xF0, 0x86, 0x2E), // orange
	borderRGB(0xE8, 0xA9, 0x3C), // gold
	borderRGB(0xE0, 0x46, 0x8C), // pink
	borderRGB(0x9B, 0x5F, 0xB0), // purple
	borderRGB(0x3C, 0x7D, 0xC4), // blue
	borderRGB(0x2F, 0xA6, 0xA6), // teal
}

func borderRGB(r, g, b uint8) colorful.Color {
	return colorful.Color{R: float64(r) / 255.0, G: float64(g) / 255.0, B: float64(b) / 255.0}
}

// thinkingBorderSpeed is how many perimeter cells the gradient advances
// per pulse tick. At activityPulseInterval (40ms) and a typical ~120-col
// terminal, 1 cell/tick is roughly a 3s lap — the same "slow and calm,
// not frantic" cadence activityPulsePeriodTicks establishes for the
// status dot's breathing.
const thinkingBorderSpeed = 1

// thinkingBorderFadeDuration is how long the border takes to disappear
// after a turn completes. Long enough to read as a deliberate fade
// rather than a flicker, short enough that it never lingers into the
// user reading the reply — the border is decoration, and this pane
// already had a permanent decorative border removed once after live
// feedback (see textareaBoxStyle's history), so erring short is the
// safe side.
const thinkingBorderFadeDuration = 900 * time.Millisecond

// thinkingBorderMinWidth/Height are the smallest pane that gets a border
// ring at all — below this, giving up 2 columns and 2 rows to decoration
// would visibly crowd real content, so the whole feature switches off
// (borderInset returns 0 and View never wraps).
const (
	thinkingBorderMinWidth  = 20
	thinkingBorderMinHeight = 8
)

// borderInset is how many cells of the pane's edge are reserved for the
// thinking border's ring: 1 all around on a normally-sized terminal, 0
// when it's too small to spare them. The ring is reserved permanently
// (recomputeLayout sizes the viewport/textarea inside it) rather than
// only while a turn is in flight — the alternative re-wraps the entire
// conversation twice per turn as the content area grows and shrinks by
// two columns, a visible reflow 900ms after every reply lands. At rest
// the reserved cells render as plain spaces with no styling whatsoever
// (see renderThinkingBorder's fade<=0 path), which is what keeps this
// compatible with the earlier "remove the side bar" feedback that killed
// the permanent colored input-box border: blank margin is invisible,
// colored glyphs were the complaint.
func (m tutorModel) borderInset() int {
	if m.width >= thinkingBorderMinWidth && m.height >= thinkingBorderMinHeight {
		return 1
	}
	return 0
}

// innerWidth/innerHeight are the content area's dimensions once the
// border ring is subtracted — every layout computation sizes against
// these, never against m.width/m.height directly, so the layout and the
// border renderer can't disagree about where content ends and ring
// begins.
func (m tutorModel) innerWidth() int  { return m.width - 2*m.borderInset() }
func (m tutorModel) innerHeight() int { return m.height - 2*m.borderInset() }

// thinkingBorderFade is the border's opacity for the current frame:
// 1 for the whole time a turn is in flight, an ease-out decay to 0 over
// thinkingBorderFadeDuration once it settles, and exactly 0 at rest.
// Wall-clock based (like pulsedCallLines' settle fade) rather than
// tick-counted, so the fade duration holds regardless of tick jitter;
// the free-running pulseTickMsg already guarantees a re-render every
// 40ms to actually show the decay.
func (m tutorModel) thinkingBorderFade() float64 {
	if m.turnInFlight {
		return 1.0
	}
	if m.turnSettledAt.IsZero() {
		return 0.0
	}
	elapsed := time.Since(m.turnSettledAt)
	if elapsed >= thinkingBorderFadeDuration {
		return 0.0
	}
	// t*t ease-out: drops quickly at first, tapers at the end — reads
	// as the border "settling" rather than a linear dimmer.
	t := 1.0 - float64(elapsed)/float64(thinkingBorderFadeDuration)
	return t * t
}

// thinkingBorderColorAt returns the gradient color for one perimeter
// cell. pos is the cell's position walking the ring clockwise from the
// top-left corner; phase is the free-running pulse phase. The 7 palette
// stops are stretched across the full perimeter and Luv-blended between
// adjacent stops (same technique and rationale as activityDotColor: RGB
// interpolation between distant hues passes through muddy midpoints,
// Luv doesn't), so the ring reads as one continuous flowing gradient
// that travels as phase advances, not rotating color bands.
func thinkingBorderColorAt(pos, perimeter, phase int) colorful.Color {
	n := len(thinkingBorderColors)
	if perimeter <= 0 {
		return thinkingBorderColors[0]
	}
	shifted := ((pos+phase*thinkingBorderSpeed)%perimeter + perimeter) % perimeter
	scaled := float64(shifted) / float64(perimeter) * float64(n)
	i0 := int(scaled) % n
	i1 := (i0 + 1) % n
	frac := scaled - math.Floor(scaled)
	return thinkingBorderColors[i0].BlendLuv(thinkingBorderColors[i1], frac)
}

// renderThinkingBorder wraps content in the w x h border ring. Every
// border glyph is individually truecolor-colored per
// thinkingBorderColorAt, dimmed toward black by fade — except at
// fade <= 0, where every ring cell is a single plain space with no
// escape sequences at all. That unstyled-at-rest path is deliberate and
// load-bearing: dimming to black would still leave glyphs occupying the
// ring on any terminal theme that isn't pure black, exactly the
// residual-border artifact the old permanent input-box border was
// removed for.
//
// content is padded (or truncated) to exactly fill the interior so the
// bottom edge stays anchored to the pane's last row; lines are padded
// ANSI-aware via lipgloss.Width, since conversation lines carry their
// own styling escapes.
func renderThinkingBorder(content string, w, h, phase int, fade float64) string {
	if w < 4 || h < 4 {
		return content
	}
	innerW, innerH := w-2, h-2
	perimeter := 2*w + 2*h - 4

	lines := strings.Split(content, "\n")
	if len(lines) > innerH {
		lines = lines[:innerH]
	}
	for len(lines) < innerH {
		lines = append(lines, "")
	}

	var b strings.Builder
	// Top edge, left to right: perimeter positions 0..w-1.
	for c := 0; c < w; c++ {
		glyph := "─"
		if c == 0 {
			glyph = "╭"
		} else if c == w-1 {
			glyph = "╮"
		}
		b.WriteString(borderCell(glyph, c, perimeter, phase, fade))
	}
	b.WriteByte('\n')
	// Sides. Walking the ring clockwise from the top-left corner: the
	// right edge continues from the top edge (top-to-bottom), and the
	// left edge is the final stretch back up (bottom-to-top), so the
	// gradient flows around the ring continuously instead of mirroring.
	for r := 1; r <= innerH; r++ {
		leftPos := (2*w + h - 3) + (h - 1 - r)
		rightPos := (w - 1) + r
		line := lines[r-1]
		pad := innerW - lipgloss.Width(line)
		if pad < 0 {
			pad = 0
		}
		b.WriteString(borderCell("│", leftPos, perimeter, phase, fade))
		b.WriteString(line)
		b.WriteString(strings.Repeat(" ", pad))
		b.WriteString(borderCell("│", rightPos, perimeter, phase, fade))
		b.WriteByte('\n')
	}
	// Bottom edge, right to left (continuing clockwise).
	for c := 0; c < w; c++ {
		glyph := "─"
		if c == 0 {
			glyph = "╰"
		} else if c == w-1 {
			glyph = "╯"
		}
		pos := (w - 1) + (h - 1) + (w - 1 - c)
		b.WriteString(borderCell(glyph, pos, perimeter, phase, fade))
	}
	return b.String()
}

// borderCell renders one ring cell: a plain space at rest (see
// renderThinkingBorder's doc comment for why not a black glyph), or the
// glyph in its gradient color blended from black up to full brightness
// by fade.
func borderCell(glyph string, pos, perimeter, phase int, fade float64) string {
	if fade <= 0 {
		return " "
	}
	c := colorful.Color{}.BlendLuv(thinkingBorderColorAt(pos, perimeter, phase), fade)
	r, g, b := c.RGB255()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, glyph)
}
