# Walls and Gates

You are given an `m x n` grid `rooms` initialized with these three
possible values:

- `-1` A wall or an obstacle.
- `0` A gate.
- `2147483647` (`INF`) An empty room, representing "infinity" since
  its distance to the nearest gate is not yet known.

Fill each empty room with the distance to its nearest gate. If it is
impossible to reach a gate, it should remain `INF`. Modify `rooms`
in place; do not return anything.

## Example

```
Input: rooms =
[[INF,-1,0,INF],
 [INF,INF,INF,-1],
 [INF,-1,INF,-1],
 [0,-1,INF,INF]]
Output:
[[3,-1,0,1],
 [2,2,1,-1],
 [1,-1,2,-1],
 [0,-1,3,4]]
```

## Constraints

- `m == rooms.length`
- `n == rooms[i].length`
- `1 <= m, n <= 250`
- `rooms[i][j]` is `-1`, `0`, or `2147483647`.
