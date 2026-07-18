# Remove while iterating: `RemoveValue`

`RemoveValue` is supposed to remove every occurrence of `target` from
a list of integers, in place, and return it. It runs without
crashing, but when `target` appears in two or more consecutive
positions, some of those matches survive.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
RemoveValue([1, 2, 2, 2, 3], 2) => [1, 3]
RemoveValue([2, 2, 5, 2], 2)    => [5]
RemoveValue([1, 3, 5], 9)       => [1, 3, 5]
```

## Constraints

- Relative order of the surviving elements is preserved.
- `target` may appear zero or more times, anywhere, including runs of
  consecutive matches.
