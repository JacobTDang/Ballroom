package tutor

import (
	"regexp"
	"strings"
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
	mdBoldOn  = "\x1b[1m"
	mdBoldOff = "\x1b[22m"
	// Inline code and fenced code share the pane's teal accent (the
	// same #2FA6A6 as the input box rule), so "code" reads as one
	// consistent signal everywhere it appears.
	mdCodeColor  = "\x1b[38;2;47;166;166m"
	mdColorReset = "\x1b[39m"
	// Fence labels ("python") render dim -- metadata, not content.
	mdDimColor = "\x1b[38;2;150;145;135m"
)

// styleMarkdown renders the tutor's markdown constructs as terminal
// styling: **bold**, `inline code`, #-headers, and ``` fenced blocks
// (code lines in the accent color, fence markers replaced by a dim
// language label / closing rule). Inside a fence nothing is
// transformed -- code legitimately contains * and ` characters.
func styleMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if !inFence {
				inFence = true
				label := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
				if label == "" {
					label = "code"
				}
				out = append(out, mdDimColor+"── "+label+" ──"+mdColorReset)
			} else {
				inFence = false
				out = append(out, mdDimColor+"──"+mdColorReset)
			}
			continue
		}
		if inFence {
			out = append(out, mdCodeColor+line+mdColorReset)
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
	return strings.Join(out, "\n")
}
