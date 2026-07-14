# Coin Change

You are given an integer array `coins` representing coin denominations
and an integer `amount` representing a target amount of money. You
have an unlimited supply of each coin denomination.

Return the fewest number of coins needed to make up `amount`. If that
amount of money cannot be made up by any combination of the coins,
return `-1`.

## Example

```
Input: coins = [1,2,5], amount = 11
Output: 3
Explanation: 11 = 5 + 5 + 1
```

## Constraints

- `1 <= coins.length <= 50`
- `1 <= coins[i] <= 10^4`
- `0 <= amount <= 10^4`
