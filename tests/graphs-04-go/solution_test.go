package main

import (
	"reflect"
	"testing"
)

func TestWallsAndGates(t *testing.T) {
	rooms := [][]int{
		{Inf, -1, 0, Inf},
		{Inf, Inf, Inf, -1},
		{Inf, -1, Inf, -1},
		{0, -1, Inf, Inf},
	}
	want := [][]int{
		{3, -1, 0, 1},
		{2, 2, 1, -1},
		{1, -1, 2, -1},
		{0, -1, 3, 4},
	}

	WallsAndGates(rooms)
	if !reflect.DeepEqual(rooms, want) {
		t.Errorf("WallsAndGates(...) -> %v, want %v", rooms, want)
	}
}

func TestWallsAndGates_UnreachableRoomStaysInf(t *testing.T) {
	rooms := [][]int{
		{0, -1, Inf},
	}
	want := [][]int{
		{0, -1, Inf},
	}
	WallsAndGates(rooms)
	if !reflect.DeepEqual(rooms, want) {
		t.Errorf("WallsAndGates(...) -> %v, want %v (unreachable room stays Inf)", rooms, want)
	}
}

func TestWallsAndGates_NoGates(t *testing.T) {
	rooms := [][]int{{Inf, Inf}}
	want := [][]int{{Inf, Inf}}
	WallsAndGates(rooms)
	if !reflect.DeepEqual(rooms, want) {
		t.Errorf("WallsAndGates(...) -> %v, want %v", rooms, want)
	}
}

func TestWallsAndGates_NearestGateWins(t *testing.T) {
	rooms := [][]int{{0, Inf, Inf, 0}}
	want := [][]int{{0, 1, 1, 0}}
	WallsAndGates(rooms)
	if !reflect.DeepEqual(rooms, want) {
		t.Errorf("WallsAndGates(...) -> %v, want %v (equidistant from both gates)", rooms, want)
	}
}
