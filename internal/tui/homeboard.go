package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// The home dashboard ("homeboard"): the block under the main menu that
// makes the home screen a status page instead of a bare list -- today's
// pick, streak and due counts, per-track progress bars, and the last
// few attempts. All of it derives from data the app already loads for
// other screens (catalog.List + the tracker); a fresh install renders
// nothing here rather than a wall of zeros.

// homeTracks are the practice tracks the progress bars cover, in
// picker order. Every category the picker lists gets a row here too
// (issue #255) -- ai-assisted used to be the one exception, on the
// theory that a single-problem bar reads as noise, but that just made
// the dashboard inconsistent with the picker (which lists it) without
// actually hiding it from practice. It's real content the user can
// work through, so it gets a row like everything else; a 0/1 or 1/1
// bar is honest, not noisy.
var homeTracks = []struct {
	group string
	label string
}{
	{exercise.CategoryDSA, "DSA"},
	{exercise.CategoryDebug, "Debug"},
	{exercise.CategoryConcurrency, "Concurrency"},
	{exercise.CategoryImplementation, "Implementation"},
	{exercise.CategoryAIAssisted, "AI-Assisted"},
	{exercise.CategorySystemDesign, "System Design"},
	{exercise.CategoryAPIDesign, "API Design"},
	{exercise.CategoryOODesign, "OO Design"},
	{exercise.CategoryBehavioral, "Behavioral"},
}

const progressBarWidth = 14

// progressBar renders solved/total as a fixed-width bar. Any progress
// at all shows at least one filled cell -- 1/149 flooring to an empty
// bar would read as "nothing done", the opposite of encouragement.
func progressBar(solved, total, width int) string {
	filled := 0
	if total > 0 {
		filled = solved * width / total
		if solved > 0 && filled == 0 {
			filled = 1
		}
	}
	return passStyle.Render(strings.Repeat("▓", filled)) +
		checkDimStyle.Render(strings.Repeat("░", width-filled))
}

// practiceStreak counts consecutive practice days ending today -- or
// ending yesterday, so the streak reads as still alive before today's
// session rather than resetting to zero at midnight. Dates that don't
// parse (hand-edited rows) are skipped.
func practiceStreak(attempts []tracker.Attempt, now time.Time) int {
	days := make(map[string]bool, len(attempts))
	for _, a := range attempts {
		if _, err := time.Parse("2006-01-02", a.Date); err == nil {
			days[a.Date] = true
		}
	}
	day := now
	if !days[day.Format("2006-01-02")] {
		day = day.AddDate(0, 0, -1) // today not practiced yet: count from yesterday
	}
	streak := 0
	for days[day.Format("2006-01-02")] {
		streak++
		day = day.AddDate(0, 0, -1)
	}
	return streak
}

// homeRecentLimit is how many recent attempts the homeboard shows --
// a glance at momentum, not the Stats screen's full history.
const homeRecentLimit = 3

// renderHomeboard renders the dashboard block. Empty inputs render an
// empty string -- the menu stands alone on a fresh install.
func renderHomeboard(problems []catalog.ProblemStatus, attempts []tracker.Attempt, now time.Time) string {
	if len(problems) == 0 {
		return ""
	}
	var b strings.Builder

	if pick, ok := catalog.DailyPick(problems, now); ok {
		b.WriteString(checkDimStyle.Render("today: "))
		b.WriteString(hintStyle.Render(pick.Title))
		b.WriteString("\n")
	}

	statusBits := make([]string, 0, 2)
	if streak := practiceStreak(attempts, now); streak > 0 {
		plural := "s"
		if streak == 1 {
			plural = ""
		}
		statusBits = append(statusBits, sparkleStyle.Render(fmt.Sprintf("streak %d day%s", streak, plural)))
	}
	mockDue, reviewDue := 0, 0
	for _, p := range problems {
		if catalog.MockDue(p) {
			mockDue++
		} else if catalog.ReviewDue(p, now) {
			reviewDue++
		}
	}
	if mockDue+reviewDue > 0 {
		statusBits = append(statusBits, dueMarkerStyle.Render(fmt.Sprintf("due: %d mock · %d review", mockDue, reviewDue)))
	}
	if len(statusBits) > 0 {
		b.WriteString(strings.Join(statusBits, checkDimStyle.Render("   ")))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	for _, track := range homeTracks {
		solved, total := groupCounts(problems, track.group)
		if total == 0 {
			continue
		}
		fmt.Fprintf(&b, "%-14s %s %s\n",
			track.label, progressBar(solved, total, progressBarWidth),
			checkDimStyle.Render(fmt.Sprintf("%d/%d", solved, total)))
	}

	if len(attempts) > 0 {
		b.WriteString("\n")
		b.WriteString(checkDimStyle.Render("recent: "))
		start := len(attempts) - homeRecentLimit
		if start < 0 {
			start = 0
		}
		parts := make([]string, 0, homeRecentLimit)
		// Newest first -- the thing you just did leads.
		for i := len(attempts) - 1; i >= start; i-- {
			a := attempts[i]
			mark := failStyle.Render("✗")
			switch a.Result {
			case tracker.ResultPass:
				mark = passStyle.Render("✓")
			case tracker.ResultGaveUp:
				mark = gaveUpStyle.Render("~")
			}
			parts = append(parts, mark+" "+a.ExerciseID)
		}
		b.WriteString(strings.Join(parts, checkDimStyle.Render("  · ")))
		b.WriteString("\n")
	}

	return b.String()
}
