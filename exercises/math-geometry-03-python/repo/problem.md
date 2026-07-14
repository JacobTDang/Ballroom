# Set Matrix Zeroes

Given an `m x n` matrix, if an element is `0`, set its entire row and
column to `0`, in place.

## Example

```
Input: matrix = [[1,1,1],
                  [1,0,1],
                  [1,1,1]]

Output: [[1,0,1],
         [0,0,0],
         [1,0,1]]
```

## Constraints

- `m == matrix.length`
- `n == matrix[0].length`
- `1 <= m, n <= 20`
- `-2^31 <= matrix[i][j] <= 2^31 - 1`
