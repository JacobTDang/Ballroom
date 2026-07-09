package main

import "testing"

func TestIsValidSudoku(t *testing.T) {
	validBoard := []string{
		"53..7....",
		"6..195...",
		".98....6.",
		"8...6...3",
		"4..8.3..1",
		"7...2...6",
		".6....28.",
		"...419..5",
		"....8..79",
	}
	invalidColumnBoard := []string{
		"83..7....",
		"6..195...",
		".98....6.",
		"8...6...3",
		"4..8.3..1",
		"7...2...6",
		".6....28.",
		"...419..5",
		"....8..79",
	}
	invalidRowBoard := []string{
		"5.......5",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
	}
	invalidBoxBoard := []string{
		"1........",
		".1.......",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
	}
	emptyBoard := []string{
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
		".........",
	}

	cases := []struct {
		name  string
		board []string
		want  bool
	}{
		{"valid", validBoard, true},
		{"invalid column", invalidColumnBoard, false},
		{"invalid row", invalidRowBoard, false},
		{"invalid box", invalidBoxBoard, false},
		{"empty", emptyBoard, true},
	}
	for _, c := range cases {
		got := IsValidSudoku(c.board)
		if got != c.want {
			t.Errorf("%s: IsValidSudoku(...) = %v, want %v", c.name, got, c.want)
		}
	}
}
