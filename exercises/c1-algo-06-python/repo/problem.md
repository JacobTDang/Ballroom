# Fair Fee Split

After a group trip, each cardholder's share of the fees is tallied into
a balance: negative means they still owe money, positive means they are
owed, and the balances always sum to exactly 0. The app settles up with
this exact procedure: while any balance is nonzero, pick the person who
owes the most (most negative balance) and the person who is owed the
most (most positive balance); the debtor sends the creditor the smaller
of the two absolute amounts — that is one transfer — and both balances
are updated. Whenever two people are tied for a pick, the one with the
lower index is chosen.

Write `min_transfers_two(balances: list[int]) -> int` returning how
many transfers this procedure makes.

## Examples

    min_transfers_two([4, -2, -2])     -> 2
    min_transfers_two([3, -1, -1, -1]) -> 3
    min_transfers_two([0, 0])          -> 0

## Constraints

- 0 <= len(balances) <= 100_000
- -1_000_000_000 <= balances[i] <= 1_000_000_000
- sum(balances) == 0
