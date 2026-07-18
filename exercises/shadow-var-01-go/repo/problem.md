# Shadowed result: `MaxBelowLimit`

`MaxBelowLimit` is supposed to return the largest value in a slice of
integers that is less than or equal to `limit`, or `-1` if no value
qualifies. It compiles and runs without crashing, but it always
returns `-1` — even when a qualifying value clearly exists.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
MaxBelowLimit([3, 7, 2, 9, 5], 7) => 7
MaxBelowLimit([3, 7, 2, 9, 5], 6) => 5
MaxBelowLimit([10, 20, 30], 5)    => -1
```

## Constraints

- Values may be negative.
- Return -1 if no value in the slice is <= limit.
