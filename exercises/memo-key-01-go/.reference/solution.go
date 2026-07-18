package main

// zoneRatePerKg is the per-kilogram shipping rate for each delivery
// zone, in cents.
var zoneRatePerKg = map[string]int{
	"A": 200,
	"B": 500,
	"C": 900,
}

const baseFeeCents = 500

// shippingKey identifies one (weight, zone) quote.
type shippingKey struct {
	weightKg int
	zone     string
}

// shippingCache memoizes quotes so repeated lookups skip the rate
// table -- the calculation is treated as expensive enough to cache.
var shippingCache = map[shippingKey]int{}

// ShippingCost returns the shipping cost, in cents, for a package of
// weightKg shipped to zone.
func ShippingCost(weightKg int, zone string) int {
	key := shippingKey{weightKg, zone}
	if cost, ok := shippingCache[key]; ok {
		return cost
	}
	cost := baseFeeCents + weightKg*zoneRatePerKg[zone]
	shippingCache[key] = cost
	return cost
}
