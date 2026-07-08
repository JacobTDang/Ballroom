package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

// menuChoice is one of the video-game-style main menu options.
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

// menuModel is the title-screen selection: Practice / Sandbox / Stats,
// with an animated mosaic banner (this is the one idle screen where
// animation doesn't compete with reading a list, so it gets the full
// shimmer treatment).
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
		return m, nil

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

func (m menuModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.MosaicBanner(m.phase))
	b.WriteString("\n")

	for i, label := range menuLabels {
		if i == m.cursor {
			line := fmt.Sprintf("▸ %d. %s", i+1, label)
			b.WriteString(cursorRowStyle.Render(line))
			b.WriteString("\n  " + checkDimStyle.Render(menuDescriptions[i]))
		} else {
			line := fmt.Sprintf("  %d. ", i+1) + categoryStyle.Render(label)
			b.WriteString(line)
		}
		b.WriteString("\n\n")
	}

	b.WriteString(checkDimStyle.Render("  ↑/↓ or j/k move · 1/2/3 jump · enter select · q quit"))
	b.WriteString("\n")

	content := b.String()
	if m.width > 0 && m.height > 0 {
		return placeBlock(m.width, m.height, content)
	}
	return content
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
