package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func fakeStatus(id string) catalog.ExerciseStatus {
	return catalog.ExerciseStatus{
		Exercise: exercise.Exercise{
			ID:           id,
			Title:        "Title for " + id,
			Category:     exercise.CategoryPattern,
			Language:     exercise.LanguageGo,
			TimeLimitMin: 20,
			TutorMode:    exercise.TutorModeHintsFirst,
			RepoPath:     "/fake/repo",
			TestCommand:  "true",
		},
	}
}

func TestPickerModel_CursorStartsAtZero(t *testing.T) {
	m := newPickerModel([]catalog.ExerciseStatus{fakeStatus("a")})
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
}

func TestPickerModel_UpStaysAtTop(t *testing.T) {
	m := newPickerModel([]catalog.ExerciseStatus{fakeStatus("a"), fakeStatus("b")})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newM.(pickerModel).cursor != 0 {
		t.Error("cursor should stay at 0 when already at the top")
	}
}

func TestPickerModel_DownMovesCursorAndStopsAtSandboxRow(t *testing.T) {
	m := newPickerModel([]catalog.ExerciseStatus{fakeStatus("a"), fakeStatus("b")})

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	pm := newM.(pickerModel)
	if pm.cursor != 1 {
		t.Fatalf("cursor = %d, want 1", pm.cursor)
	}

	newM2, _ := pm.Update(tea.KeyMsg{Type: tea.KeyDown})
	pm2 := newM2.(pickerModel)
	if pm2.cursor != 2 {
		t.Fatalf("cursor = %d, want 2 (sandbox row)", pm2.cursor)
	}

	newM3, _ := pm2.Update(tea.KeyMsg{Type: tea.KeyDown})
	if newM3.(pickerModel).cursor != 2 {
		t.Error("cursor should not move past the sandbox row")
	}
}

func TestPickerModel_VimKeysAlsoMoveCursor(t *testing.T) {
	m := newPickerModel([]catalog.ExerciseStatus{fakeStatus("a"), fakeStatus("b")})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if newM.(pickerModel).cursor != 1 {
		t.Error("expected j to move the cursor down")
	}
	newM2, _ := newM.(pickerModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if newM2.(pickerModel).cursor != 0 {
		t.Error("expected k to move the cursor up")
	}
}

func TestPickerModel_EnterSelectsHighlightedExercise(t *testing.T) {
	m := newPickerModel([]catalog.ExerciseStatus{fakeStatus("a"), fakeStatus("b")})
	m.cursor = 1

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	pm := newM.(pickerModel)
	if pm.selected == nil || pm.selected.Exercise.ID != "b" {
		t.Errorf("expected exercise \"b\" selected, got %+v", pm.selected)
	}
	if pm.sandbox {
		t.Error("should not have selected sandbox")
	}
}

func TestPickerModel_EnterOnSandboxRowSelectsSandbox(t *testing.T) {
	m := newPickerModel([]catalog.ExerciseStatus{fakeStatus("a")})
	m.cursor = 1 // sandbox row == len(statuses)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	pm := newM.(pickerModel)
	if !pm.sandbox {
		t.Error("expected sandbox=true")
	}
	if pm.selected != nil {
		t.Error("expected no exercise selected when choosing sandbox")
	}
}

func TestPickerModel_QRequestsQuit(t *testing.T) {
	m := newPickerModel(nil)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	if !newM.(pickerModel).quit {
		t.Error("expected quit=true")
	}
}

func TestPickerModel_CtrlCRequestsQuit(t *testing.T) {
	m := newPickerModel(nil)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected ctrl+c to return a quit command")
	}
	if !newM.(pickerModel).quit {
		t.Error("expected quit=true")
	}
}
