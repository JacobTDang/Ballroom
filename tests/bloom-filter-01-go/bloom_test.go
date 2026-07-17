package main

import (
	"fmt"
	"testing"
)

func TestZeroFalseNegatives(t *testing.T) {
	b := NewBloomFilter(16384, 4)
	for i := 0; i < 500; i++ {
		b.Add(fmt.Sprintf("present-%d", i))
	}
	for i := 0; i < 500; i++ {
		if !b.MightContain(fmt.Sprintf("present-%d", i)) {
			t.Fatalf("added key present-%d reported absent -- a bloom filter must never false-negative", i)
		}
	}
}

func TestFalsePositiveRateBounded(t *testing.T) {
	b := NewBloomFilter(16384, 4)
	for i := 0; i < 500; i++ {
		b.Add(fmt.Sprintf("present-%d", i))
	}
	falsePositives := 0
	const probes = 10000
	for i := 0; i < probes; i++ {
		if b.MightContain(fmt.Sprintf("absent-%d", i)) {
			falsePositives++
		}
	}
	// Theory for 16384 bits / 4 hashes / 500 keys: ~0.02%. Budget: 2%.
	if falsePositives >= probes*2/100 {
		t.Fatalf("%d of %d absent keys reported present (%.1f%%) -- the false-positive rate must stay under 2%%",
			falsePositives, probes, 100*float64(falsePositives)/probes)
	}
}

func TestEmptyFilterContainsNothing(t *testing.T) {
	b := NewBloomFilter(1024, 3)
	for i := 0; i < 100; i++ {
		if b.MightContain(fmt.Sprintf("anything-%d", i)) {
			t.Fatalf("empty filter reported anything-%d present", i)
		}
	}
}
