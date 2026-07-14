# Plus One

You are given a large integer represented as an array of digits,
where each `digits[i]` is the `i`th digit of the integer, ordered
from most significant to least significant. The large integer does
not contain any leading zeros, except the number `0` itself.

Increment the large integer by one and return the resulting array of
digits.

## Example

```
Input: digits = [1,2,3]
Output: [1,2,4]
Explanation: The array represents the integer 123. Incrementing by
one gives 124.
```

```
Input: digits = [9,9,9]
Output: [1,0,0,0]
Explanation: The array represents the integer 999. Incrementing by
one gives 1000.
```

## Constraints

- `1 <= digits.length <= 100`
- `0 <= digits[i] <= 9`
- `digits` does not contain any leading zero, except the number `0`
  itself.
