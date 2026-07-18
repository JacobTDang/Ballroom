# Per-kilogram shipping rate for each delivery zone, in cents.
_ZONE_RATE_PER_KG = {"A": 200, "B": 500, "C": 900}
_BASE_FEE_CENTS = 500

# Memoizes quotes so repeated lookups skip the rate table -- the
# calculation is treated as expensive enough to cache.
_shipping_cache: dict[tuple[int, str], int] = {}


def shipping_cost(weight_kg: int, zone: str) -> int:
    """Return the shipping cost, in cents, for a package of weight_kg
    shipped to zone."""
    key = (weight_kg, zone)
    if key in _shipping_cache:
        return _shipping_cache[key]
    cost = _BASE_FEE_CENTS + weight_kg * _ZONE_RATE_PER_KG[zone]
    _shipping_cache[key] = cost
    return cost
