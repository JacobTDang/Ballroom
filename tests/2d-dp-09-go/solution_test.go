package main

import "testing"

func TestMinDistance_Classic(t *testing.T) {
	if got := MinDistance("horse", "ros"); got != 3 {
		t.Errorf("MinDistance(%q, %q) = %d, want 3", "horse", "ros", got)
	}
}

func TestMinDistance_SecondClassic(t *testing.T) {
	if got := MinDistance("intention", "execution"); got != 5 {
		t.Errorf("MinDistance(%q, %q) = %d, want 5", "intention", "execution", got)
	}
}

func TestMinDistance_EmptyFirst(t *testing.T) {
	if got := MinDistance("", "abc"); got != 3 {
		t.Errorf("MinDistance(%q, %q) = %d, want 3", "", "abc", got)
	}
}

func TestMinDistance_Identical(t *testing.T) {
	if got := MinDistance("abc", "abc"); got != 0 {
		t.Errorf("MinDistance(%q, %q) = %d, want 0", "abc", "abc", got)
	}
}

func TestMinDistance_BothEmpty(t *testing.T) {
	if got := MinDistance("", ""); got != 0 {
		t.Errorf("MinDistance(%q, %q) = %d, want 0", "", "", got)
	}
}

func TestMinDistance_EmptySecond(t *testing.T) {
	if got := MinDistance("abc", ""); got != 3 {
		t.Errorf("MinDistance(%q, %q) = %d, want 3", "abc", "", got)
	}
}

func TestMinDistance_SingleCharReplace(t *testing.T) {
	if got := MinDistance("a", "b"); got != 1 {
		t.Errorf("MinDistance(%q, %q) = %d, want 1", "a", "b", got)
	}
}

func TestMinDistance_PureInsertion(t *testing.T) {
	if got := MinDistance("cat", "cats"); got != 1 {
		t.Errorf("MinDistance(%q, %q) = %d, want 1", "cat", "cats", got)
	}
}

func TestMinDistance_MultiOpMix(t *testing.T) {
	if got := MinDistance("sunday", "saturday"); got != 3 {
		t.Errorf("MinDistance(%q, %q) = %d, want 3", "sunday", "saturday", got)
	}
}
