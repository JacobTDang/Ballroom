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

// categoryOrder is the practice taxonomy: the broad tracks first, then
// the NeetCode 150 roadmap's own category sequence, then the design
// tracks — exercises are grouped by category in this order on the
// homepage and in the picker.
var categoryOrder = map[string]int{
	exercise.CategoryDSA:            0,
	exercise.CategoryDebug:          1,
	exercise.CategoryConcurrency:    2,
	exercise.CategoryImplementation: 3,
	exercise.CategoryAIAssisted:     4,

	exercise.CategoryArraysHashing:   5,
	exercise.CategoryTwoPointers:     6,
	exercise.CategorySlidingWindow:   7,
	exercise.CategoryStack:           8,
	exercise.CategoryBinarySearch:    9,
	exercise.CategoryLinkedList:      10,
	exercise.CategoryTrees:           11,
	exercise.CategoryTries:           12,
	exercise.CategoryHeap:            13,
	exercise.CategoryBacktracking:    14,
	exercise.CategoryGraphs:          15,
	exercise.CategoryAdvancedGraphs:  16,
	exercise.CategoryDP1D:            17,
	exercise.CategoryDP2D:            18,
	exercise.CategoryGreedy:          19,
	exercise.CategoryIntervals:       20,
	exercise.CategoryMathGeometry:    21,
	exercise.CategoryBitManipulation: 22,

	exercise.CategorySystemDesign: 23,
	exercise.CategoryAPIDesign:    24,
	exercise.CategoryOODesign:     25,
	exercise.CategoryBehavioral:   26,
}

// languageOrder ranks language variants of the same problem so Python
// sorts first (and is therefore the TUI language picker's default
// selection) — before this map existed, variants sorted by raw
// exercise ID string, which put cpp first purely as an alphabetical
// accident ("cpp" < "go" < "python").
var languageOrder = map[string]int{
	exercise.LanguagePython: 0,
	exercise.LanguageGo:     1,
	exercise.LanguageCpp:    2,

	// Design-kind session styles (the language slot carries the style --
	// see exercise.LanguageCoach's doc comment). Coach ranks first so
	// it's the picker default: the roadmap does each question coach-first,
	// interviewer on the second pass. Explicit ranks, not the
	// alphabetical accident this map exists to prevent.
	exercise.LanguageCoach:       3,
	exercise.LanguageInterviewer: 4,
}

// categoryDisplayNames holds the categories whose display label isn't
// just the raw id read as-is — an acronym (DSA), a NeetCode roadmap name
// with punctuation/casing a simple humanizer can't derive
// ("Heap / Priority Queue", "1-D Dynamic Programming"), or a
// multi-word phrase joined by "&". Anything not listed here falls back
// to title-casing each hyphen-separated word (see DisplayCategory).
var categoryDisplayNames = map[string]string{
	exercise.CategoryDSA:        "DSA",
	exercise.CategoryAIAssisted: "AI-Assisted",

	exercise.CategoryArraysHashing:   "Arrays & Hashing",
	exercise.CategoryHeap:            "Heap / Priority Queue",
	exercise.CategoryDP1D:            "1-D Dynamic Programming",
	exercise.CategoryDP2D:            "2-D Dynamic Programming",
	exercise.CategoryMathGeometry:    "Math & Geometry",
	exercise.CategoryBitManipulation: "Bit Manipulation",

	exercise.CategoryOODesign: "OO Design",
	// The title-case fallback would render "Api Design".
	exercise.CategoryAPIDesign: "API Design",
}

// DisplayCategory maps a raw category id to how it's shown in the UI.
// Exported so internal/tui can render category names the same way
// FormatSummary does below.
func DisplayCategory(category string) string {
	if name, ok := categoryDisplayNames[category]; ok {
		return name
	}
	words := strings.Split(category, "-")
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}

// dsaSubcategories is the set of NeetCode roadmap categories that nest
// under a single "DSA" entry in the practice picker, rather than being
// their own top-level categories — DSA is the practice picker's way of
// saying "algorithm problems, browse by topic", and every one of these
// is exactly that.
var dsaSubcategories = map[string]bool{
	exercise.CategoryArraysHashing:   true,
	exercise.CategoryTwoPointers:     true,
	exercise.CategorySlidingWindow:   true,
	exercise.CategoryStack:           true,
	exercise.CategoryBinarySearch:    true,
	exercise.CategoryLinkedList:      true,
	exercise.CategoryTrees:           true,
	exercise.CategoryTries:           true,
	exercise.CategoryHeap:            true,
	exercise.CategoryBacktracking:    true,
	exercise.CategoryGraphs:          true,
	exercise.CategoryAdvancedGraphs:  true,
	exercise.CategoryDP1D:            true,
	exercise.CategoryDP2D:            true,
	exercise.CategoryGreedy:          true,
	exercise.CategoryIntervals:       true,
	exercise.CategoryMathGeometry:    true,
	exercise.CategoryBitManipulation: true,
}

// TopLevelGroup returns the practice-picker top-level entry category
// belongs under: exercise.CategoryDSA for any NeetCode roadmap
// subcategory (Arrays & Hashing, Two Pointers, ...), or category itself
// for anything not grouped (debug, concurrency, implementation,
// ai-assisted).
func TopLevelGroup(category string) string {
	if dsaSubcategories[category] {
		return exercise.CategoryDSA
	}
	return category
}

// IsGroupedCategory reports whether category is nested under a
// top-level group (currently only DSA) rather than being its own
// top-level practice-picker entry.
func IsGroupedCategory(category string) bool {
	return dsaSubcategories[category]
}

// CategoryRank returns category's position in the canonical taxonomy
// order (categoryOrder) — exported so internal/tui can sort derived
// category lists (the top-level practice picker, and DSA's subcategory
// picker) consistently with FormatSummary without duplicating the
// ordering itself.
func CategoryRank(category string) int {
	return categoryOrder[category]
}

// ExerciseStatus is one exercise plus its practice history summary.
type ExerciseStatus struct {
	Exercise   exercise.Exercise
	Attempts   int
	LastResult string // tracker.ResultPass, tracker.ResultFail, or "" if never attempted
	// LastAttemptDate is the most recent attempt's date ("2006-01-02",
	// the format session submits log), or "" if never attempted --
	// what ReviewDue's spaced resurfacing runs on.
	LastAttemptDate string
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
		if exercises[i].ProblemID != exercises[j].ProblemID {
			return exercises[i].ProblemID < exercises[j].ProblemID
		}
		li, lj := languageOrder[exercises[i].Language], languageOrder[exercises[j].Language]
		if li != lj {
			return li < lj
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
			Exercise:        ex,
			Attempts:        len(a),
			LastResult:      lastResult(a),
			LastAttemptDate: lastAttemptDate(a),
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

// lastAttemptDate is lastResult's date counterpart.
func lastAttemptDate(attempts []tracker.Attempt) string {
	if len(attempts) == 0 {
		return ""
	}
	return attempts[len(attempts)-1].Date
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
		parts[i] = fmt.Sprintf("%s: %s", styled(colorBlue, DisplayCategory(cat)), fraction)
	}
	return strings.Join(parts, styled(colorDim, " · "))
}
