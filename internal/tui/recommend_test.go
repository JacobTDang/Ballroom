package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
)

func recAppModel(t *testing.T) appModel {
	t.Helper()
	return appModel{
		cfg:   config.Config{DataDir: t.TempDir()},
		stage: stageMain,
		problems: []catalog.ProblemStatus{
			codingProblem("two-sum-01", "Two Sum", "easy"),
			codingProblem("valid-anagram-01", "Valid Anagram", "easy"),
		},
	}
}

// TestUpdateMain_NOpensRecommendations: the menu's number keys are
// already the menu items, so "next up" needs its own key.
func TestUpdateMain_NOpensRecommendations(t *testing.T) {
	m := recAppModel(t)
	got, _ := m.updateMain(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	next := got.(appModel)

	if next.stage != stageRecommend {
		t.Fatalf("stage = %v, want stageRecommend", next.stage)
	}
	if len(next.recommendations) == 0 {
		t.Error("no recommendations computed on entry")
	}
}

func TestUpdateRecommend_EnterLaunchesTheSelection(t *testing.T) {
	m := recAppModel(t)
	m.recommendations = catalog.Recommend(m.problems, nil, time.Now())
	m.stage = stageRecommend
	m.recommendCursor = 0
	if len(m.recommendations) == 0 {
		t.Skip("fixture produced no recommendations")
	}
	want := m.recommendations[0].Problem.ProblemID

	got, _ := m.updateRecommend(tea.KeyMsg{Type: tea.KeyEnter})
	next := got.(appModel)
	if next.selectedProblem.ProblemID != want {
		t.Errorf("selectedProblem = %q, want %q", next.selectedProblem.ProblemID, want)
	}
	if next.stage != stageLanguage {
		t.Errorf("stage = %v, want stageLanguage", next.stage)
	}
}

func TestUpdateRecommend_EscReturnsToMain(t *testing.T) {
	m := recAppModel(t)
	m.recommendations = catalog.Recommend(m.problems, nil, time.Now())
	m.stage = stageRecommend

	got, _ := m.updateRecommend(tea.KeyMsg{Type: tea.KeyEsc})
	if next := got.(appModel); next.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", next.stage)
	}
}

func TestRenderRecommend_ShowsReasons(t *testing.T) {
	m := recAppModel(t)
	m.recommendations = catalog.Recommend(m.problems, nil, time.Now())
	m.stage = stageRecommend
	if len(m.recommendations) == 0 {
		t.Skip("fixture produced no recommendations")
	}

	out := stripAnsiTUI(m.renderRecommend())
	if !strings.Contains(out, m.recommendations[0].Problem.Title) {
		t.Errorf("renderRecommend missing the problem title:\n%s", out)
	}
	if !strings.Contains(out, m.recommendations[0].Reason) {
		t.Errorf("renderRecommend missing the reason -- a suggestion you can't evaluate is noise:\n%s", out)
	}
}

func TestRenderRecommend_EmptyStateIsHonest(t *testing.T) {
	m := recAppModel(t)
	m.recommendations = nil
	m.stage = stageRecommend

	out := stripAnsiTUI(m.renderRecommend())
	if strings.TrimSpace(out) == "" {
		t.Error("empty recommendations rendered a blank screen")
	}
}
