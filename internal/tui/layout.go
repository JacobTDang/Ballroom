package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// padLeft prepends n spaces to every line of a (possibly multi-line)
// block of text.
func padLeft(block string, n int) string {
	if n <= 0 {
		return block
	}
	pad := strings.Repeat(" ", n)
	lines := strings.Split(block, "\n")
	for i, l := range lines {
		lines[i] = pad + l
	}
	return strings.Join(lines, "\n")
}

// placeBlock centers content as a single rigid block within a
// termWidth x termHeight viewport — unlike lipgloss.Place, which centers
// each line of the content independently (measured against that line's
// own width), silently destroying any deliberate relative horizontal
// alignment between lines. Every line gets the same left-pad, computed
// once from the content's own widest line.
func placeBlock(termWidth, termHeight int, content string) string {
	lines := strings.Split(content, "\n")
	contentWidth := 0
	for _, l := range lines {
		if w := lipgloss.Width(l); w > contentWidth {
			contentWidth = w
		}
	}

	leftPad := (termWidth - contentWidth) / 2
	if leftPad > 0 {
		content = padLeft(content, leftPad)
		lines = strings.Split(content, "\n")
	}

	topPad := (termHeight - len(lines)) / 2
	if topPad > 0 {
		content = strings.Repeat("\n", topPad) + content
	}

	return content
}
