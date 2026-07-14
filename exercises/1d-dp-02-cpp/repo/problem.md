# Min Cost Climbing Stairs

You are given an integer array `cost` where `cost[i]` is the cost of
`i`th step on a staircase. Once you pay the cost, you can either climb
one or two steps.

You can either start from the step with index `0`, or the step with
index `1`. You "win" the game once you go beyond the last index in
`cost`.

Return the minimum cost to reach the top of the staircase.

## Example

```
Input: cost = [10,15,20]
Output: 15
Explanation: You will start at index 1.
- Pay 15 and climb two steps to reach the top.
The total cost is 15.
```

## Constraints

- `2 <= cost.length <= 1000`
- `0 <= cost[i] <= 999`
