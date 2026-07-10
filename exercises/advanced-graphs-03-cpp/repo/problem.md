# Network Delay Time

You are given a network of `n` nodes, labeled `1` to `n`. You are also
given `times`, a list of directed edges where `times[i] = [ui, vi, wi]`,
representing a signal that travels from node `ui` to node `vi` and takes
`wi` time units to arrive.

A signal is sent from node `k`. Return the minimum time it takes for all
`n` nodes to receive the signal. If it is impossible for all `n` nodes to
receive the signal, return `-1`.

## Example 1

```
Input: times = [[2,1,1],[2,3,1],[3,4,1]], n = 4, k = 2
Output: 2
```

## Example 2

```
Input: times = [[1,2,1]], n = 2, k = 1
Output: 1
```

## Example 3

```
Input: times = [[1,2,1]], n = 2, k = 2
Output: -1
```

## Constraints

- `1 <= k <= n <= 100`
- `1 <= times.length <= 6000`
- `times[i].length == 3`
- `1 <= ui, vi <= n`
- `ui != vi`
- `0 <= wi <= 100`
- All the pairs `(ui, vi)` are unique. (i.e., no multiple edges.)
