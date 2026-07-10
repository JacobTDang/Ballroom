package main

import "testing"

func TestKthLargest(t *testing.T) {
	kl := NewKthLargest(3, []int{4, 5, 8, 2})
	cases := []struct {
		val  int
		want int
	}{
		{3, 4},
		{5, 5},
		{10, 5},
		{9, 8},
		{4, 8},
	}
	for _, c := range cases {
		if got := kl.Add(c.val); got != c.want {
			t.Errorf("Add(%d) = %d, want %d", c.val, got, c.want)
		}
	}
}

func TestKthLargest_EmptyInitialStream(t *testing.T) {
	kl := NewKthLargest(1, []int{})
	if got := kl.Add(-3); got != -3 {
		t.Errorf("Add(-3) = %d, want -3", got)
	}
	if got := kl.Add(-2); got != -2 {
		t.Errorf("Add(-2) = %d, want -2", got)
	}
}
