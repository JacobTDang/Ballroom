package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

// categoryPickerModel is a small popup listing the practice categories —
// the first step of Practice, matching the same bordered-dialog style as
// the language popup rather than a spatial graph.
type categoryPickerModel struct {
	problems      []catalog.ProblemStatus
	categories    []string
	cursor        int
	selected      *string
	back          bool
	width, height int
}

func newCategoryPickerModel(problems []catalog.ProblemStatus) categoryPickerModel {
	var categories []string
	seen := make(map[string]bool)
	for _, p := range problems {
		if !seen[p.Category] {
			seen[p.Category] = true
			categories = append(categories, p.Category)
		}
	}
	return categoryPickerModel{problems: problems, categories: categories}
}

func (m categoryPickerModel) categoryCounts(category string) (solved, total int) {
	for _, p := range m.problems {
		if p.Category == category {
			total++
			if p.Solved {
				solved++
			}
		}
	}
	return solved, total
}

func (m categoryPickerModel) Init() tea.Cmd { return nil }

func (m categoryPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.cursor < len(m.categories)-1 {
			m.cursor++
		}
	case "enter":
		sel := m.categories[m.cursor]
		m.selected = &sel
		return m, tea.Quit
	case "q", "esc", "ctrl+c":
		m.back = true
		return m, tea.Quit
	}
	return m, nil
}

func (m categoryPickerModel) View() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Practice"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a category"))
	b.WriteString("\n\n")

	for i, cat := range m.categories {
		solved, total := m.categoryCounts(cat)
		label := fmt.Sprintf("%-16s", cat)
		status := fmt.Sprintf("%d/%d solved", solved, total)
		if i == m.cursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", label, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", categoryStyle.Render(label), checkDimStyle.Render(status)))
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

// RunCategoryPicker shows the category popup and blocks until the user
// picks a category (ok=true) or backs out (ok=false).
func RunCategoryPicker(problems []catalog.ProblemStatus) (category string, ok bool, err error) {
	final, err := tea.NewProgram(newCategoryPickerModel(problems), tea.WithAltScreen()).Run()
	if err != nil {
		return "", false, err
	}
	cm := final.(categoryPickerModel)
	if cm.selected == nil {
		return "", false, nil
	}
	return *cm.selected, true, nil
}
