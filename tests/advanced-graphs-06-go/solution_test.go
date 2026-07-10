package main

import "testing"

func TestFindCheapestPrice_OneStop(t *testing.T) {
	flights := [][]int{{0, 1, 100}, {1, 2, 100}, {2, 0, 100}, {1, 3, 600}, {2, 3, 200}}
	if got := FindCheapestPrice(4, flights, 0, 3, 1); got != 700 {
		t.Errorf("FindCheapestPrice(4, %v, 0, 3, 1) = %d, want 700", flights, got)
	}
}

func TestFindCheapestPrice_CheaperViaStop(t *testing.T) {
	flights := [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}}
	if got := FindCheapestPrice(3, flights, 0, 2, 1); got != 200 {
		t.Errorf("FindCheapestPrice(3, %v, 0, 2, 1) = %d, want 200", flights, got)
	}
}

func TestFindCheapestPrice_NoStopsAllowed(t *testing.T) {
	flights := [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}}
	if got := FindCheapestPrice(3, flights, 0, 2, 0); got != 500 {
		t.Errorf("FindCheapestPrice(3, %v, 0, 2, 0) = %d, want 500", flights, got)
	}
}
