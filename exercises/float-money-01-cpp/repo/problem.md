# The Bill That Never Settles

`settles_bill` checks whether a list of itemized receipt amounts sums
to the printed bill total -- used to flag receipts where the line
items don't add up before they go out to a customer. It's right most
of the time, but some perfectly valid receipts get flagged as broken:
three dimes on the receipt somehow don't equal thirty cents.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
settles_bill([10.0, 20.0, 70.0], 100.0)  => true
settles_bill([5.00, 5.00], 11.00)        => false
```

## Constraints

- Amounts are always non-negative.
- The list may be empty (an empty receipt settles a zero bill).
- Money is in dollars; two amounts "settle" if they match to the
  nearest cent.
