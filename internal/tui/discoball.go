package tui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// discoBallHeight/Width are fixed — the ball is a consistent, modest
// size on every screen rather than scaling with the terminal, so it
// stays proportionate next to the animated banner instead of dominating
// the panel.
const (
	discoBallHeight = 24
	discoBallWidth  = 48
)

// discoShades goes dim -> bright; used for the mirror-ball's directional
// shading (denser characters catch more "light").
var discoShades = []rune{'.', ':', '*', '%', '#'}

var (
	discoBallDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B6B6B"))
	discoBallBrightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9D3C4"))

	// The sparse "glint" palette — only a small fraction of tiles ever
	// use these; everything else stays grayscale.
	discoSparkleStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("#3ED6D6")), // cyan
		lipgloss.NewStyle().Foreground(lipgloss.Color("#E0468C")), // magenta
		lipgloss.NewStyle().Foreground(lipgloss.Color("#E8A93C")), // gold
		lipgloss.NewStyle().Foreground(lipgloss.Color("#2FA6A6")), // teal
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
// correction the "circle" reads as a tall oval), filled with a shading
// gradient simulating a single light source (denser characters closer to
// the light, sparser toward the shadowed edge) plus a sparse dotted grid
// evoking individual mirror tiles. A small, fixed fraction of tiles are
// marked sparkle=true — which ones is deterministic (not random per
// call), only their *color* animates, via renderDiscoBall's phase.
func buildDiscoBall(height, width int) [][]discoBallCell {
	grid := make([][]discoBallCell, height)
	cx, cy := float64(width)/2, float64(height)/2
	rx, ry := float64(width)/2, float64(height)/2
	lightX, lightY := cx-rx*0.5, cy-ry*0.6

	for row := 0; row < height; row++ {
		grid[row] = make([]discoBallCell, width)
		for col := 0; col < width; col++ {
			fx := (float64(col) + 0.5 - cx) / rx
			fy := (float64(row) + 0.5 - cy) / ry
			if fx*fx+fy*fy > 1.0 {
				grid[row][col] = discoBallCell{ch: ' '}
				continue
			}

			// Sparse dotted grid lines (mirror tile borders) — not a
			// solid overlay, so it still reads as shaded rather than
			// gridded-over.
			if row%2 == 0 && col%3 == 0 {
				grid[row][col] = discoBallCell{ch: '.', sparkle: sparkleAt(row, col)}
				continue
			}

			dx := float64(col) - lightX
			dy := (float64(row) - lightY) * 2 // correct for cell aspect ratio
			lightDist := math.Sqrt(dx*dx+dy*dy) / (rx * 1.4)
			idx := len(discoShades) - 1 - int(lightDist*float64(len(discoShades)))
			if idx < 0 {
				idx = 0
			}
			if idx >= len(discoShades) {
				idx = len(discoShades) - 1
			}

			grid[row][col] = discoBallCell{ch: discoShades[idx], sparkle: sparkleAt(row, col)}
		}
	}
	return grid
}

// sparkleAt deterministically marks roughly 12% of cells as glint
// candidates, scattered rather than clustered.
func sparkleAt(row, col int) bool {
	return (row*31+col*17)%100 < 12
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
