package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/mock"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func mockFixturePlan() [4]exercise.Exercise {
	mk := func(pid string, limit int) exercise.Exercise {
		return exercise.Exercise{ID: pid + "-python", ProblemID: pid, Category: exercise.CategoryCapitalOne,
			Language: exercise.LanguagePython, TimeLimitMin: limit}
	}
	return [4]exercise.Exercise{mk("c1-warmup-01", 10), mk("c1-warmup-02", 10), mk("c1-grid-01", 25), mk("c1-algo-01", 20)}
}

// scriptedClock advances by step on every call — deterministic time for
// the run loop without sleeping.
func scriptedClock(start time.Time, step time.Duration) func() time.Time {
	t := start
	return func() time.Time {
		t = t.Add(step)
		return t
	}
}

func TestRunMockSitting_AllSolved(t *testing.T) {
	var launched []string
	var limits []int
	attemptID := int64(0)
	attempts := map[string][]tracker.Attempt{}
	runFake := func(_ config.Config, ex exercise.Exercise, draftDir string) error {
		launched = append(launched, ex.ProblemID)
		limits = append(limits, ex.TimeLimitMin)
		if draftDir != "" {
			t.Errorf("mock question launched with draft dir %q, want fresh start", draftDir)
		}
		attemptID++
		attempts[ex.ID] = append(attempts[ex.ID], tracker.Attempt{ID: attemptID, Result: tracker.ResultPass})
		return nil
	}
	attemptsFake := func(id string) ([]tracker.Attempt, error) { return attempts[id], nil }
	promptFake := func(q, remaining int) mockChoice { return mockContinue }

	s, err := runMockSitting(config.Config{}, mockFixturePlan(),
		scriptedClock(time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC), time.Minute),
		runFake, attemptsFake, promptFake)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(launched) != 4 {
		t.Fatalf("launched %d questions, want 4", len(launched))
	}
	if s.Solved() != 4 || !s.Completed {
		t.Errorf("sitting = %+v, want 4 solved, completed", s)
	}
	// Every launch's TimeLimitMin must be the remaining budget, not the
	// exercise's own limit — first launch close to 70, later ones smaller.
	if limits[0] > mock.TotalMinutes || limits[0] < 60 {
		t.Errorf("first limit = %d, want ~70", limits[0])
	}
	if !(limits[0] > limits[1] && limits[1] > limits[2] && limits[2] > limits[3]) {
		t.Errorf("limits not strictly decreasing: %v", limits)
	}
}

func TestRunMockSitting_DeadlineStopsLaunches(t *testing.T) {
	launched := 0
	runFake := func(_ config.Config, ex exercise.Exercise, _ string) error { launched++; return nil }
	attemptsFake := func(string) ([]tracker.Attempt, error) { return nil, nil }
	promptFake := func(q, remaining int) mockChoice { return mockContinue }

	// 40-minute steps: the 70-minute budget is gone after the first
	// couple of clock reads.
	s, _ := runMockSitting(config.Config{}, mockFixturePlan(),
		scriptedClock(time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC), 40*time.Minute),
		runFake, attemptsFake, promptFake)

	if launched >= 4 {
		t.Errorf("launched %d, want fewer — deadline must stop the loop", launched)
	}
	if s.Completed {
		t.Error("a timed-out sitting is not Completed")
	}
	for i := launched; i < 4; i++ {
		if s.Outcomes[i] != mock.OutcomeSkipped {
			t.Errorf("slot %d = %q, want skipped", i, s.Outcomes[i])
		}
	}
}

func TestRunMockSitting_AbortRecordsRestSkipped(t *testing.T) {
	runFake := func(_ config.Config, ex exercise.Exercise, _ string) error { return nil }
	attemptsFake := func(string) ([]tracker.Attempt, error) { return nil, nil }
	promptFake := func(q, remaining int) mockChoice {
		if q == 2 {
			return mockAbort
		}
		return mockContinue
	}
	s, _ := runMockSitting(config.Config{}, mockFixturePlan(),
		scriptedClock(time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC), time.Minute),
		runFake, attemptsFake, promptFake)
	if s.Completed {
		t.Error("aborted sitting must not be Completed")
	}
	if s.Outcomes[2] != mock.OutcomeSkipped || s.Outcomes[3] != mock.OutcomeSkipped {
		t.Errorf("post-abort slots = %v, want skipped", s.Outcomes)
	}
}

func TestRunMockSitting_SkipStillCompletes(t *testing.T) {
	runFake := func(_ config.Config, ex exercise.Exercise, _ string) error { return nil }
	attemptsFake := func(string) ([]tracker.Attempt, error) { return nil, nil }
	promptFake := func(q, remaining int) mockChoice {
		if q == 3 {
			return mockSkip
		}
		return mockContinue
	}
	s, _ := runMockSitting(config.Config{}, mockFixturePlan(),
		scriptedClock(time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC), time.Minute),
		runFake, attemptsFake, promptFake)
	if !s.Completed {
		t.Error("a voluntary skip must still complete the sitting")
	}
	if s.Outcomes[2] != mock.OutcomeSkipped {
		t.Errorf("skipped slot = %q, want skipped", s.Outcomes[2])
	}
}

func TestRunMockSitting_LaunchErrorSurfaces(t *testing.T) {
	boom := fmt.Errorf("docker went away")
	calls := 0
	runFake := func(_ config.Config, ex exercise.Exercise, _ string) error {
		calls++
		if calls == 2 {
			return boom
		}
		return nil
	}
	attemptsFake := func(string) ([]tracker.Attempt, error) { return nil, nil }
	promptFake := func(q, remaining int) mockChoice { return mockContinue }

	s, err := runMockSitting(config.Config{}, mockFixturePlan(),
		scriptedClock(time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC), time.Minute),
		runFake, attemptsFake, promptFake)
	if err == nil || !strings.Contains(err.Error(), "docker went away") {
		t.Fatalf("launch error must be returned, got %v", err)
	}
	if s.Completed {
		t.Error("an error-aborted sitting must not be Completed")
	}
	if calls != 2 {
		t.Errorf("launches after the error: got %d calls, want 2", calls)
	}
}

func TestMockSummaryView_ShowsError(t *testing.T) {
	m := appModel{cfg: config.Config{}, stage: stageMockSummary, err: fmt.Errorf("docker went away")}
	m.mockSitting = &mock.Sitting{StartedAt: "2026-08-16T10:00:00Z"}
	got := stripAnsiTUI(m.renderMockSummary())
	if !strings.Contains(got, "docker went away") {
		t.Errorf("summary must surface the session error:\n%s", got)
	}
}

func TestMockStartView_TruncatesLongTitles(t *testing.T) {
	m := appModel{cfg: config.Config{}, stage: stageMockStart}
	m.mockPlan = mockFixturePlan()
	m.mockPlan[0].Title = strings.Repeat("Very Long Banking Problem Title ", 4)
	got := stripAnsiTUI(m.renderMockStart())
	for _, line := range strings.Split(got, "\n") {
		if len([]rune(line)) > 76 {
			t.Errorf("start-screen line overflows the panel (%d runes): %q", len([]rune(line)), line)
		}
	}
}

func TestMockSummaryView(t *testing.T) {
	m := appModel{cfg: config.Config{}, stage: stageMockSummary}
	m.mockSitting = &mock.Sitting{
		StartedAt:      "2026-08-16T10:00:00Z",
		ProblemIDs:     [4]string{"c1-warmup-01", "c1-warmup-02", "c1-grid-01", "c1-algo-01"},
		Outcomes:       [4]string{mock.OutcomeSolved, mock.OutcomeSolved, mock.OutcomeAttempted, mock.OutcomeSkipped},
		MinutesPerSlot: [4]float64{6.5, 8.2, 30.1, 0},
		MinutesTotal:   62.3,
		Completed:      true,
	}
	got := stripAnsiTUI(m.renderMockSummary())
	for _, want := range []string{"Solved 2/4", "c1-warmup-01", "attempted", "skipped", "62"} {
		if !strings.Contains(got, want) {
			t.Errorf("summary missing %q:\n%s", want, got)
		}
	}
}

func TestMockStartView_ShowsDrawAndHistory(t *testing.T) {
	m := appModel{cfg: config.Config{}, stage: stageMockStart}
	m.mockPlan = mockFixturePlan()
	m.mockHistory = "2 sittings · best 3/4 · last 2/4"
	got := stripAnsiTUI(m.renderMockStart())
	for _, want := range []string{"70", "c1-warmup-01", "c1-grid-01", "c1-algo-01", "2 sittings", "forward-only"} {
		if !strings.Contains(got, want) {
			t.Errorf("start screen missing %q:\n%s", want, got)
		}
	}
}
