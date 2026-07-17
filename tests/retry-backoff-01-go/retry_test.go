package main

import (
	"errors"
	"testing"
)

func recorder() (func(int64), *[]int64) {
	var slept []int64
	return func(ms int64) { slept = append(slept, ms) }, &slept
}

func failNTimes(n int) (func() error, *int) {
	calls := 0
	return func() error {
		calls++
		if calls <= n {
			return errors.New("boom")
		}
		return nil
	}, &calls
}

func TestFirstTrySuccessNeverSleeps(t *testing.T) {
	sleep, slept := recorder()
	op, calls := failNTimes(0)
	if err := Retry(op, 5, 100, 1000, sleep); err != nil {
		t.Fatalf("Retry = %v, want nil", err)
	}
	if *calls != 1 || len(*slept) != 0 {
		t.Fatalf("calls=%d sleeps=%v, want 1 call and no sleeps", *calls, *slept)
	}
}

func TestExponentialDelaysExact(t *testing.T) {
	sleep, slept := recorder()
	op, calls := failNTimes(3)
	if err := Retry(op, 5, 100, 10000, sleep); err != nil {
		t.Fatalf("Retry = %v, want eventual success", err)
	}
	if *calls != 4 {
		t.Fatalf("calls = %d, want 4 (3 failures + success)", *calls)
	}
	want := []int64{100, 200, 400}
	if len(*slept) != len(want) {
		t.Fatalf("sleeps = %v, want %v", *slept, want)
	}
	for i, w := range want {
		if (*slept)[i] != w {
			t.Fatalf("sleep %d = %d, want %d (base*2^i)", i, (*slept)[i], w)
		}
	}
}

func TestCapFlattensTheCurve(t *testing.T) {
	sleep, slept := recorder()
	op, _ := failNTimes(4)
	if err := Retry(op, 6, 100, 250, sleep); err != nil {
		t.Fatalf("Retry = %v, want eventual success", err)
	}
	want := []int64{100, 200, 250, 250}
	if len(*slept) != len(want) {
		t.Fatalf("sleeps = %v, want %v", *slept, want)
	}
	for i, w := range want {
		if (*slept)[i] != w {
			t.Fatalf("sleep %d = %d, want %d (capped)", i, (*slept)[i], w)
		}
	}
}

func TestExhaustionReturnsLastErrorNoTrailingSleep(t *testing.T) {
	sleep, slept := recorder()
	calls := 0
	op := func() error {
		calls++
		return errors.New("always")
	}
	err := Retry(op, 3, 100, 1000, sleep)
	if err == nil || err.Error() != "always" {
		t.Fatalf("Retry = %v, want the operation's own last error", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want exactly maxAttempts 3", calls)
	}
	if len(*slept) != 2 {
		t.Fatalf("sleeps = %v, want 2 (never sleep after the final failure)", *slept)
	}
}
