package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

func fakeProblemsInCategory() []catalog.ProblemStatus {
	return []catalog.ProblemStatus{
		{ProblemID: "two-pointers-01", Title: "Two Sum II", Category: "pattern"},
		{ProblemID: "off-by-one-01", Title: "Off by one", Category: "debug"},
	}
}

func TestProblemPickerModel_ShowsOnlyProblemsInGivenCategory(t *testing.T) {
	m := newProblemPickerModel(fakeProblemsInCategory(), "pattern")
	if len(m.problems) != 1 {
		t.Fatalf("expected 1 problem filtered to pattern category, got %d", len(m.problems))
	}
	if m.problems[0].ProblemID != "two-pointers-01" {
		t.Errorf("expected two-pointers-01, got %q", m.problems[0].ProblemID)
	}
}

func TestProblemPickerModel_CursorStartsAtZero(t *testing.T) {
	m := newProblemPickerModel(fakeProblemsInCategory(), "pattern")
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
}

func TestProblemPickerModel_DownStopsAtLastProblem(t *testing.T) {
	problems := []catalog.ProblemStatus{
		{ProblemID: "a", Title: "A", Category: "pattern"},
		{ProblemID: "b", Title: "B", Category: "pattern"},
	}
	m := newProblemPickerModel(problems, "pattern")
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(problemPickerModel)
	}
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 (last problem)", m.cursor)
	}
}

func TestProblemPickerModel_EnterSelectsHighlightedProblemAndQuits(t *testing.T) {
	m := newProblemPickerModel(fakeProblemsInCategory(), "pattern")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	pm := newM.(problemPickerModel)
	if pm.selected == nil || pm.selected.ProblemID != "two-pointers-01" {
		t.Errorf("expected two-pointers-01 selected, got %+v", pm.selected)
	}
}

func TestProblemPickerModel_QRequestsBackWithoutSelecting(t *testing.T) {
	m := newProblemPickerModel(fakeProblemsInCategory(), "pattern")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	pm := newM.(problemPickerModel)
	if !pm.back {
		t.Error("expected back=true")
	}
	if pm.selected != nil {
		t.Error("expected no selection when backing out")
	}
}

func TestProblemPickerModel_ViewRendersCategoryName(t *testing.T) {
	m := newProblemPickerModel(fakeProblemsInCategory(), "pattern")
	out := stripAnsiTUI(m.View())
	if !strings.Contains(out, "pattern") {
		t.Errorf("expected the category name in the view, got:\n%s", out)
	}
}
