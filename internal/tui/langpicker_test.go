package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

func fakeProblem() catalog.ProblemStatus {
	return catalog.ProblemStatus{
		ProblemID: "two-pointers-01",
		Title:     "Two Sum II (sorted input)",
		Category:  "pattern",
		Variants: []catalog.ExerciseStatus{
			fakeStatusIn("pattern", "two-pointers-01-go"),
			fakeStatusIn("pattern", "two-pointers-01-cpp"),
			fakeStatusIn("pattern", "two-pointers-01-python"),
		},
	}
}

func TestLangPickerModel_CursorStartsAtZero(t *testing.T) {
	m := newLangPickerModel(fakeProblem())
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
}

func TestLangPickerModel_UpStaysAtTop(t *testing.T) {
	m := newLangPickerModel(fakeProblem())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newM.(langPickerModel).cursor != 0 {
		t.Error("cursor should stay at 0 when already at the top")
	}
}

func TestLangPickerModel_DownStopsAtLastVariant(t *testing.T) {
	m := newLangPickerModel(fakeProblem())
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(langPickerModel)
	}
	if m.cursor != len(m.problem.Variants)-1 {
		t.Errorf("cursor = %d, want %d (last variant)", m.cursor, len(m.problem.Variants)-1)
	}
}

func TestLangPickerModel_EnterSelectsHighlightedVariantAndQuits(t *testing.T) {
	m := newLangPickerModel(fakeProblem())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown}) // move to cpp
	m = newM.(langPickerModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	lm := newM2.(langPickerModel)
	if lm.selected == nil || lm.selected.Exercise.ID != "two-pointers-01-cpp" {
		t.Errorf("expected two-pointers-01-cpp selected, got %+v", lm.selected)
	}
}

func TestLangPickerModel_QRequestsBackWithoutSelecting(t *testing.T) {
	m := newLangPickerModel(fakeProblem())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	lm := newM.(langPickerModel)
	if !lm.back {
		t.Error("expected back=true")
	}
	if lm.selected != nil {
		t.Error("expected no selection when backing out")
	}
}
