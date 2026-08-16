package mock

import (
	"time"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// TotalMinutes is the sitting budget — the real GCA's 70 minutes.
const TotalMinutes = 70

// RemainingMinutes is the whole-minute budget left before deadline,
// rounded up so a container is never launched with a zero limit while
// time genuinely remains. 0 means the sitting is over. Wall clock, not
// container uptime: between-question time (and a lid close) burns mock
// time on purpose — the real assessment's clock doesn't pause either.
func RemainingMinutes(deadline, now time.Time) int {
	left := deadline.Sub(now)
	if left <= 0 {
		return 0
	}
	return int((left + time.Minute - 1) / time.Minute)
}

// MaxAttemptID returns the largest attempt ID, 0 for none — the
// per-question baseline captured before its container launches.
func MaxAttemptID(attempts []tracker.Attempt) int64 {
	var max int64
	for _, a := range attempts {
		if a.ID > max {
			max = a.ID
		}
	}
	return max
}

// OutcomeSince classifies a question by the attempts its container
// logged: rows with ID > baselineID are this question's. Attempt.Date
// is day-granular, so the autoincrement ID is the only reliable
// "since the question started" marker.
func OutcomeSince(attempts []tracker.Attempt, baselineID int64) string {
	outcome := OutcomeSkipped
	for _, a := range attempts {
		if a.ID <= baselineID {
			continue
		}
		if a.Result == tracker.ResultPass {
			return OutcomeSolved
		}
		outcome = OutcomeAttempted
	}
	return outcome
}
