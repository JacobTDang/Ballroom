package tutor

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
)

// The thinking aurora: while a turn is in flight, the whole pane's
// background becomes a slowly drifting mesh gradient -- overlapping
// soft color blobs, like light moving under water -- that fades away
// once the reply lands. It replaced an earlier 1-cell gradient border
// ring: the user wanted the full floating-gradient look (reference: a
// blurry blue/cyan/magenta mesh wallpaper), not an outline. Because
// it's a background, it occupies no cells and the layout never changes;
// and at rest it emits nothing at all (overlayAurora's brightness<=0
// path), the same no-residue rule the ring obeyed, inherited from the
// "remove the side bar" feedback that killed a permanent colored border.

// auroraBlob is one drifting color source. Its center follows a slow
// Lissajous path through normalized pane coordinates, and its influence
// at any point falls off as a gaussian of distance -- summing a handful
// of these over a deep-blue base is what produces the soft mesh look,
// with no seams for the eye to catch.
type auroraBlob struct {
	r, g, b float64 // color, 0..1
	ax, ay  float64 // drift amplitudes
	fx, fy  float64 // drift frequencies (rad/s)
	px, py  float64 // phase offsets
	radius  float64 // gaussian falloff radius
	weight  float64 // relative strength
}

// Palette sampled from the reference image: cyan, periwinkle, magenta
// and mint drifting over a deep blue base. Frequencies are deliberately
// non-harmonic so the composition never visibly repeats.
var auroraBlobs = []auroraBlob{
	{0x19 / 255.0, 0xE3 / 255.0, 0xF0 / 255.0, 0.35, 0.30, 0.11, 0.17, 0.0, 1.0, 0.45, 1.4}, // cyan
	{0x6C / 255.0, 0x7B / 255.0, 0xE0 / 255.0, 0.40, 0.25, 0.07, 0.13, 2.1, 3.9, 0.50, 1.0}, // periwinkle
	{0xC1 / 255.0, 0x3F / 255.0, 0xD0 / 255.0, 0.30, 0.35, 0.13, 0.09, 4.2, 1.3, 0.30, 0.9}, // magenta
	{0xBF / 255.0, 0xE8 / 255.0, 0xE0 / 255.0, 0.38, 0.28, 0.09, 0.15, 1.0, 5.2, 0.40, 0.8}, // mint
	{0x19 / 255.0, 0xE3 / 255.0, 0xF0 / 255.0, 0.33, 0.36, 0.15, 0.08, 3.3, 2.6, 0.35, 0.7}, // second cyan
}

// auroraBase is the deep blue the blobs float over -- gaps between
// blobs read as ocean, never as black holes in the field.
var auroraBase = [3]float64{0x1E / 255.0, 0x50 / 255.0, 0xA2 / 255.0}

// auroraBaseWeight is how strongly the base holds its ground against
// the blobs -- higher flattens the field toward uniform blue, lower
// lets blob colors saturate their neighborhoods.
const auroraBaseWeight = 1.2

// auroraBrightness caps the field so bright conversation text stays
// crisp on top of it -- at full strength the reference-image colors
// would drown the foreground. Verified against real chat lines in the
// prototype: ~0.35 reads as a dark animated aurora with white text
// floating cleanly over it.
const auroraBrightness = 0.35

// auroraFadeDuration is how long the background takes to disappear
// after a turn completes -- same value and reasoning as the border ring
// it replaced: long enough to read as a deliberate fade, short enough
// to never linger into reading the reply.
const auroraFadeDuration = 900 * time.Millisecond

// auroraColorAt is the field color at normalized pane coordinates
// (u,v) and animation time t (seconds), before any brightness scaling.
// Weighted gaussian blend of every blob over the base.
func auroraColorAt(u, v, t float64) (r, g, b float64) {
	r = auroraBase[0] * auroraBaseWeight
	g = auroraBase[1] * auroraBaseWeight
	b = auroraBase[2] * auroraBaseWeight
	wsum := auroraBaseWeight
	for _, bl := range auroraBlobs {
		cx := 0.5 + bl.ax*math.Sin(bl.fx*t+bl.px)
		cy := 0.5 + bl.ay*math.Cos(bl.fy*t+bl.py)
		dx, dy := u-cx, v-cy
		w := bl.weight * math.Exp(-(dx*dx+dy*dy)/(bl.radius*bl.radius))
		r += bl.r * w
		g += bl.g * w
		b += bl.b * w
		wsum += w
	}
	return r / wsum, g / wsum, b / wsum
}

// auroraFade is the field's opacity for the current frame: 1 for the
// whole time a turn is in flight, an ease-out decay to 0 over
// auroraFadeDuration once it settles, exactly 0 at rest. Wall-clock
// based rather than tick-counted so the duration holds regardless of
// tick jitter; the free-running pulseTickMsg guarantees a re-render
// every 40ms to actually show the decay.
func (m tutorModel) auroraFade() float64 {
	if m.turnInFlight {
		return 1.0
	}
	if m.turnSettledAt.IsZero() {
		return 0.0
	}
	elapsed := time.Since(m.turnSettledAt)
	if elapsed >= auroraFadeDuration {
		return 0.0
	}
	// t*t ease-out: drops quickly at first, tapers at the end.
	t := 1.0 - float64(elapsed)/float64(auroraFadeDuration)
	return t * t
}

// overlayAurora paints the animated field behind content: every
// printable cell (and every padding cell out to the full w x h pane)
// gets a truecolor background escape sampled from auroraColorAt, while
// the content's own glyphs and foreground styling pass through
// untouched. brightness <= 0 returns content byte-identical -- the
// no-residue guarantee lives here, as a genuinely separate code path
// rather than a background dimmed to black (which would still tint the
// pane on any terminal theme that isn't pure black).
//
// Content is walked with ansi.DecodeSequence, which splits styled text
// into escape sequences (passed through as-is) and grapheme clusters
// (each one a cell of known width) -- the same tokenizer the Charm
// stack itself uses, so wide characters and the conversation's own
// styling can't desync the column accounting. A content reset (\x1b[0m)
// clears our injected background too, but the very next cell re-emits
// it, so no gap ever renders. Lines never carry their own background
// styling in this pane (everything is fg-styled), which is what makes
// blind per-cell injection safe.
func overlayAurora(content string, w, h int, t, brightness float64) string {
	if brightness <= 0 || w <= 0 || h <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	// Pad to full height so the field fills the pane below short
	// content; never truncate -- a background paints behind content,
	// it doesn't clip it.
	for len(lines) < h {
		lines = append(lines, "")
	}
	rows := len(lines)

	var out strings.Builder
	parser := ansi.NewParser()
	for row, line := range lines {
		if row > 0 {
			out.WriteByte('\n')
		}
		col := 0
		var state byte
		for len(line) > 0 {
			seq, width, n, newState := ansi.DecodeSequence(line, state, parser)
			if width > 0 {
				out.WriteString(auroraBG(col, row, w, rows, t, brightness))
				col += width
			}
			out.WriteString(seq)
			state = newState
			line = line[n:]
		}
		for ; col < w; col++ {
			out.WriteString(auroraBG(col, row, w, rows, t, brightness))
			out.WriteByte(' ')
		}
		// Restore the default background at each line's edge so the
		// field can't bleed past the pane via the terminal's
		// background-color-erase behavior.
		out.WriteString("\x1b[49m")
	}
	return out.String()
}

// auroraBG is the background escape for one cell.
func auroraBG(col, row, w, h int, t, brightness float64) string {
	u := float64(col) / math.Max(float64(w-1), 1)
	v := float64(row) / math.Max(float64(h-1), 1)
	r, g, b := auroraColorAt(u, v, t)
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm",
		clamp255(r*brightness), clamp255(g*brightness), clamp255(b*brightness))
}

func clamp255(x float64) int {
	v := int(x * 255.0)
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
