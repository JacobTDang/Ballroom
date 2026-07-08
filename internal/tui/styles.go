package tui

import "github.com/charmbracelet/lipgloss"

// Shared lipgloss styles used across the menu, tree, and stats screens.
// Same bold/saturated palette as internal/catalog/theme.go, in hex since
// lipgloss takes hex directly — no green in this palette, so passStyle
// reuses teal (the palette's cool "good" color) and failStyle uses red
// (the dominant hue, an unambiguous fail).
var (
	cursorRowStyle = lipgloss.NewStyle().Background(lipgloss.Color("#9B5FB0")).Foreground(lipgloss.Color("#000000")).Bold(true)
	categoryStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C7DC4"))
	langStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#9B5FB0"))
	passStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#2FA6A6"))
	failStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#F03C3C"))
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
