# Find the Duplicate Number

Given an array of integers `nums` containing `n + 1` integers where
each integer is in the range `[1, n]` inclusive, there is only one
repeated number in `nums`. Return this repeated number.

You must solve the problem without modifying the array `nums` and use
only constant extra space.

This can be solved by treating `nums` as a linked list where index
`i` points to index `nums[i]` — the duplicate value creates a cycle,
so Floyd's cycle detection (the same technique as Linked List Cycle)
finds it.

## Examples

```
Input: nums = [1,3,4,2,2]
Output: 2
```

```
Input: nums = [3,1,3,4,2]
Output: 3
```

## Constraints

- `1 <= n <= 10^5`
- `nums.length == n + 1`
- `1 <= nums[i] <= n`
- All the integers in `nums` appear only once except for precisely
  one integer which appears two or more times.
