# Burst Balloons

You are given `n` balloons, indexed from `0` to `n - 1`. Each balloon
is painted with a number on it represented by an array `nums`. You
are asked to burst all the balloons.

If you burst the `i`-th balloon, you will get
`nums[i - 1] * nums[i] * nums[i + 1]` coins. If `i - 1` or `i + 1`
goes out of bounds of the array, then treat it as if there is a
balloon with a value of `1`.

Return the maximum coins you can collect by bursting all the
balloons.

## Example

```
Input: nums = [3,1,5,8]
Output: 167
Explanation:
nums = [3,1,5,8] -> [3,5,8] -> [3,8] -> [8] -> []
coins =  3*1*5    +  3*5*8   +  1*3*8  + 1*8*1 = 167
```

## Constraints

- `1 <= nums.length <= 300`
- `0 <= nums[i] <= 100`
