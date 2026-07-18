# Wrong floor: `CountPaths`

`CountPaths` is supposed to count the number of paths from the
top-left to the bottom-right corner of a grid, moving only right or
down, that never step on a blocked cell (marked `1`; open cells are
`0`). It runs without crashing, but it always returns `0`, even when
a path clearly exists.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
CountPaths([[0, 0],
            [0, 0]]) => 2

CountPaths([[0]]) => 1

CountPaths([[0, 0],
            [0, 1]]) => 0   (the destination is blocked)
```

## Constraints

- The grid has at least one row and one column.
- The top-left cell (the start) is never blocked.
- Movement is only right or down.
