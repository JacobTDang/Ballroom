package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func treeFixture() []catalog.ExerciseStatus {
	return []catalog.ExerciseStatus{
		fakeStatusIn("pattern", "two-pointers-01"),
		fakeStatusIn("pattern", "two-pointers-01-cpp"),
		fakeStatusIn("debug", "off-by-one-01-go"),
		fakeStatusIn("debug", "off-by-one-01-cpp"),
	}
}

// fakeStatus builds a minimal ExerciseStatus for Update()-logic tests
// that don't touch the real exercise catalog or tracker DB.
func fakeStatus(id string) catalog.ExerciseStatus {
	return catalog.ExerciseStatus{
		Exercise: exercise.Exercise{
			ID:       id,
			Title:    id,
			Category: "pattern",
			Language: "go",
		},
	}
}

func fakeStatusIn(category, id string) catalog.ExerciseStatus {
	s := fakeStatus(id)
	s.Exercise.Category = category
	return s
}

func TestTreeModel_StartsOnCategoryRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	if m.inExerciseRow {
		t.Error("expected to start on the category row, not the exercise row")
	}
	if m.catCursor != 0 {
		t.Errorf("catCursor = %d, want 0", m.catCursor)
	}
	if len(m.categories) != 2 {
		t.Fatalf("expected 2 categories (pattern, debug), got %v", m.categories)
	}
}

func TestTreeModel_RightMovesCatCursorAndStopsAtEnd(t *testing.T) {
	m := newTreeModel(treeFixture())
	for i := 0; i < 5; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
		m = newM.(treeModel)
	}
	if m.catCursor != 1 { // only 2 categories: indices 0,1
		t.Errorf("catCursor = %d, want 1 (last category)", m.catCursor)
	}
}

func TestTreeModel_LeftStaysAtZero(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if newM.(treeModel).catCursor != 0 {
		t.Error("catCursor should stay at 0 when already at the first category")
	}
}

func TestTreeModel_DownEntersExerciseRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd != nil {
		t.Fatal("entering the exercise row should not quit the program")
	}
	tm := newM.(treeModel)
	if !tm.inExerciseRow {
		t.Fatal("expected down to move focus into the exercise row")
	}
	if tm.exCursor != 0 {
		t.Errorf("exCursor = %d, want 0", tm.exCursor)
	}
}

func TestTreeModel_EnterOnCategoryAlsoEntersExerciseRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("entering the exercise row should not quit the program")
	}
	if !newM.(treeModel).inExerciseRow {
		t.Fatal("expected enter on a category to also enter its exercise row")
	}
}

func TestTreeModel_LeftRightMoveWithinExerciseRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown}) // pattern has 2 exercises
	tm := newM.(treeModel)

	newM2, _ := tm.Update(tea.KeyMsg{Type: tea.KeyRight})
	tm2 := newM2.(treeModel)
	if tm2.exCursor != 1 {
		t.Fatalf("exCursor = %d, want 1", tm2.exCursor)
	}

	// stops at the last exercise (pattern has 2: index 0,1)
	newM3, _ := tm2.Update(tea.KeyMsg{Type: tea.KeyRight})
	if newM3.(treeModel).exCursor != 1 {
		t.Error("exCursor should stop at the last exercise in the category")
	}
}

func TestTreeModel_ExerciseRowStaysAlignedUnderLeftmostCategory(t *testing.T) {
	statuses := treeFixture()
	// pattern is leftmost; a 3rd exercise makes its exercise row wider
	// than needed to align under pattern's own (further-left) center —
	// exactly the reported bug scenario: the category ends up narrower/
	// further left than its own exercise row wants to sit.
	statuses = append(statuses, fakeStatusIn("pattern", "two-pointers-01-python"))
	m := newTreeModel(statuses)

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown}) // expand pattern (leftmost)
	view := stripAnsiTUI(newM.(treeModel).View())
	lines := strings.Split(view, "\n")

	rootLineIdx := -1
	for i, l := range lines {
		if strings.Contains(l, "PRACTICE") {
			rootLineIdx = i
			break
		}
	}
	if rootLineIdx == -1 {
		t.Fatal("could not find the PRACTICE root in the rendered view")
	}

	catLineIdx := -1
	for i := rootLineIdx + 1; i < len(lines); i++ {
		if strings.Contains(lines[i], "pattern") {
			catLineIdx = i
			break
		}
	}
	if catLineIdx == -1 {
		t.Fatal("could not find the pattern category row in the rendered view")
	}
	catStart := strings.Index(lines[catLineIdx], "pattern")
	catEnd := catStart + len("pattern")

	// category row -> exercise connector's single parent stem -> spine ->
	// child stems -> exercise boxes.
	stemLineIdx := catLineIdx + 2
	if stemLineIdx >= len(lines) {
		t.Fatal("expected an exercise connector stem row below the category")
	}
	stemCol := strings.IndexRune(lines[stemLineIdx], '│')
	if stemCol == -1 {
		t.Fatal("expected a '│' stem connecting the category down to its exercises")
	}
	if stemCol < catStart || stemCol > catEnd {
		t.Errorf("exercise connector stem at col %d is outside the pattern box's own span [%d,%d) — misaligned, connects to boxes positioned elsewhere",
			stemCol, catStart, catEnd)
	}
}

func TestTreeModel_UpFromExerciseRowReturnsToCategoryRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	tm := newM.(treeModel)

	newM2, _ := tm.Update(tea.KeyMsg{Type: tea.KeyUp})
	tm2 := newM2.(treeModel)
	if tm2.inExerciseRow {
		t.Fatal("expected up to return focus to the category row")
	}
	if tm2.catCursor != 0 {
		t.Errorf("catCursor should be preserved as 0, got %d", tm2.catCursor)
	}
}

func TestTreeModel_EnterOnExerciseSelectsAndQuits(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown}) // enter pattern's exercise row
	tm := newM.(treeModel)

	newM2, cmd := tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter on an exercise to return a quit command")
	}
	tm2 := newM2.(treeModel)
	if tm2.selected == nil || tm2.selected.Exercise.ID != "two-pointers-01" {
		t.Errorf("expected two-pointers-01 selected, got %+v", tm2.selected)
	}
}

func TestTreeModel_QRequestsBackFromEitherRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	if !newM.(treeModel).back {
		t.Error("expected back=true")
	}

	inExercise, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	newM2, cmd2 := inExercise.(treeModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd2 == nil {
		t.Fatal("expected q to quit even from the exercise row")
	}
	if !newM2.(treeModel).back {
		t.Error("expected back=true from the exercise row too")
	}
}

func TestPixelStatusIcon_DiffersByStatus(t *testing.T) {
	notAttempted := stripAnsiTUI(pixelStatusIcon(""))
	fail := stripAnsiTUI(pixelStatusIcon(tracker.ResultFail))
	pass := stripAnsiTUI(pixelStatusIcon(tracker.ResultPass))

	if notAttempted == fail || notAttempted == pass || fail == pass {
		t.Errorf("expected 3 visually distinct icons, got %q / %q / %q", notAttempted, fail, pass)
	}
	if !strings.Contains(pass, "✦") {
		t.Errorf("expected a sparkle in the solved icon, got %q", pass)
	}
}
