package main

import "testing"

func TestBoundaryBurstIsCaught(t *testing.T) {
	// The fixed-window flaw: 2 allowed at t=90,95 and 2 more at
	// t=105,110 puts 4 through the span (10,110] -- a sliding window
	// must deny at 105 and 110.
	s := NewSlidingWindow(2, 100)
	for _, tc := range []struct {
		at   int64
		want bool
	}{{90, true}, {95, true}, {105, false}, {110, false}, {191, true}} {
		if got := s.AllowAt(tc.at); got != tc.want {
			t.Fatalf("AllowAt(%d) = %v, want %v", tc.at, got, tc.want)
		}
	}
}

func TestRequestsAgeOutExactly(t *testing.T) {
	s := NewSlidingWindow(1, 100)
	if !s.AllowAt(1000) {
		t.Fatal("first request denied")
	}
	if s.AllowAt(1099) {
		t.Fatal("AllowAt(1099): the 1000 request is 99ms old, still in the window")
	}
	if !s.AllowAt(1100) {
		t.Fatal("AllowAt(1100): the 1000 request is exactly a window old and must no longer count")
	}
}

func TestDeniedRequestsDoNotCount(t *testing.T) {
	s := NewSlidingWindow(2, 100)
	s.AllowAt(0)
	s.AllowAt(1)
	for i := int64(2); i < 50; i++ {
		if s.AllowAt(i) {
			t.Fatalf("AllowAt(%d) allowed a third request inside the window", i)
		}
	}
	// If the denials above were (wrongly) recorded, this would still
	// see a full window and deny.
	if !s.AllowAt(101) {
		t.Fatal("AllowAt(101) denied -- denied requests must not consume the budget")
	}
}

func TestSteadyRateUnderLimitAlwaysPasses(t *testing.T) {
	s := NewSlidingWindow(2, 100)
	for at := int64(0); at < 1000; at += 60 {
		if !s.AllowAt(at) {
			t.Fatalf("AllowAt(%d) denied a steady rate well under the limit", at)
		}
	}
}
