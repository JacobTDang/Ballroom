package exercise

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// RenderProblemText converts an exercise's problem.md into clean
// structured plain text -- what the practice session's editor pane
// actually opens (as problem.txt, see orchestrator.PrepareWorkspace).
// problem.md stays the authoring format and what the tutor's
// read_problem_statement tool reads; this render exists because a raw
// markdown buffer shows its markers (**bold**, backticks, #s) instead
// of structure, and the user reading a problem statement wants the
// structure, not the syntax.
//
// Transformations, all line-based and fence-aware (nothing inside a
// ``` fence is touched except the fence markers themselves):
//   - "# Title"    -> the bare title over a ═ underline
//   - "## Section" -> the bare section name over a ─ underline
//   - **bold** and `code` markers stripped, text kept
//   - fence markers dropped, fenced content indented two spaces
var (
	problemHeaderPattern = regexp.MustCompile(`^(#{1,4})\s+(.*)$`)
	// Bold and inline-code content may each span a line break -- authors
	// hard-wrap prose, so "**at\nleast twice**" (contains-duplicate) and
	// "`target =\n[x, y, z]`" (greedy-06) are real shapes -- hence [^*]+
	// / [^`]+ rather than the \n-excluding forms. Applied to whole
	// non-fence chunks, not single lines, so a pair can straddle the
	// break; leftmost-shortest matching keeps properly paired markers
	// pairing correctly (the full catalog renders marker-free -- see the
	// 459-file sweep this was verified with).
	problemBoldPattern       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	problemInlineCodePattern = regexp.MustCompile("`([^`]+)`")
)

func RenderProblemText(md string) string {
	lines := strings.Split(md, "\n")
	out := make([]string, 0, len(lines))
	chunk := make([]string, 0, len(lines))
	flushChunk := func() {
		if len(chunk) == 0 {
			return
		}
		text := strings.Join(chunk, "\n")
		text = problemBoldPattern.ReplaceAllString(text, "$1")
		text = problemInlineCodePattern.ReplaceAllString(text, "$1")
		out = append(out, strings.Split(text, "\n")...)
		chunk = chunk[:0]
	}
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			flushChunk()
			inFence = !inFence
			continue
		}
		if inFence {
			out = append(out, "  "+line)
			continue
		}
		if m := problemHeaderPattern.FindStringSubmatch(line); m != nil {
			flushChunk()
			underline := "─"
			if len(m[1]) == 1 {
				underline = "═"
			}
			// Headers carry inline markers too ("# Off-by-one: `MaxOf`")
			// -- strip them before sizing the underline.
			title := problemBoldPattern.ReplaceAllString(m[2], "$1")
			title = problemInlineCodePattern.ReplaceAllString(title, "$1")
			out = append(out, title, strings.Repeat(underline, utf8.RuneCountInString(title)))
			continue
		}
		chunk = append(chunk, line)
	}
	flushChunk()
	return strings.Join(out, "\n")
}
