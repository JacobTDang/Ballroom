package tutor

import (
	"bytes"
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestLoadDiscoBallFrames_ProducesFramesOfTheRightSize(t *testing.T) {
	frames := loadDiscoBallFrames()
	if len(frames) != discoBallFrameCount {
		t.Fatalf("len(frames) = %d, want %d", len(frames), discoBallFrameCount)
	}
	for i, frame := range frames {
		if len(frame) != discoBallSampleRows {
			t.Fatalf("frame %d has %d rows, want %d", i, len(frame), discoBallSampleRows)
		}
		for r, row := range frame {
			if len(row) != discoBallSampleCols {
				t.Fatalf("frame %d row %d has %d cols, want %d", i, r, len(row), discoBallSampleCols)
			}
		}
	}
}

func TestDiscoBallFramePNG_ReturnsValidPNGBytesForEachFrame(t *testing.T) {
	for i := 0; i < discoBallFrameCount; i++ {
		data, err := discoBallFramePNG(i)
		if err != nil {
			t.Fatalf("discoBallFramePNG(%d): %v", i, err)
		}
		if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
			t.Errorf("discoBallFramePNG(%d) is not a decodable image: %v", i, err)
		}
	}
}

func TestDiscoBallFramePNG_OutOfRangeReturnsError(t *testing.T) {
	if _, err := discoBallFramePNG(discoBallFrameCount); err == nil {
		t.Error("expected an error for an out-of-range frame index, got nil")
	}
}

func TestLoadDiscoBallFrames_HasSomeOpaqueAndSomeBackgroundCells(t *testing.T) {
	// Sanity check against the real sprite: a disco ball inscribed in a
	// square should leave the corners empty and the center painted — if
	// every cell came back nil (or none did), the alpha/brightness
	// threshold or the embed path is almost certainly wrong. Uses a
	// bright mid-sequence frame deliberately, not frame 0 — this sprite
	// fades in from black, so frame 0's ball is legitimately almost
	// entirely background-colored (that's the source content, not a
	// detection bug); a lit frame is what actually exercises "can this
	// tell ball from background" at all.
	frame := loadDiscoBallFrames()[8]
	var opaque, background int
	for _, row := range frame {
		for _, cell := range row {
			if cell == nil {
				background++
			} else {
				opaque++
			}
		}
	}
	if opaque == 0 {
		t.Error("expected at least one opaque (ball) cell, got none")
	}
	if background == 0 {
		t.Error("expected at least one background cell, got none")
	}
}

// solidAlphaImage builds an image.Image where the left half is fully
// opaque solid color c and the right half is fully transparent — enough
// to exercise downsampleImage's premultiplied-average and
// background-threshold logic without depending on the real sprite's
// exact pixel values.
func solidAlphaImage(w, h int, c color.RGBA) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < w/2 {
				img.SetNRGBA(x, y, color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255})
			} else {
				img.SetNRGBA(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}
	return img
}

func TestDownsampleImage_OpaqueRegionKeepsItsColor(t *testing.T) {
	src := solidAlphaImage(40, 20, color.RGBA{R: 200, G: 100, B: 50, A: 255})
	frame := downsampleImage(src, 2, 4)

	cell := frame[0][0] // top-left, entirely within the opaque left half
	if cell == nil {
		t.Fatal("expected an opaque cell, got nil (background)")
	}
	if cell.R != 200 || cell.G != 100 || cell.B != 50 {
		t.Errorf("cell = %+v, want R=200 G=100 B=50", cell)
	}
}

func TestDownsampleImage_TransparentRegionIsNil(t *testing.T) {
	src := solidAlphaImage(40, 20, color.RGBA{R: 200, G: 100, B: 50, A: 255})
	frame := downsampleImage(src, 2, 4)

	cell := frame[0][3] // top-right, entirely within the transparent right half
	if cell != nil {
		t.Errorf("expected a nil (background) cell, got %+v", cell)
	}
}

// solidOpaqueBlackImage builds an image.Image where the left half is a
// fully opaque bright color and the right half is fully opaque black —
// matching this sprite's own convention (solid black background,
// alpha=255 everywhere) rather than the older sprite's genuine
// transparency, which solidAlphaImage models.
func solidOpaqueBlackImage(w, h int, c color.RGBA) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < w/2 {
				img.SetNRGBA(x, y, color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255})
			} else {
				img.SetNRGBA(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	return img
}

func TestDownsampleImage_OpaqueBlackRegionIsNil(t *testing.T) {
	src := solidOpaqueBlackImage(40, 20, color.RGBA{R: 200, G: 100, B: 50, A: 255})
	frame := downsampleImage(src, 2, 4)

	bright := frame[0][0] // top-left, entirely within the bright left half
	if bright == nil {
		t.Fatal("expected the bright region to stay opaque, got nil")
	}
	if bright.R != 200 || bright.G != 100 || bright.B != 50 {
		t.Errorf("bright cell = %+v, want R=200 G=100 B=50", bright)
	}

	black := frame[0][3] // top-right, entirely within the opaque-black right half
	if black != nil {
		t.Errorf("expected the opaque-black region to be treated as background (nil), got %+v", black)
	}
}

func TestRenderDiscoBallFrame_PacksTwoSampleRowsPerLine(t *testing.T) {
	frame := discoBallFrame{
		{&color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{&color.RGBA{R: 0, G: 255, B: 0, A: 255}},
		{&color.RGBA{R: 0, G: 0, B: 255, A: 255}}, // unpaired third row
	}
	lines := renderDiscoBallFrame(frame)
	if len(lines) != 2 {
		t.Fatalf("len(lines) = %d, want 2 (ceil(3 sample rows / 2))", len(lines))
	}
}

func TestRenderDiscoBallFrame_BothOpaqueSetsForegroundAndBackground(t *testing.T) {
	frame := discoBallFrame{
		{&color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{&color.RGBA{R: 0, G: 255, B: 0, A: 255}},
	}
	lines := renderDiscoBallFrame(frame)
	if !strings.Contains(lines[0], "\033[38;2;255;0;0m") {
		t.Errorf("expected the top sample's color as foreground, got %q", lines[0])
	}
	if !strings.Contains(lines[0], "\033[48;2;0;255;0m") {
		t.Errorf("expected the bottom sample's color as background, got %q", lines[0])
	}
	if !strings.Contains(lines[0], halfBlockUpper) {
		t.Errorf("expected the upper-half-block glyph, got %q", lines[0])
	}
}

func TestRenderDiscoBallFrame_TopOnlyOmitsBackgroundEscape(t *testing.T) {
	frame := discoBallFrame{
		{&color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{nil},
	}
	lines := renderDiscoBallFrame(frame)
	if !strings.Contains(lines[0], "\033[38;2;255;0;0m") {
		t.Errorf("expected the top sample's color as foreground, got %q", lines[0])
	}
	if strings.Contains(lines[0], "\033[48;2;") {
		t.Errorf("expected no background escape when the bottom sample is background, got %q", lines[0])
	}
}

func TestRenderDiscoBallFrame_BottomOnlyUsesLowerHalfBlock(t *testing.T) {
	frame := discoBallFrame{
		{nil},
		{&color.RGBA{R: 0, G: 255, B: 0, A: 255}},
	}
	lines := renderDiscoBallFrame(frame)
	if !strings.Contains(lines[0], halfBlockLower) {
		t.Errorf("expected the lower-half-block glyph, got %q", lines[0])
	}
	if !strings.Contains(lines[0], "\033[38;2;0;255;0m") {
		t.Errorf("expected the bottom sample's color as foreground, got %q", lines[0])
	}
}

func TestRenderDiscoBallFrame_BothBackgroundIsPlainSpace(t *testing.T) {
	frame := discoBallFrame{
		{nil},
		{nil},
	}
	lines := renderDiscoBallFrame(frame)
	if lines[0] != " " {
		t.Errorf("lines[0] = %q, want a single plain space", lines[0])
	}
}

func TestAnsiBallRenderer_ShowFrameWritesTruecolorBallLines(t *testing.T) {
	var buf bytes.Buffer
	r := ansiBallRenderer{}
	r.init(&buf) // no-op; called for interface completeness
	r.showFrame(&buf, 0)

	got := buf.String()
	if !strings.Contains(got, "\033[38;2;") && !strings.Contains(got, "\033[48;2;") {
		t.Error("expected truecolor escapes from the ANSI ball renderer, got none")
	}
	if strings.Count(got, "\n") != discoBallRenderedRows {
		t.Errorf("showFrame wrote %d lines, want exactly discoBallRenderedRows (%d)", strings.Count(got, "\n"), discoBallRenderedRows)
	}
}

func TestAnsiBallRenderer_RowsMatchesActualLinesWritten(t *testing.T) {
	var buf bytes.Buffer
	r := ansiBallRenderer{}
	r.showFrame(&buf, 0)

	if got := strings.Count(buf.String(), "\n"); got != r.rows() {
		t.Errorf("showFrame wrote %d lines, rows() reports %d — must match for redrawLocked's cursor-up math to stay correct", got, r.rows())
	}
}

func TestAnsiBallRenderer_CloseIsANoOp(t *testing.T) {
	var buf bytes.Buffer
	ansiBallRenderer{}.close(&buf)
	if buf.Len() != 0 {
		t.Errorf("close() wrote %q, want nothing — the ANSI renderer owns no resources to release", buf.String())
	}
}
