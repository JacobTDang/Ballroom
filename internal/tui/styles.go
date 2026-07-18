package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/palette"
)

// Shared lipgloss styles used across the menu, tree, and stats screens.
// Same bold/saturated palette as internal/catalog/theme.go, in hex since
// lipgloss takes hex directly — no green in this palette, so passStyle
// reuses teal (the palette's cool "good" color) and failStyle uses red
// (the dominant hue, an unambiguous fail).
var (
	cursorRowStyle = lipgloss.NewStyle().Background(palette.Lip(palette.Purple)).Foreground(palette.Lip(palette.Ink)).Bold(true)
	categoryStyle  = lipgloss.NewStyle().Foreground(palette.Lip(palette.Blue))
	langStyle      = lipgloss.NewStyle().Foreground(palette.Lip(palette.Purple))
	passStyle      = lipgloss.NewStyle().Foreground(palette.Lip(palette.Teal))
	failStyle      = lipgloss.NewStyle().Foreground(palette.Lip(palette.Red))
	// gaveUpStyle marks a tracker.ResultGaveUp attempt (issue #238) --
	// orange, distinct from passStyle's teal and failStyle's red, and
	// deliberately not gold: sparkleStyle/dueMarkerStyle already own
	// gold for "a nudge to do something", not a recorded outcome.
	gaveUpStyle  = lipgloss.NewStyle().Foreground(palette.Lip(palette.Orange))
	sparkleStyle = lipgloss.NewStyle().Foreground(palette.Lip(palette.Gold)).Bold(true)
	// dueMarkerStyle colors the picker's "· mock due" / "· review due"
	// nudges -- gold, but not sparkleStyle's bold: a nudge, not an alert.
	dueMarkerStyle = lipgloss.NewStyle().Foreground(palette.Lip(palette.Gold))
)

func truncateTitle(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
