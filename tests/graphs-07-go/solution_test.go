package main

import (
	"reflect"
	"testing"
)

func toGrid(rows []string) [][]byte {
	out := make([][]byte, len(rows))
	for i, row := range rows {
		out[i] = []byte(row)
	}
	return out
}

func TestSolve_Classic(t *testing.T) {
	board := toGrid([]string{
		"XXXX",
		"XOOX",
		"XXOX",
		"XOXX",
	})
	want := toGrid([]string{
		"XXXX",
		"XXXX",
		"XXXX",
		"XOXX",
	})
	Solve(board)
	if !reflect.DeepEqual(board, want) {
		t.Errorf("Solve(...) = %v, want %v", board, want)
	}
}

func TestSolve_AllBorderConnected(t *testing.T) {
	board := toGrid([]string{
		"OOO",
		"OXO",
		"OOO",
	})
	want := toGrid([]string{
		"OOO",
		"OXO",
		"OOO",
	})
	Solve(board)
	if !reflect.DeepEqual(board, want) {
		t.Errorf("Solve(...) = %v, want %v", board, want)
	}
}

func TestSolve_SingleCell(t *testing.T) {
	board := toGrid([]string{"O"})
	want := toGrid([]string{"O"})
	Solve(board)
	if !reflect.DeepEqual(board, want) {
		t.Errorf("Solve(...) = %v, want %v", board, want)
	}
}

func TestSolve_MixedSurroundedAndBorderConnected(t *testing.T) {
	board := toGrid([]string{
		"XXXXX",
		"XOOXX",
		"XOXXX",
		"XXXOX",
		"XXOOX",
	})
	want := toGrid([]string{
		"XXXXX",
		"XXXXX",
		"XXXXX",
		"XXXOX",
		"XXOOX",
	})
	Solve(board)
	if !reflect.DeepEqual(board, want) {
		t.Errorf("Solve(...) = %v, want %v", board, want)
	}
}
