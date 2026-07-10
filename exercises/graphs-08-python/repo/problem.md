# Course Schedule

There are `numCourses` courses labeled `0` to `numCourses - 1`. You
are given `prerequisites` where `prerequisites[i] = [a, b]` indicates
that you must take course `b` first if you want to take course `a`.

Return `true` if you can finish all courses, otherwise return
`false`.

## Example

```
Input: numCourses = 2, prerequisites = [[1,0]]
Output: true

Input: numCourses = 2, prerequisites = [[1,0],[0,1]]
Output: false
```

## Constraints

- `1 <= numCourses <= 2000`
- `0 <= prerequisites.length <= 5000`
- `prerequisites[i].length == 2`
- `0 <= a, b < numCourses`
- All pairs `[a, b]` are distinct.
