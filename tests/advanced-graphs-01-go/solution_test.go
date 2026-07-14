package main

import (
	"reflect"
	"testing"
)

func TestFindItinerary_DeadEnd(t *testing.T) {
	tickets := [][]string{
		{"JFK", "SFO"},
		{"JFK", "ATL"},
		{"SFO", "ATL"},
		{"ATL", "JFK"},
		{"ATL", "SFO"},
	}
	want := []string{"JFK", "ATL", "JFK", "SFO", "ATL", "SFO"}
	if got := FindItinerary(tickets); !reflect.DeepEqual(got, want) {
		t.Errorf("FindItinerary(%v) = %v, want %v", tickets, got, want)
	}
}

func TestFindItinerary_Simple(t *testing.T) {
	tickets := [][]string{
		{"MUC", "LHR"},
		{"JFK", "MUC"},
		{"SFO", "SJC"},
		{"LHR", "SFO"},
	}
	want := []string{"JFK", "MUC", "LHR", "SFO", "SJC"}
	if got := FindItinerary(tickets); !reflect.DeepEqual(got, want) {
		t.Errorf("FindItinerary(%v) = %v, want %v", tickets, got, want)
	}
}

func TestFindItinerary_LexicalTieBreak(t *testing.T) {
	tickets := [][]string{
		{"JFK", "KUL"},
		{"JFK", "NRT"},
		{"NRT", "JFK"},
	}
	want := []string{"JFK", "NRT", "JFK", "KUL"}
	if got := FindItinerary(tickets); !reflect.DeepEqual(got, want) {
		t.Errorf("FindItinerary(%v) = %v, want %v", tickets, got, want)
	}
}

func TestFindItinerary_SimpleTwoCycle(t *testing.T) {
	tickets := [][]string{
		{"JFK", "A"},
		{"A", "JFK"},
	}
	want := []string{"JFK", "A", "JFK"}
	if got := FindItinerary(tickets); !reflect.DeepEqual(got, want) {
		t.Errorf("FindItinerary(%v) = %v, want %v", tickets, got, want)
	}
}

func TestFindItinerary_BranchingAtOrigin(t *testing.T) {
	tickets := [][]string{
		{"JFK", "B"},
		{"JFK", "A"},
		{"B", "JFK"},
		{"A", "JFK"},
	}
	want := []string{"JFK", "A", "JFK", "B", "JFK"}
	if got := FindItinerary(tickets); !reflect.DeepEqual(got, want) {
		t.Errorf("FindItinerary(%v) = %v, want %v", tickets, got, want)
	}
}
