package tui

import "testing"

func TestDashboardBallSize_UsesFullHeightWhenWidthAllows(t *testing.T) {
	h, w := dashboardBallSize(300, 30)
	if h != 30 {
		t.Errorf("h = %d, want 30 (full available height)", h)
	}
	if w != 60 {
		t.Errorf("w = %d, want 60 (2x height, circular aspect)", w)
	}
}

func TestDashboardBallSize_ShrinksToFitWidthConstraint(t *testing.T) {
	h, w := dashboardBallSize(40, 30)
	if w > 40 {
		t.Errorf("w = %d, must not exceed available width 40", w)
	}
	if h != w/2 {
		t.Errorf("h = %d, w = %d: expected h == w/2 to keep circular aspect", h, w)
	}
	if h >= 30 {
		t.Errorf("h = %d, expected it to shrink below the full available height 30", h)
	}
}

func TestDashboardBallSize_NeverBelowMinimum(t *testing.T) {
	h, _ := dashboardBallSize(200, 2)
	if h < minBallHeight {
		t.Errorf("h = %d, want at least minBallHeight = %d", h, minBallHeight)
	}
}

func TestBallAreaSize_ReservesMarginsBorderPaddingAndRightColumn(t *testing.T) {
	maxW, maxH := ballAreaSize(120, 40, 54)
	wantW := 120 - dashboardMarginW - dashboardBorderPadW - 54 - dashboardGapWidth
	wantH := 40 - dashboardMarginH - dashboardBorderPadH
	if maxW != wantW {
		t.Errorf("maxW = %d, want %d", maxW, wantW)
	}
	if maxH != wantH {
		t.Errorf("maxH = %d, want %d", maxH, wantH)
	}
}
