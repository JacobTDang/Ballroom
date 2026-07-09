package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

var popupBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#9B5FB0")).
	Padding(1, 3)

// langPickerModel is a small popup asking which language variant of a
// problem to practice, shown after a problem is selected from the tree
// rather than baking language into the problem selection itself.
type langPickerModel struct {
	problem       catalog.ProblemStatus
	cursor        int
	selected      *catalog.ExerciseStatus
	back          bool
	width, height int
}

func newLangPickerModel(problem catalog.ProblemStatus) langPickerModel {
	return langPickerModel{problem: problem}
}

func (m langPickerModel) Init() tea.Cmd { return nil }

func (m langPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = sizeMsg.Width, sizeMsg.Height
		return m, tea.ClearScreen
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.problem.Variants)-1 {
			m.cursor++
		}
	case "enter":
		sel := m.problem.Variants[m.cursor]
		m.selected = &sel
		return m, tea.Quit
	case "q", "esc", "ctrl+c":
		m.back = true
		return m, tea.Quit
	}
	return m, nil
}

func (m langPickerModel) View() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(m.problem.Title))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a language"))
	b.WriteString("\n\n")

	for i, v := range m.problem.Variants {
		lang := fmt.Sprintf("%-8s", v.Exercise.Language)
		status := "not attempted"
		statusStyle := checkDimStyle
		if v.LastResult != "" {
			status = v.LastResult
			statusStyle = failStyle
			if v.LastResult == tracker.ResultPass {
				statusStyle = passStyle
			}
		}
		if i == m.cursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", lang, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", langStyle.Render(lang), statusStyle.Render(status)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · q back"))

	content := popupBoxStyle.Render(b.String())
	if m.width > 0 && m.height > 0 {
		return placeBlock(m.width, m.height, content)
	}
	return content
}

// RunLangPicker shows the language popup for a problem and blocks until
// the user picks a variant (ok=true) or backs out (ok=false).
func RunLangPicker(problem catalog.ProblemStatus) (sel catalog.ExerciseStatus, ok bool, err error) {
	final, err := tea.NewProgram(newLangPickerModel(problem), tea.WithAltScreen()).Run()
	if err != nil {
		return catalog.ExerciseStatus{}, false, err
	}
	lm := final.(langPickerModel)
	if lm.selected == nil {
		return catalog.ExerciseStatus{}, false, nil
	}
	return *lm.selected, true, nil
}
