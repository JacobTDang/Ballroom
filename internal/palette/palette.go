// Package palette is the one place the app's colors are defined.
//
// They used to live as string literals in four independent homes —
// internal/tui's lipgloss styles, internal/catalog's raw-ANSI banner,
// internal/tutor's pane constants, and the container's tmux/nvim
// config — using three different rendering mechanisms, with two files
// inside internal/tutor bypassing even that package's own palette. Any
// change had to be made four times and stayed correct only until
// someone forgot one.
//
// This package deliberately imports nothing from internal/: the tutor
// pane cannot import the TUI (the container builds its own binary from
// a subset of the tree), so a leaf package is the only shape all three
// Go homes can share. The container's static config files can't import
// Go at all — a drift test greps their hex values and asserts
// membership here instead.
package palette

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// The app's colors, named by role rather than by hue so a future
// re-theme changes values here without renaming every use site.
const (
	// Accents.
	Teal   = "#2FA6A6" // "good": passes, code, the tutor's own voice
	Pink   = "#E0468C" // the user's voice, separators
	Gold   = "#E8A93C" // attention: streaks, due markers, hints
	Red    = "#F03C3C" // failures
	Purple = "#9B5FB0" // selection, language column
	Blue   = "#3C7DC4" // category labels
	Orange = "#F0862E" // banner mid-tone
	Cyan   = "#3ED6D6" // disco-ball sparkle

	// Text.
	Cream    = "#F2EBDD" // brightest text
	PaleGray = "#D9D3C4" // ordinary dim text
	WarmGray = "#96918B" // pane metadata
	MidGray  = "#8B8680" // footers, subtitles
	DimGray  = "#6B6B6B" // disco-ball body

	// Structure.
	Rule      = "#3A3D4D" // borders, rules, card frames
	InputRule = "#1E5A5A" // the input box's dimmed-teal frame
	Ink       = "#000000" // text on a colored background

	// Surfaces.
	CardBg       = "#14151C" // editor-card body
	CardHeaderBg = "#1E2029" // card header bar, status bar row
	GutterFg     = "#5C5852" // line-number gutter
)

// The thinking aurora's own colors (internal/tutor/aurora.go). These
// are deliberately not the chrome accents: the glow is a decorative
// wash that has to read as light passing behind the frame, which wants
// a wider, cooler range than the app's four accents. They live here so
// there is still exactly one file to open when re-theming — not
// because they play the same role.
const (
	AuroraCyan       = "#19E3F0"
	AuroraPeriwinkle = "#6C7BE0"
	AuroraMagenta    = "#C13FD0"
	AuroraMint       = "#BFE8E0"
	AuroraBase       = "#1E50A2" // the quiet blue the blobs float over

	// PulseGlow is the activity dot's peak — a pale tint of Teal it
	// brightens toward, rather than dimming toward black (a dot that
	// dims reads as flickering off).
	PulseGlow = "#BFFCF7"
)

// Lip wraps a palette color for lipgloss styles.
func Lip(hex string) lipgloss.Color { return lipgloss.Color(hex) }

// RGB decomposes a palette color. It panics on a malformed constant
// rather than rendering garbage — every caller passes a constant from
// this file, so a failure here is a typo at build-out time, not a
// runtime condition worth degrading around.
func RGB(hex string) (r, g, b int) {
	if n, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b); n != 3 || err != nil {
		panic("palette: want #RRGGBB, got " + hex)
	}
	return r, g, b
}

// ANSIFg renders a color as a raw truecolor foreground escape, for the
// places that style strings by hand instead of through lipgloss — the
// tutor pane's markdown and cards, and the banner. Those exist because
// lipgloss routes colors through terminal-profile detection, which
// strips them entirely when there's no TTY (notably under `go test`,
// where several of these strings are asserted on).
func ANSIFg(hex string) string {
	r, g, b := RGB(hex)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// ANSIBg is ANSIFg's background counterpart.
func ANSIBg(hex string) string {
	r, g, b := RGB(hex)
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// All is every color defined here, for the drift test that checks the
// container's static config files don't invent their own.
func All() []string {
	return []string{
		Teal, Pink, Gold, Red, Purple, Blue, Orange, Cyan,
		Cream, PaleGray, WarmGray, MidGray, DimGray,
		Rule, InputRule, Ink,
		CardBg, CardHeaderBg, GutterFg,
	}
}

// Contains reports whether hex is a palette color, case-insensitively —
// the container's config files spell some of them in lowercase.
func Contains(hex string) bool {
	for _, c := range All() {
		if strings.EqualFold(c, hex) {
			return true
		}
	}
	return false
}
