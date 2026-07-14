# Detect Squares

Design an algorithm that accepts a stream of points on the X-Y plane
and queries points at a later time. Each point is an integer pair
`(x, y)`. The same point may be added more than once, and each
occurrence is counted separately.

Implement `DetectSquares`:

- `Add(point)` adds a new point to the data structure.
- `Count(point)` counts the number of ways to choose three points from
  the previously added points such that, together with `point`, they
  form an axis-aligned square with **positive area**. A square is
  only counted if all four of its corners were previously added (or
  are the query point). Squares formed by different chosen point
  combinations (including using the same physical point added
  multiple times) are all counted separately.

## Example

```
Input:
["DetectSquares", "add", "add", "add", "count", "count", "add", "count"]
[[], [[3, 10]], [[11, 2]], [[3, 2]], [[11, 10]], [[14, 8]], [[11, 2]], [[11, 10]]]

Output:
[null, null, null, null, 1, 0, null, 2]

Explanation:
DetectSquares detectSquares = new DetectSquares();
detectSquares.add([3, 10]);
detectSquares.add([11, 2]);
detectSquares.add([3, 2]);
detectSquares.count([11, 10]); // returns 1: using (3, 10), (11, 2),
                                // and (3, 2), we can form an
                                // axis-aligned square with a
                                // side length of 8.
detectSquares.count([14, 8]);  // returns 0: no axis-aligned square
                                // with all 4 corners can be found.
detectSquares.add([11, 2]);    // now (11, 2) has been added twice.
detectSquares.count([11, 10]); // returns 2: two squares can be
                                // formed, one using the first (11, 2)
                                // and one using the second.
```

## Constraints

- `point.length == 2`
- `0 <= x, y <= 1000`
- At most `3000` calls in total will be made to `Add` and `Count`.
