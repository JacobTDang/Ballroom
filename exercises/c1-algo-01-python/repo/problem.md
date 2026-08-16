# Top Accounts by Activity

Fraud analysts want a quick view of which accounts generated the most
events today. You receive the day's event log as a list of account ids;
each entry records one event for that account.

Write `top_k(events: list[str], k: int) -> list[str]` returning the `k`
account ids with the most events, ordered by event count descending.
Accounts with equal counts are ordered lexicographically ascending.

## Examples

    top_k(["acct9", "acct1", "acct9", "acct4", "acct1", "acct9"], 2) -> ["acct9", "acct1"]
    top_k(["b", "a", "c"], 2)                                        -> ["a", "b"]
    top_k(["x"], 1)                                                  -> ["x"]

## Constraints

- 0 <= len(events) <= 100_000
- account ids are non-empty strings of lowercase letters and digits
- 0 <= k <= number of distinct account ids
