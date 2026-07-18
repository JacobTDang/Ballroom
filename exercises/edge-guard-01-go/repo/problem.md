# Edge guard: `MaxAdjacentDiff`

`MaxAdjacentDiff` is supposed to return the largest absolute difference
between two adjacent elements in a slice of integers, requiring at
least two elements. It currently panics (Go) / crashes under the
sanitizer (C++) / raises `IndexError` (Python) — but only on inputs
with exactly one element.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
MaxAdjacentDiff([3, 1, 4, 1, 5, 9, 2, 6]) => 7
MaxAdjacentDiff([5, 5])                   => 0
MaxAdjacentDiff([-5, -1, -10])            => 9
```

## Constraints

- The input has at least one element.
- Fewer than two elements has no defined answer: the function should
  report that error the same controlled way it already does for an
  empty input, not crash.
