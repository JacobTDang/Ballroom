package main

import "testing"

// Weight 5kg, zones queried A -> B -> C in order: each must be correct
// even though the previous call just cached weight 5 for a different
// zone.
func TestShippingCost_SameWeightDifferentZonesInOrder(t *testing.T) {
	cases := []struct {
		zone string
		want int
	}{
		{"A", 1500},
		{"B", 3000},
		{"C", 5000},
	}
	for _, c := range cases {
		if got := ShippingCost(5, c.zone); got != c.want {
			t.Errorf("ShippingCost(5, %q) = %d, want %d", c.zone, got, c.want)
		}
	}
}

// Weight 8kg, zones queried in a different order (C -> A -> B) -- the
// bug must not survive a different call sequence either.
func TestShippingCost_SameWeightDifferentZonesReversedOrder(t *testing.T) {
	cases := []struct {
		zone string
		want int
	}{
		{"C", 7700},
		{"A", 2100},
		{"B", 4500},
	}
	for _, c := range cases {
		if got := ShippingCost(8, c.zone); got != c.want {
			t.Errorf("ShippingCost(8, %q) = %d, want %d", c.zone, got, c.want)
		}
	}
}
