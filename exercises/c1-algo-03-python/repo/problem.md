# Next Higher Spend

A budgeting feature shows cardholders, for each day of a statement
period, how long they had to wait for a day with strictly higher
spending. You get the period's daily spend totals in order.

Write `next_higher(spend: list[int]) -> list[int]` returning a list
where entry `i` is `j - i` for the nearest day `j > i` with
`spend[j] > spend[i]`, or `0` if no later day has strictly higher
spend.

## Examples

    next_higher([30, 40, 50, 20])                 -> [1, 1, 0, 0]
    next_higher([73, 74, 75, 71, 69, 72, 76, 73]) -> [1, 1, 4, 2, 1, 1, 0, 0]
    next_higher([50, 40, 30])                     -> [0, 0, 0]

## Constraints

- 0 <= len(spend) <= 100_000
- 0 <= spend[i] <= 1_000_000_000
