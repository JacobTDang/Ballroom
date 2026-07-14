package tutor

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
)

// The thinking aurora: while a turn is in flight, drifting blooms of
// color glow along the pane's borders and dissolve toward a clean
// center, fading away entirely once the reply lands. This is the third
// shape of the feature, each driven by live feedback: a 1-cell border
// ring read as an outline, and a full-pane background wash read as "you
// just put the image over it" -- static and flat. What the user
// actually wants is light moving around the frame's edges ("the aurora
// should be moving and fading on the borders"), so the field is masked
// to the border region, the interior stays on the terminal's own
// background, and the drift is fast enough to be seen, not just
// measured. At rest it emits nothing at all (overlayAurora's
// brightness<=0 path) -- the no-residue rule inherited from the
// "remove the side bar" feedback that killed a permanent colored
// border.

// auroraBlob is one drifting color source. Its center follows a
// Lissajous path that carries it along and past the pane's edges, so
// blooms of its color slide along the borders, swell as the center
// approaches, and die away as it leaves -- influence falls off as a
// gaussian of distance.
type auroraBlob struct {
	r, g, b float64 // color, 0..1
	ax, ay  float64 // drift amplitudes (wide: centers ride the edges)
	fx, fy  float64 // drift frequencies (rad/s)
	px, py  float64 // phase offsets
	radius  float64 // gaussian falloff radius
	weight  float64 // relative strength
}

// Palette sampled from the reference image: cyan, periwinkle, magenta
// and mint over a deep blue base. Frequencies sit in the 0.4-1.4 rad/s
// band -- full sweeps every handful of seconds, so the motion is
// obvious within a second of watching (the first version's 0.07-0.17
// rad/s was imperceptible, a real complaint) -- and are deliberately
// non-harmonic so the composition never visibly repeats.
var auroraBlobs = []auroraBlob{
	{0x19 / 255.0, 0xE3 / 255.0, 0xF0 / 255.0, 0.60, 0.55, 0.90, 1.40, 0.0, 1.0, 0.35, 1.4}, // cyan
	{0x6C / 255.0, 0x7B / 255.0, 0xE0 / 255.0, 0.70, 0.50, 0.55, 1.05, 2.1, 3.9, 0.40, 1.0}, // periwinkle
	{0xC1 / 255.0, 0x3F / 255.0, 0xD0 / 255.0, 0.55, 0.65, 1.10, 0.70, 4.2, 1.3, 0.28, 0.9}, // magenta
	{0xBF / 255.0, 0xE8 / 255.0, 0xE0 / 255.0, 0.65, 0.55, 0.75, 1.25, 1.0, 5.2, 0.33, 0.8}, // mint
	{0x19 / 255.0, 0xE3 / 255.0, 0xF0 / 255.0, 0.58, 0.68, 1.30, 0.45, 3.3, 2.6, 0.30, 0.7}, // second cyan
}

// auroraBase is the deep blue the blobs float over -- the border glow
// keeps a quiet blue presence even where no bloom is passing, so the
// frame never flickers fully dark mid-turn.
var auroraBase = [3]float64{0x1E / 255.0, 0x50 / 255.0, 0xA2 / 255.0}

// auroraBaseWeight is how strongly the base holds its ground against
// the blobs -- kept low so passing blooms saturate the border with
// their own color instead of always averaging back toward blue.
const auroraBaseWeight = 0.8

// auroraBrightness is the glow's peak strength at the very edge of the
// pane. Higher than the old full-wash value (0.35): the edge mask
// already protects the conversation area, so the border itself can
// afford to be vivid -- that vividness is what makes the effect read
// as light rather than tint.
const auroraBrightness = 0.65

// auroraGlowDepth returns how far the glow reaches into the pane, in
// column units (rows count double -- a terminal cell is roughly twice
// as tall as wide, so a symmetric cell count would look twice as deep
// vertically). Scales with pane width so the glow stays proportionate,
// clamped so tiny panes still show a sliver and huge panes don't drown
// half the conversation.
func auroraGlowDepth(w int) float64 {
	d := float64(w) / 6.0
	if d < 6 {
		return 6
	}
	if d > 16 {
		return 16
	}
	return d
}

// auroraDepthWave modulates how deep the glow reaches at a given point
// and moment: a sum of slow traveling waves, incommensurate in both
// frequency and speed so the undulation never settles into a visible
// repeat. Normalized to [0,1]. This is what keeps the glow's inner
// boundary alive -- a constant depth read as a rigid box frame (real
// feedback: "rn its just a box"), while a boundary that swells and
// recedes as tongues of light travel the border reads as aurora.
func auroraDepthWave(u, v, t float64) float64 {
	s := math.Sin(u*7.3+t*1.1) +
		math.Sin(u*3.1-t*0.7+1.7) +
		math.Sin(v*5.9+t*0.9+3.9) +
		math.Sin(v*2.3-t*1.3+0.6)
	return (s/4 + 1) / 2
}

// auroraEdgeMask is the glow's intensity at a cell: 1 at the pane's
// outermost cells, dissolving to 0 at a locally undulating depth (0.5x
// to 1.2x the base glowDepth, per auroraDepthWave). The falloff is a
// smoothstep rather than a plain quadratic: zero slope at both ends
// gives a brightness plateau right at the edge and a completely soft
// entry at the inner boundary, so there's no visible line where the
// glow begins. This mask is the whole difference between "border glow"
// and "background wash" -- interior cells (mask 0) are never painted
// at all, so the conversation sits on the terminal's own background.
func auroraEdgeMask(col, row, w, h int, depth, t float64) float64 {
	d := math.Min(float64(col), float64(w-1-col))
	d = math.Min(d, 2*float64(row))
	d = math.Min(d, 2*float64(h-1-row))
	u := float64(col) / math.Max(float64(w-1), 1)
	v := float64(row) / math.Max(float64(h-1), 1)
	local := depth * (0.5 + 0.7*auroraDepthWave(u, v, t))
	if d >= local {
		return 0
	}
	s := 1 - d/local
	return s * s * (3 - 2*s)
}

// auroraFadeDuration is how long the glow takes to disappear after a
// turn completes -- long enough to read as a deliberate fade, short
// enough to never linger into reading the reply.
const auroraFadeDuration = 900 * time.Millisecond

// auroraColorAt is the field color at normalized pane coordinates
// (u,v) and animation time t (seconds), before brightness or edge
// masking. Weighted gaussian blend of every blob over the base.
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

// auroraFade is the glow's opacity for the current frame: 1 for the
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

// auroraPaintFloor is the effective brightness below which a cell is
// left unpainted entirely rather than tinted imperceptibly -- both the
// clean-interior guarantee and most of the per-frame escape volume
// hang off this cutoff.
const auroraPaintFloor = 0.02

// overlayAurora paints the border glow behind content: cells within
// the edge mask get a truecolor background sampled from auroraColorAt
// (scaled by the mask and brightness), interior cells are left
// completely untouched, and the content's own glyphs and foreground
// styling pass through as-is. brightness <= 0 returns content
// byte-identical -- the no-residue guarantee, as a genuinely separate
// code path.
//
// Content is walked with ansi.DecodeSequence, which splits styled text
// into escape sequences (passed through) and grapheme clusters (each a
// cell of known width), so wide characters and the conversation's own
// styling can't desync the column accounting. Every painted->unpainted
// transition emits a background reset -- without it the glow color
// would smear across the whole line through attribute inheritance.
// Lines never carry their own background styling in this pane
// (everything is fg-styled), which is what makes injection safe.
func overlayAurora(content string, w, h int, t, brightness float64) string {
	if brightness <= 0 || w <= 0 || h <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	// Pad to full height so the glow reaches the pane's bottom edge
	// below short content; never truncate -- a background paints
	// behind content, it doesn't clip it.
	for len(lines) < h {
		lines = append(lines, "")
	}
	rows := len(lines)
	depth := auroraGlowDepth(w)

	var out strings.Builder
	parser := ansi.NewParser()
	for row, line := range lines {
		if row > 0 {
			out.WriteByte('\n')
		}
		col := 0
		painted := false
		var state byte
		emitCell := func() {
			eb := brightness * auroraEdgeMask(col, row, w, rows, depth, t)
			if eb >= auroraPaintFloor {
				out.WriteString(auroraBG(col, row, w, rows, t, eb))
				painted = true
			} else if painted {
				out.WriteString("\x1b[49m")
				painted = false
			}
		}
		for len(line) > 0 {
			seq, width, n, newState := ansi.DecodeSequence(line, state, parser)
			if width > 0 {
				emitCell()
				col += width
			}
			out.WriteString(seq)
			state = newState
			line = line[n:]
		}
		for ; col < w; col++ {
			emitCell()
			out.WriteByte(' ')
		}
		// Restore the default background at each line's edge so the
		// glow can't bleed past the pane via the terminal's
		// background-color-erase behavior.
		if painted {
			out.WriteString("\x1b[49m")
		}
	}
	return out.String()
}

// auroraBG is the background escape for one cell at effective
// brightness eb (edge mask and fade already applied).
func auroraBG(col, row, w, h int, t, eb float64) string {
	u := float64(col) / math.Max(float64(w-1), 1)
	v := float64(row) / math.Max(float64(h-1), 1)
	r, g, b := auroraColorAt(u, v, t)
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", clamp255(r*eb), clamp255(g*eb), clamp255(b*eb))
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
