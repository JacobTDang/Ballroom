package tutor

import (
	"testing"

	"github.com/JacobTDang/Ballroom/internal/palette"
)

// The activity feed's pulse colors are untyped constants rather than
// palette lookups, because the blend math uses them in both integer and
// float contexts and only untyped constants satisfy both. That buys
// flexibility at the cost of a second copy, so this test is what keeps
// the copy honest -- change the palette and this fails until the
// constants follow.
func TestActivityColorsMatchPalette(t *testing.T) {
	r, g, b := palette.RGB(palette.Teal)
	if activityPulseBaseR != r || activityPulseBaseG != g || activityPulseBaseB != b {
		t.Errorf("pulse base = %d,%d,%d, want palette.Teal %d,%d,%d",
			activityPulseBaseR, activityPulseBaseG, activityPulseBaseB, r, g, b)
	}
	r, g, b = palette.RGB(palette.PulseGlow)
	if activityPulseGlowR != r || activityPulseGlowG != g || activityPulseGlowB != b {
		t.Errorf("pulse glow = %d,%d,%d, want palette.PulseGlow %d,%d,%d",
			activityPulseGlowR, activityPulseGlowG, activityPulseGlowB, r, g, b)
	}
}
