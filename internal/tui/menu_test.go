package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuModel_CursorStartsAtZero(t *testing.T) {
	m := newMenuModel()
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
}

func TestMenuModel_UpStaysAtTop(t *testing.T) {
	m := newMenuModel()
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newM.(menuModel).cursor != 0 {
		t.Error("cursor should stay at 0 when already at the top")
	}
}

func TestMenuModel_DownStopsAtLastOption(t *testing.T) {
	m := newMenuModel()
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(menuModel)
	}
	if m.cursor != len(menuLabels)-1 {
		t.Errorf("cursor = %d, want %d (last option)", m.cursor, len(menuLabels)-1)
	}
}

func TestMenuModel_NumberKeysJumpDirectly(t *testing.T) {
	m := newMenuModel()
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	if newM.(menuModel).cursor != 2 {
		t.Errorf("pressing 3 should jump cursor to index 2, got %d", newM.(menuModel).cursor)
	}
}

func TestMenuModel_EnterChoosesHighlightedOption(t *testing.T) {
	m := newMenuModel()
	m.cursor = 1 // Sandbox

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	mm := newM.(menuModel)
	if !mm.chosen {
		t.Fatal("expected chosen=true")
	}
	if mm.choice != menuSandbox {
		t.Errorf("choice = %v, want menuSandbox", mm.choice)
	}
}

func TestMenuModel_QRequestsQuitWithoutChoosing(t *testing.T) {
	m := newMenuModel()
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	mm := newM.(menuModel)
	if !mm.quit {
		t.Error("expected quit=true")
	}
	if mm.chosen {
		t.Error("expected chosen=false when quitting")
	}
}

func TestMenuModel_TickAdvancesPhaseAndReschedules(t *testing.T) {
	m := newMenuModel()
	newM, cmd := m.Update(tickMsg{})
	if cmd == nil {
		t.Fatal("expected tick to reschedule another tick command")
	}
	if newM.(menuModel).phase != 1 {
		t.Errorf("phase = %d, want 1", newM.(menuModel).phase)
	}
}
