package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// menuChoice is one of the main menu options.
type menuChoice int

const (
	menuPractice menuChoice = iota
	menuSandbox
	menuStats
)

var menuLabels = []string{"Practice", "Sandbox", "Stats"}

var menuDescriptions = []string{
	"Pick a pattern and work through exercises",
	"Free practice, no grading, persists across sessions",
	"See your progress across categories",
}

// menuRightColWidth is the fixed content width of the right column —
// wide enough for the longest line (the keybinding hint) — so the
// selected row's highlight reads as a full-width bar rather than a
// tight box around just the label text.
const menuRightColWidth = 54

var (
	menuSubtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B8680"))
	menuRowHighlight  = lipgloss.NewStyle().Background(lipgloss.Color("#9B5FB0")).Foreground(lipgloss.Color("#000000")).Bold(true)
)

// menuModel is the title-screen selection: Practice / Sandbox / Stats, a
// two-column dashboard — a fixed-size disco ball on the left (a sparse
// subset of its mirror tiles glint with color on a timer) and the
// animated Ballroom banner + menu on the right, framed in a single
// bordered panel sized to fill most of the terminal.
type menuModel struct {
	cursor        int
	phase         int
	choice        menuChoice
	chosen        bool
	quit          bool
	width, height int
}

func newMenuModel() menuModel {
	return menuModel{}
}

func (m menuModel) Init() tea.Cmd {
	return tickCmd()
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, tea.ClearScreen

	case tickMsg:
		m.phase++
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuLabels)-1 {
				m.cursor++
			}
		case "1", "2", "3":
			n, _ := strconv.Atoi(msg.String())
			m.cursor = n - 1
		case "enter":
			m.choice = menuChoice(m.cursor)
			m.chosen = true
			return m, tea.Quit
		case "q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) renderRightColumn() string {
	var b strings.Builder
	for i, label := range menuLabels {
		numLabel := fmt.Sprintf("%d. %s", i+1, label)
		if i == m.cursor {
			row := fmt.Sprintf("❯ %-*s", menuRightColWidth-2, numLabel)
			b.WriteString(menuRowHighlight.Render(row))
			b.WriteString("\n  " + menuSubtitleStyle.Render(menuDescriptions[i]))
		} else {
			b.WriteString("  " + numLabel)
		}
		b.WriteString("\n\n\n")
	}

	b.WriteString("\n")
	b.WriteString(menuSubtitleStyle.Render("↑/↓ or j/k move · 1/2/3 jump · enter select · q quit"))
	return b.String()
}

func (m menuModel) View() string {
	right := m.renderRightColumn()
	if m.width == 0 || m.height == 0 {
		return right
	}
	panel := renderDashboardPanel(m.width, m.height, m.phase, right)
	return placeBlock(m.width, m.height, panel)
}

// RunMenu shows the main menu and blocks until the user picks an option
// (ok=true) or quits (ok=false).
func RunMenu() (choice menuChoice, ok bool, err error) {
	final, err := tea.NewProgram(newMenuModel(), tea.WithAltScreen()).Run()
	if err != nil {
		return 0, false, err
	}
	mm := final.(menuModel)
	if mm.quit || !mm.chosen {
		return 0, false, nil
	}
	return mm.choice, true, nil
}
