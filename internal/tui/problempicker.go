package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

// problemPickerModel is a small popup listing the problems within one
// category — the second step of Practice, same bordered-dialog style as
// the category and language popups.
type problemPickerModel struct {
	category      string
	problems      []catalog.ProblemStatus
	cursor        int
	selected      *catalog.ProblemStatus
	back          bool
	width, height int
}

func newProblemPickerModel(all []catalog.ProblemStatus, category string) problemPickerModel {
	var problems []catalog.ProblemStatus
	for _, p := range all {
		if p.Category == category {
			problems = append(problems, p)
		}
	}
	return problemPickerModel{category: category, problems: problems}
}

func (m problemPickerModel) Init() tea.Cmd { return nil }

func (m problemPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.cursor < len(m.problems)-1 {
			m.cursor++
		}
	case "enter":
		sel := m.problems[m.cursor]
		m.selected = &sel
		return m, tea.Quit
	case "q", "esc", "ctrl+c":
		m.back = true
		return m, tea.Quit
	}
	return m, nil
}

func (m problemPickerModel) View() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(m.category))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a problem"))
	b.WriteString("\n\n")

	for i, p := range m.problems {
		label := fmt.Sprintf("%-30s", truncateTitle(p.Title, 30))
		status := "not attempted"
		statusStyle := checkDimStyle
		if p.Attempts > 0 {
			plural := "s"
			if p.Attempts == 1 {
				plural = ""
			}
			status = fmt.Sprintf("%d attempt%s", p.Attempts, plural)
			statusStyle = failStyle
			if p.Solved {
				status = "solved (" + status + ")"
				statusStyle = passStyle
			}
		}
		if i == m.cursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", label, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", label, statusStyle.Render(status)))
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

// RunProblemPicker shows the problem popup for a category and blocks
// until the user picks a problem (ok=true) or backs out (ok=false).
func RunProblemPicker(problems []catalog.ProblemStatus, category string) (sel catalog.ProblemStatus, ok bool, err error) {
	final, err := tea.NewProgram(newProblemPickerModel(problems, category), tea.WithAltScreen()).Run()
	if err != nil {
		return catalog.ProblemStatus{}, false, err
	}
	pm := final.(problemPickerModel)
	if pm.selected == nil {
		return catalog.ProblemStatus{}, false, nil
	}
	return *pm.selected, true, nil
}
