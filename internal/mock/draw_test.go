package mock

import (
	"testing"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func status(problemID, lastResult, lastDate string, attempts int) catalog.ExerciseStatus {
	return catalog.ExerciseStatus{
		Exercise: exercise.Exercise{
			ID:        problemID + "-python",
			ProblemID: problemID,
			Category:  exercise.CategoryCapitalOne,
			Language:  exercise.LanguagePython,
		},
		Attempts:        attempts,
		LastResult:      lastResult,
		LastAttemptDate: lastDate,
	}
}

func TestSlotOf(t *testing.T) {
	cases := map[string]string{
		"c1-warmup-03": SlotWarmup,
		"c1-grid-01":   SlotGrid,
		"c1-algo-06":   SlotAlgo,
		"two-sum-01":   "",
	}
	for id, want := range cases {
		if got := SlotOf(id); got != want {
			t.Errorf("SlotOf(%q) = %q, want %q", id, got, want)
		}
	}
}

func TestDraw_FillsSlotsPreferringUnsolved(t *testing.T) {
	statuses := []catalog.ExerciseStatus{
		status("c1-warmup-01", tracker.ResultPass, "2026-08-01", 1),
		status("c1-warmup-02", "", "", 0),
		status("c1-warmup-03", "", "", 0),
		status("c1-grid-01", tracker.ResultPass, "2026-08-02", 2),
		status("c1-grid-02", tracker.ResultFail, "2026-08-03", 1),
		status("c1-algo-01", "", "", 0),
	}
	got, err := Draw(statuses)
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	// Slots 0,1: the two unsolved warmups, not the solved one.
	if got[0].ProblemID != "c1-warmup-02" || got[1].ProblemID != "c1-warmup-03" {
		t.Errorf("warmup slots = %q,%q; want c1-warmup-02,c1-warmup-03", got[0].ProblemID, got[1].ProblemID)
	}
	// Slot 2: grid-02 (failed = still unsolved) beats grid-01 (solved).
	if got[2].ProblemID != "c1-grid-02" {
		t.Errorf("grid slot = %q, want c1-grid-02", got[2].ProblemID)
	}
	if got[3].ProblemID != "c1-algo-01" {
		t.Errorf("algo slot = %q, want c1-algo-01", got[3].ProblemID)
	}
}

func TestDraw_ExhaustedPoolFallsBackToLeastRecent(t *testing.T) {
	statuses := []catalog.ExerciseStatus{
		status("c1-warmup-01", tracker.ResultPass, "2026-08-05", 1),
		status("c1-warmup-02", tracker.ResultPass, "2026-08-01", 1),
		status("c1-warmup-03", tracker.ResultPass, "2026-08-03", 1),
		status("c1-grid-01", tracker.ResultPass, "2026-08-02", 1),
		status("c1-algo-01", tracker.ResultPass, "2026-08-02", 1),
	}
	got, err := Draw(statuses)
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	// All solved: least-recently-attempted warmups first (08-01, then 08-03).
	if got[0].ProblemID != "c1-warmup-02" || got[1].ProblemID != "c1-warmup-03" {
		t.Errorf("warmup slots = %q,%q; want c1-warmup-02,c1-warmup-03", got[0].ProblemID, got[1].ProblemID)
	}
}

func TestDraw_EmptyPoolErrors(t *testing.T) {
	if _, err := Draw(nil); err == nil {
		t.Fatal("Draw with no capital-one exercises should error")
	}
}

func TestDraw_IgnoresOtherCategoriesAndLanguages(t *testing.T) {
	other := status("c1-warmup-01", "", "", 0)
	other.Exercise.Category = exercise.CategoryDSA
	statuses := []catalog.ExerciseStatus{other}
	if _, err := Draw(statuses); err == nil {
		t.Fatal("non-capital-one exercises must not enter the pools")
	}
}
