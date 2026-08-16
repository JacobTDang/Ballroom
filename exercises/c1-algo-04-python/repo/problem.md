# Two-Card Sum

A rewards promotion lets a customer pick exactly two gift cards, and
the program covers the cost as long as the combined price stays within
a budget. The app should suggest the pair that uses as much of the
budget as possible.

Write
`pair_within_budget(prices: list[int], budget: int) -> tuple[int, int] | None`
returning indices `(i, j)` with `i < j` that maximize
`prices[i] + prices[j]` without exceeding `budget`, or `None` if no
pair fits. If several pairs reach the same maximum total, return the
one with the smallest `i`; among those, the smallest `j`.

## Examples

    pair_within_budget([40, 60, 25, 30], 90) -> (1, 3)
    pair_within_budget([50, 40, 40, 50], 90) -> (0, 1)
    pair_within_budget([10, 20, 15], 12)     -> None

## Constraints

- 0 <= len(prices) <= 100_000
- 1 <= prices[i] <= 1_000_000_000
- 1 <= budget <= 2_000_000_000
