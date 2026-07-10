# Graph Valid Tree

You have a graph of `n` nodes labeled `0` to `n - 1`. You are given
an integer `n` and a list of `edges` where `edges[i] = [a, b]`
indicates that there is an undirected edge between nodes `a` and `b`
in the graph.

Return `true` if the edges of the given graph make up a valid tree,
and `false` otherwise.

A valid tree is a connected graph with no cycles, i.e. exactly
`n - 1` edges connecting all `n` nodes with no redundant connections.

## Example

```
Input: n = 5, edges = [[0,1],[0,2],[0,3],[1,4]]
Output: true

Input: n = 5, edges = [[0,1],[1,2],[2,3],[1,3],[1,4]]
Output: false
```

## Constraints

- `1 <= n <= 2000`
- `0 <= edges.length <= 5000`
- `edges[i].length == 2`
- `0 <= a, b < n`
- `a != b`
- There are no self-loops or repeated edges.
