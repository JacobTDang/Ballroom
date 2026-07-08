package tui

import "github.com/charmbracelet/lipgloss"

// Shared lipgloss styles used across the menu, tree, and stats screens.
var (
	cursorRowStyle = lipgloss.NewStyle().Background(lipgloss.Color("141")).Foreground(lipgloss.Color("0")).Bold(true)
	categoryStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	langStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	passStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("120"))
	failStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("210"))
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
