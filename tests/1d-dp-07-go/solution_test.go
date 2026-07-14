package main

import "testing"

func TestNumDecodings_TwoWays(t *testing.T) {
	if got := NumDecodings("12"); got != 2 {
		t.Errorf("NumDecodings(%q) = %d, want 2", "12", got)
	}
}

func TestNumDecodings_ThreeWays(t *testing.T) {
	if got := NumDecodings("226"); got != 3 {
		t.Errorf("NumDecodings(%q) = %d, want 3", "226", got)
	}
}

func TestNumDecodings_LeadingZero(t *testing.T) {
	if got := NumDecodings("06"); got != 0 {
		t.Errorf("NumDecodings(%q) = %d, want 0", "06", got)
	}
}

func TestNumDecodings_TwoDigitOnly(t *testing.T) {
	if got := NumDecodings("10"); got != 1 {
		t.Errorf("NumDecodings(%q) = %d, want 1", "10", got)
	}
}

func TestNumDecodings_SingleDigit(t *testing.T) {
	if got := NumDecodings("5"); got != 1 {
		t.Errorf("NumDecodings(%q) = %d, want 1", "5", got)
	}
}

func TestNumDecodings_LoneZero(t *testing.T) {
	if got := NumDecodings("0"); got != 0 {
		t.Errorf("NumDecodings(%q) = %d, want 0", "0", got)
	}
}

func TestNumDecodings_JustOverTwentySix(t *testing.T) {
	if got := NumDecodings("27"); got != 1 {
		t.Errorf("NumDecodings(%q) = %d, want 1", "27", got)
	}
}

func TestNumDecodings_UnresolvableZeroPair(t *testing.T) {
	if got := NumDecodings("100"); got != 0 {
		t.Errorf("NumDecodings(%q) = %d, want 0", "100", got)
	}
}

func TestNumDecodings_LongerMultipleWays(t *testing.T) {
	if got := NumDecodings("11106"); got != 2 {
		t.Errorf("NumDecodings(%q) = %d, want 2", "11106", got)
	}
}
