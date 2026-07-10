# Course Schedule II

There are `numCourses` courses labeled `0` to `numCourses - 1`. You
are given `prerequisites` where `prerequisites[i] = [a, b]` indicates
that you must take course `b` first if you want to take course `a`.

Return a valid order in which you could take all the courses. If
there are many valid answers, return any of them. If it is
impossible to finish all courses, return an empty array.

## Example

```
Input: numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]
Output: [0,1,2,3]

Input: numCourses = 2, prerequisites = [[1,0],[0,1]]
Output: []
```

## Constraints

- `1 <= numCourses <= 2000`
- `0 <= prerequisites.length <= 5000`
- `prerequisites[i].length == 2`
- `0 <= a, b < numCourses`
- All pairs `[a, b]` are distinct.
