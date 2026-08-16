# Alternating Card Digits

A fraud heuristic flags card digit strings whose digits swing between
odd and even the entire way through. Given the digit string, decide
whether the parity flips at every step.

Write `is_alternating(digits: str) -> bool` returning True when every
adjacent pair of digits differs in parity — the sequence reads
odd/even/odd... or even/odd/even... — and False otherwise. A single
digit counts as alternating.

## Examples

    is_alternating("2743")  -> True
    is_alternating("2744")  -> False
    is_alternating("7")     -> True

## Constraints

- 1 <= len(digits) <= 100
- every character is a digit 0-9
