package main

import "testing"

func TestCountComponents_Classic(t *testing.T) {
	edges := [][]int{{0, 1}, {1, 2}, {3, 4}}
	if got := CountComponents(5, edges); got != 2 {
		t.Errorf("CountComponents(5, %v) = %d, want 2", edges, got)
	}
}

func TestCountComponents_AllConnected(t *testing.T) {
	edges := [][]int{{0, 1}, {1, 2}, {2, 3}}
	if got := CountComponents(4, edges); got != 1 {
		t.Errorf("CountComponents(4, %v) = %d, want 1", edges, got)
	}
}

func TestCountComponents_NoEdges(t *testing.T) {
	if got := CountComponents(4, nil); got != 4 {
		t.Errorf("CountComponents(4, nil) = %d, want 4", got)
	}
}

func TestCountComponents_SingleNode(t *testing.T) {
	if got := CountComponents(1, nil); got != 1 {
		t.Errorf("CountComponents(1, nil) = %d, want 1", got)
	}
}
