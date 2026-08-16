# Merge Statement Lines

Two monthly statements each report how many times a card was charged at
each merchant. Finance wants one combined view covering both months.

Write `merge_counts(a: dict[str, int], b: dict[str, int]) -> dict[str, int]`
returning a new dict containing every merchant from either input, with
the counts summed for merchants present in both. The keys of the
returned dict must be in sorted order.

## Examples

    merge_counts({"cafe": 2, "gas": 1}, {"cafe": 3})  -> {"cafe": 5, "gas": 1}
    merge_counts({}, {"gym": 4})                      -> {"gym": 4}
    merge_counts({}, {})                              -> {}

## Constraints

- 0 <= len(a), len(b) <= 1_000
- 1 <= count <= 10_000
