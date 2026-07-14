package main

import "testing"

func TestExist(t *testing.T) {
	board := [][]byte{
		[]byte("ABCE"),
		[]byte("SFCS"),
		[]byte("ADEE"),
	}
	cases := []struct {
		word string
		want bool
	}{
		{"ABCCED", true},
		{"SEE", true},
		{"ABCB", false},
		{"ABFSAB", false},
	}

	for _, c := range cases {
		b := make([][]byte, len(board))
		for i, row := range board {
			b[i] = append([]byte(nil), row...)
		}
		got := Exist(b, c.word)
		if got != c.want {
			t.Errorf("Exist(board, %q) = %v, want %v", c.word, got, c.want)
		}
	}
}

func TestExist_SingleCell(t *testing.T) {
	board := [][]byte{[]byte("A")}
	if !Exist(board, "A") {
		t.Error("Exist single-cell board with matching word should be true")
	}
	board2 := [][]byte{[]byte("A")}
	if Exist(board2, "AA") {
		t.Error("Exist should be false when the word needs more cells than exist")
	}
}

func TestExist_DiagonalNotAllowed(t *testing.T) {
	board := [][]byte{[]byte("ab"), []byte("cd")}
	if !Exist(board, "abdc") {
		t.Error("Exist(board, \"abdc\") = false, want true")
	}
	if Exist(board, "abcd") {
		t.Error("Exist(board, \"abcd\") = true, want false (b and c are diagonal, not adjacent)")
	}
}
