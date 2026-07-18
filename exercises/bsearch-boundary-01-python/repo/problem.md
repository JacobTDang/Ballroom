# Binary search boundary: `FirstAtLeast`

`FirstAtLeast` is supposed to return the index of the first element in
a sorted slice that is greater than or equal to `target`, or the
slice's length if every element is smaller than `target` (the
insertion point at the end). It runs to completion and never crashes,
but for some targets it returns an index one too small.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
FirstAtLeast([1, 3, 5, 7, 9], 6)  => 3
FirstAtLeast([1, 3, 5, 7, 9], 1)  => 0
FirstAtLeast([1, 3, 5, 7, 9], 10) => 5
```

## Constraints

- `values` is sorted ascending and may contain duplicates.
- `target` may be smaller than every element, larger than every
  element, or anywhere in between.
