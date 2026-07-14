package main

import "testing"

func TestNumDistinct_Classic(t *testing.T) {
	if got := NumDistinct("rabbbit", "rabbit"); got != 3 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 3", "rabbbit", "rabbit", got)
	}
}

func TestNumDistinct_SecondClassic(t *testing.T) {
	if got := NumDistinct("babgbag", "bag"); got != 5 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 5", "babgbag", "bag", got)
	}
}

func TestNumDistinct_ExactMatch(t *testing.T) {
	if got := NumDistinct("abc", "abc"); got != 1 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 1", "abc", "abc", got)
	}
}

func TestNumDistinct_TargetLonger(t *testing.T) {
	if got := NumDistinct("abc", "abcd"); got != 0 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 0", "abc", "abcd", got)
	}
}

func TestNumDistinct_EmptyTarget(t *testing.T) {
	if got := NumDistinct("abc", ""); got != 1 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 1", "abc", "", got)
	}
}

func TestNumDistinct_EmptySource(t *testing.T) {
	if got := NumDistinct("", "abc"); got != 0 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 0", "", "abc", got)
	}
}

func TestNumDistinct_BothEmpty(t *testing.T) {
	if got := NumDistinct("", ""); got != 1 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 1", "", "", got)
	}
}

func TestNumDistinct_RepeatedCharsCombinatoric(t *testing.T) {
	if got := NumDistinct("aaaa", "aa"); got != 6 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 6", "aaaa", "aa", got)
	}
}

func TestNumDistinct_SingleCharManyOccurrences(t *testing.T) {
	if got := NumDistinct("aaa", "a"); got != 3 {
		t.Errorf("NumDistinct(%q, %q) = %d, want 3", "aaa", "a", got)
	}
}
