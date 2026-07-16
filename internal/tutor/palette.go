package tutor

import "fmt"

// The tutor pane's palette, defined once -- every hardcoded escape
// that used to live at its point of use references these now, so the
// pane reads as one deliberate visual system instead of accumulated
// per-feature choices. Hex values shared with the wider app palette
// (internal/tui/styles.go, docker/tmux.conf) where a color plays the
// same role there.
const (
	// paneTeal is the pane's accent: inline/fenced code, the input
	// prompt glyph, the header's live dot -- same #2FA6A6 as the status
	// bar's own teal.
	paneTeal = "#2FA6A6"
	// panePink marks the user's own voice (the › in "you ›"), same
	// #E0468C as the status bar's separators.
	panePink = "#E0468C"
	// paneDimText is metadata text: fence labels, the "you" in the echo
	// prefix, the header's endpoint.
	paneDimText = "#96918B"
	// paneRule is structural chrome: the header rule, card borders.
	paneRule = "#3A3D4D"
	// paneInputRule is the input box's top rule -- a dimmed teal, so the
	// separation reads without the line competing with real content.
	paneInputRule = "#1E5A5A"
)

// ansiFg renders a palette hex color as a raw truecolor foreground
// escape, for the strings this package styles by hand (markdown
// styling operates on raw escapes rather than lipgloss styles -- see
// markdown.go). Called at package init to build the md* escape vars;
// panics on a malformed constant rather than silently rendering
// garbage.
func ansiFg(hex string) string {
	var r, g, b int
	if n, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b); n != 3 || err != nil {
		panic("tutor: ansiFg wants #RRGGBB, got " + hex)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}
