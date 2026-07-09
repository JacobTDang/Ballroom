package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// statsModel is a read-only progress summary: overall/per-category
// solve counts plus recent attempt history. Any keypress returns to the
// main menu.
type statsModel struct {
	statuses []catalog.ExerciseStatus
	recent   []tracker.Attempt // newest first
	back     bool
}

func newStatsModel(statuses []catalog.ExerciseStatus, recent []tracker.Attempt) statsModel {
	return statsModel{statuses: statuses, recent: recent}
}

func (m statsModel) Init() tea.Cmd { return nil }

func (m statsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.back = true
		return m, tea.Quit
	}
	return m, nil
}

func (m statsModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.CompactBanner())
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("  Stats"))
	b.WriteString("\n\n")

	total, attempted, solved := 0, 0, 0
	for _, s := range m.statuses {
		total++
		if s.Attempts > 0 {
			attempted++
		}
		if s.LastResult == tracker.ResultPass {
			solved++
		}
	}
	fmt.Fprintf(&b, "  %s solved · %s attempted · %d total exercises\n\n",
		passStyle.Render(fmt.Sprintf("%d", solved)),
		checkDimStyle.Render(fmt.Sprintf("%d", attempted)),
		total)

	b.WriteString("  " + catalog.FormatSummary(m.statuses) + "\n\n")

	if len(m.recent) == 0 {
		b.WriteString(checkDimStyle.Render("  No attempts yet — go practice something!"))
		b.WriteString("\n")
	} else {
		b.WriteString(hintStyle.Render("  Recent activity"))
		b.WriteString("\n")
		for _, a := range m.recent {
			resultStyle := failStyle
			if a.Result == tracker.ResultPass {
				resultStyle = passStyle
			}
			fmt.Fprintf(&b, "  %s  %-28s %s\n", a.Date, a.ExerciseID, resultStyle.Render(a.Result))
		}
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("  press any key to go back"))
	b.WriteString("\n")
	return b.String()
}

// RunStats shows the stats screen and blocks until the user presses any
// key to go back.
func RunStats(statuses []catalog.ExerciseStatus, recent []tracker.Attempt) error {
	_, err := tea.NewProgram(newStatsModel(statuses, recent), tea.WithAltScreen()).Run()
	return err
}
