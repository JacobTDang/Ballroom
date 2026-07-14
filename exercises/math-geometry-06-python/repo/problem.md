# Pow(x, n)

Implement `pow(x, n)`, which calculates `x` raised to the power `n`
(i.e. `x^n`), without using a built-in power operator.

## Example

```
Input: x = 2.00000, n = 10
Output: 1024.00000
```

```
Input: x = 2.10000, n = 3
Output: 9.26100
```

```
Input: x = 2.00000, n = -2
Output: 0.25000
Explanation: 2^-2 = 1/2^2 = 1/4 = 0.25
```

## Constraints

- `-100.0 < x < 100.0`
- `-2^31 <= n <= 2^31 - 1`
- `n` is an integer.
- Either `x` is not zero or `n > 0`.
- Your solution must run in better than `O(n)` time.
