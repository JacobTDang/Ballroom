package catalog

import (
	"strings"
	"testing"
)

func TestMosaicBanner_DifferentPhasesProduceDifferentColorOutput(t *testing.T) {
	a := MosaicBanner(0)
	b := MosaicBanner(1)
	if a == b {
		t.Error("expected different phases to change the rendered output (animation would be a no-op otherwise)")
	}
	// but the underlying letters should be unchanged — only color shifts
	if stripAnsi(a) != stripAnsi(b) {
		t.Error("phase should only change color, not the underlying art/text")
	}
}

func TestMosaicBanner_ContainsArtAndTagline(t *testing.T) {
	out := stripAnsi(MosaicBanner(0))
	if !strings.Contains(out, "I N T E R V I E W") {
		t.Errorf("mosaic banner missing tagline:\n%s", out)
	}
	if strings.Count(out, "\n") < 6 {
		t.Errorf("mosaic banner looks too short to contain the ASCII art:\n%s", out)
	}
}

func TestCompactBanner_ContainsWordmark(t *testing.T) {
	out := stripAnsi(CompactBanner())
	if !strings.Contains(out, "BALLROOM") {
		t.Errorf("compact banner missing wordmark:\n%s", out)
	}
	if strings.Count(out, "\n") > 1 {
		t.Errorf("compact banner should be a single line, got:\n%q", out)
	}
}
