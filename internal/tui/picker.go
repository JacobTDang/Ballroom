package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

var (
	cursorRowStyle = lipgloss.NewStyle().Background(lipgloss.Color("141")).Foreground(lipgloss.Color("0")).Bold(true)
	categoryStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	langStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	passStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("120"))
	failStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("210"))
	sandboxStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
)

// Selection is what the user picked on the homepage.
type Selection struct {
	Sandbox  bool
	Exercise catalog.ExerciseStatus // zero value if Sandbox is true
}

// pickerModel is the arrow-key-navigable exercise list. The sandbox is
// modeled as one extra row at index len(statuses).
type pickerModel struct {
	statuses []catalog.ExerciseStatus
	cursor   int
	selected *catalog.ExerciseStatus
	sandbox  bool
	quit     bool
}

func newPickerModel(statuses []catalog.ExerciseStatus) pickerModel {
	return pickerModel{statuses: statuses}
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.cursor < len(m.statuses) {
			m.cursor++
		}
	case "enter":
		if m.cursor == len(m.statuses) {
			m.sandbox = true
		} else {
			sel := m.statuses[m.cursor]
			m.selected = &sel
		}
		return m, tea.Quit
	case "ctrl+c", "q":
		m.quit = true
		return m, tea.Quit
	}
	return m, nil
}

func (m pickerModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.Banner())
	b.WriteString("\n")

	header := fmt.Sprintf("    %-15s %-8s %-36s %s", "Category", "Lang", "Title", "Status")
	b.WriteString(hintStyle.Render(header))
	b.WriteString("\n")

	for i, s := range m.statuses {
		b.WriteString(formatPickerRow(s, i == m.cursor))
		b.WriteString("\n")
	}
	b.WriteString(formatSandboxRow(m.cursor == len(m.statuses)))
	b.WriteString("\n\n")
	b.WriteString(checkDimStyle.Render("  ↑/↓ or j/k move · enter select · q quit"))
	b.WriteString("\n")
	return b.String()
}

func formatPickerRow(s catalog.ExerciseStatus, highlighted bool) string {
	status := "not attempted"
	statusStyle := checkDimStyle
	if s.LastResult != "" {
		plural := "s"
		if s.Attempts == 1 {
			plural = ""
		}
		status = fmt.Sprintf("%s (%d attempt%s)", s.LastResult, s.Attempts, plural)
		statusStyle = failStyle
		if s.LastResult == tracker.ResultPass {
			statusStyle = passStyle
		}
	}

	category := fmt.Sprintf("%-15s", s.Exercise.Category)
	lang := fmt.Sprintf("%-8s", s.Exercise.Language)
	title := fmt.Sprintf("%-36s", truncateTitle(s.Exercise.Title, 36))

	if highlighted {
		plain := fmt.Sprintf("▸ %s %s %s %s", category, lang, title, status)
		return cursorRowStyle.Render(plain)
	}
	return fmt.Sprintf("  %s %s %s %s",
		categoryStyle.Render(category), langStyle.Render(lang), title, statusStyle.Render(status))
}

func formatSandboxRow(highlighted bool) string {
	text := "sandbox — free practice, no grading"
	if highlighted {
		return cursorRowStyle.Render("▸ " + text)
	}
	return "  " + sandboxStyle.Render(text)
}

func truncateTitle(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}

// RunPicker shows the exercise picker and blocks until the user selects
// something (ok=true) or quits (ok=false).
func RunPicker(statuses []catalog.ExerciseStatus) (sel Selection, ok bool, err error) {
	final, err := tea.NewProgram(newPickerModel(statuses), tea.WithAltScreen()).Run()
	if err != nil {
		return Selection{}, false, err
	}
	pm := final.(pickerModel)
	if pm.quit {
		return Selection{}, false, nil
	}
	if pm.sandbox {
		return Selection{Sandbox: true}, true, nil
	}
	return Selection{Exercise: *pm.selected}, true, nil
}
