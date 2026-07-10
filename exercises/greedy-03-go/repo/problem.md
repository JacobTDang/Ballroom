# Jump Game II

You are given a 0-indexed array of integers `nums` of length `n`. You
are initially positioned at index `0`.

Each element `nums[i]` represents the maximum length of a forward jump
from index `i`. In other words, if you are at index `i`, you can jump
to any index `i + j` where `0 <= j <= nums[i]` and `i + j < n`.

Return the minimum number of jumps to reach index `n - 1`. It is
guaranteed that you can always reach index `n - 1`.

## Example

```
Input: nums = [2,3,1,1,4]
Output: 2
Explanation: The minimum number of jumps to reach the last index is 2.
Jump 1 step from index 0 to 1, then 3 steps to the last index.
```

## Constraints

- `1 <= nums.length <= 10^4`
- `0 <= nums[i] <= 1000`
- It is guaranteed that you can reach `nums[n - 1]`.
