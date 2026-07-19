package tutor

import (
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// Display-only markdown styling for the chat pane. The tutor's replies
// arrive as markdown (bold emphasis, inline code, fenced blocks) and
// used to be shown raw -- literal ** and ``` on screen, a real
// complaint ("**<text>** is trying bold but doesnt actually"). This is
// deliberately a tiny hand-rolled pass over the handful of constructs
// the tutor actually emits, not a full markdown engine: a rendering
// library would re-wrap text (fighting refreshViewport's own wrapping)
// and add a dependency for four regexes' worth of styling.
//
// Applied ONLY at the displayLines append (model.go's turnCompleteMsg)
// -- m.history keeps the raw text, because that's what goes back to
// the model as conversation context and escape codes there would
// pollute every later turn's prompt.

var (
	// Bold: **text** (single line, no stray-asterisk false positives --
	// "2 * 3" has no ** pair and passes through untouched).
	mdBoldPattern = regexp.MustCompile(`\*\*([^*\n]+)\*\*`)
	// Inline code: `text`.
	mdInlineCodePattern = regexp.MustCompile("`([^`\n]+)`")
	// Headers: leading #s at line start.
	mdHeaderPattern = regexp.MustCompile(`^(#{1,4})\s+(.*)$`)
	// List items: the marker needs trailing whitespace, so *emphasis*
	// and **bold** at line start never read as bullets.
	mdBulletPattern  = regexp.MustCompile(`^(\s*)[-*]\s+(.*)$`)
	mdOrderedPattern = regexp.MustCompile(`^(\s*)(\d+)\.\s+(.*)$`)
	// Blockquotes: > with an optional single space, bare > included.
	mdQuotePattern = regexp.MustCompile(`^>\s?(.*)$`)
	// Horizontal rules: three or more -/*/_ alone on a line. Two
	// dashes are prose (an em-dash stand-in), never a rule.
	mdHrPattern = regexp.MustCompile(`^\s*(-{3,}|\*{3,}|_{3,})\s*$`)
	// Links: [text](url), with the optional leading ! captured so
	// image syntax can be recognized and left raw.
	mdLinkPattern = regexp.MustCompile(`(!?)\[([^\]\n]+)\]\(([^)\n]+)\)`)
)

// Styling escapes. Bold closes with 22 (bold off) and underline with
// 24 rather than a full reset so they can't clobber styling an
// enclosing construct set up; color spans close with 39 (default
// foreground) for the same reason.
const (
	mdBoldOn       = "\x1b[1m"
	mdBoldOff      = "\x1b[22m"
	mdUnderlineOn  = "\x1b[4m"
	mdUnderlineOff = "\x1b[24m"
	mdColorReset   = "\x1b[39m"
)

var (
	// Inline code and fenced code share the pane's teal accent (the
	// same accent as the input prompt glyph), so "code" reads as one
	// consistent signal everywhere it appears.
	mdCodeColor = ansiFg(paneTeal)
	// Fence labels ("python") render dim -- metadata, not content.
	mdDimColor = ansiFg(paneDimText)
)

// styleMarkdown renders the tutor's markdown constructs as terminal
// styling: **bold**, `inline code`, #-headers, and ``` fenced blocks.
// Inside a fence nothing is transformed -- code legitimately contains
// * and ` characters.
//
// width is the render-time pane width fenced blocks may size
// themselves to -- threaded through since the display pipeline now
// re-renders every block per frame at the viewport's current width
// (see displayBlock, model.go). Prose ignores it entirely;
// refreshViewport word-wraps prose after the fact, same as always.
func styleMarkdown(content string, width int) string {
	cardWidth := cardWidthFor(width)
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	inFence := false
	var fenceLabel string
	var fenceLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if !inFence {
				inFence = true
				fenceLabel = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
				fenceLines = fenceLines[:0]
			} else {
				inFence = false
				out = append(out, renderFence(fenceLabel, fenceLines, cardWidth, false)...)
			}
			continue
		}
		if inFence {
			fenceLines = append(fenceLines, line)
			continue
		}
		if m := mdHeaderPattern.FindStringSubmatch(line); m != nil {
			out = append(out, mdBoldOn+m[2]+mdBoldOff)
			continue
		}
		if mdHrPattern.MatchString(line) {
			out = append(out, renderRule(width))
			continue
		}
		if m := mdQuotePattern.FindStringSubmatch(line); m != nil {
			out = append(out, renderQuote(m[1], width))
			continue
		}
		if m := mdBulletPattern.FindStringSubmatch(line); m != nil {
			marker := mdCodeColor + "-" + mdColorReset
			out = append(out, renderListItem(m[1], marker, 1, m[2], width)...)
			continue
		}
		if m := mdOrderedPattern.FindStringSubmatch(line); m != nil {
			marker := mdDimColor + m[2] + "." + mdColorReset
			out = append(out, renderListItem(m[1], marker, len(m[2])+1, m[3], width)...)
			continue
		}
		out = append(out, styleInline(line))
	}
	if inFence {
		// Unterminated fence (a reply still streaming, or truncated):
		// still render what arrived rather than dropping it -- as a
		// bottomless card that visibly grows until the closing fence.
		out = append(out, renderFence(fenceLabel, fenceLines, cardWidth, true)...)
	}
	return strings.Join(out, "\n")
}

// styleInline styles one prose line's inline constructs: code spans
// first, splitting the line around them so nothing else ever
// transforms inside a span (a literal `[text](url)` in backticks must
// stay literal -- running the passes sequentially over the whole line
// would restyle the span's contents), then bold and links on the
// non-code segments only.
func styleInline(line string) string {
	locs := mdInlineCodePattern.FindAllStringSubmatchIndex(line, -1)
	if len(locs) == 0 {
		return styleProse(line)
	}
	var b strings.Builder
	last := 0
	for _, loc := range locs {
		b.WriteString(styleProse(line[last:loc[0]]))
		b.WriteString(mdCodeColor + line[loc[2]:loc[3]] + mdColorReset)
		last = loc[1]
	}
	b.WriteString(styleProse(line[last:]))
	return b.String()
}

// styleProse is the non-code half of styleInline: bold, then links --
// link text underlined, the URL dim in parens (terminals make plain
// URLs clickable on their own; the dim keeps it from competing with
// the prose). Image syntax (![...]) is left raw: the pane can't
// render an image, and pretending it's a link would be a lie.
func styleProse(s string) string {
	s = mdBoldPattern.ReplaceAllString(s, mdBoldOn+"$1"+mdBoldOff)
	return mdLinkPattern.ReplaceAllStringFunc(s, func(match string) string {
		m := mdLinkPattern.FindStringSubmatch(match)
		if m[1] == "!" {
			return match
		}
		return mdUnderlineOn + m[2] + mdUnderlineOff + mdDimColor + " (" + m[3] + ")" + mdColorReset
	})
}

// minListWrapWidth is the narrowest usable wrap column for a list
// item's text; below it the hang indent costs more than it's worth,
// so the item renders on one line and the outer wrap deals with it
// (same spirit as minCardWidth degrading fences to the flat style).
const minListWrapWidth = 8

// renderListItem renders one list item with a hang indent: the styled
// marker on the first line, continuation lines aligned under the item
// text. The wrapping happens HERE, not in refreshViewport's outer
// pass -- that pass is escape-aware but indent-blind, and would fold
// a long item flush-left, destroying the hang. Emitting lines that
// already fit the width keeps the outer pass a no-op for them (the
// same invariant the editor cards rely on). markerWidth is the
// marker's visual width (the styled string carries escapes, so the
// caller states it).
func renderListItem(indent, styledMarker string, markerWidth int, text string, width int) []string {
	styled := styleInline(text)
	first := indent + styledMarker + " "
	avail := width - len(indent) - markerWidth - 1
	if width <= 0 || avail < minListWrapWidth {
		return []string{first + styled}
	}
	wrapped := strings.Split(lipgloss.NewStyle().Width(avail).Render(styled), "\n")
	hang := indent + strings.Repeat(" ", markerWidth+1)
	out := make([]string, 0, len(wrapped))
	out = append(out, first+wrapped[0])
	for _, cont := range wrapped[1:] {
		out = append(out, hang+cont)
	}
	return out
}

// renderQuote renders a blockquote line: a bar in the structural rule
// color plus the quoted text in the metadata dim, self-wrapped so the
// bar lands on every wrapped line — the outer pass would fold a long
// quote flush-left and orphan the continuation from its bar (seen
// live in the preview harness). Inline spans still style; their
// closing default-foreground resets are re-armed to dim so the rest
// of the quote can't "leak" back to full brightness mid-line (the
// same re-arm trick the editor cards use on chroma's resets).
func renderQuote(text string, width int) string {
	styled := strings.ReplaceAll(styleInline(text), mdColorReset, mdDimColor)
	bar := ansiFg(paneRule) + "│" + mdColorReset + " "
	avail := width - 2
	if width <= 0 || avail < minListWrapWidth {
		return bar + mdDimColor + styled + mdColorReset
	}
	wrapped := strings.Split(lipgloss.NewStyle().Width(avail).Render(styled), "\n")
	for i, line := range wrapped {
		wrapped[i] = bar + mdDimColor + line + mdColorReset
	}
	return strings.Join(wrapped, "\n")
}

// mdRuleWidth caps a horizontal rule: full-width rules read as pane
// chrome (the old header rule), and 40 cells is plenty to say
// "section break" in a chat column.
const mdRuleWidth = 40

// renderRule draws a markdown hr (---/***/___) as a double rule, the
// same weight as the editor cards' frames (card.go) and the host
// panel's own border (lipgloss.DoubleBorder(), internal/tui/dashboard.go)
// -- so every horizontal line in the pane reads as one system.
func renderRule(width int) string {
	n := mdRuleWidth
	if width > 0 && width < n {
		n = width
	}
	return mdDimColor + strings.Repeat("═", n) + mdColorReset
}

// renderFence renders one fenced block: an editor card (card.go) when
// the pane gives it enough width, the width-independent rule style
// otherwise -- dim label rule, highlighted code, closing rule. The rule
// itself is double, matching renderRule and the cards it stands in for
// when the pane is too narrow to frame one.
func renderFence(label string, lines []string, cardWidth int, unterminated bool) []string {
	if cardWidth > 0 {
		return renderCodeCard(label, lines, cardWidth, unterminated)
	}
	shown := label
	if shown == "" {
		shown = "code"
	}
	out := make([]string, 0, len(lines)+2)
	out = append(out, mdDimColor+"══ "+shown+" ══"+mdColorReset)
	out = append(out, highlightCode(label, lines)...)
	if !unterminated {
		out = append(out, mdDimColor+"══"+mdColorReset)
	}
	return out
}

// chromaStyle renders well on this pane's dark background; the
// formatter emits the same 24-bit truecolor escapes the rest of the
// chat styling uses (tmux.conf already enables RGB passthrough).
var (
	chromaStyle     = styles.Get("monokai")
	chromaFormatter = formatters.Get("terminal16m")
)

// highlightCode renders a fenced block's lines: real per-token syntax
// highlighting when chroma knows the language label, the flat accent
// color otherwise (including unlabeled fences). Highlighting the whole
// block at once -- not line by line -- keeps multi-line constructs
// (strings, comments) tokenized correctly.
func highlightCode(label string, lines []string) []string {
	lexer := lexers.Get(label)
	if label == "" || lexer == nil {
		return flatCode(lines)
	}
	it, err := chroma.Coalesce(lexer).Tokenise(nil, strings.Join(lines, "\n"))
	if err != nil {
		return flatCode(lines)
	}
	var buf strings.Builder
	if err := chromaFormatter.Format(&buf, chromaStyle, it); err != nil {
		return flatCode(lines)
	}
	return strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
}

// flatCode is the pre-chroma rendering -- every line in the shared
// accent color -- kept as the fallback for unlabeled fences, unknown
// languages, and highlighter errors.
func flatCode(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = mdCodeColor + line + mdColorReset
	}
	return out
}
