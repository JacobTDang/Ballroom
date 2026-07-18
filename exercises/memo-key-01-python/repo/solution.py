# Per-kilogram shipping rate for each delivery zone, in cents.
_ZONE_RATE_PER_KG = {"A": 200, "B": 500, "C": 900}
_BASE_FEE_CENTS = 500

# Memoizes quotes so repeated lookups skip the rate table -- the
# calculation is treated as expensive enough to cache.
_shipping_cache: dict[int, int] = {}


def shipping_cost(weight_kg: int, zone: str) -> int:
    """Return the shipping cost, in cents, for a package of weight_kg
    shipped to zone. Currently returns a stale price when the same
    weight is quoted for a different zone right after an earlier quote
    -- find and fix the bug."""
    if weight_kg in _shipping_cache:
        return _shipping_cache[weight_kg]
    cost = _BASE_FEE_CENTS + weight_kg * _ZONE_RATE_PER_KG[zone]
    _shipping_cache[weight_kg] = cost
    return cost
