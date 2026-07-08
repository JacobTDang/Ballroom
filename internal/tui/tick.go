package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickInterval drives the main menu's banner shimmer animation.
const tickInterval = 150 * time.Millisecond

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
