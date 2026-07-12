package tutor

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"strings"
)

//go:embed discoballassets/*.png
var discoBallFS embed.FS

// discoBallSampleRows/Cols set the downsample grid resolution — how much
// detail survives from the real sprite. Equal to each other (not the
// old 2:1 width:height ratio) because renderDiscoBallFrame packs two
// sample rows into one physical terminal line via half-block characters
// (see its doc comment): each half-block glyph's top/bottom halves are
// roughly square, so a square sample grid is what maps back to a
// circle. A first pass sampled directly at physical-line resolution
// (one sample row per terminal line, doubled cell width for aspect
// correction) — confirmed in a real practice session to look like an
// unrecognizable color blob at a small physical size, since so few
// samples couldn't resolve the sprite's facet pattern at all. Packing
// two samples per line recovers real detail (discoBallSampleRows/2
// physical lines instead of discoBallSampleRows) without growing the
// on-screen footprint. Shrunk to 6x6 (3 rendered lines) at one point to
// make the indicator less prominent, but that was tuned while Kitty
// (a real image, crisp at any size) was still the renderer actually in
// use live — a since-reverted change briefly forced ANSI whenever the
// box was active (thinkingdisplay.go's newThinkingDisplay), and 6x6 read
// as an undifferentiated color blob under that path in real live
// testing, not a ball. Back to 10x10 (5 rendered lines), the last size
// confirmed live to still look
// like a ball (that round's feedback was "too big/prominent", not
// "unrecognizable") — prioritizing recognizable over compact.
const (
	discoBallFrameCount = 16
	discoBallSampleRows = 10
	discoBallSampleCols = 10

	// discoBallAlphaThreshold: a downsampled cell averaging below this
	// (out of the 65535-scale alpha image.Color.RGBA() returns) is
	// treated as background and left unpainted, rather than rendering a
	// faint/dark smear at the ball's edge. Only catches genuinely
	// transparent backgrounds (the first sprite this package used).
	discoBallAlphaThreshold = 65535 / 10

	// discoBallBrightnessThreshold: a fully opaque cell whose R+G+B sum
	// (each 0-255) falls below this is also treated as background — the
	// current sprite fills its background with solid opaque black
	// (alpha=255 everywhere) rather than transparency, so alpha alone
	// can't distinguish ball from background. 12 is deliberately tight:
	// it only catches cells that are essentially pure black, not merely
	// dark ball facets (a dark purple or navy tile is still well above
	// this). The sprite's own darkest frames (fading in/out) legitimately
	// render as mostly-or-entirely background under this rule — that's
	// the source content, not a bug.
	discoBallBrightnessThreshold = 12
)

// discoBallRenderedRows is how many physical terminal lines
// renderDiscoBallFrame actually produces — half discoBallSampleRows,
// since each line packs two sample rows via a half-block glyph. Callers
// doing cursor-math (thinkingdisplay.go's redrawLocked) need this, not
// discoBallSampleRows.
const discoBallRenderedRows = (discoBallSampleRows + 1) / 2

// discoBallFrame is one animation frame, downsampled to
// discoBallSampleRows x discoBallSampleCols cells. A nil cell is
// background (transparent in the source sprite) — left unpainted at
// render time so it blends with the pane's real background instead of a
// solid black square.
type discoBallFrame [][]*color.RGBA

// discoBallFrames holds all 16 frames, decoded and downsampled once at
// package init — cheap (a handful of ~300-500KB PNGs), so there's no
// need for lazy loading. Sliced (see discoballassets/) from a 4x4 sheet
// that fades in from dark (frame 0) to full sparkle (~frames 5-13) and
// back down to dark (frames 14-15) — not a continuous rotation like the
// first sprite this package used. Looping frame 15 back to 0 therefore
// reads as a pulse/breathe rather than a spin; that's the sheet's own
// content, not something to correct for.
var discoBallFrames = loadDiscoBallFrames()

// loadDiscoBallFrames decodes each embedded sprite frame and downsamples
// it. Panics on failure — these are compiled-in assets, so a failure
// here means the embed itself is broken, not something a caller can
// meaningfully recover from.
func loadDiscoBallFrames() []discoBallFrame {
	frames := make([]discoBallFrame, discoBallFrameCount)
	for i := 0; i < discoBallFrameCount; i++ {
		data, err := discoBallFramePNG(i)
		if err != nil {
			panic(fmt.Sprintf("tutor: embedded disco ball asset frame %d: %v", i, err))
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			panic(fmt.Sprintf("tutor: decode embedded disco ball asset frame %d: %v", i, err))
		}
		frames[i] = downsampleImage(img, discoBallSampleRows, discoBallSampleCols)
	}
	return frames
}

// discoBallFramePNG returns the raw embedded PNG bytes for frame index,
// unmodified — used by the Kitty rendering path (kittyimage.go), which
// transmits the original image data directly rather than downsampling
// it (Kitty decodes PNG natively and does its own display-time
// scaling), unlike renderDiscoBallFrame's ANSI path.
func discoBallFramePNG(index int) ([]byte, error) {
	return discoBallFS.ReadFile(fmt.Sprintf("discoballassets/frame_%d.png", index))
}

// downsampleImage box-averages img down to rows x cols cells. Averaging
// uses image.Color.RGBA()'s premultiplied-alpha values, which is what
// makes a cell straddling the ball's circular edge (part opaque, part
// transparent) average out to the correct un-premultiplied color instead
// of just darkening toward black. Per pixel, the premultiplied color is
// always <= its own alpha (both on the 0-65535 scale), so the averaged
// premultiplied color is always <= the averaged alpha too — the
// un-premultiply division below never overflows past full scale.
func downsampleImage(img image.Image, rows, cols int) discoBallFrame {
	b := img.Bounds()
	frame := make(discoBallFrame, rows)
	for row := 0; row < rows; row++ {
		frame[row] = make([]*color.RGBA, cols)
		y0 := b.Min.Y + row*b.Dy()/rows
		y1 := b.Min.Y + (row+1)*b.Dy()/rows
		for col := 0; col < cols; col++ {
			x0 := b.Min.X + col*b.Dx()/cols
			x1 := b.Min.X + (col+1)*b.Dx()/cols

			var rSum, gSum, bSum, aSum, n uint64
			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					r, g, bl, a := img.At(x, y).RGBA()
					rSum += uint64(r)
					gSum += uint64(g)
					bSum += uint64(bl)
					aSum += uint64(a)
					n++
				}
			}
			if n == 0 {
				continue
			}
			avgA := aSum / n
			if avgA < discoBallAlphaThreshold {
				continue // leave nil — transparent background (older sprite's convention)
			}
			r8 := uint8((rSum / n * 65535 / avgA) >> 8)
			g8 := uint8((gSum / n * 65535 / avgA) >> 8)
			b8 := uint8((bSum / n * 65535 / avgA) >> 8)
			if int(r8)+int(g8)+int(b8) < discoBallBrightnessThreshold {
				continue // leave nil — opaque near-black background (this
				// sprite's convention: solid black fill, alpha=255
				// everywhere, so alpha alone can't tell ball from
				// background; a near-black region is background even
				// though it's fully opaque)
			}
			frame[row][col] = &color.RGBA{
				R: r8,
				G: g8,
				B: b8,
				A: 255,
			}
		}
	}
	return frame
}

// halfBlockUpper/Lower are Block Elements (U+2580/U+2584) — a
// completely different Unicode range from the Braille spinner glyphs
// that broke the previous indicator (those were in the specialized
// Braille Patterns block, rarely bundled in programming fonts since
// they're not code-related). Block-drawing characters exist
// specifically for terminal graphics like this and are near-universal
// in monospace fonts — the same technique tools like chafa/catimg use
// to preview images in a terminal.
const (
	halfBlockUpper = "▀"
	halfBlockLower = "▄"
)

// renderDiscoBallFrame renders one frame to terminal lines, packing two
// sample rows into each physical line: the top sample's color becomes
// the half-block glyph's foreground (painting the cell's top half), the
// bottom sample's color becomes its background (painting the bottom
// half) — one truecolor-capable character cell showing two independent
// "pixels" instead of one. This is what lets the ball show real facet
// detail at a small physical size instead of the blurry few-cell blob
// an earlier one-sample-per-line version produced. A cell where only
// one of the pair is opaque falls back to a single half-block filled
// with just that color, leaving the other half as the pane's own
// background; a fully-background pair is a plain space.
func renderDiscoBallFrame(frame discoBallFrame) []string {
	sampleRows := len(frame)
	lines := make([]string, 0, (sampleRows+1)/2)
	for r := 0; r < sampleRows; r += 2 {
		top := frame[r]
		bottom := make([]*color.RGBA, len(top))
		if r+1 < sampleRows {
			bottom = frame[r+1]
		}

		var b strings.Builder
		for c := range top {
			t, bo := top[c], bottom[c]
			switch {
			case t == nil && bo == nil:
				b.WriteString(" ")
			case t != nil && bo != nil:
				fmt.Fprintf(&b, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%s\033[0m", t.R, t.G, t.B, bo.R, bo.G, bo.B, halfBlockUpper)
			case t != nil:
				fmt.Fprintf(&b, "\033[38;2;%d;%d;%dm%s\033[0m", t.R, t.G, t.B, halfBlockUpper)
			default:
				fmt.Fprintf(&b, "\033[38;2;%d;%d;%dm%s\033[0m", bo.R, bo.G, bo.B, halfBlockLower)
			}
		}
		lines = append(lines, b.String())
	}
	return lines
}

// ansiBallRenderer implements discoBallRenderer (thinkingdisplay.go)
// using the half-block ANSI mosaic above — always correct, works in any
// terminal, the guaranteed fallback for anyone not on a terminal
// kittyAvailable() recognizes (internal/tutor/kittyimage.go).
type ansiBallRenderer struct{}

func (ansiBallRenderer) init(io.Writer) {}

func (ansiBallRenderer) showFrame(w io.Writer, frame int) {
	for _, line := range renderDiscoBallFrame(discoBallFrames[frame%len(discoBallFrames)]) {
		fmt.Fprintf(w, "\033[2K\r  %s\n", line)
	}
}

func (ansiBallRenderer) rows() int { return discoBallRenderedRows }

func (ansiBallRenderer) close(io.Writer) {}
