# Longest Spending Streak

A card's daily spending totals are recorded in date order. Analysts
want the longest stretch of consecutive days where each day's total was
strictly higher than the day before it.

Write `longest_streak(amounts: list[int]) -> int` returning the length
of the longest strictly increasing run of consecutive entries. An empty
list has streak 0; a single entry counts as a streak of 1.

## Examples

    longest_streak([10, 20, 30, 5, 6])  -> 3
    longest_streak([7, 7, 7])           -> 1
    longest_streak([])                  -> 0

## Constraints

- 0 <= len(amounts) <= 10_000
- 0 <= amount <= 1_000_000
