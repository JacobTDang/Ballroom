// Package mock implements the Capital One timed mock sitting: the
// 4-question draw from the capital-one slot pools, the sitting log,
// and the shared-clock arithmetic. All host-side — the per-question
// containers are ordinary exercise sessions.
package mock

import (
	"fmt"
	"sort"
	"strings"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// The GCA's question shape (research-verified): two easy warm-ups, an
// implementation-heavy grid problem, an algorithmic medium.
const (
	SlotWarmup = "warmup"
	SlotGrid   = "grid"
	SlotAlgo   = "algo"
)

// Slots is a sitting's question sequence — index i is question i+1.
var Slots = [4]string{SlotWarmup, SlotWarmup, SlotGrid, SlotAlgo}

// SlotOf maps a capital-one problem id to its slot pool by prefix —
// the id convention (c1-warmup-*, c1-grid-*, c1-algo-*) is the slot
// encoding, deliberately not a schema field.
func SlotOf(problemID string) string {
	for _, slot := range []string{SlotWarmup, SlotGrid, SlotAlgo} {
		if strings.HasPrefix(problemID, "c1-"+slot+"-") {
			return slot
		}
	}
	return ""
}

// Draw fills Slots from the capital-one pools: unsolved problems first,
// then least-recently-attempted (never-attempted sorts as oldest), id as
// the deterministic tiebreak. No problem repeats within a sitting.
func Draw(statuses []catalog.ExerciseStatus) ([4]exercise.Exercise, error) {
	pools := make(map[string][]catalog.ExerciseStatus)
	for _, s := range statuses {
		if s.Exercise.Category != exercise.CategoryCapitalOne ||
			s.Exercise.Language != exercise.LanguagePython {
			continue
		}
		if slot := SlotOf(s.Exercise.ProblemID); slot != "" {
			pools[slot] = append(pools[slot], s)
		}
	}

	for _, pool := range pools {
		sort.Slice(pool, func(i, j int) bool {
			si, sj := pool[i].LastResult == tracker.ResultPass, pool[j].LastResult == tracker.ResultPass
			if si != sj {
				return !si // unsolved first
			}
			if pool[i].LastAttemptDate != pool[j].LastAttemptDate {
				return pool[i].LastAttemptDate < pool[j].LastAttemptDate
			}
			return pool[i].Exercise.ID < pool[j].Exercise.ID
		})
	}

	var out [4]exercise.Exercise
	used := make(map[string]int) // slot -> next index into its sorted pool
	for i, slot := range Slots {
		pool := pools[slot]
		if used[slot] >= len(pool) {
			return out, fmt.Errorf("mock: %s pool has %d problems, need %d", slot, len(pool), used[slot]+1)
		}
		out[i] = pool[used[slot]].Exercise
		used[slot]++
	}
	return out, nil
}
