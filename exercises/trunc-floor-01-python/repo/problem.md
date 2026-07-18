# Negative Buckets

`align` rounds a timestamp down to the start of its `k`-wide bucket --
the kind of thing you'd use to group events (metrics, rate-limit
windows, hourly rollups) into fixed-size time windows. It's correct
for every timestamp at or after the epoch, but silently returns the
wrong bucket for timestamps before it: some of this data is backfilled
from before 1970, and those negative timestamps land in buckets that
don't line up with the grid everything else uses.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
align(7, 4)   => 4
align(8, 4)   => 8
align(0, 4)   => 0
align(-7, 4)  => -8
align(-8, 4)  => -8
```

## Constraints

- `k` is always a positive integer (the bucket width).
- `t` may be any integer: positive, negative, or zero.
- A bucket start is always a multiple of `k`, and it never lands after
  `t` -- the bucket a timestamp falls into starts at or before it.
