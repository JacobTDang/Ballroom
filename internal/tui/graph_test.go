package tui

import "testing"

// runeAt indexes a string by rune position, not byte position — needed
// since the connector lines contain multi-byte box-drawing characters.
func runeAt(s string, i int) rune {
	return []rune(s)[i]
}

func TestBoxCenters_SingleBox(t *testing.T) {
	centers := boxCenters([]int{10}, 2)
	if len(centers) != 1 || centers[0] != 5 {
		t.Errorf("boxCenters([10], 2) = %v, want [5]", centers)
	}
}

func TestBoxCenters_MultipleBoxesAccountForGap(t *testing.T) {
	// box0: width 4, center 2. gap 2. box1 starts at 6, width 6, center 9.
	centers := boxCenters([]int{4, 6}, 2)
	want := []int{2, 9}
	if len(centers) != 2 || centers[0] != want[0] || centers[1] != want[1] {
		t.Errorf("boxCenters([4,6], 2) = %v, want %v", centers, want)
	}
}

func TestConnectorLines_SingleChildAlignedWithParent(t *testing.T) {
	lines := connectorLines(10, 5, []int{5})
	if len(lines) != 3 {
		t.Fatalf("expected 3 connector lines, got %d", len(lines))
	}
	if runeAt(lines[0], 5) != '│' {
		t.Errorf("line0 col5 = %q, want stem", string(runeAt(lines[0], 5)))
	}
	if runeAt(lines[1], 5) != '┼' {
		t.Errorf("line1 col5 = %q, want junction (parent aligns with its only child)", string(runeAt(lines[1], 5)))
	}
	if runeAt(lines[2], 5) != '│' {
		t.Errorf("line2 col5 = %q, want stem down to child", string(runeAt(lines[2], 5)))
	}
}

func TestConnectorLines_ParentCenteredBetweenTwoChildren(t *testing.T) {
	lines := connectorLines(10, 5, []int{2, 8})
	if runeAt(lines[0], 5) != '│' {
		t.Errorf("line0 col5 = %q, want parent stem", string(runeAt(lines[0], 5)))
	}
	if runeAt(lines[1], 5) != '┴' {
		t.Errorf("line1 col5 = %q, want '┴' (parent stem meets spine, no child here)", string(runeAt(lines[1], 5)))
	}
	if runeAt(lines[1], 2) != '┬' || runeAt(lines[1], 8) != '┬' {
		t.Errorf("expected '┬' branch at each child center, got %q / %q", string(runeAt(lines[1], 2)), string(runeAt(lines[1], 8)))
	}
	if runeAt(lines[2], 2) != '│' || runeAt(lines[2], 8) != '│' {
		t.Errorf("expected stems down to each child, got %q / %q", string(runeAt(lines[2], 2)), string(runeAt(lines[2], 8)))
	}
	// spine should be a continuous horizontal line between the children
	for i := 2; i <= 8; i++ {
		if runeAt(lines[1], i) == ' ' {
			t.Errorf("spine has a gap at column %d", i)
		}
	}
}

func TestConnectorLines_ManyChildren(t *testing.T) {
	lines := connectorLines(30, 15, []int{2, 8, 15, 22, 28})
	for _, c := range []int{2, 8, 22, 28} {
		if runeAt(lines[1], c) != '┬' {
			t.Errorf("expected '┬' at child col %d, got %q", c, string(runeAt(lines[1], c)))
		}
	}
	if runeAt(lines[1], 15) != '┼' {
		t.Errorf("expected '┼' where parent aligns with the middle child, got %q", string(runeAt(lines[1], 15)))
	}
}

func TestSpanAnchor_OddCountUsesTrueMiddleChild(t *testing.T) {
	got := spanAnchor([]int{7, 25, 43, 63, 82})
	if got != 43 {
		t.Errorf("spanAnchor(odd) = %d, want 43 (the middle child's own center, not the span midpoint)", got)
	}
}

func TestSpanAnchor_EvenCountUsesMidpointOfMiddleTwo(t *testing.T) {
	got := spanAnchor([]int{10, 20, 30, 40})
	if got != 25 {
		t.Errorf("spanAnchor(even) = %d, want 25", got)
	}
}

func TestCenterOffset_BoxNarrowerThanSpan(t *testing.T) {
	// span center at column 20, box width 6 -> should left-pad by 17
	// so the box's own center (offset 3) lands on column 20.
	got := centerOffset(6, 20)
	if got != 17 {
		t.Errorf("centerOffset(6, 20) = %d, want 17", got)
	}
}

func TestCenterOffset_NeverNegative(t *testing.T) {
	got := centerOffset(40, 5)
	if got < 0 {
		t.Errorf("centerOffset should clamp to 0, got %d", got)
	}
}
