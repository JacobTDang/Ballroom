# The Snapshot That Edits the Original

A `Grid` needs an undo feature: take a snapshot of its current state,
let the caller keep editing the live grid, and later compare against
(or restore) the snapshot. It mostly works, but the "snapshot" isn't
actually independent -- editing the live grid after taking a snapshot
also changes the values the snapshot reports, as if no copy were ever
made.

## Task

Find and fix the bug. You may need to change a return type to do it --
just don't change what the function takes as input, or its name.

## Examples

```
Grid g(2, 2);
g.Set(0, 0, 1);
const Grid& saved = g.Snapshot();
g.Set(0, 0, 999);
saved.Get(0, 0)   // => 1, unaffected by the edit made after the snapshot
g.Get(0, 0)       // => 999, the live grid, which was edited
```

## Constraints

- A grid is `rows` x `cols`, 0-indexed; every cell starts at `0`.
- `Snapshot()` returns an independent `Grid`.
- The snapshot must reflect the grid's state at the moment it was
  taken, and stay that way regardless of later edits to the live grid.
