package main

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestNothingDroppedNothingDuplicated(t *testing.T) {
	for run := 0; run < 3; run++ {
		inputs := make([]int, 600)
		for i := range inputs {
			inputs[i] = i % 250 // duplicates on purpose
		}
		got := FanOutIn(inputs, 8, func(v int) int {
			time.Sleep(time.Duration(rand.Intn(2)) * time.Millisecond)
			return v * 3
		})

		if len(got) != len(inputs) {
			t.Fatalf("run %d: got %d results, want %d -- the fan-in dropped work", run, len(got), len(inputs))
		}
		want := make([]int, len(inputs))
		for i, v := range inputs {
			want[i] = v * 3
		}
		sort.Ints(want)
		sort.Ints(got)
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("run %d: result multiset differs at %d: got %d want %d", run, i, got[i], want[i])
			}
		}
	}
}
