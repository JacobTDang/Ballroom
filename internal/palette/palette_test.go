package palette

import (
	"go/build"
	"strings"
	"testing"
)

// TestColorValues pins every color to the value it had before this
// package existed. The consolidation is meant to be invisible: if one
// of these changes, the refactor moved a color rather than just
// relocating it, and the "no visual change" claim is false.
func TestColorValues(t *testing.T) {
	want := map[string]string{
		"Teal": "#2FA6A6", "Pink": "#E0468C", "Gold": "#E8A93C", "Red": "#F03C3C",
		"Purple": "#9B5FB0", "Blue": "#3C7DC4", "Orange": "#F0862E", "Cyan": "#3ED6D6",
		"Cream": "#F2EBDD", "PaleGray": "#D9D3C4", "WarmGray": "#96918B",
		"MidGray": "#8B8680", "DimGray": "#6B6B6B",
		"Rule": "#3A3D4D", "InputRule": "#1E5A5A", "Ink": "#000000",
		"CardBg": "#14151C", "CardHeaderBg": "#1E2029", "GutterFg": "#5C5852",
	}
	got := map[string]string{
		"Teal": Teal, "Pink": Pink, "Gold": Gold, "Red": Red,
		"Purple": Purple, "Blue": Blue, "Orange": Orange, "Cyan": Cyan,
		"Cream": Cream, "PaleGray": PaleGray, "WarmGray": WarmGray,
		"MidGray": MidGray, "DimGray": DimGray,
		"Rule": Rule, "InputRule": InputRule, "Ink": Ink,
		"CardBg": CardBg, "CardHeaderBg": CardHeaderBg, "GutterFg": GutterFg,
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("%s = %s, want %s", name, got[name], w)
		}
	}
	if len(All()) != len(want) {
		t.Errorf("All() has %d colors, want %d — a new color needs pinning here too", len(All()), len(want))
	}
}

func TestANSIFg(t *testing.T) {
	if got, want := ANSIFg(Teal), "\x1b[38;2;47;166;166m"; got != want {
		t.Errorf("ANSIFg(Teal) = %q, want %q", got, want)
	}
}

func TestANSIBg(t *testing.T) {
	if got, want := ANSIBg(CardBg), "\x1b[48;2;20;21;28m"; got != want {
		t.Errorf("ANSIBg(CardBg) = %q, want %q", got, want)
	}
}

func TestRGB(t *testing.T) {
	r, g, b := RGB(Pink)
	if r != 224 || g != 70 || b != 140 {
		t.Errorf("RGB(Pink) = %d,%d,%d, want 224,70,140", r, g, b)
	}
}

// TestRGBPanicsOnGarbage keeps the existing contract: a malformed
// constant is a build-out typo, and rendering garbage would hide it.
func TestRGBPanicsOnGarbage(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("RGB on a malformed color did not panic")
		}
	}()
	RGB("not-a-color")
}

func TestContainsIsCaseInsensitive(t *testing.T) {
	if !Contains("#2fa6a6") {
		t.Error("Contains should match the container config's lowercase spelling")
	}
	if Contains("#123456") {
		t.Error("Contains matched a color that isn't in the palette")
	}
}

// TestNoInternalImports is the structural guarantee that makes this
// package usable from every home: the tutor pane cannot import the TUI,
// and the container builds its own binary from a subset of the tree, so
// a palette that reached back into internal/ would be unusable from at
// least one of them.
func TestNoInternalImports(t *testing.T) {
	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		t.Fatalf("import dir: %v", err)
	}
	for _, imp := range pkg.Imports {
		if strings.Contains(imp, "JacobTDang/Ballroom/internal/") {
			t.Errorf("palette imports %s — it must stay a leaf package", imp)
		}
	}
}
