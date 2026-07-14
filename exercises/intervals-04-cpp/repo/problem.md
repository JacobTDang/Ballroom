# Meeting Rooms

Given an array of meeting time intervals `intervals` where
`intervals[i] = [start_i, end_i]`, determine if a person could attend
all of the given meetings.

Two meetings that only touch at an endpoint, e.g. `[5,10]` and
`[10,15]`, do **not** count as overlapping.

## Example

```
Input: intervals = [[0,30],[5,10],[15,20]]
Output: false
```

```
Input: intervals = [[7,10],[2,4]]
Output: true
```

## Constraints

- `0 <= intervals.length <= 10^4`
- `intervals[i].length == 2`
- `0 <= start_i < end_i <= 10^6`
