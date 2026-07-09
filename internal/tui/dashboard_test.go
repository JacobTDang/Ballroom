package tui

import (
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
