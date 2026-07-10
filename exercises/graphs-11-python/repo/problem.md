# Number of Connected Components in an Undirected Graph

You have a graph of `n` nodes labeled `0` to `n - 1`. You are given
an integer `n` and an array `edges` where `edges[i] = [a, b]`
indicates that there is an undirected edge between nodes `a` and `b`
in the graph.

Return the number of connected components in the graph.

## Example

```
Input: n = 5, edges = [[0,1],[1,2],[3,4]]
Output: 2
```

## Constraints

- `1 <= n <= 2000`
- `0 <= edges.length <= 5000`
- `edges[i].length == 2`
- `0 <= a, b < n`
- `a != b`
- There are no repeated edges.
