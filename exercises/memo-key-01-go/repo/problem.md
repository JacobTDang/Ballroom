# Stale quote after a zone change: `ShippingCost`

`ShippingCost` quotes a shipping price, in cents, for a package of a
given weight shipped to a given delivery zone. The quote calculation is
treated as expensive enough to cache.

Quotes are correct in isolation. But ask for the same weight again
right after quoting a *different* zone for it, and the new quote is
wrong -- it repeats the earlier zone's price instead of pricing the
zone you actually asked for.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
ShippingCost(5, "A")  => 1500
ShippingCost(5, "B")  => 3000
ShippingCost(5, "C")  => 5000
```

Each of these must be correct even when called right after one of the
others, in any order, for the same weight.

## Constraints

- Zones are one of "A", "B", "C".
- Weight is a positive integer number of kilograms.
