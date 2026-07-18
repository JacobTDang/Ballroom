package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

// stageStatsDetail (issue #252) is the Stats drill-down: every attempt
// ever logged against one exercise (tracker.ListAttemptsFor), reached
// with Enter on a row of stageStats' recent-activity list. Esc/q back
// out one level to stageStats, not stageMain, matching every other
// nested picker in this program (stageDSACategories -> stageCategories,
// stageProblems -> stageCategories/stageDSACategories, ...).

// maxVisibleDetail bounds how many attempt rows the drill-down shows at
// once -- same budget as maxVisibleProblems: this screen, unlike
// stageStats' list, owns its whole body (no totals/weak-spot sections
// above it eating into the vertical space).
const maxVisibleDetail = 12

// statsDash renders a column whose value is unknown, as opposed to
// zero/empty. hints_used, tutor_mode, and model are nullable columns
// (migration 002 added them; any attempt logged before that has NULL
// here, not 0 or "") -- printing strconv.Itoa(0) or "" for a nil
// pointer would misreport "we don't know" as "verified zero hints",
// exactly the silently-wrong number issue #252 called out. Notes
// (never nullable, but often blank) reuses the same dash for "nothing
// was written" -- it can't be confused with the nullable columns'
// "unknown" since an empty Notes IS the recorded value, not a missing
// one, but rendering it identically keeps every "nothing here" cell
// reading the same way.
const statsDash = "-"

// formatHintsUsed renders HintsUsed: a dash for nil (unknown -- an
// attempt logged before migration 002), otherwise the real count,
// including a genuine 0.
func formatHintsUsed(n *int) string {
	if n == nil {
		return statsDash
	}
	return strconv.Itoa(*n)
}

// formatOptional renders TutorMode/Model: a dash for nil or empty,
// otherwise the value as-is.
func formatOptional(s *string) string {
	if s == nil || *s == "" {
		return statsDash
	}
	return *s
}

// formatNotes renders Notes with the same dash convention as the
// nullable columns, so a blank cell always reads as "nothing here"
// rather than looking like a rendering gap.
func formatNotes(notes string) string {
	if strings.TrimSpace(notes) == "" {
		return statsDash
	}
	return notes
}

// formatTimeSpent renders TimeSpentMin to one decimal place, e.g. "12.5m".
func formatTimeSpent(min float64) string {
	return fmt.Sprintf("%.1fm", min)
}

// loadStatsDetail enters the drill-down for exerciseID: every attempt
// logged against it, newest first (attemptsForFn -> tracker.
// ListAttemptsFor). Mirrors loadStats' own error handling -- on failure
// the stage doesn't change (stays on stageStats, whose render now
// surfaces m.err -- see renderStats) rather than swapping to a broken
// detail screen with nothing loaded.
func (m appModel) loadStatsDetail(exerciseID string) appModel {
	attempts, err := attemptsForFn(m.cfg, exerciseID)
	if err != nil {
		m.err = err
		return m
	}
	m.err = nil
	m.statsDetailExerciseID = exerciseID
	m.statsDetailAttempts = attempts
	m.statsDetailCursor = 0
	m.stage = stageStatsDetail
	return m
}

// updateStatsDetail drives the drill-down: up/down/j/k move
// statsDetailCursor, pgup/pgdn page it, "?" opens help (same carve-out
// as everywhere else), and q/esc/ctrl+c back out one level to
// stageStats -- never straight to stageMain, so backing out of a
// history you drilled into from Stats lands back on the list you
// picked it from. Enter has nothing further to do here (there's no
// second level below a single exercise's history), so it's left
// unmapped -- a deliberate no-op, same as any other unrecognized key.
func (m appModel) updateStatsDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.statsDetailCursor > 0 {
			m.statsDetailCursor--
		}
	case "down", "j":
		if m.statsDetailCursor < len(m.statsDetailAttempts)-1 {
			m.statsDetailCursor++
		}
	case "pgup":
		m.statsDetailCursor = max(0, m.statsDetailCursor-maxVisibleDetail)
	case "pgdown":
		m.statsDetailCursor = min(max(len(m.statsDetailAttempts)-1, 0), m.statsDetailCursor+maxVisibleDetail)
	case "?":
		return m.openHelp()
	case "q", "esc", "ctrl+c":
		m.stage = stageStats
	}
	return m, nil
}

// renderStatsDetail draws one exercise's full attempt history: date,
// result, time spent, hints used, tutor mode, model, and notes (issue
// #252) -- everything ListAttemptsFor returns except the grade summary
// (already shown in aggregate on the Stats list's rubric weak-spots
// section) and turns (internal detail, not asked for here).
//
// Column widths are fixed and shared between the header and every data
// row so they can't drift out of alignment. Mode and model are
// truncated before padding (truncateTitle, not just %-Ns) since either
// can run long enough to blow the layout open -- a design tutor_mode
// like "behavioral-interviewer", or an OpenRouter slug like
// "nvidia/nemotron-3-super-120b-a12b:free" -- and %-Ns only pads short
// values, it never clips long ones.
func (m appModel) renderStatsDetail() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(m.statsDetailExerciseID))
	b.WriteString("\n")

	if len(m.statsDetailAttempts) > 0 {
		first := m.statsDetailAttempts[0]
		plural := "s"
		if len(m.statsDetailAttempts) == 1 {
			plural = ""
		}
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("%s · %s · %d attempt%s",
			catalog.DisplayCategory(first.Category), first.Language, len(m.statsDetailAttempts), plural)))
	} else {
		b.WriteString(checkDimStyle.Render("history"))
	}
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(renderFriendlyError("couldn't load this exercise's history", m.err))
		b.WriteString("\n\n")
	}

	if len(m.statsDetailAttempts) == 0 {
		b.WriteString(checkDimStyle.Render("no attempts logged for this exercise yet"))
		b.WriteString("\n\n")
		b.WriteString(checkDimStyle.Render("esc/q back"))
		return b.String()
	}

	header := fmt.Sprintf("%-11s %-9s %-8s %-7s %-16s %-20s %s",
		"date", "result", "time", "hints", "mode", "model", "notes")
	b.WriteString(checkDimStyle.Render(header))
	b.WriteString("\n")

	cursor := min(m.statsDetailCursor, max(len(m.statsDetailAttempts)-1, 0))
	start, end := sliceWindow(cursor, len(m.statsDetailAttempts), maxVisibleDetail)
	if start > 0 {
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("  ↑ %d more", start)))
		b.WriteString("\n")
	}
	for i := start; i < end; i++ {
		a := m.statsDetailAttempts[i]
		dateCol := fmt.Sprintf("%-11s", a.Date)
		resultCol := fmt.Sprintf("%-9s", a.Result)
		timeCol := fmt.Sprintf("%-8s", formatTimeSpent(a.TimeSpentMin))
		hintsCol := fmt.Sprintf("%-7s", formatHintsUsed(a.HintsUsed))
		modeCol := fmt.Sprintf("%-16s", truncateTitle(formatOptional(a.TutorMode), 14))
		modelCol := fmt.Sprintf("%-20s", truncateTitle(formatOptional(a.Model), 18))
		notesCol := truncateTitle(formatNotes(a.Notes), 32)

		if i == cursor {
			row := dateCol + " " + resultCol + " " + timeCol + " " + hintsCol + " " + modeCol + " " + modelCol + " " + notesCol
			b.WriteString(cursorRowStyle.Render("❯ " + row))
		} else {
			row := dateCol + " " + resultStyle(a.Result).Render(resultCol) + " " + timeCol + " " + hintsCol + " " + modeCol + " " + modelCol + " " + notesCol
			b.WriteString("  " + row)
		}
		b.WriteString("\n")
	}
	if end < len(m.statsDetailAttempts) {
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("  ↓ %d more", len(m.statsDetailAttempts)-end)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · pgup/pgdn page · esc/q back"))
	return b.String()
}
