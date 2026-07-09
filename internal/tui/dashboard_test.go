package tui

import "testing"

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
