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

func TestBallDimensions_FloorOnSmallPanelsGrowthOnLarge(t *testing.T) {
	cases := []struct {
		name           string
		innerW, innerH int
		wantH          int
	}{
		// A small terminal gets exactly the old fixed ball.
		{"small panel keeps the 24-row floor", 100, 22, discoBallHeight},
		// Plenty of both dimensions: capped at the max.
		{"huge panel caps at discoBallMaxHeight", 250, 60, discoBallMaxHeight},
		// Tall but narrow: width is the binding constraint --
		// (innerW - gap - menu col - slack) / 2, rounded down to even.
		{"narrow panel is width-bound", 124, 60, (124 - dashboardGapWidth - menuRightColWidth - 4) / 2 / 2 * 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h, w := ballDimensions(c.innerW, c.innerH)
			if h != c.wantH {
				t.Errorf("ballDimensions(%d, %d) h = %d, want %d", c.innerW, c.innerH, h, c.wantH)
			}
			if w != 2*h {
				t.Errorf("w = %d, want the 2:1 terminal-cell aspect (%d)", w, 2*h)
			}
			if h%2 != 0 {
				t.Errorf("h = %d, want an even row count", h)
			}
		})
	}
}

func TestBallGridFor_CachesAStableGridPerSize(t *testing.T) {
	a := ballGridFor(28)
	b := ballGridFor(28)
	if &a[0] != &b[0] {
		t.Error("two calls for the same height built different grids -- the shape must stay stable across frames")
	}
	if len(a) != 28 || len(a[0]) != 56 {
		t.Errorf("grid is %dx%d, want 28x56", len(a), len(a[0]))
	}
}

// TestRenderDashboardPanel_CenteredLayoutFillsTallPanels pins the fix
// for a real complaint (screenshot, 2026-07-16): at a large window the
// menu content pinned to the panel's top-left corner over a sea of
// empty space. Centered layout must put blank rows above the content;
// the boot screen's top layout must not (its streaming log lines would
// jiggle a centered composition).
func TestRenderDashboardPanel_CenteredLayoutFillsTallPanels(t *testing.T) {
	body := "1. Practice\n2. Daily"

	centered := renderDashboardPanel(200, 60, 0, body, layoutCentered, "")
	lines := strings.Split(centered, "\n")
	// Row 1 is the border; rows 2-3 are padding+slack. Find the first
	// row with real content and require it to sit well below the top.
	firstContent := -1
	for i, l := range lines {
		inner := strings.Trim(stripAnsiTUI(l), "║╔╗╚╝═ ")
		if inner != "" {
			firstContent = i
			break
		}
	}
	if firstContent < 5 {
		t.Errorf("centered layout's first content row is %d, want it pushed down by vertical centering:\n%s", firstContent, stripAnsiTUI(centered))
	}

	top := renderDashboardPanel(200, 60, 0, body, layoutTop, "")
	lines = strings.Split(top, "\n")
	firstContent = -1
	for i, l := range lines {
		inner := strings.Trim(stripAnsiTUI(l), "║╔╗╚╝═ ")
		if inner != "" {
			firstContent = i
			break
		}
	}
	if firstContent > 4 {
		t.Errorf("top layout's first content row is %d, want it near the top (boot log stability)", firstContent)
	}
}
