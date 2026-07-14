package main

import "testing"

func TestHammingWeight_Classic(t *testing.T) {
	if got := HammingWeight(11); got != 3 {
		t.Errorf("HammingWeight(11) = %d, want 3", got)
	}
}

func TestHammingWeight_Zero(t *testing.T) {
	if got := HammingWeight(0); got != 0 {
		t.Errorf("HammingWeight(0) = %d, want 0", got)
	}
}

func TestHammingWeight_AllOnes(t *testing.T) {
	if got := HammingWeight(4294967295); got != 32 {
		t.Errorf("HammingWeight(4294967295) = %d, want 32", got)
	}
}

func TestHammingWeight_PowerOfTwo(t *testing.T) {
	if got := HammingWeight(1 << 31); got != 1 {
		t.Errorf("HammingWeight(1<<31) = %d, want 1", got)
	}
}

func TestHammingWeight_AlternatingBits(t *testing.T) {
	if got := HammingWeight(0xAAAAAAAA); got != 16 {
		t.Errorf("HammingWeight(0xAAAAAAAA) = %d, want 16", got)
	}
}
