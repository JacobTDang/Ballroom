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

	solvedBoard := []string{
		"534678912",
		"672195348",
		"198342567",
		"859761423",
		"426853791",
		"713924856",
		"961537284",
		"287419635",
		"345286179",
	}
	sameDigitDifferentUnitsBoard := []string{
		"5........",
		".........",
		".........",
		".........",
		"....5....",
		".........",
		".........",
		".........",
		".........",
	}
	singleCellBoard := []string{
		"5........",
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
		{"fully solved valid board", solvedBoard, true},
		{"same digit different row/col/box", sameDigitDifferentUnitsBoard, true},
		{"single cell filled", singleCellBoard, true},
	}
	for _, c := range cases {
		got := IsValidSudoku(c.board)
		if got != c.want {
			t.Errorf("%s: IsValidSudoku(...) = %v, want %v", c.name, got, c.want)
		}
	}
}
