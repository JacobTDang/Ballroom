# Spiral Matrix

Given an `m x n` matrix, return all elements of the matrix in spiral
order: starting at the top-left corner, walk right along the top row,
down the right column, left along the bottom row, and up the left
column, shrinking the boundary inward one layer at a time until every
element has been visited.

## Example

```
Input: matrix = [[1,2,3],
                  [4,5,6],
                  [7,8,9]]

Output: [1,2,3,6,9,8,7,4,5]
```

## Constraints

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 10`
- `-100 <= matrix[i][j] <= 100`
