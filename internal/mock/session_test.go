package mock

import (
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func TestRemainingMinutes(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	cases := []struct {
		deadline time.Time
		want     int
	}{
		{now.Add(70 * time.Minute), 70},
		{now.Add(30*time.Minute + time.Second), 31}, // ceiling
		{now.Add(30 * time.Second), 1},
		{now, 0},
		{now.Add(-5 * time.Minute), 0},
	}
	for _, c := range cases {
		if got := RemainingMinutes(c.deadline, now); got != c.want {
			t.Errorf("RemainingMinutes(%v) = %d, want %d", c.deadline.Sub(now), got, c.want)
		}
	}
}

func TestOutcomeSince(t *testing.T) {
	attempts := []tracker.Attempt{
		{ID: 5, Result: tracker.ResultPass}, // before baseline — ignored
		{ID: 11, Result: tracker.ResultFail},
		{ID: 12, Result: tracker.ResultPass},
	}
	if got := OutcomeSince(attempts, 10); got != OutcomeSolved {
		t.Errorf("solved case = %q", got)
	}
	if got := OutcomeSince(attempts[:2], 10); got != OutcomeAttempted {
		t.Errorf("attempted case = %q", got)
	}
	if got := OutcomeSince(attempts, 12); got != OutcomeSkipped {
		t.Errorf("skipped case = %q", got)
	}
	if got := OutcomeSince(nil, 0); got != OutcomeSkipped {
		t.Errorf("no attempts = %q", got)
	}
}

func TestMaxAttemptID(t *testing.T) {
	if got := MaxAttemptID(nil); got != 0 {
		t.Errorf("empty = %d", got)
	}
	if got := MaxAttemptID([]tracker.Attempt{{ID: 3}, {ID: 9}, {ID: 7}}); got != 9 {
		t.Errorf("max = %d, want 9", got)
	}
}
