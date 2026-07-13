package tutor

import (
	"fmt"
	"sync"
	"testing"
)

func TestActivityFeed_StartedAddsARunningLine(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "read_solution_file", "")
	if len(lines) != 1 || lines[0] != "● read_solution_file" {
		t.Errorf("lines = %v, want [\"● read_solution_file\"]", lines)
	}
}

func TestActivityFeed_StartedWithArgsShowsThemInParens(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "highlight_lines", `{"start_line":10,"end_line":20}`)
	if len(lines) != 1 || lines[0] != `● highlight_lines({"start_line":10,"end_line":20})` {
		t.Errorf("lines = %v, want the args shown in parens", lines)
	}
}

func TestActivityFeed_StartedWithEmptyOrNoArgsOmitsParens(t *testing.T) {
	f := &activityFeed{}
	lines := f.started("call-1", "read_solution_file", "{}")
	if len(lines) != 1 || lines[0] != "● read_solution_file" {
		t.Errorf("lines = %v, want no parens for empty/no-op args", lines)
	}
}

func TestActivityFeed_FinishedUpdatesTheMatchingCallToDone(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.finished("call-1", "312 bytes")
	if len(lines) != 1 || lines[0] != "● read_solution_file  312 bytes" {
		t.Errorf("lines = %v, want the call marked done with its result", lines)
	}
}

func TestActivityFeed_FinishedWithEmptyResultOmitsTrailingSpace(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "highlight_lines", "")
	lines := f.finished("call-1", "")
	if len(lines) != 1 || lines[0] != "● highlight_lines" {
		t.Errorf("lines = %v, want just the dot and name, no trailing separator", lines)
	}
}

func TestActivityFeed_FailedUpdatesTheMatchingCallToFailed(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_test_output", "")
	lines := f.failed("call-1", "no test run yet")
	if len(lines) != 1 || lines[0] != "● read_test_output - failed: no test run yet" {
		t.Errorf("lines = %v, want the call marked failed with the error", lines)
	}
}

func TestActivityFeed_FinishedForUnknownCallIDIsANoOp(t *testing.T) {
	// A callID that was never started (or already dropped by the cap
	// below) must not panic or fabricate a new entry -- eino's own
	// OnEnd/OnError always follow a real OnStart for the same call, but
	// this call may have aged out of the capped list already.
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.finished("call-unknown", "some result")
	if len(lines) != 1 || lines[0] != "● read_solution_file" {
		t.Errorf("lines = %v, want the existing call untouched and no new entry added", lines)
	}
}

func TestActivityFeed_MultipleCallsPreserveStartOrder(t *testing.T) {
	f := &activityFeed{}
	f.started("call-1", "read_solution_file", "")
	lines := f.started("call-2", "read_problem_statement", "")
	if len(lines) != 2 || lines[0] != "● read_solution_file" || lines[1] != "● read_problem_statement" {
		t.Errorf("lines = %v, want both calls in start order", lines)
	}
}

func TestActivityFeed_CapsAtFourDroppingTheOldest(t *testing.T) {
	f := &activityFeed{}
	for i := 1; i <= 5; i++ {
		f.started(fmt.Sprintf("call-%d", i), fmt.Sprintf("tool_%d", i), "")
	}
	lines := f.started("call-6", "tool_6", "")
	if len(lines) != activityToolLines {
		t.Fatalf("len(lines) = %d, want %d (the cap)", len(lines), activityToolLines)
	}
	if lines[0] != "● tool_3" {
		t.Errorf("lines[0] = %q, want the oldest (tool_1, tool_2) dropped, starting at tool_3", lines[0])
	}
	if lines[len(lines)-1] != "● tool_6" {
		t.Errorf("lines[last] = %q, want the newest call last", lines[len(lines)-1])
	}
}

func TestActivityFeed_ConcurrentStartedFinishedDoNotRace(t *testing.T) {
	f := &activityFeed{}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("call-%d", n)
			name := fmt.Sprintf("tool_%d", n)
			f.started(id, name, "")
			f.finished(id, "done")
		}(i)
	}
	wg.Wait()
}
