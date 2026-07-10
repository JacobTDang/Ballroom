# Largest Rectangle in Histogram

Given an array of integers `heights` representing the histogram's bar
height where the width of each bar is 1, return the area of the
largest rectangle in the histogram.

## Examples

```
Input: heights = [2,1,5,6,2,3]
Output: 10
Explanation: The largest rectangle has area 10, formed by the bars of
height 5 and 6 (indices 2-3), giving width 2 and height 5: 2*5=10.
```

```
Input: heights = [2,4]
Output: 4
```

## Constraints

- `1 <= heights.length <= 10^5`
- `0 <= heights[i] <= 10^4`
