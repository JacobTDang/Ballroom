package tutor

import "github.com/charmbracelet/lipgloss"

// The pane's shared lipgloss styles, consolidated here from model.go
// so the layout chrome reads as one system next to palette.go's
// colors. (Anything tests assert on by string uses palette.go's raw
// ansiFg/ansiBg escapes instead — see statusbar.go's doc comment for
// why lipgloss colors can vanish under `go test`.)

// textareaBoxStyle frames the input as a full double-ruled box — the
// opencode-style anchor of the 2026-07-17 restyle, chosen explicitly
// by the user, superseding the earlier top-rule-only design. That
// top-rule design itself replaced an even earlier full box whose
// *bright teal* left edge read as a persistent "sidebar" down the
// pane (a real live complaint); this frame is the dim structural
// paneRule — the same chrome as the editor cards' borders — so the
// box reads as structure, not a colored bar. The accent stays on the
// prompt glyph. PaddingLeft(1) keeps the prompt roughly aligned with
// the conversation's own left inset (viewportContentStyle's
// PaddingLeft(2)).
var textareaBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color(paneRule)).
	PaddingLeft(1)

// viewportContentStyle left-pads the scrolling conversation area — a
// real readability complaint found live: text printed flush against the
// pane's own left edge read as cramped. A lipgloss primitive instead of
// prefixing every print call site by hand (the old architecture's only
// option).
var viewportContentStyle = lipgloss.NewStyle().PaddingLeft(2)
