# Meeting Rooms II

Given an array of meeting time intervals `intervals` where
`intervals[i] = [start_i, end_i]`, return the minimum number of
conference rooms required so that all meetings can be held without any
two overlapping meetings sharing a room.

Two meetings that only touch at an endpoint, e.g. `[5,10]` and
`[10,15]`, do not require separate rooms.

## Example

```
Input: intervals = [[0,30],[5,10],[15,20]]
Output: 2
```

```
Input: intervals = [[7,10],[2,4]]
Output: 1
```

## Constraints

- `1 <= intervals.length <= 10^4`
- `0 <= start_i < end_i <= 10^6`
