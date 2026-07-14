# House Robber

You are a professional robber planning to rob houses along a street.
Each house has a certain amount of money stashed, given in the array
`nums`. The only constraint stopping you from robbing each of them is
that adjacent houses have connected security systems, and it will
automatically alert the police if two adjacent houses are broken into
on the same night.

Given `nums`, return the maximum amount of money you can rob tonight
without alerting the police.

## Example

```
Input: nums = [1,2,3,1]
Output: 4
Explanation: Rob house 0 (money = 1) and house 2 (money = 3).
Total amount you can rob = 1 + 3 = 4.
```

## Constraints

- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 1000`
