package tutor

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// The editor card: fenced code blocks render as a floating editor
// window -- rounded frame, header bar with three traffic-light dots
// and a centered language label, a dim line-number gutter, and
// syntax-highlighted code on the card's own background -- per the
// user's reference screenshot ("in the tutor terminal it could display
// a widget like this"). Cards are fixed-width constructs built per
// frame at the current pane width (see displayBlock, model.go); code
// wider than the card hard-truncates with a dim … rather than ever
// wrapping inside the gutter.

// minCardWidth is the narrowest frame that still fits borders, the
// gutter, and enough code cells to be worth framing -- below it,
// renderCodeCard falls back to the flat rendering instead of drawing
// something broken.
const minCardWidth = 16

// cardsEnabledWidth converts styleMarkdown's content width into the
// card's own width (inset from the prose flow), or 0 when the pane is
// too narrow for cards at all.
func cardWidthFor(contentWidth int) int {
	w := contentWidth - 2
	if w < minCardWidth {
		return 0
	}
	return w
}

// renderCodeCard draws one fenced block as an editor card exactly width
// cells wide. unterminated (a fence still streaming in) leaves the
// bottom border off -- the card visibly grows until the closing fence
// arrives. Below minCardWidth it degrades to flatCode.
func renderCodeCard(label string, lines []string, width int, unterminated bool) []string {
	if width < minCardWidth {
		return flatCode(lines)
	}
	interior := width - 2

	borderFg := ansiFg(paneRule)
	out := make([]string, 0, len(lines)+3)
	out = append(out, borderFg+"╭"+strings.Repeat("─", interior)+"╮\x1b[0m")
	out = append(out, borderFg+"│\x1b[0m"+cardHeaderRow(label, interior)+borderFg+"│\x1b[0m")

	gutterWidth := len(fmt.Sprintf("%d", max(len(lines), 1)))
	if gutterWidth < 3 {
		gutterWidth = 3
	}
	// interior = " " + gutter + "  " + code
	codeWidth := interior - 1 - gutterWidth - 2
	bg := ansiBg(cardBg)
	for i, code := range highlightCode(label, lines) {
		if lipgloss.Width(code) > codeWidth {
			// Hard truncate -- wrapping would break the gutter's
			// line-number alignment. The dim … says content was cut.
			code = ansi.Truncate(code, codeWidth-1, "") + mdDimColor + "…" + mdColorReset
		}
		pad := codeWidth - lipgloss.Width(code)
		row := bg + " " + ansiFg(cardGutterFg) + fmt.Sprintf("%*d", gutterWidth, i+1) + mdColorReset + "  " +
			// Chroma ends tokens with full \x1b[0m resets, which would
			// drop the card background mid-row -- re-arm it after each.
			strings.ReplaceAll(code, "\x1b[0m", "\x1b[0m"+bg) +
			bg + strings.Repeat(" ", max(pad, 0)) + "\x1b[0m"
		out = append(out, borderFg+"│\x1b[0m"+row+borderFg+"│\x1b[0m")
	}

	if !unterminated {
		out = append(out, borderFg+"╰"+strings.Repeat("─", interior)+"╯\x1b[0m")
	}
	return out
}

// cardHeaderRow is the card's title bar: three traffic-light dots on
// the left, the language label centered, all on the header background.
func cardHeaderRow(label string, interior int) string {
	if label == "" {
		label = "code"
	}
	const dotsWidth = 6 // " ● ● ●"
	maxLabel := interior - dotsWidth - 2
	if maxLabel < 1 {
		label = ""
	} else if lipgloss.Width(label) > maxLabel {
		label = ansi.Truncate(label, maxLabel, "…")
	}

	labelWidth := lipgloss.Width(label)
	start := (interior - labelWidth) / 2
	if start < dotsWidth+1 {
		start = dotsWidth + 1
	}
	rightPad := interior - start - labelWidth
	if rightPad < 0 {
		rightPad = 0
	}

	hbg := ansiBg(cardHeaderBg)
	return hbg + " " + ansiFg(trafficRed) + "●" + " " + ansiFg(trafficGold) + "●" + " " + ansiFg(paneTeal) + "●" +
		strings.Repeat(" ", start-dotsWidth) + mdDimColor + label + mdColorReset +
		strings.Repeat(" ", rightPad) + "\x1b[0m"
}
