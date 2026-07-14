# House Robber II

You are a professional robber planning to rob houses along a street.
The houses are arranged in a circle, meaning the first house and the
last house are adjacent to each other. Each house has a certain amount
of money stashed, given in the array `nums`. Adjacent houses have
connected security systems, and it will automatically alert the police
if two adjacent houses are broken into on the same night.

Given `nums`, return the maximum amount of money you can rob tonight
without alerting the police.

## Example

```
Input: nums = [2,3,2]
Output: 3
Explanation: You cannot rob house 0 (money = 2) and house 2 (money = 2)
together since they are adjacent in the circle. The best you can do is
rob house 1 (money = 3).
```

## Constraints

- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 1000`
