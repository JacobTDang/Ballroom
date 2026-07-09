# Off-by-one: `MaxOf`

`MaxOf` is supposed to return the largest value in a non-empty slice of
integers. It currently panics (Go) / crashes under the sanitizer (C++) /
raises `IndexError` (Python) on every input.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
MaxOf([3, 1, 4, 1, 5, 9, 2, 6]) => 9
MaxOf([-5, -1, -10])            => -1
MaxOf([42])                     => 42
```

## Constraints

- The input is never empty.
- Values may be negative.
