package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func searchAppModel(t *testing.T) appModel {
	t.Helper()
	return appModel{
		cfg:   config.Config{DataDir: t.TempDir()},
		stage: stageMain,
		problems: []catalog.ProblemStatus{
			codingProblem("two-sum-01", "Two Sum", "easy"),
			codingProblem("rate-limiter-01", "Fixed-Window Rate Limiter", "medium"),
			codingProblem("trapping-rain-01", "Trapping Rain Water", "hard"),
		},
	}
}

// TestUpdateMain_SlashOpensSearch: the whole point is reaching any of
// 645 problems without knowing its category first, so search has to be
// available from the top of the app.
func TestUpdateMain_SlashOpensSearch(t *testing.T) {
	m := searchAppModel(t)
	got, _ := m.updateMain(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	next := got.(appModel)

	if next.stage != stageSearch {
		t.Fatalf("stage = %v, want stageSearch", next.stage)
	}
	if next.searchQuery != "" {
		t.Errorf("searchQuery = %q, want an empty box on entry", next.searchQuery)
	}
}

func TestUpdateSearch_TypingFiltersAcrossCategories(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch

	for _, r := range "rate" {
		got, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = got.(appModel)
	}

	// Weaker fuzzy hits may follow (subsequence matching means "rate"
	// also appears in "T-r-apping R-a-in Wa-te-r"), but the real match
	// has to lead.
	results := m.searchResults()
	if len(results) == 0 || results[0].ProblemID != "rate-limiter-01" {
		t.Errorf("first result = %v, want rate-limiter-01 to rank first", results)
	}
}

func TestUpdateSearch_BackspaceTrimsAndResetsCursor(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch
	m.searchQuery = "two"
	m.searchCursor = 0

	got, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyBackspace})
	next := got.(appModel)
	if next.searchQuery != "tw" {
		t.Errorf("searchQuery = %q, want %q", next.searchQuery, "tw")
	}
	if next.searchCursor != 0 {
		t.Errorf("searchCursor = %d, want reset to 0 when the match set changes", next.searchCursor)
	}
}

// TestUpdateSearch_EnterLaunchesTheSelectedProblem pins the payoff: a
// search result goes straight to the language pick, not back to a
// category listing the user was trying to avoid.
func TestUpdateSearch_EnterLaunchesTheSelectedProblem(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch
	m.searchQuery = "trapping"

	got, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyEnter})
	next := got.(appModel)

	if next.stage != stageLanguage {
		t.Fatalf("stage = %v, want stageLanguage", next.stage)
	}
	if next.selectedProblem.ProblemID != "trapping-rain-01" {
		t.Errorf("selectedProblem = %q, want trapping-rain-01", next.selectedProblem.ProblemID)
	}
}

func TestUpdateSearch_EnterOnNoMatchesDoesNothing(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch
	m.searchQuery = "zzzznope"

	got, cmd := m.updateSearch(tea.KeyMsg{Type: tea.KeyEnter})
	next := got.(appModel)
	if next.stage != stageSearch || cmd != nil {
		t.Error("Enter with no matches should stay put")
	}
}

func TestUpdateSearch_EscReturnsToMain(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch
	m.searchQuery = "two"

	got, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyEsc})
	if next := got.(appModel); next.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", next.stage)
	}
}

func TestUpdateSearch_ArrowsMoveWithinResults(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch
	m.searchQuery = "" // everything

	got, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyDown})
	next := got.(appModel)
	if next.searchCursor != 1 {
		t.Errorf("searchCursor = %d after down, want 1", next.searchCursor)
	}
	got, _ = next.updateSearch(tea.KeyMsg{Type: tea.KeyUp})
	if back := got.(appModel); back.searchCursor != 0 {
		t.Errorf("searchCursor = %d after up, want 0", back.searchCursor)
	}
}

func TestRenderSearch_ShowsQueryCategoryAndResults(t *testing.T) {
	m := searchAppModel(t)
	m.stage = stageSearch
	m.searchQuery = "rate"

	out := stripAnsiTUI(m.renderSearch())
	for _, want := range []string{"rate", "Fixed-Window Rate Limiter", "Two Pointers"} {
		if !strings.Contains(out, want) {
			t.Errorf("renderSearch missing %q:\n%s", want, out)
		}
	}
}

// TestFilterProblems_MatchesExerciseIDs: the in-category filter shares
// the global matcher, so an id typed into the picker works there too
// rather than silently finding nothing.
func TestFilterProblems_MatchesExerciseIDs(t *testing.T) {
	problems := []catalog.ProblemStatus{codingProblem("two-sum-01", "Two Sum", "easy")}
	problems[0].Variants[0].Exercise.ID = "two-sum-01-go"
	problems[0].Variants[0].Exercise.Language = exercise.LanguageGo

	if got := filterProblems(problems, "two-sum-01-go"); len(got) != 1 {
		t.Errorf("filterProblems by exercise id matched %d, want 1", len(got))
	}
}
