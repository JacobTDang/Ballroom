package tui

import (
	"strings"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/catalog"
)

func TestPanelDimensions_SubtractsMargin(t *testing.T) {
	w, h := panelDimensions(160, 50)
	if w != 160-dashboardMarginW {
		t.Errorf("w = %d, want %d", w, 160-dashboardMarginW)
	}
	if h != 50-dashboardMarginH {
		t.Errorf("h = %d, want %d", h, 50-dashboardMarginH)
	}
}

func TestPanelDimensions_ClampsToMinimumOnTinyTerminal(t *testing.T) {
	w, h := panelDimensions(5, 3)
	if w < minPanelWidth {
		t.Errorf("w = %d, want at least minPanelWidth = %d", w, minPanelWidth)
	}
	if h < minPanelHeight {
		t.Errorf("h = %d, want at least minPanelHeight = %d", h, minPanelHeight)
	}
}

func TestMinPanelWidth_AlwaysFitsBallPlusGapPlusSomeText(t *testing.T) {
	// If the floor itself is narrower than ball+gap, every render at the
	// minimum size would still overflow and wrap — the whole point of
	// having a floor.
	const wantHeadroom = 20
	if minPanelWidth < discoBallWidth+dashboardGapWidth+wantHeadroom {
		t.Errorf("minPanelWidth = %d, too small to fit the ball (%d) + gap (%d) + %d cols of text",
			minPanelWidth, discoBallWidth, dashboardGapWidth, wantHeadroom)
	}
}

func TestDashboardBanner_UsesFullBannerWhenWidthAllows(t *testing.T) {
	got := dashboardBanner(0, 200)
	want := catalog.MosaicBannerScaled(0, 1)
	if got != want {
		t.Error("expected the full animated banner when there's plenty of width")
	}
}

func TestDashboardBanner_FallsBackToCompactWhenNarrow(t *testing.T) {
	got := dashboardBanner(0, 10)
	want := catalog.CompactBanner()
	if got != want {
		t.Error("expected the compact single-line banner when width is too tight for the full one")
	}
}

func TestCenterRightColumn_AddsHalfTheSlackAsLeftMargin(t *testing.T) {
	// centerRightColumn only owns the left margin — the caller
	// (renderDashboardPanel's outer bordered/Width()-styled box) fills
	// the remaining trailing space on the right automatically, so
	// leaving exactly half the slack (avail-contentWidth) as a left
	// margin is what makes the final rendered panel come out balanced.
	got := centerRightColumn("abc", 11) // avail=11, content=3, slack=8 -> want left=4
	line := strings.Split(got, "\n")[0]
	left := len(line) - len(strings.TrimLeft(line, " "))
	if left != 4 {
		t.Errorf("expected left margin of 4 (half of 8 slack), got %d in %q", left, line)
	}
}

func TestCenterRightColumn_NoOpWhenContentAlreadyFillsAvailWidth(t *testing.T) {
	got := centerRightColumn("hello world", 8) // content wider than avail
	if got != "hello world" {
		t.Errorf("expected content returned unchanged when it doesn't fit, got %q", got)
	}
}

func TestCenterRightColumn_PreservesRelativeAlignmentAcrossLines(t *testing.T) {
	// Every line in the block should get the same left margin, so a
	// multi-line block (banner above body text) keeps its own internal
	// alignment instead of each line being centered independently.
	got := centerRightColumn("ab\nabcde", 15)
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), lines)
	}
	left0 := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
	left1 := len(lines[1]) - len(strings.TrimLeft(lines[1], " "))
	if left0 != left1 {
		t.Errorf("expected both lines to share the same left margin, got %d and %d", left0, left1)
	}
}
