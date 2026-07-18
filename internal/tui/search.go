package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

// Search is the way into the catalog when you remember the problem but
// not where it lives. The picker's own filter only ever sees the
// category you already drilled into, which meant finding one of 645
// problems started with guessing its track -- exactly the step this
// skips. Reachable with "/" from the menu and from the pickers; a
// result goes straight to the language pick, never back through the
// category listing.

// searchResults is the live match set. Recomputed per keystroke rather
// than cached: catalog.Search over a few hundred problems is trivial,
// and a cache would be one more thing to invalidate.
func (m appModel) searchResults() []catalog.ProblemStatus {
	if strings.TrimSpace(m.searchQuery) == "" {
		return m.problems
	}
	return catalog.Search(m.problems, m.searchQuery)
}

func (m appModel) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	results := m.searchResults()
	switch msg.Type {
	case tea.KeyUp:
		if m.searchCursor > 0 {
			m.searchCursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.searchCursor < len(results)-1 {
			m.searchCursor++
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			// The match set just changed underneath the cursor; keeping
			// the old index would point at an unrelated problem.
			m.searchCursor = 0
		}
		return m, nil
	case tea.KeyEnter:
		if len(results) == 0 || m.searchCursor >= len(results) {
			return m, nil
		}
		m.selectedProblem = results[m.searchCursor]
		m.langCursor = 0
		m.stage = stageLanguage
		return m.resolveLanguageStage()
	case tea.KeyEsc, tea.KeyCtrlC:
		m.stage = stageMain
		return m, nil
	case tea.KeyRunes:
		// "?" with nothing typed yet opens help, matching the same
		// nothing-typed-yet carve-out stageProblems gives "q"/"?" (see
		// updateProblems) -- once a query is underway every rune feeds
		// it instead, since "?" could in principle be part of one.
		if m.searchQuery == "" && string(msg.Runes) == "?" {
			return m.openHelp()
		}
		m.searchQuery += string(msg.Runes)
		m.searchCursor = 0
		return m, nil
	case tea.KeySpace:
		m.searchQuery += " "
		m.searchCursor = 0
		return m, nil
	}
	return m, nil
}

func (m appModel) renderSearch() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Search"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("title, id, or track — across every category"))
	b.WriteString("\n\n")
	b.WriteString(checkDimStyle.Render("› ") + m.searchQuery)
	b.WriteString("\n\n")

	results := m.searchResults()
	if len(results) == 0 {
		b.WriteString(checkDimStyle.Render("  no matches"))
		b.WriteString("\n")
	}

	cursor := min(m.searchCursor, max(len(results)-1, 0))
	start, end := problemWindow(cursor, len(results))
	if start > 0 {
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("  ↑ %d more", start)))
		b.WriteString("\n")
	}
	for i := start; i < end; i++ {
		p := results[i]
		// The category is shown here and nowhere else in the pickers,
		// because a global result list is the one place you genuinely
		// can't infer it from context.
		if i == cursor {
			// The selected row draws its own badge unstyled, so the
			// cursor bar's background isn't punched through by the
			// badge's color (same reason the picker uses stripBadgeStyle).
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %-32s %-5s %s",
				truncateTitle(p.Title, 32), stripBadgeStyle(p), catalog.DisplayCategory(p.Category))))
			b.WriteString("\n")
			continue
		}
		b.WriteString(fmt.Sprintf("  %-32s %-5s %s\n",
			truncateTitle(p.Title, 32), renderDifficultyBadge(p), catalog.DisplayCategory(p.Category)))
	}
	if end < len(results) {
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("  ↓ %d more", len(results)-end)))
		b.WriteString("\n")
	}
	return b.String()
}
