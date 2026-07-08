package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// treeRow is one visible line in the tree — either a category header or,
// when its category is expanded, one of its exercises.
type treeRow struct {
	isCategory bool
	isLast     bool // exercise rows only: last child in its category (picks the tree connector)
	category   string
	status     catalog.ExerciseStatus
}

// treeModel is the NeetCode-roadmap-style practice picker: categories as
// collapsible branches, exercises as leaves underneath.
type treeModel struct {
	statuses []catalog.ExerciseStatus
	expanded map[string]bool
	cursor   int
	selected *catalog.ExerciseStatus
	back     bool
}

func newTreeModel(statuses []catalog.ExerciseStatus) treeModel {
	return treeModel{
		statuses: statuses,
		expanded: make(map[string]bool),
	}
}

// visibleRows groups statuses (already category-sorted by catalog.List)
// into category headers, each followed by its exercises only if that
// category is currently expanded.
func (m treeModel) visibleRows() []treeRow {
	type group struct {
		category string
		items    []catalog.ExerciseStatus
	}
	var groups []group
	for _, s := range m.statuses {
		if len(groups) == 0 || groups[len(groups)-1].category != s.Exercise.Category {
			groups = append(groups, group{category: s.Exercise.Category})
		}
		groups[len(groups)-1].items = append(groups[len(groups)-1].items, s)
	}

	var rows []treeRow
	for _, g := range groups {
		rows = append(rows, treeRow{isCategory: true, category: g.category})
		if m.expanded[g.category] {
			for i, s := range g.items {
				rows = append(rows, treeRow{status: s, isLast: i == len(g.items)-1})
			}
		}
	}
	return rows
}

func (m treeModel) rowCategory(row treeRow) string {
	if row.isCategory {
		return row.category
	}
	return row.status.Exercise.Category
}

func (m treeModel) categoryRowIndex(rows []treeRow, category string) int {
	for i, r := range rows {
		if r.isCategory && r.category == category {
			return i
		}
	}
	return 0
}

func (m treeModel) Init() tea.Cmd { return nil }

func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	rows := m.visibleRows()

	switch keyMsg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(rows)-1 {
			m.cursor++
		}
	case "right", "l":
		if m.cursor < len(rows) && rows[m.cursor].isCategory {
			m.expanded[rows[m.cursor].category] = true
		}
	case "left", "h":
		if m.cursor < len(rows) {
			cat := m.rowCategory(rows[m.cursor])
			m.expanded[cat] = false
			m.cursor = m.categoryRowIndex(m.visibleRows(), cat)
		}
	case "enter":
		if m.cursor < len(rows) {
			row := rows[m.cursor]
			if row.isCategory {
				m.expanded[row.category] = !m.expanded[row.category]
			} else {
				sel := row.status
				m.selected = &sel
				return m, tea.Quit
			}
		}
	case "q", "esc", "ctrl+c":
		m.back = true
		return m, tea.Quit
	}
	return m, nil
}

func (m treeModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.CompactBanner())
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("  Practice — pick a pattern"))
	b.WriteString("\n\n")

	rows := m.visibleRows()
	for i, row := range rows {
		highlighted := i == m.cursor
		if row.isCategory {
			solved, total := m.categoryCounts(row.category)
			b.WriteString(formatCategoryRow(row.category, m.expanded[row.category], highlighted, solved, total))
		} else {
			b.WriteString(formatTreeExerciseRow(row.status, row.isLast, highlighted))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("  ↑/↓ move · →/enter expand · ← collapse · enter (exercise) select · q back"))
	b.WriteString("\n")
	return b.String()
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

func formatCategoryRow(category string, expanded, highlighted bool, solved, total int) string {
	disclosure := "▸"
	if expanded {
		disclosure = "▾"
	}
	text := fmt.Sprintf("%s %-15s (%d/%d)", disclosure, category, solved, total)
	if highlighted {
		return cursorRowStyle.Render("❯ " + text)
	}
	return "  " + categoryStyle.Render(text)
}

func formatTreeExerciseRow(s catalog.ExerciseStatus, isLast, highlighted bool) string {
	connector := "├──"
	if isLast {
		connector = "└──"
	}

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

	lang := fmt.Sprintf("%-8s", s.Exercise.Language)
	title := fmt.Sprintf("%-36s", truncateTitle(s.Exercise.Title, 36))

	if highlighted {
		plain := fmt.Sprintf("%s %s %s %s", connector, lang, title, status)
		return cursorRowStyle.Render("❯ " + plain)
	}
	return "    " + connector + " " + langStyle.Render(lang) + " " + title + " " + statusStyle.Render(status)
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
