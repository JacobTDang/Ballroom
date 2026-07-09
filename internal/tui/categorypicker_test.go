package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

func fakeProblems() []catalog.ProblemStatus {
	return []catalog.ProblemStatus{
		{ProblemID: "two-pointers-01", Title: "Two Sum II", Category: "pattern"},
		{ProblemID: "off-by-one-01", Title: "Off by one", Category: "debug", Solved: true},
	}
}

func TestCategoryPickerModel_CursorStartsAtZero(t *testing.T) {
	m := newCategoryPickerModel(fakeProblems())
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
	if len(m.categories) != 2 {
		t.Fatalf("expected 2 categories (pattern, debug), got %v", m.categories)
	}
}

func TestCategoryPickerModel_DownStopsAtLastCategory(t *testing.T) {
	m := newCategoryPickerModel(fakeProblems())
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(categoryPickerModel)
	}
	if m.cursor != len(m.categories)-1 {
		t.Errorf("cursor = %d, want %d (last category)", m.cursor, len(m.categories)-1)
	}
}

func TestCategoryPickerModel_UpStaysAtTop(t *testing.T) {
	m := newCategoryPickerModel(fakeProblems())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newM.(categoryPickerModel).cursor != 0 {
		t.Error("cursor should stay at 0 when already at the top")
	}
}

func TestCategoryPickerModel_EnterSelectsHighlightedCategoryAndQuits(t *testing.T) {
	m := newCategoryPickerModel(fakeProblems())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown}) // move to debug
	m = newM.(categoryPickerModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	cm := newM2.(categoryPickerModel)
	if cm.selected == nil || *cm.selected != "debug" {
		t.Errorf("expected debug selected, got %+v", cm.selected)
	}
}

func TestCategoryPickerModel_QRequestsBackWithoutSelecting(t *testing.T) {
	m := newCategoryPickerModel(fakeProblems())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	cm := newM.(categoryPickerModel)
	if !cm.back {
		t.Error("expected back=true")
	}
	if cm.selected != nil {
		t.Error("expected no selection when backing out")
	}
}
