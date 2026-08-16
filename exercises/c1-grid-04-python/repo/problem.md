# Queue Ticker

The lobby display animates waiting customers. Each row of the grid is a
lane; a cell is `"."` (empty) or `">"` (a customer marker). One tick
updates every marker simultaneously:

- a marker in the rightmost column leaves the grid (its cell becomes `"."`);
- any other marker moves one cell right if that cell was `"."` before
  the tick started; otherwise it stays put.

Moves are judged against the pre-tick grid, so a marker whose right
neighbor held a marker at the start of the tick stays put even if that
neighbor moves or leaves during the same tick. Two markers therefore
never share a cell.

Write `step(grid: list[list[str]], n: int) -> list[list[str]]` returning
the grid after `n` ticks as a new list of lists; the input grid must not
be modified.

## Examples

    step([[">", ".", "."]], 1)   -> [[".", ">", "."]]

    step([[">", ">", "."]], 1)   -> [[">", ".", ">"]]   # left marker was blocked at tick start

    step([[">", ">"],
          [".", ">"]], 1)        -> [[">", "."],
                                     [".", "."]]

## Constraints

- 1 <= rows, cols <= 100
- 0 <= n <= 1000
- every cell is "." or ">"
