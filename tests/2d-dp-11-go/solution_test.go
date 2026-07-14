package main

import "testing"

func TestIsMatch_NoStarMismatch(t *testing.T) {
	if got := IsMatch("aa", "a"); got != false {
		t.Errorf("IsMatch(%q, %q) = %v, want false", "aa", "a", got)
	}
}

func TestIsMatch_StarRepeat(t *testing.T) {
	if got := IsMatch("aa", "a*"); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "aa", "a*", got)
	}
}

func TestIsMatch_Classic(t *testing.T) {
	if got := IsMatch("aab", "c*a*b"); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "aab", "c*a*b", got)
	}
}

func TestIsMatch_LongNoMatch(t *testing.T) {
	if got := IsMatch("mississippi", "mis*is*p*."); got != false {
		t.Errorf("IsMatch(%q, %q) = %v, want false", "mississippi", "mis*is*p*.", got)
	}
}

func TestIsMatch_BothEmpty(t *testing.T) {
	if got := IsMatch("", ""); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "", "", got)
	}
}

func TestIsMatch_EmptyStringStarZero(t *testing.T) {
	if got := IsMatch("", "a*"); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "", "a*", got)
	}
}

func TestIsMatch_DotMatchesAny(t *testing.T) {
	if got := IsMatch("a", "."); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "a", ".", got)
	}
}

func TestIsMatch_DotStarMatchesAll(t *testing.T) {
	if got := IsMatch("ab", ".*"); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "ab", ".*", got)
	}
}

func TestIsMatch_LongerMatch(t *testing.T) {
	if got := IsMatch("mississippi", "mis*is*ip*."); got != true {
		t.Errorf("IsMatch(%q, %q) = %v, want true", "mississippi", "mis*is*ip*.", got)
	}
}

func TestIsMatch_DotStarTrailingLiteralFails(t *testing.T) {
	if got := IsMatch("ab", ".*c"); got != false {
		t.Errorf("IsMatch(%q, %q) = %v, want false", "ab", ".*c", got)
	}
}
