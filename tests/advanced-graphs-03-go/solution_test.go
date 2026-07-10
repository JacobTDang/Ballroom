package main

import "testing"

func TestNetworkDelayTime_Classic(t *testing.T) {
	times := [][]int{{2, 1, 1}, {2, 3, 1}, {3, 4, 1}}
	if got := NetworkDelayTime(times, 4, 2); got != 2 {
		t.Errorf("NetworkDelayTime(%v, 4, 2) = %d, want 2", times, got)
	}
}

func TestNetworkDelayTime_SingleEdgeReachable(t *testing.T) {
	times := [][]int{{1, 2, 1}}
	if got := NetworkDelayTime(times, 2, 1); got != 1 {
		t.Errorf("NetworkDelayTime(%v, 2, 1) = %d, want 1", times, got)
	}
}

func TestNetworkDelayTime_Unreachable(t *testing.T) {
	times := [][]int{{1, 2, 1}}
	if got := NetworkDelayTime(times, 2, 2); got != -1 {
		t.Errorf("NetworkDelayTime(%v, 2, 2) = %d, want -1", times, got)
	}
}
