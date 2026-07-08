package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// pixelBarWidth is how many blocks wide a category's mini progress bar is.
const pixelBarWidth = 5

// pixelStatusIcon renders a small "sprite" per exercise that lights up as
// you make progress on it: dim/hollow when untouched, half-lit in the
// fail color after a failed attempt, fully lit with a sparkle once
// solved. Always 3 runes wide (before styling) so rows stay aligned.
func pixelStatusIcon(result string) string {
	switch result {
	case tracker.ResultPass:
		return passStyle.Render("▓▓") + sparkleStyle.Render("✦")
	case tracker.ResultFail:
		return failStyle.Render("▓░") + " "
	default:
		return checkDimStyle.Render("░░") + " "
	}
}

// treeModel is the NeetCode-roadmap-style practice picker: a real node
// graph — one root node branching to the 5 categories, and (one at a
// time, to keep the layout width-bounded) a category branching down to
// its exercises when you drill into it.
type treeModel struct {
	statuses      []catalog.ExerciseStatus
	categories    []string
	catCursor     int
	inExerciseRow bool
	exCursor      int
	selected      *catalog.ExerciseStatus
	back          bool
}

func newTreeModel(statuses []catalog.ExerciseStatus) treeModel {
	var categories []string
	seen := make(map[string]bool)
	for _, s := range statuses {
		if !seen[s.Exercise.Category] {
			seen[s.Exercise.Category] = true
			categories = append(categories, s.Exercise.Category)
		}
	}
	return treeModel{statuses: statuses, categories: categories}
}

func (m treeModel) exercisesFor(category string) []catalog.ExerciseStatus {
	var out []catalog.ExerciseStatus
	for _, s := range m.statuses {
		if s.Exercise.Category == category {
			out = append(out, s)
		}
	}
	return out
}

func (m treeModel) categoryCounts(category string) (solved, total int) {
	for _, s := range m.statuses {
		if s.Exercise.Category == category {
			total++
			if s.LastResult == tracker.ResultPass {
				solved++
			}
		}
	}
	return solved, total
}

func (m treeModel) Init() tea.Cmd { return nil }

func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "q", "esc", "ctrl+c":
		m.back = true
		return m, tea.Quit
	}

	if m.inExerciseRow {
		exs := m.exercisesFor(m.categories[m.catCursor])
		switch keyMsg.String() {
		case "left", "h":
			if m.exCursor > 0 {
				m.exCursor--
			}
		case "right", "l":
			if m.exCursor < len(exs)-1 {
				m.exCursor++
			}
		case "up", "k":
			m.inExerciseRow = false
		case "enter":
			sel := exs[m.exCursor]
			m.selected = &sel
			return m, tea.Quit
		}
		return m, nil
	}

	switch keyMsg.String() {
	case "left", "h":
		if m.catCursor > 0 {
			m.catCursor--
		}
	case "right", "l":
		if m.catCursor < len(m.categories)-1 {
			m.catCursor++
		}
	case "down", "j", "enter":
		m.inExerciseRow = true
		m.exCursor = 0
	}
	return m, nil
}

func (m treeModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.CompactBanner())
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("  Practice — pick a pattern"))
	b.WriteString("\n\n")

	catBoxes := make([]string, len(m.categories))
	catWidths := make([]int, len(m.categories))
	for i, cat := range m.categories {
		solved, total := m.categoryCounts(cat)
		catBoxes[i] = renderCategoryBox(cat, solved, total, !m.inExerciseRow && i == m.catCursor)
		catWidths[i] = lipgloss.Width(catBoxes[i])
	}
	catRow := joinBoxesHorizontal(catBoxes, boxGap)
	catCenters := boxCenters(catWidths, boxGap)
	totalWidth := lipgloss.Width(catRow)

	rootBox := renderRootBox()
	rootWidth := lipgloss.Width(rootBox)
	rootPad := centerOffset(rootWidth, spanAnchor(catCenters))
	rootRow := padLeft(rootBox, rootPad)
	if w := lipgloss.Width(rootRow); w > totalWidth {
		totalWidth = w
	}

	var exRowPadded string
	var exConn []string
	var exs []catalog.ExerciseStatus
	if m.inExerciseRow {
		exs = m.exercisesFor(m.categories[m.catCursor])
		exBoxes := make([]string, len(exs))
		exWidths := make([]int, len(exs))
		for i, s := range exs {
			exBoxes[i] = renderExerciseBox(s, i == m.exCursor)
			exWidths[i] = lipgloss.Width(exBoxes[i])
		}
		exRow := joinBoxesHorizontal(exBoxes, boxGap)
		exCenters := boxCenters(exWidths, boxGap)
		parentCenter := catCenters[m.catCursor]
		exPad := parentCenter - spanAnchor(exCenters)
		if exPad < 0 {
			exPad = 0
		}
		exRowPadded = padLeft(exRow, exPad)
		if w := lipgloss.Width(exRowPadded); w > totalWidth {
			totalWidth = w
		}
		shifted := make([]int, len(exCenters))
		for i, c := range exCenters {
			shifted[i] = c + exPad
		}
		exConn = connectorLines(totalWidth, parentCenter, shifted)
	}

	b.WriteString(rootRow + "\n")
	for _, l := range connectorLines(totalWidth, rootPad+rootWidth/2, catCenters) {
		b.WriteString(l + "\n")
	}
	b.WriteString(catRow + "\n")

	if m.inExerciseRow {
		for _, l := range exConn {
			b.WriteString(l + "\n")
		}
		b.WriteString(exRowPadded + "\n\n")
		b.WriteString(exerciseDetailLine(exs[m.exCursor]))
		b.WriteString("\n\n")
		b.WriteString(checkDimStyle.Render("  ←/→ move · ↑ back to categories · enter select · q back"))
	} else {
		b.WriteString("\n\n")
		b.WriteString(checkDimStyle.Render("  ←/→ move · ↓/enter expand · q back"))
	}
	b.WriteString("\n")
	return b.String()
}

func exerciseDetailLine(s catalog.ExerciseStatus) string {
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
	title := truncateTitle(s.Exercise.Title, 50)
	return fmt.Sprintf("  %s %s — %s", pixelStatusIcon(s.LastResult), hintStyle.Render(title), statusStyle.Render(status))
}

// RunTree shows the practice tree and blocks until the user selects an
// exercise (ok=true) or goes back to the main menu (ok=false).
func RunTree(statuses []catalog.ExerciseStatus) (sel catalog.ExerciseStatus, ok bool, err error) {
	final, err := tea.NewProgram(newTreeModel(statuses), tea.WithAltScreen()).Run()
	if err != nil {
		return catalog.ExerciseStatus{}, false, err
	}
	tm := final.(treeModel)
	if tm.selected == nil {
		return catalog.ExerciseStatus{}, false, nil
	}
	return *tm.selected, true, nil
}
