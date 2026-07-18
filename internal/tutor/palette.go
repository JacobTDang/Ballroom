package tutor

import "github.com/JacobTDang/Ballroom/internal/palette"

// The tutor pane's colors, named for their role in this pane and
// sourced from internal/palette so the pane, the host TUI, the banner,
// and the container's chrome can never drift apart. The local names
// stay because they carry pane-specific meaning ("the user's voice",
// "the card gutter") that the palette's role names deliberately don't.
const (
	// paneTeal is the pane's accent: inline/fenced code, the input
	// prompt glyph, the status bar's live dot.
	paneTeal = palette.Teal
	// panePink marks the user's own voice -- the accent bar down the
	// left of every message they send.
	panePink = palette.Pink
	// paneDimText is metadata: fence labels, tool-call detail, the
	// status bar's right half.
	paneDimText = palette.WarmGray
	// paneRule is structural chrome: rules and card borders.
	paneRule = palette.Rule
	// paneInputRule is the input box's frame -- a dimmed teal, so the
	// box reads as structure without competing with real content.
	paneInputRule = palette.InputRule

	// Editor-card colors (card.go): the card floats on its own
	// near-black background with a slightly lighter header bar.
	cardBg       = palette.CardBg
	cardHeaderBg = palette.CardHeaderBg
	// cardGutterFg is the line-number gutter -- warm and dim, present
	// without competing with the code.
	cardGutterFg = palette.GutterFg
	// The card header bar's three traffic-light dots.
	trafficRed  = palette.Red
	trafficGold = palette.Gold

	// paneStatusBg is the bottom status bar's row background -- the
	// same near-black as the card header, so the pane's two pieces of
	// fixed chrome read as one system.
	paneStatusBg = cardHeaderBg
)

// ansiFg and ansiBg render a palette color as a raw truecolor escape.
// The pane styles many strings by hand rather than through lipgloss
// because lipgloss routes colors through terminal-profile detection,
// which strips them when there's no TTY -- including under `go test`,
// where these strings are asserted on.
func ansiFg(hex string) string { return palette.ANSIFg(hex) }

func ansiBg(hex string) string { return palette.ANSIBg(hex) }
