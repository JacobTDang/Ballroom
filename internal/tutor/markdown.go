package tutor

import (
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
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
)

// Styling escapes. Bold closes with 22 (bold off) rather than a full
// reset so it can't clobber styling an enclosing construct set up;
// color spans close with 39 (default foreground) for the same reason.
const (
	mdBoldOn     = "\x1b[1m"
	mdBoldOff    = "\x1b[22m"
	mdColorReset = "\x1b[39m"
)

var (
	// Inline code and fenced code share the pane's teal accent (the
	// same accent as the input prompt glyph), so "code" reads as one
	// consistent signal everywhere it appears.
	mdCodeColor = ansiFg(paneTeal)
	// Fence labels ("python") render dim -- metadata, not content.
	mdDimColor = ansiFg(paneDimText)

	// userEchoPrefix marks the user's own submitted lines in the
	// transcript (display-only; history keeps the raw text): a dim
	// "you" so the speaker label never competes with the message, with
	// the › glyph in the palette's pink -- the user's own accent color,
	// distinct from the tutor's teal.
	userEchoPrefix = mdDimColor + "you " + mdColorReset + mdBoldOn + ansiFg(panePink) + "› " + mdColorReset + mdBoldOff
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
		line = mdBoldPattern.ReplaceAllString(line, mdBoldOn+"$1"+mdBoldOff)
		line = mdInlineCodePattern.ReplaceAllString(line, mdCodeColor+"$1"+mdColorReset)
		out = append(out, line)
	}
	if inFence {
		// Unterminated fence (a reply still streaming, or truncated):
		// still render what arrived rather than dropping it -- as a
		// bottomless card that visibly grows until the closing fence.
		out = append(out, renderFence(fenceLabel, fenceLines, cardWidth, true)...)
	}
	return strings.Join(out, "\n")
}

// renderFence renders one fenced block: an editor card (card.go) when
// the pane gives it enough width, the width-independent rule style
// otherwise -- dim label rule, highlighted code, closing rule.
func renderFence(label string, lines []string, cardWidth int, unterminated bool) []string {
	if cardWidth > 0 {
		return renderCodeCard(label, lines, cardWidth, unterminated)
	}
	shown := label
	if shown == "" {
		shown = "code"
	}
	out := make([]string, 0, len(lines)+2)
	out = append(out, mdDimColor+"── "+shown+" ──"+mdColorReset)
	out = append(out, highlightCode(label, lines)...)
	if !unterminated {
		out = append(out, mdDimColor+"──"+mdColorReset)
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
