package main

import "testing"

func TestFindRedundantConnection_Triangle(t *testing.T) {
	edges := [][]int{{1, 2}, {1, 3}, {2, 3}}
	got := FindRedundantConnection(edges)
	want := []int{2, 3}
	if got[0] != want[0] || got[1] != want[1] {
		t.Errorf("FindRedundantConnection(%v) = %v, want %v", edges, got, want)
	}
}

func TestFindRedundantConnection_LaterCycle(t *testing.T) {
	edges := [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 4}, {1, 5}}
	got := FindRedundantConnection(edges)
	want := []int{1, 4}
	if got[0] != want[0] || got[1] != want[1] {
		t.Errorf("FindRedundantConnection(%v) = %v, want %v", edges, got, want)
	}
}

func TestFindRedundantConnection_MergingComponents(t *testing.T) {
	edges := [][]int{{1, 4}, {3, 4}, {1, 3}, {1, 2}, {4, 5}}
	got := FindRedundantConnection(edges)
	want := []int{1, 3}
	if got[0] != want[0] || got[1] != want[1] {
		t.Errorf("FindRedundantConnection(%v) = %v, want %v", edges, got, want)
	}
}
