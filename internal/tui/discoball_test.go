package tui

import (
	"strings"
	"testing"
)

func TestBuildDiscoBall_ProducesRequestedDimensions(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	if len(grid) != 24 {
		t.Fatalf("expected 24 rows, got %d", len(grid))
	}
	for i, row := range grid {
		if len(row) != 48 {
			t.Fatalf("row %d: expected 48 cols, got %d", i, len(row))
		}
	}
}

func TestBuildDiscoBall_CornersAreOutsideTheCircle(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	if grid[0][0].ch != ' ' {
		t.Errorf("expected top-left corner to be outside the circle (blank), got %q", grid[0][0].ch)
	}
	if grid[0][47].ch != ' ' {
		t.Errorf("expected top-right corner to be outside the circle (blank), got %q", grid[0][47].ch)
	}
}

func TestBuildDiscoBall_CenterIsInsideTheCircle(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	center := grid[12][24]
	if center.ch == ' ' {
		t.Error("expected the ball's center to be inside the circle (non-blank)")
	}
}

func TestBuildDiscoBall_OnlyUsesShadingCharacters(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	allowed := map[rune]bool{' ': true}
	for _, ch := range discoShades {
		allowed[ch] = true
	}
	for r, row := range grid {
		for c, cell := range row {
			if !allowed[cell.ch] {
				t.Fatalf("cell (%d,%d) has unexpected char %q, not in the shading set", r, c, cell.ch)
			}
		}
	}
}

func TestBuildDiscoBall_SparkleFractionIsSparse(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	total, sparkles := 0, 0
	for _, row := range grid {
		for _, cell := range row {
			if cell.ch == ' ' {
				continue
			}
			total++
			if cell.sparkle {
				sparkles++
			}
		}
	}
	if total == 0 {
		t.Fatal("expected a non-empty ball")
	}
	frac := float64(sparkles) / float64(total)
	if frac < 0.05 || frac > 0.20 {
		t.Errorf("sparkle fraction = %.2f, want roughly 8-15%% (allowing some margin)", frac)
	}
}

func TestRenderDiscoBall_ShapeStableAcrossPhases(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	a := stripAnsiTUI(renderDiscoBall(grid, 0))
	b := stripAnsiTUI(renderDiscoBall(grid, 1))
	if a != b {
		t.Error("expected the underlying shape/characters to stay identical across phases — only color should change")
	}
}

func TestSparkleColorIndex_ChangesAcrossPhases(t *testing.T) {
	// lipgloss silently drops ANSI codes outside a real terminal (true
	// in `go test`), so comparing rendered strings can't detect a color
	// change here — test the underlying color-selection logic directly
	// instead of the string lipgloss happens to produce in this
	// environment.
	a := sparkleColorIndex(0, 3, 5)
	b := sparkleColorIndex(1, 3, 5)
	if a == b {
		t.Error("expected sparkleColorIndex to pick a different color for a different phase")
	}
}

func TestSparkleColorIndex_StaggeredByCellPosition(t *testing.T) {
	// Two different cells at the same phase should land on different
	// colors, so sparkles don't all change in lockstep.
	a := sparkleColorIndex(0, 1, 1)
	b := sparkleColorIndex(0, 2, 0)
	if a == b {
		t.Error("expected different cell positions to offset into different colors at the same phase")
	}
}

func TestRenderDiscoBall_NoTrailingBlankLine(t *testing.T) {
	grid := buildDiscoBall(24, 48)
	out := renderDiscoBall(grid, 0)
	if strings.HasSuffix(out, "\n") {
		t.Error("expected no trailing newline, callers control their own line joining")
	}
}
