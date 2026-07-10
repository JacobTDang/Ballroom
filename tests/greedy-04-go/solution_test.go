package main

import "testing"

func TestCanCompleteCircuit_Classic(t *testing.T) {
	gas := []int{1, 2, 3, 4, 5}
	cost := []int{3, 4, 5, 1, 2}
	if got := CanCompleteCircuit(gas, cost); got != 3 {
		t.Errorf("CanCompleteCircuit(%v, %v) = %d, want 3", gas, cost, got)
	}
}

func TestCanCompleteCircuit_Impossible(t *testing.T) {
	gas := []int{2, 3, 4}
	cost := []int{3, 4, 3}
	if got := CanCompleteCircuit(gas, cost); got != -1 {
		t.Errorf("CanCompleteCircuit(%v, %v) = %d, want -1", gas, cost, got)
	}
}

func TestCanCompleteCircuit_SingleExact(t *testing.T) {
	gas := []int{5}
	cost := []int{4}
	if got := CanCompleteCircuit(gas, cost); got != 0 {
		t.Errorf("CanCompleteCircuit(%v, %v) = %d, want 0", gas, cost, got)
	}
}

func TestCanCompleteCircuit_SingleInsufficient(t *testing.T) {
	gas := []int{3}
	cost := []int{4}
	if got := CanCompleteCircuit(gas, cost); got != -1 {
		t.Errorf("CanCompleteCircuit(%v, %v) = %d, want -1", gas, cost, got)
	}
}

func TestCanCompleteCircuit_StartAtLastIndex(t *testing.T) {
	gas := []int{5, 1, 2, 3, 4}
	cost := []int{4, 4, 1, 5, 1}
	if got := CanCompleteCircuit(gas, cost); got != 4 {
		t.Errorf("CanCompleteCircuit(%v, %v) = %d, want 4", gas, cost, got)
	}
}
