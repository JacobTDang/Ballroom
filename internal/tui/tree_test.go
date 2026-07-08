package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
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

func TestTreeModel_StartsAllCollapsedShowingOnlyCategories(t *testing.T) {
	m := newTreeModel(treeFixture())
	rows := m.visibleRows()
	if len(rows) != 2 {
		t.Fatalf("expected 2 visible rows (pattern, debug categories only), got %d", len(rows))
	}
	for _, r := range rows {
		if !r.isCategory {
			t.Errorf("expected only category rows while collapsed, got exercise row %+v", r)
		}
	}
}

func TestTreeModel_RightExpandsCategoryRevealingChildren(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	tm := newM.(treeModel)

	rows := tm.visibleRows()
	if len(rows) != 4 { // pattern (expanded, 2 children) + debug (collapsed)
		t.Fatalf("expected 4 visible rows after expanding pattern, got %d: %+v", len(rows), rows)
	}
	if !rows[0].isCategory || rows[0].category != "pattern" {
		t.Fatalf("row 0 should be the pattern category, got %+v", rows[0])
	}
	if rows[1].isCategory || rows[1].status.Exercise.Category != "pattern" {
		t.Fatalf("row 1 should be a pattern exercise, got %+v", rows[1])
	}
}

func TestTreeModel_EnterOnCategoryTogglesExpandWithoutQuitting(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("enter on a category row should not quit the program")
	}
	tm := newM.(treeModel)
	if len(tm.visibleRows()) != 4 {
		t.Fatal("expected enter on category to expand it")
	}
}

func TestTreeModel_DownMovesIntoExpandedChildren(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight}) // expand pattern
	tm := newM.(treeModel)

	newM2, _ := tm.Update(tea.KeyMsg{Type: tea.KeyDown})
	tm2 := newM2.(treeModel)
	if tm2.cursor != 1 {
		t.Fatalf("cursor = %d, want 1 (first child of pattern)", tm2.cursor)
	}
}

func TestTreeModel_LeftCollapsesAndMovesCursorToParent(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight}) // expand pattern
	tm := newM.(treeModel)
	newM2, _ := tm.Update(tea.KeyMsg{Type: tea.KeyDown}) // move to first child
	tm2 := newM2.(treeModel)

	newM3, _ := tm2.Update(tea.KeyMsg{Type: tea.KeyLeft}) // collapse from child
	tm3 := newM3.(treeModel)

	if len(tm3.visibleRows()) != 2 {
		t.Fatalf("expected collapse to hide children again, got %d rows", len(tm3.visibleRows()))
	}
	if tm3.cursor != 0 {
		t.Fatalf("cursor = %d, want 0 (back on the pattern category row)", tm3.cursor)
	}
}

func TestTreeModel_EnterOnExerciseSelectsAndQuits(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight}) // expand pattern
	tm := newM.(treeModel)
	newM2, _ := tm.Update(tea.KeyMsg{Type: tea.KeyDown}) // move to first child
	tm2 := newM2.(treeModel)

	newM3, cmd := tm2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter on an exercise row to return a quit command")
	}
	tm3 := newM3.(treeModel)
	if tm3.selected == nil || tm3.selected.Exercise.ID != "two-pointers-01" {
		t.Errorf("expected two-pointers-01 selected, got %+v", tm3.selected)
	}
}

func TestTreeModel_QRequestsBack(t *testing.T) {
	m := newTreeModel(treeFixture())
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	if !newM.(treeModel).back {
		t.Error("expected back=true")
	}
}

func TestTreeModel_DownStopsAtLastVisibleRow(t *testing.T) {
	m := newTreeModel(treeFixture())
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(treeModel)
	}
	if m.cursor != 1 { // only 2 rows while collapsed: indices 0,1
		t.Errorf("cursor = %d, want 1 (last category while collapsed)", m.cursor)
	}
}
