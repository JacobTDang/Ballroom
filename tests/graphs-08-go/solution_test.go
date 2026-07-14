package main

import "testing"

func TestCanFinish_NoCycle(t *testing.T) {
	if !CanFinish(2, [][]int{{1, 0}}) {
		t.Errorf("CanFinish(2, [[1,0]]) = false, want true")
	}
}

func TestCanFinish_Cycle(t *testing.T) {
	if CanFinish(2, [][]int{{1, 0}, {0, 1}}) {
		t.Errorf("CanFinish(2, [[1,0],[0,1]]) = true, want false")
	}
}

func TestCanFinish_NoPrerequisites(t *testing.T) {
	if !CanFinish(5, nil) {
		t.Errorf("CanFinish(5, nil) = false, want true")
	}
}

func TestCanFinish_LongerCycle(t *testing.T) {
	prereqs := [][]int{{1, 0}, {2, 1}, {3, 2}, {0, 3}}
	if CanFinish(4, prereqs) {
		t.Errorf("CanFinish(4, %v) = true, want false", prereqs)
	}
}

func TestCanFinish_DiamondDAGNoCycle(t *testing.T) {
	prereqs := [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}
	if !CanFinish(4, prereqs) {
		t.Errorf("CanFinish(4, %v) = false, want true (diamond DAG, no cycle)", prereqs)
	}
}
