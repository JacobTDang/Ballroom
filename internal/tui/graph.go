package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// boxGap is the horizontal spacing (in columns) between adjacent boxes
// in a row.
const boxGap = 3

// Solid-filled boxes (background color + bold light text), matching
// NeetCode's roadmap look — colored "buttons", not just outlined frames.
var (
	rootBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#4A4A4A")).
			Foreground(lipgloss.Color("#F2EBDD")).
			Bold(true).
			Padding(0, 2)

	categoryBoxStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#3C7DC4")).
				Foreground(lipgloss.Color("#F2EBDD")).
				Bold(true).
				Padding(0, 1)

	categorySolvedBoxStyle = categoryBoxStyle.
				Background(lipgloss.Color("#2FA6A6"))

	highlightedBoxStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#9B5FB0")).
				Foreground(lipgloss.Color("#F2EBDD")).
				Bold(true).
				Padding(0, 1)

	exerciseNotAttemptedBoxStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#6B6B6B")).
					Foreground(lipgloss.Color("#F2EBDD")).
					Bold(true).
					Padding(0, 1)

	exerciseFailBoxStyle = exerciseNotAttemptedBoxStyle.
				Background(lipgloss.Color("#F03C3C"))

	exercisePassBoxStyle = exerciseNotAttemptedBoxStyle.
				Background(lipgloss.Color("#2FA6A6"))
)

// plainProgressBar is pixelProgressBar without per-block coloring, for
// embedding inside a box whose border/text color already signals state —
// nesting separately-styled ANSI segments inside a lipgloss-rendered box
// causes the inner reset codes to clobber the box's own color for
// anything after them, so box content stays plain and lets the box color
// carry the signal instead.
func plainProgressBar(solved, total int) string {
	filled := 0
	if total > 0 {
		filled = solved * pixelBarWidth / total
		if filled == 0 && solved > 0 {
			filled = 1
		}
		if filled > pixelBarWidth {
			filled = pixelBarWidth
		}
	}
	var b strings.Builder
	for i := 0; i < pixelBarWidth; i++ {
		if i < filled {
			b.WriteRune('▓')
		} else {
			b.WriteRune('░')
		}
	}
	return b.String()
}

func renderRootBox() string {
	return rootBoxStyle.Render("PRACTICE")
}

func renderCategoryBox(category string, solved, total int, highlighted bool) string {
	content := category + "\n" + fmt.Sprintf("%s (%d/%d)", plainProgressBar(solved, total), solved, total)
	style := categoryBoxStyle
	if solved > 0 && solved == total {
		style = categorySolvedBoxStyle
	}
	if highlighted {
		style = highlightedBoxStyle
	}
	return style.Render(content)
}

func renderExerciseBox(s catalog.ExerciseStatus, highlighted bool) string {
	icon := "░░"
	style := exerciseNotAttemptedBoxStyle
	switch s.LastResult {
	case tracker.ResultPass:
		icon = "▓▓✦"
		style = exercisePassBoxStyle
	case tracker.ResultFail:
		icon = "▓░"
		style = exerciseFailBoxStyle
	}
	if highlighted {
		style = highlightedBoxStyle
	}
	content := s.Exercise.Language + "\n" + icon
	return style.Render(content)
}

// spanAnchor returns the column a parent should center itself over given
// a row of children's center columns: the true middle child's own center
// for an odd count (visually cleaner — the stem lands exactly on that
// child instead of merely near it), or the midpoint between the two
// middle children for an even count.
func spanAnchor(childCenters []int) int {
	n := len(childCenters)
	if n%2 == 1 {
		return childCenters[n/2]
	}
	return (childCenters[n/2-1] + childCenters[n/2]) / 2
}

// boxCenters returns each box's horizontal center column, given their
// rendered widths and the gap placed between adjacent boxes. Columns are
// relative to column 0 of the row.
func boxCenters(widths []int, gap int) []int {
	centers := make([]int, len(widths))
	pos := 0
	for i, w := range widths {
		centers[i] = pos + w/2
		pos += w + gap
	}
	return centers
}

// centerOffset returns how many columns to left-pad a box of the given
// width so its own center lands on spanCenter. Never negative.
func centerOffset(boxWidth, spanCenter int) int {
	offset := spanCenter - boxWidth/2
	if offset < 0 {
		return 0
	}
	return offset
}

// connectorLines draws the 3-line org-chart-style connector between a
// parent's center column and a row of children's center columns: a stem
// down from the parent, a horizontal spine across the children, and a
// stem down into each child. Where the parent's stem lands exactly on a
// child, that's a '┼' (both); otherwise a '┴' where it just meets the
// spine.
func connectorLines(width, parentCenter int, childCenters []int) []string {
	line1 := blankLine(width)
	line2 := blankLine(width)
	line3 := blankLine(width)

	if parentCenter >= 0 && parentCenter < width {
		line1[parentCenter] = '│'
	}

	minC, maxC := childCenters[0], childCenters[0]
	for _, c := range childCenters {
		if c < minC {
			minC = c
		}
		if c > maxC {
			maxC = c
		}
	}
	for i := minC; i <= maxC; i++ {
		if i >= 0 && i < width {
			line2[i] = '─'
		}
	}

	isChild := make(map[int]bool, len(childCenters))
	for _, c := range childCenters {
		isChild[c] = true
		if c >= 0 && c < width {
			line2[c] = '┬'
			line3[c] = '│'
		}
	}
	if parentCenter >= 0 && parentCenter < width {
		if isChild[parentCenter] {
			line2[parentCenter] = '┼'
		} else if parentCenter >= minC && parentCenter <= maxC {
			line2[parentCenter] = '┴'
		}
	}

	return []string{string(line1), string(line2), string(line3)}
}

func blankLine(width int) []rune {
	line := make([]rune, width)
	for i := range line {
		line[i] = ' '
	}
	return line
}

// joinBoxesHorizontal lays out a row of (possibly multi-line) boxes side
// by side with gap spaces between them.
func joinBoxesHorizontal(boxes []string, gap int) string {
	if len(boxes) == 0 {
		return ""
	}
	gapCol := strings.Repeat(" ", gap)

	// Split each box into its lines; assume uniform height across boxes
	// in the same row (true for our fixed-style single-line-label boxes).
	splitBoxes := make([][]string, len(boxes))
	height := 0
	for i, b := range boxes {
		lines := strings.Split(b, "\n")
		splitBoxes[i] = lines
		if len(lines) > height {
			height = len(lines)
		}
	}

	rows := make([]string, height)
	for row := 0; row < height; row++ {
		var b strings.Builder
		for i, lines := range splitBoxes {
			if i > 0 {
				b.WriteString(gapCol)
			}
			if row < len(lines) {
				b.WriteString(lines[row])
			}
		}
		rows[row] = b.String()
	}
	return strings.Join(rows, "\n")
}

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
