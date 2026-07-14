# Best Time to Buy and Sell Stock with Cooldown

You are given an array `prices` where `prices[i]` is the price of a
given stock on the `i`-th day.

Find the maximum profit you can achieve. You may complete as many
transactions as you like (buy one and sell one share of the stock
multiple times) with the following restrictions:

- After you sell your stock, you cannot buy stock on the next day
  (i.e., there is a mandatory one-day cooldown).
- You may not engage in multiple transactions simultaneously (you
  must sell the stock before you buy again).

## Example

```
Input: prices = [1,2,3,0,2]
Output: 3
Explanation: buy on day 0 (price = 1), sell on day 1 (price = 2),
cooldown on day 2, buy on day 3 (price = 0), sell on day 4
(price = 2). Total profit = (2 - 1) + (2 - 0) = 3.
```

## Constraints

- `1 <= prices.length <= 5000`
- `0 <= prices[i] <= 1000`
