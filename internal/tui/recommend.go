package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

// The "next up" picker turns the dashboard's recommendations into
// something you can act on. It gets its own key ("n") rather than a
// number because 1-5 are already the menu entries -- and it routes
// through launchExercise like every other path, so a saved draft still
// prompts before anything is overwritten.

func (m appModel) updateRecommend(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.recommendCursor > 0 {
			m.recommendCursor--
		}
	case "down", "j":
		if m.recommendCursor < len(m.recommendations)-1 {
			m.recommendCursor++
		}
	case "1", "2", "3":
		n := int(msg.String()[0] - '1')
		if n < len(m.recommendations) {
			m.recommendCursor = n
			return m.chooseRecommendation()
		}
	case "enter":
		return m.chooseRecommendation()
	case "esc", "q", "ctrl+c":
		m.stage = stageMain
	}
	return m, nil
}

func (m appModel) chooseRecommendation() (tea.Model, tea.Cmd) {
	if len(m.recommendations) == 0 || m.recommendCursor >= len(m.recommendations) {
		return m, nil
	}
	m.selectedProblem = m.recommendations[m.recommendCursor].Problem
	m.langCursor = 0
	m.stage = stageLanguage
	return m.resolveLanguageStage()
}

func (m appModel) renderRecommend() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(heading("Next up")))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("picked from what's due, what you've practiced least, and where you lose points"))
	b.WriteString("\n\n")

	if len(m.recommendations) == 0 {
		b.WriteString(checkDimStyle.Render("  nothing to suggest yet — practice a few problems first"))
		b.WriteString("\n")
		return b.String()
	}

	for i, r := range m.recommendations {
		label := fmt.Sprintf("%d. %s", i+1, r.Problem.Title)
		if i == m.recommendCursor {
			b.WriteString(cursorRowStyle.Render("❯ " + label))
		} else {
			b.WriteString("  " + label)
		}
		b.WriteString("\n")
		b.WriteString(checkDimStyle.Render("     " + r.Reason))
		b.WriteString("\n\n")
	}
	return b.String()
}

// refreshRecommendations recomputes on entry rather than caching: the
// underlying signals (due dates, solved counts) change as you practice,
// and a stale suggestion is worse than none.
func (m *appModel) refreshRecommendations() {
	m.recommendations = catalog.Recommend(m.problems, m.homeAttempts, time.Now())
	m.recommendCursor = 0
}
