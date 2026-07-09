package tui

import (
	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// fakeStatus builds a minimal ExerciseStatus for Update()-logic tests
// that don't touch the real exercise catalog or tracker DB. ProblemID
// defaults to ID, so each fixture call is its own standalone problem
// unless a test deliberately groups several under a shared ProblemID.
func fakeStatus(id string) catalog.ExerciseStatus {
	return catalog.ExerciseStatus{
		Exercise: exercise.Exercise{
			ID:        id,
			ProblemID: id,
			Title:     id,
			Category:  "pattern",
			Language:  "go",
		},
	}
}

func fakeStatusIn(category, id string) catalog.ExerciseStatus {
	s := fakeStatus(id)
	s.Exercise.Category = category
	return s
}

func treeFixture() []catalog.ExerciseStatus {
	return []catalog.ExerciseStatus{
		fakeStatusIn("pattern", "two-pointers-01"),
		fakeStatusIn("pattern", "two-pointers-01-cpp"),
		fakeStatusIn("debug", "off-by-one-01-go"),
		fakeStatusIn("debug", "off-by-one-01-cpp"),
	}
}
