# Surrounded Regions

You are given an `m x n` matrix `board` containing letters `'X'` and
`'O'`. Capture all regions that are 4-directionally surrounded by
`'X'`.

A region is captured by flipping all `'O'`s into `'X'`s in that
surrounded region. A region is NOT surrounded if it touches the edge
of the board, or is connected (4-directionally) to a region that
touches the edge of the board.

Modify `board` in place.

## Example

```
Input: board = [["X","X","X","X"],["X","O","O","X"],["X","X","O","X"],["X","O","X","X"]]
Output: [["X","X","X","X"],["X","X","X","X"],["X","X","X","X"],["X","O","X","X"]]
```

## Constraints

- `m == board.length`
- `n == board[r].length`
- `1 <= m, n <= 200`
- `board[r][c]` is `'X'` or `'O'`
