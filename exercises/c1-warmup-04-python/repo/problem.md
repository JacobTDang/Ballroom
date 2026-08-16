# Rewards Rounding

A rewards program rounds each purchase up to the next multiple of 100
cents and moves the difference into a savings pot. A purchase that is
already an exact multiple of 100 moves nothing.

Write `round_up_total(prices: list[int]) -> int` returning the total
amount moved into the pot for a list of purchase prices in cents.

## Examples

    round_up_total([250, 300, 199])  -> 51
    round_up_total([1])              -> 99
    round_up_total([])               -> 0

## Constraints

- 0 <= len(prices) <= 10_000
- 1 <= price <= 100_000
