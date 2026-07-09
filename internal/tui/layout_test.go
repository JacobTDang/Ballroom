package tui

import (
	"strings"
	"testing"
)

func runeIndex(s string, r rune) int {
	for i, c := range []rune(s) {
		if c == r {
			return i
		}
	}
	return -1
}

func TestPlaceBlock_PreservesRelativeAlignmentBetweenLines(t *testing.T) {
	// A narrow line directly under a wider one, both meant to line up on
	// the same column. If placeBlock centered each line independently
	// (lipgloss.Place's actual bug), this relationship would break.
	content := "1234567890\n     X"
	placed := placeBlock(30, 10, content)
	lines := strings.Split(strings.TrimLeft(placed, "\n"), "\n")

	line1Start := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
	line2XCol := runeIndex(lines[1], 'X')
	wantXCol := line1Start + 5 // 'X' sat at column 5 in the original content

	if line2XCol != wantXCol {
		t.Errorf("X at col %d, want %d — relative alignment between lines was not preserved", line2XCol, wantXCol)
	}
}

func TestPlaceBlock_CentersWithinViewport(t *testing.T) {
	// Only top-padding is added — a terminal's alt-screen is already
	// blank below the content, so there's nothing to gain from padding
	// the bottom too.
	placed := placeBlock(20, 5, "hi")
	lines := strings.Split(placed, "\n")
	wantTopPad := (5 - 1) / 2 // 1 line of content
	if len(lines) != wantTopPad+1 {
		t.Fatalf("expected %d blank lines + 1 content line, got %d lines: %q", wantTopPad, len(lines), lines)
	}
	if !strings.Contains(lines[len(lines)-1], "hi") {
		t.Fatalf("expected content on the last line, got %q", lines)
	}
	for _, l := range lines[:len(lines)-1] {
		if strings.TrimSpace(l) != "" {
			t.Errorf("expected blank top-padding lines, got %q", l)
		}
	}
}
