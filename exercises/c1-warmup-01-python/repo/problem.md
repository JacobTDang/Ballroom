# Transaction Ledger Balance

A checking account starts at 0. You receive the day's transactions as a
list of strings, each either `"D <amount>"` (deposit) or `"W <amount>"`
(withdrawal), amounts are positive integers. A withdrawal that would
take the balance below 0 is declined and skipped.

Write `final_balance(transactions: list[str]) -> tuple[int, int]`
returning `(balance, declined)` — the end-of-day balance and how many
withdrawals were declined.

## Examples

    final_balance(["D 100", "W 30", "W 90"])  -> (70, 1)
    final_balance(["W 10"])                   -> (0, 1)
    final_balance([])                         -> (0, 0)

## Constraints

- 0 <= len(transactions) <= 10_000
- 1 <= amount <= 10_000
