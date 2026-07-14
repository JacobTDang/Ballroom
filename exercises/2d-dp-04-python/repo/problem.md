# Coin Change II

You are given an integer array `coins` representing coins of
different denominations and an integer `amount` representing a total
amount of money.

Return the number of combinations that make up that amount. If that
amount of money cannot be made up by any combination of the coins,
return `0`.

You may assume that you have an infinite number of each kind of
coin. Two combinations are the same if they use the same coins with
the same multiplicities, regardless of order.

## Example

```
Input: amount = 5, coins = [1,2,5]
Output: 4
Explanation: There are four ways to make up the amount:
5
2+2+1
2+1+1+1
1+1+1+1+1
```

## Constraints

- `1 <= coins.length <= 300`
- `1 <= coins[i] <= 5000`
- All the values of `coins` are unique.
- `0 <= amount <= 5000`
