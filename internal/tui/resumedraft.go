package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/draft"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// The resume prompt stands between picking a problem and launching it,
// and exists because drafts persist now (internal/draft): reopening a
// problem you have unfinished work on must never silently overwrite
// that work, nor silently hand you yesterday's half-solution when you
// meant to start over. Both choices are non-destructive -- "start
// fresh" archives rather than deletes -- so neither answer can cost the
// user anything.

// resumeChoice indexes the prompt's two rows.
const (
	resumeChoiceResume = iota
	resumeChoiceFresh
)

// problemHasDraft reports whether any of p's language variants has a
// saved draft under dataDir — used by the problem picker (renderProblems)
// to mark rows so the resume prompt above is never a surprise. Backed
// by draft.Exists, a cheap existence check (no file content reads), so
// this is safe to call once per visible row on every render.
func problemHasDraft(dataDir string, p catalog.ProblemStatus) bool {
	for _, v := range p.Variants {
		if draft.Exists(dataDir, v.Exercise.ID) {
			return true
		}
	}
	return false
}

// launchExercise is the single door into a session. Every path that
// starts an exercise goes through it (the language picker's Enter and
// resolveLanguageStage's default-language fast path), which is what
// keeps the draft check from being bypassable -- an earlier design that
// only hooked the visible picker would have skipped the prompt entirely
// for anyone with a default language set.
func (m appModel) launchExercise(ex exercise.Exercise) (tea.Model, tea.Cmd) {
	d, ok, err := draft.Load(m.cfg.DataDir, ex.ID)
	if err != nil || !ok {
		// A draft that can't be read is treated as absent: the starter
		// is always a safe launch, and refusing to start a session over
		// an unreadable snapshot would be worse than ignoring it.
		return m.runExercise(ex, "")
	}
	m.pendingExercise = ex
	m.pendingDraft = d
	m.resumeCursor = resumeChoiceResume
	m.stage = stageResumeDraft
	return m, nil
}

// runExercise hands off to the orchestrator. draftDir is the directory
// whose solution files overlay the starter ("" for a clean start).
func (m appModel) runExercise(ex exercise.Exercise, draftDir string) (tea.Model, tea.Cmd) {
	m.exerciseToRun = ex
	m.draftDirToUse = draftDir
	m.outcome = outcomeRunExercise
	return m, tea.Quit
}

func (m appModel) updateResumeDraft(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.resumeCursor > resumeChoiceResume {
			m.resumeCursor--
		}
	case "down", "j":
		if m.resumeCursor < resumeChoiceFresh {
			m.resumeCursor++
		}
	case "enter":
		if m.resumeCursor == resumeChoiceFresh {
			// Archive before launching: the starter overlays nothing,
			// and the draft survives as previous.* if this was a
			// mistake.
			if err := draft.Archive(m.cfg.DataDir, m.pendingExercise.ID); err != nil {
				m.err = fmt.Errorf("set the draft aside: %w", err)
				return m, nil
			}
			return m.runExercise(m.pendingExercise, "")
		}
		return m.runExercise(m.pendingExercise, draft.Dir(m.cfg.DataDir, m.pendingExercise.ID))
	case "esc", "q", "ctrl+c":
		m.stage = stageLanguage
	}
	return m, nil
}

func (m appModel) renderResumeDraft() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("  Unfinished work"))
	b.WriteString("\n\n")
	b.WriteString(checkDimStyle.Render("  " + m.pendingExercise.Title))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render(fmt.Sprintf("  saved %s · %s",
		humanizeDraftAge(m.pendingDraft.Meta.SavedAt), pluralLines(m.pendingDraft.Meta.Lines))))
	b.WriteString("\n\n")

	for _, line := range m.pendingDraft.Preview {
		b.WriteString(langStyle.Render("  │ ") + checkDimStyle.Render(line) + "\n")
	}
	b.WriteString("\n")

	for i, row := range []string{"Resume where you left off", "Start fresh (keeps a copy)"} {
		if i == m.resumeCursor {
			b.WriteString(cursorRowStyle.Render("  ❯ "+row) + "\n")
			continue
		}
		b.WriteString("    " + row + "\n")
	}
	return b.String()
}

// humanizeDraftAge turns the snapshot timestamp into the phrasing a
// person would use. An unparseable timestamp degrades to "recently"
// rather than showing raw RFC3339 or an error -- the age is context,
// never the decision.
func humanizeDraftAge(savedAt string) string {
	t, err := time.Parse(time.RFC3339, savedAt)
	if err != nil {
		return "recently"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hr ago", int(d.Hours()))
	case d < 48*time.Hour:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}

func pluralLines(n int) string {
	if n == 1 {
		return "1 line"
	}
	return fmt.Sprintf("%d lines", n)
}
