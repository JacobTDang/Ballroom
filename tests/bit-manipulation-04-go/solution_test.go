package main

import "testing"

func TestReverseBits_One(t *testing.T) {
	if got := ReverseBits(1); got != 2147483648 {
		t.Errorf("ReverseBits(1) = %d, want 2147483648", got)
	}
}

func TestReverseBits_Zero(t *testing.T) {
	if got := ReverseBits(0); got != 0 {
		t.Errorf("ReverseBits(0) = %d, want 0", got)
	}
}

func TestReverseBits_AllOnes(t *testing.T) {
	if got := ReverseBits(4294967295); got != 4294967295 {
		t.Errorf("ReverseBits(4294967295) = %d, want 4294967295", got)
	}
}

func TestReverseBits_Two(t *testing.T) {
	if got := ReverseBits(2); got != 1073741824 {
		t.Errorf("ReverseBits(2) = %d, want 1073741824", got)
	}
}

func TestReverseBits_Classic(t *testing.T) {
	if got := ReverseBits(43261596); got != 964176192 {
		t.Errorf("ReverseBits(43261596) = %d, want 964176192", got)
	}
}
