package tui

import (
	"strings"

	"github.com/JacobTDang/Ballroom/internal/palette"
	"github.com/charmbracelet/lipgloss"
)

// discoBallHeight/Width are the ball's floor size — what every
// terminal used to get unconditionally, and still the guaranteed
// minimum the panel arithmetic (minPanelWidth) is sized around. Large
// terminals now scale the ball up from here (see ballDimensions,
// dashboard.go) after a real complaint that the fixed ball left the
// home screen looking like everything was crammed into one corner.
// discoBallMaxHeight caps the growth: past ~36 rows the ball stops
// reading as an accent and starts dominating the panel.
const (
	discoBallHeight    = 24
	discoBallWidth     = 48
	discoBallMaxHeight = 36
)

// discoShades goes dim -> bright; the full set of characters a rendered
// ball ever uses (see discoTileSequence for the order they're tiled in).
var discoShades = []rune{'.', ':', '*', '%', '#'}

// discoTileSequence is the repeating bright->dim mirror-tile unit, tiled
// diagonally across the ball (see buildDiscoBall) to produce banded
// facets rather than a single smooth light-source gradient — modeled on
// a hand-drawn ASCII disco-ball reference. Each character repeats twice
// so the bands read clearly at terminal resolution instead of dissolving
// into per-column flicker.
var discoTileSequence = []rune{'%', '%', '#', '#', '*', '*', ':', ':', '.', '.'}

var (
	discoBallDimStyle    = lipgloss.NewStyle().Foreground(palette.Lip(palette.DimGray))
	discoBallBrightStyle = lipgloss.NewStyle().Foreground(palette.Lip(palette.PaleGray))

	// The sparse "glint" palette — only a small fraction of tiles ever
	// use these; everything else stays grayscale.
	discoSparkleStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(palette.Lip(palette.Cyan)), // cyan
		lipgloss.NewStyle().Foreground(palette.Lip(palette.Pink)), // magenta
		lipgloss.NewStyle().Foreground(palette.Lip(palette.Gold)), // gold
		lipgloss.NewStyle().Foreground(palette.Lip(palette.Teal)), // teal
	}
)

// discoBallCell is one character position in a generated ball.
type discoBallCell struct {
	ch      rune
	sparkle bool
}

// buildDiscoBall computes a static circular mirror-ball texture: an
// ellipse mask (width intentionally ~2x height, since typical terminal
// character cells are roughly twice as tall as wide — without that
// correction the "circle" reads as a tall oval), filled by tiling
// discoTileSequence diagonally (indexed by row+col) so the facets read as
// banded mirror tiles catching light at an angle, rather than a single
// smooth gradient. A small set of tiles in the equator band (see
// sparkleAt) are marked sparkle=true — which ones is deterministic (not
// random per call), only their *color* animates, via renderDiscoBall's
// phase.
func buildDiscoBall(height, width int) [][]discoBallCell {
	grid := make([][]discoBallCell, height)
	cx, cy := float64(width)/2, float64(height)/2
	rx, ry := float64(width)/2, float64(height)/2

	for row := 0; row < height; row++ {
		grid[row] = make([]discoBallCell, width)
		for col := 0; col < width; col++ {
			fx := (float64(col) + 0.5 - cx) / rx
			fy := (float64(row) + 0.5 - cy) / ry
			if fx*fx+fy*fy > 1.0 {
				grid[row][col] = discoBallCell{ch: ' '}
				continue
			}

			idx := (row + col) % len(discoTileSequence)
			grid[row][col] = discoBallCell{
				ch:      discoTileSequence[idx],
				sparkle: sparkleAt(row, col, height),
			}
		}
	}
	return grid
}

// discoClusterRowSpacing/ColSpacing/Size lay shimmer clusters out on a
// coarse grid within the equator band — small discoClusterSize x
// discoClusterSize clumps of glint, spaced apart, rather than scattered
// single cells, so it reads as a few mirror tiles catching light instead
// of random color noise.
const (
	discoClusterRowSpacing = 4
	discoClusterColSpacing = 6
	discoClusterSize       = 2
)

// sparkleAt marks cells belonging to a small cluster on that coarse grid,
// restricted to the equator band (the middle third of the ball's height)
// — where a spinning mirror ball actually catches and scatters the most
// light. Alternating cluster rows are horizontally offset so the clusters
// don't line up into a rigid, obviously-mechanical grid.
func sparkleAt(row, col, height int) bool {
	top, bottom := height/3, height-height/3
	if row < top || row >= bottom {
		return false
	}
	bandRow := row - top
	rowBand := bandRow / discoClusterRowSpacing
	colOffset := (rowBand % 2) * (discoClusterColSpacing / 2)
	localCol := col + colOffset
	return bandRow%discoClusterRowSpacing < discoClusterSize &&
		localCol%discoClusterColSpacing < discoClusterSize
}

// sparkleColorIndex picks a discoSparkleStyles index for a glinting
// cell, offset by its own position so cells cycle at staggered moments
// (driven by phase) instead of all changing color in lockstep.
func sparkleColorIndex(phase, row, col int) int {
	return (phase + row*3 + col) % len(discoSparkleStyles)
}

// renderDiscoBall renders a pre-built ball grid to a string. Only
// sparkle cells get colored, cycling through discoSparkleStyles offset
// by each cell's own position (so they glint at staggered moments
// rather than all changing color in lockstep) — every other character
// stays a static grayscale, so it reads as light catching a few mirror
// tiles, not the whole ball flashing.
func renderDiscoBall(grid [][]discoBallCell, phase int) string {
	var b strings.Builder
	for row, cells := range grid {
		for col, cell := range cells {
			switch {
			case cell.ch == ' ':
				b.WriteRune(' ')
			case cell.sparkle:
				idx := sparkleColorIndex(phase, row, col)
				b.WriteString(discoSparkleStyles[idx].Render(string(cell.ch)))
			case cell.ch == '#' || cell.ch == '%':
				b.WriteString(discoBallBrightStyle.Render(string(cell.ch)))
			default:
				b.WriteString(discoBallDimStyle.Render(string(cell.ch)))
			}
		}
		if row < len(grid)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()
}
