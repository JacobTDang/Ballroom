package main

import "testing"

func TestValidTree_Valid(t *testing.T) {
	edges := [][]int{{0, 1}, {0, 2}, {0, 3}, {1, 4}}
	if !ValidTree(5, edges) {
		t.Errorf("ValidTree(5, %v) = false, want true", edges)
	}
}

func TestValidTree_HasCycle(t *testing.T) {
	edges := [][]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}}
	if ValidTree(5, edges) {
		t.Errorf("ValidTree(5, %v) = true, want false", edges)
	}
}

func TestValidTree_Disconnected(t *testing.T) {
	edges := [][]int{{0, 1}, {2, 3}}
	if ValidTree(4, edges) {
		t.Errorf("ValidTree(4, %v) = true, want false", edges)
	}
}

func TestValidTree_SingleNode(t *testing.T) {
	if !ValidTree(1, nil) {
		t.Errorf("ValidTree(1, nil) = false, want true")
	}
}

func TestValidTree_NoEdgesMultipleNodes(t *testing.T) {
	if ValidTree(3, nil) {
		t.Errorf("ValidTree(3, nil) = true, want false (disconnected)")
	}
}
