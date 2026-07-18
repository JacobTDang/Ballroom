#include <map>
#include <string>
#include <utility>

// Per-kilogram shipping rate for each delivery zone, in cents.
static const std::map<std::string, int> kZoneRatePerKg = {
    {"A", 200}, {"B", 500}, {"C", 900},
};
static const int kBaseFeeCents = 500;

// Memoizes quotes so repeated lookups skip the rate table -- the
// calculation is treated as expensive enough to cache.
static std::map<std::pair<int, std::string>, int> shipping_cache;

// Returns the shipping cost, in cents, for a package of weight_kg
// shipped to zone.
int shipping_cost(int weight_kg, const std::string& zone) {
    auto key = std::make_pair(weight_kg, zone);
    auto it = shipping_cache.find(key);
    if (it != shipping_cache.end()) {
        return it->second;
    }
    int cost = kBaseFeeCents + weight_kg * kZoneRatePerKg.at(zone);
    shipping_cache[key] = cost;
    return cost;
}
