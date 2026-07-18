package main

// zoneRatePerKg is the per-kilogram shipping rate for each delivery
// zone, in cents.
var zoneRatePerKg = map[string]int{
	"A": 200,
	"B": 500,
	"C": 900,
}

const baseFeeCents = 500

// shippingCache memoizes quotes so repeated lookups skip the rate
// table -- the calculation is treated as expensive enough to cache.
var shippingCache = map[int]int{}

// ShippingCost returns the shipping cost, in cents, for a package of
// weightKg shipped to zone. Currently returns a stale price when the
// same weight is quoted for a different zone right after an earlier
// quote -- find and fix the bug.
func ShippingCost(weightKg int, zone string) int {
	if cost, ok := shippingCache[weightKg]; ok {
		return cost
	}
	cost := baseFeeCents + weightKg*zoneRatePerKg[zone]
	shippingCache[weightKg] = cost
	return cost
}
