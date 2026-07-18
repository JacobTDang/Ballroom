#include <map>
#include <string>

// Per-kilogram shipping rate for each delivery zone, in cents.
static const std::map<std::string, int> kZoneRatePerKg = {
    {"A", 200}, {"B", 500}, {"C", 900},
};
static const int kBaseFeeCents = 500;

// Memoizes quotes so repeated lookups skip the rate table -- the
// calculation is treated as expensive enough to cache.
static std::map<int, int> shipping_cache;

// Returns the shipping cost, in cents, for a package of weight_kg
// shipped to zone. Currently returns a stale price when the same
// weight is quoted for a different zone right after an earlier quote
// -- find and fix the bug.
int shipping_cost(int weight_kg, const std::string& zone) {
    auto it = shipping_cache.find(weight_kg);
    if (it != shipping_cache.end()) {
        return it->second;
    }
    int cost = kBaseFeeCents + weight_kg * kZoneRatePerKg.at(zone);
    shipping_cache[weight_kg] = cost;
    return cost;
}
