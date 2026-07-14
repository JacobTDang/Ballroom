# Longest Increasing Path in a Matrix

Given an `m x n` integer matrix, return the length of the longest
increasing path in the matrix.

From each cell, you can move in four directions: left, right, up, or
down. You may not move diagonally or move outside the boundary (i.e.,
wrap-around is not allowed).

## Example

```
Input: matrix = [[9,9,4],[6,6,8],[2,1,1]]
Output: 4
Explanation: The longest increasing path is [1, 2, 6, 9].
```

## Constraints

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 200`
- `0 <= matrix[i][j] <= 2^31 - 1`
