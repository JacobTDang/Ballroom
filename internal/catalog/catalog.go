// Package catalog lists exercises and their practice status (from the
// tracker DB) for the homepage view.
package catalog

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// categoryOrder matches the taxonomy in interview_prep_plan.md Part 2 —
// exercises are grouped by category in this order on the homepage.
var categoryOrder = map[string]int{
	exercise.CategoryPattern:        0,
	exercise.CategoryDebug:          1,
	exercise.CategoryConcurrency:    2,
	exercise.CategoryImplementation: 3,
	exercise.CategoryAIAssisted:     4,
}

// ExerciseStatus is one exercise plus its practice history summary.
type ExerciseStatus struct {
	Exercise   exercise.Exercise
	Attempts   int
	LastResult string // tracker.ResultPass, tracker.ResultFail, or "" if never attempted
}

// List returns every valid exercise under cfg.ExercisesDir (skipping
// "_template" and any directory that isn't a well-formed exercise —
// one broken exercise shouldn't take down the whole homepage), sorted by
// category then id, each annotated with its attempt count and most recent
// result from the tracker DB.
func List(cfg config.Config) ([]ExerciseStatus, error) {
	entries, err := os.ReadDir(cfg.ExercisesDir)
	if err != nil {
		return nil, fmt.Errorf("catalog: read exercises dir: %w", err)
	}

	var exercises []exercise.Exercise
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "_template" {
			continue
		}
		ex, err := exercise.Load(cfg.ExercisePath(e.Name()))
		if err != nil {
			continue
		}
		exercises = append(exercises, ex)
	}

	sort.Slice(exercises, func(i, j int) bool {
		ci, cj := categoryOrder[exercises[i].Category], categoryOrder[exercises[j].Category]
		if ci != cj {
			return ci < cj
		}
		return exercises[i].ID < exercises[j].ID
	})

	attemptsByExercise, err := loadAttempts(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	statuses := make([]ExerciseStatus, len(exercises))
	for i, ex := range exercises {
		a := attemptsByExercise[ex.ID]
		statuses[i] = ExerciseStatus{
			Exercise:   ex,
			Attempts:   len(a),
			LastResult: lastResult(a),
		}
	}
	return statuses, nil
}

func loadAttempts(dbPath string) (map[string][]tracker.Attempt, error) {
	tr, err := tracker.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("catalog: open tracker: %w", err)
	}
	defer tr.Close()

	all, err := tr.ListAttempts()
	if err != nil {
		return nil, fmt.Errorf("catalog: list attempts: %w", err)
	}

	byExercise := make(map[string][]tracker.Attempt)
	for _, a := range all {
		byExercise[a.ExerciseID] = append(byExercise[a.ExerciseID], a)
	}
	return byExercise, nil
}

// lastResult returns the result of the most recent attempt (highest id,
// since ListAttempts returns them id-ascending), or "" if there are none.
func lastResult(attempts []tracker.Attempt) string {
	if len(attempts) == 0 {
		return ""
	}
	return attempts[len(attempts)-1].Result
}

// FormatTable renders a numbered table of exercises for the homepage.
//
// Fields are padded to their column width BEFORE being wrapped in ANSI
// color codes — styling first would make the invisible escape-code bytes
// count toward the padding width and break column alignment.
func FormatTable(statuses []ExerciseStatus) string {
	var b strings.Builder
	header := fmt.Sprintf("  %-3s %-15s %-8s %-36s %s", "#", "Category", "Lang", "Title", "Status")
	fmt.Fprintln(&b, styled(ansiBold+colorTeal, header))

	for i, s := range statuses {
		status := "not attempted"
		statusColor := colorDim
		if s.LastResult != "" {
			plural := "s"
			if s.Attempts == 1 {
				plural = ""
			}
			status = fmt.Sprintf("%s (%d attempt%s)", s.LastResult, s.Attempts, plural)
			statusColor = colorFail
			if s.LastResult == tracker.ResultPass {
				statusColor = colorPass
			}
		}

		num := fmt.Sprintf("%-3d", i+1)
		category := styled(colorBlue, fmt.Sprintf("%-15s", s.Exercise.Category))
		lang := styled(colorPurple, fmt.Sprintf("%-8s", s.Exercise.Language))
		title := fmt.Sprintf("%-36s", truncate(s.Exercise.Title, 36))

		fmt.Fprintf(&b, "  %s %s %s %s %s\n", num, category, lang, title, styled(statusColor, status))
		fmt.Fprintf(&b, "      %s\n", styled(colorDim, s.Exercise.ID))
	}
	return b.String()
}

// FormatSummary renders a "solved/total" count per category, in
// categoryOrder. Solved means the most recent attempt passed.
func FormatSummary(statuses []ExerciseStatus) string {
	type counts struct{ solved, total int }
	byCategory := make(map[string]*counts)
	var order []string

	for _, s := range statuses {
		c, ok := byCategory[s.Exercise.Category]
		if !ok {
			c = &counts{}
			byCategory[s.Exercise.Category] = c
			order = append(order, s.Exercise.Category)
		}
		c.total++
		if s.LastResult == tracker.ResultPass {
			c.solved++
		}
	}

	sort.Slice(order, func(i, j int) bool {
		return categoryOrder[order[i]] < categoryOrder[order[j]]
	})

	parts := make([]string, len(order))
	for i, cat := range order {
		c := byCategory[cat]
		fraction := styled(colorCream, fmt.Sprintf("%d/%d", c.solved, c.total))
		parts[i] = fmt.Sprintf("%s: %s", styled(colorBlue, cat), fraction)
	}
	return strings.Join(parts, styled(colorDim, " · "))
}

// FormatSandboxRow renders the sandbox menu option, styled consistently
// with FormatTable's rows.
func FormatSandboxRow(n int) string {
	num := fmt.Sprintf("%-3d", n)
	return fmt.Sprintf("  %s %s\n", num, styled(colorPink, "sandbox — free practice, no grading"))
}

// Prompt styles the input-prompt line shown at the bottom of the homepage.
func Prompt(s string) string {
	return styled(ansiBold+colorCream, s)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
