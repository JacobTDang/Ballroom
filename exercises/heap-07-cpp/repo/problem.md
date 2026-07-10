# Find Median from Data Stream

The median is the middle value in an ordered integer list. If the
size of the list is even, there is no middle value, and the median is
the mean of the two middle values.

Implement the `MedianFinder` class:

- `MedianFinder()` initializes the object.
- `AddNum(num int)` adds the integer `num` to the data structure.
- `FindMedian() float64` returns the median of all elements so far.

## Example

```
Input:
["MedianFinder", "addNum", "addNum", "findMedian", "addNum", "findMedian"]
[[], [1], [2], [], [3], []]

Output:
[null, null, null, 1.5, null, 2.0]

Explanation:
MedianFinder medianFinder = new MedianFinder();
medianFinder.addNum(1);
medianFinder.addNum(2);
medianFinder.findMedian(); // return 1.5 ((1 + 2) / 2)
medianFinder.addNum(3);
medianFinder.findMedian(); // return 2.0
```

## Constraints

- `-10^5 <= num <= 10^5`
- There will be at least one element in the data structure before
  calling `FindMedian`.
- At most `5 * 10^4` calls will be made to `AddNum` and
  `FindMedian`.
